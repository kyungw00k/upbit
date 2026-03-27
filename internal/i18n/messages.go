package i18n

// ── Root / Global ──

const (
	MsgRootShort Key = "MsgRootShort"
	MsgRootLong  Key = "MsgRootLong"

	// Group titles
	GroupQuotation Key = "GroupQuotation"
	GroupTrading   Key = "GroupTrading"
	GroupWallet    Key = "GroupWallet"
	GroupRealtime  Key = "GroupRealtime"
	GroupUtil      Key = "GroupUtil"

	// Global flags
	FlagOutputUsage     Key = "FlagOutputUsage"
	FlagJSONFieldsUsage Key = "FlagJSONFieldsUsage"
	FlagForceUsage      Key = "FlagForceUsage"

	// Root errors
	ErrConfigLoad  Key = "ErrConfigLoad"
	ErrAuthRequired Key = "ErrAuthRequired"

	// Args helper
	MsgUsagePrefix Key = "MsgUsagePrefix"

	// ── Balance ──

	HdrCurrency   Key = "HdrCurrency"
	HdrBalance    Key = "HdrBalance"
	HdrLocked     Key = "HdrLocked"
	HdrAvgBuyPrice Key = "HdrAvgBuyPrice"
	HdrEvalKRW    Key = "HdrEvalKRW"

	MsgBalanceShort      Key = "MsgBalanceShort"
	ErrBalanceNotFound   Key = "ErrBalanceNotFound"
	MsgBalanceEmpty      Key = "MsgBalanceEmpty"

	// ── Buy ──

	MsgBuyShort             Key = "MsgBuyShort"
	ErrBuyArgsRequired      Key = "ErrBuyArgsRequired"
	ErrBuyParamsRequired    Key = "ErrBuyParamsRequired"
	MsgTickCheckFailed      Key = "MsgTickCheckFailed"
	MsgBuyOrderAdjusted     Key = "MsgBuyOrderAdjusted"
	MsgBuyOrderNormal       Key = "MsgBuyOrderNormal"
	MsgTickAdjusted         Key = "MsgTickAdjusted"
	MsgOrderCancelled       Key = "MsgOrderCancelled"
	MsgDescLimitOrder       Key = "MsgDescLimitOrder"
	MsgDescPriceOrder       Key = "MsgDescPriceOrder"
	FlagPriceUsage          Key = "FlagPriceUsage"
	FlagVolumeUsage         Key = "FlagVolumeUsage"
	FlagTotalUsage          Key = "FlagTotalUsage"
	FlagSMPUsage            Key = "FlagSMPUsage"
	FlagIdentifierUsage     Key = "FlagIdentifierUsage"
	FlagTestUsage           Key = "FlagTestUsage"
	FlagBestUsage           Key = "FlagBestUsage"
	FlagWatchPriceUsage     Key = "FlagWatchPriceUsage"
	MsgDescBestOrder        Key = "MsgDescBestOrder"
	MsgDescReservedOrder    Key = "MsgDescReservedOrder"

	// ── Percent ──

	ErrTickerFetch         Key = "ErrTickerFetch"
	ErrTickerEmpty         Key = "ErrTickerEmpty"
	MsgPriceKeywordResolved Key = "MsgPriceKeywordResolved"

	ErrPercentParse        Key = "ErrPercentParse"
	ErrPercentRange        Key = "ErrPercentRange"
	ErrPercentBuyNeedsPrice Key = "ErrPercentBuyNeedsPrice"
	ErrChanceFetch         Key = "ErrChanceFetch"
	MsgPercentResolved     Key = "MsgPercentResolved"

	// ── Sell ──

	MsgSellShort          Key = "MsgSellShort"
	ErrSellArgsRequired   Key = "ErrSellArgsRequired"
	ErrSellParamsRequired Key = "ErrSellParamsRequired"
	MsgSellOrderAdjusted  Key = "MsgSellOrderAdjusted"
	MsgSellOrderNormal    Key = "MsgSellOrderNormal"
	MsgDescMarketSell     Key = "MsgDescMarketSell"

	// ── Market ──

	HdrMarket      Key = "HdrMarket"
	HdrKoreanName  Key = "HdrKoreanName"
	HdrEnglishName Key = "HdrEnglishName"

	MsgMarketShort       Key = "MsgMarketShort"
	MsgMarketFilterEmpty Key = "MsgMarketFilterEmpty"
	FlagQuoteUsage       Key = "FlagQuoteUsage"

	// ── Ticker ──

	HdrPrice        Key = "HdrPrice"
	HdrChange       Key = "HdrChange"
	HdrChangeRate   Key = "HdrChangeRate"
	HdrVolume24h    Key = "HdrVolume24h"
	HdrTradePrice24h Key = "HdrTradePrice24h"
	HdrHigh         Key = "HdrHigh"
	HdrLow          Key = "HdrLow"

	MsgTickerShort        Key = "MsgTickerShort"
	ErrTickerNoMarket     Key = "ErrTickerNoMarket"
	FlagTickerQuoteUsage  Key = "FlagTickerQuoteUsage"

	// ── API Keys ──

	HdrAccessKey Key = "HdrAccessKey"
	HdrExpireAt  Key = "HdrExpireAt"

	MsgApiKeysShort    Key = "MsgApiKeysShort"
	MsgApiKeysEmpty    Key = "MsgApiKeysEmpty"
	FlagApiKeysAllUsage Key = "FlagApiKeysAllUsage"

	// ── Candle ──

	HdrCandleTime   Key = "HdrCandleTime"
	HdrOpen         Key = "HdrOpen"
	HdrClose        Key = "HdrClose"
	HdrCandleVolume Key = "HdrCandleVolume"

	MsgCandleShort         Key = "MsgCandleShort"
	ErrCandleMarketRequired Key = "ErrCandleMarketRequired"
	FlagIntervalUsage      Key = "FlagIntervalUsage"
	FlagCountUsage         Key = "FlagCountUsage"
	FlagFromUsage          Key = "FlagFromUsage"
	FlagAscUsage           Key = "FlagAscUsage"
	FlagDescUsage          Key = "FlagDescUsage"
	FlagNoCacheUsage       Key = "FlagNoCacheUsage"
	MsgCandleConfirm       Key = "MsgCandleConfirm"
	MsgCancelled           Key = "MsgCancelled"
	ErrFromParse           Key = "ErrFromParse"
	ErrUnrecognizedTime    Key = "ErrUnrecognizedTime"

	// Cache (candle.go internal messages)
	MsgCacheInitFailed      Key = "MsgCacheInitFailed"
	MsgCacheRangeError      Key = "MsgCacheRangeError"
	MsgCacheSaveError       Key = "MsgCacheSaveError"
	MsgCacheUpdateError     Key = "MsgCacheUpdateError"
	ErrOlderRangeFetch      Key = "ErrOlderRangeFetch"
	ErrNewerRangeFetch      Key = "ErrNewerRangeFetch"
	ErrCacheQuery           Key = "ErrCacheQuery"

	// ── Orderbook ──

	MsgOrderbookShort     Key = "MsgOrderbookShort"
	ErrOrderbookMarket    Key = "ErrOrderbookMarket"
	MsgOrderbookMarket    Key = "MsgOrderbookMarket"
	MsgOrderbookTotalSizes Key = "MsgOrderbookTotalSizes"
	HdrAskSize            Key = "HdrAskSize"
	HdrAskPrice           Key = "HdrAskPrice"
	HdrBidPrice           Key = "HdrBidPrice"
	HdrBidSize            Key = "HdrBidSize"

	// ── Orderbook Levels ──

	HdrSupportedLevels       Key = "HdrSupportedLevels"
	MsgOrderbookLevelsShort  Key = "MsgOrderbookLevelsShort"

	// ── Trades ──

	HdrTradeTime   Key = "HdrTradeTime"
	HdrTradePrice  Key = "HdrTradePrice"
	HdrTradeVolume Key = "HdrTradeVolume"
	HdrAskBid      Key = "HdrAskBid"
	HdrChangePrice Key = "HdrChangePrice"

	MsgTradesShort  Key = "MsgTradesShort"

	// ── Order Root ──

	MsgOrderShort Key = "MsgOrderShort"

	// ── Order List ──

	HdrSide         Key = "HdrSide"
	HdrOrdType      Key = "HdrOrdType"
	HdrOrderPrice   Key = "HdrOrderPrice"
	HdrOrderVolume  Key = "HdrOrderVolume"
	HdrExecutedVol  Key = "HdrExecutedVol"
	HdrState        Key = "HdrState"
	HdrCreatedAt    Key = "HdrCreatedAt"

	MsgOrderListShort     Key = "MsgOrderListShort"
	MsgOrderListClosed    Key = "MsgOrderListClosed"
	MsgOrderListOpen      Key = "MsgOrderListOpen"
	FlagClosedUsage       Key = "FlagClosedUsage"
	FlagMarketFilterUsage Key = "FlagMarketFilterUsage"
	FlagPageUsage         Key = "FlagPageUsage"

	// ── Order Show ──

	HdrRemainingVol    Key = "HdrRemainingVol"
	MsgOrderShowShort  Key = "MsgOrderShowShort"
	ErrOrderShowArgs   Key = "ErrOrderShowArgs"
	FlagOrderShowIDUsage Key = "FlagOrderShowIDUsage"

	// ── Order Cancel ──

	MsgOrderCancelShort     Key = "MsgOrderCancelShort"
	MsgCancelAllOrders      Key = "MsgCancelAllOrders"
	MsgCancelAllMarket      Key = "MsgCancelAllMarket"
	MsgCancelAborted        Key = "MsgCancelAborted"
	MsgCancelSuccess        Key = "MsgCancelSuccess"
	MsgCancelFailed         Key = "MsgCancelFailed"
	ErrCancelNoUUID         Key = "ErrCancelNoUUID"
	MsgCancelSingleOrder    Key = "MsgCancelSingleOrder"
	FlagCancelAllUsage      Key = "FlagCancelAllUsage"
	FlagCancelMarketUsage   Key = "FlagCancelMarketUsage"

	// ── Order Replace ──

	MsgOrderReplaceShort    Key = "MsgOrderReplaceShort"
	ErrReplaceArgs          Key = "ErrReplaceArgs"
	ErrReplaceNoParam       Key = "ErrReplaceNoParam"
	ErrReplaceLimitNoPrice  Key = "ErrReplaceLimitNoPrice"
	MsgReplaceConfirm       Key = "MsgReplaceConfirm"
	MsgReplacePrice         Key = "MsgReplacePrice"
	MsgReplaceVolume        Key = "MsgReplaceVolume"
	MsgReplaceRemain        Key = "MsgReplaceRemain"
	MsgReplaceCancelled     Key = "MsgReplaceCancelled"
	FlagNewPriceUsage       Key = "FlagNewPriceUsage"
	FlagNewVolumeUsage      Key = "FlagNewVolumeUsage"
	FlagOrdTypeUsage        Key = "FlagOrdTypeUsage"

	// ── Order Chance ──

	MsgOrderChanceShort Key = "MsgOrderChanceShort"
	LblMarket           Key = "LblMarket"
	LblState            Key = "LblState"
	LblBidFee           Key = "LblBidFee"
	LblAskFee           Key = "LblAskFee"
	LblBidAvailable     Key = "LblBidAvailable"
	LblAskAvailable     Key = "LblAskAvailable"
	LblMinOrder         Key = "LblMinOrder"
	LblBidTypes         Key = "LblBidTypes"
	LblAskTypes         Key = "LblAskTypes"

	// ── Watch ──

	MsgWatchShort          Key = "MsgWatchShort"
	MsgWatchTickerShort    Key = "MsgWatchTickerShort"
	MsgWatchOrderbookShort Key = "MsgWatchOrderbookShort"
	MsgWatchTradeShort     Key = "MsgWatchTradeShort"
	MsgWatchCandleShort    Key = "MsgWatchCandleShort"
	MsgWatchMyOrderShort   Key = "MsgWatchMyOrderShort"
	MsgWatchMyAssetShort   Key = "MsgWatchMyAssetShort"
	FlagWatchIntervalUsage Key = "FlagWatchIntervalUsage"

	ErrMarketNotFound     Key = "ErrMarketNotFound"
	ErrSubscribeBuild     Key = "ErrSubscribeBuild"
	ErrSubscribeSend      Key = "ErrSubscribeSend"
	ErrMessageReceive     Key = "ErrMessageReceive"
	ErrServerError        Key = "ErrServerError"

	// Watch formatters
	WatchVolume    Key = "WatchVolume"
	WatchSell      Key = "WatchSell"
	WatchBuy       Key = "WatchBuy"
	WatchSpread    Key = "WatchSpread"
	WatchQty       Key = "WatchQty"
	WatchOpen      Key = "WatchOpen"
	WatchHigh      Key = "WatchHigh"
	WatchLow       Key = "WatchLow"
	WatchClose     Key = "WatchClose"
	WatchState     Key = "WatchState"
	WatchLocking   Key = "WatchLocking"

	// ── Wallet / Status ──

	HdrWalletState   Key = "HdrWalletState"
	HdrBlockState    Key = "HdrBlockState"
	HdrNetwork       Key = "HdrNetwork"
	HdrNetworkType   Key = "HdrNetworkType"
	MsgWalletShort   Key = "MsgWalletShort"

	// ── Tick Size ──

	HdrQuoteCurrency  Key = "HdrQuoteCurrency"
	HdrTickSize       Key = "HdrTickSize"
	MsgTickSizeShort  Key = "MsgTickSizeShort"

	// ── Cache command ──

	HdrPath             Key = "HdrPath"
	HdrFiles            Key = "HdrFiles"
	HdrSize             Key = "HdrSize"
	MsgCacheShort       Key = "MsgCacheShort"
	MsgCacheFileCount   Key = "MsgCacheFileCount"
	ErrCachePathFailed  Key = "ErrCachePathFailed"
	MsgCacheConfirmClear Key = "MsgCacheConfirmClear"
	ErrCacheOpenFailed  Key = "ErrCacheOpenFailed"
	ErrCacheClearFailed Key = "ErrCacheClearFailed"
	MsgCacheCleared     Key = "MsgCacheCleared"
	FlagCacheClearUsage Key = "FlagCacheClearUsage"

	// ── Deposit ──

	MsgDepositShort         Key = "MsgDepositShort"
	MsgDepositListShort     Key = "MsgDepositListShort"
	MsgDepositShowShort     Key = "MsgDepositShowShort"
	MsgDepositAddressShort  Key = "MsgDepositAddressShort"
	MsgDepositCreateShort   Key = "MsgDepositCreateShort"
	MsgDepositListEmpty     Key = "MsgDepositListEmpty"
	MsgDepositAddressEmpty  Key = "MsgDepositAddressEmpty"
	MsgDepositAddrCreated   Key = "MsgDepositAddrCreated"
	MsgDepositAddrCheck     Key = "MsgDepositAddrCheck"
	ErrDepositUUIDRequired  Key = "ErrDepositUUIDRequired"
	ErrDepositCurrRequired  Key = "ErrDepositCurrRequired"

	HdrAmount        Key = "HdrAmount"
	HdrFee           Key = "HdrFee"
	HdrTransType     Key = "HdrTransType"
	HdrTXID          Key = "HdrTXID"
	HdrDoneAt        Key = "HdrDoneAt"
	HdrDepositAddr   Key = "HdrDepositAddr"
	HdrTag           Key = "HdrTag"

	FlagCurrencyUsage     Key = "FlagCurrencyUsage"
	FlagStateUsage        Key = "FlagStateUsage"
	FlagNetTypeUsage      Key = "FlagNetTypeUsage"

	// ── Withdraw ──

	MsgWithdrawShort         Key = "MsgWithdrawShort"
	MsgWithdrawListShort     Key = "MsgWithdrawListShort"
	MsgWithdrawShowShort     Key = "MsgWithdrawShowShort"
	MsgWithdrawRequestShort  Key = "MsgWithdrawRequestShort"
	MsgWithdrawCancelShort   Key = "MsgWithdrawCancelShort"
	MsgWithdrawListEmpty     Key = "MsgWithdrawListEmpty"
	MsgWithdrawCancelled     Key = "MsgWithdrawCancelled"
	ErrWithdrawUUIDRequired  Key = "ErrWithdrawUUIDRequired"
	ErrWithdrawCurrRequired  Key = "ErrWithdrawCurrRequired"
	ErrAmountRequired        Key = "ErrAmountRequired"
	ErrTwoFactorRequired     Key = "ErrTwoFactorRequired"
	MsgWithdrawConfirmKRW    Key = "MsgWithdrawConfirmKRW"
	ErrWithdrawAddrRequired  Key = "ErrWithdrawAddrRequired"
	MsgWithdrawConfirmCoin   Key = "MsgWithdrawConfirmCoin"
	MsgWithdrawCancelConfirm Key = "MsgWithdrawCancelConfirm"

	FlagAmountUsage        Key = "FlagAmountUsage"
	FlagToUsage            Key = "FlagToUsage"
	FlagSecondaryAddrUsage Key = "FlagSecondaryAddrUsage"
	FlagTxTypeUsage        Key = "FlagTxTypeUsage"
	FlagTwoFactorUsage     Key = "FlagTwoFactorUsage"
	FlagWithdrawStateUsage Key = "FlagWithdrawStateUsage"

	// ── Travelrule ──

	MsgTravelruleShort       Key = "MsgTravelruleShort"
	MsgTravelruleVaspsShort  Key = "MsgTravelruleVaspsShort"
	MsgTravelruleTxIDShort   Key = "MsgTravelruleTxIDShort"
	MsgTravelruleUUIDShort   Key = "MsgTravelruleUUIDShort"
	ErrTravelruleTxIDArgs    Key = "ErrTravelruleTxIDArgs"
	ErrTravelruleUUIDArgs    Key = "ErrTravelruleUUIDArgs"

	HdrVaspName      Key = "HdrVaspName"
	HdrVaspNameEn    Key = "HdrVaspNameEn"
	HdrDepositable   Key = "HdrDepositable"
	HdrWithdrawable  Key = "HdrWithdrawable"
	HdrDepositUUID   Key = "HdrDepositUUID"
	HdrVerifyResult  Key = "HdrVerifyResult"
	HdrDepositState  Key = "HdrDepositState"

	FlagVaspUsage         Key = "FlagVaspUsage"
	FlagTRCurrencyUsage   Key = "FlagTRCurrencyUsage"
	FlagTRNetTypeUsage    Key = "FlagTRNetTypeUsage"

	// ── Tool Schema ──

	MsgToolSchemaShort   Key = "MsgToolSchemaShort"
	ErrToolSchemaNotFound Key = "ErrToolSchemaNotFound"

	// Tool schema arg descriptions
	ArgMarketCode     Key = "ArgMarketCode"
	ArgOrderbookMarket Key = "ArgOrderbookMarket"
	ArgTradesMarket   Key = "ArgTradesMarket"
	ArgCandleMarket   Key = "ArgCandleMarket"
	ArgBuyMarket      Key = "ArgBuyMarket"
	ArgSellMarket     Key = "ArgSellMarket"
	ArgBalanceCurrency Key = "ArgBalanceCurrency"
	ArgShowUUID       Key = "ArgShowUUID"
	ArgCancelUUID     Key = "ArgCancelUUID"
	ArgReplaceUUID    Key = "ArgReplaceUUID"
	ArgRequestCurrency Key = "ArgRequestCurrency"
	ArgAddressCurrency Key = "ArgAddressCurrency"
	ArgChanceMarket   Key = "ArgChanceMarket"

	// ── Price Adjust ──

	ErrTickSizeFetch  Key = "ErrTickSizeFetch"
	ErrTickSizeEmpty  Key = "ErrTickSizeEmpty"
	ErrTickSizeParse  Key = "ErrTickSizeParse"
	ErrTickSizeZero   Key = "ErrTickSizeZero"
	ErrPriceParse     Key = "ErrPriceParse"
	ErrUnknownSide    Key = "ErrUnknownSide"

	// ── Update ──

	MsgUpdateShort        Key = "MsgUpdateShort"
	MsgUpdateChecking     Key = "MsgUpdateChecking"
	MsgUpdateLatest       Key = "MsgUpdateLatest"
	MsgUpdateAvailable    Key = "MsgUpdateAvailable"
	MsgUpdateDownloading  Key = "MsgUpdateDownloading"
	MsgUpdateComplete     Key = "MsgUpdateComplete"
	MsgUpdateAlreadyLatest Key = "MsgUpdateAlreadyLatest"
	ErrUpdateFetch        Key = "ErrUpdateFetch"
	ErrUpdateNoAsset      Key = "ErrUpdateNoAsset"
	ErrUpdateDownload     Key = "ErrUpdateDownload"
	ErrUpdateReplace      Key = "ErrUpdateReplace"
	ErrUpdateChecksum     Key = "ErrUpdateChecksum"
	MsgUpdateVerifying    Key = "MsgUpdateVerifying"
	FlagUpdateCheckUsage  Key = "FlagUpdateCheckUsage"

	// ── TUI ──

	TUIQuitHint       Key = "TUIQuitHint"
	TUITabHint        Key = "TUITabHint"
	TUITickerHeader   Key = "TUITickerHeader"
	TUIOrderbookTitle Key = "TUIOrderbookTitle"
	TUITotalAsk       Key = "TUITotalAsk"
	TUITotalBid       Key = "TUITotalBid"
	TUITradeTitle     Key = "TUITradeTitle"
	TUICandleTitle    Key = "TUICandleTitle"
	TUIHighest        Key = "TUIHighest"
	TUILowest         Key = "TUILowest"
	TUIVar            Key = "TUIVar"
	TUIAvg            Key = "TUIAvg"
	TUICumVol         Key = "TUICumVol"
	TUIConfirmYes     Key = "TUIConfirmYes"
	TUIConfirmNo      Key = "TUIConfirmNo"

	// ── Confirm ──

	MsgConfirmNonTTY Key = "MsgConfirmNonTTY"
	ErrInputRead     Key = "ErrInputRead"
)

