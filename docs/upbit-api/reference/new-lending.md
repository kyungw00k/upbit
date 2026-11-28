---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 신청

지정한 자산과 담보를 기준으로 코인빌리기를 신청합니다.

원화(KRW) 또는 지원되는 디지털 자산을 담보로 코인빌리기를 신청합니다. 한 번의 신청에서 원화만 단독으로 설정하거나 여러 종류의 지원 자산을 담보로 설정할 수 있습니다.

<Callout icon="fad fa-lightbulb-exclamation" theme="warn">
  ### **신청 요청이 접수되었다고 해서 코인빌리기가 완료된 것은 아닙니다.**

  요청이 접수되면 `202 Accepted`와 `request_uuid`가 반환되며, 실제 신청 처리는 이후에 진행됩니다. 따라서 요청 접수 시점에는 포지션이 아직 생성되지 않았을 수 있습니다. 신청 처리 결과는 **코인빌리기 요청 상태 조회** API로 확인해 주세요.
</Callout>

Rate Limit은 API 요청 횟수에 대한 제한으로, 제한 범위 내에서 요청하더라도 서비스 이용 조건에 따라 실제 신청이 완료되지 않을 수 있습니다.

#### 신청 흐름

1. <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/available-lending-information">**코인빌리기 가능 정보 조회**</Anchor>로 빌릴 자산과 담보자산의 정책, 계정별 이용 한도를 확인합니다. 대여한도 및 레벨별 이용 조건은 <Anchor target="_blank" href="https://support.upbit.com/hc/ko/articles/51996138070809-%EC%BD%94%EC%9D%B8%EB%B9%8C%EB%A6%AC%EA%B8%B0-%EB%8C%80%EC%97%AC-%ED%95%9C%EB%8F%84-%EC%95%88%EB%82%B4">코인빌리기 대여 한도 안내</Anchor>를 참고해 주세요.
2. <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-test">**코인빌리기 신청 테스트**</Anchor>로 신청 가능 여부와 예상 LTV(렌딩비율)를 확인합니다.
3. **코인빌리기 신청**을 요청합니다. 요청이 접수되면 `request_uuid`가 반환되며, 요청 상태(`state`)는 `processing`으로 표시됩니다.
4. <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending-request">**코인빌리기 요청 상태 조회**</Anchor>로 최종 처리 결과(`done` 또는 `failed`)를 확인합니다.
5. 신청이 완료되면 `lending_uuid`를 이용해 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">**개별 코인빌리기 조회**</Anchor>에서 빌린 수량, 담보 및 만기일을 확인합니다.

#### 담보 설정

* 원화(KRW)를 포함한 지원 자산을 담보로 설정할 수 있습니다. 담보 가능 자산은 코인빌리기 가능 정보 조회의 `collateral_assets[].currency`에서 확인할 수 있으며, 한 번의 신청에 설정할 수 있는 최대 담보자산 종류 수는 `max_collateral_count`에서 확인해 주세요.
* 담보로 맡길 자산의 원화 평가 금액이 자산별 최소 담보 평가 금액(`collateral_assets[].min_fiat`) 이상이 되도록 수량을 입력해야 합니다.
* 빌릴 자산과 동일한 종류의 자산이나 이미 빌리고 있는 자산은 담보로 설정할 수 없습니다.
* 맡긴 담보는 전체 상환이 완료되기 전까지 출금할 수 없습니다. 상환 및 수수료 정산 후 남은 담보는 반환됩니다.

#### 빌릴 수량

* 신청 처리 시점의 업비트 기준 평가가격과 LTV 계산 결과에 따라 실제 확정 수량이 요청 수량과 다를 수 있습니다.
* 최대 LTV(`max_ltv`)를 소폭 초과하더라도 보정 허용 범위 이내라면, 빌릴 수량을 최대 LTV에 맞춰 줄여 처리할 수 있습니다. 최소 LTV(`min_ltv`)에 미달하거나 보정 허용 범위를 초과하면 `invalid_lending_ratio` 에러가 반환됩니다.
* <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-test">코인빌리기 신청 테스트</Anchor>의 예상 LTV(`estimated_ltv`)는 수량 보정이 적용된 경우 보정 후 비율을 나타냅니다. 다만 테스트 응답에는 예상 확정 수량이 포함되지 않습니다. 요청 수량, 예상 LTV, 확정 수량이 달라지는 예시는 신청 테스트 문서의 '예상 LTV'를 참고해 주세요.
* 신청이 정상 처리(`done`)되더라도 확정 수량이 요청 수량보다 적을 수 있으며, 이는 오류가 아닙니다. 확정된 수량은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor>의 `borrowed_amount` 또는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/list-lending-transactions">코인빌리기 거래 내역 조회</Anchor>의 `borrow_amount_confirmation` 내역에서 확인할 수 있습니다.

#### 이용기간과 상환

