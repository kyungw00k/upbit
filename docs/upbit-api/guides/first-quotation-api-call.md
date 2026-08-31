---
updatedAt: 2026-07-14T09:05:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Public API 호출

REST API 호출 테스트 가이드를 따라 업비트 API를 호출해보고 성공 응답을 확인할 수 있습니다.

## 시작하기

본 가이드에서는 cURL을 사용하여 업비트의 시세 REST API를 호출하고 성공 응답을 확인해 보는 초심자 과정을 안내합니다. 간단한 API 요청부터 복잡한 요청까지 따라하면서 업비트 REST API의 동작 방식을 확인해보세요.

<br />

## cURL 이란?

cURL은 HTTP 요청을 CLI(Command Line Interface) 환경에서 전송하고 응답을 확인할 수 있는 경량 커맨드라인 도구입니다. GUI 없이도 명령어만으로 요청을 실행할 수 있어, 별도의 설치 없이 대부분의 운영체제에서 기본 터미널을 통해 편리하게 사용할 수 있습니다.

<div className="callout-section">
  <div className="callout-title">
    <i className="fa-solid fa-compass"></i> 호출 예제가 아닌, 코드 기반의 연동 예제를 확인하고 싶습니다.
    </div>
    호출 테스트 없이 각 <a href="../reference/list-trading-pairs">API Reference 페이지</a>의 예제 코드 또는 <a href="../recipes">Recipes</a> 메뉴로 이동하여 주요 언어별 예제 코드와 함께 즉시 코드 구현을 시작할 수 있습니다.
    API 사용 전 반드시 알아야 할 사항들은 <a href="../reference/rest-api-guide">REST API 사용 안내</a> 및 <a href="../reference/websocket-guide">WebSocket 사용 안내</a> 페이지에서 확인해주시기 바랍니다.
</div>

<br />

## cURL을 활용한 REST API 호출 테스트

### 1. cURL 실행 환경 확인

사용 중인 운영체제(OS)에 따라 아래 안내 가이드를 클릭하여 cURL 실행 환경을 확인해주세요.

<div className="accordion">
  <input type="checkbox" id="macos-curl" />
  <label for="macos-curl">
    <i className="fa-solid fa-circle-info"></i> macOS cURL 실행 환경 확인하기
  </label>
  <div className="accordion-content">
    cURL이 기본 내장되어 있어 바로 실행 가능합니다. 터미널(Terminal)을 열고, 아래 명령어를 입력합니다.<br /><br />

    <pre><code>curl --version</code></pre>

아래와 같이 버전 정보가 출력되면 정상적으로 사용이 가능한 상태입니다. 사용자 PC 환경에 따라 버전은 상이할 수 있습니다.<br /><br />

    <pre><code>% curl --version
curl 8.7.1 (x86_64-apple-darwin24.0) libcurl/8.7.1 (SecureTransport) LibreSSL/3.3.6 zlib/1.2.12 nghttp2/1.62.0
Release-Date: 2024-03-27
Protocols: dict file ftp ftps gopher gophers http https imap imaps ipfs ipns ldap ldaps mqtt pop3 pop3s rtsp smb smbs smtp smtps telnet tftp
Features: alt-svc AsynchDNS GSS-API HSTS HTTP2 HTTPS-proxy IPv6 Kerberos Largefile libz MultiSSL NTLM SPNEGO SSL threadsafe UnixSockets</code></pre>

  </div>
</div>

<div className="accordion">
    <input type="checkbox" id="windows-curl" />
    <label for="windows-curl">
      <i className="fa-solid fa-circle-info"></i> Windows cURL 실행 환경 확인하기
    </label>
    <div className="accordion-content">
      Windows 10 이후 버전부터는 대부분 기본적으로 cURL이 내장되어 있지만, 일부 환경에서는 설치가 필요할 수 있습니다. 명령 프롬프트(cmd) 또는 PowerShell을 실행한 뒤 아래 명령어를 실행합니다. 
      명령 프롬프트(cmd)와 PowerShell은 윈도우 키 + R을 눌러 “실행” 창을 열고 "cmd" 또는 "powershell"를 입력하여 실행할 수 있습니다. 
      <br /><br />

      <pre><code>curl --version</code></pre>

