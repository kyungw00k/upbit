---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 진행 중 코인빌리기 목록 조회

이용 중인 코인빌리기 목록을 조회합니다.

해당 포켓에서 신청한 코인빌리기 중 종료되지 않은 포지션 목록을 조회합니다. `borrowed_currency`를 지정하여 빌린 자산별로 조회할 수 있습니다.

#### 조회 대상

| 포지션 상태        | 설명              |
| ------------- | --------------- |
| `active`      | 이용 중            |
| `suspended`   | 이용 제한           |
| `liquidating` | 청산 또는 담보상환 진행 중 |

* `liquidating` 상태에는 강제상환(청산)과 담보상환 요청의 처리가 모두 포함됩니다. 처리 중에는 두 유형을 구분할 수 없으며, 종료 후 **종료된 코인빌리기 목록 조회** 또는 **개별 코인빌리기 조회**의 `closure_type`으로 확인할 수 있습니다.
* 신청 요청이 접수되었더라도 포지션이 아직 생성되지 않았다면 조회 결과에 포함되지 않습니다. 신청 처리 결과는 **코인빌리기 요청 상태 조회**에서 확인해 주세요.
* 종료된 포지션은 **종료된 코인빌리기 목록 조회**에서 확인할 수 있습니다.
* 목록 조회 응답에는 현재 LTV(렌딩비율)가 포함되지 않습니다. 현재 LTV는 **개별 코인빌리기 조회**에서 확인해 주세요.

#### 담보 인정 비율

`lendings[].collaterals[].factor`는 해당 코인빌리기 포지션에 적용된 담보 인정 비율입니다.

담보 인정 금액은 담보 평가 금액에 담보 인정 비율을 곱하여 계산합니다.

**담보 인정 금액 = 담보 평가 금액 × 담보 인정 비율**

#### 조회 기준 및 제약 사항

1. 별도의 조회 조건을 지정하지 않으면 해당 포켓의 종료되지 않은 포지션을 생성 시각(`created_at`) 기준 최신순으로 20건 조회합니다. `limit`을 지정하면 한 번에 최대 100건까지 조회할 수 있습니다.
2. 목록은 커서 기반 페이지네이션 방식으로 조회합니다.
   * 최초 요청에서는 `cursor`를 지정하지 않습니다.
   * 다음 페이지가 있는 경우 `next_cursor`에 현재 페이지의 마지막 포지션 UUID가 반환됩니다. 반환된 값을 다음 요청의 `cursor`에 입력하면 해당 포지션 이후의 목록을 조회할 수 있습니다.
   * 다음 페이지를 조회할 때는 `borrowed_currency` 등 다른 조회 조건을 이전 요청과 동일하게 유지해야 합니다.
   * `next_cursor`가 `null`이면 마지막 페이지입니다.
3. 여러 페이지에 걸쳐 조회하는 경우 동일한 시점의 데이터가 보장되지 않습니다. 조회 중 새로운 포지션이 생성되거나 기존 포지션의 상태가 변경될 수 있습니다.
4. 조회 조건에 맞는 포지션이 없는 경우 다음과 같이 빈 목록을 반환합니다.

```json
{
  "lendings": [],
  "next_cursor": null
}
```

***

#### 에러 안내

코인빌리기 API 이용 중 에러가 발생하면 반환된 에러 코드와 발생 사유를 확인하고, 해당 조건에 맞게 조치한 뒤 다시 요청해 주세요.

