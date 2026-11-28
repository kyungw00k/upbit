---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 요청 상태 조회

코인빌리기 요청의 처리 상태와 결과를 조회합니다.

코인빌리기 POST 요청은 접수 결과만 바로 응답하고, 실제 처리는 그 뒤에 진행됩니다. 이 API로 요청의 처리 상태와 처리 결과를 확인합니다. 요청 응답의 `request_uuid` 또는 요청 시 지정한 `identifier`로 조회하며, 조회할 수 있는 요청은 다음과 같습니다.

* [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral)

<Callout icon="far fa-triangle-exclamation" theme="warn">
  **조회 요청 시&#x20;**`request_uuid`**&#x20;또는&#x20;**`identifier`**&#x20;중 하나만 입력해야 합니다.** 두 파라미터 모두 선택(Optional) 파라미터이지만, 조회할 요청을 지정하기 위해 하나는 반드시 입력해야 합니다. 두 파라미터를 동시에 입력하거나 모두 생략하면 `invalid_parameter` 에러가 반환됩니다.
</Callout>

#### 요청 처리 상태

| 상태           | 설명                                    | 다음 상태            |
| ------------ | ------------------------------------- | ---------------- |
| `processing` | 요청을 처리하고 있습니다. 요청이 접수되면 바로 이 상태가 됩니다. | `done`, `failed` |
| `done`       | 요청이 정상 처리되었습니다.                       | -                |
| `failed`     | 요청이 실패했습니다. 실패 사유는 `failure`에서 확인합니다. | -                |

#### 처리 결과 확인

