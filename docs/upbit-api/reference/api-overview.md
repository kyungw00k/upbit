---
updatedAt: 2026-08-31T02:33:46.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 개요

업비트 API의 주요 기능과 연동 순서를 한눈에 확인하세요.

## 업비트 API 기능 개요

업비트 API는 제공하는 정보의 성격에 따라 시세 조회(Quotation)와 거래 및 자산 관리(Exchange)의 두 가지 카테고리로 분류됩니다.

* **시세 조회 기능**: 업비트 거래소에서 지원하는 모든 페어별 시세 정보(실시간 정보와 이력)를 조회할 수 있습니다.
* **거래 및 자산 관리 기능**: 개인 업비트 계정과 연동하여 주문, 입출금, 자산 관리를 실행할 수 있습니다.

각 API 카테고리별 주요 지원 기능 및 범위는 아래와 같습니다.

|                 | 시세 조회(Quotation)                 | 거래 및 자산 관리(Exchange)                 |
| --------------- | -------------------------------- | ------------------------------------ |
| **세부 기능**       | 페어, 캔들(OHLCV), 체결 이력, 현재가, 호가 조회 | 계정 자산 조회, 주문 관리, 입출금 관리 등            |
| **Open API 권한** | Public API로, **인증 없이 조회 가능**     | Private API로, **API Key를 사용한 인증 필수** |
| **API 동작 범위**   | 조회만 지원 (과거 이력 및 실시간 조회 포함)       | 요청 생성(실행), 취소, 조회 지원                 |

<br />

## 업비트 API 연동 방식

업비트 API는 REST API 방식과 WebSocket 방식의 연동을 모두 지원합니다. 아래 두 프로토콜을 비교한 표를 참고하여 프로그램 구현 환경 및 용도에 따라 연동하시기 바랍니다.

<Table>
  <thead>
    <tr>
      <th>

      </th>

      <th>
        REST API
      </th>

      <th>
        WebSocket
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        **통신 방식**
      </td>

      <td>
        요청(Request)-응답(Response)으로 동작하며 필요 시점에 요청하는 방식
      </td>

      <td>
        최초 연결 이후 서버와의 지속적인 통신을 통해 실시간 데이터를 수신하는 스트림 방식
      </td>
    </tr>

    <tr>
      <td>
        **장점**
      </td>

      <td>
        - 구현 및 테스트가 직관적이고 쉬움
        - HTTP 기반으로 서버 환경 연동 지원
        - 요청 시점에 명확한 데이터 확보 가능
      </td>

      <td>
        - 빠른 데이터 반영과 낮은 지연 시간으로 시세 등 실시간 데이터 수신에 최적화
        - 압축 데이터 형식 지원 등을 활용하여 트래픽 최소화 가능
      </td>
    </tr>

    <tr>
      <td>
        **단점**
      </td>

      <td>
        - 요청에 의해 응답하므로 실시간성 낮음
        - 정보가 필요한 시점마다 매번 새로 요청 필요
      </td>

      <td>
        - 비교적 높은 구현 난이도: ping/pong, 연결 유지 및 재연결 등 관리 필요
        - 서버와의 연결 유지에 지속적인 리소스 할당 필요
      </td>
    </tr>

    <tr>
      <td>
        **추천 용도**
      </td>

      <td>
        주문 생성 및 취소, 입출금 요청 등 조회를 제외한 실행 작업, 비교적 긴 주기의 정보 갱신을 위한 조회 시
      </td>

      <td>
        실시간 시세/체결 데이터 구독을 통한 자동 매매 전략 반영 및 모니터링
      </td>
    </tr>
  </tbody>
</Table>

<br />

## 업비트 REST API 및 WebSocket 기능 목록

현재 업비트가 지원하는 전체 REST API 및 WebSocket 기능 목록은 아래와 같습니다.

