---
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# [업데이트] Upbit SDK v1.1.0 - WebSocket 데이터 타입 지원 확대

안녕하세요. 업비트 개발자센터입니다.

업비트 Python 및 TypeScript SDK가 업데이트됩니다.

이번 업데이트를 통해 기존에 제공 중인 공지사항 및 기술적 지표 WebSocket 기능을 Python, TypeScript SDK에서도 사용할 수 있도록 지원 범위를 확대합니다.

## 1.적용 일시

* 적용일: 2026-10-20 (화)

## 2.주요 업데이트 내용

### 업데이트 버전

* [Python SDK](https://github.com/upbit-official/upbit-sdk-python): `v1.1.0`
* [TypeScript SDK](https://github.com/upbit-official/upbit-sdk-typescript): `v1.1.0`

신규 WebSocket 데이터 타입을 사용하려면 해당 버전 이상으로 업데이트해 주세요.

### WebSocket 데이터 타입 지원 추가

Python 및 TypeScript SDK에서 다음 3가지 Private WebSocket 데이터 타입을 지원합니다.

* `announcement`: 신규 공지사항 게시 및 기존 공지사항의 수정 정보를 실시간으로 수신
* `indicator`: 선택한 캔들 주기의 마감 시점마다 계산된 기술적 지표 값을 수신
* `indicator_signal`: 설정한 기술적 지표 조건이 충족되는 시점에 이벤트를 수신

### 공지사항 (`announcement`)

* 신규 공지사항 게시(`CREATED`) 및 기존 공지사항 수정(`UPDATED`) 이벤트를 수신할 수 있습니다.
* 공지사항 카테고리를 선택하여 구독하거나 공지 본문을 함께 수신할 수 있습니다.

### 기술적 지표 (`indicator`, `indicator_signal`)

* 거래 페어 및 캔들 주기를 지정하여 거래량, RSI, 이동평균선, MACD 등 6종의 기술적 지표를 구독할 수 있습니다.
* `indicator`를 통해 캔들 마감 시점의 기술적 지표 값을 수신할 수 있습니다.
* `indicator_signal`을 통해 설정한 지표 조건이 충족되는 시점에 이벤트를 수신할 수 있습니다.

## 3.이용 안내

이번에 추가된 데이터 타입은 인증이 필요한 Private WebSocket을 통해 제공됩니다.

기술적 지표 WebSocket은 KRW 마켓에서만 지원하며, 공지사항 및 기술적 지표 데이터는 스냅샷 없이 실시간 이벤트만 제공합니다.

자세한 사용 방법은 아래 문서를 참고해 주세요.

* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk)
* [공지사항 WebSocket](https://docs.upbit.com/kr/reference/websocket-announcement)
* [기술적 지표 WebSocket](https://docs.upbit.com/kr/reference/websocket-indicator)

***

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림