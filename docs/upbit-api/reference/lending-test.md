---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 신청 테스트

코인빌리기 신청 조건을 검증하고 예상 LTV(렌딩비율)를 조회합니다. 실제 신청은 실행되지 않습니다.

코인빌리기 신청 테스트 API는 실제로 신청하지 않고, 신청 API와 동일한 검증을 거쳐 요청 형식과 현재 시점의 신청 가능 여부를 확인하는 API입니다. 요청한 빌릴 수량과 담보를 기준으로 계산한 예상 LTV(`estimated_ltv`)도 함께 반환합니다.

* 실제 코인빌리기 신청 전에 이 API로 신청 가능 여부를 확인하고, 발생할 수 있는 에러를 미리 확인할 수 있습니다.
* 정상 응답이 반환되면 요청 형식이 올바르고 현재 신청 가능한 상태라는 뜻입니다. 에러가 반환되면 에러 코드로 신청할 수 없는 사유를 확인할 수 있습니다.
* 신청 API의 응답에는 LTV가 포함되지 않습니다. 신청 전 예상 LTV는 코인빌리기 신청 테스트 API로 확인해 주세요.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  **테스트 요청은 실제 신청을 생성하지 않습니다.** 검증은 응답 시점에 끝나며, 이후에 처리되는 작업은 없습니다. 응답의 `request_uuid`는 코인빌리기 요청 상태 조회에 사용할 수 없습니다.
</Callout>

#### 테스트 활용 흐름

1. [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information)로 빌릴 자산과 담보자산의 정책, 계정 한도를 확인합니다.
2. 코인빌리기 신청 테스트로 신청 가능 여부와 예상 LTV를 확인합니다. 오류가 발생하면 반환된 에러 코드와 발생 사유를 확인하고, 해당 조건을 해결한 뒤 다시 테스트해 주세요. 자세한 내용은 코인빌리기 API 안내의 에러 안내를 참고해 주세요.
3. 같은 요청 본문으로 [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending)을 요청합니다. 신청 테스트와 신청 API의 요청 파라미터는 같습니다.

#### 예상 LTV

* 예상 LTV(`estimated_ltv`)는 테스트 요청 시점의 시세와 담보를 기준으로 계산합니다.
* `estimated_ltv`는 소수 비율 문자열이며, 소수 자릿수는 응답마다 다를 수 있습니다.
* 신청 가능한 최소 LTV(`min_ltv`)에 미달하거나 최대 LTV(`max_ltv`)의 보정 허용 범위를 초과하면 `invalid_lending_ratio` 에러가 반환됩니다. 최대 LTV를 소폭 초과하더라도 허용 범위 이내라면 빌릴 수량을 줄여 처리할 수 있습니다.
* 수량이 보정된 경우 `estimated_ltv`는 보정 후 예상 LTV를 반환합니다. 다만, 신청 테스트 응답에는 보정된 예상 빌릴 수량이 포함되지 않습니다.
* 실제 신청 시점의 시세와 LTV 계산 결과에 따라 최종 확정 수량이 요청 수량과 다를 수 있습니다.

신청 가능한 LTV 범위는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-lending-information">코인빌리기 가능 정보 조회</Anchor>에서 확인할 수 있습니다.

#### 실제 신청과 다른 점

<Table>
  <thead>
    <tr>
      <th>
        구분
      </th>

      <th style={{ width: "34.9%" }}>
        신청 테스트
      </th>

      <th style={{ width: "45.3%" }}>
        코인빌리기 신청
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        포지션 생성
      </td>

      <td>
        생성하지 않음
      </td>

      <td>
        요청 처리 후 생성
      </td>
    </tr>

    <tr>
      <td>
        응답
      </td>

      <td>
        예상 LTV(`estimated_ltv`) 포함
      </td>

      <td>
        신청 접수 결과 (`state`: `processing`)
      </td>
    </tr>

    <tr>
      <td>
        `request_uuid`
      </td>

      <td>
        `코인빌리기 요청 상태 조회`에 사용할 수 없음
      </td>

      <td>
        `코인빌리기 요청 상태 조회`에 사용
      </td>
    </tr>

    <tr>
      <td>
        `identifier`
      </td>

      <td>
        형식과 기존 사용 여부만 검증하며, 사용된 것으로 간주하지 않음
      </td>

      <td>
        요청이 접수되면 사용된 것으로 간주
      </td>
    </tr>

    <tr>
      <td>
        Rate Limit 그룹
      </td>

      <td>
        코인빌리기 테스트 그룹 (`lending-test`)
      </td>

      <td>
        코인빌리기 그룹(`lending`)
      </td>
    </tr>
  </tbody>
