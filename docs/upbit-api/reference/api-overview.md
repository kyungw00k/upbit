---
updatedAt: 2026-11-26T04:31:22.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 개요

업비트 API의 주요 기능과 연동 순서를 한눈에 확인하세요.

## 업비트 API 기능 개요

업비트 API는 제공하는 정보의 성격에 따라 시세 조회(Quotation)와 거래 및 자산 관리(Exchange)의 두 가지 카테고리로 분류됩니다.

* **시세 조회 기능**: 업비트 거래소에서 지원하는 페어별 시세 정보(실시간 정보 및 과거 이력)를 조회할 수 있습니다.
* **거래 및 자산 관리 기능**: 개인 업비트 계정과 연동하여 주문, 입출금, 자산 관리 및 코인빌리기 등의 기능을 이용할 수 있습니다.

각 API 카테고리별 주요 지원 기능과 범위는 다음과 같습니다.

| 구분            | 시세 조회(Quotation)                 | 거래 및 자산 관리(Exchange)             |
| ------------- | -------------------------------- | -------------------------------- |
| **세부 기능**     | 페어, 캔들(OHLCV), 체결 이력, 현재가, 호가 조회 | 계정 자산 조회, 주문 관리, 입출금 관리, 코인빌리기 등 |
| **API 인증**    | Public API로, 인증 없이 조회 가능         | Private API로, API Key를 사용한 인증 필수 |
| **API 동작 범위** | 조회만 지원(과거 이력 및 실시간 조회 포함)        | 요청 생성 및 실행, 취소, 조회 지원            |

## 업비트 API 연동 방식

업비트 API는 REST API와 WebSocket 방식의 연동을 모두 지원합니다. 두 방식의 특징을 비교하고, 프로그램 구현 환경과 이용 목적에 적합한 방식을 선택해 주세요.

| 구분        | REST API                                    | WebSocket                                      |
| --------- | ------------------------------------------- | ---------------------------------------------- |
| **통신 방식** | 요청(Request)에 대해 응답(Response)을 반환하는 방식       | 연결을 유지하면서 실시간 데이터를 지속적으로 수신하는 방식               |
| **장점**    | 구현과 테스트가 비교적 간단하며, 필요한 시점에 데이터를 요청할 수 있음    | 낮은 지연 시간으로 실시간 데이터를 수신할 수 있으며, 지속적인 데이터 갱신에 적합 |
| **단점**    | 데이터를 갱신하려면 반복 요청이 필요하며, 요청 수 제한의 영향을 받음     | 연결 유지, 재연결 및 Ping/Pong 처리 등 별도의 연결 관리가 필요      |
| **추천 용도** | 주문 생성 및 취소, 입출금 요청, 코인빌리기 신청 및 상환, 계정 정보 조회 | 실시간 시세·체결·호가·주문 정보 구독 및 모니터링                   |

## 업비트 REST API 및 WebSocket 기능 목록

현재 업비트에서 제공하는 주요 REST API 및 WebSocket 기능은 다음과 같습니다.

