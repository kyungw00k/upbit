---
updatedAt: 2026-11-25T18:26:15.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 입금 TxID로 계정주 검증 요청

입금의 TxID로 트래블룰 검증(입금 계정주 확인)을 요청합니다.

<Callout icon="far fa-circle-check" theme="success">
  #### 계정주 검증이 완료된 이후 입금 반영이 진행됩니다.

  계정주 검증이 완료되지 않는 경우 업비트 신원 확인 과정에서 입력한 이름, 생년월일 정보와 상대 거래소에 입력된 정보가 동일한지 여부를 확인해주세요.
  검증이 완료 시 즉시 입금 반영이 진행되나, 블록체인 네트워크 지연등으로 검증 이후에도 입금 자체가 아직 처리중인 경우 입금이 완료된 이후 입금 반영이 진행됩니다.
</Callout>

<Callout icon="far fa-octagon-exclamation" theme="error">
  #### POST API에 대한 Form 방식 요청은 2022년 3월 1일부로 지원이 종료되었습니다.

  Form 방식 지원 종료에 따라 Urlencoded Form 방식으로 전송하는 POST 요청에 대한 정상적인 동작을 보장하지 않습니다.**&#x20;반드시 JSON 형식으로 요청 본문(Body)을 전송**해주시기 바랍니다.
</Callout>

***

#### Rate Limit

초당 최대 30회 호출할 수 있습니다.&#x20;

