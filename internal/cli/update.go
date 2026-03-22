package cli

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/i18n"
)

// maxExtractSize tar 엔트리 최대 크기 (100MB) — zip bomb 방지
const maxExtractSize = 100 * 1024 * 1024

// GitHub Release API 응답 구조체
type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

const githubReleaseURL = "https://api.github.com/repos/kyungw00k/upbit/releases/latest"

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   i18n.T(i18n.MsgUpdateShort),
	GroupID: "util",
	Example: `  upbit update          # 최신 버전으로 업데이트
  upbit update --check  # 업데이트 확인만`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().BoolP("check", "c", false, i18n.T(i18n.FlagUpdateCheckUsage))
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check")

	fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgUpdateChecking))

	// 1. GitHub API에서 최신 릴리스 조회
	release, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrUpdateFetch), err)
	}

	latestVersion := release.TagName
	currentVersion := Version

	// 버전 비교를 위해 "v" 접두사 정규화
	latestClean := strings.TrimPrefix(latestVersion, "v")
	currentClean := strings.TrimPrefix(currentVersion, "v")

	fmt.Fprintf(os.Stderr, "%s\n", i18n.Tf(i18n.MsgUpdateLatest, latestVersion))

	// 2. 시맨틱 버전 비교 (Q-4)
	if compareVersions(currentClean, latestClean) >= 0 {
		fmt.Fprintln(os.Stderr, i18n.Tf(i18n.MsgUpdateAlreadyLatest, currentVersion))
		return nil
	}

	fmt.Fprintln(os.Stderr, i18n.Tf(i18n.MsgUpdateAvailable, currentVersion, latestVersion))

	if checkOnly {
		return nil
	}

	// 3. OS/Arch에 맞는 asset 찾기
	asset := findAsset(release.Assets)
	if asset == nil {
		return fmt.Errorf("%s", i18n.Tf(i18n.ErrUpdateNoAsset, runtime.GOOS, runtime.GOARCH))
	}

	// 4. 다운로드
	fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgUpdateDownloading))

	tmpFile, err := downloadAsset(asset)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrUpdateDownload), err)
	}
	defer os.Remove(tmpFile)

	// 4-1. 체크섬 검증 (S-2)
	fmt.Fprintln(os.Stderr, i18n.T(i18n.MsgUpdateVerifying))
	if checksumAsset := findChecksumAsset(release.Assets); checksumAsset != nil {
		expectedHash, err := fetchExpectedChecksum(checksumAsset, asset.Name)
		if err == nil && expectedHash != "" {
			if err := verifyChecksum(tmpFile, expectedHash); err != nil {
				os.Remove(tmpFile)
				return fmt.Errorf("%s: %w", i18n.T(i18n.ErrUpdateChecksum), err)
			}
		}
	}

	// 5. tar.gz인 경우 압축 해제
	binaryPath := tmpFile
	if strings.HasSuffix(asset.Name, ".tar.gz") || strings.HasSuffix(asset.Name, ".tgz") {
		extracted, err := extractTarGz(tmpFile)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T(i18n.ErrUpdateDownload), err)
		}
		defer os.Remove(extracted)
		binaryPath = extracted
	}

	// 6. 바이너리 교체
	if err := replaceBinary(binaryPath); err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.ErrUpdateReplace), err)
	}

	fmt.Fprintln(os.Stderr, i18n.Tf(i18n.MsgUpdateComplete, currentVersion, latestVersion))
	return nil
}

