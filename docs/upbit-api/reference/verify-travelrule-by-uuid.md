---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 입금 UUID로 계정주 검증 요청

입금의 UUID로 트래블룰 검증(입금 계정주 확인)을 요청합니다.

[block:html]
{
  "html": "<div class=\"callout-section\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i>  계정주 검증이 완료된 이후 입금 반영이 진행됩니다.\n      </div>\n    계정주 검증이 완료되지 않는 경우 업비트 신원 확인 과정에서 입력한 이름, 생년월일 정보와 상대 거래소에 입력된 정보가 동일한지 여부를 확인해주세요. <br>\n    검증이 완료 시 즉시 입금 반영이 진행되나, 블록체인 네트워크 지연등으로 검증 이후에도 입금 자체가 아직 처리중인 경우 입금이 완료된 이후 입금 반영이 진행됩니다.\n  </div>"
}
[/block]

[block:html]
{
  "html": "<div class=\"callout-section callout-section--danger\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i> POST API에 대한 Form 방식 요청은 2022년 3월 1일부로 지원이 종료되었습니다.\n      </div>\n    Form 방식 지원 종료에 따라 Urlencoded Form 방식으로 전송하는 POST 요청에 대한 정상적인 동작을 보장하지 않습니다. <b>반드시 JSON 형식으로 요청 본문(Body)을 전송</b>해주시기 바랍니다.\n  </div>"
}
[/block]

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
            <td>2024-04-24</td>
            <td><a href="https://docs.upbit.com/kr/changelog/travelrule_verification"> '입금 UUID로 트래블룰 검증 요청' 기능 신규 지원</a></td>
        </tr>
    </tbody>
</table>


    </div>
</div>

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 \[Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. 단, 동일 입금건에 대해서는 <span style="color: #C92532;">10분당 최대 1회</span> 호출할 수 있습니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

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
      "java",
      "node"
    ]
  },
  "paths": {
    "/v1/travel_rule/deposit/uuid": {
      "post": {
        "operationId": "verify-travelrule-by-uuid",
        "summary": "입금 UUID로 계정주 검증 요청",
        "description": "입금의 UUID로 트래블룰 검증(입금 계정주 확인)을 요청합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request POST \\\n  --url 'https://api.upbit.com/v1/travel_rule/deposit/uuid' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Content-Type: application/json' \\\n  --data '\n{\n\"deposit_uuid\": \"5b871d34-fe38-4025-8f5c-9b22028f85d3\",\n\"vasp_uuid\": \"8d4fe968-82b2-42e5-822f-3840a245f802\"\n}\n'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/travel_rule/deposit/uuid\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"deposit_uuid\": \"5b871d34-fe38-4025-8f5c-9b22028f85d3\",\n  \"vasp_uuid\": \"8d4fe968-82b2-42e5-822f-3840a245f802\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/travel_rule/deposit/uuid\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  deposit_uuid: \"5b871d34-fe38-4025-8f5c-9b22028f85d3\",\n  vasp_uuid: \"8d4fe968-82b2-42e5-822f-3840a245f802\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"POST\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n  data: params,\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\nimport com.google.gson.Gson;\n\npublic class CreateDepositAddress {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/travel_rule/deposit/uuid\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"deposit_uuid\", \"5b871d34-fe38-4025-8f5c-9b22028f85d3\");\n        params.put(\"vasp_uuid\", \"8d4fe968-82b2-42e5-822f-3840a245f802\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        String jsonBody = new Gson().toJson(params);\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .post(RequestBody.create(jsonBody, okhttp3.MediaType.parse(\"application/json; charset=utf-8\")))\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "requestBody": {
          "description": "UUID로 입금 트랜잭션 트래블룰 검증",
          "content": {
            "application/json": {
              "examples": {
                "default": {
                  "value": {
                    "deposit_uuid": "5b871d34-fe38-4025-8f5c-9b22028f85d3",
                    "vasp_uuid": "8d4fe968-82b2-42e5-822f-3840a245f802"
                  }
                }
              },
              "schema": {
                "type": "object",
                "required": [
                  "deposit_uuid",
                  "vasp_uuid"
                ],
                "properties": {
                  "deposit_uuid": {
                    "type": "string",
                    "description": "입금의 유일식별자(UUID)",
                    "example": "5b871d34-fe38-4025-8f5c-9b22028f85d3"
                  },
                  "vasp_uuid": {
                    "type": "string",
                    "description": "자산을 출금한 상대 거래소의 유일식별자(UUID)",
                    "example": "8d4fe968-82b2-42e5-822f-3840a245f802"
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
                  "empty parameter error": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "deposit_uuid is empty"
                      }
                    }
                  },
                  "invalid value error": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "\"field name\" does not have a valid value"
                      }
                    }
                  },
                  "missing required parameter error": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "currency is missing, currency does not have a valid value, currency is empty, net_type is missing, net_type does not have a valid value, net_type is empty"
                      }
                    }
                  },
                  "already have deposit transaction error": {
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
                  "not found deposit error": {
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

# Sibling pages

* [계정주 확인 서비스 지원 거래소 목록 조회](https://docs.upbit.com/kr/reference/list-travelrule-vasps.md)
* [입금 TxID로 계정주 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-txid.md)