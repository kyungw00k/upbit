---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 담보 추가

이용 중인 코인빌리기에 담보자산을 추가합니다.

이용 중인 코인빌리기 포지션에 담보를 추가합니다. 빌린 수량은 변경되지 않으며, 담보를 추가하여 현재 LTV(렌딩비율)를 낮추고 강제상환 기준 LTV(`liquidation_ltv`)와의 여유를 확보할 수 있습니다.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### **담보추가 요청이 접수되었다고 해서 담보추가가 완료된 것은 아닙니다.**

  요청이 접수되면 `202 Accepted`와 `request_uuid`가 반환되며, 실제 담보추가 처리는 이후에 진행됩니다. 담보추가 처리 결과는 [**코인빌리기 요청 상태 조회**](https://docs.upbit.com/kr/reference/get-lending-request) API로 확인해 주세요.
</Callout>

#### 담보추가 흐름

1. [**개별 코인빌리기 조회**](https://docs.upbit.com/kr/reference/get-lending)로 대상 포지션의 현재 LTV(`current_ltv`)와 설정된 담보자산(`collaterals[].currency`)을 확인합니다.
2. **코인빌리기 담보 추가**를 요청합니다. 요청이 접수되면 `request_uuid`가 반환되며, 요청 상태(`state`)는 `processing`으로 표시됩니다.
3. [**코인빌리기 요청 상태 조회**](https://docs.upbit.com/kr/reference/get-lending-request)로 최종 처리 결과(`done` 또는 `failed`)를 확인합니다.
4. 담보추가가 완료되면 [**개별 코인빌리기 조회**](https://docs.upbit.com/kr/reference/get-lending)에서 변경된 담보 수량(`collaterals[].amount`)과 현재 LTV를 확인합니다.

#### 담보추가 기준 및 제한 사항

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### **담보를 추가한 포지션은 이후 코인빌리기 추가 신청을 할 수 없습니다.**

  추가로 빌릴 계획이 있다면 담보추가 전에 확인해 주세요.
</Callout>

* 최초 신청 시 설정한 담보자산에만 담보를 추가할 수 있습니다. 설정되지 않은 자산을 추가하면 `not_related_collateral` 에러가 반환됩니다.
* 추가할 담보 수량만큼 해당 자산을 포켓에 보유하고 있어야 합니다. 잔고가 부족하면 `insufficient_main_balance` 에러가 반환됩니다.
* 담보를 추가해도 빌린 수량은 변경되지 않으며, LTV만 낮아집니다.
* `type`은 현재 `add`(담보추가)만 지원합니다.
* 청산이 진행 중인 포지션에는 담보를 추가할 수 없습니다.
* 추가한 담보는 전체 상환이 완료되기 전까지 출금할 수 없습니다.

#### 질권 설정 계약서

질권 설정 계약서는 업비트 웹·앱에서 확인할 수 있습니다. 계약서 확인에 필요한 신청번호(`lending_number`)는 [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)에서 확인해 주세요.

#### 수량 입력 규칙

담보로 추가할 수량은 숫자 형식의 문자열로 입력하며, 소수점 이하 최대 8자리까지 입력할 수 있습니다. 소수점 이하 자릿수가 8자리를 초과하면 에러가 반환됩니다.

#### 요청 식별자(identifier)

`identifier`는 클라이언트가 요청마다 직접 지정하는 식별자입니다. 요청 처리 결과를 추적하거나 중복 요청을 식별하는 데 사용하며, 코인빌리기 요청 상태 조회에서 `request_uuid` 대신 사용할 수 있습니다.

* 영문 대소문자, 숫자, 언더스코어(`_`), 점(`.`), 하이픈(`-`)을 조합한 1\~64자 문자열로 입력하며, 대소문자를 구분합니다.
* 동일한 포켓의 모든 코인빌리기 요청(신청, 추가신청, 담보추가, 상환, 담보상환)에서 중복되지 않아야 합니다.
* 요청이 접수되어 `request_uuid`가 발급되면 처리 결과와 관계없이 해당 `identifier`를 다시 사용할 수 없습니다. 요청이 실패(`failed`)한 경우에도 재시도 시 새로운 값을 지정해 주세요.
* 응답을 받지 못했거나 `duplicated_identifier` 에러를 받았다면, 새 요청을 보내기 전에 원래 `identifier`로 [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)를 먼저 호출해 주세요.
* `identifier`를 사용하지 않는 경우 요청 Body에서 해당 필드를 생략해 주세요. 이 경우 응답에서는 `"identifier": null`로 반환됩니다.

#### query\_hash 생성

`query_hash`를 계산할 때 `collaterals`처럼 배열 또는 객체 형태의 파라미터는 공백 없는 JSON 문자열로 직렬화하여 쿼리 문자열에 포함해야 합니다. 나머지 인증 과정은 [인증 가이드](https://docs.upbit.com/kr/reference/auth)를 참고해 주세요.

* 직렬화 예시: `collaterals=[{"currency":"KRW","amount":"500000"}]`
* Python의 `json.dumps`는 `separators=(",", ":")`를 지정해야 합니다. JavaScript의 `JSON.stringify`와 Java Gson의 `toJson`은 기본적으로 공백 없이 직렬화합니다. 자세한 구현 방법은 언어별 요청 예제를 참고해 주세요.

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
      "AdjustCollateralResponse": {
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
            "description": "담보추가 요청의 UUID"
          },
          "identifier": {
            "type": "string",
            "description": "요청 시 지정한 식별자",
            "nullable": true
          },
          "action": {
            "type": "string",
            "enum": [
              "add_collateral"
            ],
            "description": "요청 유형\n\n* `add_collateral`: 담보추가"
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
            "description": "담보를 추가할 포지션의 UUID"
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
            "description": "포지션 상태\n\n* `active`: 이용 중\n* `suspended`: 이용 제한\n* `liquidating`: 청산·담보상환 진행 중\n* `repaid`: 정상 상환 완료\n* `liquidated`: 청산 완료 (담보상환에 따른 종료 포함)"
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
    "/v1/lendings/adjust_collateral": {
      "post": {
        "operationId": "adjust-collateral",
        "summary": "담보 추가",
        "description": "이용 중인 코인빌리기에 담보자산을 추가합니다.",
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
                  "uuid",
                  "type",
                  "collaterals"
                ],
                "properties": {
                  "uuid": {
                    "type": "string",
                    "description": "담보를 추가할 포지션의 UUID"
                  },
                  "type": {
                    "type": "string",
                    "enum": [
                      "add"
                    ],
                    "description": "담보 조정 방향\n\n* `add`: 담보추가"
                  },
                  "collaterals": {
                    "type": "array",
                    "items": {
                      "type": "object",
                      "required": [
                        "currency",
                        "amount"
                      ],
                      "properties": {
                        "currency": {
                          "type": "string",
                          "description": "추가할 담보자산의 통화 코드"
                        },
                        "amount": {
                          "type": "string",
                          "description": "추가할 담보 수량"
                        }
                      }
                    },
                    "description": "추가할 담보 목록"
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
                    "uuid": "b03d6818-efc2-4dcb-b7cf-fca26f6ba81d",
                    "type": "add",
                    "collaterals": [
                      {
                        "currency": "KRW",
                        "amount": "500000"
                      }
                    ],
                    "identifier": "add_collateral_20260810_001"
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
                  "$ref": "#/components/schemas/AdjustCollateralResponse"
                },
                "examples": {
                  "Example": {
                    "description": "요청 접수 완료 (처리 결과는 코인빌리기 요청 상태 조회로 확인)",
                    "value": {
                      "request_uuid": "052d008d-254c-4578-8048-c402f6d79bcc",
                      "identifier": "add_collateral_20260810_001",
                      "action": "add_collateral",
                      "state": "processing",
                      "failure": null,
                      "created_at": "2026-08-10T00:05:12+09:00",
                      "completed_at": null,
                      "lending_uuid": "b03d6818-efc2-4dcb-b7cf-fca26f6ba81d",
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
                    "description": "* identifier가 빈 문자열\n* identifier 길이 초과 (65자 이상)\n* identifier에 허용되지 않는 문자 포함\n* type이 add 외 값\n* collaterals가 빈 배열\n* 동일한 담보자산을 중복 입력",
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
                  "insufficient_main_balance": {
                    "description": "추가할 담보 수량만큼 포켓 잔고가 없음",
                    "value": {
                      "error": {
                        "name": "insufficient_main_balance",
                        "message": "추가할 담보자산의 보유 수량이 부족해요."
                      }
                    }
                  },
                  "not_related_collateral": {
                    "description": "최초 신청 시 설정하지 않은 담보자산을 추가",
                    "value": {
                      "error": {
                        "name": "not_related_collateral",
                        "message": "설정된 담보자산만 입금할 수 있어요."
                      }
                    }
                  },
                  "insufficient_level": {
                    "description": "자격 요건(레벨) 미충족",
                    "value": {
                      "error": {
                        "name": "insufficient_level",
                        "message": "코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요."
                      }
                    }
                  },
                  "terminated_state": {
                    "description": "이미 상환이 완료된 포지션에 요청",
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
                  },
                  "processing_transaction_exists": {
                    "description": "처리 중인 요청 존재",
                    "value": {
                      "error": {
                        "name": "processing_transaction_exists",
                        "message": "이미 처리 중인 요청이 있어요.\n잠시 후 다시 시도해 주세요."
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
                    "description": "`uuid`에 해당하는 포지션이 없거나 해당 포켓의 포지션이 아님",
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
          "409": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "command_preparation_stale": {
                    "description": "처리 중 정보 변경",
                    "value": {
                      "error": {
                        "name": "command_preparation_stale",
                        "message": "처리 중 정보가 변경되었어요.\n잠시 후 다시 시도해 주세요."
                      }
                    }
                  },
                  "amount_limit_price_snapshot_changed": {
                    "description": "처리 중 가격 정보 변경",
                    "value": {
                      "error": {
                        "name": "amount_limit_price_snapshot_changed",
                        "message": "가격 정보가 변경되었어요.\n잠시 후 다시 시도해 주세요."
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
              "code": "# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\ncurl --request POST \\\n  --url 'https://api.upbit.com/v1/lendings/adjust_collateral' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json' \\\n  --header 'Content-Type: application/json' \\\n  --data '{\n    \"uuid\": \"b03d6818-efc2-4dcb-b7cf-fca26f6ba81d\",\n    \"type\": \"add\",\n    \"collaterals\": [\n      {\"currency\": \"KRW\", \"amount\": \"500000\"}\n    ],\n    \"identifier\": \"add_collateral_20260810_001\"\n  }'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport json\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/adjust_collateral\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nbody = {\n    \"uuid\": \"b03d6818-efc2-4dcb-b7cf-fca26f6ba81d\",\n    \"type\": \"add\",\n    \"collaterals\": [\n        {\"currency\": \"KRW\", \"amount\": \"500000\"},\n    ],\n    \"identifier\": \"add_collateral_20260810_001\",\n}\n\n\ndef to_hash_value(value):\n    # query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다(separators로 공백 제거).\n    if isinstance(value, (list, dict)):\n        return json.dumps(value, separators=(\",\", \":\"), ensure_ascii=False)\n    return value\n\n\nhash_params = {key: to_hash_value(value) for key, value in body.items()}\nquery_string = unquote(urlencode(hash_params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n    \"Content-Type\": \"application/json\",\n}\n# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=body)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/adjust_collateral\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst body = {\n  uuid: \"b03d6818-efc2-4dcb-b7cf-fca26f6ba81d\",\n  type: \"add\",\n  collaterals: [{ currency: \"KRW\", amount: \"500000\" }],\n  identifier: \"add_collateral_20260810_001\",\n};\n\n// query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다(JSON.stringify 기본 출력).\nconst toHashValue = (value) =>\n  typeof value === \"object\" && value !== null ? JSON.stringify(value) : value;\nconst queryString = Object.entries(body)\n  .map(([key, value]) => `${key}=${toHashValue(value)}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\n// 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\naxios\n  .post(`${BASE_URL}${PATH}`, body, {\n    headers: {\n      Authorization: `Bearer ${jwtToken}`,\n      Accept: \"application/json\",\n      \"Content-Type\": \"application/json\",\n    },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport com.google.gson.Gson;\nimport com.google.gson.GsonBuilder;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.List;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.MediaType;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\n\npublic class AdjustCollateral {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/adjust_collateral\";\n    // Gson은 공백 없이 직렬화합니다. 문자열 값이 그대로 유지되도록 HTML 이스케이프는 끕니다.\n    private static final Gson GSON = new GsonBuilder().disableHtmlEscaping().create();\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, Object> body = new LinkedHashMap<>();\n        body.put(\"uuid\", \"b03d6818-efc2-4dcb-b7cf-fca26f6ba81d\");\n        body.put(\"type\", \"add\");\n        body.put(\"collaterals\", List.of(collateral(\"KRW\", \"500000\")));\n        body.put(\"identifier\", \"add_collateral_20260810_001\");\n        String queryString = body.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + toHashValue(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .post(RequestBody.create(GSON.toJson(body), MediaType.parse(\"application/json; charset=utf-8\")))\n            .build();\n\n        // 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static Map<String, String> collateral(String currency, String amount) {\n        Map<String, String> item = new LinkedHashMap<>();\n        item.put(\"currency\", currency);\n        item.put(\"amount\", amount);\n        return item;\n    }\n\n    // query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다.\n    private static String toHashValue(Object value) {\n        return (value instanceof List || value instanceof Map) ? GSON.toJson(value) : String.valueOf(value);\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)