발생 가능한 에러와 상세 응답 명세는 [코인빌리기 API 안내](https://docs.upbit.com/kr/reference/lending-overview)의 **에러 안내** 항목을 참고해 주세요.

#### Rate Limit

초당 최대 30회 호출할 수 있습니다.&#x20;

포켓 단위로 측정되며 \[`default` 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

#### API Key Permission

[인증](https://docs.upbit.com/kr/reference/auth)이 필요한 API로, 메인포켓 API Key만 사용할 수 있습니다. \[코인빌리기] 권한이 설정된 API Key를 사용해야 합니다. 권한 오류(`out_of_scope`)가 발생하면 <Anchor target="_blank" href="https://www.upbit.com/mypage/open_api_management">API Key 관리</Anchor> 메뉴에서 권한 설정을 확인해 주세요.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### 코인빌리기 API 이용 권한

  코인빌리기 권한은 API Key 발급 시 이용 자격(고객 확인, 약관 동의 등)을 확인한 뒤 설정할 수 있습니다. DAXA 온라인 교육을 이수하지 않아도 권한 설정과 조회 API 이용은 가능합니다. 이후 동의를 철회하거나 이용 자격이 바뀌면 공통 에러가 반환됩니다.
</Callout>

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
      },
      "OngoingLendings": {
        "type": "object",
        "required": [
          "lendings",
          "next_cursor"
        ],
        "properties": {
          "lendings": {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "uuid",
                "state",
                "liquidation_ltv",
                "borrowed_currency",
                "borrowed_amount",
                "collaterals",
                "created_at",
                "expire_at",
                "updated_at"
              ],
              "properties": {
                "uuid": {
                  "type": "string",
                  "description": "포지션의 UUID"
                },
                "state": {
                  "type": "string",
                  "enum": [
                    "active",
                    "liquidating",
                    "suspended"
                  ],
                  "description": "포지션 상태\n\n* `active`: 이용 중\n* `liquidating`: 청산·담보상환 중\n* `suspended`: 이용제한"
                },
                "liquidation_ltv": {
                  "type": "string",
                  "description": "해당 포지션에 적용된 강제상환 기준 LTV(렌딩비율)"
                },
                "borrowed_currency": {
                  "type": "string",
                  "description": "빌린 자산의 자산 코드"
                },
                "borrowed_amount": {
                  "type": "string",
                  "description": "조회 시점의 남은 빌린 수량"
                },
                "collaterals": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "required": [
                      "currency",
                      "amount",
                      "factor"
                    ],
                    "properties": {
                      "currency": {
                        "type": "string",
                        "description": "담보자산의 자산 코드"
                      },
                      "amount": {
                        "type": "string",
                        "description": "조회 시점의 담보 수량"
                      },
                      "factor": {
                        "type": "string",
                        "description": "해당 포지션에 적용된 담보 인정 비율"
                      }
                    }
                  },
                  "description": "담보 목록"
                },
                "created_at": {
                  "type": "string",
                  "description": "빌리기 시작 시각 (ISO 8601)"
                },
                "expire_at": {
                  "type": "string",
                  "description": "최초 약정된 만기일 (ISO 8601)"
                },
                "updated_at": {
                  "type": "string",
                  "description": "정보가 마지막으로 변경된 시각 (ISO 8601)"
                }
              }
            },
            "description": "조회 조건을 만족하는 포지션 목록"
          },
          "next_cursor": {
            "type": "string",
            "description": "다음 페이지 조회에 사용할 커서",
            "nullable": true
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
    "/v1/lendings/ongoing": {
      "get": {
        "operationId": "list-ongoing-lendings",
        "summary": "진행 중 코인빌리기 목록 조회",
        "description": "이용 중인 코인빌리기 목록을 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "parameters": [
          {
            "name": "borrowed_currency",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회할 빌린 자산의 자산 코드"
            }
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "한 번에 조회할 포지션 수",
              "default": 20,
              "minimum": 1,
              "maximum": 100
            }
          },
          {
            "name": "cursor",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "다음 페이지 조회의 기준이 되는 포지션의 UUID"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/OngoingLendings"
                },
                "examples": {
                  "Example": {
                    "description": "정상 응답",
                    "value": {
                      "lendings": [
                        {
                          "uuid": "5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b",
                          "state": "active",
                          "liquidation_ltv": "0.9500",
                          "borrowed_currency": "USDT",
                          "borrowed_amount": "50000.00000000",
                          "collaterals": [
                            {
                              "currency": "BTC",
                              "amount": "1.00000000",
                              "factor": "0.90"
                            }
                          ],
                          "created_at": "2026-09-01T10:15:32+09:00",
                          "expire_at": "2026-10-01T10:15:32+09:00",
                          "updated_at": "2026-09-06T21:42:18+09:00"
                        },
                        {
                          "uuid": "942662c8-4587-489d-b775-6a58aa9ce47f",
                          "state": "liquidating",
                          "liquidation_ltv": "0.9500",
                          "borrowed_currency": "USDT",
                          "borrowed_amount": "35000.00000000",
                          "collaterals": [
                            {
                              "currency": "BTC",
                              "amount": "0.41000000",
                              "factor": "0.90"
                            },
                            {
                              "currency": "ETH",
                              "amount": "2.50000000",
                              "factor": "0.90"
                            }
                          ],
                          "created_at": "2026-08-18T13:42:17+09:00",
                          "expire_at": "2026-09-17T13:42:17+09:00",
                          "updated_at": "2026-09-06T22:03:41+09:00"
                        },
                        {
                          "uuid": "2b214968-5e62-4414-9dc4-9bc66f024b15",
                          "state": "suspended",
                          "liquidation_ltv": "0.9500",
                          "borrowed_currency": "BTC",
                          "borrowed_amount": "0.25000000",
                          "collaterals": [
                            {
                              "currency": "ETH",
                              "amount": "8.50000000",
                              "factor": "0.90"
                            },
                            {
                              "currency": "USDT",
                              "amount": "12000.00000000",
                              "factor": "0.95"
                            }
                          ],
                          "created_at": "2026-08-10T16:25:03+09:00",
                          "expire_at": "2026-09-09T16:25:03+09:00",
                          "updated_at": "2026-09-06T18:31:27+09:00"
                        }
                      ],
                      "next_cursor": "2b214968-5e62-4414-9dc4-9bc66f024b15"
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
                  "invalid_parameter": {
                    "description": "* limit이 100 초과\n* limit이 0 또는 음수\n* limit이 정수가 아님\n* borrowed_currency - 잘못된 형식\n* 허용되지 않는 자산 코드 - 빈 문자열, 소문자 - 존재하지 않는 코드 등\n* cursor가 UUID 형식이 아님",
                    "value": {
                      "error": {
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
                      }
                    }
                  }
                }
              }
            }
          },
          "403": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "out_of_scope": {
                    "description": "코인빌리기(`lending`) 권한이 없는 API Key로 호출",
                    "value": {
                      "error": {
                        "name": "out_of_scope",
                        "message": "권한이 부족합니다."
                      }
                    }
                  }
                }
              }
            }
          }
        },
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/lendings/ongoing?limit=3' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/ongoing\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\"limit\": 3}\nquery_string = unquote(urlencode(params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"}\nres = requests.get(f\"{BASE_URL}{PATH}?{query_string}\", headers=headers)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/ongoing\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = { limit: 3 };\nconst queryString = Object.entries(params)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .get(`${BASE_URL}${PATH}?${queryString}`, {\n    headers: { Authorization: `Bearer ${jwtToken}`, Accept: \"application/json\" },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListOngoingLendings {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/ongoing\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new LinkedHashMap<>();\n        params.put(\"limit\", \"3\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .get()\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
            }
          ]
        }
      }
    }
  }
}
```

# 같은 수준의 페이지

* [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information.md)
* [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test.md)
* [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending.md)
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request.md)
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)