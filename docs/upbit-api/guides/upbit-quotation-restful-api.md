---
updatedAt: 2026-05-04T05:48:42.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# REST API를 이용한 시세 조회

업비트 REST API를 이용한 시세 조회 방법을 안내합니다.

[block:html]
{
  "html": "<meta name=\"robots\" content=\"noindex\">"
}
[/block]

## 시작하기

업비트 Quotation API를 이용하여 다음과 같은 시세 정보를 조회할 수 있습니다.

* **페어 목록**: 업비트에서 지원하는 모든 페어(거래 쌍) 목록
* **캔들(OHLCV)**: 시간 단위(초, 분, 일, 주, 월, 연)별 시가/고가/저가/종가/누적 거래량 등
* **체결(Trade)**: 특정 페어의 최근 체결 내역. 체결 방향, 체결 시각, 체결 가격, 체결량, 전일 종가 등
* **현재가(Ticker)**: 최근 거래 시각, 시가/고가/저가/종가/전일 종가, 가격 변화, 거래량 등
* **호가(Orderbook)**: 현재 매수/매도 호가, 매수/매도 잔량 등

이번 튜토리얼에서는 cURL과 주요 프로그래밍 언어(Python, Java, Node.js)를 사용해 시세 정보를 조회하는 방법을 안내합니다.

<br />

## 개발 환경 설정 안내

