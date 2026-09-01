---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 출금 허용 주소 목록 조회

등록된 출금 허용 주소 목록을 조회합니다.

### 출금 허용 주소 등록

출금 기능을 이용하기 위해서는 출금 주소 등록이 반드시 필요합니다.

[출금 허용 주소 등록 가이드](https://docs.upbit.com/kr/docs/open-api-withdraw_access_register) 를 참고하여 업비트 PC Web에서 \[마이페이지 > Open API 관리 > <Anchor target="_blank" href="www.upbit.com/mypage/opn_apimanagement?tab=fund_source">디지털 자산 출금주소 관리</Anchor>] 메뉴를 통해 출금 허용 주소를 등록해 주세요.

<br />

### 트래블룰에 따른 주소 검증 안내

출금 허용 주소 등록시 업비트는 트래블룰 준수를 위해 해당 주소를 발급한 상대 거래소로 주소 검증을 요청합니다. 상대 거래소가 이를 수신 또는 승인하지 않거나 검증 실패인 경우 출금이 제한될 수 있으므로, 출금 허용 주소 등록 전 해당 거래소에서 지원하는 주소인지 여부를 반드시 확인하시기 바랍니다.

  <div className="callout-section">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  네트워크 타입("net_type")과 네트워크 이름("network_name")
      </div>
      네트워크 타입(net_type)은 디지털 자산 입출금시 실제 자산이 이동되는 블록체인 네트워크(대상 체인)를 지정하기 위한 식별자 필드(예: BTC)입니다. 디지털 자산 출금 시 필수 파라미터로, 정상적인 입출금 진행을 위해 정확한 식별자 값을 사용해야 합니다.
      디지털 자산 출금 API 호출 시 사전에 출금 허용 주소 목록 조회 API를 호출한 뒤 응답으로부터 정확한 네트워크 타입 값을 참조하여 사용하시기 바랍니다. 
      <br /><br />
      네트워크 이름(network_name)은 블록체인 네트워크의 전체 이름(예: Bitcoin)을 나타내는 필드로서, 사람이 인식할 수 있는 정보이며 식별자로 사용할 수 없습니다. 서비스 UI 등에서 블록체인 네트워크를 표현하는 용도로 사용할 수 있습니다.
  </div>

<div className="callout-section">
  <div className="callout-title">
    <i className="fa-solid fa-circle-exclamation"></i> 출금 허용 주소의 종류에 따라 일부 응답 필드가 null로 반환될 수 있습니다.</div>
  <b>[주소 종류에 따른 응답 필드 차이]</b><br />
    <li><b>개인 지갑 주소</b>인 경우 <code>beneficiary_name</code>필드에는 회원의 이름이, <code>wallet_type</code>필드에는 개인 지갑 이름이 반환됩니다. <code>exchange_name</code>, <code>beneficiary_type</code> 필드가 null로 반환됩니다. </li>
    <li><b>거래소 지갑</b>인 경우  <code>exchange_name</code>필드에는 해당 지갑의 거래소 이름이,<code>beneficiary_name</code>필드에는 회원의 이름이 반환되며 <code>wallet_type</code> 필드가 null로 반환됩니다. <code>beneficiary_type</code>은 소유주 유형에 따라 아래와 같이 반환됩니다.</li> 
    <br /><br />
    <b>[주소 소유주 유형에 따른 응답 필드 차이]</b><br />
    <li><b>개인 소유 주소</b>인 경우 <code>beneficiary_type</code>이 <code>individual</code>로 반환되며 이때 <code>beneficiary_company_name</code> 필드는 null로 반환됩니다. </li>
    <li><b>법인 소유 주소</b>인 경우 <code>beneficiary_type</code>이 <code>corporate</code>로 반환되며 이때 <code>beneficiary_company_name</code> 필드에 해당 법인명이 반환됩니다. </li>
</div>

<div className="accordion-changelog">
    <input type="checkbox" id="api-changelog" />
    <label for="api-changelog">
        <div className="APISectionHeader-heading4MUMLbp4_nLs">Revision History <i className="fa-solid fa-angle-right"></i> </div>
    </label>

    <div className="accordion-changelog-content">
        



<table className="custom-table">
    <thead>
        <tr>
            <th>반영 버전</th>
            <th>반영 일자</th>
            <th>변경 사항</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td className="code-col">v1.5.8</td>
            <td>2025-07-07</td>
            <td><a href="https://docs.upbit.com/kr/changelog/coin_addresses_update">수신 계정 관련 필드 추가<br />(exchange_name, wallet_type, beneficiary_type,beneficiary_name, beneficiary_company_name)</a></td>
     	  </tr>
<tr>
            <td className="code-col">-</td>
            <td>2023-11-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/wallet_network_name">네트워크 명(network_name) 필드 추가</a></td>
        </tr>
<tr>
            <td className="code-col">-</td>
            <td>2023-05-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/net_type">네트워크 타입(net_type) 필드 추가</a></td>
        </tr>
    </tbody>
</table>




    </div>
</div>

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 \[Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>

  <br />
  <div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>
  <div className="box-rate-limit">
    <a href="auth">인증</a>이 필요한 API로, [출금조회] 권한이 설정된 API Key를 사용해야 합니다. <br />
    권한 오류(out_of_scope) 오류가 발생한다면, <a href="https://upbit.com/mypage/open_api_management">API Key 관리 메뉴</a>에서 권한 설정을 확인해주세요.
  </div>

<br />

# OpenAPI definition

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "EXCHANGE API",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.upbit.com"
    }
  ],
  "components": {
    "schemas": {
      "WithdrawAddress": {
        "type": "object",
        "required": [
          "currency",
          "net_type",
          "network_name",
          "withdraw_address"
        ],
        "properties": {
          "currency": {
            "type": "string",
            "description": "출금하고자 하는 디지털 자산의 통화 코드",
            "example": "BTC"
          },
          "net_type": {
            "type": "string",
            "description": "출금 네트워크 유형.\n\n업비트에서 사용하는 블록체인 네트워크 구분자입니다. 출금 요청 시 사용되는 `net_type` 파라미터는 이 필드와 같은 값을 사용해야 합니다.\n\n[예시] \"ETH\", \"TRX\", \"SOL\"\n",
            "example": "BTC"
          },
          "network_name": {
            "type": "string",
            "description": "출금 네트워크 이름.\n업비트에서 사용자에게 표시되는 블록체인 네트워크 이름입니다.\n\n[예시] \"Ethereum\", \"Bitcoin\", \"Tron\", \"Solana\"\n",
            "example": "Bitcoin"
          },
          "withdraw_address": {
            "type": "string",
            "description": "디지털 자산 출금 시 수신 주소.\n\n출금 가능 주소 목록에 등록된 주소만 사용 가능합니다.\n",
            "example": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
          },
          "secondary_address": {
            "type": "string",
            "nullable": true,
            "description": "2차 출금 주소.\n일부 디지털 자산의 경우 입출금 주소가 Destination Tag, Memo, 또는 Message와 같은 2차 주소를 포함합니다.\n"
          },
          "beneficiary_name": {
            "type": "string",
            "description": "수취 지갑 소유주명.\n\n출금된 자산을 받을 개인 또는 법인 대표의 성명입니다. 거래소 지갑인 경우에만 반환되며, 개인 지갑인 경우 null로 반환됩니다.\n"
          },
          "beneficiary_company_name": {
            "type": "string",
            "nullable": true,
            "description": "출금된 자산을 받을 법인명."
          },
          "beneficiary_type": {
            "type": "string",
            "nullable": true,
            "description": "계정주 타입.\n  - `individual`: 개인 지갑\n  - `corporate`: 법인 지갑\n"
          },
          "exchange_name": {
            "type": "string",
            "nullable": true,
            "description": "출금 허용 주소가 등록된 거래소명.\n\n[예시] \"바이낸스\", \"바이비트\"\n"
          },
          "wallet_type": {
            "type": "string",
            "nullable": true,
            "description": "개인 지갑 종류.\n\n[예시] \"메타마스크\"\n"
          }
        }
      }
    }
  },
  "x-readme": {
    "explorer-enabled": false,
    "samples-languages": [
      "shell",
      "python",
      "java",
      "node"
    ]
  },
  "paths": {
    "/v1/withdraws/coin_addresses": {
      "get": {
        "operationId": "list-withdrawal-addresses",
        "summary": "출금 허용 주소 목록 조회",
        "description": "등록된 출금 허용 주소 목록을 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n--url 'https://api.upbit.com/v1/withdraws/coin_addresses' \\\n--header 'Authorization: Bearer {JWT_TOKEN}' \\\n--header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport jwt\nimport requests\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/withdraws/coin_addresses\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/withdraws/coin_addresses\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.Objects;\nimport java.util.UUID;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListCoinAddresses {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/withdraws/coin_addresses\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "responses": {
          "200": {
            "description": "List of withdrawal allowed addresses",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/WithdrawAddress"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "currency": "BTC",
                        "net_type": "BTC",
                        "network_name": "Bitcoin",
                        "withdraw_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
                        "beneficiary_name": "John",
                        "beneficiary_type": "individual",
                        "exchange_name": "바이낸스"
                      },
                      {
                        "currency": "ETH",
                        "net_type": "ETH",
                        "network_name": "Ethereum",
                        "withdraw_address": "0x1234615148db0926d76bde31d420abcd5439484fd",
                        "beneficiary_name": "John",
                        "wallet_type": "메타마스크"
                      }
                    ]
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

# Sibling pages

* [출금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-withdrawal-information.md)
* [디지털 자산 출금 요청](https://docs.upbit.com/kr/reference/withdraw.md)
* [원화 출금 요청](https://docs.upbit.com/kr/reference/withdraw-krw.md)
* [개별 출금 조회](https://docs.upbit.com/kr/reference/get-withdrawal.md)
* [출금 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawals.md)
* [디지털 자산 출금 취소 요청](https://docs.upbit.com/kr/reference/cancel-withdrawal.md)