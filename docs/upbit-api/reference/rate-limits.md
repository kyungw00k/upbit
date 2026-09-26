---
updatedAt: 2026-09-08T13:57:51.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 요청 수 제한(Rate Limits)

업비트 API의 요청 수 제한(Rate Limits) 정책 안내 및 구현 가이드입니다.

### 기본 정책 안내

* 요청 수 제한은 초(Second) 단위로 적용됩니다.
* 같은 그룹의 API는 요청 수를 함께 차감합니다. 각 API의 그룹과 정책은 본 문서 하단과 API Reference의 Rate Limit 영역에서 확인할 수 있습니다.
* 초당 최대 허용 요청 수는 공지 후 변경될 수 있습니다. 서비스 상황에 따라 추가 제한도 발생할 수 있습니다. 잔여 요청 수를 확인하여 과도한 요청을 보내지 않도록 유의해 주세요.
* Origin 헤더를 포함한 요청에는 별도 정책이 적용됩니다. 시세 조회(Quotation) REST API와 WebSocket 요청은 모두 10초당 1회만 허용됩니다. 자세한 내용은 [관련 공지](https://docs.upbit.com/kr/changelog/origin_rate_limit)를 확인해 주세요.

<br />

## 측정 단위

측정 단위는 기능 분류에 따라 달라집니다. 시세 조회는 IP, 거래·자산 관리는 포켓 단위로 측정됩니다. WebSocket은 연결 요청과 데이터 요청을 구분해 측정됩니다.

<style>
.mu-table th.unit-col, .mu-table td.unit-col {
  width: 150px;
  min-width: 150px;
  max-width: 150px;
  word-break: keep-all;
}
.mu-table th.cat-col, .mu-table td.cat-col {
  width: 200px;
  min-width: 200px;
  max-width: 200px;
  word-break: keep-all;
}
</style>

<table className="custom-table mu-table">
  <thead>
    <tr>
      <th className="cat-col">기능 분류</th>
      <th className="unit-col">측정 단위</th>
      <th>설명</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td className="cat-col"><b>시세 조회 REST API</b><br />(Quotation)</td>
      <td className="unit-col">IP</td>
      <td>한도는 IP 주소 단위로 적용됩니다. 같은 IP의 모든 요청은 하나의 한도를 공유합니다.</td>
    </tr>
    <tr>
      <td className="cat-col"><b>거래·자산 관리 REST API</b><br />(Exchange)</td>
      <td className="unit-col">포켓(Pocket)</td>
      <td>한도는 포켓 단위로 적용됩니다. 같은 포켓의 여러 API Key는 하나의 한도를 공유합니다.</td>
    </tr>
    <tr>
      <td className="cat-col"><b>WebSocket 연결 요청</b></td>
      <td className="unit-col">IP 또는 포켓</td>
      <td>한도는 연결 방식에 따라 다르게 적용됩니다. 인증 없이 연결하면 <b>IP</b>, 인증 정보를 포함해 연결하면 <b>포켓</b> 단위로 측정됩니다.</td>
    </tr>
    <tr>
      <td className="cat-col"><b>WebSocket 데이터 요청</b></td>
      <td className="unit-col">커넥션(Connection)</td>
      <td>한도는 WebSocket 연결 1개마다 적용됩니다. 연결을 여러 개 맺으면 각 연결이 별도의 한도를 가집니다.</td>
    </tr>
  </tbody>
</table>

<br />

### 포켓 단위 — 독립 한도와 처리량 확장

거래·자산 관리(Exchange) 요청은 포켓 단위로 차감됩니다. 같은 포켓의 여러 API Key는 하나의 한도를 공유하고, 서로 다른 포켓은 각자 독립된 한도를 가집니다. 포켓별 한도는 서로 간섭하지 않습니다. 따라서 운영 중인 포켓 수만큼 계정 전체의 처리량이 늘어납니다. 예를 들어 메인포켓 1개와 서브포켓 5개가 각각 `exchange.default`(초당 30회)를 사용하면, 계정 전체로는 초당 최대 180회까지 처리할 수 있습니다. 개별 포켓의 한도는 30회로 동일합니다.

<style>
.pk { padding:8px 0; }
.pk-cap { font-size:0.74rem; color:#868e96; margin-bottom:10px; }
.pk-cap b { color:#444950; }
.pk-grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:10px; }
.pk-card { border:1px solid #e0e0e0; border-radius:10px; padding:12px 14px; background:#fafbfc; }
.pk-card.main { grid-column:1 / -1; background:#f1f3f5; border-color:#ced4da; }
.pk-head { display:flex; align-items:center; justify-content:space-between; margin-bottom:8px; }
.pk-name { font-weight:700; font-size:0.84rem; color:#222; }
.pk-badge { font-size:0.68rem; font-weight:600; padding:2px 9px; border-radius:9px; background:#fff; color:#495057; border:1px solid #ced4da; }
.pk-gauge { height:8px; border-radius:5px; background:#e9ecef; overflow:hidden; margin-bottom:5px; }
.pk-gauge > span { display:block; height:100%; border-radius:5px; background:#adb5bd; }
.pk-num { font-size:0.74rem; color:#495057; }
.pk-num b { color:#222; }
.pk-keys { margin-top:7px; font-size:0.68rem; color:#868e96; }
.pk-keychip { font-family:monospace; background:#444950; color:#fff; padding:1px 7px; border-radius:5px; margin-right:4px; }
</style>

<div className="pk">
  <div className="pk-cap">계정 = <b>메인포켓 1개 + 서브포켓 5개</b> · 각 포켓이 <b>exchange.default 초당 30회</b>를 독립 보유 → 계정 합산 최대 <b>180 req/s</b></div>
  <div className="pk-grid">
    <div className="pk-card main">
      <div className="pk-head"><span className="pk-name">메인포켓</span><span className="pk-badge">차단됨 · HTTP 429</span></div>
      <div className="pk-gauge"><span style="width:100%"></span></div>
      <div className="pk-num"><b>30 / 30회</b> 사용 — 한도 도달</div>
      <div className="pk-keys"><span className="pk-keychip">Key A</span><span className="pk-keychip">Key B</span>↳ 같은 포켓의 Key는 한도 공유</div>
    </div>
    <div className="pk-card">
      <div className="pk-head"><span className="pk-name">서브포켓 1</span><span className="pk-badge">요청 가능</span></div>
      <div className="pk-gauge"><span style="width:40%"></span></div>
      <div className="pk-num"><b>12 / 30회</b> 사용</div>
    </div>
    <div className="pk-card">
      <div className="pk-head"><span className="pk-name">서브포켓 2</span><span className="pk-badge">요청 가능</span></div>
      <div className="pk-gauge"><span style="width:20%"></span></div>
      <div className="pk-num"><b>6 / 30회</b> 사용</div>
    </div>
    <div className="pk-card">
      <div className="pk-head"><span className="pk-name">서브포켓 3</span><span className="pk-badge">요청 가능</span></div>
      <div className="pk-gauge"><span style="width:0%"></span></div>
      <div className="pk-num"><b>0 / 30회</b> 사용</div>
    </div>
    <div className="pk-card">
      <div className="pk-head"><span className="pk-name">서브포켓 4</span><span className="pk-badge">요청 가능</span></div>
      <div className="pk-gauge"><span style="width:60%"></span></div>
      <div className="pk-num"><b>18 / 30회</b> 사용</div>
    </div>
    <div className="pk-card">
      <div className="pk-head"><span className="pk-name">서브포켓 5</span><span className="pk-badge">요청 가능</span></div>
      <div className="pk-gauge"><span style="width:10%"></span></div>
      <div className="pk-num"><b>3 / 30회</b> 사용</div>
    </div>
  </div>
</div>

메인포켓이 한도에 도달해 차단되어도 서브포켓 5개는 영향을 받지 않습니다. 각 서브포켓은 자신의 한도만큼 요청할 수 있습니다. 처리량을 늘리려면 API Key가 아니라 포켓을 분리해 발급해야 합니다.

<br />

### 커넥션 단위 — 커넥션별 독립 한도

WebSocket은 한 번 연결을 맺으면 그 연결을 통해 계속 메시지를 주고받습니다. 연결 이후의 데이터 요청은 각 연결마다 따로 한도가 적용됩니다. 따라서 처리량이 부족하면 연결을 나누어 늘릴 수 있습니다.

<style>
.cnx { padding:8px 0; }
.cnx-cap { font-size:0.74rem; color:#868e96; margin-bottom:10px; }
.cnx-cap b { color:#444950; }
.cnx-grid { display:grid; grid-template-columns:repeat(3, 1fr); gap:10px; }
.cnx-card { border:1px solid #e0e0e0; border-radius:10px; padding:12px 14px; background:#fafbfc; }
.cnx-head { display:flex; align-items:center; justify-content:space-between; margin-bottom:8px; }
.cnx-name { font-weight:700; font-size:0.8rem; color:#222; }
.cnx-badge { font-size:0.66rem; font-weight:600; padding:2px 8px; border-radius:9px; background:#fff; color:#495057; border:1px solid #ced4da; }
.cnx-gauge { height:8px; border-radius:5px; background:#e9ecef; overflow:hidden; margin-bottom:5px; }
.cnx-gauge > span { display:block; height:100%; border-radius:5px; background:#adb5bd; }
.cnx-num { font-size:0.74rem; color:#495057; }
.cnx-num b { color:#222; }
</style>

<div className="cnx">
  <div className="cnx-cap">동일 포켓에서 맺은 3개 커넥션 · 각 <b>커넥션당 초당 5회</b>를 독립 보유</div>
  <div className="cnx-grid">
    <div className="cnx-card">
      <div className="cnx-head"><span className="cnx-name">Connection #1</span><span className="cnx-badge">차단됨</span></div>
      <div className="cnx-gauge"><span style="width:100%"></span></div>
      <div className="cnx-num"><b>5 / 5회</b></div>
    </div>
    <div className="cnx-card">
      <div className="cnx-head"><span className="cnx-name">Connection #2</span><span className="cnx-badge">요청 가능</span></div>
      <div className="cnx-gauge"><span style="width:40%"></span></div>
      <div className="cnx-num"><b>2 / 5회</b></div>
    </div>
    <div className="cnx-card">
      <div className="cnx-head"><span className="cnx-name">Connection #3</span><span className="cnx-badge">요청 가능</span></div>
      <div className="cnx-gauge"><span style="width:0%"></span></div>
      <div className="cnx-num"><b>0 / 5회</b></div>
    </div>
  </div>
</div>

Connection #1이 한도(초당 5회 / 분당 100회)에 도달해 차단되어도 #2, #3은 영향을 받지 않습니다. 각 연결은 자신의 한도만큼 요청을 보낼 수 있습니다.

<br />

## Rate Limit 그룹별 정책

API는 Rate Limit 그룹 단위로 요청 수를 제한합니다. 같은 그룹에 속한 API는 요청 한도를 함께 차감합니다.

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
| `default`          | 초당 최대 30회 | 포켓    | **기타**<br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-service-status">통화별 입출금 서비스 상태 조회</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-api-keys">API Key 목록 조회</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `order`            | 초당 최대 12회 | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/new-order">주문 생성</Anchor><br /><Anchor target="_blank" href="https://docs.upbit.com/kr/reference/cancel-and-new-order">취소 후 재주문</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| `order-test`       | 초당 최대 8회  | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/order-test">주문 생성 테스트</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `order-cancel-all` | 2초당 최대 1회 | 포켓    | <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/batch-cancel-orders">주문 일괄 취소</Anchor>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |

<Callout icon="❗️" theme="error">
  ### **트래블룰 검증 요청 제한**

  입금 UUID 또는 TXID를 이용한 트래블룰 검증 요청은 동일 입금 건에 대해<span style="color: #D92D20;">**10분당 최대 1회**</span> 요청할 수 있습니다.
</Callout>

### WebSocket

| Rate Limit 그룹       | 요청 한도                    | 적용 단위                     | 대상 요청                   |
| ------------------- | ------------------------ | ------------------------- | ----------------------- |
| `websocket-connect` | 초당 최대 5회                 | 인증 미포함: IP<br />인증 포함: 포켓 | WebSocket 연결 요청         |
| `websocket-message` | 초당 최대 5회<br />분당 최대 100회 | 커넥션                       | WebSocket 데이터 요청 메시지 전송 |

> Rate Limit 그룹별 최대 요청 수는 서비스 정책에 따라 공지 후 변경되거나, 서비스 상황에 따라 추가 제한이 적용될 수 있습니다.

<br />

## 잔여 요청 수 확인 방법

REST API 응답의 Remaining-Req 헤더로 잔여 요청 수를 확인할 수 있습니다.

```http
Remaining-Req: group=default; min=1800; sec=29
```

| 필드      | 설명                                                                                                                               |
| ------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `group` | 요청이 속한 Rate Limit 그룹                                                                                                             |
| `min`   | <span style="color: #D92D20;"><strong> Deprecated</strong>— 분 단위 필드. <br />더 이상 사용되지 않는 필드 입니다. 응답 처리 시 참조하지 않는 것을 권장합니다.</span> |
| `sec`   | 현재 잔여 요청 수.<br />0이면 일정 시간 후에 다시 요청해야 합니다.                                                                                       |

<style>
.custom-table { width: 100%; table-layout: fixed; }
.custom-table th:first-child, .custom-table td:first-child { width: 200px; }
</style>

| HTTP 상태                            | 의미                               | 권장 조치                                 |
| ---------------------------------- | -------------------------------- | ------------------------------------- |
| <code>429 Too Many Requests</code> | 초당 한도 초과                         | 다음 초 경계까지 대기 후 재시도                    |
| <code>418 I'm a teapot</code>      | 429 누적으로 일시 차단 (동일 IP·포켓·커넥션 단위) | 응답에 포함된 차단 시간 정보를 확인 후, 안내된 시간 이후 재시도 |

> 정책 위반이 반복되면 차단 시간은 점진적으로 증가합니다.

<br />

## 관련 문서

* <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/rest-api-best-practice">REST API Best Practice</Anchor> — `Remaining-Req` 활용 관리 방법과 Python 예제 코드(`update_from_header`)
* <Anchor target="_blank" href="https://docs.upbit.com/kr/changelog/origin_rate_limit">Origin 헤더 정책 공지</Anchor>
* <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/pocket-overview">포켓(Pocket) 안내</Anchor>

# Sibling pages

* [개요](https://docs.upbit.com/kr/reference/api-overview.md)
* [인증](https://docs.upbit.com/kr/reference/auth.md)
* [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide.md)
* [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide.md)