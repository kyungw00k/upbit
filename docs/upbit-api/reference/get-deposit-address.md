---
updatedAt: 2026-11-25T18:18:15.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 개별 입금 주소 조회

지정한 통화의 입금 주소 정보를 조회힙니다.

<Callout icon="far fa-circle-info" theme="info">
  #### 조회 결과 deposit_address 필드가 null인 경우

  입금 주소 생성 요청 후 주소가 생성되기 전에 이 API를 호출하는 경우, 응답의 deposit_address 필드가 null로 반환될 수 있습니다. 일정 시간 이후 재요청 해주시기 바랍니다.
</Callout>

***

#### Rate Limit

초당 최대 30회 호출할 수 있습니다.&#x20;

포켓 단위로 측정되며 \[<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits#:~:text=%ED%98%B8%EA%B0%80%20%EC%A0%95%EC%B1%85%20%EC%A1%B0%ED%9A%8C-,Exchange%20REST%20API,-Rate%20Limit%20%EA%B7%B8%EB%A3%B9">`default` 그룹</Anchor>] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

#### API Key Permission

[인증](https://docs.upbit.com/kr/reference/auth)이 필요한 API로, \[입금조회] 권한이 설정된 API Key를 사용해야 합니다. 권한 오류(`out_of_scope`)가 발생하면 <Anchor target="_blank" href="https://www.upbit.com/mypage/open_api_management">API Key 관리</Anchor> 메뉴에서 권한 설정을 확인해 주세요.

# OpenAPI 정의

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "EXCHANGE API",
    "version": "1.0.1"
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
      "node",
      "java"
    ]
  },
  "paths": {
    "/v1/deposits/coin_address": {
      "get": {
        "operationId": "get-deposit-address",
        "summary": "개별 입금 주소 조회",
        "description": "지정한 통화의 입금 주소 정보를 조회힙니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n    --url 'https://api.upbit.com/v1/deposits/coin_address?currency=BTC&net_type=BTC' \\\n    --header 'Authorization: Bearer {JWT_TOKEN}' \\\n    --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/deposits/coin_address\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"currency\": \"BTC\",\n  \"net_type\": \"BTC\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/deposits/coin_address\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  currency: \"BTC\",\n  net_type: \"BTC\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}?${queryString}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class GetDepositAddress {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/deposits/coin_address\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"currency\", \"BTC\");\n        params.put(\"net_type\", \"BTC\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "parameters": [
          {
            "name": "currency",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "조회하고자 하는 통화 코드",
              "example": "BTC"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "BTC"
              }
            }
          },
          {
            "name": "net_type",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "디지털 자산 입출금에 사용되는 블록체인 네트워크 식별자.\n\n조회 대상을 네트워크 식별자로 한정하기 위한 필터 파라미터입니다.\n",
              "example": "BTC"
            },
            "examples": {
              "default": {
                "value": "BTC"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Object of deposit address",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/DepositAddress"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "currency": "BTC",
                      "net_type": "BTC",
                      "deposit_address": "3EusRwybuZUhVDeHL7gh3HSLmbhLcy7NqD"
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
                  "validation_error (invalid enum)": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "\"field name\" does not have a valid value"
                      }
                    }
                  },
                  "validation_error (missing parameter)": {
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
          },
          "401": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "invalid_query_payload": {
                    "value": {
                      "error": {
                        "name": "invalid_query_payload",
                        "message": "Jwt의 query를 검증하는데 실패하였습니다."
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

# 같은 수준의 페이지

* [디지털 자산 입금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-deposit-information.md)
* [입금 주소 생성 요청](https://docs.upbit.com/kr/reference/create-deposit-address.md)
* [입금 주소 목록 조회](https://docs.upbit.com/kr/reference/list-deposit-addresses.md)
* [원화 입금](https://docs.upbit.com/kr/reference/deposit-krw.md)
* [개별 입금 조회](https://docs.upbit.com/kr/reference/get-deposit.md)
* [입금 목록 조회](https://docs.upbit.com/kr/reference/list-deposits.md)
* [트래블룰 검증](https://docs.upbit.com/kr/reference/travelrule-guide.md)