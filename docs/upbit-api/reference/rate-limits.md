---
updatedAt: 2026-11-25T15:47:54.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 요청 수 제한(Rate Limits)

업비트 API의 요청 수 제한(Rate Limits) 정책 안내 및 구현 가이드입니다.

### 기본 정책 안내

* 요청 수 제한은 초(Second) 단위로 적용됩니다.
* REST API는 같은 그룹에 속한 API의 요청 수를 함께 차감합니다. 각 API의 그룹과 정책은 본 문서 하단과 API Reference의 Rate Limit 영역에서 확인할 수 있습니다.
* 초당 최대 허용 요청 수는 공지 후 변경될 수 있습니다. 서비스 상황에 따라 추가 제한도 발생할 수 있습니다. 잔여 요청 수를 확인하여 과도한 요청을 보내지 않도록 유의해 주세요.
* Origin 헤더를 포함한 요청에는 별도 정책이 적용됩니다. 시세 조회(Quotation) REST API와 WebSocket 요청은 모두 10초당 1회만 허용됩니다. 자세한 내용은 [관련 공지](https://docs.upbit.com/kr/changelog/origin_rate_limit)를 확인해 주세요.
* 코인빌리기 권한은 메인포켓 API Key에만 설정할 수 있습니다. 권한 설정 시 이용 자격(고객 확인, 약관 동의 등)을 확인하며, 신청 테스트·신청·추가신청은 DAXA 온라인 교육을 이수해야 이용할 수 있습니다.

<br />

## REST API 요청 수 제한

시세 조회(Quotation)는 IP 단위, 거래·자산 관리(Exchange)는 포켓 단위로 적용됩니다. 같은 포켓의 여러 API Key는 한도를 공유하며, 서로 다른 포켓에는 독립된 한도가 적용됩니다.

<Accordion title="측정 단위 - 포켓(Pocket)" icon="fa-info-circle">
  ### 포켓 단위 — 독립 한도와 처리량 확장

  거래·자산 관리(Exchange) 요청은 포켓 단위로 차감됩니다. 같은 포켓에서 여러 API Key를 사용하더라도 요청 횟수는 합산하며, 서로 다른 포켓의 요청 한도는 독립적으로 적용됩니다.

  한 계정에는 메인포켓 1개가 제공되며, 서브포켓은 최대 5개까지 생성할 수 있습니다. 예를 들어 메인포켓 1개와 서브포켓 5개에서 각각 `default`(초당 30회)를 사용하면, 포켓별 한도를 합산한 계정 전체의 최대 처리량은 초당 180회입니다. 각 포켓의 개별 요청 한도는 초당 30회로 동일합니다.

  <style>
  .pocket-limit { margin: 20px 0; color: #344054; }
  .pocket-limit * { box-sizing: border-box; }
  .pocket-limit-summary { display: flex; justify-content: space-between; align-items: center; gap: 16px; padding: 16px 18px; border: 1px solid #e2e8f0; border-radius: 10px; margin-bottom: 20px; }
  .pocket-limit-label { font-size: 12px; color: #7b8794; margin-bottom: 5px; }
  .pocket-limit-summary-text { font-size: 14px; font-weight: 600; }
  .pocket-limit-total { font-size: 22px; font-weight: 700; text-align: right; white-space: nowrap; }
  .pocket-limit-summary-right { text-align: right; }
  .pocket-limit-heading { display: flex; justify-content: space-between; align-items: center; gap: 12px; margin-bottom: 12px; }
  .pocket-limit-heading strong { font-size: 14px; }
  .pocket-limit-heading span { font-size: 12px; color: #7b8794; }
  .pocket-limit-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 12px; }
  .pocket-limit-card { padding: 16px; border: 1px solid #e2e8f0; border-radius: 10px; background: #fff; min-width: 0; }
  .pocket-limit-card.main { background: #f8fafc; border-color: #ccd6e2; }
  .pocket-limit-name { font-size: 13px; font-weight: 600; margin-bottom: 12px; }
  .pocket-limit-count { display: flex; align-items: baseline; gap: 5px; margin-bottom: 10px; }
  .pocket-limit-count strong { font-size: 22px; line-height: 1.3; }
  .pocket-limit-count span { font-size: 12px; color: #7b8794; }
  .pocket-limit-track { height: 7px; border-radius: 6px; background: #e9eef4; overflow: hidden; margin-bottom: 10px; }
  .pocket-limit-fill { height: 100%; background: #8295ad; border-radius: 6px; }
  .pocket-limit-fill.fill-100 { width: 100%; }
  .pocket-limit-fill.fill-60 { width: 60%; }
  .pocket-limit-fill.fill-40 { width: 40%; }
  .pocket-limit-fill.fill-20 { width: 20%; }
  .pocket-limit-fill.fill-10 { width: 10%; }
  .pocket-limit-fill.fill-0 { width: 0%; }
  .pocket-limit-card.main .pocket-limit-fill { background: #566b83; }
  .pocket-limit-status { font-size: 12px; color: #66768a; }
  .pocket-limit-card.main .pocket-limit-status { color: #52677e; }
  .pocket-limit-note { font-size: 12px; color: #7b8794; line-height: 1.7; margin-top: 14px; }
  @media (max-width: 768px) { .pocket-limit-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
  @media (max-width: 480px) { .pocket-limit-grid { grid-template-columns: 1fr; } .pocket-limit-summary { align-items: flex-start; } .pocket-limit-total { font-size: 18px; } }
  </style>

  <div className="pocket-limit">
    <div className="pocket-limit-summary">
      <div>
        <div className="pocket-limit-label">포켓 구성</div>
        <div className="pocket-limit-summary-text">메인포켓 1개 + 서브포켓 최대 5개</div>
      </div>
      <div className="pocket-limit-summary-right">
        <div className="pocket-limit-label">계정 합산 최대 처리량</div>
        <div className="pocket-limit-total">180 req/s</div>
      </div>
    </div>
    <div className="pocket-limit-heading">
      <strong>포켓별 요청 한도</strong>
      <span>각 30 req/s</span>
    </div>
    <div className="pocket-limit-grid">
      <div className="pocket-limit-card main">
        <div className="pocket-limit-name">메인포켓</div>
        <div className="pocket-limit-count"><strong>30</strong><span>/ 30회</span></div>
        <div className="pocket-limit-track"><div className="pocket-limit-fill fill-100"></div></div>
        <div className="pocket-limit-status">한도 도달 · 추가 요청 시 HTTP 429</div>
      </div>
      <div className="pocket-limit-card">
        <div className="pocket-limit-name">서브포켓 1</div>
        <div className="pocket-limit-count"><strong>12</strong><span>/ 30회</span></div>
        <div className="pocket-limit-track"><div className="pocket-limit-fill fill-40"></div></div>
        <div className="pocket-limit-status">요청 가능</div>
      </div>
      <div className="pocket-limit-card">
        <div className="pocket-limit-name">서브포켓 2</div>
        <div className="pocket-limit-count"><strong>6</strong><span>/ 30회</span></div>
        <div className="pocket-limit-track"><div className="pocket-limit-fill fill-20"></div></div>
        <div className="pocket-limit-status">요청 가능</div>
      </div>
      <div className="pocket-limit-card">
        <div className="pocket-limit-name">서브포켓 3</div>
        <div className="pocket-limit-count"><strong>0</strong><span>/ 30회</span></div>
        <div className="pocket-limit-track"><div className="pocket-limit-fill fill-0"></div></div>
        <div className="pocket-limit-status">요청 가능</div>
      </div>
      <div className="pocket-limit-card">
        <div className="pocket-limit-name">서브포켓 4</div>
        <div className="pocket-limit-count"><strong>18</strong><span>/ 30회</span></div>
        <div className="pocket-limit-track"><div className="pocket-limit-fill fill-60"></div></div>
        <div className="pocket-limit-status">요청 가능</div>
      </div>
      <div className="pocket-limit-card">
        <div className="pocket-limit-name">서브포켓 5</div>
        <div className="pocket-limit-count"><strong>3</strong><span>/ 30회</span></div>
        <div className="pocket-limit-track"><div className="pocket-limit-fill fill-10"></div></div>
        <div className="pocket-limit-status">요청 가능</div>
      </div>
    </div>
    <div className="pocket-limit-note">※ 예시 사용량입니다. 180 req/s는 각 포켓의 최대 요청 한도를 합산한 값으로, 계정 전체에 별도로 부여된 요청 한도가 아닙니다.</div>
  </div>

  메인포켓이 요청 한도에 도달하더라도 서브포켓의 요청에는 영향을 주지 않습니다. 각 서브포켓은 남은 요청 횟수만큼 API를 호출할 수 있습니다.

  단, 코인빌리기 API는 메인포켓 API Key만 지원하므로 서브포켓으로 처리량을 늘릴 수 없습니다.
</Accordion>

### Quotation REST API

| Rate Limit 그룹 | 요청 한도     | 적용 단위 | 대상 API                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| ------------- | --------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `market`      | 초당 최대 10회 | IP    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-trading-pairs">페어 목록 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| `candle`      | 초당 최대 10회 | IP    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-seconds">초(Second) 캔들 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-minutes">분(Minute) 캔들 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-days">일(Day) 캔들 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-weeks">주(Week) 캔들 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-months">월(Month) 캔들 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-candles-years">연(Year) 캔들 조회</Anchor> |
| `trade`       | 초당 최대 10회 | IP    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/recent-trades-history">최근 체결 내역 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `ticker`      | 초당 최대 10회 | IP    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-tickers">페어 단위 현재가 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-quote-tickers">마켓 단위 현재가 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| `orderbook`   | 초당 최대 10회 | IP    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orderbooks">호가 정보 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orderbook-instruments">호가 정책 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |

### Exchange REST API

| Rate Limit 그룹      | 요청 한도     | 적용 단위 | 대상 API                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ------------------ | --------- | ----- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `default`          | 초당 최대 30회 | 포켓    | **포켓**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-pockets">포켓 정보 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-pocket-api-keys">포켓별 API Key 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-balance">포켓 잔고 조회</Anchor><br />[서브포켓 잔고 조회](https://docs.upbit.com/kr/reference/get-sub-pocket-balance)<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/universal-transfer">메인포켓 자산 이전</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-universal-transfers">메인포켓 자산 이전 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/transfer">서브포켓 자산 이전</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-transfers">서브포켓 자산 이전 목록 조회</Anchor> |
| `default`          | 초당 최대 30회 | 포켓    | **주문 조회 및 취소**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-order-information">주문 가능정보 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-order">개별 주문 취소</Anchor><br />[지정 주문 목록 취소](https://docs.upbit.com/kr/reference/cancel-orders-by-ids)<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-order">개별 주문 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-orders-by-ids">주문 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-open-orders">체결 대기 주문 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-closed-orders">종료 주문 조회</Anchor>                                                                                                                    |
| `default`          | 초당 최대 30회 | 포켓    | **출금**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/withdraw">디지털 자산 출금하기</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/withdraw-krw">원화 출금하기</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-withdrawal">디지털 자산 출금 취소 접수</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-withdrawal-information">출금 가능 정보 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-withdrawal-addresses">출금 허용 주소 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-withdrawal">개별 출금 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-withdrawals">출금 목록 조회</Anchor>                                                                         |
| `default`          | 초당 최대 30회 | 포켓    | **입금**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/create-deposit-address">입금 주소 생성 요청</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-deposit-address">개별 입금 주소 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-deposit-addresses">입금 주소 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-deposit-information">입금 가능 통화 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-deposit">개별 입금 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-deposits">입금 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/deposit-krw">원화 입금</Anchor>                                                                               |
| `default`          | 초당 최대 30회 | 포켓    | **트래블룰 및 서비스**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-travelrule-vasps">트래블룰 지원 거래소 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid">입금 UUID로 트래블룰 검증 요청</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/verify-travelrule-by-txid">입금 TXID로 트래블룰 검증 요청</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| `default`          | 초당 최대 30회 | 포켓    | **코인빌리기 조회**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-lending-information">코인빌리기 가능 정보 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending-request">코인빌리기 요청 상태 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-ongoing-lendings">진행 중 코인빌리기 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-closed-lendings">종료된 코인빌리기 목록 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-lending-transactions">코인빌리기 거래 내역 조회</Anchor>                                                                                                                                                                                                                                                    |
| `default`          | 초당 최대 30회 | 포켓    | **기타**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-service-status">통화별 입출금 서비스 상태 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-api-keys">API Key 목록 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `order`            | 초당 최대 12회 | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/new-order">주문 생성</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-and-new-order">취소 후 재주문</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| `order-test`       | 초당 최대 8회  | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/order-test">주문 생성 테스트</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `order-cancel-all` | 2초당 최대 1회 | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/batch-cancel-orders">주문 일괄 취소</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `lending`          | 초당 최대 10회 | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/new-lending">코인빌리기 신청</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/borrow-more">코인빌리기 추가 신청</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/adjust-collateral">담보 추가</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/repay-lending">코인빌리기 상환</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/repay-lending-with-collateral">담보로 상환</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor>                                                                                                                                                                                                                                 |
| `lending-test`     | 초당 최대 10회 | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-test">코인빌리기 신청 테스트</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |

<Callout icon="❗️" theme="error">
  ### **트래블룰 검증 요청 제한**

  입금 UUID 또는 TXID를 이용한 트래블룰 검증 요청은 동일 입금 건에 대해<span style="color: #D92D20;"> **10분당 최대 1회**</span> 요청할 수 있습니다.
</Callout>

<br />

## WebSocket 요청 수 제한

* WebSocket 연결 요청 (Handshake)은 초당 최대 5회로 제한됩니다. 인증하지 않은 연결은 IP 단위, 인증한 연결은 포켓 단위로 적용됩니다.
* 연결 이후 데이터 요청 메시지는 커넥션당 초당 최대 5회, 분당 최대 100회로 제한됩니다.
* 한도를 초과하면 연결 요청은 HTTP `429 Too Many Requests`로 거부되고, 데이터 요청 메시지는 `TOO_MANY_REQUEST` 에러 응답 후 연결이 종료됩니다. 자세한 내용은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 및 에러 안내</Anchor> 문서를 참고해 주세요.
* Private WebSocket(`/websocket/v1/private`, `/websocket/v1/info`)의 반복적인 연결과 해제를 지양하고, 기존 연결을 유지하여 사용하는 것을 권장합니다.

<br />

## 잔여 요청 수 확인 방법

REST API 응답에 포함된 `Remaining-Req` 헤더로 잔여 요청 수를 확인할 수 있습니다.

```http
Remaining-Req: group=default; min=1800; sec=29
```

| 필드      | 설명                                                                                                  |
| ------- | --------------------------------------------------------------------------------------------------- |
| `group` | 요청이 속한 Rate Limit 그룹                                                                                |
| `min`   | <span style="color: #D92D20;">**Deprecated** — 더 이상 사용되지 않는 필드입니다. 응답 처리 시 참조하지 않는 것을 권장합니다.</span> |
| `sec`   | 초당 한도 기준 잔여 요청 수. `0`이면 다음 초에 다시 요청해야 합니다.                                                          |

## 관련 문서

* <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/rest-api-best-practice">REST API Best Practice</Anchor> — `Remaining-Req` 활용 관리 방법과 Python 예제 코드(`update_from_header`)
* <Anchor target="_blank" href="https://docs.upbit.com/kr/changelog/origin_rate_limit">Origin 헤더 정책 공지</Anchor>
* <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/pocket-overview">포켓(Pocket) 안내</Anchor>
* <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/websocket-guide">WebSocket 사용 및 에러 안내</Anchor>

# 같은 수준의 페이지

* [개요](https://docs.upbit.com/kr/reference/api-overview.md)
* [인증](https://docs.upbit.com/kr/reference/auth.md)
* [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide.md)
* [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide.md)