---
updatedAt: 2026-07-07T15:49:56.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Upbit API 이용 준비

Upbit API를 사용하기 위해 필요한 사전 준비 절차를 안내합니다.

## 사용 목적 정의

Upbit API 사용 전, 사용 목적을 먼저 확인하세요.

사용 목적에 따라 API Key 필요 여부와 권한(Scope)이 달라집니다.<br />필요한 기능에 맞게 최소 권한만 설정하는 것을 권장합니다.<br />API Key는 Upbit PC Web의 [Open API 관리 페이지](https://upbit.com/mypage/open_api_management)에서 발급할 수 있습니다.

<br />

## 사용 목적별 권한 안내

<Table align={["left",null,null,"left","left"]}>
  <thead>
    <tr>
      <th>
        사용 목적
      </th>

      <th>
        메인포켓<br />(API Key)
      </th>

      <th>
        서브포켓<br />(API Key)
      </th>

      <th>
        필요 권한
      </th>

      <th>
        비고
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        [시세 및 마켓 데이터 조회](https://docs.upbit.com/kr/reference/list-trading-pairs)
      </td>

      <td>
        불필요
      </td>

      <td>
        불필요
      </td>

      <td>

      </td>

      <td>
        Public API로 조회 가능
      </td>
    </tr>

    <tr>
      <td>
        [과거 마켓 데이터 수집](https://www.upbit.com/historical_data/main)
      </td>

      <td>
        불필요
      </td>

      <td>
        불필요
      </td>

      <td>

      </td>

      <td>
        Historical Market Data 문서 참고
      </td>
    </tr>

    <tr>
      <td>
        메인 포켓 정보 조회
      </td>

      <td>
        필요
      </td>

      <td>
        필요
      </td>

      <td>
        - 메인포켓: 포켓관리
        - 서브포켓: 자산이전
      </td>

      <td>
        - 메인포켓 :&#x20;
          - <Anchor target="_blank" href="ref: get-pocket-information">포켓 정보 조회</Anchor>
          - <Anchor target="_blank" href="ref: get-pocket-api-keys">포켓별 API Key 목록 조회</Anchor>
          - <Anchor target="_blank" href="ref: get-sub-pocket-balance">서브포켓 잔고 조회</Anchor>
          - <Anchor target="_blank" href="ref: get-universal-transfer">메인포켓 자산 이전 목록 조회</Anchor>
        - 서브포켓 :&#x20;
          - <Anchor target="_blank" href="ref: get-transfer">서브포켓 자산 이전 목록 조회</Anchor>
      </td>
    </tr>

    <tr>
      <td>
        포켓 간 자산 이전
      </td>

      <td>
        필요
      </td>

      <td>
        필요
      </td>

      <td>
        - 메인포켓: 포켓관리
        - 서브포켓: 자산이전
      </td>

      <td>
        - 메인포켓 :&#x20;
          - <Anchor target="_blank" href="ref : post-universal-transfer">메인포켓 자산 이전</Anchor>
        - 서브포켓 :&#x20;
          - <Anchor target="_blank" href="ref: post-transfer">서브포켓 자산 이전</Anchor>
      </td>
    </tr>

    <tr>
      <td>
        [자산 조회 및 관리](https://docs.upbit.com/kr/reference/get-balance)
      </td>

      <td>
        필요
      </td>

      <td>
        필요
      </td>

      <td>
        자산조회
      </td>

      <td>
        잔고 및 보유 자산 조회
      </td>
    </tr>

    <tr>
      <td>
        [주문](https://docs.upbit.com/kr/reference/available-order-information)
      </td>

      <td>
        필요
      </td>

      <td>
        필요
      </td>

      <td>
        주문조회, 주문하기
      </td>

      <td>
        주문 API 사용
      </td>
    </tr>

    <tr>
      <td>
        [입금](https://docs.upbit.com/kr/reference/available-deposit-information)
      </td>

      <td>
        필요
      </td>

      <td>
        미지원
      </td>

      <td>
        입금조회, 입금하기
      </td>

      <td>
        입금 및 조회 가능
      </td>
    </tr>

    <tr>
      <td>
        [출금](https://docs.upbit.com/kr/reference/available-withdrawal-information)
      </td>

      <td>
        필요
      </td>

      <td>
        미지원
      </td>

      <td>
        출금조회, 출금하기
      </td>

      <td>
        - 디지털 자산 출금 시 허용 주소 등록 필요<br />[거래소](https://docs.upbit.com/kr/docs/open-api-withdraw_access_register) / [개인지갑](https://docs.upbit.com/kr/docs/open-api-withdraw-private-wallet)
        - 출금 안심차단 1회<br />해지 필요
      </td>
    </tr>
  </tbody>
</Table>

<br />

<Callout icon="⚠️" theme="warn">
  ### 오류 해결

  권한 또는 출금허용주소가 사전에 설정되지 않은 경우 API 호출 시 오류가 발생할 수 있습니다.

  - **401 Unauthorized(out\_of\_scope)**
    - API Key에 필요한 권한이 포함되어 있지 않은 경우 발생합니다. 필요한 권한 선택 후 새로 Key를 발급 받길 바랍니다.
  - **400 Bad Request(withdraw\_address\_not\_registered)**
    - 출금허용주소가 등록되지 않은 경우 발생합니다. 출금주소를 등록하시길 바랍니다.
  - **403 Forbidden (open\_api\_withdraw\_locked)**
    - 해당 API Key는 출금 안심차단이 설정되어 있어 출금이 불가능합니다.<br />\[업비트 앱 > 더보기 > 인증/보안 > Open API 관리]에서 해당 키의 출금 안심차단을 해지해주세요.
</Callout>

<br />

## 다음 단계

* [API Key 발급 받기](https://docs.upbit.com/kr/docs/api-key)
* [거래소 지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw_access_register)
* [개인지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw-private-wallet)
* [개발환경 설정](https://docs.upbit.com/kr/docs/dev-environment)

<br />

# Sub pages

* [API Key 발급 받기](https://docs.upbit.com/kr/docs/api-key.md)
* [거래소 지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw_access_register.md)
* [개인지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw-private-wallet.md)

# Sibling pages

* [개발환경 설정](https://docs.upbit.com/kr/docs/dev-environment.md)
* [Public API 호출](https://docs.upbit.com/kr/docs/first-quotation-api-call.md)
* [인증 API 호출](https://docs.upbit.com/kr/docs/first-exchange-api-call.md)