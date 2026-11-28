---
updatedAt: 2026-11-26T09:16:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기(Lending)

코인빌리기는 보유한 원화(KRW) 또는 지원되는 디지털 자산을 담보로 맡기고 디지털 자산을 빌리는 서비스입니다.

## 시작하기 전에

<Table>
  <thead>
    <tr>
      <th>
        항목
      </th>

      <th>
        내용
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        API Key 권한
      </td>

      <td>
        메인포켓 API Key의 코인빌리기(`lending`) 권한이 필요합니다. <br />서브포켓 API Key는 지원하지 않습니다.
      </td>
    </tr>

    <tr>
      <td>
        인증
      </td>

      <td>
        모든 API는 JWT 인증이 필요합니다. 자세한 내용은 [인증 가이드](https://docs.upbit.com/kr/reference/auth)를 참고해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        회원 조건
      </td>

      <td>
        * 코인빌리기를 이용하려면 고객 확인(KYC), 약관 동의 등 필수 요건을 충족하고, 레벨 1 이상이어야 합니다.
        * 보유 자산, 최근 매매 이력, DAXA 교육 등 세부 조건은 아래 "회원 조건"을 참고해 주세요.

        자세한 대여한도 및 레벨별 이용 조건은 <Anchor target="_blank" href="https://support.upbit.com/hc/ko/articles/51996138070809-%EC%BD%94%EC%9D%B8%EB%B9%8C%EB%A6%AC%EA%B8%B0-%EB%8C%80%EC%97%AC-%ED%95%9C%EB%8F%84-%EC%95%88%EB%82%B4">이용 조건</Anchor>을 참고해 주세요.
      </td>
    </tr>
  </tbody>
</Table>

## 회원 조건

코인빌리기는 다음 조건을 모두 충족한 회원만 이용할 수 있습니다. API Key에 코인빌리기 권한을 설정할 때 이용 자격을 확인하며, 이후 조건을 충족하지 않게 되면 API 요청 시 에러가 반환됩니다. 단, DAXA 온라인 교육 이수 여부는 권한을 설정할 때가 아니라 신청 테스트·신규 신청·추가신청을 요청할 때 확인합니다.

| 조건            | 내용                                                                                    | 미충족 시 에러                                         |
| ------------- | ------------------------------------------------------------------------------------- | ------------------------------------------------ |
| 이용 대상         | 대한민국 국적의 개인회원만 이용할 수 있습니다. 법인회원은 이용할 수 없습니다.                                          | -                                                |
| 고객 확인(KYC)    | 고객 확인을 완료해야 하며, 갱신이 필요한 상태가 아니어야 합니다.                                                 | `invalid_kyc`, `invalid_kyc_validity`            |
| 약관 및 정보 제공 동의 | 코인빌리기 서비스 이용약관과 신용정보 제3자 제공에 동의해야 합니다.                                                | `not_agreed_term`, `not_agreed_daxa_use_privacy` |
| DAXA 교육 이수    | 신청 테스트·신규 신청·추가신청을 이용하려면 DAXA 온라인 교육을 이수해야 합니다. 이수하지 않은 상태에서도 권한 설정과 조회는 가능합니다.       | `not_agreed_daxa_education`                      |
| 최소 자격 조건      | 보유 원화(KRW)와 산정 대상 디지털 자산의 원화 환산액 합계가 500만 원 이상이고, 최근 1년간 디지털 자산 매매 이력이 10회 이상이어야 합니다. | `insufficient_level`                             |

위 조건을 충족하더라도 업비트 정책에 따라 일부 회원은 이용이 제한될 수 있습니다. 자세한 내용은 <Anchor target="_blank" href="https://www.upbit.com/lending/guide?tab=faq">코인빌리기 가이드</Anchor> 및 코인빌리기 서비스 소개 하단에 있는 <Anchor target="_blank" href="https://www.upbit.com/lending#:~:text=%EB%B9%8C%EB%A6%AC%EA%B8%B0%20%EB%82%B4%EC%97%AD-,%EC%BD%94%EC%9D%B8%EB%B9%8C%EB%A6%AC%EA%B8%B0%20%EC%84%9C%EB%B9%84%EC%8A%A4%20%EC%9C%A0%EC%9D%98%EC%82%AC%ED%95%AD,-%EC%84%9C%EB%B9%84%EC%8A%A4%20%EC%9D%B4%EC%9A%A9%20%EC%9E%90%EA%B2%A9">서비스 유의사항</Anchor>을 확인해 주세요.

## 이용 흐름

코인빌리기는 **이용 전 확인 → 신청 → 이용 중 → 종료** 순서로 진행됩니다.

코인빌리기를 신청하면 포지션이 생성됩니다. 포지션은 빌린 자산과 맡긴 담보를 함께 관리하는 단위로, 상환 또는 강제상환으로 종료될 때까지 유지됩니다. 추가신청, 담보추가 및 상환은 대상 포지션의 `uuid`를 지정해 요청합니다.

<Table>
  <thead>
    <tr>
      <th style={{ width: "14%" }}>
        단계
      </th>

      <th style={{ width: "40.2%" }}>
        API
      </th>

      <th>
        설명
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        이용 전 확인
      </td>

      <td>
        [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information)​<br />`GET /v1/lendings/chance`
      </td>

      <td>
        빌릴 수 있는 자산, LTV 기준, 수수료율, 이용 한도를 확인합니다.
      </td>
    </tr>

    <tr>
      <td>
        이용 전 확인
      </td>

      <td>
        [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test)<br />`POST /v1/lendings/test`
      </td>

      <td>
        실제로 신청하지 않고 신청 가능 여부와 예상 LTV를 확인합니다.
      </td>
    </tr>

    <tr>
      <td>
        신청
      </td>

      <td>
        [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending)<br />`POST /v1/lendings`
      </td>

      <td>
        빌릴 자산과 담보를 지정해 코인빌리기를 신청합니다.
      </td>
    </tr>

    <tr>
      <td>
        신청
      </td>

      <td>
        [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)<br />`GET /v1/lendings/request`
      </td>

      <td>
        신청, 상환 등 요청의 처리 결과를 확인합니다.
      </td>
    </tr>

    <tr>
      <td>
        이용 중
      </td>

      <td>
        [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)<br />`GET /v1/lending`
      </td>

      <td>
        포지션의 상태와 현재 LTV를 확인합니다.
      </td>
    </tr>

    <tr>
      <td>
        이용 중
      </td>

      <td>
        [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more)<br />`POST /v1/lendings/borrow_more`
      </td>

      <td>
        이용 중인 포지션에서 자산을 더 빌립니다.
      </td>
    </tr>

    <tr>
      <td>
        이용 중
      </td>

      <td>
        [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral)<br />`POST /v1/lendings/adjust_collateral`
      </td>

      <td>
        담보를 추가해 LTV를 낮춥니다.
      </td>
    </tr>

    <tr>
      <td>
        종료
      </td>

      <td>
        [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending)

        `POST /v1/lendings/repay`
      </td>

      <td>
        빌린 자산과 동일한 종류의 자산으로 일부 또는 전액을 상환합니다.
      </td>
    </tr>

    <tr>
      <td>
        종료
      </td>

      <td>
        [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral)​<br />`POST /v1/lendings/repay_with_collateral`
      </td>

      <td>
        맡긴 담보를 활용해 빌린 자산을 전액 상환하며, 담보자산과 빌린 자산의 종류에 따라 담보자산을 매도하거나 상환에 필요한 자산을 매수하는 거래가 발생할 수 있습니다.
      </td>
    </tr>

    <tr>
      <td>
        내역 조회
      </td>

      <td>
        [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings)<br />`GET /v1/lendings/ongoing`
      </td>

      <td>
        종료되지 않은 포지션 목록을 조회합니다.  이용 중(active), 이용 제한(suspended), 강제상환 또는 담보상환 진행 중(liquidating) 상태가 포함됩니다.
      </td>
    </tr>

    <tr>
      <td>
        내역 조회
      </td>

      <td>
        [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings)<br />`GET /v1/lendings/closed`
      </td>

      <td>
        상환 또는 강제상환으로 종료된 포지션 목록을 조회합니다.
      </td>
    </tr>

    <tr>
      <td>
        내역 조회
      </td>

      <td>
        [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions)<br />`GET /v1/lendings/transactions`
      </td>

      <td>
        신청, 담보추가, 상환 등 거래 내역을 조회합니다.
      </td>
    </tr>
  </tbody>
</Table>

## 응답 값 표기 규칙

코인빌리기 API의 금액, 수량 및 비율 관련 값은 숫자 형식의 문자열로 제공되며, 필드에 따라 소수점 자릿수가 다를 수 있습니다.

* **자산 수량:** 소수점 이하 최대 20자리까지 반환될 수 있습니다. 예: `"30.00000000000000000000"`
* **원화 환산 금액:** 코인빌리기 거래 내역 조회의 `transactions[].entries[].fiat_value`는 소수점 이하 8자리로 반올림하여 반환합니다.

## 요청 처리 방식

신청, 추가신청, 담보추가, 상환, 담보상환 요청은 접수 시 `202 Accepted`와 `request_uuid`를 반환합니다. 실제 처리는 접수 이후에 진행되며, 처리 결과는 **코인빌리기 요청 상태 조회** API의 `state`로 확인할 수 있습니다.

| `state`      | 의미                                               |
| ------------ | ------------------------------------------------ |
| `processing` | 요청을 처리하고 있습니다.                                   |
| `done`       | 요청이 정상적으로 처리되었습니다.                               |
| `failed`     | 요청 처리에 실패했습니다. 실패 사유는 `failure` 필드에서 확인할 수 있습니다. |

<Callout icon="far fa-circle-info" theme="info">
  #### `202 Accepted`**는 요청 접수를 의미하며, 처리 완료를 의미하지 않습니다.**

  응답을 받지 못했거나 요청이 `processing` 상태인 경우, 같은 요청을 다시 보내기 전에 요청 상태를 먼저 조회해 주세요. 확인 없이 재요청하면 코인빌리기가 중복으로 신청되거나 상환될 수 있습니다.
</Callout>

## 포지션 상태

포지션의 현재 상태는 `state` 필드로 확인할 수 있습니다. **요청 처리 상태**와 **포지션 상태**는 서로 다른 개념입니다.

| 구분   | `state`       | 의미                         |
| ---- | ------------- | -------------------------- |
| 진행 중 | `active`      | 정상적으로 이용 중인 상태입니다.         |
| 진행 중 | `suspended`   | 이용이 제한된 상태입니다.             |
| 진행 중 | `liquidating` | 강제상환 또는 담보상환이 진행 중인 상태입니다. |
| 종료   | `repaid`      | 빌린 자산을 모두 상환해 종료된 상태입니다.   |
| 종료   | `liquidated`  | 강제상환 또는 담보상환해 종료된 상태입니다.   |

구체적인 종료 방식은 `closure_type`으로 확인해 주세요.

## 라이프사이클별 LTV

LTV(렌딩비율)는 **담보 인정금액 대비 빌린 자산의 가치 비율**입니다. 각 담보의 인정금액은 해당 자산의 평가금액에 담보 인정 비율을 적용해 계산합니다.

LTV는 소수 비율로 반환합니다. 예를 들어 `"0.85"`는 85%를 의미합니다. 단계별 확인 필드는 다음과 같습니다.

<Table>
  <thead>
    <tr>
      <th>
        필드
      </th>

      <th style={{ width: "29.5%" }}>
        의미
      </th>

      <th style={{ width: "48.3%" }}>
        응답 API
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        `estimated_ltv`
      </td>

      <td>
        신청 테스트 시점의 정보로 계산한 예상 LTV입니다.
      </td>

      <td>
        `POST /v1/lendings/test`
      </td>
    </tr>

    <tr>
      <td>
        `current_ltv`
      </td>

      <td>
        조회 시점에 계산한 현재 LTV입니다. 포지션 상태가 `active`인 경우에만 제공하며, 그 외에는 `null`로 반환합니다.
      </td>

      <td>
        `GET /v1/lending`
      </td>
    </tr>

    <tr>
      <td>
        `liquidation_ltv`
      </td>

      <td>
        포지션에 적용되는 강제상환 기준 LTV입니다.
      </td>

      <td>
        신청 전 정책: `GET /v1/lendings/chance`<br />포지션별 기준: `GET /v1/lendings/ongoing`, `GET /v1/lending`, `GET /v1/lendings/closed`
      </td>
    </tr>
  </tbody>
</Table>

<br />

<Callout icon="far fa-triangle-exclamation" theme="warn">
  ### **LTV 관리 및 강제상환 유의사항**

  코인빌리기 API는 LTV 변동에 따른 알림 기능을 제공하지 않습니다. 다만, LTV 위험 등 특정 상황에서는 별도로 안내될 수 있습니다.

  **안내 여부와 관계없이 이용 중에는&#x20;**`current_ltv`**를 주기적으로 확인해 주세요.** LTV가 상승하여 강제상환 위험이 높아지는 경우, 담보를 추가하거나 빌린 자산을 상환하여 LTV를 낮춰야 합니다.

  강제상환 시에는 맡긴 담보가 처분되고 강제상환수수료가 부과될 수 있습니다.
</Callout>

## 요청 수 제한

코인빌리기 API의 요청 수 제한은 **포켓 단위**로 적용되며, 같은 그룹에 속한 API는 요청 횟수를 함께 차감합니다. 자세한 내용은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits">요청 수 제한(Rate Limits)</Anchor>을 참고해 주세요.

<Table>
  <thead>
    <tr>
      <th style={{ width: "24.5%" }}>
        그룹
      </th>

      <th style={{ width: "26.6%" }}>
        제한
      </th>

      <th style={{ width: "48.8%" }}>
        대상 API
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        `default`
      </td>

      <td>
        초당 30회
      </td>

      <td>
        가능 정보 조회, 요청 상태 조회, 진행 중·종료된 목록 조회, 거래 내역 조회
      </td>
    </tr>

    <tr>
      <td>
        `lending`
      </td>

      <td>
        초당 10회
      </td>

      <td>
        신청, 추가신청, 담보추가, 상환, 담보상환, 개별 코인빌리기 조회
      </td>
    </tr>

    <tr>
      <td>
        `lending-test`
      </td>

      <td>
        초당 10회
      </td>

      <td>
        코인빌리기 신청 테스트
      </td>
    </tr>
  </tbody>
</Table>

### 요청 수 제한과 서비스 이용 조건

`lending` 그룹의 초당 10회는 API 요청 횟수의 상한입니다. 이 한도를 준수하더라도 다음 서비스 정책에 따라 요청이 처리되지 않을 수 있습니다.

* **최대 이용 건수:** 동시에 이용할 수 있는 포지션은 현재 최대 1개입니다. 변경될 수 있으므로 코인빌리기 가능 정보 조회의 `account_limits.max_lending_count`를 확인해 주세요.

개별 코인빌리기 조회도 `lending` 그룹에 포함됩니다. 짧은 주기로 반복 조회하면 신청이나 상환 요청의 이용 가능 횟수가 줄어들 수 있습니다.

## 에러 안내

API 응답의 `error.name`을 기준으로 오류를 처리해 주세요. `error.message`는 사용자 안내 문구이므로 프로그램의 분기 조건으로 사용하지 않는 것을 권장합니다.

```json
{
  "error": {
    "name": "invalid_lending_ratio",
    "message": "신청 가능한 렌딩비율은 {N}%~{M}%예요."
  }
}
```

### 400 Bad Request

요청 형식이나 파라미터가 올바르지 않은 경우 발생합니다.

<Table>
  <thead>
    <tr>
      <th style={{ width: "40.9%" }}>
        에러 코드 (`error.name`)
      </th>

      <th style={{ width: "29%" }}>
        발생 이유
      </th>

      <th style={{ width: "25.4%" }}>
        해결 방법
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        `invalid_parameter`
      </td>

      <td>
        필수 값이 누락되었거나 형식이 올바르지 않습니다.
      </td>

      <td>
        요청 값을 확인해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `duplicated_identifier`
      </td>

      <td>
        이미 사용한 `identifier`입니다.
      </td>

      <td>
        응답을 받지 못했거나 `duplicated_identifier`를 받았다면 원래 `identifier`로 요청 상태를 먼저 조회해 주세요. 별도 요청을 보낼 때는 새 `identifier`를 사용해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `request_collateral_amount_less_than_min_amount`
      </td>

      <td>
        담보자산의 원화 평가금액이 최소 담보 평가금액보다 적습니다.
      </td>

      <td>
        최소 담보 평가금액 이상이 되도록 담보 수량을 늘려 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `below_min_lending_valuation`
      </td>

      <td>
        빌리려는 자산의 원화 평가금액이 최소 대여금액보다 적습니다.
      </td>

      <td>
        최소 대여금액 이상이 되도록 신청 수량을 늘려 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `invalid_request`
      </td>

      <td>
        일시적인 요청 오류입니다.
      </td>

      <td>
        잠시 후 다시 요청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `insufficient_collateral_main_balance`
      </td>

      <td>
        담보로 설정할 자산의 보유 수량이 부족합니다.
      </td>

      <td>
        포켓의 보유 수량을 확인해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `insufficient_main_balance`
      </td>

      <td>
        담보추가 또는 상환에 필요한 수량이 부족합니다.
      </td>

      <td>
        포켓의 보유 수량을 확인해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `exceed_lending_repay_amount`
      </td>

      <td>
        상환 수량이 남은 빌린 수량을 초과합니다.
      </td>

      <td>
        남은 수량 이하로 입력하거나 `repayment_mode: full`을 사용해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `invalid_lending_ratio`
      </td>

      <td>
        예상 LTV가 신청 가능한 범위를 벗어났습니다.
      </td>

      <td>
        * 신규 신청: 에러 메시지의 `{N}%~{M}%` 범위에 맞게 빌릴 수량 또는 신청에 사용할 담보 구성을 조정해 주세요.
        * 추가신청: 개별 코인빌리기 조회(`GET /v1/lending`)에서 포지션 상태와 현재 LTV를 확인한 뒤 추가로 빌릴 수량을 조정해 주세요. 담보를 추가한 이력이 있는 포지션은 추가신청이 제한되므로, 담보를 추가하는 대신 신규 신청을 검토해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `exceed_collateral_count_limit`
      </td>

      <td>
        담보자산 종류 수가 최대치를 초과했습니다.
      </td>

      <td>
        메시지에 안내된 `{N}`개 이하로 설정해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `already_lent_collateral_currency`
      </td>

      <td>
        대여 중인 자산을 담보로 설정했습니다.
      </td>

      <td>
        다른 자산으로 담보를 구성해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `scheduled_liquidation_in_progress`
      </td>

      <td>
        예약청산 대상 자산이 포함되어 있습니다.
      </td>

      <td>
        해당 자산을 제외해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `insufficient_lending_budget`
      </td>

      <td>
        자산별 신청 가능 수량이 소진되었습니다.
      </td>

      <td>
        `available_budget_amount`를 확인한 뒤 잠시 후 다시 요청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `not_related_collateral`
      </td>

      <td>
        최초 신청 때 설정하지 않은 담보자산입니다.
      </td>

      <td>
        최초 신청 시 설정한 담보자산만 사용해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `validate_increment_for_lending_ratio`
      </td>

      <td>
        추가신청 렌딩비율 기준을 충족하지 않습니다.
      </td>

      <td>
        메시지에 안내된 `{M}%` 미만이 된 뒤 다시 신청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `block_lending_increment_if_deposit_collateral`
      </td>

      <td>
        담보금 입금 이력이 있어 추가신청이 제한됩니다.
      </td>

      <td>
        신규 신청을 검토해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `exceed_total_lending_count_limit`
      </td>

      <td>
        계정의 최대 이용 건수를 초과했습니다.
      </td>

      <td>
        이용 중인 포지션을 상환해 종료한 뒤 신청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `exceed_currency_lending_count_limit`
      </td>

      <td>
        자산별 최대 이용 건수를 초과했습니다.
      </td>

      <td>
        같은 자산으로 이용 중인 포지션을 상환해 종료한 뒤 신청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `exceed_loan_cap_limit`
      </td>

      <td>
        계정의 남은 대여한도를 초과했습니다.
      </td>

      <td>
        `remaining_lending_fiat`을 확인해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `exceed_collateral_amount_limit`
      </td>

      <td>
        계정의 담보 한도를 초과했습니다.
      </td>

      <td>
        담보 수량을 조정해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `insufficient_level`
      </td>

      <td>
        최소 자격 조건을 충족하지 않습니다.
      </td>

      <td>
        회원 조건을 확인해 주세요. 조건 충족 직후라면 자주 묻는 질문을 참고해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `terminated_state`
      </td>

      <td>
        이미 종료된 포지션입니다.
      </td>

      <td>
        종료된 코인빌리기 목록에서 확인해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `invalid_state`
      </td>

      <td>
        현재 요청할 수 없는 포지션 상태입니다.
      </td>

      <td>
        포지션 상태를 조회한 뒤 잠시 후 다시 요청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `processing_transaction_exists`
      </td>

      <td>
        같은 포지션에 처리 중인 요청이 있습니다.
      </td>

      <td>
        진행 중인 요청의 처리가 끝난 뒤 다시 요청해 주세요.
      </td>
    </tr>

    <tr>
      <td>
        `invalid_price_data`
      </td>

      <td>
        시세 정보를 확인할 수 없습니다.
      </td>

      <td>
        잠시 후 다시 요청해 주세요.
      </td>
    </tr>
  </tbody>
</Table>

### 403 Forbidden

API Key 권한이나 회원 조건 때문에 요청이 허용되지 않는 경우 발생합니다.

| 에러 코드 (`error.name`)          | 발생 이유                                                                                                                                  | 해결 방법                                                            |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| `out_of_scope`                | 코인빌리기 권한이 없는 API Key입니다.                                                                                                               | 메인포켓 API Key의 `lending` 권한을 확인해 주세요.                             |
| `invalid_kyc`                 | 고객 확인(KYC)이 완료되지 않았습니다.                                                                                                                | 앱 또는 웹에서 고객 확인을 완료해 주세요.                                         |
| `invalid_kyc_validity`        | 고객 확인 갱신이 필요합니다.                                                                                                                       | 앱 또는 웹에서 고객 확인을 갱신한 뒤 다시 신청해 주세요.                                |
| `not_agreed_term`             | 코인빌리기 서비스 이용약관에 동의하지 않았거나 동의를 철회했습니다.                                                                                                  | 앱 또는 웹에서 약관에 다시 동의해 주세요.                                         |
| `not_agreed_daxa_use_privacy` | 신용정보 제3자 제공에 동의하지 않았거나 동의를 철회했습니다.                                                                                                     | 동의한 뒤 다시 요청해 주세요.                                                |
| `not_agreed_daxa_education`   | DAXA 온라인 교육을 이수하지 않은 상태에서 신청 테스트(`POST /v1/lendings/test`), 신규 신청(`POST /v1/lendings`), 추가신청(`POST /v1/lendings/borrow_more`)을 요청했습니다. | DAXA 온라인 교육을 이수한 뒤 다시 요청해 주세요. 이수하지 않은 상태에서도 조회 API는 이용할 수 있습니다. |

### 404 Not Found

요청한 포지션이나 요청을 찾을 수 없는 경우 발생합니다.

| 에러 코드 (`error.name`) | 발생 이유                      | 해결 방법                                    |
| -------------------- | -------------------------- | ---------------------------------------- |
| `lending_not_found`  | 포지션이 없거나 다른 포켓에 속한 포지션입니다. | `uuid`를 확인해 주세요.                         |
| `request_not_found`  | 요청이 없거나 다른 포켓에 속한 요청입니다.   | `request_uuid` 또는 `identifier`를 확인해 주세요. |

### 409 Conflict

요청 처리 중 시세나 정보가 변경된 경우 발생합니다. 동일한 조건으로 재요청하면 다시 실패할 수 있으므로, 에러 메시지와 요청 조건을 확인해 주세요.

| 에러 코드 (`error.name`)                  | 발생 이유                                                                                          | 해결 방법                                                                                                                                                                                                                                                                              |
| ------------------------------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `exceed_ltv_deviation_limit`          | 신규 신청(`POST /v1/lendings`) 또는 추가신청(`POST /v1/lendings/borrow_more`) 처리 중 시세 변동으로 LTV가 변경되었습니다. | 신규 신청(`POST /v1/lendings`): 빌릴 수량 또는 신청에 사용할 담보 구성을 조정해 다시 요청해 주세요.<br /><br />추가신청(`POST /v1/lendings/borrow_more`): 개별 코인빌리기 조회(`GET /v1/lending`)에서 포지션 상태와 현재 LTV를 확인한 뒤 추가로 빌릴 수량을 조정해 주세요. 담보를 추가한 이력이 있는 포지션은 추가신청이 제한되므로, 담보를 추가하는 대신 신규 신청(`POST /v1/lendings`)을 검토해 주세요. |
| `amount_limit_price_snapshot_changed` | 처리 중 가격 정보가 변경되었습니다.                                                                           | 잠시 후 다시 요청해 주세요.                                                                                                                                                                                                                                                                   |
| `command_preparation_stale`           | 처리 중 정보가 변경되었습니다.                                                                              | 잠시 후 다시 요청해 주세요.                                                                                                                                                                                                                                                                   |

### 500 Internal Server Error

서버 내부 오류로 요청을 처리할 수 없는 경우 발생합니다.

| 에러 코드 (`error.name`)          | 발생 이유              | 해결 방법            |
| ----------------------------- | ------------------ | ---------------- |
| `base_price_must_be_positive` | 기준 시세를 확인할 수 없습니다. | 잠시 후 다시 요청해 주세요. |

### 503 Service Unavailable

서비스 점검 등으로 코인빌리기를 일시적으로 이용할 수 없는 경우 발생합니다. `500` 및 `503` 오류가 발생했다면 동일한 요청을 연속으로 보내지 마세요.

| 에러 코드 (`error.name`)  | 발생 이유                        | 해결 방법            |
| --------------------- | ---------------------------- | ---------------- |
| `service_unavailable` | 코인빌리기 서비스를 일시적으로 이용할 수 없습니다. | 잠시 후 다시 요청해 주세요. |

### 요청 실패 처리

POST 요청이 `202 Accepted`로 접수되었더라도 실제 처리 과정에서 실패할 수 있습니다. 코인빌리기 요청 상태 조회의 `state`가 `failed`이면 `failure` 필드에서 실패 사유를 확인할 수 있습니다.

```json
{
  "state": "failed",
  "action": "borrow",
  "failure": {
    "name": "insufficient_level",
    "message": "코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요."
  },
  "lending_uuid": null
}
```

접수 단계에서 거절된 요청과 접수된 뒤 처리 단계에서 실패한 요청은 확인하는 위치가 다릅니다.

| 구분             | 접수 거절               | 접수 후 처리 실패                      |
| -------------- | ------------------- | ------------------------------- |
| 확인 위치          | 신청·변경 API(POST)의 응답 | 코인빌리기 요청 상태 조회의 응답              |
| HTTP 상태 코드     | `400` 등 에러 코드       | `200`                           |
| 실패 사유          | `error.name`        | `state: failed`와 `failure.name` |
| `request_uuid` | 발급되지 않음             | 발급됨                             |

예를 들어 `processing_transaction_exists`는 접수 단계에서 거절되면 HTTP 400의 `error.name`으로, 접수된 뒤 처리 단계에서 실패하면 요청 상태 조회의 `failure.name`으로 반환됩니다.

실패 사유는 `failure.name`을 기준으로 구분해 주세요. 요청 유형(`action`)별로 반환될 수 있는 값과 조치 방법은 아래 표를 참고해 주세요.

#### 신규 신청 (`borrow`)

| `failure.name`                                   | 안내 메시지 (`failure.message`)                                | 조치 방법                                                            |
| ------------------------------------------------ | --------------------------------------------------------- | ---------------------------------------------------------------- |
| `service_unavailable`                            | 코인빌리기 서비스를 이용할 수 없어요. 잠시 후 다시 시도해 주세요.                    | 잠시 후 다시 요청해 주세요.                                                 |
| `base_price_must_be_positive`                    | 기준 시세를 정상적으로 확인할 수 없어 요청을 처리하지 못했어요. 잠시 후 다시 요청해 주세요.     | 잠시 후 다시 요청해 주세요.                                                 |
| `exceed_collateral_count_limit`                  | 담보자산은 최대 {N}개까지 설정할 수 있어요. 담보 설정을 확인해 주세요.                | 메시지에 안내된 `{N}`개 이하로 설정해 주세요.                                     |
| `insufficient_level`                             | 코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요.                  | 회원 조건을 확인해 주세요. 조건 충족 직후라면 자주 묻는 질문을 참고해 주세요.                    |
| `insufficient_collateral_main_balance`           | 담보로 설정할 자산의 보유 수량이 부족해요.                                  | 포켓의 보유 수량을 확인해 주세요.                                              |
| `request_collateral_amount_less_than_min_amount` | 최소 담보금 이상 입력해 주세요.                                        | 최소 담보 평가금액 이상이 되도록 담보 수량을 늘려 주세요.                                |
| `invalid_price_data`                             | 시세 정보를 확인할 수 없어요. 잠시 후 다시 시도해 주세요.                        | 잠시 후 다시 요청해 주세요.                                                 |
| `exceed_ltv_deviation_limit`                     | 시세 변동으로 LTV가 변경되었어요. 빌릴 수량을 변경해 다시 요청해 주세요.               | 빌릴 수량 또는 신청에 사용할 담보 구성을 조정해 다시 요청해 주세요.                          |
| `invalid_lending_ratio`                          | 신청 가능한 렌딩비율은 {N}%\~{M}%예요.                                | 에러 메시지의 `{N}%~{M}%` 범위에 맞게 빌릴 수량 또는 신청에 사용할 담보 구성을 조정해 주세요.      |
| `below_min_lending_valuation`                    | 신청 금액이 최소 기준에 미달해요.                                       | 최소 대여금액 이상이 되도록 신청 수량을 늘려 주세요.                                   |
| `scheduled_liquidation_in_progress`              | 현재 예약청산이 진행 중인 자산을 포함하고 있어요.                              | 해당 자산을 제외해 주세요.                                                  |
| `exceed_total_lending_count_limit`               | 최대 이용 건 수를 초과했어요. 신규 신청을 원하시면 이용 중인 건을 먼저 상환해 주세요.        | 이용 중인 포지션을 상환해 종료한 뒤 신청해 주세요.                                    |
| `exceed_currency_lending_count_limit`            | 해당 코인의 최대 이용 건 수를 초과했어요. 신규 신청을 원하시면 이용 중인 건을 먼저 상환해 주세요. | 같은 자산으로 이용 중인 포지션을 상환해 종료한 뒤 신청해 주세요.                            |
| `already_lent_collateral_currency`               | 대여 중인 자산은 담보로 설정할 수 없어요.                                  | 다른 자산으로 담보를 구성해 주세요.                                             |
| `exceed_loan_cap_limit`                          | 최대 대여한도를 초과하여 신청할 수 없어요.                                  | `remaining_lending_fiat`을 확인해 주세요.                               |
| `exceed_collateral_amount_limit`                 | 최대 설정 가능한 담보금을 초과하여 신청할 수 없어요.                            | 담보 수량을 조정해 주세요.                                                  |
| `amount_limit_price_snapshot_changed`            | 가격 정보가 변경되었어요. 잠시 후 다시 시도해 주세요.                           | 잠시 후 다시 요청해 주세요.                                                 |
| `insufficient_lending_budget`                    | 현재 신청 가능한 수량이 모두 소진되었어요.                                  | `available_budget_amount`를 확인한 뒤 잠시 후 다시 요청해 주세요.                |
| `not_agreed_term`                                | 코인빌리기 서비스 이용약관 동의가 필요해요. 약관 동의를 완료해 주세요.                  | 앱 또는 웹에서 약관에 다시 동의해 주세요.                                         |
| `invalid_kyc_validity`                           | 고객 확인(KYC) 갱신 후 신청이 가능해요. 고객 확인 절차를 다시 완료해 주세요.           | 앱 또는 웹에서 고객 확인을 갱신한 뒤 다시 신청해 주세요.                                |
| `invalid_kyc`                                    | 고객 확인(KYC)을 완료한 후 이용할 수 있어요. 고객 확인 절차를 완료해 주세요.           | 앱 또는 웹에서 고객 확인을 완료해 주세요.                                         |
| `not_agreed_daxa_education`                      | DAXA 가이드라인에 따른 필수 교육을 이수해야 코인을 빌릴 수 있어요. 교육을 먼저 이수해 주세요   | DAXA 온라인 교육을 이수한 뒤 다시 요청해 주세요. 이수하지 않은 상태에서도 조회 API는 이용할 수 있습니다. |
| `transaction_discarded`                          | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                           | 잠시 후 다시 요청해 주세요.                                                 |
| `transaction_creation_unknown`                   | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                           | 잠시 후 다시 요청해 주세요.                                                 |

#### 추가신청 (`borrow_more`)

| `failure.name`                                   | 안내 메시지 (`failure.message`)                              | 조치 방법                                                                      |
| ------------------------------------------------ | ------------------------------------------------------- | -------------------------------------------------------------------------- |
| `service_unavailable`                            | 코인빌리기 서비스를 이용할 수 없어요. 잠시 후 다시 시도해 주세요.                  | 잠시 후 다시 요청해 주세요.                                                           |
| `base_price_must_be_positive`                    | 기준 시세를 정상적으로 확인할 수 없어 요청을 처리하지 못했어요. 잠시 후 다시 요청해 주세요.   | 잠시 후 다시 요청해 주세요.                                                           |
| `insufficient_level`                             | 코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요.                | 회원 조건을 확인해 주세요. 조건 충족 직후라면 자주 묻는 질문을 참고해 주세요.                              |
| `command_preparation_stale`                      | 처리 중 정보가 변경되었어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.                                                           |
| `terminated_state`                               | 이미 상환이 완료됐어요.                                           | 종료된 코인빌리기 목록에서 확인해 주세요.                                                    |
| `invalid_state`                                  | 일시적인 오류가 발생했어요. 잠시 후 다시 시도해 주세요.                        | 포지션 상태를 조회한 뒤 잠시 후 다시 요청해 주세요.                                             |
| `validate_increment_for_lending_ratio`           | 현재 렌딩비율은 {N}%예요. {M}% 미만일 때 추가신청할 수 있어요.                | 메시지에 안내된 `{M}%` 미만이 된 뒤 다시 신청해 주세요.                                        |
| `scheduled_liquidation_in_progress`              | 현재 예약청산이 진행 중인 자산을 포함하고 있어요.                            | 예약청산 대상 자산이 포함된 포지션에는 추가신청을 할 수 없습니다.                                      |
| `block_lending_increment_if_deposit_collateral`  | 담보자산 입금 이력이 있는 경우 추가신청이 제한돼요.                           | 신규 신청을 검토해 주세요.                                                            |
| `processing_transaction_exists`                  | 이미 처리 중인 요청이 있어요. 잠시 후 다시 시도해 주세요.                      | 진행 중인 요청의 처리가 끝난 뒤 다시 요청해 주세요.                                             |
| `invalid_price_data`                             | 시세 정보를 확인할 수 없어요. 잠시 후 다시 시도해 주세요.                      | 잠시 후 다시 요청해 주세요.                                                           |
| `exceed_ltv_deviation_limit`                     | 시세 변동으로 LTV가 변경되었어요. 빌릴 수량을 변경해 다시 요청해 주세요.             | 개별 코인빌리기 조회(`GET /v1/lending`)에서 포지션 상태와 현재 LTV를 확인한 뒤 추가로 빌릴 수량을 조정해 주세요. |
| `invalid_lending_ratio`                          | 신청 가능한 렌딩비율은 {N}%\~{M}%예요.                              | 개별 코인빌리기 조회(`GET /v1/lending`)에서 포지션 상태와 현재 LTV를 확인한 뒤 추가로 빌릴 수량을 조정해 주세요. |
| `validate_increment_lending_amount_not_positive` | 추가 코인빌리기 신청에 실패했어요. 다시 시도해 주세요.                         | 추가로 빌릴 수량을 확인한 뒤 다시 요청해 주세요.                                               |
| `insufficient_lending_budget`                    | 현재 신청 가능한 수량이 모두 소진되었어요.                                | `available_budget_amount`를 확인한 뒤 잠시 후 다시 요청해 주세요.                          |
| `exceed_loan_cap_limit`                          | 최대 대여한도를 초과하여 신청할 수 없어요.                                | `remaining_lending_fiat`을 확인해 주세요.                                         |
| `invalid_kyc`                                    | 고객 확인(KYC)을 완료한 후 이용할 수 있어요. 고객 확인 절차를 완료해 주세요.         | 앱 또는 웹에서 고객 확인을 완료해 주세요.                                                   |
| `not_agreed_term`                                | 코인빌리기 서비스 이용약관 동의가 필요해요. 약관 동의를 완료해 주세요.                | 앱 또는 웹에서 약관에 다시 동의해 주세요.                                                   |
| `invalid_kyc_validity`                           | 고객 확인(KYC) 갱신 후 신청이 가능해요. 고객 확인 절차를 다시 완료해 주세요.         | 앱 또는 웹에서 고객 확인을 갱신한 뒤 다시 신청해 주세요.                                          |
| `not_agreed_daxa_education`                      | DAXA 가이드라인에 따른 필수 교육을 이수해야 코인을 빌릴 수 있어요. 교육을 먼저 이수해 주세요 | DAXA 온라인 교육을 이수한 뒤 다시 요청해 주세요. 이수하지 않은 상태에서도 조회 API는 이용할 수 있습니다.           |
| `transaction_discarded`                          | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                         | 잠시 후 다시 요청해 주세요.                                                           |
| `transaction_creation_unknown`                   | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                         | 잠시 후 다시 요청해 주세요.                                                           |

#### 담보추가 (`add_collateral`)

| `failure.name`                                   | 안내 메시지 (`failure.message`)                            | 조치 방법                                         |
| ------------------------------------------------ | ----------------------------------------------------- | --------------------------------------------- |
| `insufficient_main_balance`                      | 추가할 담보자산의 보유 수량이 부족해요.                                | 포켓의 보유 수량을 확인해 주세요.                           |
| `service_unavailable`                            | 코인빌리기 서비스를 이용할 수 없어요. 잠시 후 다시 시도해 주세요.                | 잠시 후 다시 요청해 주세요.                              |
| `base_price_must_be_positive`                    | 기준 시세를 정상적으로 확인할 수 없어 요청을 처리하지 못했어요. 잠시 후 다시 요청해 주세요. | 잠시 후 다시 요청해 주세요.                              |
| `insufficient_level`                             | 코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요.              | 회원 조건을 확인해 주세요. 조건 충족 직후라면 자주 묻는 질문을 참고해 주세요. |
| `command_preparation_stale`                      | 처리 중 정보가 변경되었어요. 잠시 후 다시 시도해 주세요.                     | 잠시 후 다시 요청해 주세요.                              |
| `request_collateral_amount_less_than_min_amount` | 최소 담보금 이상 입력해 주세요.                                    | 최소 담보 평가금액 이상이 되도록 담보 수량을 늘려 주세요.             |
| `terminated_state`                               | 이미 상환이 완료됐어요.                                         | 종료된 코인빌리기 목록에서 확인해 주세요.                       |
| `invalid_state`                                  | 일시적인 오류가 발생했어요. 잠시 후 다시 시도해 주세요.                      | 포지션 상태를 조회한 뒤 잠시 후 다시 요청해 주세요.                |
| `processing_transaction_exists`                  | 이미 처리 중인 요청이 있어요. 잠시 후 다시 시도해 주세요.                    | 진행 중인 요청의 처리가 끝난 뒤 다시 요청해 주세요.                |
| `exceed_collateral_amount_limit`                 | 최대 설정 가능한 담보금을 초과하여 신청할 수 없어요.                        | 담보 수량을 조정해 주세요.                               |
| `amount_limit_price_snapshot_changed`            | 가격 정보가 변경되었어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.                              |
| `invalid_kyc`                                    | 고객 확인(KYC)을 완료한 후 이용할 수 있어요. 고객 확인 절차를 완료해 주세요.       | 앱 또는 웹에서 고객 확인을 완료해 주세요.                      |
| `not_agreed_term`                                | 코인빌리기 서비스 이용약관 동의가 필요해요. 약관 동의를 완료해 주세요.              | 앱 또는 웹에서 약관에 다시 동의해 주세요.                      |
| `transaction_discarded`                          | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.                              |
| `transaction_creation_unknown`                   | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.                              |

#### 상환 (`repay`)

| `failure.name`                  | 안내 메시지 (`failure.message`)                            | 조치 방법                                            |
| ------------------------------- | ----------------------------------------------------- | ------------------------------------------------ |
| `service_unavailable`           | 코인빌리기 서비스를 이용할 수 없어요. 잠시 후 다시 시도해 주세요.                | 잠시 후 다시 요청해 주세요.                                 |
| `base_price_must_be_positive`   | 기준 시세를 정상적으로 확인할 수 없어 요청을 처리하지 못했어요. 잠시 후 다시 요청해 주세요. | 잠시 후 다시 요청해 주세요.                                 |
| `command_preparation_stale`     | 처리 중 정보가 변경되었어요. 잠시 후 다시 시도해 주세요.                     | 잠시 후 다시 요청해 주세요.                                 |
| `invalid_state`                 | 현재 상환할 수 없어요. 잠시 후 다시 시도해 주세요.                        | 포지션 상태를 조회한 뒤 잠시 후 다시 요청해 주세요.                   |
| `insufficient_main_balance`     | 상환에 필요한 수량이 부족해요.                                     | 포켓의 보유 수량을 확인해 주세요.                              |
| `terminated_state`              | 이미 상환이 완료됐어요.                                         | 종료된 코인빌리기 목록에서 확인해 주세요.                          |
| `processing_transaction_exists` | 이미 처리 중인 요청이 있어요. 잠시 후 다시 시도해 주세요.                    | 진행 중인 요청의 처리가 끝난 뒤 다시 요청해 주세요.                   |
| `exceed_lending_repay_amount`   | 상환 수량을 남은 대여 수량 이하로 입력해 주세요.                          | 남은 수량 이하로 입력하거나 `repayment_mode: full`을 사용해 주세요. |
| `invalid_kyc`                   | 고객 확인(KYC)을 완료한 후 이용할 수 있어요. 고객 확인 절차를 완료해 주세요.       | 앱 또는 웹에서 고객 확인을 완료해 주세요.                         |
| `not_agreed_term`               | 코인빌리기 서비스 이용약관 동의가 필요해요. 약관 동의를 완료해 주세요.              | 앱 또는 웹에서 약관에 다시 동의해 주세요.                         |
| `transaction_discarded`         | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.                                 |
| `transaction_creation_unknown`  | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.                                 |

#### 담보상환 (`repay_with_collateral`)

| `failure.name`                 | 안내 메시지 (`failure.message`)                            | 조치 방법                          |
| ------------------------------ | ----------------------------------------------------- | ------------------------------ |
| `base_price_must_be_positive`  | 기준 시세를 정상적으로 확인할 수 없어 요청을 처리하지 못했어요. 잠시 후 다시 요청해 주세요. | 잠시 후 다시 요청해 주세요.               |
| `invalid_state`                | 현재 상환할 수 없어요. 잠시 후 다시 시도해 주세요.                        | 포지션 상태를 조회한 뒤 잠시 후 다시 요청해 주세요. |
| `transaction_creation_unknown` | 요청을 처리하지 못했어요. 잠시 후 다시 시도해 주세요.                       | 잠시 후 다시 요청해 주세요.               |

요청이 실패했더라도 이미 접수된 요청의 `identifier`는 다시 사용할 수 없습니다. 재시도할 때는 새로운 `identifier`를 지정해 주세요.

## 서비스 정책 안내

수수료, 이용기간, 강제상환 등 코인빌리기 서비스 정책은 다음 문서를 참고해 주세요.

| 구분             | 링크                                                             |
| -------------- | -------------------------------------------------------------- |
| 코인빌리기 가이드      | [코인빌리기 가이드 및 FAQ](https://www.upbit.com/lending/guide?tab=faq) |
| 코인빌리기 약관       | [서비스 이용약관](https://www.upbit.com/terms/lending)                |
| 코인빌리기 서비스 유의사항 | [서비스 유의사항](https://www.upbit.com/lending)                      |

## 자주 묻는 질문

<Accordion title="`202 Accepted`가 반환되면 코인빌리기가 완료된 건가요?" icon="fa-info-circle">
  아니요. `202 Accepted`는 요청이 접수되었다는 뜻입니다. 응답의 `request_uuid`로 코인빌리기 요청 상태를 조회해 `state`가 `done`인지 확인해 주세요. `failed`라면 `failure`에서 사유를 확인할 수 있습니다.
</Accordion>

<Accordion title="요청에 대한 응답을 받지 못했습니다. 다시 요청해도 되나요?" icon="fa-info-circle">
  즉시 다시 요청하지 말고, 요청 시 지정한 `identifier`로 코인빌리기 요청 상태를 먼저 조회해 주세요. 요청이 이미 접수되었는데 다시 요청하면 중복 신청이나 상환이 발생할 수 있습니다.

  응답을 받지 못한 경우에도 요청 상태를 조회할 수 있도록, 신청·상환 요청 시 `identifier`를 지정하는 것을 권장합니다.
</Accordion>

<Accordion title="신청 테스트가 성공하면 실제 신청도 반드시 성공하나요?" icon="fa-info-circle">
  아니요. 코인빌리기 신청 테스트는 테스트 요청 시점을 기준으로 검증합니다. 이후 시세나 신청 가능 수량이 달라지면 실제 신청 결과도 달라질 수 있습니다.

  테스트 응답의 `request_uuid`는 테스트용 값이므로 요청 상태 조회에 사용할 수 없습니다.
</Accordion>

<Accordion title="이용 중 LTV는 어떻게 관리하나요?" icon="fa-info-circle">
  개별 코인빌리기 조회의 `current_ltv`가 강제상환 기준인 `liquidation_ltv`에 도달하지 않도록 관리해 주세요. LTV를 낮추려면 담보추가 또는 코인빌리기 상환을 이용할 수 있습니다.
</Accordion>

<Accordion title="코인빌리기를 일부만 상환할 수 있나요?" icon="fa-info-circle">
  네. 코인빌리기 상환 API에서 `repayment_amount`에 상환할 수량을 입력하면 일부 상환할 수 있습니다. 남은 수량을 모두 상환하려면 `repayment_mode`를 `full`로 지정해 주세요.

  담보상환 API는 전액 상환만 지원합니다.
</Accordion>

<Accordion title="담보를 추가한 뒤에도 자산을 더 빌릴 수 있나요?" icon="fa-info-circle">
  아니요. 담보를 추가한 포지션은 추가신청할 수 없으며, `block_lending_increment_if_deposit_collateral` 에러가 반환됩니다. 추가로 자산을 빌리려면 새로운 코인빌리기 신청을 검토해 주세요.
</Accordion>

<Accordion title="이용 조건을 충족했는데 `insufficient_level` 에러가 반환됩니다." icon="fa-info-circle">
  이용 조건을 충족한 직후에도 `insufficient_level`이 반환될 수 있습니다. 이 경우 업비트 앱 또는 웹의 **코인빌리기 '내 빌리기한도' 페이지에서 우측 상단 한도 업데이트**한 뒤 요청을 다시 보내 주세요.

  * 앱 : 코인빌리기 > 자격 확인 - 알아보기 > 우측상단 '한도 업데이트'
  * PC웹 : 코인빌리기 > 자격 확인 - 알아보기 > 우측상단 '한도 업데이트'

  새로고침은 한 번만 하면 됩니다. 이후에도 같은 에러가 반복되면 고객센터로 문의해 주세요.
</Accordion>

<Accordion title="수수료율과 이용 한도는 어디에서 확인할 수 있나요?" icon="fa-info-circle">
  **신청 전 현재 정책의 수수료율과 이용 한도**는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-lending-information">코인빌리기 가능 정보 조회</Anchor> API에서 확인할 수 있습니다.

  | 확인할 내용         | 응답 필드                                                                      |
  | -------------- | -------------------------------------------------------------------------- |
  | 이용수수료율(일 단위)   | `lending_assets[].borrowing_fee_rate`                                      |
  | 강제상환수수료율       | `lending_assets[].liquidation_fee_rate`                                    |
  | 담보상환수수료율       | `lending_assets[].collateral_repayment_fee_rate`                           |
  | 계정의 최대·남은 대여한도 | `account_limits.max_lending_fiat`, `account_limits.remaining_lending_fiat` |
  | 최대·현재 이용 건수    | `account_limits.max_lending_count`, `account_limits.ongoing_lending_count` |

  **이미 생성된 코인빌리기 포지션에 적용된 수수료율**은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor> API에서 확인할 수 있습니다. 해당 포지션의 UUID로 조회해 주세요.

  | 확인할 내용           | 응답 필드                           |
  | ---------------- | ------------------------------- |
  | 적용된 이용수수료율(일 단위) | `borrowing_fee_rate`            |
  | 적용된 강제상환수수료율     | `liquidation_fee_rate`          |
  | 적용된 담보상환수수료율     | `collateral_repayment_fee_rate` |

  신규 신청 시 적용되는 현재 정책과 기존 포지션에 적용된 요율은 다를 수 있으므로, **기존 포지션의 수수료율은 개별 코인빌리기 조회 결과를 기준으로 확인해 주세요.**
</Accordion>

# 하위 페이지

* [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information.md)
* [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test.md)
* [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending.md)
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request.md)
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)

# 같은 수준의 페이지

* [포켓(Pocket)](https://docs.upbit.com/kr/reference/pocket-overview.md)
* [자산(Asset)](https://docs.upbit.com/kr/reference/자산asset.md)
* [주문(Order)](https://docs.upbit.com/kr/reference/exchange.md)
* [출금(Withdrawal)](https://docs.upbit.com/kr/reference/출금withdrawal.md)
* [입금(Deposit)](https://docs.upbit.com/kr/reference/입금deposit.md)
* [서비스 정보(Service)](https://docs.upbit.com/kr/reference/서비스-정보service.md)