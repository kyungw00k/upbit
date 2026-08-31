---
updatedAt: 2026-07-28T11:49:23.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 포켓 잔고 조회

포켓이 보유하고 있는 자산 목록과 잔고를 조회합니다.

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 \[Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>
<br />
<div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>
<div className="box-rate-limit">
  <a href="auth">인증</a>이 필요한 API로, [자산조회] 권한이 설정된 API Key를 사용해야 합니다. <br />
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
      "Account": {
        "type": "object",
        "required": [
          "currency",
          "balance",
          "locked",
          "avg_buy_price",
          "avg_buy_price_modified",
          "unit_currency"
        ],
        "properties": {
          "currency": {
            "type": "string",
            "description": "조회하고자 하는 통화 코드",
            "example": "BTC"
          },
          "balance": {
            "type": "string",
            "format": "decimal",
            "description": "주문 가능 수량 또는 금액.\n\n디지털 자산의 경우 수량, 법정 통화(KRW)의 경우 금액입니다.\n",
            "example": 1000000
          },
          "locked": {
            "type": "string",
            "format": "decimal",
            "description": "출금이나 주문 등에 잠겨 있는 잔액",
            "example": 0
          },
          "avg_buy_price": {
            "type": "string",
            "format": "decimal",
            "description": "매수 평균가",
            "example": 0
          },
          "avg_buy_price_modified": {
            "type": "boolean",
            "description": "매수 평균가 수정 여부",
            "example": false
          },
          "unit_currency": {
            "type": "string",
            "description": "평균가 기준 통화.\n\n\"avg_buy_price\"가 기준하는 단위입니다.\n\n[예시] KRW, BTC, USDT\n",
            "example": "KRW"
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
    "/v1/accounts": {
      "get": {
        "operationId": "get-balance",
        "summary": "포켓 잔고 조회",
        "description": "포켓이 보유하고 있는 자산 목록과 잔고를 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n--url 'https://api.upbit.com/v1/accounts' \\\n--header 'Authorization: Bearer {JWT_TOKEN}' \\\n--header 'accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport jwt\nimport requests\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/accounts\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid\n",
              "code": "const axios = require(\"axios\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/accounts\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.Objects;\nimport java.util.UUID;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class GetBalance {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/accounts\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "responses": {
          "200": {
            "description": "List of account balances",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Account"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "currency": "KRW",
                        "balance": "1000000.0",
                        "locked": "0.0",
                        "avg_buy_price": "0",
                        "avg_buy_price_modified": false,
                        "unit_currency": "KRW"
                      },
                      {
                        "currency": "BTC",
                        "balance": "2.0",
                        "locked": "0.0",
                        "avg_buy_price": "140000000",
                        "avg_buy_price_modified": false,
                        "unit_currency": "KRW"
                      }
                    ]
                  }
                }
              }
            }
          },
          "401": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "no authorization ip error": {
                    "value": {
                      "error": {
                        "name": "no_authorization_ip",
                        "message": "This is not a verified IP."
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