---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 개별 출금 조회

최신 단일 출금 정보를 조회합니다. 특정 출금 정보를 조회하고자 하는 경우 출금의 UUID 또는 트랜잭션 ID(TXID)로 지정할 수 있습니다.

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
            <td className="code-col">-</td>
            <td>2023-05-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/net_type">네트워크 타입(net_type) 필드 추가</a></td>
      </tr>
tr>
            <td className="code-col">-</td>
            <td>2020-05-29</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_feature_0529"><code>transaction_type</code> 필드 추가</a></td>
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
      "Withdraw": {
        "type": "object",
        "required": [
          "type",
          "uuid",
          "currency",
          "net_type",
          "txid",
          "state",
          "created_at",
          "amount",
          "fee",
          "transaction_type",
          "is_cancelable"
        ],
        "properties": {
          "type": {
            "type": "string",
            "description": "입출금 종류",
            "example": "withdraw",
            "default": "withdraw"
          },
          "uuid": {
            "type": "string",
            "description": "출금의 유일 식별자",
            "example": "9f432943-54e0-40b7-825f-b6fec8b42b79"
          },
          "currency": {
            "type": "string",
            "description": "조회하고자 하는 통화 코드",
            "example": "BTC"
          },
          "net_type": {
            "type": "string",
            "nullable": true,
            "description": "출금 네트워크 유형.\n\n업비트에서 사용하는 블록체인 네트워크 구분자입니다.\n\n[예시] \"ETH\", \"TRX\", \"SOL\"\n",
            "example": "BTC"
          },
          "txid": {
            "type": "string",
            "nullable": true,
            "description": "트랜잭션 ID",
            "example": 1234567890
          },
          "state": {
            "type": "string",
            "enum": [
              "WAITING",
              "PROCESSING",
              "DONE",
              "FAILED",
              "CANCELLED",
              "REJECTED"
            ],
            "description": "출금 처리 상태\n  - `WAITING`: 대기 중\n  - `PROCESSING`: 처리 중\n  - `DONE`: 완료\n  - `FAILED`: 실패\n  - `CANCELLED`: 취소됨\n  - `REJECTED`: 거절됨\n",
            "example": "WAITING"
          },
          "created_at": {
            "type": "string",
            "description": "출금 생성 시간",
            "example": "2024-01-01T00:00:00"
          },
          "done_at": {
            "type": "string",
            "nullable": true,
            "description": "출금 완료 시간. 출금이 완료되지 않은 경우 null로 반환됩니다.",
            "example": "2024-01-01T00:00:00"
          },
          "amount": {
            "type": "string",
            "format": "decimal",
            "description": "출금하고자 하는 자산의 수량",
            "example": "0.01"
          },
          "fee": {
            "type": "string",
            "format": "decimal",
            "description": "출금 수수료",
            "example": 0
          },
          "transaction_type": {
            "type": "string",
            "enum": [
              "default",
              "internal"
            ],
            "description": "출금 유형.\n\n* `default`: 일반출금\n* `internal`: 바로출금\n",
            "example": "default",
            "default": "default"
          },
          "is_cancelable": {
            "type": "boolean",
            "nullable": true,
            "description": "출금 취소 가능 여부",
            "example": false
          }
        }
      },
      "Error": {
        "type": "object",
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "name",
              "message"
            ],
            "properties": {
              "name": {
                "type": "string",
                "description": "에러명"
              },
              "message": {
                "type": "string",
                "description": "에러 메세지"
              }
            }
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
    "/v1/withdraw": {
      "get": {
        "operationId": "get-withdrawal",
        "summary": "개별 출금 조회",
        "description": "최신 단일 출금 정보를 조회합니다. 특정 출금 정보를 조회하고자 하는 경우 출금의 UUID 또는 트랜잭션 ID(TXID)로 지정할 수 있습니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n    --url 'https://api.upbit.com/v1/withdraw?currency=BTC' \\\n    --header 'Authorization: Bearer {JWT_TOKEN}' \\\n    --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/withdraw\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"currency\": \"BTC\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\nprint(jwt_token)\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/withdraw\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  currency: \"BTC\",\n};\n\nqueryString = new URLSearchParams(params).toString();\n\nqueryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\npayload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}?${queryString}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class GetWithdraw {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/withdraw\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"currency\", \"BTC\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "parameters": [
          {
            "name": "uuid",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회하고자 하는 출금의 유일식별자(UUID).\n\nuuid와 txid를 모두 입력하지 않는 경우 최신 출금 정보가 반환됩니다.\n",
              "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
            },
            "allowReserved": true
          },
          {
            "name": "txid",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회하고자 하는 출금의 트랜잭션 ID.\n\nuuid와 txid를 모두 입력하지 않는 경우 최신 출금 정보가 반환됩니다.\n",
              "example": "98c15999f0bdc4ae0e8a-ed35868bb0c204fe6ec29e4058a3451e-88636d1040f4baddf943274ce37cf9cc"
            },
            "allowReserved": true
          },
          {
            "name": "currency",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회하고자 하는 통화 코드.\n통화 코드로 조회 대상을 한정하기 위한 필터 파라미터입니다.\n미입력시 최신 출금 내역이 조회됩니다.\n",
              "example": "KRW"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "BTC"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Object of withdrawal",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Withdraw"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "type": "withdraw",
                      "uuid": "9f432943-54e0-40b7-825f-b6fec8b42b79",
                      "currency": "BTC",
                      "state": "PROCESSING",
                      "created_at": "2025-07-04T15:00:00+09:00",
                      "amount": "0.01",
                      "fee": "0.0",
                      "transaction_type": "default",
                      "is_cancelable": false
                    }
                  }
                }
              }
            }
          },
          "404": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "not found withdraw error": {
                    "value": {
                      "error": {
                        "name": "withdraw_not_found",
                        "message": "출금 정보를 찾지 못했습니다."
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
  }
}
```

# Sibling pages

* [출금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-withdrawal-information.md)
* [출금 허용 주소 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawal-addresses.md)
* [디지털 자산 출금 요청](https://docs.upbit.com/kr/reference/withdraw.md)
* [원화 출금 요청](https://docs.upbit.com/kr/reference/withdraw-krw.md)
* [출금 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawals.md)
* [디지털 자산 출금 취소 요청](https://docs.upbit.com/kr/reference/cancel-withdrawal.md)