package cli

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/i18n"
)

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

	// 2. 버전 비교
	if latestClean == currentClean {
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

			if _, err := io.Copy(tmpFile, tr); err != nil {
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