* 요청이 정상 처리(`done`)되면 `lending_uuid`와 `transaction_uuid`로 [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)와 [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions)에서 결과를 확인할 수 있습니다.
* 신청(`borrow`)은 포지션이 생성된 뒤(`done`)에만 `lending_uuid`와 `lending_state`가 채워집니다. 그 외 요청은 대상 포지션의 `lending_uuid`와 `lending_state`를 항상 반환합니다.
* `lending_state`는 요청 시점이 아닌 조회 시점의 포지션 상태입니다.
* 처리에 실패한 요청도 이력이 유지되며 조회할 수 있습니다. 실패 사유는 `failure`에서 확인합니다. 접수 단계의 에러와 처리 단계의 실패를 구분하는 방법은 [코인빌리기 개요](https://docs.upbit.com/kr/reference/lending-overview)의 **요청 실패 처리**를 참고해 주세요.

#### 응답을 받지 못했거나 `duplicated_identifier`를 받은 경우

새 요청을 보내기 전에 원래 `identifier`로 요청 상태를 먼저 조회해 주세요. 예: `GET /v1/lendings/request?identifier=borrow_btc_001`

| `state`      | 의미                | 조치                                                                 |
| ------------ | ----------------- | ------------------------------------------------------------------ |
| `processing` | 요청이 접수되어 처리 중입니다. | 새 요청을 보내지 말고 같은 `identifier`로 상태를 다시 확인해 주세요.                      |
| `done`       | 요청이 정상 처리되었습니다.   | 기존 결과(`lending_uuid`, `transaction_uuid`)를 사용해 주세요.                |
| `failed`     | 요청이 실패했습니다.       | `failure.name`으로 실패 원인을 확인하고, 별도 요청이 필요하면 새 `identifier`를 사용해 주세요. |

#### 조회 주기

* 요청 상태가 `done` 또는 `failed`가 될 때까지 조회합니다.
* 이 API는 \[Exchange 기본 그룹]에서 다른 조회 API와 요청 수를 함께 차감합니다.**&#x20;짧은 간격으로 반복 조회하지 않도록 조회 주기를 조절해 주세요.**

#### 유의사항

* 코인빌리기 신청 테스트 응답의 `request_uuid`는 조회할 수 없습니다.
* 해당 포켓의 요청만 조회할 수 있습니다. 존재하지 않거나 다른 포켓의 요청을 조회하면 `request_not_found` 에러가 반환됩니다.

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
      "LendingRequest": {
        "type": "object",
        "required": [
          "request_uuid",
          "identifier",
          "action",
          "state",
          "failure",
          "created_at",
          "completed_at",
          "lending_uuid",
          "lending_state",
          "transaction_uuid"
        ],
        "properties": {
          "request_uuid": {
            "type": "string",
            "description": "요청의 UUID"
          },
          "identifier": {
            "type": "string",
            "description": "요청 시 지정한 식별자",
            "nullable": true
          },
          "action": {
            "type": "string",
            "enum": [
              "borrow",
              "borrow_more",
              "add_collateral",
              "repay",
              "repay_with_collateral"
            ],
            "description": "요청 유형\n\n* `borrow`: 코인빌리기 신청\n* `borrow_more`: 코인빌리기 추가 신청\n* `add_collateral`: 담보추가\n* `repay`: 빌린 자산으로 상환\n* `repay_with_collateral`: 담보상환"
          },
          "state": {
            "type": "string",
            "enum": [
              "processing",
              "failed",
              "done"
            ],
            "description": "요청 처리 상태\n\n* `processing`: 처리 중\n* `failed`: 실패\n* `done`: 정상 처리 완료"
          },
          "failure": {
            "type": "object",
            "description": "요청이 정상적으로 완료되지 않은 사유",
            "nullable": true,
            "properties": {
              "name": {
                "type": "string",
                "description": "실패 사유 코드. 요청 유형(action)별로 반환될 수 있는 값과 조치 방법은 코인빌리기 개요의 '요청 실패 처리'를 참고해 주세요."
              },
              "message": {
                "type": "string",
                "description": "실패 사유 설명"
              }
            }
          },
          "created_at": {
            "type": "string",
            "description": "요청 접수 시각 (ISO 8601)"
          },
          "completed_at": {
            "type": "string",
            "description": "요청이 성공 또는 실패로 종료된 시각 (ISO 8601)",
            "nullable": true
          },
          "lending_uuid": {
            "type": "string",
            "description": "요청과 연결된 포지션의 UUID",
            "nullable": true
          },
          "lending_state": {
            "type": "string",
            "enum": [
              "active",
              "suspended",
              "liquidating",
              "repaid",
              "liquidated",
              null
            ],
            "description": "포지션 상태\n\n* `active`: 이용 중\n* `suspended`: 이용 제한\n* `liquidating`: 청산·담보상환 진행 중\n* `repaid`: 정상 상환 완료\n* `liquidated`: 청산 완료 (담보상환에 따른 종료 포함)",
            "nullable": true
          },
          "transaction_uuid": {
            "type": "string",
            "description": "요청으로 생성된 거래의 UUID",
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
    "/v1/lendings/request": {
      "get": {
        "operationId": "get-lending-request",
        "summary": "코인빌리기 요청 상태 조회",
        "description": "코인빌리기 요청의 처리 상태와 결과를 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "parameters": [
          {
            "name": "request_uuid",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회할 요청의 UUID"
            },
            "examples": {
              "default": {
                "value": "80e7f7c8-d383-42c6-9206-64674a6e07f4"
              }
            }
          },
          {
            "name": "identifier",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회할 요청의 클라이언트 지정 식별자"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/LendingRequest"
                },
                "examples": {
                  "Processing borrow request": {
                    "description": "요청이 접수되어 처리 중인 경우",
                    "value": {
                      "request_uuid": "3f6b2d1e-8c4a-4e7b-9a52-1d0c7e9b4f21",
                      "identifier": "borrow_btc_001",
                      "action": "borrow",
                      "state": "processing",
                      "failure": null,
                      "created_at": "2026-09-06T20:58:02+09:00",
                      "completed_at": null,
                      "lending_uuid": null,
                      "lending_state": null,
                      "transaction_uuid": null
                    }
                  },
                  "Failed borrow request": {
                    "description": "요청은 접수되었으나 신청이 실패한 경우",
                    "value": {
                      "request_uuid": "80e7f7c8-d383-42c6-9206-64674a6e07f4",
                      "identifier": "borrow-usdt-20260906-001",
                      "action": "borrow",
                      "state": "failed",
                      "failure": {
                        "name": "insufficient_level",
                        "message": "코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요."
                      },
                      "created_at": "2026-09-06T21:03:14+09:00",
                      "completed_at": "2026-09-06T21:08:14+09:00",
                      "lending_uuid": null,
                      "lending_state": null,
                      "transaction_uuid": null
                    }
                  },
                  "Completed full repayment": {
                    "description": "빌린 자산으로 전액 상환이 완료된 경우",
                    "value": {
                      "request_uuid": "bcfbec63-e47c-4469-9ca6-9143f87dc8eb",
                      "identifier": "repay-usdt-20260906-003",
                      "action": "repay",
                      "state": "done",
                      "failure": null,
                      "created_at": "2026-09-06T22:11:37+09:00",
                      "completed_at": "2026-09-06T22:11:39+09:00",
                      "lending_uuid": "5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b",
                      "lending_state": "repaid",
                      "transaction_uuid": "89ebc7b6-fcbd-438d-a4a0-8a5b22c09e26"
                    }
                  },
                  "Processing repayment with collateral": {
                    "description": "담보상환이 진행 중인 경우",
                    "value": {
                      "request_uuid": "c1a93e63-862c-4236-94af-a72cc63becbd",
                      "identifier": "repay_with_collateral_20260809_001",
                      "action": "repay_with_collateral",
                      "state": "processing",
                      "failure": null,
                      "created_at": "2026-08-09T20:00:00+09:00",
                      "completed_at": null,
                      "lending_uuid": "98e408b9-2d78-4a2a-b738-42ae07610410",
                      "lending_state": "liquidating",
                      "transaction_uuid": null
                    }
                  },
                  "Completed repayment with collateral": {
                    "description": "담보상환이 완료된 경우",
                    "value": {
                      "request_uuid": "c1a93e63-862c-4236-94af-a72cc63becbd",
                      "identifier": "repay_with_collateral_20260809_001",
                      "action": "repay_with_collateral",
                      "state": "done",
                      "failure": null,
                      "created_at": "2026-08-09T20:00:00+09:00",
                      "completed_at": "2026-08-09T20:00:05+09:00",
                      "lending_uuid": "98e408b9-2d78-4a2a-b738-42ae07610410",
                      "lending_state": "liquidated",
                      "transaction_uuid": "7d2e9a14-5b3c-4f8e-a1d6-0c9b8e7f6a53"
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
                    "description": "* request_uuid와 identifier를 모두 생략\n* request_uuid와 identifier를 동시에 입력",
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
          },
          "404": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "request_not_found": {
                    "description": "* request_uuid에 해당하는 요청이 없거나 해당 포켓의 코인빌리기 요청이 아님\n* identifier에 해당하는 요청이 없거나 해당 포켓의 코인빌리기 요청이 아님",
                    "value": {
                      "error": {
                        "name": "request_not_found",
                        "message": "빌리기 관련 요청 건을 찾지 못했습니다."
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
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/lendings/request?request_uuid=80e7f7c8-d383-42c6-9206-64674a6e07f4' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/request\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\"request_uuid\": \"80e7f7c8-d383-42c6-9206-64674a6e07f4\"}\nquery_string = unquote(urlencode(params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"}\nres = requests.get(f\"{BASE_URL}{PATH}?{query_string}\", headers=headers)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/request\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = { request_uuid: \"80e7f7c8-d383-42c6-9206-64674a6e07f4\" };\nconst queryString = Object.entries(params)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .get(`${BASE_URL}${PATH}?${queryString}`, {\n    headers: { Authorization: `Bearer ${jwtToken}`, Accept: \"application/json\" },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class GetLendingRequest {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/request\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new LinkedHashMap<>();\n        params.put(\"request_uuid\", \"80e7f7c8-d383-42c6-9206-64674a6e07f4\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .get()\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)