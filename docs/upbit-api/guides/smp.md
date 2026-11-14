---
updatedAt: 2026-10-07T11:14:13.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 자전거래 체결 방지(Self-Match Prevention, SMP)

자전거래 체결 방지 옵션의 개념과 사용 방법을 안내합니다.

## 자전거래 체결 방지 기능이 무엇인가요?

자전거래 체결 방지(Self-Match Prevention, 이하 SMP)는 동일 회원의 서로 반대 방향 주문이 체결되는 자전거래를 방지하기 위한 기능입니다.

이를 통해 의도하지 않은 자전거래와 불필요한 수수료 발생을 방지하고, 관련 규제를 보다 효과적이고 안정적으로 준수할 수 있도록 지원합니다.

<br />

## 어떻게 설정하나요?

다음 API 엔드포인트의 SMP 파라미터를 통해 설정할 수 있습니다.

| **대상**                                                                                                                                              | **파라미터**       |
| --------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- |
| <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/new-order">주문 생성</Anchor> (POST /v1/orders)                                | `smp_type`     |
| <Anchor target="_blank" href="https://docs.upbit.com/kr/v1.6.4/reference/cancel-and-new-order">취소 후 재주문</Anchor> (POST /v1/orders/cancel\_and\_new) | `new_smp_type` |

<br />

## `smp_type`에는 어떤 값을 설정할 수 있나요?

SMP는 Taker 주문에 설정된 값을 기준으로 동작합니다.<br />따라서 호가창에 이미 남아 있는 Maker 주문에 설정된 SMP 값은 이후 체결 시 적용되지 않습니다.<br />`smp_type`에는 다음의 허용 값을 설정할 수 있습니다.

| `smp_type`허용 값 | 설명                                                                                                                                                                                          |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `cancel_maker` | `smp_type`**&#x20;**&#xC744; 지정하지 않으&#xBA74;**&#x20;**`cancel_maker`&#xAC00;**&#x20;**<span style="color: #D92D20;">**기본값으로 적용**</span>됩니다.<br />Taker 주문을 유지하고 Maker 주문을 취소하여 자전거래를 방지합니다. |
| `cancel_taker` | Maker 주문을 유지하고 Taker 주문을 취소하여 자전거래를 방지합니다.                                                                                                                                                  |
| `reduce`       | 자기 주문이 겹치는 수량만큼 Maker와 Taker 주문의 수량을 각각 차감하여 자전거래를 방지합니다. 차감 후 수량이 남아 있는 주문은 유지됩니다.                                                                                                         |
| `none`         | SMP를 적용하지 않습니다.<br />동일 회원의 주문 간 자전거래가 발생할 수 있습니다.                                                                                                                                          |

<Accordion title="SMP 추가 설명" icon="fa-info-circle">
  ![](https://files.readme.io/faab1112e7b27d040603acca01069b98fba4d8678d4094114ba6402165cf6be3-d775474a9c5185e82350c0bdc2d5eb96b0275b564d255fc76cb640d2c85b725d-kr_smp_head.png)

  ![](https://files.readme.io/09c455b169253a958bbe9d5ee4c40a07291d7301e0dc9a9c66db49759d2664e8-7c30a02fb8a977e6a25951296f182f7581bac36c502e49a1cc16cc4b5f6e41da-KR_2.png)
</Accordion>

<br />

## `post_only` 주문에서는 어떻게 동작하나요?

&#x20;`time_in_force=post_only` 주문에서는 다음과 같이 처리됩니다.

* `smp_type=none`은 함께 사용할 수 있습니다.&#x20;
* `smp_type=reduce`, `cancel_maker`, `cancel_taker`는 함께 사용할 수 없습니다.

<br />

## SMP 설정으로 인해 주문이 취소되었는지 어떻게 확인하나요?

SMP로 인해 주문이 취소되거나 수량이 차감된 경우 주문 관련 응답의 다음 필드를 통해 확인할 수 있습니다.

| 필드                 | 설명                                                              |
| ------------------ | --------------------------------------------------------------- |
| `smp_type`         | 주문에 적용된 SMP 값                                                   |
| `prevented_volume` | SMP로 인해 취소된 주문 수량                                               |
| `prevented_locked` | **매수 주문**의 경우 SMP로 인해 취소된 금액,<br />**매도 주문**의 경우 SMP로 인해 취소된 수량 |

WebSocket으로 [내 주문 및 체결(MyOrder)](https://docs.upbit.com/kr/reference/websocket-myorder) 데이터를 구독하는 경우 state 필드를 통해 SMP 적용 결과를 확인할 수 있습니다.

* `prevented`: SMP로 인해 주문 수량 또는 금액이 차감된 경우
* `cancel`: SMP 적용 결과 주문 수량 또는 금액이 더 이상 남아 있지 않아 주문이 취소된 경우

<br />

## 유의사항

* 자전거래 여부는 주문 시점이 아닌 체결 시점에 확인합니다.
* 시장가, 지정가, 예약가, IOC/FOK, 취소 후 재주문 등 지원되는 모든 주문 타입(`ord_type`)과 함께 사용할 수 있습니다.
* `prevented_locked`는 다음과 같이 제공됩니다.
  * 매수 주문: SMP로 인해 체결이 방지된 주문 금액이며 수수료를 포함합니다.<br />예를 들어 수수료가 0.25%인 경우 체결이 방지된 주문 금액에 1.0025를 곱한 값입니다.
  * 매도 주문: 수수료와 관계없이 SMP로 인해 체결이 방지된 주문 수량이며, `prevented_volume`과 동일합니다.
* 주문 관련 REST API 응답에서 `smp_type`, `prevented_volume`, `prevented_locked` 필드는 항상 제공됩니다.
  * `smp_type`은 주문에 적용된 값을 반환합니다.
  * 2026-11-09 이전에 생성된 주문에서 `smp_type`을 지정하지 않은 경우 `none`으로 반환됩니다.
* WebSocket으로 내 주문 및 체결(MyOrder) 데이터를 구독하는 경우에도 `smp_type`, `prevented_volume`, `prevented_locked` 필드가 응답에 포함됩니다.
* `“smp_type" : "none"`인 경우 `prevented_volume` 및 `prevented_locked`는 0으로 반환됩니다.

# Sibling pages

* [원화(KRW) 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/krw-market-info.md)
* [BTC 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/btc-market-info.md)
* [USDT 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/usdt-market-info.md)