포켓 단위로 측정되며 \[<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits#:~:text=%ED%98%B8%EA%B0%80%20%EC%A0%95%EC%B1%85%20%EC%A1%B0%ED%9A%8C-,Exchange%20REST%20API,-Rate%20Limit%20%EA%B7%B8%EB%A3%B9">`default` 그룹</Anchor>] 내에서 요청 가능 횟수를 공유합니다. 단, <span style="color: #C92532;">동일 입금건에 대해서는 10분당 최대 1회 호출할 수 있습니다</span>. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

#### API Key Permission

[인증](https://docs.upbit.com/kr/reference/auth)이 필요한 API로, \[입금하기] 권한이 설정된 API Key를 사용해야 합니다. 권한 오류(`out_of_scope`)가 발생하면 <Anchor target="_blank" href="https://www.upbit.com/mypage/open_api_management">API Key 관리</Anchor> 메뉴에서 권한 설정을 확인해 주세요.

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
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "요청 쿼리 파라미터 또는 본문(Body)을 API Key로 서명한 JWT Token을 입력합니다.\n자세한 사항은 https://docs.upbit.com/reference/auth 를 참고해주세요.\n\n`예시` Bearer <ACCESS_TOKEN>\n"
      }
    },
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
      "TravelRuleVerificationResult": {
        "type": "object",
        "required": [
          "deposit_uuid",
          "verification_result",
          "deposit_state"
        ],
        "properties": {
          "deposit_uuid": {
            "type": "string",
            "description": "입금의 유일식별자(UUID)",
            "example": "5b871d34-fe38-4025-8f5c-9b22028f85d3"
          },
          "deposit_state": {
            "type": "string",
            "enum": [
              "PROCESSING",
              "ACCEPTED",
              "CANCELLED",
              "REJECTED",
              "TRAVEL_RULE_SUSPECTED",
              "REFUNDING",
              "REFUNDED"
            ],
            "description": "입금 처리 상태\n  - `PROCESSING`: 입금 진행 중\n  - `ACCEPTED`: 완료\n  - `CANCELLED`: 취소됨\n  - `REJECTED`: 거절됨\n  - `TRAVEL_RULE_SUSPECTED`: 트래블룰 추가 인증 대기중\n  - `REFUNDING`: 반환 절차 중\n  - `REFUNDED`: 반환됨\n",
            "example": "PROCESSING"
          },
          "verification_result": {
            "type": "string",
            "description": "검증 결과",
            "example": "verified"
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
    "/v1/travel_rule/deposit/txid": {
      "post": {
        "operationId": "verify-travelrule-by-txid",
        "summary": "입금 TxID로 계정주 검증 요청",
        "description": "입금의 TxID로 트래블룰 검증(입금 계정주 확인)을 요청합니다.",
        "tags": [
          "Exchange"
        ],
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request POST \\\n  --url 'https://api.upbit.com/v1/travel_rule/deposit/txid' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Content-Type: application/json' \\\n  --data '\n{\n\"txid\": \"5b871d34-fe38-4025-8f5c-9b22028f85d3\",\n\"vasp_uuid\": \"8d4fe968-82b2-42e5-822f-3840a245f802\",\n\"currency\": \"ETH\",\n\"net_type\": \"ETH\"\n}\n'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/travel_rule/deposit/txid\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"txid\": \"5b871d34-fe38-4025-8f5c-9b22028f85d3\",\n  \"vasp_uuid\": \"8d4fe968-82b2-42e5-822f-3840a245f802\",\n  \"currency\": \"ETH\",\n  \"net_type\": \"ETH\"\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/travel_rule/deposit/txid\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  txid: \"5b871d34-fe38-4025-8f5c-9b22028f85d3\",\n  vasp_uuid: \"8d4fe968-82b2-42e5-822f-3840a245f802\",\n  currency: \"ETH\",\n  net_type: \"ETH\"\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"POST\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n  data: params,\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\nimport com.google.gson.Gson;\n\npublic class CreateDepositAddress {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/travel_rule/deposit/txid\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"txid\", \"5b871d34-fe38-4025-8f5c-9b22028f85d3\");\n        params.put(\"vasp_uuid\", \"8d4fe968-82b2-42e5-822f-3840a245f802\");\n        params.put(\"currency\", \"ETH\");\n        params.put(\"net_type\", \"ETH\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        String jsonBody = new Gson().toJson(params);\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .post(RequestBody.create(jsonBody, okhttp3.MediaType.parse(\"application/json; charset=utf-8\")))\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "requestBody": {
          "description": "txid로 입금 트랜잭션 트래블룰 검증",
          "content": {
            "application/json": {
              "examples": {
                "default": {
                  "value": {
                    "txid": "5b871d34-fe38-4025-8f5c-9b22028f85d3",
                    "vasp_uuid": "8d4fe968-82b2-42e5-822f-3840a245f802",
                    "currency": "ETH",
                    "net_type": "ETH"
                  }
                }
              },
              "schema": {
                "type": "object",
                "required": [
                  "vasp_uuid",
                  "txid",
                  "currency",
                  "net_type"
                ],
                "properties": {
                  "vasp_uuid": {
                    "type": "string",
                    "description": "자산을 출금한 상대 거래소의 유일식별자(UUID)",
                    "example": "8d4fe968-82b2-42e5-822f-3840a245f802"
                  },
                  "txid": {
                    "type": "string",
                    "description": "검증할 입금의 트랜잭션 ID",
                    "example": 1234567890
                  },
                  "currency": {
                    "type": "string",
                    "description": "검증하고자 하는 통화 코드",
                    "example": "BTC"
                  },
                  "net_type": {
                    "type": "string",
                    "description": "디지털 자산 입출금에 사용되는 블록체인 네트워크 식별자.\n\n검증 대상을 네트워크 식별자로 한정하기 위한 필터 파라미터입니다.\n",
                    "example": "BTC"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "검증 결과 객체",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/TravelRuleVerificationResult"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "deposit_uuid": "9f432943-54e0-40b7-825f-b6fec8b42b79",
                      "verification_result": "verified",
                      "deposit_state": "PROCESSING"
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "에러 객체",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "validation_error (empty parameter)": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "deposit_uuid is empty"
                      }
                    }
                  },
                  "validation_error (invalid value)": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "\"field name\" does not have a valid value"
                      }
                    }
                  },
                  "validation_error (missing required parameter)": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "currency is missing, currency does not have a valid value, currency is empty, net_type is missing, net_type does not have a valid value, net_type is empty"
                      }
                    }
                  },
                  "already_have_deposit_transaction": {
                    "value": {
                      "error": {
                        "name": "already_have_deposit_transaction",
                        "message": "이미 입금 이력이 존재합니다."
                      }
                    }
                  }
                }
              }
            }
          },
          "404": {
            "description": "에러 객체",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "deposit_not_found": {
                    "value": {
                      "error": {
                        "name": "deposit_not_found",
                        "message": "입출금 정보를 찾지 못했습니다."
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

* [계정주 확인 서비스 지원 거래소 목록 조회](https://docs.upbit.com/kr/reference/list-travelrule-vasps.md)
* [입금 UUID로 계정주 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid.md)