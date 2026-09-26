Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내] 기술적 지표 WebSocket 지원

안녕하세요. 업비트 개발자센터입니다.

업비트 Open API WebSocket에서 RSI, 이동평균선, 볼린저밴드 등 기술적 지표 데이터를 수신할 수 있는 신규 데이터 타입을 제공합니다.

`indicator`를 구독하면 선택한 캔들 단위의 캔들이 마감될 때마다 업비트에서 계산한 지표 값을 수신할 수 있으며, `indicator_signal`을 구독하면 지정한 조건이 미충족 상태에서 충족 상태로 전환되는 시점에 이벤트를 수신할 수 있습니다.

자세한 내용은 아래를 참고해 주세요.

## 1.적용 일시

* 적용일: 2026-09-23 (수)

## 2.신규 기능 내용

* 기능명: <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-indicator">기술적 지표</Anchor>
* Endpoint: `wss://api.upbit.com/websocket/v1/info `
* 데이터 타입
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-indicator">`indicator`</Anchor>: 선택한 캔들 단위의 캔들이 마감될 때마다 계산된 지표 값을 수신합니다.
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-indicator-signal">`indicator_signal`</Anchor>: 지정한 조건이 미충족(false)에서 충족(true)으로 전환되는 시점에 이벤트를 수신합니다.
* 지원 지표
  * 거래량
  * RSI
  * 이동평균선(단일)
  * 이동평균선(교차)
  * 볼린저밴드
  * MACD

※ 데이터 타입별로 지원하는 지표와 설정 가능한 계산 파라미터 및 조건이 다를 수 있습니다.&#x20;

* 지원 범위
  * 지원 마켓: 업비트 KRW 마켓
  * 지원 봉 단위: 5분봉, 15분봉, 60분봉(1시간), 240분봉(4시간), 일봉<br />※ 1분봉 및 틱 단위 데이터는 지원하지 않습니다.
  * 데이터 전송 기준: 선택한 시간 단위의 캔들이 마감될 때 전송
  * 계산 기준: 마감된 캔들의 종가(Close)

## 3.이용 방법

* Private WebSocket Endpoint에 연결한 후 API Key를 이용해 인증합니다.
* 구독하려는 데이터에 따라 `type`을 지정하고, 마켓 코드(`codes`), 캔들 단위(`timeframe`), 지표(`indicators`) 등을 함께 설정하여 구독을 요청합니다.
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-indicator">지표 값(Indicator)</Anchor>을 수신하려는 경우: indicator

    <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-indicator-signal">지표 조건 이벤트 (Indicator Signal)</Anchor>를 수신하려는 경우: indicator\_signal
* `indicator_signal`을 이용하는 경우 지표별 `trigger`를 함께 설정합니다.
* 기술적 지표 WebSocket은 Private WebSocket을 사용하므로 API Key 인증이 필요합니다. 다만, 본 기능을 이용하기 위한 별도의 API Key 권한(Scope)은 필요하지 않습니다.
* 자세한 요청·응답 형식과 지원 파라미터는 문서를 참고해 주세요.

## 유의사항

* 본 기능은 지표의 산식 결과 및 사전에 설정한 조건 충족 여부를 기술적으로 제공하는 기능이며, 투자 자문 및 권유의 목적이 아닙니다. 회사는 특정 종목이나 매매 시점을 추천하거나 매수 또는 매도를 권유하지 않습니다.
* 모든 투자 판단과 그 결과에 대한 책임은 이용자 본인에게 있습니다.
* 구독 시점의 최근 지표 값을 제공하는 스냅샷은 제공되지 않으며, 구독 이후 처음 마감되는 캔들부터 데이터가 전송됩니다.
* 신규 구독 요청에 대한 별도의 응답은 제공되지 않습니다. 현재 등록된 구독 정보는 [구독 중인 스트림 목록 조회(LIST\_SUBSCRIPTIONS)](https://docs.upbit.com/kr/reference/list-subscriptions)를 통해 확인할 수 있습니다.
* 조건 이벤트(indicator\_signal)는 조건이 미충족 상태에서 충족 상태로 전환되는 시점에만 전송됩니다. 조건을 계속 충족하는 동안에는 반복 전송되지 않으며, 조건을 이탈한 후 다시 충족하는 경우 다시 전송됩니다.
* 하나의 WebSocket 연결에서 활성화할 수 있는 구독 조합은 지표 값(indicator)과 조건 이벤트(indicator\_signal)를 합산하여 최대 100개입니다. 동일한 연결에서 새로운 구독을 요청하면 기존 구독 정보는 새로운 요청 내용으로 대체됩니다.
* 지원하지 않는 지표, 타임프레임, 파라미터 또는 조건 조합이 포함된 경우 요청 전체가 거부됩니다.
* 지원하지 않거나 존재하지 않는 마켓 코드를 요청한 경우에도 별도 오류가 발생하지 않으며, 해당 마켓의 데이터는 전송되지 않습니다.
* 조건 이벤트(indicator\_signal)에서 지원하는 조건의 범위는 지표 값(indicator)보다 제한적입니다. 필요한 조건이 지원 범위에 없는 경우 지표 값을 수신한 뒤 직접 조건 충족 여부를 판단할 수 있습니다.
* 캔들 마감 직전 짧은 순간에 종가가 결정되므로, 극단적인 변동이 있는 경우 실제 종가가 아닌 직전 가격으로 일시적으로 계산될 수 있습니다.
* 지표 값은 시장 상황에 따라 일시적으로 실제 종가를 반영하지 못할 수 있으며, 네트워크 및 시스템 환경 등에 따라 데이터 수신이 지연되거나 이용자별 수신 시점에 차이가 발생할 수 있습니다.
* 지원하는 마켓, 지표, 타임프레임 및 조건의 범위가 변경되거나 종료되는 경우 이용약관에서 정한 방법으로 공지 또는 통지할 예정입니다.거래지원 종료 또는 지원 범위 변경 등으로 구독 중인 데이터의 생성이 중단되는 경우, 별도의 WebSocket 알림 없이 해당 데이터의 전송이 중단될 수 있습니다.

***

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림