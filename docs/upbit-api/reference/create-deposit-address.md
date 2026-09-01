---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 입금 주소 생성 요청

개인 지갑 또는 타 거래소 자산을 업비트로 입금 하기 위한 입금 주소 생성을 요청합니다.

### 비동기 방식 주소 생성으로 인한 API 응답 객체 구분

입금 주소 생성은 비동기 방식으로 동작합니다. API 호출 시점의 입금 주소 생성 완료 여부에 따라 아래 두가지 응답을 반환할 수 있습니다.

1. 최초 API 요청 직후 반환되는 응답은 주소 생성 요청의 **접수 성공 여부를 반환**하며, 응답은 `success`, `message` 필드만 반환됩니다. 해당 응답은 API 추가 호출 시 주소 생성이 완료되기 이전까지 반환됩니다.
2. 비동기 방식으로 주소 생성이 완료된 이후 API 호출 시 "currency", "net\_type", "deposit\_address"를 포함하는 **생성된 주소 정보가 반환**됩니다. 해당 정보는 통화당 최초 1회 생성 되며, 이후 생성 요청의 응답으로는 기존에 생성된 주소 정보가 반환됩니다.

일정 시간이 지난 후에도 입금 주소가 정상적으로 생성되지 않는 경우, 시간 간격을 두고 이 API를 다시 호출해주시기 바랍니다.

<div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> POST API에 대한 Form 방식 요청은 2022년 3월 1일부로 지원이 종료되었습니다.
      </div>
    Form 방식 지원 종료에 따라 Urlencoded Form 방식으로 전송하는 POST 요청에 대한 정상적인 동작을 보장하지 않습니다. <b>반드시 JSON 형식으로 요청 본문(Body)을 전송</b>해주시기 바랍니다.
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
    <a href="auth">인증</a>이 필요한 API로, [입금하기] 권한이 설정된 API Key를 사용해야 합니다. <br />
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
      "DepositAddress": {
        "type": "object",
        "required": [
          "currency",
          "net_type",
          "deposit_address"
        ],
        "properties": {
          "currency": {
            "type": "string",
            "description": "입금 주소가 생성된 통화 코드",
            "example": "BTC"
          },
          "net_type": {
            "type": "string",
            "nullable": true,
            "description": "입금 네트워크 유형.\n\n업비트에서 사용하는 블록체인 네트워크 구분자입니다.\n\n[예시] ETH, TRX, SOL\n",
            "example": "BTC"
          },
          "deposit_address": {
            "type": "string",
            "description": "입금 주소",
            "example": "3GXAGnqLWpZWiChDU2AsJBaVxpnPiLBaxU"
          },
          "secondary_address": {
            "type": "string",
            "nullable": true,
            "description": "2차 출금 주소.\n일부 디지털 자산의 경우 입출금 주소가 Destination Tag, Memo, 또는 Message와 같은 2차 주소를 포함합니다.\n"
          }
        }
      },
      "GenerateCoinAddressResponse": {
        "type": "object",
        "required": [
          "success",
          "message"
        ],
        "properties": {
          "success": {
            "type": "boolean",
            "description": "입금 주소 생성 요청의 성공 여부",
            "example": true,
            "default": true
          },
          "message": {
            "type": "string",
            "description": "입금 주소 생성 요청 결과에 대한 메시지",
            "example": "BTC 입금 주소를 생성 중 입니다."
          }
        }
      },
      "CreateDepositAddressResponse": {
        "oneOf": [
          {
            "$ref": "#/components/schemas/DepositAddress"
          },
          {
            "$ref": "#/components/schemas/GenerateCoinAddressResponse"
          }
        ]
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
    "/v1/deposits/generate_coin_address": {
      "post": {
        "operationId": "create-deposit-address",
        "summary": "입금 주소 생성 요청",
        "description": "개인 지갑 또는 타 거래소 자산을 업비트로 입금 하기 위한 입금 주소 생성을 요청합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request POST \\\n  --url 'https://api.upbit.com/v1/deposits/generate_coin_address' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Content-Type: application/json' \\\n  --data '\n{\n\"currency\": \"BTC\",\n\"net_type\": \"BTC\"\n}\n'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/deposits/generate_coin_address\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"currency\": \"BTC\",\n  \"net_type\": \"BTC\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/deposits/generate_coin_address\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  currency: \"BTC\",\n  net_type: \"BTC\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"POST\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n  data: params,\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\nimport com.google.gson.Gson;\n\npublic class CreateDepositAddress {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/deposits/generate_coin_address\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"currency\", \"BTC\");\n        params.put(\"net_type\", \"BTC\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        String jsonBody = new Gson().toJson(params);\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .post(RequestBody.create(jsonBody, okhttp3.MediaType.parse(\"application/json; charset=utf-8\")))\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "requestBody": {
          "description": "입금 주소 생성 요청",
          "content": {
            "application/json": {
              "examples": {
                "default": {
                  "value": {
                    "currency": "BTC",
                    "net_type": "BTC"
                  }
                }
              },
              "schema": {
                "type": "object",
                "required": [
                  "currency",
                  "net_type"
                ],
                "properties": {
                  "currency": {
                    "type": "string",
                    "description": "입금 주소가 생성된 통화 코드",
                    "example": "BTC"
                  },
                  "net_type": {
                    "type": "string",
                    "description": "네트워크 유형",
                    "example": "BTC"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Deposit address (existing or creation in progress)",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/CreateDepositAddressResponse"
                },
                "examples": {
                  "Address ready": {
                    "value": {
                      "currency": "BTC",
                      "net_type": "BTC",
                      "deposit_address": "3EusRwybuZUhVDeHL7gh3HSLmbhLcy7NqD"
                    }
                  },
                  "Address generating": {
                    "value": {
                      "success": true,
                      "message": "BTC 입금 주소를 생성 중 입니다."
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "invalid enum error": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "\"field name\" does not have a valid value"
                      }
                    }
                  },
                  "invalid parameter error": {
                    "value": {
                      "error": {
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
                      }
                    }
                  },
                  "missing parameter error": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "\"field name\" is missing"
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

* [디지털 자산 입금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-deposit-information.md)
* [개별 입금 주소 조회](https://docs.upbit.com/kr/reference/get-deposit-address.md)
* [입금 주소 목록 조회](https://docs.upbit.com/kr/reference/list-deposit-addresses.md)
* [원화 입금](https://docs.upbit.com/kr/reference/deposit-krw.md)
* [개별 입금 조회](https://docs.upbit.com/kr/reference/get-deposit.md)
* [입금 목록 조회](https://docs.upbit.com/kr/reference/list-deposits.md)
* [트래블룰 검증](https://docs.upbit.com/kr/reference/travelrule-guide.md)