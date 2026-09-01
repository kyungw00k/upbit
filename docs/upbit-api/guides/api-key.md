---
updatedAt: 2026-07-07T15:19:22.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# API Key 발급 받기

API Key를 발급받는 방법을 안내합니다.

## 시작하기

API Key는 사용자의 신원과 API 호출 권한을 확인하기 위해 사용되는 인증 정보입니다.

자산 조회, 주문, 입출금 관리와 같이 API를 호출하기 위해서는 반드시 사전에 API Key를 발급받고, 요청에 인증 정보를 포함해야 합니다. 본 가이드에서는 업비트 API Key 발급 방법과 관리 주의사항을 안내합니다.

<br />

## API Key 발급 정책 및 주의사항

* 업비트 API Key 발급은 **PC 웹 환경**에서만 가능합니다.
* 포켓당 발급 가능한 API Key는 아래와 같습니다.
  * 메인포켓에서 발급가능한 API Key 개수는 최대 10개입니다.
  * 서브포켓 당 발급 가능한  API Key 개수는 최대 5개입니다.
* 업비트 회원가입 후 고객 확인 및 2채널 인증(2FA)을 완료한 뒤에 API Key 발급이 가능합니다.

<br />

## API Key 발급 절차

### 1. API 관리 페이지로 이동

업비트 PC 웹에 로그인한 후, \[마이페이지 > <Anchor target="_blank" href="https://upbit.com/mypage/open_api_management">Open API 관리</Anchor>]로 이동합니다.

<br />

### 2. 생성할 API Key 권한 선택 및 허용 IP 주소 등록

API Key로 수행하고자 하는 API 기능 권한을 모두 선택합니다. 최소 하나 이상의 기능을 선택한 후 IP 주소 등록 란에 해당 API Key를 사용할 공개 IP 주소 목록을 입력합니다. 개인정보 수집 및 이용 동의란에 체크한 뒤 `Open API Key 발급받기` 버튼을 클릭하여 API Key 발급을 요청합니다.

#### 메인포켓

