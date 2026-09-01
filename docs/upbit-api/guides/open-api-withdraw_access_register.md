---
updatedAt: 2026-06-01T02:38:44.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 거래소 지갑 주소 등록

타 거래소에서 발급 받은 입금 주소를 업비트의 출금 허용 주소로 등록하는 방법을 안내합니다.

## 출금허용주소 등록

업비트 PC 웹에 로그인한 후, \[마이페이지 > [Open API 관리](https://upbit.com/mypage/open_api_management)]에서 \[디지털 자산 출금주소 관리] 탭으로 이동합니다. 이동 후 \[출금허용주소 등록하기 → 주소 추가] 버튼을 클릭합니다.

![](https://files.readme.io/3c5fd82f080400b5059f101a2345570b1027703f94e157c95cd739e49bec2c6e-Frame_29.png)

<br />

## 거래소 지갑 주소 등록

Open API를 통해 거래소 지갑 주소로 출금하려면, 해당 주소가 출금허용주소목록에 사전 등록되어 있어야 합니다.

### 개인

개인에게 출금할 경우 아래의 항목을 확인해 주세요. 주소 소유주와 '받는 사람' 항목의 이름이 다를 경우 인증에서 실패할 수 있습니다. 등록 전 이름이 일치하는 지 확인 바랍니다.

| 항목     | 선택                          | 예시                        |
| :----- | :-------------------------- | :------------------------ |
| 지갑종류   | 거래소                         | 바이낸스, 바이비트, OKX 등         |
| 출금허용주소 | 출금할 \[네트워크 - 디지털 자산]의 입금 주소 |                           |
| 출금네트워크 | 네트워크                        | Bitcoin, Ethereum, Tron 등 |
| 디지털자산  | 디지털 자산                      | BTC,, ETH, USDT 등         |
| 받는 사람  | 수취인 실명                      |                           |

<Image src="https://files.readme.io/5637004af1234f1c5dc8cfa183e41aee23d5bea9ac7365786b599dc9129fc49d-Frame_6.png" align="center" />

<br />

### 법인

법인에게 출금할 경우 \[지갑종류\~디지털자산] 의 등록 방법은 개인과 동일합니다.

단, '받는 사람'의 항목을 \[법인]으로 선택해 주세요.

받는 사람 이름 입력창에 \[법인 대표명]과 \[법인명]을 정확하게 입력해주세요. 주소 소유주(법인대표 이름)와 '받는 사람' 항목의 이름이 다를 경우 인증에서 실패할 수 있습니다. 등록 전 이름이 일치하는 지 확인 바랍니다.

<Image src="https://files.readme.io/974ce32e1daadc6479337337af71215562bd344375c50144056f7996a33f3f97-Frame_7.png" align="center" />

<br />

## 인증 및 출금 주소 등록 완료

주소 등록 페이지에서 \[확인] 버튼을 클릭한 후, ‘인증 대기’ 화면에서 인증 대기 중인 주소를 확인할 수 있습니다.<br />오른쪽 상단의 \[2채널 인증하기] 버튼을 클릭해 인증할 항목을 선택해 주세요. 인증이 완료되면 ‘디지털 자산 출금 주소 등록 성공’ 팝업이 표시됩니다.<br />법인 주소인 경우, ‘받는 사람’ 항목에는 법인대표명(법인명) 이 표시됩니다.

※ 주소 인증이 실패했을 경우 상대 거래소 해당 디지털 자산의 입금주소가 발급되었는지 확인바랍니다.

<Image src="https://files.readme.io/8d4fd6d0fd44a1cdb9478ae4f30e5b11ac766d602cae7f08671bc6eba375659c-Frame_14.png" align="center" />

<br />

## 출금허용주소 삭제

주소 관리는 [디지털 자산 출금 주소 관리](https://upbit.com/mypage/open_api_management?tab=fund_source) 메뉴에서 가능합니다.<br />삭제하려는 주소를 선택한 후, 좌측 상단의 \[삭제] 버튼을 클릭하고 2채널 인증을 완료하면, 해당 출금 주소가 삭제 처리됩니다.

<Image src="https://files.readme.io/47351be98bff4a936350870b025989a92c308820f927326c17167e25e7740c7e-Frame_28.png" align="center" />

<br />

# Sibling pages

* [API Key 발급 받기](https://docs.upbit.com/kr/docs/api-key.md)
* [개인지갑 주소 등록](https://docs.upbit.com/kr/docs/open-api-withdraw-private-wallet.md)