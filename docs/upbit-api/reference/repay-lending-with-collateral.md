---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 담보로 상환

맡긴 담보로 코인빌리기를 전액 상환합니다.

맡긴 담보를 활용해 빌린 자산을 전액 상환합니다. 빌린 자산을 별도로 준비할 필요가 없으며, 담보자산과 빌린 자산의 종류에 따라 담보자산을 매도하거나 상환에 필요한 자산을 매수할 수 있습니다. 상환 및 수수료 정산 후 남은 담보는 반환됩니다.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### 담보상환 요청이 접수되었다고 해서 상환이 완료된 것은 아닙니다.

  요청이 접수되면 `202 Accepted`와 `request_uuid`가 반환되며, 실제 상환 처리는 이후에 진행됩니다. 상환 처리 결과는 [**코인빌리기 요청 상태 조회**](https://docs.upbit.com/kr/reference/get-lending-request) API로 확인해 주세요.
</Callout>

#### 처리 과정

1. **담보상환**을 요청합니다. 요청이 접수되면 `request_uuid`가 반환되며, 요청 상태(`state`)는 `processing`입니다. 접수 응답의 포지션 상태(`lending_state`)는 `active`일 수 있으며, 담보상환 처리가 진행되면 `liquidating`으로 바뀔 수 있습니다.
2. 남은 빌린 수량과 수수료 정산에 필요한 자산을 확보하기 위해 담보를 사용하거나 매도하고, 필요한 경우 빌린 자산을 매수합니다.
3. 빌린 자산의 상환과 수수료 정산을 완료한 후 남은 담보를 반환합니다.
4. [**코인빌리기 요청 상태 조회**](https://docs.upbit.com/kr/reference/get-lending-request)로 최종 처리 결과(`done` 또는 `failed`)를 확인합니다. 상환이 완료되면 포지션 상태는 `liquidated`, 종료 유형(`closure_type`)은 `repay_with_collateral`이 됩니다.

#### 상환 범위와 수수료

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### **담보상환은 전액 상환만 지원합니다.**

  일부 수량만 상환하려면 [**코인빌리기 상환**](https://docs.upbit.com/kr/reference/repay-lending) API를 이용해 주세요. 담보상환할 경우 이용수수료와 담보상환수수료가 부과되며, 담보자산 매도 또는 빌린 자산 매수가 발생하면 해당 거래에 대한 거래수수료도 부과됩니다.
</Callout>

* 담보상환수수료는 담보상환 처리 직전 남은 빌린 수량에 해당 포지션의 담보상환수수료율(`collateral_repayment_fee_rate`)을 곱하여 계산합니다. 적용 수수료율은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">**개별 코인빌리기 조회**</Anchor> API에서 확인할 수 있습니다.
  * 예: 담보상환수수료율이 1%인 포지션에서 100 USDT를 빌리고 10 USDT를 일부 상환한 뒤 담보상환하면, 남은 90 USDT에 1%를 적용한 0.9 USDT가 담보상환수수료로 부과됩니다.
* 이용수수료는 일별 최대 이용 수량에 해당 포지션의 일 단위 이용수수료율(`borrowing_fee_rate`)을 적용하여 계산합니다. 부분 상환 시에는 정산하지 않으며, 최종 상환 시 정산합니다.
* 수수료는 반대매매로 확보한 빌린 자산에서 차감하며, 부족한 경우 담보에서 차감합니다.
* 상환 및 수수료 정산 후 남은 담보는 반환됩니다.

#### 상환 결과 확인

* 반대매매 체결 결과는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-lending-transactions">**코인빌리기 거래 내역 조회**</Anchor>에서 `include_entries=true`로 요청하여 확인할 수 있습니다. 담보 매도(`collateral_sale`)와 빌린 자산 매수(`borrowed_asset_purchase`) 내역의 `execution`에는 체결 수량, 체결 총액, 평균 체결 가격 및 거래 수수료가 포함됩니다.
* 최종 확정된 수수료는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">**개별 코인빌리기 조회**</Anchor>의 `final_fee`에서 확인할 수 있습니다.
* 담보상환과 강제상환은 모두 포지션 상태가 `liquidated`로 종료됩니다. 종료 사유는<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending"> **개별 코인빌리기 조회**</Anchor>의 `closure_type`으로 구분해 주세요.

#### 상환 전 확인 사항

* 청산이 진행 중인 포지션은 담보상환할 수 없습니다.
* 빌린 자산을 보유하고 있다면 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/repay-lending">**코인빌리기 상환**</Anchor> API를 이용해 직접 상환할 수 있습니다. 이 경우 담보상환수수료는 부과되지 않습니다.

#### 요청 식별자(identifier)

`identifier`는 클라이언트가 요청마다 직접 지정하는 식별자입니다. 요청 처리 결과를 추적하거나 중복 요청을 식별하는 데 사용하며, 코인빌리기 요청 상태 조회에서 `request_uuid` 대신 사용할 수 있습니다.

* 영문 대소문자, 숫자, 언더스코어(`_`), 점(`.`), 하이픈(`-`)을 조합한 1\~64자 문자열로 입력하며, 대소문자를 구분합니다.
* 동일한 포켓의 모든 코인빌리기 요청(신청, 추가신청, 담보추가, 상환, 담보상환)에서 중복되지 않아야 합니다.
* 요청이 접수되어 `request_uuid`가 발급되면 처리 결과와 관계없이 해당 `identifier`를 다시 사용할 수 없습니다. 요청이 실패(`failed`)한 경우에도 재시도 시 새로운 값을 지정해 주세요.
* 응답을 받지 못했거나 `duplicated_identifier` 에러를 받았다면, 새 요청을 보내기 전에 원래 `identifier`로 [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)를 먼저 호출해 주세요.
* `identifier`를 사용하지 않는 경우 요청 Body에서 해당 필드를 생략해 주세요. 이 경우 응답에서는 `"identifier": null`로 반환됩니다.

***

#### 에러 안내

코인빌리기 API 이용 중 에러가 발생하면 반환된 에러 코드와 발생 사유를 확인하고, 해당 조건에 맞게 조치한 뒤 다시 요청해 주세요.

발생 가능한 에러와 상세 응답 명세는 [코인빌리기 API 안내](https://docs.upbit.com/kr/reference/lending-overview)의 **에러 안내** 항목을 참고해 주세요.

#### Rate Limit

초당 최대 10회 호출할 수 있습니다.&#x20;

포켓 단위로 측정되며 \[`lending` 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

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
      "RepayLendingWithCollateralResponse": {
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
            "description": "담보상환 요청의 UUID"
          },
          "identifier": {
            "type": "string",
            "description": "요청 시 지정한 식별자",
            "nullable": true
          },
          "action": {
            "type": "string",
            "enum": [
              "repay_with_collateral"
            ],
            "description": "요청 유형\n\n* `repay_with_collateral`: 담보상환"
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
            "description": "요청 처리가 종료된 시각 (ISO 8601)",
            "nullable": true
          },
          "lending_uuid": {
            "type": "string",
            "description": "상환할 포지션의 UUID"
          },
          "lending_state": {
            "type": "string",
            "enum": [
              "active",
              "suspended",
              "liquidating",
              "repaid",
              "liquidated"
            ],
            "description": "포지션 상태\n\n* `active`: 이용 중\n* `suspended`: 이용 제한\n* `liquidating`: 청산 또는 담보상환 진행 중\n* `repaid`: 정상 상환 완료\n* `liquidated`: 청산 또는 담보상환 완료"
          },
          "transaction_uuid": {
            "type": "string",
            "description": "해당 요청으로 생성된 거래의 UUID",
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
    "/v1/lendings/repay_with_collateral": {
      "post": {
        "operationId": "repay-lending-with-collateral",
        "summary": "담보로 상환",
        "description": "맡긴 담보로 코인빌리기를 전액 상환합니다.",
        "tags": [
          "Exchange"
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": [
                  "uuid"
                ],
                "properties": {
                  "uuid": {
                    "type": "string",
                    "description": "상환할 포지션의 UUID"
                  },
                  "identifier": {
                    "type": "string",
                    "description": "클라이언트가 지정하는 요청 식별자",
                    "pattern": "^[A-Za-z0-9_.-]{1,64}$"
                  }
                }
              },
              "examples": {
                "default": {
                  "value": {
                    "uuid": "98e408b9-2d78-4a2a-b738-42ae07610410",
                    "identifier": "repay_with_collateral_20260809_001"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "202": {
            "description": "요청 접수 완료. 처리 결과는 코인빌리기 요청 상태 조회에서 확인합니다.",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/RepayLendingWithCollateralResponse"
                },
                "examples": {
                  "Example": {
                    "description": "요청 접수 완료 (처리 결과는 코인빌리기 요청 상태 조회로 확인)",
                    "value": {
                      "request_uuid": "c1a93e63-862c-4236-94af-a72cc63becbd",
                      "identifier": "repay_with_collateral_20260809_001",
                      "action": "repay_with_collateral",
                      "state": "processing",
                      "failure": null,
                      "created_at": "2026-08-09T20:00:00+09:00",
                      "completed_at": null,
                      "lending_uuid": "98e408b9-2d78-4a2a-b738-42ae07610410",
                      "lending_state": "active",
                      "transaction_uuid": null
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
                    "description": "* uuid 누락 또는 UUID 형식 오류\n* identifier 빈 문자열\n* identifier 길이 초과 (65자 이상)\n* identifier 허용 외 문자",
                    "value": {
                      "error": {
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
                      }
                    }
                  },
                  "duplicated_identifier": {
                    "description": "identifier 재사용",
                    "value": {
                      "error": {
                        "name": "duplicated_identifier",
                        "message": "이미 등록된 identifier입니다."
                      }
                    }
                  },
                  "terminated_state": {
                    "description": "이미 종료된 포지션",
                    "value": {
                      "error": {
                        "name": "terminated_state",
                        "message": "이미 상환이 완료됐어요."
                      }
                    }
                  },
                  "invalid_state": {
                    "description": "요청할 수 없는 포지션 상태",
                    "value": {
                      "error": {
                        "name": "invalid_state",
                        "message": "일시적인 오류가 발생했어요.\n잠시 후 다시 시도해 주세요."
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
                  },
                  "invalid_kyc": {
                    "description": "고객 확인(KYC) 미완료",
                    "value": {
                      "error": {
                        "name": "invalid_kyc",
                        "message": "고객 확인(KYC)을 완료한 후 이용할 수 있어요.\n고객 확인 절차를 완료해 주세요."
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
                  "lending_not_found": {
                    "description": "존재하지 않거나 해당 포켓의 빌리기 건이 아님",
                    "value": {
                      "error": {
                        "name": "lending_not_found",
                        "message": "빌리기 건을 찾지 못했습니다."
                      }
                    }
                  }
                }
              }
            }
          },
          "503": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "service_unavailable": {
                    "description": "코인빌리기 서비스 일시 이용 불가",
                    "value": {
                      "error": {
                        "name": "service_unavailable",
                        "message": "코인빌리기 서비스를 이용할 수 없어요.\n잠시 후 다시 시도해 주세요."
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
              "code": "# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\ncurl --request POST \\\n  --url 'https://api.upbit.com/v1/lendings/repay_with_collateral' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json' \\\n  --header 'Content-Type: application/json' \\\n  --data '{\n    \"uuid\": \"98e408b9-2d78-4a2a-b738-42ae07610410\",\n    \"identifier\": \"repay_with_collateral_20260809_001\"\n  }'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/repay_with_collateral\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nbody = {\n    \"uuid\": \"98e408b9-2d78-4a2a-b738-42ae07610410\",\n    \"identifier\": \"repay_with_collateral_20260809_001\",\n}\nquery_string = unquote(urlencode(body, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n    \"Content-Type\": \"application/json\",\n}\n# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=body)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/repay_with_collateral\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst body = {\n  uuid: \"98e408b9-2d78-4a2a-b738-42ae07610410\",\n  identifier: \"repay_with_collateral_20260809_001\",\n};\nconst queryString = Object.entries(body)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\n// 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\naxios\n  .post(`${BASE_URL}${PATH}`, body, {\n    headers: {\n      Authorization: `Bearer ${jwtToken}`,\n      Accept: \"application/json\",\n      \"Content-Type\": \"application/json\",\n    },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport com.google.gson.Gson;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.MediaType;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\n\npublic class RepayLendingWithCollateral {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/repay_with_collateral\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> body = new LinkedHashMap<>();\n        body.put(\"uuid\", \"98e408b9-2d78-4a2a-b738-42ae07610410\");\n        body.put(\"identifier\", \"repay_with_collateral_20260809_001\");\n        String queryString = body.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .post(RequestBody.create(new Gson().toJson(body), MediaType.parse(\"application/json; charset=utf-8\")))\n            .build();\n\n        // 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)