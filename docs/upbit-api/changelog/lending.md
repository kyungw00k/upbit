---
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# [안내] 코인빌리기(Lending) API 지원 안내

안녕하세요. 업비트 개발자센터입니다.

업비트 코인빌리기(Lending) 서비스를 API로 이용할 수 있도록 신규 기능을 지원합니다.

**코인빌리기**는 보유한 원화(KRW) 또는 지원되는 디지털 자산을 담보로 맡기고 디지털 자산을 빌리는 서비스입니다.

이번 업데이트를 통해 코인빌리기 가능 정보 조회부터 신청, 추가신청, 담보추가, 상환 및 이용 내역 조회까지 API로 이용할 수 있습니다.

자세한 내용은 아래를 참고해 주세요.

## 1. 적용 일시

* 적용일: 2026-11-26 (목)

## 2. 신규 기능 내용

* 기능명: 코인빌리기(Lending) API
* 지원 대상: 업비트 대한민국(KR) 개인회원
* 지원 기능: 코인빌리기 관련 API 총 12종

| 구분      | API                                                                                 | Endpoint                                  |
| ------- | ----------------------------------------------------------------------------------- | ----------------------------------------- |
| 이용 전 확인 | [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information) | `GET /v1/lendings/chance`                 |
| 이용 전 확인 | [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test)                    | `POST /v1/lendings/test`                  |
| 신청      | [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending)                         | `POST /v1/lendings`                       |
| 신청 및 변경 | [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more)                      | `POST /v1/lendings/borrow_more`           |
| 신청 및 변경 | [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral)                      | `POST /v1/lendings/adjust_collateral`     |
| 상환      | [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending)                       | `POST /v1/lendings/repay`                 |
| 상환      | [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral)         | `POST /v1/lendings/repay_with_collateral` |
| 조회      | [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)           | `GET /v1/lendings/request`                |
| 조회      | [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)                      | `GET /v1/lending`                         |
| 조회      | [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings)       | `GET /v1/lendings/ongoing`                |
| 조회      | [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings)         | `GET /v1/lendings/closed`                 |
| 조회      | [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions)     | `GET /v1/lendings/transactions`           |

* **API Key 권한**
  * 신규 권한: 코인빌리기(`lending`)
  * 지원 포켓: 메인포켓
  * 인증 방식: API Key를 이용한 JWT 인증

* **요청 수 제한(Rate Limit)**
  | 그룹             | 초당 최대 요청 횟수 | 대상 API                                      |
  | -------------- | ----------- | ------------------------------------------- |
  | `lending`      | 10회         | 신청, 추가신청, 담보추가, 상환, 담보상환, 개별 코인빌리기 조회       |
  | `lending-test` | 10회         | 코인빌리기 신청 테스트                                |
  | `default`      | 30회         | 가능 정보 조회, 요청 상태 조회, 진행 중·종료 목록 조회, 거래 내역 조회 |

※ 요청 수 제한은 포켓 단위로 적용되며, 동일한 그룹에 속한 API는 요청 가능 횟수를 공유합니다.

## 3. 이용 방법

* 업비트 PC 웹 API Key 관리 화면에서 메인포켓을 선택하고, 코인빌리기(`lending`) 권한을 추가합니다.
  * 권한 추가 시 표시되는 **\[필수] 질권설정 계약 동의** 내용을 확인한 &#xD6C4;**&#x20;\[동의하고 권한 추가]**&#xB97C; 선택합니다.
  * API Key 발급에 필요한 인증 및 등록 절차를 완료합니다.
* <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-lending-information">코인빌리기 가능 정보 조회</Anchor> API를 통해 이용 가능한 자산, 담보 조건, 대여한도 및 수수료율 등을 확인합니다.
* <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-test">신청 테스트</Anchor> API를 이용하여 요청 조건과 예상 LTV를 확인한 후, <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/new-lending">코인빌리기 신청</Anchor> API로 실제 신청을 진행합니다.
* 이용 중인 코인빌리기는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/borrow-more">추가신청</Anchor>, <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/adjust-collateral">담보추가</Anchor>, [상환](https://docs.upbit.com/kr/reference/repay-lending) 및 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/repay-lending-with-collateral">담보상환</Anchor> API를 통해 관리할 수 있습니다.
* 신청, 추가신청, 담보추가, 상환, 담보상환 요청은 비동기로 처리됩니다.
  * 요청이 접수되면 `202 Accepted` 응답이 반환됩니다.
  * 실제 처리 결과는 `request_uuid` 또는 요청 시 지정한 `identifier`를 이용해 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending-request">코인빌리기 요청 상태 조회</Anchor> API에서 확인합니다.

자세한 이용 조건과 API별 요청·응답 형식은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-overview">코인빌리기(Lending) 개요</Anchor> 및 각 API 레퍼런스 문서를 참고해 주세요.

## 유의사항

**이용 조건**

* 코인빌리기 API는 메인포켓 API Key로만 이용할 수 있으며, 서브포켓 API Key는 지원하지 않습니다.
* 실제 신청은 고객확인(KYC), 약관 동의, 보유 자산, 매매 이력, 회원 레벨 등 서비스 이용 조건을 충족한 경우에만 가능합니다.
* 신청 테스트, 신규 신청, 추가신청은 DAXA 온라인 교육을 이수해야 이용할 수 있으며, 미이수 시 `not_agreed_daxa_education` 에러가 반환됩니다. 코인빌리기(`lending`) 권한 설정과 조회 API는 교육 이수 여부와 관계없이 이용할 수 있습니다.
* 코인빌리기를 이용하면 맡긴 담보에 질권이 설정됩니다. 질권설정 계약서는 업비트 웹·앱에서 신청번호(`lending_number`)로 확인할 수 있습니다.

**요청 처리**

* `202 Accepted`는 요청이 접수되었다는 의미이며, 처리가 완료되었음을 보장하지 않습니다. 최종 처리 결과는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending-request">코인빌리기 요청 상태 조회</Anchor> API로 확인해 주세요.
* 요청 후 응답을 받지 못했다면 같은 요청을 바로 다시 보내지 말고, 먼저 요청 시 지정한 `identifier`로 요청 상태를 조회해 주세요. 새로운 `identifier`로 다시 요청하면 별도의 요청으로 접수되어 신청 또는 상환이 중복될 수 있습니다.

**LTV 관리**

* 코인빌리기 API는 LTV 변동 알림을 제공하지 않습니다. LTV가 강제상환 기준에 도달하면 강제상환이 진행될 수 있으므로, 이용 중에는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor> API의 `current_ltv`를 주기적으로 확인해 주세요. `current_ltv`는 포지션 상태가 `active`인 경우에만 제공됩니다.

**요청 수 제한**

* 개별 코인빌리기 조회는 신청·상환과 같은 `lending` 요청 수 제한 그룹을 사용합니다. 짧은 주기로 반복 조회하면 신청과 상환에 사용할 수 있는 요청 횟수가 줄어들 수 있으니 유의해 주세요.
* 요청 수 제한을 지키더라도 이용 자격, 대여한도, 담보 조건 및 서비스 운영 상황에 따라 요청이 제한되거나 처리되지 않을 수 있습니다.

***

감사합니다.

업비트 개발자센터 드림