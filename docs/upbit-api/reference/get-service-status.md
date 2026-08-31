---
updatedAt: 2026-07-28T12:03:06.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 입출금 서비스 상태 조회

전체 통화에 대해 입출금 서비스 상태를 조회합니다.

<div class="callout-section callout-section--danger">
    <div class="callout-title">
      <i class="fa-solid fa-circle-exclamation"></i> 입출금 서비스 상태 조회 API는 실시간 상태 조회를 보장하지 않습니다.
      </div>
      입출금 서비스 상태 조회 API가 반환하는 입출금 가능 여부는 서비스 상태를 실시간으로 반영하지 않으며 반영은 수 분 정도 지연될 수 있습니다. 
      따라서 <b>거래 전략 용도가 아닌 참고 용도로의 사용만을 권장</b>하며, 실제 입금을 수행하기 전에는 반드시 <a href="https://upbit.com/service_center/notice">업비트 공지사항</a> 및 <a href="https://upbit.com/service_center/wallet_status">실시간 입출금 현황</a> 페이지를 참고해 주시기를 바랍니다.
  </div>

  <div class="callout-section">
    <div class="callout-title">
      <i class="fa-solid fa-circle-exclamation"></i>  네트워크 타입("net_type")과 네트워크 이름("network_name")
      </div>
      네트워크 타입("net_type")은 디지털 자산 입출금시 실제 자산이 이동되는 블록체인 네트워크(대상 체인)를 지정하기 위한 식별자 필드(예: BTC)입니다. 디지털 자산 출금 시 필수 파라미터로, 정상적인 입출금 진행을 위해 정확한 식별자 값을 사용해야 합니다.
      디지털 자산 출금 API 호출 시 사전에 출금 허용 주소 목록 조회 API를 호출한 뒤 응답으로부터 정확한 네트워크 타입 값을 참조하여 사용하시기 바랍니다. 
      <br><br>
      네트워크 이름("network_name")은 블록체인 네트워크의 전체 이름(예: Bitcoin)을 나타내는 필드로서, 사람이 인식할 수 있는 정보이며 식별자로 사용할 수 없습니다. 서비스 UI 등에서 블록체인 네트워크를 표현하는 용도로 사용할 수 있습니다.
  </div>

<div class="accordion-changelog">
    <input type="checkbox" id="api-changelog">
    <label for="api-changelog">
        <div class="APISectionHeader-heading4MUMLbp4_nLs">Revision History <i class="fa-solid fa-angle-right"></i> </div>
    </label>

    <div class="accordion-changelog-content">
        <table class="custom-table">
            <thead>
                <tr>
                    <th>반영 버전</th>
                    <th>반영 일자</th>
                    <th>변경 사항</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td class="code-col">-</td>
                    <td>2023-11-22</td>
                    <td><a href="https://docs.upbit.com/kr/changelog/wallet_network_name">네트워크 명(network_name) 필드 추가</a></td>
              	</tr>
								<tr>
                    <td class="code-col">-</td>
                    <td>2023-11-22</td>
                    <td><a href="https://docs.upbit.com/kr/changelog/net_type">네트워크 타입(net_type) 필드 추가</a></td>
                </tr>
								<tr>
                    <td class="code-col">-</td>
                    <td></td>
                    <td><a href="https://docs.upbit.com/kr/changelog/wallet_state">'입출금 서비스 상태 조회' 기능 신규 지원</a></td>
                </tr>
            </tbody>
        </table>
    </div>
</div>

