---
updatedAt: 2026-08-20T10:03:37.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 자전거래 체결 방지(Self-Match Prevention, SMP)

자전거래 체결 방지 옵션의 개념과 사용 방법을 안내합니다.

## 자전거래 체결 방지 기능이 무엇인가요?

의도치 않게 동일 회원의 주문이 호가창에서 만날 때 주문 체결을 방지하여 **불필요한 수수료 발생을 줄이고**, **보다 효과적이고 안정적으로 관련 규제를 준수**할 수 있도록 지원하는 기능입니다.

## 어떻게 설정하나요?

* 다음의 API 엔드포인트와 파라미터에서 지원합니다:

|                              대상                             |      파라미터      |
| :---------------------------------------------------------: | :------------: |
|               [POST/v1/orders](https://docs.upbit.com/kr/reference/new-order)               |   `smp_type`   |
| [POST/v1/orders/cancel\_and\_new](https://docs.upbit.com/kr/reference/cancel-and-new-order) | `new_smp_type` |

## 설정가능한 모드는 무엇이 있나요?

다음의 3가지 타입으로 설정하실 수 있으며 **<span style="color:red">taker 주문의 설정 기준으로 동작합니다. 때문에 호가창에 이미 남겨진 주문의 SMP 모드는 더 이상 유효하지 않습니다. </span>**

* `cancel_taker`: maker 주문을 유지하고, **taker 주문을 취소**하여 자전거래를 방지합니다.
* `cancel_maker`: taker 주문을 유지하고, **maker 주문을 취소**하여 자전거래를 방지합니다.
* `reduce`: 자기주문이 겹치는 수량만큼만 maker, taker **양쪽 주문 수량을 줄여** 자전거래를 방지합니다.
  * **<span style="color:red">즉, 수량이 남아 있으면 주문은 유지됩니다.</span>**

## 자전거래 체결 방지 설정으로 인해 주문이 취소되었는지 어떻게 아나요?

* 자전거래 방지 설정으로 인해 주문이 취소된 경우, 주문 관련 응답에서 다음 필드들을 통해 확인할 수 있습니다:

<Table align={["left","left"]}>
  <thead>
    <tr>
      <th>
        필드
      </th>

      <th>
        설명
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        smp_type
      </td>

      <td>
        적용된 자전거래 체결 방지 모드를 알려줍니다.

        \- `reduce`<br />- `cancel_maker`<br />- `cancel_taker`
      </td>
    </tr>

    <tr>
      <td>
        prevented_volume
      </td>

      <td>
        자전거래 체결 방지 설정으로 인해 취소된 주문 수량
      </td>
    </tr>

    <tr>
      <td>
        prevented_locked
      </td>

      <td>
        (매수 시) 자전거래 체결 방지 설정으로 인해 취소된 금액<br />(매도 시) 자전거래 체결 방지 설정으로 인해 취소된 수량
      </td>
    </tr>
  </tbody>
</Table>

* 웹소켓으로 내 주문 및 체결 (MyOrder) 데이터 구독 시, 주문 상태를 표시하는 `state` 필드가 `prevented`로 전송될 경우 자전거래 체결 방지 설정으로 인해 취소된 수량(또는 금액) 입니다.
  * 자전거래 체결 방지 설정으로 인해 주문의 수량(또는 금액)이 더 이상 남아있지 않아 주문이 취소되는 경우 `cancel`이 전송됩니다.

## 유의사항

* 해당 기능은 선택(Optional) 사항으로, `smp_type`을 명시하지 않으면 자전거래 체결 방지 기능은 적용되지 않습니다. 즉, 동일 회원 간 주문이 취소되지 않고 체결될 수 있습니다.
* 주문 시점이 아닌 체결 시점으로 자전거래를 체크합니다.
* 시장가, 지정가, 예약가, IOC/FOK, 취소 후 재주문 등 모든 주문 타입(`ord_type`)과 함께 사용 가능합니다.
* `time_in_force`파라미터 에서 `post_only` 설정 시 함께 사용 할 수 없습니다.
* `prevented_locked`는
  * 매수 주문 시엔 해당기능으로 체결 방지된 주문 금액이며, 수수료를 포함한 금액입니다. 예를 들어 수수료가 0.25%인 경우, 체결 방지된 주문 금액에 1.0025를 곱한 값으로 계산됩니다.
  * 매도 주문 시엔 수수료와 상관없이 자전거래 체결 방지 기능에 의해 취소된 수량으로, `prevented_volume`과 동일합니다.
* 주문관련 Rest API에서의 응답은 null인 경우 아예 해당 필드가 내려가지 않기 때문에, 값이 있을 때만 다음과 같이 관련 정보를 확인할 수 있습니다.
  * 단, `prevented_locked` 및 `prevented_volume` 은 기본 값이 0으로, 값이 항상 존재합니다.

```json
{  
  "uuid": "53afa136-8882-46e5-8119-614ae10e623b",  
  "side": "bid",  
  "...": "...",
  "smp_type": "cancel_maker", // null이 아닐때만 응답에 포함
  "prevented_volume": 1.174291929,
  "prevented_locked": 0.001706246173
}
```

* 웹소켓으로 내 주문 및 체결 (MyOrder) 데이터 구독 시 null인 경우에도 해당 필드는 내려가기 때문에, 값이 없을 때도 다음과 같이 관련 정보를 확인할 수 있습니다.

```json
{
  "type": "myOrder",
  "code": "KRW-BTC",
  "uuid": "ac2dc2a3-fce9-40a2-a4f6-5987c25c438f",
  "ask_bid": "BID",
  "order_type": "limit",
  "...": "...",
  "smp_type": null, // null이여도 응답에 포함
  "prevented_volume": "0",  
  "prevented_locked": "0"
  ...
}
```

## 예시

<Image src="https://files.readme.io/d775474a9c5185e82350c0bdc2d5eb96b0275b564d255fc76cb640d2c85b725d-kr_smp_head.png" alt="smp_guide_img1" align="center" />

<Image src="https://files.readme.io/7c30a02fb8a977e6a25951296f182f7581bac36c502e49a1cc16cc4b5f6e41da-KR_2.png" alt="smp_guide_img2" align="center" />

# Sibling pages

* [원화(KRW) 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/krw-market-info.md)
* [BTC 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/btc-market-info.md)
* [USDT 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/usdt-market-info.md)