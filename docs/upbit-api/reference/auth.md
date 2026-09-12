---
updatedAt: 2026-09-01T10:28:35.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 인증

REST API 및 WebSocket 사용을 위한 인증 가이드입니다.

거래 및 자산 관리(Exchange) REST API를 호출하거나 WebSocket 데이터를 수신하려면 API Key 기반의 인증이 반드시 필요합니다. 본 가이드를 통해 인증 방식을 확인하고 언어별 구현 예시를 참고 할 수 있습니다.

## API Key

### API Key 발급 및 구조 안내

인증 구현에 앞서 [API Key 발급 가이드](https://docs.upbit.com/kr/docs/api-key)를 참고해 API Key를 발급받고, 호출지의 IP를 허용 목록에 등록해야 합니다. API Key 당 최대 10개의 IP를 등록할 수 있습니다.

API Key는 Access Key와 Secret Key 쌍으로 구성되며, Secret Key는 발급 시에만 확인 가능합니다. Secret Key는 **외부에 노출되지 않도록 반드시 안전하게 보관**해야 하며, 토큰 생성 시 반드시 Access Key와 짝을 이루는 Secret Key를 사용해야 합니다.

### API Key의 권한 그룹

업비트 API Key는 발급 시 사용 용도에 따라 필요한 권한만 선택적으로 부여할 수 있습니다. API 호출 시 권한 오류가 발생했다면, 해당 API Key에 필요한 권한이 포함되어 있는지 확인하세요. 각 권한 그룹이 지원하는 기능 목록은 아래 박스를 클릭해 확인하거나, API Reference 문서 하단의 “API Key Permission” 항목을 참고하시기 바랍니다.

<Accordion title="Exchange REST API 및 WebSocket 타입 별 필요 API Key 권한 그룹" icon="fa-info-circle">
  ###

  | 권한 그룹    | 허용 REST API                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        | 허용 WebSocket 타입                                                                            |
  | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
  | -        | [서비스 - 통화별 입출금 서비스 상태 조회](https://docs.upbit.com/kr/reference/get-service-status)<br />[서비스 - API Key 목록 조회](https://docs.upbit.com/kr/reference/list-api-keys)                                                                                                                                                                                                                                                                                                                                                                                    | [공지사항(Announcement) 타입 데이터 구독](https://docs.upbit.com/kr/reference/websocket-announcement) |
  | **자산조회** | [자산 - 포켓 잔고 조회](https://docs.upbit.com/kr/reference/get-balance)                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   | [내 자산(MyAsset) 타입 데이터 구독](https://docs.upbit.com/kr/reference/websocket-myasset)           |
  | **주문하기** | [주문 생성](https://docs.upbit.com/kr/reference/new-order)<br />[주문 생성 테스트](https://docs.upbit.com/kr/reference/order-test)<br />[단일 주문 취소](https://docs.upbit.com/kr/reference/cancel-order)<br />[지정 주문 목록 취소](https://docs.upbit.com/kr/reference/cancel-orders-by-ids)<br />[주문 일괄 취소](https://docs.upbit.com/kr/reference/batch-cancel-orders)<br />[취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order)                                                                                                                          | -                                                                                          |
  | **주문조회** | [주문 가능정보 조회](https://docs.upbit.com/kr/reference/available-order-information)<br />[단일 주문 조회](https://docs.upbit.com/kr/reference/get-order)<br />[주문 목록 조회](https://docs.upbit.com/kr/reference/list-orders-by-ids)<br />[체결 대기 주문 조회](https://docs.upbit.com/kr/reference/list-open-orders)<br />[종료 주문 조회](https://docs.upbit.com/kr/reference/list-closed-orders)                                                                                                                                                                              | [내 주문 및 체결(MyOrder) 타입 데이터 구독](https://docs.upbit.com/kr/reference/websocket-myorder)      |
  | **출금하기** | [디지털 자산 출금하기](https://docs.upbit.com/kr/reference/withdraw)<br />[원화 출금하기](https://docs.upbit.com/kr/reference/withdraw-krw)<br />[디지털 자산 출금 취소 접수](https://docs.upbit.com/kr/reference/cancel-withdrawal)                                                                                                                                                                                                                                                                                                                                         | -                                                                                          |
  | **출금조회** | [출금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-withdrawal-information)<br />[출금 허용 주소 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawal-addresses)<br />[단일 출금 조회](https://docs.upbit.com/kr/reference/get-withdrawal)<br />[출금 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawals)                                                                                                                                                                                                                                 | -                                                                                          |
  | **입금하기** | [원화 입금](https://docs.upbit.com/kr/reference/deposit-krw)<br />[입금 UUID로 트래블룰 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid)<br />[입금 TXID로 트래블룰 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-txid)                                                                                                                                                                                                                                                                                                       | -                                                                                          |
  | **입금조회** | [입금 주소 생성 요청](https://docs.upbit.com/kr/reference/create-deposit-address)<br />[입금 가능 통화 조회](https://docs.upbit.com/kr/reference/available-deposit-information)<br />[단일 입금 주소 조회](https://docs.upbit.com/kr/reference/get-deposit-address)<br />[입금 주소 목록 조회](https://docs.upbit.com/kr/reference/list-deposit-addresses)<br />[단일 입금 조회](https://docs.upbit.com/kr/reference/get-deposit)<br />[입금 목록 조회](https://docs.upbit.com/kr/reference/list-deposits)<br />[트래블룰 지원 거래소 목록 조회](https://docs.upbit.com/kr/reference/list-travelrule-vasps) | -                                                                                          |
</Accordion>

## 인증 토큰

인증 토큰이란 서버에 요청을 보낼 때 사용자의 신원 및 권한을 증명하기 위해 전달하는 문자열 정보입니다. 업비트 API는 [JWT 토큰](https://jwt.io) 기반 인증을 사용합니다. 인증이 필요한 API를 요청하는 경우 요청을 전송하기 전에 유효한 JWT 토큰을 생성하고, 요청 헤더에 생성된 토큰을 포함하여 전송해야 합니다.

### JWT 토큰 구조

JWT는 헤더(Header), 페이로드(Payload), 서명(Signature)의 세 부분으로 구성됩니다. 각 부분은 온점(.)으로 구분되며, 헤더에는 JWT 서명 생성시 사용한 알고리즘 정보가, 페이로드에는 검증하고자 하는 내용 본문이, 서명에는 페이로드에 대한 서명 값이 포함됩니다.

<style>
  .jwt-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  @media (min-width: 768px) {
    .jwt-container {
      flex-direction: row;
    }
  }

  .jwt-box {
    border: 1px dashed #ccc;
    border-radius: 8px;
    padding: 16px;
    background-color: #fdfdfd;
    font-family: monospace;
    font-size: 11px;
    word-break: break-word;
    white-space: normal;
  }

  .jwt-string {
    flex: 1;
  }

  .jwt-parts {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .jwt-box span {
    display: inline;
  }

  .jwt-header {
    color: #1976d2;
  }

  .jwt-payload {
    color: #388e3c;
  }

  .jwt-signature {
    color: #f57c00;
  }
</style>

<div className="accordion">
  <input type="checkbox" id="jwt-example" />
  <label for="jwt-example">
    <i className="fa-solid fa-circle-info"></i> JWT 토큰 예시와 구조
  </label>
	<div className="accordion-content">
    <div className="jwt-container">
      {/* 전체 JWT 문자열 */}
      <div className="jwt-box jwt-string">
        <span className="jwt-header">eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9</span>.
        <span className="jwt-payload">eyJhY2Nlc3Nfa2V5IjoiYTdYZDkyTG1RVzN2QnRSellwTWo1Q3hOS2VUMUh1VnMwZkZnSmNBdyIsIm5vbmNlIjoiYjJmMWUzZjgtMmRjMS00ZDZmLWE4MzgtYzc0YzQ5YjBlMzlhIiwicXVlcnlfaGFzaCI6IjBiM2U4ODRkNDBjYzk5MmE4NTczMGNmNDcwYjRhMzI4NmYxM2Q5YzQ2MjA0Mjc5ZWYzMjE1M2JjZGNkNGVkYjdjMTI5MjVlNzI2NjYzNmU4NmE2ZDZhZTU4MDRhNmJiMzk0ZTYzMmU0ZGJhOWI0MDQ1YWQ0NzBjOTM3MjBlNWU2IiwicXVlcnlfaGFzaF9hbGciOiJTSEE1MTIifQ</span>.
        <span className="jwt-signature">jFOftW_qV7loHsaeZj_hrpkP2iEgdK62D9rwseXSGbwvlNaLG2-iJ-WBFUEG7c1tch0AA-jsHsonhEhezd8sTw</span>
      </div>

      <div className="jwt-parts">
        <div className="jwt-box jwt-header">
          {"alg": "HS512","typ": "JWT"}
        </div>
        <div className="jwt-box jwt-payload">
          {"access_key": "a7Xd92LmQW3vBtRzYpMj5CxNKeT1HuVs0fFgJcAw","nonce": "b2f1e3f8-2dc1-4d6f-a838-c74c49b0e39a","query_hash":"0b3e884d40cc992a85730cf470b4a3286f13d9c46204279ef32153bcdcd4edb7c12925e7266636e86a6d6ae5804a6bb394e632e4dba9b4045ad470c93720e5e6","query_hash_alg": "SHA512"}
        </div>
        <div className="jwt-box jwt-signature">

HMAC-SHA512(base64UrlEncode(<span className="jwt-header">header</span>) + "." + base64UrlEncode(<span className="jwt-payload">payload</span>), Secret Key)

        </div>
      </div>
    </div>
 </div>
</div>

#### 헤더(Header)

헤더는 다음과 같은 구조의 JSON 데이터를 Base64 인코딩한 결과 문자열을 사용합니다. 이때, `alg` 필드에는 실제 토큰 서명 생성에 사용한 알고리즘을 명시해야 하며, `HS512`(HMAC with SHA-512) 사용을 권장합니다.

```json
{
  "alg": "HS512",
  "typ": "JWT"
}
```

#### 페이로드(Payload)

페이로드는 다음과 같은 구조의 JSON 데이터를 Base64 인코딩한 결과 문자열을 사용합니다.

<table className="custom-table">
  <thead>
    <tr>
      <th>Key</th>
      <th>설명</th>
      <th>필수 여부(필수/선택)</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td className="code-col"><code>access_key</code></td>
      <td>API Key의 Access key. API Key를 지정하기 위해 사용합니다.</td>
      <td>필수</td>
    </tr>
    <tr>
      <td className="code-col"><code>nonce</code></td>
      <td>무작위의 UUID 문자열. 동일한 요청을 반복하더라도 항상 토큰 값이 변경되도록, 매 요청마다 새로운 값을 입력해야 합니다.
      </td>
      <td>필수</td>
    </tr>
    <tr>
      <td className="code-col"><code>query_hash</code></td>
      <td>요청 쿼리 파라미터 또는 본문(body)의 쿼리 스트링(Query String) Hash값.</td>
      <td>REST API 요청에 쿼리 파라미터 또는 본문이 있는 경우 필수</td>
    </tr>
    <tr>
      <td className="code-col"><code>query_hash_alg</code></td>
      <td>query_hash_alg 생성시 사용한 Hash 알고리즘. <code>query_hash</code> 값이 없는 경우 입력하지 않습니다. 기본값은 <code>SHA512</code>로, <code>query_hash</code> 생성시 다른 알고리즘을 사용한 경우 필수로 입력해야 합니다.</td>
      <td>선택</td>
    </tr>
  </tbody>
</table>

```json
// 파라미터 또는 본문(Body)이 있는 REST API 인증 토큰 생성을 위한 Payload 예시
{
  "access_key": "a7Xd92LmQW3vBtRzYpMj5CxNKeT1HuVs0fFgJcAw",
  "nonce": "b2f1e3f8-2dc1-4d6f-a838-c74c49b0e39a",
  "query_hash": "0b3e884d40cc992a85730cf470b4a3286f13d9c46204279ef32153bcdcd4edb7c12925e7266636e86a6d6ae5804a6bb394e632e4dba9b4045ad470c93720e5e6",
  "query_hash_alg": "SHA512"
}

// WebSocket 및 파라미터 또는 본문이 없는 REST API인증 토큰 생성을 위한 Payload 예시
{
  "access_key": "a7Xd92LmQW3vBtRzYpMj5CxNKeT1HuVs0fFgJcAw",
  "nonce": "b2f1e3f8-2dc1-4d6f-a838-c74c49b0e39a"
}
```

#### 서명(Signature)

API Key의 Secret Key를 사용하여 헤더와 페이로드를 서명합니다. 코드로 구현하는 경우 언어별로 지원되는 JWT 라이브러리를 활용하여 생성하는 것을 권장합니다.

<Callout icon="fad fa-circle-exclamation" theme="warn">
  ### **Secret Key는 Base64 인코딩 되어있지 않습니다.**

  JWT 토큰의 서명 데이터 생성 시 키 사용을 위해 별도의 Base64 디코딩을 수행할 필요가 없습니다. JWT 생성을 위한 라이브러리 사용 시 해당 옵션을 참고하여 구현하시기 바랍니다.
</Callout>

<br />

## JWT 생성 가이드 - 페이로드의 `query_hash` 값 생성 안내

페이로드의 `query_hash` 값은 쿼리 문자열 또는 쿼리 문자열 형태로 가공한 JSON 본문(Body)을 Hash한 값입니다. 실제 요청 내용과 인증 토큰에 포함된 내용 및 문자열 구성 순서가 다른 경우 토큰 검증 오류가 발생할 수 있으므로 아래 주의사항을 반드시 확인한 뒤 사용하시기 바랍니다.

### `GET` 또는 `DELETE` REST API 요청 시

* 실제 요청에 포함된 쿼리 문자열을 그대로 발췌하여 Hash합니다. 파라미터 간 순서를 변경하거나 재정렬하지 않습니다.
* 배열 형식의 쿼리 파라미터인 `uuids[]` 또는 `states[]`처럼 이름에 `[]`가 포함된 파라미터에 여러 값을 사용하는 경우, `uuids[]=UUID1&uuids[]=UUID2...`와 같이 동일한 Key를 반복하여 값을 연결합니다.
* 쉼표로 구분된 문자열 값을 사용하는 쿼리 파라미터인 `pairs`, `quote_currencies`처럼 이름에 `[]`가 포함되지 않은 경우, 쿼리 문자열을 `pairs=KRW-BTC,KRW-ETH...`형식으로 작성합니다.
* URL 인코딩되지 않은 쿼리 문자열을 기준으로 Hash 값을 생성해야 합니다.

### `POST` REST API 요청 시

* JSON 형식의 요청 본문(Body)의 모든 Key-Value 쌍을 쿼리 문자열 형식(`=`로 Key-Value 연결, `&`로 구분)으로 가공해야 합니다.

<br />

### 인증 토큰 전송 방식

생성한 JWT 인증 토큰은 요청의 Authorization 헤더를 이용하여 전송합니다. REST API 호출과 WebSocket 연결 요청에 동일하게 적용됩니다. Bearer 방식으로 처리하기 위해 다음과 같이 요청합니다.

* **Key**: `Authorization`
* **Value**: Bearer {공백 뒤 생성한 JWT 토큰 값을 입력}

<br />

### JWT 토큰 생성 가이드 - 코드 예시

주요 언어별 JWT 토큰 생성 코드 예시는 아래와 같습니다.

<br />

#### REST API 호출 예시

```python
from urllib.parse import quote, unquote, urlencode
from typing import Any, Dict
import hashlib
import uuid
import jwt # PyJWT
import requests

def _build_query_string(params: Dict[str, Any]) -> str:
    return unquote(urlencode(params, doseq=True))

def _create_jwt(access_key: str, secret_key: str, query_string: str = "") -> str:
    payload = {"access_key": access_key, "nonce": str(uuid.uuid4())}

    if query_string:
        query_hash = hashlib.sha512(query_string.encode("utf-8")).hexdigest()
        payload["query_hash"] = query_hash
        payload["query_hash_alg"] = "SHA512"

    token = jwt.encode(payload, secret_key, algorithm="HS512")
    return token if isinstance(token, str) else token.decode('utf-8')  

if __name__ == "__main__":
    base_url = "https://api.upbit.com"
    access_key = "YOUR_ACCESS_KEY"
    secret_key = "YOUR_SECRET_KEY" # 실제로는 안전한 방식으로 로드하거나 주입하세요.
    
    # 파라미터가 없는 요청 예시
    jwt_token = _create_jwt(access_key, secret_key)
    headers = {"Authorization": f"Bearer {jwt_token}"}
                
    response = requests.get(f"{base_url}/v1/accounts", headers=headers)
        
    print(response.json())
    
    # 파라미터가 있는 GET 요청 예시
    params = {
        "market": "KRW-BTC",
        "states[]": ["wait", "watch"],
        "limit": 10
    }
    query_string = _build_query_string(params)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    headers = {"Authorization": f"Bearer {jwt_token}"}
        
    response = requests.get(f"{base_url}/v1/orders/open?{query_string}", headers=headers)
    print(response.json())    
    
    # POST 요청 예시
    order_data = {
        "market": "KRW-BTC",
        "side": "bid",
        "volume": "0.001",
        "price": "50000000",
        "ord_type": "limit"
    }
        
    query_string = _build_query_string(order_data)
    jwt_token = _create_jwt(access_key, secret_key, query_string)
    headers = {
        "Authorization": f"Bearer {jwt_token}",
        "Content-Type": "application/json"
    }
    
    # 아래 주석처리된 부분 실행시 실제 주문이 발생하므로 실행 전 반드시 확인하세요.
    # response = requests.post(f"{base_url}/v1/orders", json=order_data, headers=headers)
    # print(response.json())     
```
```java Java - OkHttp
package com.upbit.openapi.test;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTCreator;
import com.auth0.jwt.algorithms.Algorithm;
import com.google.gson.Gson;
import java.io.IOException;
import java.math.BigInteger;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.UUID;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class Auth {

    private static Gson gson = new Gson();

    /**
     * 문자열을 입력 받아 Hash를 생성하는 예제
     */
    public static String sha512(String input) throws NoSuchAlgorithmException {
        // input은 non-null의 유효한 값으로 가정
        MessageDigest md = MessageDigest.getInstance("SHA-512");
        md.update(input.getBytes(StandardCharsets.UTF_8));

        return HexFormat.of().formatHex(md.digest());
    }

    /**
     * Access Key, Secret Key, QueryString을 입력받아 JWT Token을 생성하는 예제
     */
    private static String createJwt(String accessKey, String secretKey, String queryString)
        throws NoSuchAlgorithmException {
        // accessKey, secretKey, payload는 모두 non-null의 유효한 값으로 가정
        byte[] secretKeyBytes = secretKey.getBytes(StandardCharsets.UTF_8);
        Algorithm algorithm;

        try {
            algorithm = Algorithm.HMAC512(secretKeyBytes);
        } finally {
            Arrays.fill(secretKeyBytes, (byte) 0);
        }

        // Build JWT with claims
        JWTCreator.Builder builder = JWT.create()
            .withHeader(Collections.singletonMap("alg", "HS512"))
            .withClaim("access_key", accessKey)
            .withClaim("nonce", UUID.randomUUID().toString());

        // queryString이 있는 경우 페이로드에 Hash 추가
        if (queryString != null && !queryString.isEmpty()) {
            String queryHash = sha512(queryString);
            builder.withClaim("query_hash", queryHash);
            builder.withClaim("query_hash_alg", "SHA512");
        }
        return builder.sign(algorithm);
    }

    /**
     * 요청의 Json Body를 QueryString으로 변환하는 예제
     */
    public static String jsonToQueryString(String jsonString) {
        if (jsonString == null || jsonString.isEmpty()) {
            return "";
        }

        Map<String, Object> bodyMap = gson.fromJson(jsonString, Map.class);
        if (bodyMap != null && !bodyMap.isEmpty()) {
            List<String> queryElements = new ArrayList<>();
            for (Map.Entry<String, Object> entry : bodyMap.entrySet()) {
                if (entry.getValue() != null) {
                    try {
                        String encodedKey = URLEncoder.encode(entry.getKey(), "UTF-8");
                        String encodedValue = URLEncoder.encode(String.valueOf(entry.getValue()),
                            "UTF-8");

                        encodedKey = encodedKey.replace("%5B", "[").replace("%5D", "]");

                        queryElements.add(encodedKey + "=" + encodedValue);
                    } catch (Exception e) {
                        throw new RuntimeException("Encoding failed", e);
                    }
                }
            }
            return String.join("&", queryElements);
        }
        return "";
    }

    /**
     * Map<String, Object>를 URL 인코딩된 Query String 변환하는 예제
     */
    public static String buildQueryString(Map<String, Object> params) {
        if (params == null || params.isEmpty()) {
            return "";
        }

        List<String> components = new ArrayList<>();

        for (Map.Entry<String, Object> entry : params.entrySet()) {
            String key = entry.getKey();
            Object value = entry.getValue();

            if (value == null) {
                continue;
            }

            List<Object> values;
            if (value instanceof List) {
                values = (List<Object>) value;
            } else {
                values = Collections.singletonList(value);
            }

            for (Object val : values) {
                try {
                    String encodedKey = URLEncoder.encode(
                        key.endsWith("[]") ? key : key + "[]", StandardCharsets.UTF_8
                    ).replace("%5B", "[").replace("%5D", "]");

                    String encodedVal = URLEncoder.encode(String.valueOf(val),
                        StandardCharsets.UTF_8);
                    components.add(encodedKey + "=" + encodedVal);
                } catch (Exception e) {
                    throw new RuntimeException("Encoding failed", e);
                }
            }
        }

        return String.join("&", components);
    }


    public static void main(String[] args) throws IOException, NoSuchAlgorithmException {
        String baseUrl = "https://api.upbit.com";
        String accessKey = "YOUR_ACCESS_KEY";
        String secretKey = "YOUR_SECRET_KEY"; // 실제로는 안전하게 로드하거나 주입하세요.

        OkHttpClient client = new OkHttpClient();

        // 파라미터가 있는 GET 요청 예시
        Map<String, Object> queryParams = new HashMap<>();
        queryParams.put("states[]", Arrays.asList("wait", "watch"));
        queryParams.put("limit", 100);

        String queryString = buildQueryString(queryParams);
        String jwtTokenGet = createJwt(accessKey, secretKey, queryString);

        Request getRequest = new Request.Builder()
            .url(baseUrl + "/v1/orders/open?" + queryString)
            .get()
            .addHeader("Accept", "application/json")
            .addHeader("Authorization", "Bearer " + jwtTokenGet)
            .build();

        Response response = client.newCall(getRequest).execute();
        System.out.println(response.code());
        System.out.println(Objects.requireNonNull(response.body()).string());

        // POST 요청 예시
        final MediaType JSON = MediaType.parse("application/json; charset=utf-8");
        String jsonBody = "{\"market\":\"KRW-BTC\",\"side\":\"bid\",\"volume\":\"0.0001\",\"price\":\"50000000\",\"ord_type\":\"limit\"}";
        String queryStringBody = jsonToQueryString(jsonBody);
        String jwtTokenPost = createJwt(accessKey, secretKey, queryStringBody);

        Request postRequest = new Request.Builder()
            .url(baseUrl + "/v1/orders")
            .post(RequestBody.create(jsonBody, JSON))
            .addHeader("Accept", "application/json")
            .addHeader("Authorization", "Bearer " + jwtTokenPost)
            .build();

        // 아래 주석처리된 부분 실행시 실제 주문이 발생하므로 실행 전 반드시 확인하세요.
        response = client.newCall(postRequest).execute();
        System.out.println(response.code());
        System.out.println(Objects.requireNonNull(response.body()).string());
    }

}

```
```node
const axios = require('axios');
const crypto = require('crypto');
const jwt = require('jsonwebtoken');
const { v4: uuidv4 } = require('uuid');
const qs = require('querystring');

// 필수 정보 입력
const ACCESS_KEY = 'YOUR_ACCESS_KEY';
const SECRET_KEY = 'YOUR_SECRET_KEY'; // 안전하게 관리하세요
const BASE_URL = 'https://api.upbit.com';

/**
 * SHA512 해시 함수
 */
function sha512(text) {
  return crypto.createHash('sha512').update(text, 'utf8').digest('hex');
}

/**
 * JWT 토큰 생성
 */
function createJwt(accessKey, secretKey, queryString = '') {
  const payload = {
    access_key: accessKey,
    nonce: uuidv4(),
  };

  if (queryString) {
    payload.query_hash = sha512(queryString);
    payload.query_hash_alg = 'SHA512';
  }

  return jwt.sign(payload, secretKey, { algorithm: 'HS512' });
}

/**
 * GET 요청 예제
 */
async function getOpenOrders() {
  const query = {
    states: ['wait', 'watch'],
    page: 1,
    limit: 10,
  };

  const queryString = qs.stringify(query, '&', '=', {
    encodeURIComponent: (str) =>
      str.endsWith('[]') ? encodeURIComponent(str).replace('%5B%5D', '[]') : encodeURIComponent(str),
  });

  const token = createJwt(ACCESS_KEY, SECRET_KEY, queryString);
  const headers = {
    Authorization: `Bearer ${token}`,
    Accept: 'application/json',
  };

  try {
    const res = await axios.get(`${BASE_URL}/v1/orders/open?${queryString}`, { headers });
    console.log('[GET] Status:', res.status);
    console.log(res.data);
  } catch (err) {
    console.error('[GET] Error:', err.response?.data || err.message);
  }
}

/**
 * POST 요청 예제
 */
async function placeOrder() {
  const body = {
    market: 'KRW-BTC',
    side: 'bid',
    volume: '0.001',
    price: '50000000',
    ord_type: 'limit',
  };

  const queryString = qs.stringify(body);
  const token = createJwt(ACCESS_KEY, SECRET_KEY, queryString);
  const headers = {
    Authorization: `Bearer ${token}`,
    Accept: 'application/json',
    'Content-Type': 'application/json',
  };

  try {
    // 아래 주석처리된 부분 실행시 실제 주문이 발생하므로 실행 전 반드시 확인하세요.
    // const res = await axios.post(`${BASE_URL}/v1/orders`, body, { headers });
    // console.log('[POST] Status:', res.status);
    // console.log(res.data);
    console.log('[POST] Request prepared but not sent (order disabled for safety).');
  } catch (err) {
    console.error('[POST] Error:', err.response?.data || err.message);
  }
}

// 메인 실행
(async () => {
  await getOpenOrders();
  await placeOrder();
})();
```

#### WebSocket 연결 요청 예시

```python
import jwt  # PyJWT
import uuid
import websocket  # websocket-client

def on_message(ws, message):
    # do something
    data = message.decode('utf-8')
    print(data)


def on_connect(ws):
    print("connected!")
    # Request after connection
    ws.send('[{"ticket":"test example"},{"type":"myOrder"}]')


def on_error(ws, err):
    print(err)


def on_close(ws, status_code, msg):
  print("closed!")

payload = {
    'access_key': "YOUR_ACCESS_KEY",
    'nonce': str(uuid.uuid4()),
}

jwt_token = jwt.encode(payload, "YOUR_SECRET_KEY");
authorization_token = 'Bearer {}'.format(jwt_token)
headers = {"Authorization": authorization_token}

ws_app = websocket.WebSocketApp("wss://api.upbit.com/websocket/v1/private",
                                header=headers,
                                on_message=on_message,
                                on_open=on_connect,
                                on_error=on_error,
                                on_close=on_close)
ws_app.run_forever()
```
```java
package com.upbit.openapi.test;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import java.util.UUID;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;
import okio.ByteString;
import org.jetbrains.annotations.NotNull;

public class AuthWebSocket {

    /**
     * JWT Token 생성 메소드 예시
     */
    private static String createJwt(String accessKey, String secretKey) {
        try {
            Algorithm algorithm = Algorithm.HMAC256(secretKey);
            return JWT.create().withClaim("access_key", accessKey)
                .withClaim("nonce", UUID.randomUUID().toString()).sign(algorithm);
        } catch (Exception e) {
            throw new RuntimeException("JWT token generation failed", e);
        }
    }

    static WebSocketListener createWebSocketListener() {
        return new WebSocketListener() {
            @Override
            public void onOpen(@NotNull WebSocket webSocket, @NotNull Response response) {
              // onOpen 이벤트 구현
							// TODO: 데이터 요청 
            }

            @Override
            public void onMessage(@NotNull WebSocket webSocket, @NotNull ByteString bytes) {
                // onMessage 이벤트 구현
            }

            @Override
            public void onClosing(@NotNull WebSocket webSocket, int code, @NotNull String reason) {
                // onClosing 이벤트 구현
            }

            @Override
            public void onFailure(@NotNull WebSocket webSocket, @NotNull Throwable t, Response response) {
                // onFailure 이벤트 구현
            }
        };
    }

    public static void main(String[] args){

        String accessKey = "YOUR_ACCESS_KEY";
        String secretKey = "YOUR_SECRET_KEY"; // 실제로는 안전하게 로드하거나 주입하세요.
      	OkHttpClient httpClient = new OkHttpClient();

        try {
            // Generate JWT token for authentication
            String jwtToken = createJwt(accessKey, secretKey);

            // Build request with authentication header
            Request request = new Request.Builder()
                .url("wss://api.upbit.com/websocket/v1/private")
                .addHeader("Authorization", "Bearer " + jwtToken)
                .build();

          WebSocket webSocket = httpClient.newWebSocket(request, createWebSocketListener());
					//
        } catch (Exception e) {
            throw new RuntimeException("Failed to connect to private WebSocket", e);
        }
    }

}

```
```node
const jwt = require("jsonwebtoken");
const {v4: uuidv4} = require('uuid');
const WebSocket = require("ws");

const payload = {
    access_key: "YOUR_ACCESS_KEY", 
    nonce: uuidv4(),
};

const jwtToken = jwt.sign(payload, "YOUR_SECRET_KEY");

const ws = new WebSocket("wss://api.upbit.com/websocket/v1/private", {
    headers: {
        authorization: `Bearer ${jwtToken}`
    }
});


ws.on("open", () => {
    console.log("connected!");
    ws.send('[{"ticket":"test example"},{"type":"myOrder"}]');

});

ws.on("error", console.error);

ws.on("message", (data) => console.log(data.toString()));

ws.on("close", () => console.log("closed!"));
```

# Sibling pages

* [개요](https://docs.upbit.com/kr/reference/api-overview.md)
* [요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits.md)
* [REST API 사용 및 에러 안내](https://docs.upbit.com/kr/reference/rest-api-guide.md)
* [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide.md)