</Table>

#### **수량 입력 규칙**

수량 파라미터는 숫자 형식의 문자열로 입력하며, 소수점 이하 최대 8자리까지 허용합니다. 소수점 이하 자릿수가 8자리를 초과하면 에러가 반환됩니다.

#### 요청 식별자(identifier)

`identifier`는 클라이언트가 요청마다 직접 지정하는 식별자입니다. 실제 코인빌리기 요청에서는 처리 결과를 추적하거나 중복 요청을 식별할 때 사용하며, 코인빌리기 요청 상태 조회에서 `request_uuid` 대신 사용할 수 있습니다. 신청 테스트에서는 형식과 기존 사용 여부만 검증하므로, **테스트에 사용한 값은 실제 신청에 다시 사용할 수 있습니다.**

* `identifier`는 포켓의 모든 코인빌리기 요청(신청, 추가신청, 담보추가, 상환, 담보상환)에서 고유해야 합니다.
* 실제 코인빌리기 요청이 접수되어 `request_uuid`가 발급되면, 처리 결과와 관계없이 해당 `identifier`를 다시 사용할 수 없습니다. 요청이 실패(`failed`)한 경우에도 재사용할 수 없으므로, 재시도할 때는 새로운 값을 지정해 주세요.
* `identifier`를 지정하는 경우, 영문 대소문자, 숫자, 언더스코어(`_`), 점(`.`), 하이픈(`-`)을 조합하여 1\~64자로 입력하며, 대소문자를 구분합니다. 자세한 규칙은 **identifier 사용 규칙**을 참고하세요.
* `identifier`를 사용하지 않는 경우, 요청 Body에서 해당 필드를 생략해 주세요. 이 경우 응답에서는 `"identifier": null`로 반환됩니다.

#### query\_hash 생성