<Table>
  <thead>
    <tr>
      <th>
        분류
      </th>

      <th>
        주요 기능
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
        **Quotation** 페어(Trading Pairs)
      </td>

      <td>
        지원 페어 목록 조회
      </td>

      <td>
        [페어 목록 조회](https://docs.upbit.com/kr/reference/list-trading-pairs)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Quotation** 캔들(OHLCV)
      </td>

      <td>
        기간별 OHLCV 조회
      </td>

      <td>
        [초 캔들](https://docs.upbit.com/kr/reference/list-candles-seconds)

        [분 캔들](https://docs.upbit.com/kr/reference/list-candles-minutes)

        [일 캔들](https://docs.upbit.com/kr/reference/list-candles-days)

        [주 캔들](https://docs.upbit.com/kr/reference/list-candles-weeks)

        [월 캔들](https://docs.upbit.com/kr/reference/list-candles-months)

        [연 캔들](https://docs.upbit.com/kr/reference/list-candles-years)
      </td>

      <td>
        [실시간 캔들 수신](https://docs.upbit.com/kr/reference/websocket-candle)
      </td>
    </tr>

    <tr>
      <td>
        **Quotation** 체결(Trade)
      </td>

      <td>
        체결 정보 조회
      </td>

      <td>
        [최근 체결 내역](https://docs.upbit.com/kr/reference/recent-trades-history)
      </td>

      <td>
        [실시간 체결 수신](https://docs.upbit.com/kr/reference/websocket-trade)
      </td>
    </tr>

    <tr>
      <td>
        **Quotation** 현재가(Ticker)
      </td>

      <td>
        현재가 및 변동 정보 조회
      </td>

      <td>
        [페어별 현재가](https://docs.upbit.com/kr/reference/list-tickers)

        [마켓별 현재가](https://docs.upbit.com/kr/reference/list-quote-tickers)
      </td>

      <td>
        [실시간 현재가 수신](https://docs.upbit.com/kr/reference/websocket-ticker)
      </td>
    </tr>

    <tr>
      <td>
        **Quotation** 호가(Orderbook)
      </td>

      <td>
        호가 및 잔량 조회
      </td>

      <td>
        [호가 정보](https://docs.upbit.com/kr/reference/list-orderbooks)

        [호가 정책](https://docs.upbit.com/kr/reference/list-orderbook-instruments)
      </td>

      <td>
        [실시간 호가 수신](https://docs.upbit.com/kr/reference/websocket-orderbook)
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 포켓(Pocket)
      </td>

      <td>
        포켓 및 자산 이전 관리
      </td>

      <td>
        [포켓 정보](https://docs.upbit.com/kr/reference/list-pockets)

        [포켓별 API Key](https://docs.upbit.com/kr/reference/list-pocket-api-keys)

        [서브포켓 잔고](https://docs.upbit.com/kr/reference/get-sub-pocket-balance)

        [메인포켓 자산 이전](https://docs.upbit.com/kr/reference/universal-transfer)

        [메인포켓 이전 목록](https://docs.upbit.com/kr/reference/list-universal-transfers)

        [서브포켓 자산 이전](https://docs.upbit.com/kr/reference/transfer)

        [서브포켓 이전 목록](https://docs.upbit.com/kr/reference/list-transfers)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 자산(Asset)
      </td>

      <td>
        포켓 잔고 조회
      </td>

      <td>
        [포켓 잔고](https://docs.upbit.com/kr/reference/get-balance)
      </td>

      <td>
        [실시간 잔고 변동](https://docs.upbit.com/kr/reference/websocket-myasset)
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 주문(Order)
      </td>

      <td>
        주문 생성 및 취소
      </td>

      <td>
        [주문 생성](https://docs.upbit.com/kr/reference/new-order)

        [주문 생성 테스트](https://docs.upbit.com/kr/reference/order-test)

        [개별 주문 취소](https://docs.upbit.com/kr/reference/cancel-order)

        [지정 주문 목록 취소](https://docs.upbit.com/kr/reference/cancel-orders-by-ids)

        [주문 일괄 취소](https://docs.upbit.com/kr/reference/batch-cancel-orders)

        [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 주문(Order)
      </td>

      <td>
        주문 조회
      </td>

      <td>
        [주문 가능 정보](https://docs.upbit.com/kr/reference/available-order-information)

        [개별 주문](https://docs.upbit.com/kr/reference/get-order)

        [주문 목록](https://docs.upbit.com/kr/reference/list-orders-by-ids)

        [체결 대기 주문](https://docs.upbit.com/kr/reference/list-open-orders)

        [종료 주문](https://docs.upbit.com/kr/reference/list-closed-orders)​
      </td>

      <td>
        [실시간 주문·체결 정보](https://docs.upbit.com/kr/reference/websocket-myorder)
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 출금(Withdrawal)
      </td>

      <td>
        출금 요청 및 취소
      </td>

      <td>
        [디지털 자산 출금](https://docs.upbit.com/kr/reference/withdraw)

        [원화 출금](https://docs.upbit.com/kr/reference/withdraw-krw), [출금 취소](https://docs.upbit.com/kr/reference/cancel-withdrawal)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 출금(Withdrawal)
      </td>

      <td>
        출금 조회
      </td>

      <td>
        [출금 가능 정보](https://docs.upbit.com/kr/reference/available-withdrawal-information)

        [출금 허용 주소](https://docs.upbit.com/kr/reference/list-withdrawal-addresses)

        [개별 출금](https://docs.upbit.com/kr/reference/get-withdrawal)

        [출금 목록](https://docs.upbit.com/kr/reference/list-withdrawals)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 입금(Deposit)
      </td>

      <td>
        입금 주소 관리
      </td>

      <td>
        [입금 주소 생성](https://docs.upbit.com/kr/reference/create-deposit-address)

        [개별 입금 주소](https://docs.upbit.com/kr/reference/get-deposit-address)

        [입금 주소 목록](https://docs.upbit.com/kr/reference/list-deposit-addresses)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 입금(Deposit)
      </td>

      <td>
        디지털 자산 입금 조회
      </td>

      <td>
        [입금 가능 통화](https://docs.upbit.com/kr/reference/available-deposit-information)

        [개별 입금](https://docs.upbit.com/kr/reference/get-deposit)

        [입금 목록](https://docs.upbit.com/kr/reference/list-deposits)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 입금(Deposit)
      </td>

      <td>
        트래블룰 관리
      </td>

      <td>
        [지원 거래소 목록](https://docs.upbit.com/kr/reference/list-travelrule-vasps)

        [UUID 검증](https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid)

        [TXID 검증](https://docs.upbit.com/kr/reference/verify-travelrule-by-txid)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 입금(Deposit)
      </td>

      <td>
        원화 입금
      </td>

      <td>
        [원화 입금](https://docs.upbit.com/kr/reference/deposit-krw)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 서비스 정보(Service)
      </td>

      <td>
        서비스 상태 및 API Key 조회
      </td>

      <td>
        [입출금 서비스 상태](https://docs.upbit.com/kr/reference/get-service-status)

        [API Key 목록](https://docs.upbit.com/kr/reference/list-api-keys)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 코인빌리기(Lending)
      </td>

      <td>
        이용 가능 정보 및 신청 조건 확인
      </td>

      <td>
        [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information)

        [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 코인빌리기(Lending)
      </td>

      <td>
        신청 및 담보 관리
      </td>

      <td>
        [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending)

        [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more)

        [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 코인빌리기(Lending)
      </td>

      <td>
        상환
      </td>

      <td>
        [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending)

        [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Exchange** 코인빌리기(Lending)
      </td>

      <td>
        요청 상태 및 이용 내역 조회
      </td>

      <td>
        [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)

        [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)

        [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings)

        [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings)

        [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions)
      </td>

      <td>
        \-
      </td>
    </tr>

    <tr>
      <td>
        **Service** 공지사항(Announcement)
      </td>

      <td>
        공지 게시·변경 정보 수신
      </td>

      <td>
        \-
      </td>

      <td>
        [실시간 공지사항 수신](https://docs.upbit.com/kr/reference/websocket-announcement)
      </td>
    </tr>

    <tr>
      <td>
        **Service** 기술적 지표
      </td>

      <td>
        지표 값 및 조건 이벤트 수신
      </td>

      <td>
        \-
      </td>

      <td>
        [지표 값 수신](https://docs.upbit.com/kr/reference/websocket-indicator)

        [지표 조건 이벤트 수신](https://docs.upbit.com/kr/reference/websocket-indicator-signal)
      </td>
    </tr>
  </tbody>
</Table>

# 같은 수준의 페이지

* [인증](https://docs.upbit.com/kr/reference/auth.md)
* [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits.md)
* [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide.md)
* [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide.md)