---
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# [안내] 자전거래 체결 방지(Self-Match Prevention, SMP) 기본값 적용 안내

안녕하세요. 업비트 개발자센터입니다.

주문 API의 자전거래 체결 방지(Self-Match Prevention, 이하 SMP) 기능에 기본값이 추가될 예정입니다.

SMP는 동일 계정에서 서로 반대 방향의 주문이 체결되는 자전거래를 방지하기 위한 기능입니다.

**업데이트 이후 일반 주문에서&#x20;**`smp_type`**을 별도로 지정하지 않으면&#x20;**`cancel_maker`**가 기본값으로 적용됩니다.&#x20;**&#xAE30;존 요청 형식은 그대로 사용할 수 있으나, `smp_type`을 지정하지 않은 주문은 기존과 주문 처리 방식이 달라질 수 있으므로 아래 변경 사항을 확인해 주세요.

## 1.적용 일시

* 적용 예정일 : 2026-11-09 (월)

## 2.주요 업데이트 내용

* `smp_type`을 지정하지 않은 경우 기본값으로 `cancel_maker`가 적용됩니다.
* `smp_type`에 `none`이 추가됩니다. `“smp_type"="none"`으로 요청하면 SMP가 적용되지 않으며, 동일 계정의 주문 간 체결이 발생할 수 있습니다.
* `"time_in_force"="post_only"` 주문에서 `smp_type`을 지정하지 않은 경우는 이번 변경 대상이 아니므로, `none`과 동일하게 처리됩니다.
* `"time_in_force"="post_only"`와 `“smp_type"="reduce"`, `“cancel_maker"`, `“cancel_taker"`는 함께 지정할 수 없습니다.
* 주문 생성 및 조회 관련 응답에서 `smp_type` 필드가 항상 제공되도록 변경됩니다.
* 업데이트 된 문서는 아래와 같습니다.
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/docs/smp">자전거래 체결 방지(Self-Match Prevention, SMP)</Anchor>
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/rate-limits">REST API 공통 안내</Anchor>
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/new-order">주문 생성하기</Anchor> ㅣ <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/order-test">주문 생성 테스트하기</Anchor>
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/get-order">주문 조회</Anchor> ㅣ <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/list-orders-by-ids">주문 목록 조회하기</Anchor> ㅣ <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/list-open-orders">체결 대기 주문 목록 조회</Anchor> ㅣ <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/list-closed-orders">종료 주문 목록 조회</Anchor>
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/cancel-and-new-order">취소 후 재주문하기</Anchor>
  * <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/websocket-myorder">내 주문 및 체결(MyOrder)</Anchor>

## 3.유의사항

* 일반 주문에서 `smp_type`을 지정하지 않는 경우의 기본 동작이 변경됩니다. 현재 주문 요청 및 처리 로직을 사전에 확인해 주세요.
* 기존과 동일하게 SMP를 적용하지 않으려면 `smp_type=none`을 명시적으로 지정해 주세요.
* 적용 일정은 내부 사정에 따라 변경될 수 있으며, 변경 시 본 공지를 통해 안내드리겠습니다.

***

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림

<br />