![](https://files.readme.io/3fb0fd456aea02e499f8a6b628e2e49b2fcc3cab642d04dc5d1d07506e1ee4b4-openapi_key01.png)

<Callout icon="fad fa-circle-exclamation" theme="warn">
  ### '출금하기' 권한 포함하여 API Key 발급 시 **출금 안심차단**이 완료되어야 출금이 가능합니다.

  출금하기 권한이 포함된 API Key는 발급 시 출금 안심차단이 자동으로 적용됩니다. **안심차단이 적용된 Key로는 출금 API를 호출할 수 없습니다.**

  출금 기능을 사용하려면 업비트 모바일 앱에서 안심차단을 해제해 주세요.&#x20;

  - \[더보기 > 인증/보안 > Open API 관리 > 활성] 페이지에서 해제할 Key 선택 후 활성화 되어 있는 출금 안심차단 토글을 비활성화 해주세요.
</Callout>

#### 서브포켓

![](https://files.readme.io/60acb732ebc3b1f503bd79006e87c787e351ce5315c5553f3a274694e2951d84-Frame_33.png)

<br />

### 3. API Key 발급을 위한 2채널 인증 수행

API Key 발급 인증을 위해 표시되는 인증 안내 팝업에서 사용하고자 하는 2채널 인증 수단을 선택합니다. 선택한 인증 수단을 통해 본인 인증을 완료한 후 API Key 발급이 완료됩니다. 2채널 인증 관련한 자세한 사항은 [업비트 고객 센터 > 자주하는 질문 > 2채널 인증](https://support.upbit.com/hc/ko/sections/900000896026-2%EC%B1%84%EB%84%90-%EC%9D%B8%EC%A6%9D) 문서를 참고해주시기 바랍니다.

<Image src="https://files.readme.io/95d8adb6dad32bcd8cb0ac89edb91d837633e9cd8b84b655642623d9612a6037-img2.png" align="center" />

<br />

### 4. API Key 발급

인증이 정상적으로 수행된 경우 API Key 발급이 완료되어 Access key와 Secret key를  확인할 수 있습니다. **Secret key는 최초 발급 화면에서만 확인 가능하므로 반드시 안전한 곳에 별도로 보관**하시기 바라며, 타인에게 노출 시 해킹의 위험에 놓일 수 있으니 노출되지 않도록 보안에 유의해주시길 바랍니다.

<br />

### 5. API Key 관리

\[변경] 버튼을 클릭하여 허용 IP 주소 목록을 수정할 수 있습니다.

\[삭제] 버튼을 클릭하여 발급된 API Key를 삭제할 수 있습니다. API Key 분실 혹은 유출 시 기존 키를 삭제하고 신규 키를 발급받아 사용하시기 바랍니다.

![](https://files.readme.io/b71a1082b8306302411b4f687c366b4248185bfe39b29e8d88d605acb153ff53-Frame_36.png)

## 마치며

본 가이드에서는 API 호출에 필요한 API Key를 발급받고 관리하는 방법을 알아보았습니다. 발급받은 API Key를 활용한 인증 방법은 아래 What's Next 영역을 참고하여 확인해주시기 바랍니다.

<br />

***

## 자주 들어오는 질문

**Q. 출금 안심차단은 어디서 해지 가능한가요?**

* 업비트 모바일 앱에 로그인 완료 후 \[인증/보안 > Open API 관리 > 활성] 페이지에서 차단을 해지할 Key를 선택해 안심차단을 해지해주세요.

**Q. 출금 안심차단 출금 차단 범위는 어디까지인가요?**

* 원화 및 디지털 자산 모두 출금이 불가능합니다.

**Q. 에러 해결 방법이 궁금해요?**

* `open_api_withdraw_locked`&#x20;
  * 해당 API Key는 출금 안심차단이 설정되어 있어 출금이 불가능합니다.
  * \[업비트 앱 > 더보기 > 인증/보안 > Open API 관리]에서 해당 키의 출금 안심차단을 해지해주세요.

- `expired_access_key`
  * API 키가 만료되었습니다. Open API 관리 페이지에서 만료된 키를 삭제하고, 신규키를 발급받길 바랍니다.
- `no_authorization_ip`
  * 등록되지 않은 클라이언트 IP에서 요청이 발생했습니다. 요청 환경을 확인한 후 허용 IP 주소를 변경(추가)하고 다시 요청해 주세요.

**Q. 등록 가능한 IP는 무엇인가요?**

* 허용 IP 주소는 외부 네트워크 통신에 이용되는 공개(공인)IP 주소를 입력해주셔야 합니다.
  * 개인 PC 환경에서는 사용중인 네트워크의 공개 IP 주소를 등록합니다.
  * 서버 환경에서는 해당 서버의 외부망 통신에 사용되는 IP 주소를 등록합니다.
* 사설 IP는 외부 접근이 불가능하므로 등록할 수 없습니다.
* 업비트 API는 유동 IP 환경에서의 이용을 권장하지 않습니다. 고정 IP(Static IP)를 사용하여 등록해주시기 바랍니다.
* 구글 등의 검색엔진에서 "what is my ip" 혹은 "내 IP 주소" 등을 검색하여 접속한 후 현재 사용중인 IP 주소를 확인할 수 있습니다.

<br />

**Q. Open API Key 발급 버튼이 활성화되지 않습니다.**

API Key 발급 조건이 만족되지 않은 경우 Open API Key 발급 버튼이 활성화되지 않습니다. 아래 조건을 모두 만족하였는지 확인해주시기 바랍니다.

* Open API에서 이용할 기능 권한 중 하나 이상의 기능을 선택
* ‘IP 주소 등록’창에 하나 이상의 IP를 입력
  * 등록하려는 IP 주소가 올바른 양식인지 확인 부탁드립니다.<br />(공백 포함 여부 확인, 각기 다른 IP 주소 간 구분이 아닌 한 IP주소 내에 콤마 포함 여부 확인)
  * IPv4형식의 IP 주소만 지원합니다.
* '개인정보 수집 및 이용 동의(필수)' 동의 체크박스를 클릭

<br />

**Q. 만료된 API Key는 어떻게 삭제하나요?**

1. 업비트 웹 PC에 로그인한 후 Open API 관리 페이지로 이동합니다.
2. 스크롤을 내려 페이지 하단의 ‘만료된 Key’ 버튼을 클릭합니다.
3. 삭제할 API Key의 체크박스를 클릭한 후 우측 상단의 삭제 버튼을 눌러 삭제합니다.

<br />

**Q. 등록된 API Key가 없는데 "최대 Open API 토큰 발급 수량을 초과했습니다"라는 안내가 뜹니다.**

해당 팝업이 노출되는 경우, \[만료된 Key] 탭으로 이동하여 조회되는 만료된 Key를 삭제하시고 다시 등록해보시길 바랍니다.

# Sibling pages

* [거래소 지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw_access_register.md)
* [개인지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw-private-wallet.md)