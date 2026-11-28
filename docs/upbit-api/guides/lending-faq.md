---
updatedAt: 2026-11-26T04:15:47.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 관련 문의

코인빌리기 API를 이용할 때 자주 묻는 질문을 모았습니다. API별 상세 명세는 각 API 문서를, 서비스 정책은 <Anchor target="_blank" href="https://support.upbit.com/hc/ko/search?utf8=%E2%9C%93&query=%EC%BD%94%EC%9D%B8%EB%B9%8C%EB%A6%AC%EA%B8%B0">코인빌리기 이용자 가이드</Anchor> 및 [FAQ](https://www.upbit.com/lending/guide?tab=faq)를 참고해 주세요.

## 이용 자격과 권한

<Accordion title="코인빌리기 API를 사용하려면 API Key에 어떤 권한이 필요한가요?" icon="fa-info-circle">
  메인포켓 API Key에 코인빌리기(`lending`) 권한이 필요합니다. 서브포켓 API Key는 지원하지 않습니다.

  권한은 업비트 PC 웹의 <Anchor target="_blank" href="https://upbit.com/mypage/open_api_management">Open API 관리</Anchor>에서 API Key를 발급할 때 설정할 수 있으며, 권한 설정 과정에서 코인빌리기 이용 조건을 안내합니다. 권한을 발급받았더라도 실제 신청 시에는 이용 조건을 충족해야 합니다. 권한이 없는 API Key로 호출하면 HTTP 403 `out_of_scope` 에러가 반환됩니다.
</Accordion>

<Accordion title="질권 설정 계약에 동의해야 하나요?" icon="fa-info-circle">
  네. 코인빌리기는 맡긴 담보에 질권을 설정하므로, 코인빌리기 권한이 포함된 메인포켓 API Key를 새로 발급할 때마다 질권 설정 계약서 전문에 동의해야 합니다.

  질권 설정 계약서는 업비트 웹·앱에서 확인할 수 있습니다. 계약서 확인에 필요한 신청번호는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor>의 `lending_number`에서 확인해 주세요.
</Accordion>

<Accordion title="DAXA 교육을 이수하지 않아도 API를 사용할 수 있나요?" icon="fa-info-circle">
  네. DAXA 온라인 교육을 이수하지 않아도 메인포켓 API Key에 코인빌리기(`lending`) 권한을 설정하고 조회 API를 이용할 수 있습니다.

  신청 테스트, 신규 신청, 추가신청은 교육을 이수해야 이용할 수 있으며, 이수하지 않은 상태에서 호출하면 HTTP 403 `not_agreed_daxa_education` 에러가 반환됩니다.
</Accordion>

<Accordion title="약관 동의를 철회하면 어떤 API를 사용할 수 없나요?" icon="fa-info-circle">
  기존 내역을 확인하는 조회 API는 계속 이용할 수 있습니다. 가능 정보 조회, 신청 테스트, 신규 신청, 추가신청, 담보추가, 상환, 담보상환은 HTTP 403 `not_agreed_term` 에러가 반환됩니다.
</Accordion>

<Accordion title="고객 확인(KYC)이 만료되면 이용 중인 코인빌리기는 어떻게 되나요?" icon="fa-info-circle">
  고객 확인(KYC)이 만료될 경우, 이용 중인 코인빌리기는 강제상환 대상이 됩니다.
  포지션 조회는 계속 이용할 수 있으며, 현재 상태와 상환 결과는 개별 코인빌리기 조회에서 확인해 주세요.
</Accordion>

<Accordion title="이용 조건을 충족했는데 `insufficient_level` 에러가 반환됩니다." icon="fa-info-circle">
  조건을 충족한 직후에는 `insufficient_level`이 반환될 수 있습니다. 업비트 앱 또는 웹의 **코인빌리기 한도 조회 화면을 한 번 새로고침**한 뒤 다시 요청해 주세요. 같은 에러가 반복되면 고객센터로 문의해 주세요.
</Accordion>

<Accordion title="질권 설정 계약서는 어디에서 확인하나요?" icon="fa-info-circle">
  계약서는 API로 제공하지 않으며, 업비트 웹·앱에서 신청번호(`lending_number`)로 찾아 확인할 수 있습니다.

  1. 개별 코인빌리기 조회로 해당 포지션의 `lending_number`를 확인합니다.
  2. 업비트 웹·앱에서 같은 신청번호의 코인빌리기 건을 찾습니다.
  3. 해당 건에서 질권 설정 계약서를 확인합니다.

  `lending_number`는 웹·앱에 표시되는 신청번호와 같은 값입니다. 포지션을 식별하는 `uuid`와는 다르니 구분해 주세요.
</Accordion>

## 신청과 요청 처리

<Accordion title="`202 Accepted`가 반환되면 처리가 완료된 건가요?" icon="fa-info-circle">
  아니요. `202 Accepted`는 요청이 접수되었다는 뜻입니다. 응답의 `request_uuid`로 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending-request">코인빌리기 요청 상태 조회</Anchor>를 호출해 `state`가 `done`인지 확인해 주세요. `failed`라면 `failure.name`에서 실패 사유를 확인할 수 있습니다.
</Accordion>

<Accordion title="요청에 대한 응답을 받지 못했습니다. 다시 요청해도 되나요?" icon="fa-info-circle">
  바로 다시 요청하지 말고, 요청 시 지정한 `identifier`로 요청 상태를 먼저 조회해 주세요. 이미 접수된 요청을 새 `identifier`로 다시 보내면 별도 신청이나 상환이 접수될 수 있습니다.

  응답을 받지 못한 경우에도 조회할 수 있도록 신청·상환 요청에는 `identifier`를 지정하는 것을 권장합니다.
</Accordion>

<Accordion title="신청 테스트가 성공하면 실제 신청도 반드시 성공하나요?" icon="fa-info-circle">
  아니요. 신청 테스트는 테스트 요청 시점을 기준으로 검증합니다. 이후 시세나 신청 가능 수량이 바뀌면 실제 신청 결과도 달라질 수 있습니다. 테스트 응답의 `request_uuid`는 요청 상태 조회에 사용할 수 없습니다.
</Accordion>

<Accordion title="요청한 수량과 실제로 빌린 수량이 다를 수 있나요?" icon="fa-info-circle">
  네. 요청 수량보다 많이 빌려주는 경우는 없지만, 시세 변동에 따라 적게 확정되거나 거절될 수 있습니다.신청 처리 시점의 LTV(렌딩비율)가 최대 LTV를 소폭 넘으면 최대 LTV에 맞춰 수량을 줄여 처리하고, 허용 범위를 넘으면 `invalid_lending_ratio` 에러로 거절합니다.

  확정된 수량은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor>의 `borrowed_amount`에서 확인해 주세요.
</Accordion>

<Accordion title="신청 가능 수량(`available_budget_amount`)은 실시간 값인가요?" icon="fa-info-circle">
  준실시간 값입니다. 조회 시점과 실제 신청 시점의 수량이 다를 수 있으며, 신청 시 수량이 소진되었으면 `insufficient_lending_budget` 에러가 반환됩니다.
</Accordion>

## LTV와 강제상환

<Accordion title="LTV는 어떻게 계산하나요? 수수료도 포함되나요?" icon="fa-info-circle">
  LTV는 빌린 자산의 평가금액을 담보 인정금액으로 나눈 값입니다. 담보 인정금액은 담보 평가금액에 담보 인정 비율(`factor`)을 곱해 계산합니다. 빌린 자산은 남은 빌린 수량만 반영하며, 이용수수료는 포함하지 않습니다.

  평가금액은 업비트 평가가격을 기준으로 합니다. 자세한 가격 기준은 코인빌리기 가이드를 참고해 주세요.
</Accordion>

<Accordion title="강제상환은 언제 진행되나요?" icon="fa-info-circle">
  LTV가 강제상환 기준(`liquidation_ltv`)에 도달하거나 이용기간이 만료되면 강제상환이 진행됩니다. 코인빌리기 API는 LTV 변동 알림을 제공하지 않으므로, 이용 중에는 `current_ltv`를 주기적으로 확인해 주세요.
</Accordion>

<Accordion title="강제상환으로 종료된 포지션의 사유는 어디에서 확인하나요?" icon="fa-info-circle">
  개별 코인빌리기 조회 또는 종료된 코인빌리기 목록 조회의 `liquidation_reason`에서 확인할 수 있습니다. 담보상환으로 종료된 포지션은 `closure_type`이 `repay_with_collateral`이며, `liquidation_reason`은 `null`입니다.
</Accordion>

## 상환과 담보

<Accordion title="일부만 상환할 수 있나요?" icon="fa-info-circle">
  네. 코인빌리기 상환 API의 `repayment_amount`에 상환할 수량을 입력하면 일부 상환할 수 있습니다. 남은 수량을 모두 상환하려면 `repayment_mode`를 `full`로 지정해 주세요. 담보상환은 전액 상환만 지원합니다.
</Accordion>

<Accordion title="담보를 추가한 뒤에도 더 빌릴 수 있나요?" icon="fa-info-circle">
  아니요. 담보를 추가한 포지션은 추가신청할 수 없으며, `block_lending_increment_if_deposit_collateral` 에러가 반환됩니다. 더 빌리려면 신규 신청을 검토해 주세요.
</Accordion>

<Accordion title="담보상환 처리 내역은 어디에서 확인하나요?" icon="fa-info-circle">
  <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-lending-transactions">코인빌리기 거래 내역 조회</Anchor>를 `include_entries=true`로 요청하면 담보 매도, 빌린 자산 매수, 수수료 정산, 담보 반환 내역을 확인할 수 있습니다. 완료된 처리만 기록되며, 분할 체결의 회차별 내역은 제공하지 않습니다.
</Accordion>

## 수수료

<Accordion title="수수료율은 어디에서 확인하나요? 정책 요율이 바뀌면 기존 포지션에도 적용되나요?" icon="fa-info-circle">
  신청 전 현재 정책 요율은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-lending-information">코인빌리기 가능 정보 조회</Anchor>에서, 이미 생성된 포지션에 적용된 요율은 개별 코인빌리기 조회에서 확인할 수 있습니다.

  이용수수료율, 강제상환수수료율, 담보상환수수료율은 모두 최초 신청 시점의 요율이 해당 포지션에 계속 적용됩니다. 이후 정책 요율이 바뀌어도 기존 포지션에는 적용되지 않습니다.
</Accordion>

<Accordion title="이용수수료는 어떻게 계산되나요?" icon="fa-info-circle">
  신청 다음 날부터 하루 단위로, 그날 이용한 최대 수량에 일 단위 이용수수료율(`borrowing_fee_rate`)을 곱해 계산합니다. 당일 신청 후 당일 상환해도 1일치가 부과됩니다. 이용 중에는 차감하지 않고 최종 상환 시 정산합니다.
</Accordion>

<Accordion title="강제상환수수료와 담보상환수수료는 어떤 수량을 기준으로 계산하나요?" icon="fa-info-circle">
  처리 직전의 남은 빌린 수량을 기준으로 계산합니다. 예를 들어 100 USDT를 빌린 뒤 10 USDT를 상환하고 담보상환하면, 담보상환수수료율이 1%일 때 남은 90 USDT의 1%인 0.9 USDT가 부과됩니다.
</Accordion>

## 응답 값

<Accordion title="수량은 소수점 몇 자리까지 입력하고, 몇 자리로 반환되나요?" icon="fa-info-circle">
  요청에 입력하는 수량은 소수점 이하 최대 8자리까지 허용하며, 8자리를 초과하면 에러가 반환됩니다. 응답의 자산 수량은 소수점 이하 최대 20자리까지 반환될 수 있습니다. 거래 내역의 원화 환산 금액(`entries[].fiat_value`)은 소수점 이하 8자리로 반올림하여 반환합니다.
</Accordion>

# 같은 수준의 페이지

* [API 공통 문의](https://docs.upbit.com/kr/docs/faq-api.md)
* [주문 관련 문의](https://docs.upbit.com/kr/docs/faq-order.md)
* [입출금 관련 문의](https://docs.upbit.com/kr/docs/faq-withdraw-deposit.md)
* [Upbit Strategy Toolkit 관련 문의](https://docs.upbit.com/kr/docs/faq-upbit-strategy-toolkit.md)