* 이용기간은 30일이며, 만기일(`expire_at`)은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor>에서 확인할 수 있습니다.
* 이용기간 중 언제든지 일부 또는 전액을 상환할 수 있습니다. 빌린 자산으로 직접 상환하거나, 맡긴 담보를 활용해 전액 상환할 수 있습니다.
* LTV가 강제상환 기준(`liquidation_ltv`)에 도달하거나 이용기간이 만료되면 강제상환이 진행되며, 강제상환수수료가 부과됩니다.
* 이용수수료는 이용 중에 차감되지 않고 최종 상환 시 정산됩니다. 직접 상환하는 경우 맡긴 담보에서 차감하며, 담보상환하거나 강제상환하는 경우에는 반대매매로 확보한 자산 또는 남은 담보 등으로 정산할 수 있습니다.

#### 질권 설정 계약서

질권 설정 계약서는 업비트 웹·앱에서 확인할 수 있습니다. 계약서 확인에 필요한 신청번호(`lending_number`)는 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-lending">개별 코인빌리기 조회</Anchor>에서 확인해 주세요.

#### 수량 입력 규칙

수량 파라미터는 숫자 형식의 문자열로 입력하며, 소수점 이하 최대 8자리까지 입력할 수 있습니다. 소수점 이하 자릿수가 8자리를 초과하면 에러가 반환됩니다.

#### 응답 값 표기 규칙

코인빌리기 API의 금액, 수량 및 비율 관련 값은 숫자 형식의 문자열로 반환됩니다.

* 자산 수량은 소수점 이하 최대 20자리까지 반환될 수 있습니다.
* 코인빌리기 거래 내역 조회의 `entries[].fiat_value`는 소수점 이하 8자리로 반올림하여 반환됩니다.
* 필드의 용도에 따라 소수점 이하 자릿수가 다를 수 있으므로, 각 필드의 응답 명세를 확인해 주세요.

#### 요청 식별자(identifier)

`identifier`는 클라이언트가 요청마다 직접 지정하는 식별자입니다. 요청 처리 결과를 추적하거나 중복 요청을 식별하는 데 사용하며, 코인빌리기 요청 상태 조회에서 `request_uuid` 대신 사용할 수 있습니다.

* 영문 대소문자, 숫자, 언더스코어(`_`), 점(`.`), 하이픈(`-`)을 조합한 1\~64자 문자열로 입력하며, 대소문자를 구분합니다.
* 동일한 포켓의 모든 코인빌리기 요청(신청, 추가신청, 담보추가, 상환, 담보상환)에서 중복되지 않아야 합니다.
* 요청이 접수되어 `request_uuid`가 발급되면 처리 결과와 관계없이 해당 `identifier`를 다시 사용할 수 없습니다. 요청이 실패(`failed`)한 경우에도 재시도 시 새로운 값을 지정해 주세요.
* 응답을 받지 못했거나 `duplicated_identifier` 에러를 받았다면, 새 요청을 보내기 전에 원래 `identifier`로 [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)를 먼저 호출해 주세요.
* `identifier`를 사용하지 않는 경우 요청 Body에서 해당 필드를 생략해 주세요. 이 경우 응답에서는 `"identifier": null`로 반환됩니다.

#### query\_hash 생성