var ko = map[Key]string{
	// Root
	MsgRootShort: "Upbit 거래소 CLI",
	MsgRootLong:  "Upbit 거래소 CLI — 시세 조회, 거래, 입출금, 실시간 구독을 지원합니다.\n\n  https://github.com/kyungw00k/upbit\n  Sponsor: https://github.com/sponsors/kyungw00k",

	// Groups
	GroupQuotation: "시세 명령어:",
	GroupTrading:   "거래 명령어:",
	GroupWallet:    "입출금 명령어:",
	GroupRealtime:  "실시간 명령어:",
	GroupUtil:      "유틸 명령어:",

	// Global flags
	FlagOutputUsage:     "출력 포맷: table, json, jsonl, csv, auto",
	FlagJSONFieldsUsage: "JSON 출력 필드 목록 (쉼표 구분, 예: market,trade_price)",
	FlagForceUsage:      "확인 프롬프트 스킵",

	// Root errors
	ErrConfigLoad:  "설정 로드 실패",
	ErrAuthRequired: "인증이 필요합니다: ACCESS_KEY 및 SECRET_KEY를 설정하세요",

	// Args helper
	MsgUsagePrefix: "사용법",

	// Balance
	HdrCurrency:   "통화",
	HdrBalance:    "잔고",
	HdrLocked:     "주문중",
	HdrAvgBuyPrice: "평균 매수가",
	HdrEvalKRW:    "평가금액(KRW)",

	MsgBalanceShort:    "계정 잔고 조회",
	ErrBalanceNotFound: "%s 잔고를 찾을 수 없습니다",
	MsgBalanceEmpty:    "보유 자산이 없습니다",

	// Buy
	MsgBuyShort:          "매수 주문",
	ErrBuyArgsRequired:   "마켓 코드를 지정하세요 (예: KRW-BTC)",
	ErrBuyParamsRequired: "매수 주문에는 --price와 --volume (지정가) 또는 --total (시장가)이 필요합니다",
	MsgTickCheckFailed:   "호가 단위 확인 실패: %v\n",
	MsgBuyOrderAdjusted:  "매수 주문: %s 단가=%s (호가 보정: %s→%s), 수량=%s (유형: %s)",
	MsgBuyOrderNormal:    "매수 주문: %s %s (유형: %s)",
	MsgTickAdjusted:      "호가 보정: %s → %s\n",
	MsgOrderCancelled:    "주문이 취소되었습니다",
	MsgDescLimitOrder:    "단가=%s, 수량=%s",
	MsgDescPriceOrder:    "총액=%s",
	FlagPriceUsage:       "주문 단가",
	FlagVolumeUsage:      "주문 수량 (50%와 같이 퍼센트 지정 가능)",
	FlagTotalUsage:       "주문 총액 (시장가 매수, 50%와 같이 퍼센트 지정 가능)",
	FlagSMPUsage:         "자기 거래 방지 (cancel_maker, cancel_taker, reduce)",
	FlagIdentifierUsage:  "클라이언트 지정 주문 식별자",
	FlagTestUsage:        "테스트 주문 (실제 체결 안됨)",
	FlagBestUsage:        "최유리 지정가 주문",
	FlagWatchPriceUsage:  "예약 주문 감시가 (해당 가격 도달 시 주문 발동)",
	MsgDescBestOrder:     "최유리, 수량=%s",
	MsgDescReservedOrder: "예약, 감시가=%s, 주문가=%s, 수량=%s",

	// Price keyword
	ErrTickerFetch:          "현재가 조회 실패",
	ErrTickerEmpty:          "%s 현재가 정보가 없습니다",
	MsgPriceKeywordResolved: "가격 키워드 %s → %s\n",

	// Percent
	ErrPercentParse:        "퍼센트 값을 파싱할 수 없습니다",
	ErrPercentRange:        "퍼센트는 0 초과 100 이하여야 합니다",
	ErrPercentBuyNeedsPrice: "퍼센트 수량 매수 시 --price가 필요합니다",
	ErrChanceFetch:         "주문 가능 정보 조회 실패",
	MsgPercentResolved:     "%s의 %s%% → %s\n",

	// Sell
	MsgSellShort:          "매도 주문",
	ErrSellArgsRequired:   "마켓 코드를 지정하세요 (예: KRW-BTC)",
	ErrSellParamsRequired: "매도 주문에는 --price와 --volume (지정가) 또는 --volume만 (시장가)이 필요합니다",
	MsgSellOrderAdjusted:  "매도 주문: %s 단가=%s (호가 보정: %s→%s), 수량=%s (유형: %s)",
	MsgSellOrderNormal:    "매도 주문: %s %s (유형: %s)",
	MsgDescMarketSell:     "수량=%s",

	// Market
	HdrMarket:      "마켓",
	HdrKoreanName:  "한글명",
	HdrEnglishName: "영문명",

	MsgMarketShort:       "마켓 목록 조회",
	MsgMarketFilterEmpty: "--quote 필터 결과가 없습니다",
	FlagQuoteUsage:       "호가 통화 필터 (KRW, BTC, USDT)",

	// Ticker
	HdrPrice:        "현재가",
	HdrChange:       "전일 대비",
	HdrChangeRate:   "변동률",
	HdrVolume24h:    "거래량(24h)",
	HdrTradePrice24h: "거래대금(24h)",
	HdrHigh:         "고가",
	HdrLow:          "저가",

	MsgTickerShort:       "현재가 조회",
	ErrTickerNoMarket:    "마켓 코드를 지정하세요 (예: KRW-BTC) 또는 --quote 플래그를 사용하세요\n\n사용법: %s",
	FlagTickerQuoteUsage: "마켓 통화 코드로 전체 시세 조회 (예: KRW, BTC, USDT)",

	// API Keys
	HdrAccessKey: "키(마스킹)",
	HdrExpireAt:  "만료일",

	MsgApiKeysShort:    "API 키 목록 조회",
	MsgApiKeysEmpty:    "유효한 API 키가 없습니다 (--all로 만료 포함 확인)",
	FlagApiKeysAllUsage: "만료된 키 포함 전체 표시",

	// Candle
	HdrCandleTime:   "시각",
	HdrOpen:         "시가",
	HdrClose:        "종가",
	HdrCandleVolume: "거래량",

	MsgCandleShort:         "캔들 조회",
	ErrCandleMarketRequired: "마켓 코드를 지정하세요 (예: KRW-BTC)",
	FlagIntervalUsage:      "캔들 간격 (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m, 1d, 1w, 1M, 1y)",
	FlagCountUsage:         "조회 개수",
	FlagFromUsage:          "시작 시각 (예: 2025-01-01, 2025-01-01T09:00:00+09:00)",
	FlagAscUsage:           "오래된 순 정렬 (기본)",
	FlagDescUsage:          "최신 순 정렬",
	FlagNoCacheUsage:       "캐시 무시",
	MsgCandleConfirm:       "약 %d개 캔들, %d회 API 호출이 필요합니다. 계속하시겠습니까?",
	MsgCancelled:           "취소되었습니다.",
	ErrFromParse:           "--from 파싱 실패",
	ErrUnrecognizedTime:    "인식할 수 없는 시각 형식: %s (예: 2025-01-01, 2025-01-01T09:00:00+09:00)",

	MsgCacheInitFailed:  "캐시 초기화 실패, API 직접 호출:",
	MsgCacheRangeError:  "캐시 범위 조회 실패:",
	MsgCacheSaveError:   "캐시 저장 실패:",
	MsgCacheUpdateError: "캐시 업데이트 실패:",
	ErrOlderRangeFetch:  "이전 구간 조회 실패",
	ErrNewerRangeFetch:  "최신 구간 조회 실패",
	ErrCacheQuery:       "캐시 조회 실패",

	// Orderbook
	MsgOrderbookShort:      "호가 조회",
	ErrOrderbookMarket:     "마켓 코드를 지정하세요 (예: KRW-BTC)",
	MsgOrderbookMarket:     "마켓: %s\n",
	MsgOrderbookTotalSizes: "총매도량: %s  총매수량: %s\n",
	HdrAskSize:             "매도잔량",
	HdrAskPrice:            "매도호가",
	HdrBidPrice:            "매수호가",
	HdrBidSize:             "매수잔량",

	// Orderbook Levels
	HdrSupportedLevels:      "지원 단위",
	MsgOrderbookLevelsShort: "호가 모아보기 단위 조회",

	// Trades
	HdrTradeTime:   "체결 시각",
	HdrTradePrice:  "체결가",
	HdrTradeVolume: "체결량",
	HdrAskBid:      "매수/매도",
	HdrChangePrice: "전일 대비",

	MsgTradesShort: "체결 내역 조회",

	// Order
	MsgOrderShort: "주문 관리 (조회, 취소, 정정)",

	// Order List
	HdrSide:        "방향",
	HdrOrdType:     "유형",
	HdrOrderPrice:  "가격",
	HdrOrderVolume: "수량",
	HdrExecutedVol: "체결량",
	HdrState:       "상태",
	HdrCreatedAt:   "생성시각",

	MsgOrderListShort:     "주문 목록 조회",
	MsgOrderListClosed:    "종료된 주문이 없습니다",
	MsgOrderListOpen:      "대기 중인 주문이 없습니다",
	FlagClosedUsage:       "종료 주문 조회",
	FlagMarketFilterUsage: "마켓 필터 (예: KRW-BTC)",
	FlagPageUsage:         "페이지 번호",

	// Order Show
	HdrRemainingVol:      "잔여수량",
	MsgOrderShowShort:    "주문 상세 조회",
	ErrOrderShowArgs:     "주문 UUID를 지정하세요",
	FlagOrderShowIDUsage: "Identifier로 조회 (UUID 대신)",

	// Order Cancel
	MsgOrderCancelShort:   "주문 취소",
	MsgCancelAllOrders:    "모든 대기 주문을 취소합니다",
	MsgCancelAllMarket:    "%s 마켓의 모든 대기 주문을 취소합니다",
	MsgCancelAborted:      "취소가 중단되었습니다",
	MsgCancelSuccess:      "취소 성공  %d건\n",
	MsgCancelFailed:       "취소 실패  %d건\n",
	ErrCancelNoUUID:       "취소할 주문의 UUID를 지정하거나 --all 플래그를 사용하세요",
	MsgCancelSingleOrder:  "주문 %s을(를) 취소합니다",
	FlagCancelAllUsage:    "모든 대기 주문 일괄 취소",
	FlagCancelMarketUsage: "마켓 필터 (--all과 함께 사용)",

	// Order Replace
	MsgOrderReplaceShort:   "주문 정정 (취소 후 재주문)",
	ErrReplaceArgs:         "정정할 주문 UUID를 지정하세요",
	ErrReplaceNoParam:      "정정할 --price 또는 --volume을 지정하세요",
	ErrReplaceLimitNoPrice: "limit 주문 정정에는 --price가 필요합니다",
	MsgReplaceConfirm:      "주문 %s 정정: ",
	MsgReplacePrice:        "단가=%s ",
	MsgReplaceVolume:       "수량=%s",
	MsgReplaceRemain:       "(잔량 유지)",
	MsgReplaceCancelled:    "정정이 취소되었습니다",
	FlagNewPriceUsage:      "신규 주문 단가",
	FlagNewVolumeUsage:     "신규 주문 수량",
	FlagOrdTypeUsage:       "주문 유형 (limit, best)",

	// Order Chance
	MsgOrderChanceShort: "마켓별 주문 가능 정보 조회",
	LblMarket:           "마켓",
	LblState:            "상태",
	LblBidFee:           "매수 수수료",
	LblAskFee:           "매도 수수료",
	LblBidAvailable:     "매수 가능",
	LblAskAvailable:     "매도 가능",
	LblMinOrder:         "최소 주문",
	LblBidTypes:         "매수 유형",
	LblAskTypes:         "매도 유형",

	// Watch
	MsgWatchShort:          "실시간 WebSocket 스트림 구독",
	MsgWatchTickerShort:    "현재가 실시간 스트림",
	MsgWatchOrderbookShort: "호가 실시간 스트림",
	MsgWatchTradeShort:     "체결 실시간 스트림",
	MsgWatchCandleShort:    "캔들 실시간 스트림",
	MsgWatchMyOrderShort:   "내 주문 실시간 스트림 (인증 필요)",
	MsgWatchMyAssetShort:   "내 자산 실시간 스트림 (인증 필요)",
	FlagWatchIntervalUsage: "캔들 단위 (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m)",

	ErrMarketNotFound: "존재하지 않는 마켓: %s\n\n유효한 마켓 목록은 'upbit market' 명령으로 확인하세요",
	ErrSubscribeBuild: "구독 메시지 생성 실패",
	ErrSubscribeSend:  "구독 전송 실패",
	ErrMessageReceive: "메시지 수신 오류",
	ErrServerError:    "서버 에러: %s",

	WatchVolume:  "거래량",
	WatchSell:    "매도",
	WatchBuy:     "매수",
	WatchSpread:  "스프레드",
	WatchQty:     "수량",
	WatchOpen:    "시",
	WatchHigh:    "고",
	WatchLow:     "저",
	WatchClose:   "종",
	WatchState:   "상태",
	WatchLocking: "주문중",

	// Wallet / Status
	HdrWalletState: "지갑 상태",
	HdrBlockState:  "블록 상태",
	HdrNetwork:     "네트워크",
	HdrNetworkType: "네트워크 타입",
	MsgWalletShort: "입출금 서비스 상태 조회",

	// Tick Size
	HdrQuoteCurrency: "호가통화",
	HdrTickSize:      "호가단위",
	MsgTickSizeShort: "호가 단위 조회",

	// Cache command
	HdrPath:             "경로",
	HdrFiles:            "파일",
	HdrSize:             "크기",
	MsgCacheShort:       "캐시 관리",
	MsgCacheFileCount:   "%d개",
	ErrCachePathFailed:  "캐시 경로 확인 실패",
	MsgCacheConfirmClear: "캐시를 삭제하시겠습니까?",
	ErrCacheOpenFailed:  "캐시 열기 실패",
	ErrCacheClearFailed: "캐시 삭제 실패",
	MsgCacheCleared:     "캐시가 삭제되었습니다.",
	FlagCacheClearUsage: "캐시 삭제",

	// Deposit
	MsgDepositShort:        "입금 관리 (조회, 주소)",
	MsgDepositListShort:    "입금 목록 조회",
	MsgDepositShowShort:    "개별 입금 조회",
	MsgDepositAddressShort: "입금 주소 조회",
	MsgDepositCreateShort:  "입금 주소 생성 요청",
	MsgDepositListEmpty:    "입금 내역이 없습니다",
	MsgDepositAddressEmpty: "등록된 입금 주소가 없습니다 (upbit deposit address create <currency>로 생성)",
	MsgDepositAddrCreated:  "주소 생성이 요청되었습니다.\n",
	MsgDepositAddrCheck:    "'upbit deposit address %s --net-type %s' 명령으로 생성 결과를 확인하세요.\n",
	ErrDepositUUIDRequired: "입금 UUID를 지정하세요",
	ErrDepositCurrRequired: "통화 코드를 지정하세요 (예: BTC)",

	HdrAmount:      "금액",
	HdrFee:         "수수료",
	HdrTransType:   "유형",
	HdrTXID:        "TXID",
	HdrDoneAt:      "완료일",
	HdrDepositAddr: "주소",
	HdrTag:         "태그",

	FlagCurrencyUsage: "통화 코드 필터 (예: BTC, KRW)",
	FlagStateUsage:    "입금 상태 필터 (PROCESSING, ACCEPTED, CANCELLED, REJECTED)",
	FlagNetTypeUsage:  "네트워크 유형 (기본: currency와 동일)",

	// Withdraw
	MsgWithdrawShort:         "출금 관리 (조회, 요청, 취소)",
	MsgWithdrawListShort:     "출금 목록 조회",
	MsgWithdrawShowShort:     "개별 출금 조회",
	MsgWithdrawRequestShort:  "출금 요청",
	MsgWithdrawCancelShort:   "출금 취소",
	MsgWithdrawListEmpty:     "출금 내역이 없습니다",
	MsgWithdrawCancelled:     "출금이 취소되었습니다",
	ErrWithdrawUUIDRequired:  "출금 UUID를 지정하세요",
	ErrWithdrawCurrRequired:  "출금할 통화 코드를 지정하세요 (예: BTC, KRW)",
	ErrAmountRequired:        "--amount 플래그는 필수입니다",
	ErrTwoFactorRequired:     "KRW 출금 시 --two-factor 플래그가 필요합니다 (kakao, naver, hana)",
	MsgWithdrawConfirmKRW:    "출금: KRW %s",
	ErrWithdrawAddrRequired:  "디지털 자산 출금에는 --to (수신 주소) 플래그가 필요합니다",
	MsgWithdrawConfirmCoin:   "출금: %s %s -> %s",
	MsgWithdrawCancelConfirm: "출금 %s을(를) 취소합니다",

	FlagAmountUsage:        "출금 수량/금액 (필수)",
	FlagToUsage:            "수신 주소 (디지털 자산 필수)",
	FlagSecondaryAddrUsage: "2차 주소 (태그/메모)",
	FlagTxTypeUsage:        "트랜잭션 유형 (default, internal)",
	FlagTwoFactorUsage:     "2차 인증 수단 (kakao, naver, hana) — KRW 전용",
	FlagWithdrawStateUsage: "출금 상태 필터 (WAITING, PROCESSING, DONE, FAILED, CANCELLED, REJECTED)",

	// Travelrule
	MsgTravelruleShort:      "트래블룰 검증 관리",
	MsgTravelruleVaspsShort: "트래블룰 지원 거래소 목록 조회",
	MsgTravelruleTxIDShort:  "TxID 기반 트래블룰 검증 요청",
	MsgTravelruleUUIDShort:  "UUID 기반 트래블룰 검증 요청",
	ErrTravelruleTxIDArgs:   "검증할 입금의 TxID를 지정하세요",
	ErrTravelruleUUIDArgs:   "검증할 입금의 UUID를 지정하세요",

	HdrVaspName:     "이름",
	HdrVaspNameEn:   "이름(영문)",
	HdrDepositable:  "입금가능",
	HdrWithdrawable: "출금가능",
	HdrDepositUUID:  "입금UUID",
	HdrVerifyResult: "검증결과",
	HdrDepositState: "입금상태",

	FlagVaspUsage:       "상대 거래소 UUID (필수)",
	FlagTRCurrencyUsage: "통화 코드 (예: ETH) (필수)",
	FlagTRNetTypeUsage:  "네트워크 식별자 (예: ETH) (필수)",

	// Tool Schema
	MsgToolSchemaShort:    "명령어 JSON Schema 출력 (LLM/MCP용)",
	ErrToolSchemaNotFound: "명령을 찾을 수 없습니다: %s",

	ArgMarketCode:      "마켓 코드 (예: KRW-BTC, KRW-ETH)",
	ArgOrderbookMarket: "마켓 코드 (예: KRW-BTC)",
	ArgTradesMarket:    "마켓 코드 (예: KRW-BTC)",
	ArgCandleMarket:    "마켓 코드 (예: KRW-BTC)",
	ArgBuyMarket:       "매수할 마켓 코드 (예: KRW-BTC)",
	ArgSellMarket:      "매도할 마켓 코드 (예: KRW-BTC)",
	ArgBalanceCurrency: "조회할 통화 코드 (예: KRW, BTC)",
	ArgShowUUID:        "주문/입금/출금 UUID",
	ArgCancelUUID:      "취소할 주문/출금 UUID",
	ArgReplaceUUID:     "정정할 주문 UUID",
	ArgRequestCurrency: "출금할 통화 코드 (예: BTC, KRW)",
	ArgAddressCurrency: "입금 주소 조회할 통화 코드 (예: BTC)",
	ArgChanceMarket:    "주문 가능 정보 조회할 마켓 코드 (예: KRW-BTC)",

	// Price Adjust
	ErrTickSizeFetch: "호가 단위 조회 실패",
	ErrTickSizeEmpty: "호가 단위 정보 없음: %s",
	ErrTickSizeParse: "호가 단위 파싱 실패",
	ErrTickSizeZero:  "호가 단위가 0: %s",
	ErrPriceParse:    "가격 파싱 실패",
	ErrUnknownSide:   "알 수 없는 side: %s",

	// Update
	MsgUpdateShort:         "최신 버전 확인 및 업데이트",
	MsgUpdateChecking:      "최신 버전 확인 중...",
	MsgUpdateLatest:        "최신 버전: %s",
	MsgUpdateAvailable:     "업데이트 가능: %s → %s",
	MsgUpdateDownloading:   "다운로드 중...",
	MsgUpdateComplete:      "업데이트 완료: %s → %s",
	MsgUpdateAlreadyLatest: "이미 최신 버전입니다 (%s)",
	ErrUpdateFetch:         "최신 버전 확인 실패",
	ErrUpdateNoAsset:       "현재 플랫폼(%s/%s)에 맞는 바이너리를 찾을 수 없습니다",
	ErrUpdateDownload:      "다운로드 실패",
	ErrUpdateReplace:       "바이너리 교체 실패",
	ErrUpdateChecksum:      "체크섬 불일치",
	MsgUpdateVerifying:     "무결성 검증 중...",
	FlagUpdateCheckUsage:   "업데이트 확인만 (다운로드 안 함)",

	// TUI
	TUIQuitHint:       "q: 종료",
	TUITabHint:        "←/→: 마켓 전환  q: 종료",
	TUITickerHeader:   "마켓          현재가           변동률       거래량",
	TUIOrderbookTitle: "호가창",
	TUITotalAsk:       "총매도",
	TUITotalBid:       "총매수",
	TUITradeTitle:     "체결 내역",
	TUICandleTitle:    "캔들 차트",
	TUIHighest:        "최고",
	TUILowest:         "최저",
	TUIVar:            "등락",
	TUIAvg:            "평균",
	TUICumVol:         "누적거래량",
	TUIConfirmYes:     "[Y] 확인",
	TUIConfirmNo:      "[N] 취소",

	// Confirm
	MsgConfirmNonTTY: "확인 프롬프트: non-tty 환경에서는 --force 플래그가 필요합니다.",
	ErrInputRead:     "입력 읽기 실패",
}

