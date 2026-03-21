package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/cache"
	"github.com/kyungw00k/upbit/internal/output"
)

var cacheCmd = &cobra.Command{
	Use:     "cache",
	Short:   "캐시 관리",
	GroupID: "util",
	Example: `  upbit cache            # 캐시 정보 (경로, 크기)
  upbit cache --clear    # 캐시 삭제
  upbit cache --clear -f # 캐시 즉시 삭제`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clearFlag, _ := cmd.Flags().GetBool("clear")

		if clearFlag {
			return clearCache(cmd)
		}
		return showCacheInfo()
	},
}

var cacheInfoColumns = []output.TableColumn{
	{Header: "경로", Key: "path"},
	{Header: "파일", Key: "files"},
	{Header: "크기", Key: "size"},
}

type cacheInfo struct {
	Path  string `json:"path"`
	Files string `json:"files"`
	Size  string `json:"size"`
}

func showCacheInfo() error {
	dir, err := cache.CandleCacheDir()
	if err != nil {
		return fmt.Errorf("캐시 경로 확인 실패: %w", err)
	}

	var totalSize int64
	var fileCount int
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		totalSize += info.Size()
		fileCount++
		return nil
	})

	info := cacheInfo{
		Path:  dir,
		Files: fmt.Sprintf("%d개", fileCount),
		Size:  humanSize(totalSize),
	}

	formatter := GetFormatterWithColumns(cacheInfoColumns)
	return formatter.Format(info)
}

func clearCache(cmd *cobra.Command) error {
	force := GetForce()

	ok, err := output.Confirm("캐시를 삭제하시겠습니까?", force)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Fprintln(os.Stderr, "취소되었습니다.")
		return nil
	}

	cc, err := cache.NewCandleCache()
	if err != nil {
		return fmt.Errorf("캐시 열기 실패: %w", err)
	}
	defer cc.Close()

	if err := cc.Clear(); err != nil {
		return fmt.Errorf("캐시 삭제 실패: %w", err)
	}

	dir, _ := cache.CandleCacheDir()
	if dir != "" {
		_ = os.Remove(filepath.Join(dir, "markets.json"))
	}

	fmt.Fprintln(os.Stderr, "캐시가 삭제되었습니다.")
	return nil
}

func humanSize(b int64) string {
	switch {
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func init() {
	cacheCmd.Flags().Bool("clear", false, "캐시 삭제")
	AddForceFlag(cacheCmd)
	rootCmd.AddCommand(cacheCmd)
}
