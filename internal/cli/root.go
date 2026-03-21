package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kyungw00k/upbit/internal/api"
	"github.com/kyungw00k/upbit/internal/config"
	"github.com/kyungw00k/upbit/internal/output"
)

// Version Makefile LDFLAGS로 주입됨
var Version = "dev"

// 전역 플래그 값을 저장하는 변수
var (
	flagOutput     string
	flagJSONFields string
	flagForce      bool
)

// rootCmd 루트 Cobra 명령
var rootCmd = &cobra.Command{
	Use:   "upbit",
	Short: "Upbit 거래소 CLI",
	Long: `Upbit 거래소 CLI — 시세 조회, 거래, 입출금, 실시간 구독을 지원합니다.`,
	Version: Version,
	// 인자 없이 실행 시 도움말 표시
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	// cobra 에러 메시지 중복 방지 (에러는 Execute() 호출자가 직접 출력)
	rootCmd.SilenceErrors = true
	// RunE 에러 시 도움말(Usage) 출력 방지
	rootCmd.SilenceUsage = true

	// 카테고리 그룹 등록
	rootCmd.AddGroup(
		&cobra.Group{ID: "quotation", Title: "시세 명령어:"},
		&cobra.Group{ID: "trading", Title: "거래 명령어:"},
		&cobra.Group{ID: "wallet", Title: "입출금 명령어:"},
		&cobra.Group{ID: "realtime", Title: "실시간 명령어:"},
		&cobra.Group{ID: "util", Title: "유틸 명령어:"},
	)

	// 글로벌 플래그 등록
	pf := rootCmd.PersistentFlags()

	pf.StringVarP(&flagOutput, "output", "o", "auto",
		"출력 포맷: table, json, jsonl, csv, auto")

	pf.StringVar(&flagJSONFields, "json", "",
		"JSON 출력 필드 목록 (쉼표 구분, 예: market,trade_price)")

	// --force는 확인 프롬프트가 있는 명령에만 로컬로 등록 (글로벌 X)

	// --version 플래그는 cobra가 Version 설정 시 자동으로 추가함

	// 오타 시 "Did you mean?" 제안 활성화
	rootCmd.SuggestionsMinimumDistance = 2
}

// Execute 루트 명령 실행
func Execute() error {
	return rootCmd.Execute()
}

// RootCmd 루트 명령 반환 (서브커맨드 등록용)
func RootCmd() *cobra.Command {
	return rootCmd
}

// GetClient config에서 키를 로드하여 API 클라이언트 반환
// 설정 로드 실패 시 빈 키로 클라이언트를 반환 (Quotation API 등 인증 불필요 경로 허용)
func GetClient() *api.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "설정 로드 실패:", err)
		return api.NewClient("", "")
	}
	return api.NewClient(cfg.AccessKey, cfg.SecretKey)
}

// GetClientE config에서 키를 로드하여 API 클라이언트와 에러를 반환
// requireAuth가 true이면 키 누락 시 조기 에러 반환
func GetClientE(requireAuth bool) (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		if requireAuth {
			return nil, fmt.Errorf("설정 로드 실패: %w", err)
		}
		return api.NewClient("", ""), nil
	}
	if requireAuth && (cfg.AccessKey == "" || cfg.SecretKey == "") {
		return nil, fmt.Errorf("인증이 필요합니다: ACCESS_KEY 및 SECRET_KEY를 설정하세요")
	}
	return api.NewClient(cfg.AccessKey, cfg.SecretKey), nil
}

// GetFormatter 현재 플래그 값으로 출력 포맷터 반환
func GetFormatter() output.Formatter {
	return output.NewFormatter(flagOutput, flagJSONFields)
}

// GetFormatterWithColumns 현재 플래그 값과 컬럼 정의로 출력 포맷터 반환
// 테이블 출력 시 지정된 컬럼만 표시, JSON/CSV는 전체 데이터 유지
func GetFormatterWithColumns(columns []output.TableColumn) output.Formatter {
	return output.NewFormatterWithColumns(flagOutput, flagJSONFields, columns)
}

// GetForce cmd에서 --force 플래그 값 반환 (로컬 플래그)
func GetForce() bool {
	return flagForce
}

// AddForceFlag 명령에 --force 로컬 플래그를 등록하는 헬퍼
func AddForceFlag(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&flagForce, "force", "f", false, "확인 프롬프트 스킵")
}