튜토리얼을 시작하기에 앞서, [개발 환경 설정 가이드](https://docs.upbit.com/kr/docs/dev-environment) 문서를 참고하여 사용하고자 하는 프로그래밍 언어에 맞는 개발 환경을 설정하시기 바랍니다. 본 튜토리얼은 가이드에서 안내한 언어별 주요 HTTP 클라이언트 라이브러리를 활용한 예제 코드를 포함합니다. cURL로 호출하는 경우 별도의 개발 환경 설정은 필요 없이 macOS 또는 Windows에서 기본으로 지원하는 터미널 또는 cmd/PowerShell을 사용할 수 있습니다. [첫 업비트 API 호출하기](https://docs.upbit.com/kr/docs/first-quotation-api-call) 가이드에 포함된 cURL 환경 설정 안내를 참조하세요.

<br />

## 1. 페어 목록 조회 후 특정 페어의 현재가(Ticker) 스냅샷 조회

업비트에서 제공하는 모든 페어 목록을 확인한 뒤, 하나의 페어를 선택하여 해당 페어의 현재가 정보를 조회해보겠습니다.

[페어 목록 조회](https://docs.upbit.com/kr/reference/list-trading-pairs) API는 업비트에서 지원하는 모든 페어 목록을 반환합니다. 응답의 market 필드 값을 다른 API에서 페어를 지정하기 위한 파라미터로 사용할 수 있습니다. 아래 예제 코드 탭에서 원하는 언어를 선택하거나, API Reference의 예제 코드 영역을 참고하여 API를 호출하세요. 각 페어의 상세 정보를 함께 조회하는 경우 is\_details 파라미터를 true로 설정합니다.

```curl
curl --request GET \
  --url https://api.upbit.com/v1/market/all
```
```python
import requests

url = "https://api.upbit.com/v1/market/all"

headers = {"accept": "application/json"}

response = requests.get(url, headers=headers)

print(response.text)
```
```java Java - OkHTTP
OkHttpClient client = new OkHttpClient();

Request request = new Request.Builder()
  .url("https://api.upbit.com/v1/market/all")
  .get()
  .addHeader("accept", "application/json")
  .build();

Response response = client.newCall(request).execute();
```
```java Java - java.net
HttpRequest request = HttpRequest.newBuilder()
    .uri(URI.create("https://api.upbit.com/v1/market/all"))
    .header("accept", "application/json")
    .method("GET", HttpRequest.BodyPublishers.noBody())
    .build();
HttpResponse<String> response = HttpClient.newHttpClient().send(request, HttpResponse.BodyHandlers.ofString());
System.out.println(response.body());
```
```node
import axios from 'axios';

const options = {
  method: 'GET',
  url: 'https://api.upbit.com/v1/market/all',
  headers: {accept: 'application/json'}
};

axios
  .request(options)
  .then(res => console.log(res.data))
  .catch(err => console.error(err));
```

응답은 아래와 같이 JSON 배열 형식으로 반환됩니다. (... 표기로 중략)

```json
[
    {
        "market": "KRW-BTC",
        "korean_name": "비트코인",
        "english_name": "Bitcoin"
    },
  	...
		{
        "market": "KRW-POLYX",
        "korean_name": "폴리매쉬",
        "english_name": "Polymesh"
    }
]

```

배열의 첫 번째 market 값인 KRW-BTC 코드를 markets 파라미터에 지정해 해당 페어의 현재가 스냅샷을 조회할 수 있습니다. 여러 페어 코드를 함께 입력하여 동시에 조회할 수 있습니다.

```curl
curl --request GET \
     --url 'https://api.upbit.com/v1/ticker?markets=KRW-BTC' \
     --header 'accept: application/json'
```
```python
import requests

url = "https://api.upbit.com/v1/ticker?markets=KRW-BTC"

headers = {"accept": "application/json"}

response = requests.get(url, headers=headers)

print(response.text)
```
```java Java - OkHTTP
OkHttpClient client = new OkHttpClient();

Request request = new Request.Builder()
  .url("https://api.upbit.com/v1/ticker?markets=KRW-BTC")
  .get()
  .addHeader("accept", "application/json")
  .build();

Response response = client.newCall(request).execute();
```
```java Java - java.net
HttpRequest request = HttpRequest.newBuilder()
    .uri(URI.create("https://api.upbit.com/v1/ticker?markets=KRW-BTC"))
    .header("accept", "application/json")
    .method("GET", HttpRequest.BodyPublishers.noBody())
    .build();
HttpResponse<String> response = HttpClient.newHttpClient().send(request, HttpResponse.BodyHandlers.ofString());
System.out.println(response.body());
```
```node
import axios from 'axios';

const options = {
  method: 'GET',
  url: 'https://api.upbit.com/v1/ticker?markets=KRW-BTC',
  headers: {accept: 'application/json'}
};

axios
  .request(options)
  .then(res => console.log(res.data))
  .catch(err => console.error(err));
```

아래와 같이 JSON 배열 형식으로 반환된 현재가 응답을 확인할 수 있습니다.

```json
[
  {
    "market": "KRW-BTC",
    "trade_date": "20250704",
    "trade_time": "051400",
    "trade_date_kst": "20250704",
    "trade_time_kst": "141400",
    "trade_timestamp": 1751606040365,
    "opening_price": "148737000.00000000,",
    "high_price": "149360000.00000000,",
    "low_price": "148288000.00000000,",
    "trade_price": "148601000.00000000,",
    "prev_closing_price": "148737000.00000000,",
    "change": "FALL",
    "change_price": 136000,
    "change_rate": 0.0009143656,
    "signed_change_price": "-136000.00000000,",
    "signed_change_rate": -0.0009143656,
    "trade_volume": 0.00016823,
    "acc_trade_price": 31615925234.05438,
    "acc_trade_price_24h": 178448329314.96686,
    "acc_trade_volume": 212.38911576,
    "acc_trade_volume_24h": 1198.26954807,
    "highest_52_week_price": 163325000,
    "highest_52_week_date": "2025-01-20",
    "lowest_52_week_price": 72100000,
    "lowest_52_week_date": "2024-08-05",
    "timestamp": 1751606040403
  }
]
```

<br />

<br />

## 2. 캔들 조회

캔들이란 일반적으로 봉이라고 불리기도 하며 차트를 구성하는 기본적인 막대입니다. 한 번에 최대 200개의 캔들을 요청할 수 있으며 이를 초과하는 경우 페이지네이션을 이용해 순차적으로 요청하는 것을 권장합니다.

<br />

### 최근 캔들 조회하기

KRW-BTC 페어에서 가장 최근 3개의 5분봉을 조회하기 위해서는 다음과 같이 요청합니다. 언어별 예제 코드는 [API Reference](https://docs.upbit.com/kr/reference/list-candles-minutes) 에를 참고하시기 바랍니다.

```curl
curl --request GET \
  --url 'https://api.upbit.com/v1/candles/minutes/5?market=KRW-BTC&count=3'
```

응답 예제는 아래와 같습니다.

```json
[
    {
        "market": "KRW-BTC",
        "candle_date_time_utc": "2025-08-01T14:00:00",
        "candle_date_time_kst": "2025-08-01T23:00:00",
        "opening_price": 159399000.00000000,
        "high_price": 159525000.00000000,
        "low_price": 159001000.00000000,
        "trade_price": 159090000.00000000,
        "timestamp": 1754057055712,
        "candle_acc_trade_price": 3784659174.62006000,
        "candle_acc_trade_volume": 23.77492422,
        "unit": 5
    },
    {
        "market": "KRW-BTC",
        "candle_date_time_utc": "2025-08-01T13:55:00",
        "candle_date_time_kst": "2025-08-01T22:55:00",
        "opening_price": 159382000.00000000,
        "high_price": 159678000.00000000,
        "low_price": 159255000.00000000,
        "trade_price": 159399000.00000000,
        "timestamp": 1754056799898,
        "candle_acc_trade_price": 4231889234.92549000,
        "candle_acc_trade_volume": 26.53149963,
        "unit": 5
    },
    {
        "market": "KRW-BTC",
        "candle_date_time_utc": "2025-08-01T13:50:00",
        "candle_date_time_kst": "2025-08-01T22:50:00",
        "opening_price": 159500000.00000000,
        "high_price": 159615000.00000000,
        "low_price": 159205000.00000000,
        "trade_price": 159355000.00000000,
        "timestamp": 1754056499639,
        "candle_acc_trade_price": 4466832748.99644000,
        "candle_acc_trade_volume": 28.02236287,
        "unit": 5
    }
]
```

<br />

### to 파라미터로 특정 시간대의 캔들 조회하기

캔들 API에 `to` 파라미터로 특정 시각을 지정하면 해당 시각 이전의 캔들을 조회할 수 있습니다. 지정하고자 하는 시각 ISO8601 표준 형식으로 입력합니다.

1. UTC 기준 0시까지 1분봉을 1개 요청\
   <https://api.upbit.com/v1/candles/minutes/1?market=KRW-BTC&to=2025-07-24T00:00:00Z>

2. KST 기준(UTC + 09:00) 7시까지 1분봉을 1개 요청하기\
   <https://api.upbit.com/v1/candles/minutes/1?market=KRW-BTC&to=2025-07-27T07:00:00+09:00>

`to` 파라미터를 입력한 경우 전송 전 반드시 URL 인코딩을 적용해 아래와 같이 요청해야 합니다.

<https://api.upbit.com/v1/candles/minutes/1?market=KRW-BTC&to=2025-09-27T07:00:00%2B09:00>

### 캔들 생성 기준 안내

* 캔들은 해당 시간대에 체결이 발생한 경우에만 생성됩니다.
* 해당 시간대에 체결이 발생하지 않은 경우, 해당 시간대의 캔들은 생성되지 않으며 응답 데이터에도 포함되지 않습니다.

> **주의사항**
>
> 캔들 데이터 활용시 비어있는 시간 구간을 고려한 구현이 필요합니다.

<br />

## 3. 시세 체결 조회

최근 체결 이력을 조회하여 해당 자산의 최근 거래 가격 및 추이를 확인할 수 있습니다. [페어 체결 이력 조회](https://docs.upbit.com/kr/reference/list-pair-trades) API로 최근 7일의 체결 이력을 최대 200건 조회할 수 있습니다. 페이지네이션을 활용하여 순차 조회합니다.

```curl
curl --request GET \
     --url 'https://api.upbit.com/v1/trades/ticks?market=KRW-BTC' \
     --header 'accept: application/json'
```

위와 같이 페어만 지정하여 요청하는 경우 해당 페어의 최근 체결 데이터가 반환됩니다. 특정 일자의 체결 이력을 조회하고자 하는 경우 아래와 같이 `days_ago` 파라미터를 사용합니다. 2로 설정하는 경우 UTC 기준 2일 전 체결 데이터가 내림차순으로 반환됩니다.

```curl
curl --request GET \
     --url 'https://api.upbit.com/v1/trades/ticks?market=KRW-BTC&days_ago=2' \
     --header 'accept: application/json'
```

2일 전 특정 시간대의 체결 이력을 지정하여 조회하고 싶은 경우 `to` 파라미터를 사용할 수 있습니다. `to`에 입력한 시간을 기준으로 해당 시간 이전의 체결 이력이 시간 역순으로 반환됩니다. 아래 예제 요청시 UTC 22:00:00 이전 체결 이력을 반환합니다.

```curl
curl --request GET \
     --url 'https://api.upbit.com/v1/trades/ticks?market=KRW-BTC&to=220000&days_ago=2' \
     --header 'accept: application/json'
```

반환되는 응답 예시는 아래와 같습니다.

```json
[
    {
        "market": "KRW-BTC",
        "trade_date_utc": "2025-07-30",
        "trade_time_utc": "21:59:59",
        "timestamp": 1753912799945,
        "trade_price": 162806000.00000000,
        "trade_volume": 0.00012284,
        "prev_closing_price": 163158000.00000000,
        "change_price": -352000.00000000,
        "ask_bid": "BID",
        "sequential_id": 17539127999450000
    }
]
```

<br />

## 4. 호가 조회

거래소의 실시간 호가 정보를 API로 조회할 수 있습니다. [호가 조회](https://docs.upbit.com/kr/reference/list-orderbooks) API를 아래와 같이 요청합니다.

```curl
curl --request GET \
     --url 'https://api.upbit.com/v1/orderbook?markets=KRW-BTC&count=2' \
     --header 'accept: application/json'
```

반환된 호가 정보는 아래와 같습니다. count 파라미터를 2로 요청하여 2개의 오더북 쌍을 조회했습니다.

```json
[
    {
        "market": "KRW-BTC",
        "timestamp": 1754057310152,
        "total_ask_size": 1.44700033,
        "total_bid_size": 1.83572538,
        "orderbook_units": [
            {
                "ask_price": 159399000,
                "bid_price": 159396000,
                "ask_size": 0.50775343,
                "bid_size": 0.30813376
            },
            {
                "ask_price": 159400000,
                "bid_price": 159385000,
                "ask_size": 0.03768668,
                "bid_size": 0.0003137
            }
        ],
        "level": 0
    }
]
```

조회 시점을 기준으로 최고 매수 호가와 최저 매도 호가의 쌍이 `orderbook_units` 필드의 첫번째 객체로 반환됩니다. 위 데이터의 경우 `159399000` 원의 최저 매도 호가에 0.50775343 개의 BTC 수량이, `159396000` 원의 최고 매수 호가에 0.30813376 개의 BTC 수량이 체결 대기 중임을 의미합니다. 두 번째 객체는 각각 두번째 순위의 매수 호가와 매도 호가의 쌍입니다.

# Sibling pages

* [WebSocket으로 실시간 캔들 차트 생성](https://docs.upbit.com/kr/docs/how-to-subscribe-quotation-stream.md)
* [캔들 API로 지정한 기간의 시세 데이터를 CSV 파일로 출력하기](https://docs.upbit.com/kr/docs/how-to-download-candle-data.md)