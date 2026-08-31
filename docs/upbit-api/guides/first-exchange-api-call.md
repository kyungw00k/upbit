---
updatedAt: 2026-07-14T09:05:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 인증 API 호출

발급받은 API Key의 Access Key와 Secret Key를 사용하여 인증 토큰을 생성하고 Exchange API를 호출할 수 있습니다.

## 시작하기

본 가이드에서는 cURL 요청과 함께 활용할 수 있는 Shell 스크립트 기반 인증 토큰 생성방법 및 인증 API 호출 방법을 알아봅니다.

<br />

## 인증 방식 알아보기

[인증](https://docs.upbit.com/kr/reference/auth) 문서에서 업비트의 REST API와 WebSocket에 사용되는 인증 토큰의 구조 및 생성 방법을 참고할 수 있습니다. 본 가이드 진행 시 토큰의 구조와 생성 방식에 대한 이해가 있는 경우 더욱 원활한 진행이 가능하므로, 필요에 따라 해당 내용을 먼저 참고해주시기 바랍니다.

<br />

## 인증 토큰 생성을 위한 Shell 스크립트

cURL은 매우 강력한 도구이지만, 해시나 파라미터 처리와 같은 기능은 내장하고 있지 않아 자체적으로 인증 토큰을 생성하기 어렵습니다. 본 가이드에서는 Shell 스크립트를 사용한 인증 토큰 생성 방법을 안내합니다.

Windows 환경에서는 PowerShell이나 CMD에서 직접 Shell 스크립트를 실행할 수 없습니다. Git Bash 또는 WSL(Windows Subsystem for Linux)와 같은 Bash 환경을 설치한 후 실행하세요. 아래 환경 설정 가이드 중 하나를 참고하여 Bash 환경을 구성해주세요.

  <div className="accordion">
    <input type="checkbox" id="gitbash" />
    <label for="gitbash">
      <i className="fa-solid fa-circle-info"></i> Git Bash 설치 방법
    </label>
    <div className="accordion-content">
      <a href="https://gitforwindows.org/">Git for Windows 공식 다운로드 페이지</a>에서 Git을 설치하세요. 설치 완료 후, Git Bash라는 프로그램이 시작 메뉴에 추가됩니다.
    </div>
	</div>

	<div className="accordion">
    <input type="checkbox" id="wsl" />
    <label for="wsl">
      <i className="fa-solid fa-circle-info"></i> WSL(Windows Subsystem for Linux) 설치 방법
    </label>
    <div className="accordion-content">
      PowerShell을 관리자 권한으로 실행하고 아래 명령어를 실행하여 WSL을 설치합니다. 
      <pre><code>wsl --install</code></pre>
    </div>
	</div>

Bash 환경 구성이 완료되면, Bash 명령어를 사용하여 아래와 같이 Shell 스크립트를 실행할 수 있습니다.

```shell
bash auth_curl.sh "curl --request GET --url 'https://api.upbit.com/v1/accounts' --header 'accept: application/json'"
```

<br />

Shell 스크립트 작성 및 실행에 앞서, json 처리 및 hash등 암호화 관련 처리를 위해 사용할 `jq`와 `openssl` 유틸리티가 PC에 설치되어 있어야 합니다. 아래 OS별 설치 여부 확인 및 설치 가이드를 참고하여 실행 환경을 확인해주세요.

 <div className="accordion">
    <input type="checkbox" id="macos-openssl" />
    <label for="macos-openssl">
      <i className="fa-solid fa-circle-info"></i> macOS jq, openssl 설치 확인 방법
    </label>
    <div className="accordion-content">
      macOS의 경우 대부분 jq와 openssl이 기본적으로 설치되어있으나, 아래 명령어를 실행하여 설치 여부를 확인할 수 있습니다.<br /><br />

      <pre><code>jq --version</code><code>openssl --version</code></pre>

버전 정보가 출력되면 정상적으로 사용이 가능한 상태입니다. 사용자 PC 환경에 따라 버전은 상이할 수 있습니다.<br />
설치되어있지 않은 경우, 아래 명령어를 입력하어 Homebrew를 통해 설치할 수 있습니다.

      <pre><code>brew install jq</code><code>brew install openssl</code></pre>

    </div>
</div>

<div className="accordion">
    <input type="checkbox" id="windows-openssl" />
    <label for="windows-openssl">
      <i className="fa-solid fa-circle-info"></i> Windows jq, openssl 설치 확인 방법
    </label>
    <div className="accordion-content">
      Windows cmd 또는 powerShell에서 아래 명령어를 실행하여 설치 여부를 확인할 수 있습니다.<br /><br />

      <pre><code>jq --version</code><code>openssl --version</code></pre>

버전 정보가 출력되면 정상적으로 사용이 가능한 상태입니다. 사용자 PC 환경에 따라 버전은 상이할 수 있습니다.<br />
설치되어있지 않은 경우, 아래와 가이드를 따라 설치할 수 있습니다. <br /><br />

<code>jq 설치</code>

      <ol>
        <li>https://stedolan.github.io/jq/download/ 에서 Windows용 jq 바이너리 다운로드</li>
        <li>원하는 폴더에 저장한 후, 해당 폴더를 환경변수 PATH에 추가</li>
      </ol>

<code>openssl 설치</code>

      <ol>
        <li>https://slproweb.com/products/Win32OpenSSL.html 에서 Windows용 OpenSSL 설치파일 다운로드	
        <li>설치 후, C:\Program Files\OpenSSL-Win64\bin 등의 경로를 환경변수 PATH에 추가</li>
      </ol>

      <br />	
			또는, WSL 환경 사용시 아래 설치 명령어를 통해 jq, openssl을 일괄 설치할 수 있습니다.
      <pre><code>sudo apt update<br />
sudo apt install jq openssl</code></pre>
    </div>
  </div>

<br />

아래 Shell 스크립트 예제를 복사한 뒤 `auth_curl.sh` 과 같은 이름의 파일로 생성합니다. **ACCESS\_KEY와 SECRET\_KEY에 발급받은 API Key 값으로 변경한 뒤 파일을 저장**합니다. 아래 스크립트는 업비트 인증 토큰을 생성하는 예시 코드로 다음과 같은 역할을 수행합니다.

* 스크립트 실행 인자로 호출하고자 하는 cURL 명령을 전달하면, 필요한 파라미터와 API Key를 사용하여 인증 헤더를 생성합니다.
* 인증 헤더를 포함한 cURL을 생성하여 출력합니다. 해당 cURL 명령을 복사하여 실행하여 인증 API를 호출할 수 있습니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> Secret Key 보안 안내
      </div>
    API Key는 안전하게 저장 및 관리되어야 하므로, Secret Key를 포함한 문구 또는 스크립트를 타인에게 전송하거나 온라인에 업로드하지 않도록 주의하시기 바랍니다. 반드시 테스트를 위한 목적으로 로컬 환경에서만 사용하시기 바랍니다.
  </div>

```shell
#!/bin/bash

ACCESS_KEY="YOUR_ACCESS_KEY"
SECRET_KEY="YOUR_SECRET_KEY"

INPUT_CURL=$1

METHOD=$(echo "$INPUT_CURL" | sed -n "s/.*-X \([A-Z]*\).*/\1/p")
URL=$(echo "$INPUT_CURL" | sed -n "s/.*--url '\([^']*\)'.*/\1/p")

BODY=$(echo "$INPUT_CURL" | sed -n "s/.*--data '\([^']*\)'.*/\1/p")

if [[ "$METHOD" == "POST" ]]; then
  QUERY_STRING=$(echo "$BODY" | tr -d '{}' | sed 's/"//g' | tr ',' '&' | tr ':' '=' | sed 's/ //g')
elif [[ "$METHOD" == "GET" || "$METHOD" == "DELETE" ]]; then
  QUERY_STRING=$(echo "$URL" | awk -F '?' '{print $2}')
fi

if [ -z "$QUERY_STRING" ]; then
  QUERY_HASH=""
else
  QUERY_HASH=$(echo -n "$QUERY_STRING" | openssl dgst -sha512 | sed 's/^.* //')
fi

NONCE=$(uuidgen)
HEADER='{"alg":"HS512","typ":"JWT"}'
PAYLOAD=$(jq -n --arg ak "$ACCESS_KEY" --arg n "$NONCE" --arg qh "$QUERY_HASH" --arg alg "SHA512" \
  '{access_key: $ak, nonce: $n, query_hash: $qh, query_hash_alg: $alg}')

HEADER_BASE64=$(echo -n "$HEADER" | openssl base64 -A | tr '+/' '-_' | tr -d '=')
PAYLOAD_BASE64=$(echo -n "$PAYLOAD" | openssl base64 -A | tr '+/' '-_' | tr -d '=')

SIGNATURE=$(echo -n "$HEADER_BASE64.$PAYLOAD_BASE64" | \
  openssl dgst -sha512 -hmac "$SECRET_KEY" -binary | \
  openssl base64 -A | tr '+/' '-_' | tr -d '=')

JWT="$HEADER_BASE64.$PAYLOAD_BASE64.$SIGNATURE"

CLEANED_CURL=$(echo "$INPUT_CURL" | sed "s/-H 'Authorization: Bearer [^']*'//g")

if echo "$CLEANED_CURL" | grep -q "Content-Type"; then
    FINAL_CURL=$(echo "$CLEANED_CURL" | sed "s/curl /curl -H 'Authorization: Bearer $JWT' /")
else
    FINAL_CURL=$(echo "$CLEANED_CURL" | sed "s/curl /curl -H 'Authorization: Bearer $JWT' -H 'Content-Type: application\/json' /")
fi

echo "[+] Signed curl command:"
echo "$FINAL_CURL"
```

인증 토큰을 생성하는 과정에서 jq, openssl 등의 유틸리티를 사용하기 때문에, macOS 및 기타 UNIX 계열 환경에서는 아래와 같이 명령어를 입력하여 실행 권한을 추가해야 합니다.

```shell
chmod +x auth_curl.sh
```

<br />

### 잔고 조회 하기

인증이 필요한 첫 API로 <Anchor target="_blank" href="ref:get-balance">포켓 잔고 조회</Anchor> API를 호출해보겠습니다. 위에서 저장한 Shell 스크립트의 인자로 해당 API의 cURL 명령을 전달하여 실행합니다. 해당 API의 cURL 예제는 Reference 우측 예제 코드 영역에서 확인할 수 있습니다.

```shell
./auth_curl.sh "curl --request GET --url https://api.upbit.com/v1/accounts --header 'accept: application/json'"
[+] Signed curl command:
curl -H 'Authorization: Bearer eyJhbGciO...TvcqiPUoNTg' --request GET --url https://api.upbit.com/v1/accounts --header 'accept: application/json'
```

스크립트의 실행 결과로 인증 헤더를 포함한 cURL 명령 전체가 출력됩니다. 전체 요청을 복사 후 붙여넣기로 실행하여 아래와 같은 API 응답을 확인할 수 있습니다.

```json
[{"currency": "KRW","balance": "1000000.0","locked": "0.0","avg_buy_price": "0","avg_buy_price_modified": false,"unit_currency": "KRW"},{"currency": "BTC","balance": "2.0","locked": "0.0","avg_buy_price": "140000000","avg_buy_price_modified": false,"unit_currency": "KRW"}]
```

<br />

### 매수 주문 생성하기

이번에는 POST API를 요청해보겠습니다. KRW-BTC 마켓에서 체결되지 않을 만한 충분히 낮은 금액의 지정가 매수 주문을 생성하여 테스트해볼 수 있습니다. Shell 스크립트에 아래와 같이 [주문 생성](https://docs.upbit.com/kr/reference/new-order)을 위한 cURL 요청을 인자로 호출하면, 아래와 같이 인증 토큰이 포함된 요청이 반환됩니다. (인증 토큰 부분은 `...`으로 일부 요약 처리 되어있습니다.)

```shell
./auth_curl.sh "curl -X POST --url 'https://api.upbit.com/v1/orders' --header 'accept: application/json' --header 'content-type: application/json' --data '{\"market\":\"KRW-BTC\",\"side\":\"bid\",\"volume\":\"0.0001\",\"price\":\"50000000\",\"ord_type\":\"limit\"}'"


[+] Signed curl command:
curl -H 'Authorization: Bearer eyJhbGciOi...a0uejONDLQ' -X POST --url 'https://api.upbit.com/v1/orders' --header 'accept: application/json' --header 'content-type: application/json' --data '{"market":"KRW-BTC","side":"bid","volume":"0.0001","price":"50000000","ord_type":"limit"}'

```

출력된 결과 cURL 커맨드를 복사하여 터미널에 입력한 뒤 실행하면, KRW-BTC 페어에 주문이 생성됩니다.

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> 위 cURL 요청 실행시 실제 주문이 생성됩니다. 
      </div>
    디지털 화폐의 급격한 가격 변동으로 해당 주문이 체결될 수 있으므로, 실제 주문을 생성 테스트하는 경우에만 실행해주시기 바랍니다. 미체결 주문은 업비트 서비스 또는 주문 취소 API를 호출하여 취소할 수 있습니다.
  </div>

<br />

## 인증 토큰 생성 방식 이해하기

제시된 Shell 스크립트의 아래 코드에 의해 JWT가 생성됩니다.

```shell
if [ -z "$QUERY_STRING" ]; then
  QUERY_HASH=""
else
  QUERY_HASH=$(echo -n "$QUERY_STRING" | openssl dgst -sha512 | sed 's/^.* //')
fi

NONCE=$(uuidgen)
HEADER='{"alg":"HS512","typ":"JWT"}'
PAYLOAD=$(jq -n --arg ak "$ACCESS_KEY" --arg n "$NONCE" --arg qh "$QUERY_HASH" --arg alg "SHA512" \
  '{access_key: $ak, nonce: $n, query_hash: $qh, query_hash_alg: $alg}')

HEADER_BASE64=$(echo -n "$HEADER" | openssl base64 -A | tr '+/' '-_' | tr -d '=')
PAYLOAD_BASE64=$(echo -n "$PAYLOAD" | openssl base64 -A | tr '+/' '-_' | tr -d '=')

SIGNATURE=$(echo -n "$HEADER_BASE64.$PAYLOAD_BASE64" | \
  openssl dgst -sha512 -hmac "$SECRET_KEY" -binary | \
  openssl base64 -A | tr '+/' '-_' | tr -d '=')

JWT="$HEADER_BASE64.$PAYLOAD_BASE64.$SIGNATURE"
```

위 코드는 입력한 요청의 쿼리 문자열을 Hash한 후, JWT 헤더, 페이로드, 서명을 생성하는 예시 코드에 해당합니다. 페이로드 생성 시 업비트 인증 토큰 규격에 따라 `access_key`, `nonce`, `query_hash`, `query_hash_alg` 필드가 사용되는 부분을 확인하실 수 있습니다.

본 요청에 쿼리 문자열이 존재하지 않는 POST 요청 인증 토큰 생성 과정은 어떻게 진행될까요? Shell 스크립트 실행시 아래  코드에 의해 POST 요청의 JSON 본문이 쿼리 문자열 형태로 변경됩니다. 이후 해당 쿼리 문자열을 Hash하여 GET 요청과 동일한 방식으로 인증 토큰이 생성됩니다.

[block:html]
{
  "html": "<div class=\"callout-section\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i>  URL 인코딩이란?\n      </div>\nURL 인코딩은 통신 프로토콜에서 URL 내에 포함할 수 없는 문자를 전송 가능한 문자로 변환하는 인코딩 방식입니다. 특수 문자를 포함한 대상 문자들은 인코딩 시 % 기호와 2자리 16진수로 이루어진 문자열로 변환됩니다.<br><br>\n(예시) <code>:</code>은 인코딩 후 <code>%3A</code>로, <code>+</code>는 <code>%2B</code>로 변환\n  </div>"
}
[/block]

```shell
BODY=$(echo "$INPUT_CURL" | sed -n "s/.*--data '\([^']*\)'.*/\1/p")

if [[ "$METHOD" == "POST" ]]; then
  QUERY_STRING=$(echo "$BODY" | tr -d '{}' | sed 's/"//g' | tr ',' '&' | tr ':' '=' | sed 's/ //g')
elif [[ "$METHOD" == "GET" || "$METHOD" == "DELETE" ]]; then
  QUERY_STRING=$(echo "$URL" | awk -F '?' '{print $2}')
fi
```

<br />

## 마치며

cURL을 활용하여 쉽고 빠르게 업비트 API를 호출해보았습니다. 충분한 테스트를 통해 업비트 API를 이해하셨다면, 이제 실제로 매매 전략을 구현하기 위한 개발환경을 설정해보세요. [개발 환경 설정 가이드](https://docs.upbit.com/kr/docs/dev-environment) 로 이동하여 cURL이 아닌 언어별 개발 환경을 설정하거나, [REST API 연동 Best Practice](https://docs.upbit.com/kr/docs/rest-api-best-practice) 로 이동하여 보다 고도화된 REST API 연동 방식을 확인할 수 있습니다.

# Sibling pages

* [Upbit API 이용 준비](https://docs.upbit.com/kr/docs/api-setup.md)
* [개발환경 설정](https://docs.upbit.com/kr/docs/dev-environment.md)
* [Public API 호출](https://docs.upbit.com/kr/docs/first-quotation-api-call.md)