`query_hash`를 계산할 때 `collaterals`처럼 배열 또는 객체 형태의 파라미터는 공백 없는 JSON 문자열로 직렬화하여 쿼리 문자열에 포함해야 합니다. 나머지 인증 과정은 [인증 가이드](https://docs.upbit.com/kr/reference/auth)를 참고해 주세요.

* **직렬화 예시:** `collaterals=[{"currency":"KRW","amount":"1500000"},{"currency":"ETH","amount":"0.20000000"}]`
* **언어별 처리:** Python의 `json.dumps`는 `separators=(",", ":")`를 지정해야 합니다. JavaScript의 `JSON.stringify`와 Java Gson의 `toJson`은 기본적으로 공백 없이 직렬화합니다. 자세한 구현 방법은 언어별 요청 예제를 참고해 주세요.

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
      "NewLendingResponse": {
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
            "description": "신청 요청의 UUID"
          },
          "identifier": {
            "type": "string",
            "description": "요청 시 지정한 식별자",
            "nullable": true
          },
          "action": {
            "type": "string",
            "enum": [
              "borrow"
            ],
            "description": "요청 유형\n\n* `borrow`: 코인빌리기 최초 신청"
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
    "/v1/lendings": {
      "post": {
        "operationId": "new-lending",
        "summary": "코인빌리기 신청",
        "description": "지정한 자산과 담보를 기준으로 코인빌리기를 신청합니다.",
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
                    "identifier": "borrow_btc_20260808_001"
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
                  "$ref": "#/components/schemas/NewLendingResponse"
                },
                "examples": {
                  "Example": {
                    "description": "요청 접수 완료 (처리 결과는 코인빌리기 요청 상태 조회로 확인)",
                    "value": {
                      "request_uuid": "d861af68-5b9d-4ee0-b5a1-cc22b79ac8ea",
                      "identifier": "borrow_btc_20260808_001",
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
                  },
                  "invalid_price_data": {
                    "description": "시세 정보 확인 불가",
                    "value": {
                      "error": {
                        "name": "invalid_price_data",
                        "message": "시세 정보를 확인할 수 없어요. 잠시 후 다시 시도해 주세요."
                      }
                    }
                  },
                  "invalid_request": {
                    "description": "일시적인 요청 오류",
                    "value": {
                      "error": {
                        "name": "invalid_request",
                        "message": "일시적인 오류가 발생했어요. 잠시 후 다시 시도해 주세요."
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
                  },
                  "exceed_ltv_deviation_limit": {
                    "description": "시세 변동으로 LTV 변경",
                    "value": {
                      "error": {
                        "name": "exceed_ltv_deviation_limit",
                        "message": "시세 변동으로 LTV가 변경되었어요. 빌릴 수량을 변경해 다시 요청해 주세요."
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
              "code": "# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\ncurl --request POST \\\n  --url 'https://api.upbit.com/v1/lendings' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json' \\\n  --header 'Content-Type: application/json' \\\n  --data '{\n    \"borrow_currency\": \"BTC\",\n    \"borrow_amount\": \"0.01000000\",\n    \"collaterals\": [\n      {\"currency\": \"KRW\", \"amount\": \"1500000\"},\n      {\"currency\": \"ETH\", \"amount\": \"0.20000000\"}\n    ],\n    \"identifier\": \"borrow_btc_20260808_001\"\n  }'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport json\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nbody = {\n    \"borrow_currency\": \"BTC\",\n    \"borrow_amount\": \"0.01000000\",\n    \"collaterals\": [\n        {\"currency\": \"KRW\", \"amount\": \"1500000\"},\n        {\"currency\": \"ETH\", \"amount\": \"0.20000000\"},\n    ],\n    \"identifier\": \"borrow_btc_20260808_001\",\n}\n\n\ndef to_hash_value(value):\n    # query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다(separators로 공백 제거).\n    if isinstance(value, (list, dict)):\n        return json.dumps(value, separators=(\",\", \":\"), ensure_ascii=False)\n    return value\n\n\nhash_params = {key: to_hash_value(value) for key, value in body.items()}\nquery_string = unquote(urlencode(hash_params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n    \"Content-Type\": \"application/json\",\n}\n# 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=body)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst body = {\n  borrow_currency: \"BTC\",\n  borrow_amount: \"0.01000000\",\n  collaterals: [\n    { currency: \"KRW\", amount: \"1500000\" },\n    { currency: \"ETH\", amount: \"0.20000000\" },\n  ],\n  identifier: \"borrow_btc_20260808_001\",\n};\n\n// query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다(JSON.stringify 기본 출력).\nconst toHashValue = (value) =>\n  typeof value === \"object\" && value !== null ? JSON.stringify(value) : value;\nconst queryString = Object.entries(body)\n  .map(([key, value]) => `${key}=${toHashValue(value)}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\n// 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\naxios\n  .post(`${BASE_URL}${PATH}`, body, {\n    headers: {\n      Authorization: `Bearer ${jwtToken}`,\n      Accept: \"application/json\",\n      \"Content-Type\": \"application/json\",\n    },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport com.google.gson.Gson;\nimport com.google.gson.GsonBuilder;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.List;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.MediaType;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\n\npublic class NewLending {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings\";\n    // Gson은 공백 없이 직렬화합니다. 문자열 값이 그대로 유지되도록 HTML 이스케이프는 끕니다.\n    private static final Gson GSON = new GsonBuilder().disableHtmlEscaping().create();\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, Object> body = new LinkedHashMap<>();\n        body.put(\"borrow_currency\", \"BTC\");\n        body.put(\"borrow_amount\", \"0.01000000\");\n        body.put(\"collaterals\", List.of(collateral(\"KRW\", \"1500000\"), collateral(\"ETH\", \"0.20000000\")));\n        body.put(\"identifier\", \"borrow_btc_20260808_001\");\n        String queryString = body.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + toHashValue(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .post(RequestBody.create(GSON.toJson(body), MediaType.parse(\"application/json; charset=utf-8\")))\n            .build();\n\n        // 실행하면 실제 요청이 접수됩니다. 요청 내용을 확인한 뒤 실행하세요.\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static Map<String, String> collateral(String currency, String amount) {\n        Map<String, String> item = new LinkedHashMap<>();\n        item.put(\"currency\", currency);\n        item.put(\"amount\", amount);\n        return item;\n    }\n\n    // query_hash 계산 시 배열·객체 값은 공백 없는 JSON 문자열로 바꿉니다.\n    private static String toHashValue(Object value) {\n        return (value instanceof List || value instanceof Map) ? GSON.toJson(value) : String.valueOf(value);\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request.md)
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)