버전 정보가 출력되면 설치가 완료된 상태입니다. 설치가 되어있지 않은 경우, 아래 설치 가이드를 클릭하여 펼친 후 순서를 따라 cURL을 다운로드 및 설치합니다.<br /><br />

      <ol>
        <li>cURL 공식 사이트의 Windows 버전 다운로드 페이지로 이동합니다. (https://curl.se/windows) </li>
        <li>curl for Windows 섹션에서 본인 운영체제에 맞는 버전(일반적으로 64비트 Win64 - Generic)을 선택하여 다운로드합니다.</li>
        <li>.zip 파일을 다운로드한 후 압축을 해제합니다. 예: C:\Program Files\curl</li>
        <li>압축을 해제한 폴더 위치를 환경 변수에 등록해야 합니다. 윈도우 키 > “환경 변수” 검색 → 시스템 환경 변수 편집을 클릭한 뒤, 시스템 변수 또는 사용자 변수 영역에서 Path 항목을 선택하고 새로 만들기에 아래와 같은 bin 폴더 위치를 입력하여 환경 변수에 추가합니다. 압축을 해제한 폴더 위치에 따라 실제 위치를 입력합니다.
        <pre><code>C:\Program Files\curl\bin</code></pre></li>
        <li>확인을 눌러 저장한 뒤 새로 실행한 cmd창 또는 PowerShell에서 위 버전 확인 명령어를 다시 실행하여 정상 설치를 확인할 수 있습니다.</li>
      </ol>

    </div>
  </div>

<br />

### 2. 첫 REST API 호출하기: 업비트 지원 페어 조회

cURL 실행 환경이 확인되었다면 업비트에서 지원하는 모든 페어 목록을 조회하는 API를 호출해보겠습니다. 해당 API Reference는 [페어 목록 조회](https://docs.upbit.com/kr/reference/list-trading-pairs) 에서 확인하실 수 있습니다. 다음과 같은 명령어를 터미널 또는 cmd/PowerShell 창에 입력해주세요.

```curl
curl -i --request GET \
     --url https://api.upbit.com/v1/market/all \
     --header 'Accept: application/json'
```

cURL 명령어의 각 옵션의 의미는 아래와 같습니다.

* **-i (또는 --include)**: 응답 출력 시 HTTP 응답의 상태 코드(Status Code)와 헤더를 함께 보여줍니다. 본 가이드에서 이해를 돕기 위해 설정한 옵션으로 제외할 수 있습니다.
* **--request GET**: GET 방식의 HTTP 요청을 전송합니다. API별 HTTP 메소드에 따라 POST, DELETE 등으로 지정합니다.
* **--url**: 호출하고자 하는 Endpoint 주소입니다. API 도메인과 Path를 입력합니다.
* **--header**: HTTP 요청 헤더를 설정합니다. `Accept: application/json`헤더를 설정하여 서버에게 '*응답은 JSON 형식으로 보내달라*'는 요구사항을 포함하였습니다.

위 요청을 전송한 직후 서버로부터 첫 **200 OK**응답이 반환됨을 아래와 같이 확인할 수 있습니다. (응답은 ...로 축약하였습니다.)

```
HTTP/2 200 
date: Wed, 30 Jul 2025 16:30:09 GMT
content-type: application/json;charset=UTF-8
content-length: 41776
remaining-req: group=market; min=600; sec=9
limit-by-ip: Yes
vary: origin,access-control-request-method,access-control-request-headers,accept-encoding
etag: W/"0e180e0744924ed95538e4d57fecb833f"
cache-control: no-cache, no-store, max-age=0, must-revalidate
pragma: no-cache
expires: 0

[{"market":"BTC-BERA","korean_name":"베라체인","english_name":"Berachain"},{"market":"BTC-FIL","korean_name":"파일코인","english_name":"Filecoin"},{"market":"BTC-SIGN","korean_name":"사인","english_name":"Sign"},{"market":"BTC-VIRTUAL","korean_name":"버추얼프로토콜","english_name":"Virtuals Protocol"},{"market":"BTC-BAT","korean_name":"베이직어텐션토큰","english_name":"Basic Attention Token"}...
```

응답의 주요 항목과 설명은 아래와 같습니다.

* **HTTP/2 200**: 응답의 HTTP 상태 코드가 200 OK로 반환되었음을 의미합니다. 오류 없이 정상 조회되었음을 알 수 있습니다.
* **content-type: application/json;charset=UTF-8**: 요청에 Accept 헤더로 지정한대로, 서버가 JSON 형식의 응답을 반환하였음을 나타내는 Content-Type 응답 헤더입니다.
* **remaining-req: group=market; min=600; sec=9**: 업비트 API는 IP 또는 API Key별 요청 수 제한 정책을 적용하고 있습니다. 위 값은 market API 그룹에 대해 동일 초(sec) 내의 잔여 허용 호출 수가 9번 남았음을 의미합니다.
* **limit-by-ip: Yes**: 본 요청에 IP 단위의 요청 수 제한 정책이 적용되고 있음을 의미합니다.
* **응답 JSON 데이터**: \[{"market":"BTC-BERA","korean_name":"베라체인","english_name":"Berachain"}..]와 같은 JSON Array 형식의 응답 데이터가 반환되었습니다. 업비트가 지원하는 모든 페어에 대해 `market`(페어 코드), `korean_name`(한글명), `english_name`(영문명) 정보가 JSON Object의 목록으로 반환됩니다.

<br />

### 3. 파라미터를 포함한 요청 전송하기: 페어 현재가 조회

업비트에서 지원하는 모든 페어 목록을 확인했으니, 관심이 있는 페어에 대한 현재가 정보를 조회해보겠습니다. 해당 API 명세는 [페어 단위 현재가 조회](https://docs.upbit.com/kr/reference/list-tickers)에서 확인하실 수 있습니다.

`KRW-BTC` 페어에 대한 현재가를 조회하기 위해 다음과 같은 명령어를 터미널 또는 cmd/PowerShell 창에 입력해주세요.

```curl
curl -i --request GET \
     --url 'https://api.upbit.com/v1/ticker?markets=KRW-BTC' \
     --header 'accept: application/json'
```

조회하고자 하는 페어인 `KRW-BTC`를 지정하기 위해 API Path인 `/v1/ticker` 뒤에 `?markets=KRW-BTC` 에 해당하는 쿼리 파라미터를 추가하였습니다. 동시에 여러 페어의 현재가를 조회하고 싶은 경우 콤마(,)로 구분하여 아래와 같이 페어를 추가로 지정할 수 있습니다.

```curl
curl -i --request GET \
     --url 'https://api.upbit.com/v1/ticker?markets=KRW-BTC,KRW-ETH' \
     --header 'accept: application/json'
```

아래와 같은 응답을 확인할 수 있습니다.

```
HTTP/2 200 
date: Wed, 30 Jul 2025 17:03:15 GMT
content-type: application/json;charset=UTF-8
content-length: 1619
remaining-req: group=ticker; min=600; sec=9
limit-by-ip: Yes
vary: origin,access-control-request-method,access-control-request-headers,accept-encoding
etag: W/"0f53ce9e4b1d051fc7f7b586f9789a25a"
cache-control: no-cache, no-store, max-age=0, must-revalidate
pragma: no-cache
expires: 0

[{"market":"KRW-BTC","trade_date":"20250730","trade_time":"170315","trade_date_kst":"20250731","trade_time_kst":"020315","trade_timestamp":1753894995076,"opening_price":163248000.00000000,"high_price":163708000.00000000,"low_price":162248000.00000000,"trade_price":162876000.00000000,"prev_closing_price":163158000.00000000,"change":"FALL","change_price":282000.00000000,"change_rate":0.0017283860,"signed_change_price":-282000.00000000,"signed_change_rate":-0.0017283860,"trade_volume":0.02543652,"acc_trade_price":137728155832.9616000000000000,"acc_trade_price_24h":167740032580.71786000,"acc_trade_volume":844.94369319,"acc_trade_volume_24h":1029.27369910,"highest_52_week_price":166800000.00000000,"highest_52_week_date":"2025-07-14","lowest_52_week_price":72100000.00000000,"lowest_52_week_date":"2024-08-05","timestamp":1753894995113},{"market":"KRW-ETH","trade_date":"20250730","trade_time":"170312","trade_date_kst":"20250731","trade_time_kst":"020312","trade_timestamp":1753894992075,"opening_price":5252000.00000000,"high_price":5286000.00000000,"low_price":5189000.00000000,"trade_price":5252000.00000000,"prev_closing_price":5252000.00000000,"change":"EVEN","change_price":0,"change_rate":0,"signed_change_price":0,"signed_change_rate":0,"trade_volume":0.90519935,"acc_trade_price":173631118966.8993200000000000,"acc_trade_price_24h":219851544443.94134000,"acc_trade_volume":33157.31864212,"acc_trade_volume_24h":42001.31603004,"highest_52_week_price":5900000.00000000,"highest_52_week_date":"2024-12-16","lowest_52_week_price":2096000.00000000,"lowest_52_week_date":"2025-04-09","timestamp":1753894995076}]
```

<br />

### 4. 복잡한 파라미터 처리: 캔들 조회

보다 복잡한 파라미터를 요구하는 [분(Minute) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-minutes) API를 호출해보겠습니다. 이 API는 특정 페어의 캔들(OHLCV) 데이터를 조회하는 API로, 다음과 같이 요청할 수 있습니다.

```
curl -i --request GET \
     --url 'https://api.upbit.com/v1/candles/days?market=KRW-BTC&to=2025-07-30T00%3A00%3A00%2B09%3A00&count=2' \
     --header 'accept: application/json'
```

위 요청 예시에는 `market`, `to`, `count`의 총 세 개의 파라미터가 포함되어 있으며, 각각의 파라미터는 &로 연결되어 있습니다. `to` 파라미터는 캔들 조회 기간을 지정하는 파라미터로, 위 예시에는 `2025-07-30T00:00:00+09:00` 으로 지정되어 있습니다.

**모든 쿼리 파라미터는 URL 인코딩되어야 하므로** 위의 예제와 같이 `to=2025-07-30T00%3A00%3A00%2B09%3A00`로 요청합니다. 이 요청은 `KRW-BTC` 페어에 대해 2025년 7월 30일 00시 정각을 기준으로 이전 일 캔들 2개를 조회하는 요청으로, 응답은 아래와 같습니다.

```
HTTP/2 200 
date: Wed, 30 Jul 2025 17:36:32 GMT
content-type: application/json;charset=UTF-8
content-length: 909
remaining-req: group=candles; min=600; sec=9
limit-by-ip: Yes
vary: origin,access-control-request-method,access-control-request-headers,accept-encoding
expires: 0
etag: W/"003a2ba517a62794576f650d33975cbbb"
cache-control: no-cache, no-store, max-age=0, must-revalidate
pragma: no-cache

[{"market":"KRW-BTC","candle_date_time_utc":"2025-07-29T00:00:00","candle_date_time_kst":"2025-07-29T09:00:00","opening_price":162483000.00000000,"high_price":163800000.00000000,"low_price":162034000.00000000,"trade_price":163158000.00000000,"timestamp":1753833599704,"candle_acc_trade_price":242634976840.02426000,"candle_acc_trade_volume":1487.53432374,"prev_closing_price":162483000.00000000,"change_price":675000.00000000,"change_rate":0.0041542808},{"market":"KRW-BTC","candle_date_time_utc":"2025-07-28T00:00:00","candle_date_time_kst":"2025-07-28T09:00:00","opening_price":162548000.00000000,"high_price":162985000.00000000,"low_price":161391000.00000000,"trade_price":162483000.00000000,"timestamp":1753747199919,"candle_acc_trade_price":237170267365.09124000,"candle_acc_trade_volume":1461.59937798,"prev_closing_price":162717000.00000000,"change_price":-234000.00000000,"change_rate":-0.0014380796}]
```

요청에 대한 정상적인 응답으로, 7/28일과 7/29일에 대한 캔들 데이터가 조회된 것을 확인할 수 있습니다.

[block:html]
{
  "html": "<div class=\"callout-section\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i>  URL 인코딩이란?\n      </div>\nURL 인코딩은 통신 프로토콜에서 URL 내에 포함할 수 없는 문자를 전송 가능한 문자로 변환하는 인코딩 방식입니다. 특수 문자를 포함한 대상 문자들은 인코딩 시 % 기호와 2자리 16진수로 이루어진 문자열로 변환됩니다.<br><br>\n(예시) <code>:</code>은 인코딩 후 <code>%3A</code>로, <code>+</code>는 <code>%2B</code>로 변환\n  </div>"
}
[/block]

<br />

### 5. 예제 cURL을 활용한 테스트

cURL을 사용하여 간편하게 업비트가 제공하는 일부 시세 조회 API를 사용해보았습니다. 현재가, 캔들 외에도 체결, 호가 API를 시세 API 카테고리에서 제공하고 있으니, 각 API Reference 페이지를 방문하여 쿼리 파라미터를 입력하고 우측 Shell 예제 코드 영역의 cURL 커맨드 예시를 참고하여 더 다양한 예제를 테스트해보시기 바랍니다.

<br />

## 마치며

본 가이드에서는 업비트의 시세 API를 호출하고 성공 응답을 확인하는 방법을 알아보았습니다. 아래 What's Next 링크 중 필요한 페이지로 이동하여 다음 API 호출을 진행하실 수 있습니다.

* 시세 조회 기능에 이어 자산 조회, 주문 관리, 입출금 관리 기능을 구현하고 싶다면 다음 시작하기 컨텐츠인 [API Key 발급 받기](https://docs.upbit.com/kr/docs/api-key) 가이드로 이동하여 업비트 API Key를 발급받으시기 바랍니다.
* 더 다양한 업비트 시세 조회 기능 연동이 필요한 경우 [Quotation API](https://docs.upbit.com/kr/reference/list-trading-pairs) Reference로 이동하여 다양한 API 명세와 요청 예제를 확인하실 수 있습니다.

<br />

# Sibling pages

* [Upbit API 이용 준비](https://docs.upbit.com/kr/docs/api-setup.md)
* [개발환경 설정](https://docs.upbit.com/kr/docs/dev-environment.md)
* [인증 API 호출](https://docs.upbit.com/kr/docs/first-exchange-api-call.md)