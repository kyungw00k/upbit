---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 개별 코인빌리기 조회

특정 코인빌리기의 이용 상태와 상세 정보를 조회합니다.

포지션의 UUID로 개별 코인빌리기의 상세 정보를 조회합니다. 이용 중인 포지션은 조회 시점의 현재 LTV(렌딩비율)를 확인하여 강제상환 위험을 관리할 수 있습니다.

조회할 포지션의 UUID는 [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)의 `lending_uuid` 또는 [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings)의 `lendings[].uuid`에서 확인할 수 있습니다.

#### 현재 LTV 확인

현재 LTV(`current_ltv`)는 포지션 상태가 `active`(이용 중)인 경우에만 반환됩니다. `suspended`, `liquidating`, `repaid`, `liquidated` 상태에서는 `null`을 반환합니다.

#### LTV(렌딩비율) 계산

LTV는 빌린 자산의 평가금액을 전체 담보 인정금액으로 나눈 비율입니다.

**LTV = 빌린 자산 평가금액 ÷ 전체 담보 인정금액**

담보 인정금액은 각 담보자산의 평가금액에 담보 인정비율을 적용하여 계산합니다.

**담보 인정금액 = 담보 수량 × 담보 평가가격 × 담보 인정비율**

여러 자산을 담보로 설정한 경우, 각 담보자산의 인정금액을 합산하여 전체 담보 인정금액을 계산합니다.

LTV는 백분율이 아닌 소수 형태의 문자열로 반환됩니다. 예를 들어 `"0.6643"`은 LTV 66.43%를 의미합니다.

