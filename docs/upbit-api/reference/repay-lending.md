---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 상환

빌린 자산의 전부 또는 일부를 상환합니다.

빌린 자산으로 코인빌리기를 상환합니다. 상환 수량을 지정하여 일부 또는 전액을 상환할 수 있으며, 이용기간 중 언제든지 상환할 수 있습니다.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### **상환 요청이 접수되었다고 해서 상환이 완료된 것은 아닙니다.**

  요청이 접수되면 `202 Accepted`와 `request_uuid`가 반환되며, 실제 상환 처리는 이후에 진행됩니다. 상환 처리 결과는 [**코인빌리기 요청 상태 조회**](https://docs.upbit.com/kr/reference/get-lending-request) API로 확인해 주세요.
</Callout>

#### 상환 방식

`repayment_mode`와 `repayment_amount`를 이용하여 상환 방식을 선택할 수 있습니다.

| `repayment_mode` | 상환 방식       | 요청 방법                                                                            | 처리                                                 |
| ---------------- | ----------- | -------------------------------------------------------------------------------- | -------------------------------------------------- |
| `amount` (기본값)   | 수량 지정 상환    | `repayment_amount`에 상환할 수량을 입력합니다. `repayment_mode`는 생략하거나 `amount`로 지정할 수 있습니다. | 지정한 수량을 상환합니다. 처리 시점의 남은 빌린 수량과 같으면 전체 상환으로 처리됩니다. |
| `full`           | 남은 수량 전체 상환 | `repayment_mode`를 `full`로 지정하고, `repayment_amount`는 입력하지 않습니다.                   | 처리 시점의 남은 빌린 수량 전체를 상환합니다.                         |

* **전체 상환 시에는&#x20;**`repayment_mode`**를&#x20;**`full`**로 지정하는 것을 권장합니다.** 남은 빌린 수량을 별도로 입력하지 않고 전체 상환할 수 있습니다.
* 상환 수량은 남은 빌린 수량을 초과할 수 없습니다. 남은 빌린 수량은 [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)의 `borrowed_amount`에서 확인해 주세요.

#### 부분 상환과 전체 상환

| 구분    | 부분 상환                    | 전체 상환             |
| ----- | ------------------------ | ----------------- |
| 이용수수료 | 정산하지 않음                  | 정산                |
| 담보    | 반환하지 않으며, LTV(렌딩비율)가 낮아짐 | 수수료 정산 후 남은 담보 반환 |
| 포지션   | 종료되지 않음                  | 종료(`repaid`)      |

#### 수수료 정산

* 이용수수료는 이용 중에 차감되지 않으며, 전체 상환 시 맡긴 담보에서 정산됩니다. 따라서 빌린 수량을 모두 상환하더라도 수수료 정산 후 남은 담보만 반환될 수 있습니다.
* 최종 확정된 수수료는 [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending)의 `final_fee`에서 확인할 수 있습니다. 수수료 차감 및 담보 반환 내역은 [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions)의 `entries`에서 확인해 주세요.

#### 상환 전 확인 사항

* 상환할 수량만큼 빌린 자산을 포켓에 보유하고 있어야 합니다. 잔고가 부족하면 `insufficient_main_balance` 에러가 반환됩니다.
* 청산이 진행 중인 포지션은 상환할 수 없습니다.
* 빌린 자산이 부족한 경우 [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral) API를 이용해 맡긴 담보로 전액 상환할 수 있습니다. 이 경우 담보상환수수료가 추가로 부과됩니다.

#### 파라미터 조합별 처리

`repayment_amount`와 `repayment_mode`의 입력 조합에 따라 다음과 같이 처리됩니다.

| `repayment_amount` | `repayment_mode` | 처리 결과                       |
| ------------------ | ---------------- | --------------------------- |
| 입력                 | 미입력              | 지정한 수량 상환 (기본값 `amount` 적용) |
| 입력                 | `amount`         | 지정한 수량 상환                   |
| 입력                 | `full`           | 요청 실패 (`invalid_parameter`) |
| 미입력                | 미입력              | 요청 실패 (`invalid_parameter`) |
| 미입력                | `amount`         | 요청 실패 (`invalid_parameter`) |
| 미입력                | `full`           | 처리 시점의 남은 빌린 수량 전체 상환       |
| 입력 또는 미입력          | 허용되지 않는 값 또는 타입  | 요청 실패 (`invalid_parameter`) |

`repayment_amount`를 지정한 경우, 처리 시점의 남은 빌린 수량에 따라 다음과 같이 처리됩니다.

| 입력한 `repayment_amount` | 처리 결과                                 |
| ---------------------- | ------------------------------------- |
| 남은 빌린 수량보다 작음          | 부분 상환                                 |
| 남은 빌린 수량과 같음           | 전체 상환                                 |
| 남은 빌린 수량보다 큼           | 요청 실패 (`exceed_lending_repay_amount`) |
| 0 이하                   | 요청 실패 (`invalid_parameter`)           |

#### 요청 Body 예시

1. **수량 지정 상환 (**`repayment_mode`**&#x20;생략)**

```json
{
  "uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "repayment_amount": "0.00200000",
  "identifier": "repay_amount_001"
}
```

1. **수량 지정 상환 (**`repayment_mode`**를&#x20;**`amount`**로 지정)**

첫 번째 예시와 동일하게 처리됩니다.

```json
{
  "uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "repayment_mode": "amount",
  "repayment_amount": "0.00200000",
  "identifier": "repay_amount_002"
}
```

1. **남은 수량 전체 상환 (**`repayment_mode`**를&#x20;**`full`**로 지정)**

`repayment_amount`는 입력하지 않습니다.

```json
{
  "uuid": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "repayment_mode": "full",
  "identifier": "repay_full_003"
}
```

#### 수량 입력 규칙

상환 수량은 숫자 형식의 문자열로 입력하며, 소수점 이하 최대 8자리까지 입력할 수 있습니다. 소수점 이하 자릿수가 8자리를 초과하면 에러가 반환됩니다.

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
      "RepayLendingResponse": {
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
            "description": "상환 요청의 UUID"
          },
          "identifier": {
            "type": "string",
            "description": "요청 시 지정한 식별자",
            "nullable": true
          },
          "action": {
            "type": "string",
            "enum": [
              "repay"
            ],
            "description": "요청 유형\n\n* `repay`: 빌린 자산으로 상환"
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
    "/v1/lendings/repay": {
      "post": {
        "operationId": "repay-lending",
        "summary": "코인빌리기 상환",
        "description": "빌린 자산의 전부 또는 일부를 상환합니다.",
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
                  "repayment_amount": {
                    "type": "string",
                    "description": "상환 수량"
                  },
                  "repayment_mode": {
                    "type": "string",
                    "enum": [
                      "amount",
                      "full"
                    ],
                    "description": "상환 수량을 정하는 방식\n\n* `amount` (기본값): 지정한 수량 상환. 지정한 수량이 처리 시점의 남은 빌린 수량과 같으면 전체 상환으로 처리합니다.\n* `full`: 처리 시점의 남은 빌린 수량 전체 상환. `repayment_amount`를 함께 입력할 수 없습니다.",
                    "default": "amount"
                  },
                  "identifier": {
                    "type": "string",
                    "description": "클라이언트가 지정하는 요청 식별자",
                    "pattern": "^[A-Za-z0-9_.-]{1,64}$"
                  }
                }
              },
              "examples": {
                "수량 지정 상환": {
                  "value": {
                    "uuid": "34e0c499-2bce-43f7-8eac-26852e79b78f",
                    "repayment_amount": "0.00200000",
                    "identifier": "repay_partial_20260810_001"
                  }
                },
                "수량 지정 상환 (repayment_mode: amount)": {
                  "value": {
                    "uuid": "34e0c499-2bce-43f7-8eac-26852e79b78f",
                    "repayment_mode": "amount",
                    "repayment_amount": "0.00200000",
                    "identifier": "repay_partial_20260810_002"
                  }
                },
                "남은 수량 전체 상환 (repayment_mode: full)": {
                  "value": {
                    "uuid": "34e0c499-2bce-43f7-8eac-26852e79b78f",
                    "repayment_mode": "full",
                    "identifier": "repay_full_20260810_003"
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
                  "$ref": "#/components/schemas/RepayLendingResponse"
                },
                "examples": {
                  "Example": {
                    "description": "요청 접수 완료 (처리 결과는 코인빌리기 요청 상태 조회로 확인)",
                    "value": {
                      "request_uuid": "79c7755f-4627-4427-82ae-f2ddbd893afe",
                      "identifier": "repay_partial_20260810_001",
                      "action": "repay",
                      "state": "processing",
                      "failure": null,
                      "created_at": "2026-08-10T00:21:03+09:00",
                      "completed_at": null,
                      "lending_uuid": "34e0c499-2bce-43f7-8eac-26852e79b78f",
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
                    "description": "* repayment_amount 입력 + repayment_mode: \"full\"\n* repayment_amount 미입력 + repayment_mode 미입력 또는 \"amount\"\n* repayment_mode 허용 외 값\n* repayment_mode 타입 오류 (boolean, number, null 등)\n* repayment_amount ≤ 0\n* repayment_amount가 숫자로 해석되지 않음\n* identifier 빈 문자열\n* identifier 길이 초과 (65자 이상)\n* identifier 허용 외 문자\n* identifier에 null 명시",
                    "value": {
                      "error": {
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
                      }
                    }
                  },
                  "exceed_lending_repay_amount": {
                    "description": "repayment_amount가 잔여 빌린 수량을 초과",
                    "value": {
                      "error": {
                        "name": "exceed_lending_repay_amount",
                        "message": "상환 수량을 남은 대여 수량 이하로 입력해 주세요."
                      }
                    }
                  },
                  "insufficient_main_balance": {
                    "description": "포켓에 상환 가능한 수량만큼 잔고가 없음",
                    "value": {
                      "error": {
                        "name": "insufficient_main_balance",
                        "message": "상환에 필요한 수량이 부족해요."
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
                  "invalid_state": {
                    "description": "현재 상환할 수 없음",
                    "value": {
                      "error": {
                        "name": "invalid_state",
                        "message": "현재 상환할 수 없어요.\n잠시 후 다시 시도해 주세요."
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
                  },
                  "terminated_state": {
                    "description": "이미 종료된 포지션",
                    "value": {
                      "error": {
                        "name": "terminated_state",
                        "message": "이미 상환이 완료됐어요."
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
                  }
                }
              }
            }
          },
          "500": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "base_price_must_be_positive": {
                    "description": "기준 시세 확인 불가",
                    "value": {
                      "error": {
                        "name": "base_price_must_be_positive",
                        "message": "기준 시세를 정상적으로 확인할 수 없어 요청을 처리하지 못했어요. 잠시 후 다시 요청해 주세요."
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
              "code": "# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\ncurl --request POST \\\n  --url 'https://api.upbit.com/v1/lendings/repay' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json' \\\n  --header 'Content-Type: application/json' \\\n  --data '{\n    \"uuid\": \"34e0c499-2bce-43f7-8eac-26852e79b78f\",\n    \"repayment_amount\": \"0.00200000\",\n    \"identifier\": \"repay_partial_20260810_001\"\n  }'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/repay\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nbody = {\n    \"uuid\": \"34e0c499-2bce-43f7-8eac-26852e79b78f\",\n    \"repayment_amount\": \"0.00200000\",\n    \"identifier\": \"repay_partial_20260810_001\",\n}\nquery_string = unquote(urlencode(body, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n    \"Content-Type\": \"application/json\",\n}\n# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=body)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/repay\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst body = {\n  uuid: \"34e0c499-2bce-43f7-8eac-26852e79b78f\",\n  repayment_amount: \"0.00200000\",\n  identifier: \"repay_partial_20260810_001\",\n};\nconst queryString = Object.entries(body)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\n// 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\naxios\n  .post(`${BASE_URL}${PATH}`, body, {\n    headers: {\n      Authorization: `Bearer ${jwtToken}`,\n      Accept: \"application/json\",\n      \"Content-Type\": \"application/json\",\n    },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport com.google.gson.Gson;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.MediaType;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\n\npublic class RepayLending {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/repay\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> body = new LinkedHashMap<>();\n        body.put(\"uuid\", \"34e0c499-2bce-43f7-8eac-26852e79b78f\");\n        body.put(\"repayment_amount\", \"0.00200000\");\n        body.put(\"identifier\", \"repay_partial_20260810_001\");\n        String queryString = body.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .post(RequestBody.create(new Gson().toJson(body), MediaType.parse(\"application/json; charset=utf-8\")))\n            .build();\n\n        // 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)