| 분류                                    | 주요 기능               | REST API                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | WebSocket                                                                                                      |
| ------------------------------------- | ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| **Quotation**<br />페어 (Trading Pairs) | 지원 페어 목록 조회         | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-trading-pairs">페어 목록 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         | -                                                                                                              |
| **Quotation**<br />캔들 (OHLCV)         | 기간별 OHLCV 조회        | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-seconds">초 캔들</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-minutes">분 캔들</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-days">일 캔들</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-weeks">주 캔들</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-months">월 캔들</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-years">연 캔들</Anchor>                                                                                            | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-candle">실시간 캔들 수신</Anchor>         |
| **Quotation**<br />체결 (Trade)         | 체결 정보 조회            | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/recent-trades-history">최근 체결 내역</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-trade">실시간 체결 수신</Anchor>          |
| **Quotation**<br />현재가 (Ticker)       | 현재가 및 변동 정보 조회      | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-tickers">페어별 현재가</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-quote-tickers">마켓별 현재가</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-ticker">실시간 현재가 수신</Anchor>        |
| **Quotation**<br />호가 (Orderbook)     | 호가 및 잔량 조회          | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orderbooks">호가 정보</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orderbook-instruments">호가 정책</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-orderbook">실시간 호가 수신</Anchor>      |
| **Exchange**<br />포켓 (Pocket)         | 포켓 및 자산 이전 관리       | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-pockets">포켓 정보</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-pocket-api-keys">포켓별 API Key</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-sub-pocket-balance">서브포켓 잔고</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/universal-transfer">메인포켓 자산 이전</Anchor><br />[메인포켓 이전 목록](https://docs.upbit.com/kr/reference/list-universal-transfers)<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/transfer">서브포켓 자산 이전</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-transfers">서브포켓 이전 목록</Anchor> | -                                                                                                              |
| **Exchange**<br />자산 (Asset)          | 포켓 잔고 조회            | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-balance">포켓 잔고</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-myasset">실시간 잔고 변동</Anchor>        |
| **Exchange**<br />주문 (Order)          | 주문 생성 및 취소          | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/new-order">주문 생성</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/order-test">주문 생성 테스트</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-order">개별 주문 취소</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-orders-by-ids">지정 주문 목록 취소</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/batch-cancel-orders">주문 일괄 취소</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-and-new-order">취소 후 재주문</Anchor>                                                                                         | -                                                                                                              |
| **Exchange**<br />주문 (Order)          | 주문 조회               | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-order-information">주문 가능 정보</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-order">개별 주문</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orders-by-ids">주문 목록</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-open-orders">체결 대기 주문</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-closed-orders">종료 주문</Anchor>                                                                                                                                                                                                | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-myorder">실시간 주문·체결 정보</Anchor>     |
| **Exchange**<br />출금 (Withdrawal)     | 출금 요청 및 취소          | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/withdraw">디지털 자산 출금</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/withdraw-krw">원화 출금</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-withdrawal">출금 취소</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                     | -                                                                                                              |
| **Exchange**<br />출금 (Withdrawal)     | 출금 조회               | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-withdrawal-information">출금 가능 정보</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-withdrawal-addresses">출금 허용 주소</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-withdrawal">개별 출금</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-withdrawals">출금 목록</Anchor>                                                                                                                                                                                                                                                                                         | -                                                                                                              |
| **Exchange**<br />입금 (Deposit)        | 입금 주소 관리            | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/create-deposit-address">입금 주소 생성</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-deposit-address">개별 입금 주소</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-deposit-addresses">입금 주소 목록</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                      | -                                                                                                              |
| **Exchange**<br />입금 (Deposit)        | 디지털 자산 입금 조회        | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-deposit-information">입금 가능 통화</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-deposit">개별 입금</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-deposits">입금 목록</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                      | -                                                                                                              |
| **Exchange**<br />입금 (Deposit)        | 트래블룰 관리             | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-travelrule-vasps">지원 거래소 목록</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid">UUID 검증</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/verify-travelrule-by-txid">TXID 검증</Anchor>                                                                                                                                                                                                                                                                                                                                                                                               | -                                                                                                              |
| **Exchange**<br />입금 (Deposit)        | 원화 입금               | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/deposit-krw">원화 입금</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | -                                                                                                              |
| **Exchange**<br />서비스 정보 (Service)    | 서비스 상태 및 API Key 조회 | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-service-status">입출금 서비스 상태</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-api-keys">API Key 목록</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | -                                                                                                              |
| **Service**<br />공지사항 (Announcement)  | 공지 게시·변경 정보 수신      | -                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-announcement">실시간 공지사항 수신</Anchor> |

# Sibling pages

* [인증](https://docs.upbit.com/kr/reference/auth.md)
* [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits.md)
* [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide.md)
* [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide.md)