현재 LTV는 조회 시점의 업비트 평가가격을 기준으로 계산합니다. 자세한 가격 기준은 [코인빌리기 가이드 FAQ](https://www.upbit.com/lending/guide?tab=faq)를 참고해 주세요"

| 필드                | 설명                      | 값이 제공되는 상태 |
| ----------------- | ----------------------- | ---------- |
| `current_ltv`     | 조회 시점에 계산한 현재 LTV       | `active`   |
| `liquidation_ltv` | 해당 포지션에 적용된 강제상환 기준 LTV | 모든 상태      |

<Callout icon="far fa-triangle-exclamation" theme="warn">
  **이용 중에는&#x20;**`current_ltv`**를 주기적으로 확인해 주세요.**

  현재 LTV가 강제상환 기준 LTV(`liquidation_ltv`)에 도달하면 강제상환이 진행되며, 강제상환수수료가 부과됩니다.

  강제상환 위험이 높아지는 경우 담보를 추가하거나 빌린 자산을 일부 상환하여 LTV를 낮출 수 있습니다. 단, 담보를 추가한 포지션은 이후 코인빌리기 추가 신청을 할 수 없습니다.
</Callout>

#### 포지션 상태에 따른 응답

포지션 상태에 따라 제공되는 수량, 수수료 및 종료 관련 정보가 달라집니다.

| 구분        | 포지션 상태                                            | 제공되는 정보                                                                                                    |
| --------- | ------------------------------------------------- | ---------------------------------------------------------------------------------------------------------- |
| 종료 전      | `active`<br />`suspended`<br />`liquidating`      | 남은 빌린 수량(`borrowed_amount`)<br />예상 이용수수료(`estimated_fee`)<br />현재 LTV(`current_ltv`, `active`만 제공)        |
| 종료 후      | `repaid`<br />`liquidated`                        | 종료 직전 빌린 수량(`borrowed_amount`)<br />종료 유형(`closure_type`)<br />최종 수수료(`final_fee`)<br />종료 시각(`closed_at`) |
| 강제상환으로 종료 | `liquidated`<br />(`closure_type`: `liquidation`) | 강제상환 사유(`liquidation_reason`)<br />강제상환 후 미수금(`unpaid_amount`, 미수금 발생 시)                                   |

담보상환과 강제상환은 모두 `liquidated` 상태로 종료되며, 종료 유형은 `closure_type`으로 구분합니다.

* `repay`: 빌린 자산으로 전액 상환
* `repay_with_collateral`: 담보로 전액 상환
* `liquidation`: 강제상환

#### 종료된 포지션의 수량 및 미수금

| 필드                | 설명                                                                                                                            |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `borrowed_amount` | 종료 직전 남아 있던 빌린 수량입니다.<br /><br />전액 상환 후에도 0이 아닌 값이 반환될 수 있으며, 현재 갚아야 할 빌린 수량을 의미하지 않습니다.                                     |
| `unpaid_amount`   | 강제상환 정산 후 갚지 못한 빌린 수량과 미수취 수수료의 합계입니다.<br /><br />빌린 자산(`borrowed_currency`) 기준 수량을 숫자 형식의 문자열로 반환하며, 미수금이 없으면 `null`을 반환합니다. |

* 수량은 소수점 이하 최대 20자리까지 반환될 수 있습니다. 응답 예제의 수량은 소수점 이하 8자리로 표기했습니다.

#### 수수료 확인

* 종료 전에는 `estimated_fee`에서 예상 이용수수료를 확인할 수 있습니다.
* 종료 후에는 `final_fee`에서 확정된 최종 수수료를 확인할 수 있습니다. 최종 수수료에는 이용수수료와 발생한 경우 강제상환수수료 또는 담보상환수수료가 포함됩니다.
* 수수료는 빌린 자산을 기준으로 한 수량으로 반환됩니다. 실제로 차감된 담보자산과 수량은 **코인빌리기 거래 내역 조회**에서 확인해 주세요.

#### 신청번호(lending\_number)

`lending_number`**는 업비트 웹·앱에서 확인할 수 있는 코인빌리기 신청번호입니다.**

질권 설정 계약서와 강제상환 후 미수금의 추심 내역은 업비트 웹·앱에서 해당 신청번호로 확인할 수 있습니다.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  개별 코인빌리기 조회는 코인빌리기 Rate Limit 그룹에 속하며, 신청·추가신청·담보추가·상환 등 다른 코인빌리기 API와 요청 횟수 한도를 공유합니다.

  현재 LTV를 주기적으로 조회하는 경우, 신청·상환 등 다른 요청에 필요한 호출 한도를 고려해 주세요.
</Callout>

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
      "Lending": {
        "type": "object",
        "required": [
          "uuid",
          "state",
          "closure_type",
          "liquidation_reason",
          "current_ltv",
          "liquidation_ltv",
          "borrowed_currency",
          "borrowed_amount",
          "borrowed_price",
          "unpaid_amount",
          "collaterals",
          "lending_number",
          "estimated_fee",
          "final_fee",
          "borrowing_fee_rate",
          "liquidation_fee_rate",
          "collateral_repayment_fee_rate",
          "created_at",
          "expire_at",
          "updated_at",
          "closed_at"
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
              "suspended",
              "liquidating",
              "repaid",
              "liquidated"
            ],
            "description": "포지션 상태\n\n* `active`: 이용 중\n* `suspended`: 이용 제한\n* `liquidating`: 청산·담보상환 진행 중\n* `repaid`: 정상 상환 완료\n* `liquidated`: 청산 완료 (담보상환에 따른 종료 포함)"
          },
          "closure_type": {
            "type": "string",
            "enum": [
              "repay",
              "repay_with_collateral",
              "liquidation",
              null
            ],
            "description": "포지션이 종료된 방식\n\n* `repay`: 빌린 자산으로 전액 상환\n* `repay_with_collateral`: 담보로 전액 상환\n* `liquidation`: 청산",
            "nullable": true
          },
          "liquidation_reason": {
            "type": "string",
            "enum": [
              "ltv_exceeded",
              "expired",
              "scheduled",
              "account_status",
              null
            ],
            "description": "청산 사유\n\n* `ltv_exceeded`: 강제상환 기준 LTV 도달\n* `expired`: 이용기간 만료\n* `scheduled`: 거래지원 종료 등 자산 이벤트에 따른 예약청산\n* `account_status`: 회원 상태, 고객 확인 등 계정 관련 사유",
            "nullable": true
          },
          "current_ltv": {
            "type": "string",
            "description": "조회 시점에 계산한 현재 LTV",
            "nullable": true
          },
          "liquidation_ltv": {
            "type": "string",
            "description": "강제상환 기준 LTV"
          },
          "borrowed_currency": {
            "type": "string",
            "description": "빌린 자산의 자산 코드"
          },
          "borrowed_amount": {
            "type": "string",
            "description": "이용 중에는 남은 빌린 수량, 종료 후에는 종료 직전에 남아 있던 빌린 수량. 전액 상환 후에도 0이 아닌 값이 반환될 수 있으며, 현재 갚아야 할 빌린 수량을 의미하지 않습니다."
          },
          "borrowed_price": {
            "type": "string",
            "description": "빌린 가격 (빌린 자산 1단위당, 단위: KRW)"
          },
          "unpaid_amount": {
            "type": "string",
            "description": "강제상환 정산 후 갚지 못한 빌린 수량과 미수취 수수료의 합계. 빌린 자산(`borrowed_currency`) 기준 수량을 숫자 형식의 문자열로 반환하며, 미수금이 없으면 `null`을 반환합니다.",
            "nullable": true
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
                  "description": "담보 수량"
                },
                "factor": {
                  "type": "string",
                  "description": "담보 인정 비율"
                }
              }
            },
            "description": "담보 목록"
          },
          "lending_number": {
            "type": "string",
            "description": "업비트 웹·앱에서 확인할 수 있는 신청번호"
          },
          "estimated_fee": {
            "type": "object",
            "description": "종료 전 예상 이용수수료",
            "nullable": true,
            "properties": {
              "currency": {
                "type": "string",
                "description": "예상 이용수수료의 자산 코드"
              },
              "amount": {
                "type": "string",
                "description": "빌린 자산 기준 예상 이용수수료 수량"
              }
            }
          },
          "final_fee": {
            "type": "object",
            "description": "종료 후 확정된 최종 수수료",
            "nullable": true,
            "properties": {
              "currency": {
                "type": "string",
                "description": "최종 수수료의 자산 코드"
              },
              "amount": {
                "type": "string",
                "description": "빌린 자산 기준 최종 수수료 합계"
              }
            }
          },
          "borrowing_fee_rate": {
            "type": "string",
            "description": "해당 코인빌리기 건의 최초 신청 시점에 확정된 일 단위 이용 수수료율입니다. 이후 정책 요율이 변경되더라도 해당 건에는 최초 확정된 요율이 계속 적용됩니다."
          },
          "liquidation_fee_rate": {
            "type": "string",
            "description": "해당 코인빌리기 건의 최초 신청 시점에 확정된 강제상환 수수료율입니다. 강제상환 시 적용되며, 이후 정책 요율이 변경되더라도 해당 건에는 최초 확정된 요율이 적용됩니다."
          },
          "collateral_repayment_fee_rate": {
            "type": "string",
            "description": "해당 코인빌리기 건의 최초 신청 시점에 확정된 담보상환 수수료율입니다. 담보상환 시 적용되며, 이후 정책 요율이 변경되더라도 해당 건에는 최초 확정된 요율이 적용됩니다."
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
            "description": "포지션 정보가 마지막으로 변경된 시각 (ISO 8601)"
          },
          "closed_at": {
            "type": "string",
            "description": "포지션의 종료 시각 (ISO 8601)",
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
    "/v1/lending": {
      "get": {
        "operationId": "get-lending",
        "summary": "개별 코인빌리기 조회",
        "description": "특정 코인빌리기의 이용 상태와 상세 정보를 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "parameters": [
          {
            "name": "uuid",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "조회할 포지션의 UUID"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Lending"
                },
                "examples": {
                  "Active lending": {
                    "description": "종료 전 (이용 중)",
                    "value": {
                      "uuid": "5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b",
                      "state": "active",
                      "closure_type": null,
                      "liquidation_reason": null,
                      "current_ltv": "0.6643",
                      "liquidation_ltv": "0.9500",
                      "borrowed_currency": "USDT",
                      "borrowed_amount": "50000.00000000",
                      "borrowed_price": "1390.00",
                      "unpaid_amount": null,
                      "collaterals": [
                        {
                          "currency": "BTC",
                          "amount": "1.00000000",
                          "factor": "0.9000"
                        }
                      ],
                      "lending_number": "12345",
                      "estimated_fee": {
                        "currency": "USDT",
                        "amount": "125.00000000"
                      },
                      "final_fee": null,
                      "borrowing_fee_rate": "0.0005",
                      "liquidation_fee_rate": "0.0100",
                      "collateral_repayment_fee_rate": "0.0100",
                      "created_at": "2026-09-01T10:15:32+09:00",
                      "expire_at": "2026-10-01T10:15:32+09:00",
                      "updated_at": "2026-09-06T09:40:12+09:00",
                      "closed_at": null
                    }
                  },
                  "Liquidated lending": {
                    "description": "종료 후 (강제상환 기준 LTV 도달로 청산된 경우)",
                    "value": {
                      "uuid": "5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b",
                      "state": "liquidated",
                      "closure_type": "liquidation",
                      "liquidation_reason": "ltv_exceeded",
                      "current_ltv": null,
                      "liquidation_ltv": "0.9500",
                      "borrowed_currency": "USDT",
                      "borrowed_amount": "50000.00000000",
                      "borrowed_price": "1390.00",
                      "unpaid_amount": "137.45000000",
                      "collaterals": [
                        {
                          "currency": "BTC",
                          "amount": "1.00000000",
                          "factor": "0.9000"
                        }
                      ],
                      "lending_number": "12345",
                      "estimated_fee": null,
                      "final_fee": {
                        "currency": "USDT",
                        "amount": "625.00000000"
                      },
                      "borrowing_fee_rate": "0.0005",
                      "liquidation_fee_rate": "0.0100",
                      "collateral_repayment_fee_rate": "0.0100",
                      "created_at": "2026-09-01T10:15:32+09:00",
                      "expire_at": "2026-10-01T10:15:32+09:00",
                      "updated_at": "2026-09-06T08:30:44+09:00",
                      "closed_at": "2026-09-06T02:14:08+09:00"
                    }
                  },
                  "Repaid lending": {
                    "description": "종료 후 (빌린 자산으로 전액 상환한 경우)",
                    "value": {
                      "uuid": "3d9f0a6e-1b7c-4e2d-8f5a-9c6b2e4d7a10",
                      "state": "repaid",
                      "closure_type": "repay",
                      "liquidation_reason": null,
                      "current_ltv": null,
                      "liquidation_ltv": "0.9500",
                      "borrowed_currency": "USDT",
                      "borrowed_amount": "50000.00000000",
                      "borrowed_price": "1390.00",
                      "unpaid_amount": null,
                      "collaterals": [
                        {
                          "currency": "BTC",
                          "amount": "1.00000000",
                          "factor": "0.9000"
                        }
                      ],
                      "lending_number": "12346",
                      "estimated_fee": null,
                      "final_fee": {
                        "currency": "USDT",
                        "amount": "125.00000000"
                      },
                      "borrowing_fee_rate": "0.0005",
                      "liquidation_fee_rate": "0.0100",
                      "collateral_repayment_fee_rate": "0.0100",
                      "created_at": "2026-09-01T10:15:32+09:00",
                      "expire_at": "2026-10-01T10:15:32+09:00",
                      "updated_at": "2026-09-06T02:14:08+09:00",
                      "closed_at": "2026-09-06T02:14:08+09:00"
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
                    "description": "`uuid` 누락 또는 UUID 형식 오류",
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
                  "lending_not_found": {
                    "description": "uuid에 해당하는 포지션이 없거나 해당 포켓의 포지션이 아님",
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
          }
        },
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/lending?uuid=5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lending\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\"uuid\": \"5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b\"}\nquery_string = unquote(urlencode(params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"}\nres = requests.get(f\"{BASE_URL}{PATH}?{query_string}\", headers=headers)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lending\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = { uuid: \"5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b\" };\nconst queryString = Object.entries(params)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .get(`${BASE_URL}${PATH}?${queryString}`, {\n    headers: { Authorization: `Bearer ${jwtToken}`, Accept: \"application/json\" },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class GetLending {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lending\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new LinkedHashMap<>();\n        params.put(\"uuid\", \"5cb7cf1e-5de4-4b4f-b4f4-6a5d995bc73b\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .get()\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)