query\_hash를 계산할 때 `collaterals`처럼 값이 배열·객체인 파라미터는 공백 없는 JSON 문자열로 바꿔 쿼리 문자열에 넣습니다. 나머지 과정은 [인증 가이드](https://docs.upbit.com/kr/reference/auth)와 같습니다.

* 예: `collaterals=[{"currency":"KRW","amount":"1500000"},{"currency":"ETH","amount":"0.20000000"}]`
* 공백 없이 직렬화하는 방법은 언어마다 다릅니다. Python `json.dumps`는 `separators=(",", ":")`를 지정해야 하고, JavaScript `JSON.stringify`와 Java Gson `toJson`은 기본 출력에 공백이 없습니다. 언어별 코드는 요청 예제를 참고하세요.
* 코인빌리기 API 요청 시 오류가 발생하면 반환된 에러 코드와 발생 사유를 확인하고, 해당 조건을 해결한 뒤 다시 요청해 주세요. 발생 가능한 오류와 상세 응답 명세는 [코인빌리기 API 안내](https://docs.upbit.com/kr/reference/lending-overview)의 에러 안내 항목을 참고해 주세요.

***

#### 에러 안내

코인빌리기 API 이용 중 에러가 발생하면 반환된 에러 코드와 발생 사유를 확인하고, 해당 조건에 맞게 조치한 뒤 다시 요청해 주세요.

발생 가능한 에러와 상세 응답 명세는 [코인빌리기 API 안내](https://docs.upbit.com/kr/reference/lending-overview)의 **에러 안내** 항목을 참고해 주세요.

#### Rate Limit

초당 최대 10회 호출할 수 있습니다.&#x20;

포켓 단위로 측정되며 \[`lending-test` 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

#### API Key Permission

[인증](https://docs.upbit.com/kr/reference/auth)이 필요한 API로, 메인포켓 API Key만 사용할 수 있습니다. \[코인빌리기] 권한이 설정된 API Key를 사용해야 합니다. 권한 오류(`out_of_scope`)가 발생하면 <Anchor target="_blank" href="https://www.upbit.com/mypage/open_api_management">API Key 관리</Anchor> 메뉴에서 권한 설정을 확인해 주세요.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### 코인빌리기 API 이용 권한

  코인빌리기 권한은 API Key 발급 시 이용 자격(고객 확인, 약관 동의 등)을 확인한 뒤 설정할 수 있습니다. DAXA 온라인 교육을 이수하지 않아도 권한 설정과 조회 API 이용은 가능합니다. 이후 동의를 철회하거나 이용 자격이 바뀌면 공통 에러가 반환됩니다.

  DAXA 온라인 교육을 이수하지 않은 상태에서 이 API를 호출하면 HTTP 403 `not_agreed_daxa_education` 에러가 반환됩니다.
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
      "LendingTestResponse": {
        "type": "object",
        "required": [
          "request_uuid",
          "identifier",
          "estimated_ltv",
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
            "description": "테스트 요청의 식별자"
          },
          "identifier": {
            "type": "string",
            "description": "요청 시 지정한 식별자",
            "nullable": true
          },
          "estimated_ltv": {
            "type": "string",
            "description": "요청 시점 기준으로 계산한 예상 LTV"
          },
          "action": {
            "type": "string",
            "enum": [
              "borrow"
            ],
            "description": "요청 유형\n\n* `borrow`: 코인빌리기 신청"
          },
          "state": {
            "type": "string",
            "enum": [
              "processing"
            ],
            "description": "요청 처리 상태\n\n* `processing`: 테스트 응답에서는 항상 이 값이 반환됩니다. 검증은 응답 시점에 끝나며, 이후에 처리되는 작업은 없습니다."
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
            "description": "생성된 포지션의 UUID",
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
            "description": "생성된 거래의 UUID",
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
    "/v1/lendings/test": {
      "post": {
        "operationId": "lending-test",
        "summary": "코인빌리기 신청 테스트",
        "description": "코인빌리기 신청 조건을 검증하고 예상 LTV(렌딩비율)를 조회합니다. 실제 신청은 실행되지 않습니다.",
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
                  "borrow_currency",
                  "borrow_amount",
                  "collaterals"
                ],
                "properties": {
                  "borrow_currency": {
                    "type": "string",
                    "description": "빌릴 자산의 자산 코드"
                  },
                  "borrow_amount": {
                    "type": "string",
                    "description": "빌릴 수량"
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
                          "description": "담보자산의 자산 코드"
                        },
                        "amount": {
                          "type": "string",
                          "description": "담보 수량"
                        }
                      }
                    },
                    "description": "맡길 담보 목록"
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
                    "borrow_currency": "BTC",
                    "borrow_amount": "0.01000000",
                    "collaterals": [
                      {
                        "currency": "KRW",
                        "amount": "1500000"
                      },
                      {
                        "currency": "ETH",
                        "amount": "0.20000000"
                      }
                    ],
                    "identifier": "borrow_btc_20261225_001"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "202": {
            "description": "신청 가능 여부 검증이 완료되었습니다. 테스트 응답의 `request_uuid`는 요청 상태 조회에 사용할 수 없습니다.",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/LendingTestResponse"
                },
                "examples": {
                  "Example": {
                    "description": "검증 완료 (포지션은 생성되지 않으며, 이후에 처리되는 작업 없음)",
                    "value": {
                      "request_uuid": "d861af68-5b9d-4ee0-b5a1-cc22b79ac8ea",
                      "identifier": "borrow_btc_20261225_001",
                      "estimated_ltv": "0.5043",
                      "action": "borrow",
                      "state": "processing",
                      "failure": null,
                      "created_at": "2026-08-08T14:59:40+09:00",
                      "completed_at": null,
                      "lending_uuid": null,
                      "lending_state": null,
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
                    "description": "* identifier가 빈 문자열\n* identifier 길이 초과 (65자 이상)\n* identifier에 허용되지 않는 문자 포함\n* borrow_currency와 같은 자산이 collaterals[].currency에 포함\n* collaterals가 빈 배열\n* collaterals에 같은 자산 코드를 중복 입력\n* borrow_currency 또는 담보자산의 자산 코드를 지원 대상 자산으로 해석할 수 없음\n* 수량의 소수점 자릿수가 8자리를 초과",
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
                  "request_collateral_amount_less_than_min_amount": {
                    "description": "담보 수량이(collaterals[].amount) 최소 기준 미만",
                    "value": {
                      "error": {
                        "name": "request_collateral_amount_less_than_min_amount",
                        "message": "최소 담보금 이상 입력해 주세요."
                      }
                    }
                  },
                  "below_min_lending_valuation": {
                    "description": "빌릴 수량이(borrow_amount) 최소 기준 미만",
                    "value": {
                      "error": {
                        "name": "below_min_lending_valuation",
                        "message": "신청 금액이 최소 기준에 미달해요."
                      }
                    }
                  },
                  "invalid_lending_ratio": {
                    "description": "LTV가 신청 가능 범위(min_ltv~max_ltv)를 벗어남",
                    "value": {
                      "error": {
                        "name": "invalid_lending_ratio",
                        "message": "신청 가능한 렌딩비율은 {0}%~{1}%예요."
                      }
                    }
                  },
                  "insufficient_lending_budget": {
                    "description": "대여 자산별 신청 가능 수량 소진",
                    "value": {
                      "error": {
                        "name": "insufficient_lending_budget",
                        "message": "현재 신청 가능한 수량이 모두 소진되었어요."
                      }
                    }
                  },
                  "scheduled_liquidation_in_progress": {
                    "description": "예약청산이 진행 중인 자산 포함 (리브랜딩, 거래지원 종료, 액면병합·액면분할 등)",
                    "value": {
                      "error": {
                        "name": "scheduled_liquidation_in_progress",
                        "message": "현재 예약청산이 진행 중인 자산을 포함하고 있어요."
                      }
                    }
                  },
                  "insufficient_collateral_main_balance": {
                    "description": "포켓의 담보 잔고 부족",
                    "value": {
                      "error": {
                        "name": "insufficient_collateral_main_balance",
                        "message": "담보로 설정할 자산의 보유 수량이 부족해요."
                      }
                    }
                  },
                  "exceed_total_lending_count_limit": {
                    "description": "계정의 최대 이용 건수 초과 (`account_limits.max_lending_count` 참고)",
                    "value": {
                      "error": {
                        "name": "exceed_total_lending_count_limit",
                        "message": "최대 이용 건 수를 초과했어요. 신규 신청을 원하시면 이용 중인 건을 먼저 상환해 주세요."
                      }
                    }
                  },
                  "exceed_currency_lending_count_limit": {
                    "description": "자산별 최대 이용 건수 초과 (`lending_assets[].max_lending_count` 참고)",
                    "value": {
                      "error": {
                        "name": "exceed_currency_lending_count_limit",
                        "message": "해당 코인의 최대 이용 건 수를 초과했어요. 신규 신청을 원하시면 이용 중인 건을 먼저 상환해 주세요."
                      }
                    }
                  },
                  "exceed_collateral_count_limit": {
                    "description": "담보자산 종류 수 초과",
                    "value": {
                      "error": {
                        "name": "exceed_collateral_count_limit",
                        "message": "담보자산은 최대 {0}개까지 설정할 수 있어요. 담보 설정을 확인해 주세요."
                      }
                    }
                  },
                  "already_lent_collateral_currency": {
                    "description": "이미 빌리고 있는 자산을 담보로 설정",
                    "value": {
                      "error": {
                        "name": "already_lent_collateral_currency",
                        "message": "대여 중인 자산은 담보로 설정할 수 없어요."
                      }
                    }
                  },
                  "exceed_loan_cap_limit": {
                    "description": "계정의 남은 대여한도 초과",
                    "value": {
                      "error": {
                        "name": "exceed_loan_cap_limit",
                        "message": "최대 대여한도를 초과하여 신청할 수 없어요."
                      }
                    }
                  },
                  "exceed_collateral_amount_limit": {
                    "description": "계정의 남은 담보 한도 초과",
                    "value": {
                      "error": {
                        "name": "exceed_collateral_amount_limit",
                        "message": "최대 설정 가능한 담보금을 초과하여 신청할 수 없어요."
                      }
                    }
                  },
                  "insufficient_level": {
                    "description": "신청 가능한 최소 조건(레벨) 미충족",
                    "value": {
                      "error": {
                        "name": "insufficient_level",
                        "message": "코인빌리기를 신청할 수 없어요. 신청 가능한 최소 조건을 확인해 주세요."
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
                        "message": "고객 확인(KYC)을 완료한 후 이용할 수 있어요. 고객 확인 절차를 완료해 주세요."
                      }
                    }
                  },
                  "not_agreed_daxa_education": {
                    "description": "DAXA 필수 교육 미이수",
                    "value": {
                      "error": {
                        "name": "not_agreed_daxa_education",
                        "message": "DAXA 가이드라인에 따른 필수 교육을 이수해야 코인을 빌릴 수 있어요. 교육을 먼저 이수해 주세요"
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
                  "amount_limit_price_snapshot_changed": {
                    "description": "처리 중 가격 정보 변경",
                    "value": {
                      "error": {
                        "name": "amount_limit_price_snapshot_changed",
                        "message": "가격 정보가 변경되었어요. 잠시 후 다시 시도해 주세요."
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
                        "message": "코인빌리기 서비스를 이용할 수 없어요. 잠시 후 다시 시도해 주세요."
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
              "code": "curl --request POST \\\n  --url 'https://api.upbit.com/v1/lendings/test' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json' \\\n  --header 'Content-Type: application/json' \\\n  --data '{\n    \"borrow_currency\": \"BTC\",\n    \"borrow_amount\": \"0.01000000\",\n    \"collaterals\": [\n      {\"currency\": \"KRW\", \"amount\": \"1500000\"},\n      {\"currency\": \"ETH\", \"amount\": \"0.20000000\"}\n    ],\n    \"identifier\": \"borrow_btc_20261225_001\"\n  }'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport json\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/test\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nbody = {\n    \"borrow_currency\": \"BTC\",\n    \"borrow_amount\": \"0.01000000\",\n    \"collaterals\": [\n        {\"currency\": \"KRW\", \"amount\": \"1500000\"},\n        {\"currency\": \"ETH\", \"amount\": \"0.20000000\"},\n    ],\n    \"identifier\": \"borrow_btc_20261225_001\",\n}\n\n\ndef to_hash_value(value):\n    # query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다(separators로 공백 제거).\n    if isinstance(value, (list, dict)):\n        return json.dumps(value, separators=(\",\", \":\"), ensure_ascii=False)\n    return value\n\n\nhash_params = {key: to_hash_value(value) for key, value in body.items()}\nquery_string = unquote(urlencode(hash_params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n    \"Content-Type\": \"application/json\",\n}\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=body)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/test\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst body = {\n  borrow_currency: \"BTC\",\n  borrow_amount: \"0.01000000\",\n  collaterals: [\n    { currency: \"KRW\", amount: \"1500000\" },\n    { currency: \"ETH\", amount: \"0.20000000\" },\n  ],\n  identifier: \"borrow_btc_20261225_001\",\n};\n\n// query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다(JSON.stringify 기본 출력).\nconst toHashValue = (value) =>\n  typeof value === \"object\" && value !== null ? JSON.stringify(value) : value;\nconst queryString = Object.entries(body)\n  .map(([key, value]) => `${key}=${toHashValue(value)}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .post(`${BASE_URL}${PATH}`, body, {\n    headers: {\n      Authorization: `Bearer ${jwtToken}`,\n      Accept: \"application/json\",\n      \"Content-Type\": \"application/json\",\n    },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport com.google.gson.Gson;\nimport com.google.gson.GsonBuilder;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.List;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.MediaType;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\n\npublic class TestLending {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/test\";\n    // Gson은 공백 없이 직렬화합니다. 문자열 값이 그대로 유지되도록 HTML 이스케이프는 끕니다.\n    private static final Gson GSON = new GsonBuilder().disableHtmlEscaping().create();\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, Object> body = new LinkedHashMap<>();\n        body.put(\"borrow_currency\", \"BTC\");\n        body.put(\"borrow_amount\", \"0.01000000\");\n        body.put(\"collaterals\", List.of(collateral(\"KRW\", \"1500000\"), collateral(\"ETH\", \"0.20000000\")));\n        body.put(\"identifier\", \"borrow_btc_20261225_001\");\n        String queryString = body.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + toHashValue(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .post(RequestBody.create(GSON.toJson(body), MediaType.parse(\"application/json; charset=utf-8\")))\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static Map<String, String> collateral(String currency, String amount) {\n        Map<String, String> item = new LinkedHashMap<>();\n        item.put(\"currency\", currency);\n        item.put(\"amount\", amount);\n        return item;\n    }\n\n    // query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다.\n    private static String toHashValue(Object value) {\n        return (value instanceof List || value instanceof Map) ? GSON.toJson(value) : String.valueOf(value);\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending.md)
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request.md)
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)