var en = map[Key]string{
	// Root
	MsgRootShort: "Upbit Exchange CLI",
	MsgRootLong:  "Upbit Exchange CLI — market data, trading, deposits/withdrawals, and real-time streaming.\n\n  https://github.com/kyungw00k/upbit\n  Sponsor: https://github.com/sponsors/kyungw00k",

	// Groups
	GroupQuotation: "Quotation commands:",
	GroupTrading:   "Trading commands:",
	GroupWallet:    "Wallet commands:",
	GroupRealtime:  "Real-time commands:",
	GroupUtil:      "Utility commands:",

	// Global flags
	FlagOutputUsage:     "Output format: table, json, jsonl, csv, auto",
	FlagJSONFieldsUsage: "JSON output fields (comma-separated, e.g. market,trade_price)",
	FlagForceUsage:      "Skip confirmation prompt",

	// Root errors
	ErrConfigLoad:  "Failed to load config",
	ErrAuthRequired: "Authentication required: set ACCESS_KEY and SECRET_KEY",

	// Args helper
	MsgUsagePrefix: "Usage",

	// Balance
	HdrCurrency:   "Currency",
	HdrBalance:    "Balance",
	HdrLocked:     "Locked",
	HdrAvgBuyPrice: "Avg Buy Price",
	HdrEvalKRW:    "Value(KRW)",

	MsgBalanceShort:    "View account balance",
	ErrBalanceNotFound: "Balance not found for %s",
	MsgBalanceEmpty:    "No assets held",

	// Buy
	MsgBuyShort:          "Place buy order",
	ErrBuyArgsRequired:   "Specify market code (e.g. KRW-BTC)",
	ErrBuyParamsRequired: "Buy order requires --price and --volume (limit) or --total (market)",
	MsgTickCheckFailed:   "Tick size check failed: %v\n",
	MsgBuyOrderAdjusted:  "Buy order: %s price=%s (tick adjusted: %s->%s), volume=%s (type: %s)",
	MsgBuyOrderNormal:    "Buy order: %s %s (type: %s)",
	MsgTickAdjusted:      "Tick adjusted: %s -> %s\n",
	MsgOrderCancelled:    "Order cancelled",
	MsgDescLimitOrder:    "price=%s, volume=%s",
	MsgDescPriceOrder:    "total=%s",
	FlagPriceUsage:       "Order price",
	FlagVolumeUsage:      "Order volume (supports percentage like 50%)",
	FlagTotalUsage:       "Order total for market buy (supports percentage like 50%)",
	FlagSMPUsage:         "Self-trade prevention (cancel_maker, cancel_taker, reduce)",
	FlagIdentifierUsage:  "Client-specified order identifier",
	FlagTestUsage:        "Test order (no actual execution)",
	FlagBestUsage:        "Best price limit order",
	FlagWatchPriceUsage:  "Reserved order watch price (triggers when reached)",
	MsgDescBestOrder:     "best price, volume=%s",
	MsgDescReservedOrder: "reserved, watch=%s, price=%s, volume=%s",

	// Price keyword
	ErrTickerFetch:          "Failed to fetch ticker",
	ErrTickerEmpty:          "No ticker data for %s",
	MsgPriceKeywordResolved: "Price keyword %s → %s\n",

	// Percent
	ErrPercentParse:        "Cannot parse percent value",
	ErrPercentRange:        "Percent must be between 0 (exclusive) and 100 (inclusive)",
	ErrPercentBuyNeedsPrice: "--price is required for percent volume buy",
	ErrChanceFetch:         "Failed to fetch order chance",
	MsgPercentResolved:     "%s %s%% → %s\n",

	// Sell
	MsgSellShort:          "Place sell order",
	ErrSellArgsRequired:   "Specify market code (e.g. KRW-BTC)",
	ErrSellParamsRequired: "Sell order requires --price and --volume (limit) or --volume only (market)",
	MsgSellOrderAdjusted:  "Sell order: %s price=%s (tick adjusted: %s->%s), volume=%s (type: %s)",
	MsgSellOrderNormal:    "Sell order: %s %s (type: %s)",
	MsgDescMarketSell:     "volume=%s",

	// Market
	HdrMarket:      "Market",
	HdrKoreanName:  "Korean Name",
	HdrEnglishName: "English Name",

	MsgMarketShort:       "List markets",
	MsgMarketFilterEmpty: "No results for --quote filter",
	FlagQuoteUsage:       "Quote currency filter (KRW, BTC, USDT)",

	// Ticker
	HdrPrice:        "Price",
	HdrChange:       "Change",
	HdrChangeRate:   "Change%",
	HdrVolume24h:    "Volume(24h)",
	HdrTradePrice24h: "Trade Amt(24h)",
	HdrHigh:         "High",
	HdrLow:          "Low",

	MsgTickerShort:       "View current price",
	ErrTickerNoMarket:    "Specify market code (e.g. KRW-BTC) or use --quote flag\n\nUsage: %s",
	FlagTickerQuoteUsage: "View all prices by quote currency (e.g. KRW, BTC, USDT)",

	// API Keys
	HdrAccessKey: "Key(masked)",
	HdrExpireAt:  "Expires",

	MsgApiKeysShort:    "List API keys",
	MsgApiKeysEmpty:    "No valid API keys (use --all to include expired)",
	FlagApiKeysAllUsage: "Show all including expired",

	// Candle
	HdrCandleTime:   "Time",
	HdrOpen:         "Open",
	HdrClose:        "Close",
	HdrCandleVolume: "Volume",

	MsgCandleShort:         "View candles",
	ErrCandleMarketRequired: "Specify market code (e.g. KRW-BTC)",
	FlagIntervalUsage:      "Candle interval (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m, 1d, 1w, 1M, 1y)",
	FlagCountUsage:         "Number of results",
	FlagFromUsage:          "Start time (e.g. 2025-01-01, 2025-01-01T09:00:00+09:00)",
	FlagAscUsage:           "Sort oldest first (default)",
	FlagDescUsage:          "Sort newest first",
	FlagNoCacheUsage:       "Skip cache",
	MsgCandleConfirm:       "Approximately %d candles, %d API calls required. Continue?",
	MsgCancelled:           "Cancelled.",
	ErrFromParse:           "Failed to parse --from",
	ErrUnrecognizedTime:    "Unrecognized time format: %s (e.g. 2025-01-01, 2025-01-01T09:00:00+09:00)",

	MsgCacheInitFailed:  "Cache init failed, calling API directly:",
	MsgCacheRangeError:  "Cache range query failed:",
	MsgCacheSaveError:   "Cache save failed:",
	MsgCacheUpdateError: "Cache update failed:",
	ErrOlderRangeFetch:  "Failed to fetch older range",
	ErrNewerRangeFetch:  "Failed to fetch newer range",
	ErrCacheQuery:       "Cache query failed",

	// Orderbook
	MsgOrderbookShort:      "View orderbook",
	ErrOrderbookMarket:     "Specify market code (e.g. KRW-BTC)",
	MsgOrderbookMarket:     "Market: %s\n",
	MsgOrderbookTotalSizes: "Total Ask: %s  Total Bid: %s\n",
	HdrAskSize:             "Ask Size",
	HdrAskPrice:            "Ask Price",
	HdrBidPrice:            "Bid Price",
	HdrBidSize:             "Bid Size",

	// Orderbook Levels
	HdrSupportedLevels:      "Supported Levels",
	MsgOrderbookLevelsShort: "View orderbook level units",

	// Trades
	HdrTradeTime:   "Trade Time",
	HdrTradePrice:  "Trade Price",
	HdrTradeVolume: "Trade Vol",
	HdrAskBid:      "Ask/Bid",
	HdrChangePrice: "Change",

	MsgTradesShort: "View recent trades",

	// Order
	MsgOrderShort: "Order management (list, cancel, replace)",

	// Order List
	HdrSide:        "Side",
	HdrOrdType:     "Type",
	HdrOrderPrice:  "Price",
	HdrOrderVolume: "Volume",
	HdrExecutedVol: "Exec Vol",
	HdrState:       "State",
	HdrCreatedAt:   "Created",

	MsgOrderListShort:     "List orders",
	MsgOrderListClosed:    "No closed orders",
	MsgOrderListOpen:      "No pending orders",
	FlagClosedUsage:       "Show closed orders",
	FlagMarketFilterUsage: "Market filter (e.g. KRW-BTC)",
	FlagPageUsage:         "Page number",

	// Order Show
	HdrRemainingVol:      "Remaining",
	MsgOrderShowShort:    "View order details",
	ErrOrderShowArgs:     "Specify order UUID",
	FlagOrderShowIDUsage: "Query by identifier (instead of UUID)",

	// Order Cancel
	MsgOrderCancelShort:   "Cancel order",
	MsgCancelAllOrders:    "Cancel all pending orders",
	MsgCancelAllMarket:    "Cancel all pending orders for %s",
	MsgCancelAborted:      "Cancellation aborted",
	MsgCancelSuccess:      "Cancelled  %d order(s)\n",
	MsgCancelFailed:       "Failed     %d order(s)\n",
	ErrCancelNoUUID:       "Specify order UUID to cancel or use --all flag",
	MsgCancelSingleOrder:  "Cancel order %s",
	FlagCancelAllUsage:    "Cancel all pending orders",
	FlagCancelMarketUsage: "Market filter (use with --all)",

	// Order Replace
	MsgOrderReplaceShort:   "Replace order (cancel and re-order)",
	ErrReplaceArgs:         "Specify order UUID to replace",
	ErrReplaceNoParam:      "Specify --price or --volume to replace",
	ErrReplaceLimitNoPrice: "--price is required for limit order replacement",
	MsgReplaceConfirm:      "Replace order %s: ",
	MsgReplacePrice:        "price=%s ",
	MsgReplaceVolume:       "volume=%s",
	MsgReplaceRemain:       "(keep remaining)",
	MsgReplaceCancelled:    "Replacement cancelled",
	FlagNewPriceUsage:      "New order price",
	FlagNewVolumeUsage:     "New order volume",
	FlagOrdTypeUsage:       "Order type (limit, best)",

	// Order Chance
	MsgOrderChanceShort: "View order chance by market",
	LblMarket:           "Market",
	LblState:            "State",
	LblBidFee:           "Buy Fee",
	LblAskFee:           "Sell Fee",
	LblBidAvailable:     "Buy Available",
	LblAskAvailable:     "Sell Available",
	LblMinOrder:         "Min Order",
	LblBidTypes:         "Buy Types",
	LblAskTypes:         "Sell Types",

	// Watch
	MsgWatchShort:          "Subscribe to real-time WebSocket stream",
	MsgWatchTickerShort:    "Real-time ticker stream",
	MsgWatchOrderbookShort: "Real-time orderbook stream",
	MsgWatchTradeShort:     "Real-time trade stream",
	MsgWatchCandleShort:    "Real-time candle stream",
	MsgWatchMyOrderShort:   "Real-time my orders stream (auth required)",
	MsgWatchMyAssetShort:   "Real-time my assets stream (auth required)",
	FlagWatchIntervalUsage: "Candle unit (1s, 1m, 3m, 5m, 10m, 15m, 30m, 60m, 240m)",

	ErrMarketNotFound: "Market not found: %s\n\nCheck available markets with 'upbit market' command",
	ErrSubscribeBuild: "Failed to build subscribe message",
	ErrSubscribeSend:  "Failed to send subscription",
	ErrMessageReceive: "Message receive error",
	ErrServerError:    "Server error: %s",

	WatchVolume:  "vol",
	WatchSell:    "SELL",
	WatchBuy:     "BUY",
	WatchSpread:  "spread",
	WatchQty:     "qty",
	WatchOpen:    "O",
	WatchHigh:    "H",
	WatchLow:     "L",
	WatchClose:   "C",
	WatchState:   "state",
	WatchLocking: "locked",

	// Wallet / Status
	HdrWalletState: "Wallet State",
	HdrBlockState:  "Block State",
	HdrNetwork:     "Network",
	HdrNetworkType: "Network Type",
	MsgWalletShort: "View wallet service status",

	// Tick Size
	HdrQuoteCurrency: "Quote Currency",
	HdrTickSize:      "Tick Size",
	MsgTickSizeShort: "View tick size",

	// Cache command
	HdrPath:             "Path",
	HdrFiles:            "Files",
	HdrSize:             "Size",
	MsgCacheShort:       "Cache management",
	MsgCacheFileCount:   "%d",
	ErrCachePathFailed:  "Failed to get cache path",
	MsgCacheConfirmClear: "Delete cache?",
	ErrCacheOpenFailed:  "Failed to open cache",
	ErrCacheClearFailed: "Failed to clear cache",
	MsgCacheCleared:     "Cache cleared.",
	FlagCacheClearUsage: "Clear cache",

	// Deposit
	MsgDepositShort:        "Deposit management (list, address)",
	MsgDepositListShort:    "List deposits",
	MsgDepositShowShort:    "View deposit details",
	MsgDepositAddressShort: "View deposit address",
	MsgDepositCreateShort:  "Request deposit address creation",
	MsgDepositListEmpty:    "No deposit history",
	MsgDepositAddressEmpty: "No deposit addresses (create with: upbit deposit address create <currency>)",
	MsgDepositAddrCreated:  "Address creation requested.\n",
	MsgDepositAddrCheck:    "Check result with: 'upbit deposit address %s --net-type %s'\n",
	ErrDepositUUIDRequired: "Specify deposit UUID",
	ErrDepositCurrRequired: "Specify currency code (e.g. BTC)",

	HdrAmount:      "Amount",
	HdrFee:         "Fee",
	HdrTransType:   "Type",
	HdrTXID:        "TXID",
	HdrDoneAt:      "Completed",
	HdrDepositAddr: "Address",
	HdrTag:         "Tag",

	FlagCurrencyUsage: "Currency code filter (e.g. BTC, KRW)",
	FlagStateUsage:    "Deposit state filter (PROCESSING, ACCEPTED, CANCELLED, REJECTED)",
	FlagNetTypeUsage:  "Network type (default: same as currency)",

	// Withdraw
	MsgWithdrawShort:         "Withdrawal management (list, request, cancel)",
	MsgWithdrawListShort:     "List withdrawals",
	MsgWithdrawShowShort:     "View withdrawal details",
	MsgWithdrawRequestShort:  "Request withdrawal",
	MsgWithdrawCancelShort:   "Cancel withdrawal",
	MsgWithdrawListEmpty:     "No withdrawal history",
	MsgWithdrawCancelled:     "Withdrawal cancelled",
	ErrWithdrawUUIDRequired:  "Specify withdrawal UUID",
	ErrWithdrawCurrRequired:  "Specify currency code (e.g. BTC, KRW)",
	ErrAmountRequired:        "--amount flag is required",
	ErrTwoFactorRequired:     "--two-factor flag required for KRW withdrawal (kakao, naver, hana)",
	MsgWithdrawConfirmKRW:    "Withdraw: KRW %s",
	ErrWithdrawAddrRequired:  "--to (recipient address) flag required for digital asset withdrawal",
	MsgWithdrawConfirmCoin:   "Withdraw: %s %s -> %s",
	MsgWithdrawCancelConfirm: "Cancel withdrawal %s",

	FlagAmountUsage:        "Withdrawal amount (required)",
	FlagToUsage:            "Recipient address (required for digital assets)",
	FlagSecondaryAddrUsage: "Secondary address (tag/memo)",
	FlagTxTypeUsage:        "Transaction type (default, internal)",
	FlagTwoFactorUsage:     "2FA method (kakao, naver, hana) — KRW only",
	FlagWithdrawStateUsage: "Withdrawal state filter (WAITING, PROCESSING, DONE, FAILED, CANCELLED, REJECTED)",

	// Travelrule
	MsgTravelruleShort:      "Travel rule verification",
	MsgTravelruleVaspsShort: "List travel rule supported exchanges",
	MsgTravelruleTxIDShort:  "Travel rule verification by TxID",
	MsgTravelruleUUIDShort:  "Travel rule verification by UUID",
	ErrTravelruleTxIDArgs:   "Specify deposit TxID to verify",
	ErrTravelruleUUIDArgs:   "Specify deposit UUID to verify",

	HdrVaspName:     "Name",
	HdrVaspNameEn:   "Name(EN)",
	HdrDepositable:  "Depositable",
	HdrWithdrawable: "Withdrawable",
	HdrDepositUUID:  "Deposit UUID",
	HdrVerifyResult: "Result",
	HdrDepositState: "Deposit State",

	FlagVaspUsage:       "Counterparty exchange UUID (required)",
	FlagTRCurrencyUsage: "Currency code (e.g. ETH) (required)",
	FlagTRNetTypeUsage:  "Network identifier (e.g. ETH) (required)",

	// Tool Schema
	MsgToolSchemaShort:    "Output command JSON Schema (for LLM/MCP)",
	ErrToolSchemaNotFound: "Command not found: %s",

	ArgMarketCode:      "Market code (e.g. KRW-BTC, KRW-ETH)",
	ArgOrderbookMarket: "Market code (e.g. KRW-BTC)",
	ArgTradesMarket:    "Market code (e.g. KRW-BTC)",
	ArgCandleMarket:    "Market code (e.g. KRW-BTC)",
	ArgBuyMarket:       "Market code to buy (e.g. KRW-BTC)",
	ArgSellMarket:      "Market code to sell (e.g. KRW-BTC)",
	ArgBalanceCurrency: "Currency code to query (e.g. KRW, BTC)",
	ArgShowUUID:        "Order/deposit/withdrawal UUID",
	ArgCancelUUID:      "UUID to cancel",
	ArgReplaceUUID:     "Order UUID to replace",
	ArgRequestCurrency: "Currency code to withdraw (e.g. BTC, KRW)",
	ArgAddressCurrency: "Currency code for deposit address (e.g. BTC)",
	ArgChanceMarket:    "Market code for order chance (e.g. KRW-BTC)",

	// Price Adjust
	ErrTickSizeFetch: "Failed to fetch tick size",
	ErrTickSizeEmpty: "No tick size info: %s",
	ErrTickSizeParse: "Failed to parse tick size",
	ErrTickSizeZero:  "Tick size is 0: %s",
	ErrPriceParse:    "Failed to parse price",
	ErrUnknownSide:   "Unknown side: %s",

	// Update
	MsgUpdateShort:         "Check and install latest version",
	MsgUpdateChecking:      "Checking for updates...",
	MsgUpdateLatest:        "Latest version: %s",
	MsgUpdateAvailable:     "Update available: %s → %s",
	MsgUpdateDownloading:   "Downloading...",
	MsgUpdateComplete:      "Updated: %s → %s",
	MsgUpdateAlreadyLatest: "Already up to date (%s)",
	ErrUpdateFetch:         "Failed to check latest version",
	ErrUpdateNoAsset:       "No binary found for %s/%s",
	ErrUpdateDownload:      "Download failed",
	ErrUpdateReplace:       "Failed to replace binary",
	ErrUpdateChecksum:      "Checksum mismatch",
	MsgUpdateVerifying:     "Verifying integrity...",
	FlagUpdateCheckUsage:   "Check only, don't download",

	// TUI
	TUIQuitHint:       "q: quit",
	TUITabHint:        "←/→: switch market  q: quit",
	TUITickerHeader:   "MARKET        PRICE            CHANGE       VOLUME",
	TUIOrderbookTitle: "Orderbook",
	TUITotalAsk:       "Total Ask",
	TUITotalBid:       "Total Bid",
	TUITradeTitle:     "Recent Trades",
	TUICandleTitle:    "Candle Chart",
	TUIHighest:        "Highest",
	TUILowest:         "Lowest",
	TUIVar:            "Var",
	TUIAvg:            "Avg",
	TUICumVol:         "CumVol",
	TUIConfirmYes:     "[Y] Confirm",
	TUIConfirmNo:      "[N] Cancel",

	// Confirm
	MsgConfirmNonTTY: "Confirmation prompt: --force flag required in non-tty environment.",
	ErrInputRead:     "Failed to read input",
}
