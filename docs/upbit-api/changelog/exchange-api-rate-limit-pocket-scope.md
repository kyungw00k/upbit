Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내] 거래소 API Rate Limit 적용 단위 변경 (계정 → 포켓)

안녕하세요. 업비트 개발자센터입니다.

거래·자산 관리(Exchange) API의 Rate Limit 측정 단위가 **계정 단위**에서 **포켓(Pocket) 단위**로 변경됩니다. 변경 후에는 각 포켓이 독립된 한도를 가집니다. 한 포켓의 사용량은 다른 포켓에 영향을 주지 않습니다.

이에 따라 운영 중인 포켓 수만큼 계정 전체의 처리량이 늘어납니다. 포켓 간 간섭이 사라지므로 용도별로 한도를 분리해 운영할 수 있습니다.

## 1. 적용 일시

* 2026-06-25 (목)

## 2. 변경 내용

| 항목               | 변경 전            | 변경 후           |
| ---------------- | --------------- | -------------- |
| Rate Limit 측정 단위 | 계정 단위 (Account) | 포켓 단위 (Pocket) |
| 그룹별 한도           | 변경 없음           | 변경 없음          |

> 그룹별 초당 최대 요청 수(예: `exchange.default` 그룹은 초당 30회)는 그대로입니다. 한도를 차감하는 기준만 계정에서 포켓으로 바뀝니다.

## 3. 기대 효과

* **포켓별 독립 한도**: 각 포켓이 독립된 한도를 가집니다. 서브포켓을 추가로 생성하면 그만큼 더 많은 요청을 처리할 수 있습니다. 전체 처리량은 포켓 수에 비례해 늘어납니다.
* **포켓 간 간섭 제거**: 한 포켓이 한도를 초과해도 다른 포켓은 영향을 받지 않습니다. 일시 차단(HTTP 418)도 포켓 단위로만 적용됩니다.

## 4. 영향 범위

| 기능 분류                                   | 변경 전 | 변경 후   | 비고   |
| --------------------------------------- | ---- | ------ | ---- |
| 거래·자산 관리(Exchange) REST API             | 계정   | **포켓** | ✅ 변경 |
| 인증 포함 WebSocket 연결(`websocket-connect`) | 계정   | **포켓** | ✅ 변경 |
| 시세 조회(Quotation) REST API               | IP   | IP     | 유지   |
| 인증 미포함 WebSocket 연결                     | IP   | IP     | 유지   |
| WebSocket 데이터 요청(`websocket-message`)   | 커넥션  | 커넥션    | 유지   |

> ✅ 표시된 항목만 측정 단위가 **계정 → 포켓**으로 변경됩니다. 변경 대상 그룹은 `exchange.default` / `exchange.order` / `exchange.order-test` / `exchange.order-cancel-all`입니다.

## 5. 유의사항

* 응답 헤더 `Remaining-Req`의 `sec` 값으로 한도를 관리하시는 경우, 별도의 코드 변경 없이 정상 동작합니다.
* 자체 카운터로 잔여 요청 수를 관리하시는 경우, 기준을 계정이 아니라 **포켓 단위**로 분리해 주시기 바랍니다.
* 같은 포켓의 여러 API Key는 하나의 한도를 공유합니다. 처리량을 늘리시려면 API Key를 추가 발급하는 대신 **서브포켓을 분리해 운영**해 주시기 바랍니다.
* 각 API의 상세 측정 단위와 한도는 [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 문서를 확인해 주시기 바랍니다.

***

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림