// fetchLatestRelease GitHub API에서 최신 릴리스 정보 조회
func fetchLatestRelease() (*githubRelease, error) {
	req, err := http.NewRequest("GET", githubReleaseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "upbit-cli/"+Version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

// findAsset OS/Arch에 맞는 asset 탐색
func findAsset(assets []githubAsset) *githubAsset {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// goreleaser 기본 네이밍 패턴들
	patterns := []string{
		fmt.Sprintf("upbit_%s_%s.tar.gz", goos, goarch),
		fmt.Sprintf("upbit_%s_%s.zip", goos, goarch),
		fmt.Sprintf("upbit_%s_%s", goos, goarch),
		fmt.Sprintf("upbit-%s-%s.tar.gz", goos, goarch),
		fmt.Sprintf("upbit-%s-%s", goos, goarch),
	}

	for _, pattern := range patterns {
		for i, a := range assets {
			if strings.EqualFold(a.Name, pattern) {
				return &assets[i]
			}
		}
	}
	return nil
}

// downloadAsset asset을 임시 파일에 다운로드
func downloadAsset(asset *githubAsset) (string, error) {
	resp, err := http.Get(asset.BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "upbit-update-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

// extractTarGz tar.gz에서 upbit 바이너리를 추출
func extractTarGz(tarGzPath string) (string, error) {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// 바이너리 파일 찾기: "upbit" 이름이거나 실행 가능한 파일
		name := filepath.Base(header.Name)
		if name == "upbit" || name == "upbit.exe" {
			tmpFile, err := os.CreateTemp("", "upbit-bin-*")
			if err != nil {
				return "", err
			}
			defer tmpFile.Close()

			// S-3: io.LimitReader로 100MB 제한 — zip bomb / decompression bomb 방지
			limited := io.LimitReader(tr, maxExtractSize)
			if _, err := io.Copy(tmpFile, limited); err != nil {
				os.Remove(tmpFile.Name())
				return "", err
			}

			return tmpFile.Name(), nil
		}
	}

	return "", fmt.Errorf("binary not found in archive")
}

// replaceBinary 현재 실행 파일을 새 바이너리로 교체 (rollback 지원)
func replaceBinary(newBinary string) error {
	currentPath, err := os.Executable()
	if err != nil {
		return err
	}
	currentPath, err = filepath.EvalSymlinks(currentPath)
	if err != nil {
		return err
	}

	oldPath := currentPath + ".old"

	// 기존 → .old
	if err := os.Rename(currentPath, oldPath); err != nil {
		return err
	}

	// 새 바이너리 → 원래 경로
	if err := copyFile(newBinary, currentPath, 0755); err != nil {
		// 복원
		_ = os.Rename(oldPath, currentPath)
		return err
	}

	// .old 삭제
	_ = os.Remove(oldPath)
	return nil
}

// compareVersions 시맨틱 버전 비교 (Q-4)
// a > b: 1, a == b: 0, a < b: -1
func compareVersions(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		var va, vb int
		if i < len(partsA) {
			va, _ = strconv.Atoi(partsA[i])
		}
		if i < len(partsB) {
			vb, _ = strconv.Atoi(partsB[i])
		}
		if va > vb {
			return 1
		}
		if va < vb {
			return -1
		}
	}
	return 0
}

// findChecksumAsset release assets에서 checksums.txt 파일을 탐색
func findChecksumAsset(assets []githubAsset) *githubAsset {
	for i, a := range assets {
		name := strings.ToLower(a.Name)
		if name == "checksums.txt" || name == "sha256sums.txt" {
			return &assets[i]
		}
	}
	return nil
}

// fetchExpectedChecksum checksums.txt를 다운로드하고 대상 파일의 해시를 반환
func fetchExpectedChecksum(checksumAsset *githubAsset, targetFilename string) (string, error) {
	resp, err := http.Get(checksumAsset.BrowserDownloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checksum download returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB 제한
	if err != nil {
		return "", err
	}

	// goreleaser 기본 형식: "abc123def456...  upbit_darwin_arm64.tar.gz"
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == targetFilename {
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("checksum for %s not found", targetFilename)
}

// verifyChecksum 파일의 SHA256 해시를 검증
func verifyChecksum(filePath string, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	actual := hex.EncodeToString(h.Sum(nil))
	if actual != expectedHash {
		return fmt.Errorf("expected %s, got %s", expectedHash, actual)
	}
	return nil
}

// copyFile src를 dst로 복사하고 권한 설정
func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}