<div class="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div class="box-rate-limit">
  초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 [Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.
</div>

 <br>
  <div class="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>
  <div class="box-rate-limit">
    <a href="auth">인증</a>이 필요한 API 입니다. 별도 권한은 필요하지 않습니다.
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
      "WalletStatus": {
        "type": "object",
        "required": [
          "currency",
          "wallet_state",
          "block_elapsed_minutes",
          "net_type",
          "network_name"
        ],
        "properties": {
          "currency": {
            "type": "string",
            "description": "조회하고자 하는 통화 코드",
            "example": "BTC"
          },
          "wallet_state": {
            "type": "string",
            "enum": [
              "working",
              "withdraw_only",
              "deposit_only",
              "paused",
              "unsupported"
            ],
            "description": "입출금 상태.\n  - `working` : 입출금 가능\n  - `withdraw_only` : 출금만 가능\n  - `deposit_only` : 입금만 가능\n  - `paused` : 입출금 중단\n  - `unsupported` : 입출금 미지원\n",
            "example": "working"
          },
          "block_state": {
            "type": "string",
            "enum": [
              "normal",
              "delayed",
              "inactive"
            ],
            "description": "블록체인 네트워크의 상태.\n지갑 또는 거래소 상태에 따라 null로 반환 될 수 있습니다.\n  - `normal`: 정상\n  - `delayed`: 지연\n  - `inactive`: 비활성\n",
            "example": "normal"
          },
          "block_height": {
            "type": "integer",
            "description": "현재 확인된 블록의 높이.\n지갑 또는 거래소 상태에 따라 null로 반환 될 수 있습니다.\n",
            "example": 902656
          },
          "block_updated_at": {
            "type": "string",
            "description": "마지막으로 블록 높이가 갱신된 시각 (UTC).\n지갑 또는 거래소 상태에 따라 null로 반환 될 수 있습니다.\n\n[형식] yyyy-MM-dd'T'HH:mm:sss+00:00\n",
            "example": "2024-01-01T00:00:000+00:00"
          },
          "block_elapsed_minutes": {
            "type": "integer",
            "description": "마지막 블록 업데이트 이후 현재까지 경과한 시간(분).\n\n지갑 또는 거래소 상태에 따라 null로 반환 될 수 있습니다.\n",
            "example": 31
          },
          "net_type": {
            "type": "string",
            "description": "입출금 네트워크 유형.\n\n업비트에서 사용하는 블록체인 네트워크 구분자입니다.\n\n[예시] \"ETH\", \"TRX\", \"SOL\"\n",
            "example": "BTC"
          },
          "network_name": {
            "type": "string",
            "description": "입출금 네트워크 이름.\n업비트에서 사용자에게 표시되는 블록체인 네트워크 이름입니다.\n\n[예시] \"Ethereum\", \"Bitcoin\", \"Tron\", \"Solana\"\n",
            "example": "Bitcoin"
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
    "/v1/status/wallet": {
      "get": {
        "operationId": "get-service-status",
        "summary": "입출금 서비스 상태 조회",
        "description": "전체 통화에 대해 입출금 서비스 상태를 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n    --url 'https://api.upbit.com/v1/status/wallet' \\\n    --header 'Authorization: Bearer {JWT_TOKEN}' \\\n    --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport jwt\nimport requests\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/status/wallet\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers={\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"})\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/status/wallet\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => console.log(response.data))\n  .catch((error) => console.error(error));\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class StatusService {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/status/wallet\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"access_key\", accessKey);\n        params.put(\"nonce\", UUID.randomUUID().toString());\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "responses": {
          "200": {
            "description": "List of deposit/withdrawal service status",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/WalletStatus"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "currency": "BTC",
                        "wallet_state": "working",
                        "block_state": "normal",
                        "block_height": 903942,
                        "block_updated_at": "2025-07-04T08:02:05.526+00:00",
                        "block_elapsed_minutes": 6,
                        "net_type": "BTC",
                        "network_name": "Bitcoin"
                      },
                      {
                        "currency": "ETH",
                        "wallet_state": "working",
                        "block_state": "normal",
                        "block_height": 22844550,
                        "block_updated_at": "2025-07-04T08:06:44.375+00:00",
                        "block_elapsed_minutes": 2,
                        "net_type": "ETH",
                        "network_name": "Ethereum"
                      },
                      {
                        "currency": "XRP",
                        "wallet_state": "working",
                        "block_state": "normal",
                        "block_height": 97241570,
                        "block_updated_at": "2025-07-04T08:06:53.213+00:00",
                        "block_elapsed_minutes": 2,
                        "net_type": "XRP",
                        "network_name": "XRP Ledger"
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

* [API Key 목록 조회](https://docs.upbit.com/kr/reference/list-api-keys.md)