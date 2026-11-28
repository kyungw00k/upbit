---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 거래 내역 조회

코인빌리기의 신청, 담보추가, 상환 등 거래 이력을 조회합니다.

지정한 코인빌리기 포지션에서 완료된 거래 내역을 조회합니다. 코인빌리기 신청, 추가신청, 담보추가, 상환, 담보상환, 강제상환은 각각 별도의 거래로 기록됩니다.

여기서 거래는 코인빌리기 요청의 처리 내역을 의미하며, 반대매매로 발생한 자산 매수·매도 체결 내역은 거래별 상세 내역(`entries`)에서 확인할 수 있습니다.

#### 거래 유형

| `type`                  | 설명          |
| ----------------------- | ----------- |
| `borrow`                | 코인빌리기 신청    |
| `borrow_more`           | 코인빌리기 추가 신청 |
| `add_collateral`        | 담보추가        |
| `repay`                 | 빌린 자산으로 상환  |
| `repay_with_collateral` | 담보상환        |
| `liquidation`           | 강제상환        |

* 완료된 거래만 반환하므로 거래 상태 필드는 제공하지 않습니다. 처리 중인 요청의 상태는 [**코인빌리기 요청 상태 조회**](https://docs.upbit.com/kr/reference/get-lending-request)에서 확인해 주세요.
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request)의 `transaction_uuid`를 이용해 해당 요청으로 생성된 거래(`transactions[].uuid`)를 확인할 수 있습니다.

#### 상세 내역(entries)

`include_entries=true`로 요청하면 거래별 자산 이전 및 정산 내역(`entries`)을 함께 조회할 수 있습니다. 기본값은 `false`이며, 이 경우 응답에 `entries` 필드가 포함되지 않습니다.

* `entries`는 최신순으로 정렬됩니다. 거래 내 실제 처리 순서는 `sequence`로 확인할 수 있으며, 1부터 시작합니다.
* 반대매매가 발생한 경우 담보 매도(`collateral_sale`) 또는 빌린 자산 매수(`borrowed_asset_purchase`) 내역으로 기록됩니다. 해당 내역의 `execution`에서 체결 수량, 체결 총액, 평균 체결 가격 및 거래 수수료를 확인할 수 있습니다. 단, 분할 체결(TWAP)의 회차별 내역은 제공하지 않습니다.
* `entries[].fiat_value`는 서로 다른 자산의 가치를 비교하기 위한 원화 환산 참고 금액입니다. 실제 매매가 발생한 경우 체결 금액은 `execution.traded_fiat`에서 확인해 주세요.

| 확인할 내용             | 상세 내역 유형(`entries[].type`)                                    |
| ------------------ | ------------------------------------------------------------- |
| 실제 확정된 빌린 수량       | `borrow_amount_confirmation`                                  |
| 포켓으로 지급된 빌린 자산     | `borrowed_asset_transfer`                                     |
| 담보자산 입금            | `collateral_deposit`                                          |
| 이용수수료 정산           | `borrowing_fee`                                               |
| 강제상환수수료·담보상환수수료 정산 | `collateral_liquidation_fee`                                  |
| 반대매매 체결 결과         | `collateral_sale`, `borrowed_asset_purchase` (`execution` 포함) |
| 정산 후 반환된 자산        | `asset_return`                                                |

#### 상세 내역 유형(entries\[].type)

| 구분           | `entries[].type`             | 설명                                             |
| ------------ | ---------------------------- | ---------------------------------------------- |
| 신청·추가신청·담보추가 | `collateral_deposit`         | 신청 또는 담보추가 시 담보자산 입금                           |
| 신청·추가신청      | `borrow_amount_confirmation` | 실제로 빌릴 수량 확정                                   |
| 신청·추가신청      | `borrowed_asset_transfer`    | 빌린 자산을 신청한 포켓으로 지급                             |
| 상환·청산        | `funding_for_repay`          | 상환에 사용할 빌린 자산 이전                               |
| 상환·청산        | `borrowing_fee`              | 이용수수료 정산                                       |
| 상환·청산        | `collateral_liquidation_fee` | 강제상환수수료 또는 담보상환수수료 정산                          |
| 상환·청산        | `repay`                      | 빌린 수량 정산. 빌린 자산으로 상환하거나 반대매매·담보·남은 원화를 활용하여 상환 |
| 상환·청산        | `asset_return`               | 정산 후 남은 담보·원화·초과 매수 자산 반환                      |
| 상환·청산        | `fiat_funding_for_repay`     | 빌린 자산 매수에 필요한 원화 담보 확보                         |
| 상환·청산        | `collateral_sale`            | 상환에 필요한 자산을 확보하기 위한 담보자산 매도                    |
| 상환·청산        | `borrowed_asset_purchase`    | 상환에 필요한 빌린 자산 매수                               |

#### 조회 기준 및 제약 사항

1. 별도의 정렬 조건을 지정하지 않으면 거래 생성 시각 기준 최신순(`order_by=desc`)으로 조회합니다. `order_by=asc`를 지정하면 오래된순으로 조회할 수 있습니다.
2. 기본 조회 건수는 20건이며, `limit`을 지정하면 한 번에 최대 100건까지 조회할 수 있습니다.
3. 목록은 커서 기반 페이지네이션 방식으로 조회합니다.
   * 최초 요청에서는 `cursor`를 지정하지 않습니다.
   * 다음 페이지가 있는 경우 응답의 `next_cursor`에 현재 페이지의 마지막 거래 UUID가 반환됩니다. 이 값을 다음 요청의 `cursor`에 입력하면 해당 거래 이후의 목록을 조회할 수 있습니다.
   * 다음 페이지를 조회할 때는 `order_by`를 이전 요청과 동일하게 유지해야 합니다.
   * `next_cursor`가 `null`이면 마지막 페이지입니다.
4. 여러 페이지에 걸쳐 조회하는 경우 동일한 시점의 데이터가 보장되지 않습니다. 조회 중 새로운 거래가 생성되거나 거래 정보가 변경될 수 있습니다.

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
      "LendingTransactions": {
        "type": "object",
        "required": [
          "transactions",
          "next_cursor"
        ],
        "properties": {
          "transactions": {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "uuid",
                "type",
                "created_at",
                "updated_at"
              ],
              "properties": {
                "uuid": {
                  "type": "string",
                  "description": "거래의 UUID"
                },
                "type": {
                  "type": "string",
                  "enum": [
                    "borrow",
                    "borrow_more",
                    "add_collateral",
                    "repay",
                    "repay_with_collateral",
                    "liquidation"
                  ],
                  "description": "거래 유형\n\n* `borrow`: 코인빌리기 최초 신청\n* `borrow_more`: 코인빌리기 추가 신청\n* `add_collateral`: 담보추가\n* `repay`: 빌린 자산으로 상환\n* `repay_with_collateral`: 담보를 이용한 상환\n* `liquidation`: 청산"
                },
                "created_at": {
                  "type": "string",
                  "description": "거래 생성 시각 (ISO 8601)"
                },
                "updated_at": {
                  "type": "string",
                  "description": "거래 정보가 마지막으로 변경된 시각 (ISO 8601)"
                },
                "entries": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "required": [
                      "uuid",
                      "sequence",
                      "type",
                      "source",
                      "currency",
                      "amount",
                      "fiat_value",
                      "created_at",
                      "updated_at",
                      "execution"
                    ],
                    "properties": {
                      "uuid": {
                        "type": "string",
                        "description": "상세 내역의 식별자"
                      },
                      "sequence": {
                        "type": "integer",
                        "description": "거래 안에서의 처리 순서"
                      },
                      "type": {
                        "type": "string",
                        "enum": [
                          "collateral_deposit",
                          "borrow_amount_confirmation",
                          "borrowed_asset_transfer",
                          "funding_for_repay",
                          "borrowing_fee",
                          "collateral_liquidation_fee",
                          "repay",
                          "asset_return",
                          "fiat_funding_for_repay",
                          "collateral_sale",
                          "borrowed_asset_purchase"
                        ],
                        "description": "자산 이전·정산 유형\n\n* `collateral_deposit`: 담보자산 입금\n* `borrow_amount_confirmation`: 빌리기 수량 확정\n* `borrowed_asset_transfer`: 빌린 자산 입금\n\n* `funding_for_repay`: 빌린 자산을 상환하기 위해 이동\n* `borrowing_fee`: 빌리기 이용수수료 정산\n* `collateral_liquidation_fee`: 청산 / 담보상환 시 추가 수수료 정산\n* `repay`: 빌린 수량 정산\n* `asset_return`: 정산 후 남은 담보·원화·초과 매수 자산 반환\n* `fiat_funding_for_repay`: 대여자산을 매수하기 위해 필요한 원화 담보 확보\n* `collateral_sale`: 상환에 필요한 자산을 확보하기 위해 담보자산을 매도\n* `borrowed_asset_purchase`: 상환을 위해 빌린 자산을 매수"
                      },
                      "source": {
                        "type": "string",
                        "enum": [
                          "borrowed_asset",
                          "collateral",
                          "fiat",
                          null
                        ],
                        "description": "이전·정산에 사용된 재원\n\n* `borrowed_asset`: 빌린 자산(매수해 확보한 빌린 자산 포함)\n* `collateral`: 담보자산\n* `fiat`: 반대매매 후 확보한 원화",
                        "nullable": true
                      },
                      "currency": {
                        "type": "string",
                        "description": "실제로 거래·차감·반환된 자산의 자산 코드"
                      },
                      "amount": {
                        "type": "string",
                        "description": "실제로 거래·차감·반환된 자산 수량"
                      },
                      "fiat_value": {
                        "type": "string",
                        "description": "내역 발생 시점의 가격으로 `amount`를 원화로 환산한 금액 (단위: KRW, 소수점 이하 8자리로 반올림하며 끝자리 0을 포함해 항상 8자리로 반환)"
                      },
                      "created_at": {
                        "type": "string",
                        "description": "상세 내역 생성 시각 (ISO 8601)"
                      },
                      "updated_at": {
                        "type": "string",
                        "description": "상세 내역이 마지막으로 변경된 시각 (ISO 8601)"
                      },
                      "execution": {
                        "type": "object",
                        "description": "반대매매 실행 결과 집계",
                        "nullable": true,
                        "required": [
                          "traded_amount",
                          "traded_fiat",
                          "average_price",
                          "paid_fee"
                        ],
                        "properties": {
                          "traded_amount": {
                            "type": "string",
                            "description": "체결 수량"
                          },
                          "traded_fiat": {
                            "type": "string",
                            "description": "체결 총액 (단위: KRW)"
                          },
                          "average_price": {
                            "type": "string",
                            "description": "평균 체결 가격 (단위: KRW)",
                            "nullable": true
                          },
                          "paid_fee": {
                            "type": "string",
                            "description": "거래 수수료 총액 (단위: KRW)"
                          }
                        }
                      }
                    }
                  },
                  "description": "거래에서 발생한 자산 이전·정산 내역"
                }
              }
            },
            "description": "지정한 포지션의 거래 목록"
          },
          "next_cursor": {
            "type": "string",
            "description": "다음 페이지 조회에 사용할 cursor",
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
    "/v1/lendings/transactions": {
      "get": {
        "operationId": "list-lending-transactions",
        "summary": "코인빌리기 거래 내역 조회",
        "description": "코인빌리기의 신청, 담보추가, 상환 등 거래 이력을 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "parameters": [
          {
            "name": "lending_uuid",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "조회할 포지션의 UUID"
            }
          },
          {
            "name": "include_entries",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean",
              "description": "거래별 자산 이전·정산 내역(entries)의 포함 여부를 지정",
              "default": false
            }
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "조회할 거래 수",
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
              "description": "다음 페이지 조회의 기준이 되는 거래의 UUID"
            }
          },
          {
            "name": "order_by",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "enum": [
                "desc",
                "asc"
              ],
              "description": "거래 생성 시각 기준 정렬 방식\n\n* `desc`: 최신순\n* `asc`: 오래된순",
              "default": "desc"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/LendingTransactions"
                },
                "examples": {
                  "Without entries (include_entries=false)": {
                    "description": "include_entries=false(기본값)로 요청한 경우",
                    "value": {
                      "transactions": [
                        {
                          "uuid": "dc8af831-e277-49bd-a76b-b13dde769f70",
                          "type": "repay_with_collateral",
                          "created_at": "2026-08-09T20:00:00+09:00",
                          "updated_at": "2026-08-09T20:07:42+09:00"
                        },
                        {
                          "uuid": "06c78c57-4ae8-4825-a8eb-fe462c840a41",
                          "type": "borrow",
                          "created_at": "2026-08-06T15:00:00+09:00",
                          "updated_at": "2026-08-06T15:00:19+09:00"
                        }
                      ],
                      "next_cursor": null
                    }
                  },
                  "With entries (include_entries=true)": {
                    "description": "include_entries=true로 요청한 경우",
                    "value": {
                      "transactions": [
                        {
                          "uuid": "dc8af831-e277-49bd-a76b-b13dde769f70",
                          "type": "repay_with_collateral",
                          "created_at": "2026-08-09T20:00:00+09:00",
                          "updated_at": "2026-08-09T20:07:42+09:00",
                          "entries": [
                            {
                              "uuid": "813513eb-dfa5-4f83-b1f5-e91c64d16087",
                              "sequence": 7,
                              "type": "asset_return",
                              "source": null,
                              "currency": "ETH",
                              "amount": "0.19454380",
                              "fiat_value": "671176.11000000",
                              "created_at": "2026-08-09T20:07:31+09:00",
                              "updated_at": "2026-08-09T20:07:31+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "748e4120-40bc-4b85-841a-97059d31ab02",
                              "sequence": 6,
                              "type": "asset_return",
                              "source": null,
                              "currency": "KRW",
                              "amount": "31171.17",
                              "fiat_value": "31171.17000000",
                              "created_at": "2026-08-09T20:07:27+09:00",
                              "updated_at": "2026-08-09T20:07:27+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "e67ee54d-62f4-409f-9695-c346914c0681",
                              "sequence": 5,
                              "type": "repay",
                              "source": "borrowed_asset",
                              "currency": "BTC",
                              "amount": "0.01000000",
                              "fiat_value": "1010000.00000000",
                              "created_at": "2026-08-09T20:07:20+09:00",
                              "updated_at": "2026-08-09T20:07:20+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "7ed1c47e-02dd-458e-9142-f9c724eeb2ab",
                              "sequence": 4,
                              "type": "collateral_liquidation_fee",
                              "source": "borrowed_asset",
                              "currency": "BTC",
                              "amount": "0.00010000",
                              "fiat_value": "10100.00000000",
                              "created_at": "2026-08-09T20:07:15+09:00",
                              "updated_at": "2026-08-09T20:07:15+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "e3f4c369-f67e-4089-bca0-a61b25b57348",
                              "sequence": 3,
                              "type": "borrowing_fee",
                              "source": "borrowed_asset",
                              "currency": "BTC",
                              "amount": "0.00001500",
                              "fiat_value": "1515.00000000",
                              "created_at": "2026-08-09T20:07:10+09:00",
                              "updated_at": "2026-08-09T20:07:10+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "01392643-d7f5-43c0-ac34-0f734cc18940",
                              "sequence": 2,
                              "type": "borrowed_asset_purchase",
                              "source": null,
                              "currency": "BTC",
                              "amount": "0.01011500",
                              "fiat_value": "1021615.00000000",
                              "created_at": "2026-08-09T20:03:20+09:00",
                              "updated_at": "2026-08-09T20:06:58+09:00",
                              "execution": {
                                "traded_amount": "0.01011500",
                                "traded_fiat": "1021615",
                                "average_price": "101000000",
                                "paid_fee": "510.81"
                              }
                            },
                            {
                              "uuid": "725fed83-4dc5-4e91-be10-3b11bc7e3602",
                              "sequence": 1,
                              "type": "collateral_sale",
                              "source": null,
                              "currency": "ETH",
                              "amount": "0.30545620",
                              "fiat_value": "1053823.89000000",
                              "created_at": "2026-08-09T20:00:21+09:00",
                              "updated_at": "2026-08-09T20:03:01+09:00",
                              "execution": {
                                "traded_amount": "0.30545620",
                                "traded_fiat": "1053823.89",
                                "average_price": "3450000",
                                "paid_fee": "526.91"
                              }
                            }
                          ]
                        },
                        {
                          "uuid": "06c78c57-4ae8-4825-a8eb-fe462c840a41",
                          "type": "borrow",
                          "created_at": "2026-08-06T15:00:00+09:00",
                          "updated_at": "2026-08-06T15:00:19+09:00",
                          "entries": [
                            {
                              "uuid": "cc41e97b-a2a6-4591-aeb8-8ef271db760c",
                              "sequence": 3,
                              "type": "borrowed_asset_transfer",
                              "source": null,
                              "currency": "BTC",
                              "amount": "0.01000000",
                              "fiat_value": "1000000.00000000",
                              "created_at": "2026-08-06T15:00:18+09:00",
                              "updated_at": "2026-08-06T15:00:18+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "73690ffd-a69b-41ae-b081-0fa36116fc76",
                              "sequence": 2,
                              "type": "borrow_amount_confirmation",
                              "source": null,
                              "currency": "BTC",
                              "amount": "0.01000000",
                              "fiat_value": "1000000.00000000",
                              "created_at": "2026-08-06T15:00:12+09:00",
                              "updated_at": "2026-08-06T15:00:12+09:00",
                              "execution": null
                            },
                            {
                              "uuid": "86317a71-c345-4313-88b7-6ac2ccdbad2b",
                              "sequence": 1,
                              "type": "collateral_deposit",
                              "source": null,
                              "currency": "ETH",
                              "amount": "0.50000000",
                              "fiat_value": "1500000.00000000",
                              "created_at": "2026-08-06T15:00:03+09:00",
                              "updated_at": "2026-08-06T15:00:03+09:00",
                              "execution": null
                            }
                          ]
                        }
                      ],
                      "next_cursor": "06c78c57-4ae8-4825-a8eb-fe462c840a41"
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
                    "description": "* lending_uuid 누락\n* include_entries가 boolean 값이 아님\n* limit이 1 미만 또는 100 초과\n* limit이 정수가 아님\n* cursor가 UUID 형식이 아님\n* order_by가 asc, desc 외의 값",
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
                    "description": "lending_uuid에 해당하는 포지션이 없거나 해당 포켓의 포지션이 아님",
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
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/lendings/transactions?lending_uuid=98e408b9-2d78-4a2a-b738-42ae07610410&include_entries=true' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/transactions\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\n# Python의 True는 \"True\"로 변환되므로 문자열 \"true\"로 입력합니다.\nparams = {\n    \"lending_uuid\": \"98e408b9-2d78-4a2a-b738-42ae07610410\",\n    \"include_entries\": \"true\",\n}\nquery_string = unquote(urlencode(params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"}\nres = requests.get(f\"{BASE_URL}{PATH}?{query_string}\", headers=headers)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/transactions\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  lending_uuid: \"98e408b9-2d78-4a2a-b738-42ae07610410\",\n  include_entries: true,\n};\nconst queryString = Object.entries(params)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .get(`${BASE_URL}${PATH}?${queryString}`, {\n    headers: { Authorization: `Bearer ${jwtToken}`, Accept: \"application/json\" },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListLendingTransactions {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/transactions\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new LinkedHashMap<>();\n        params.put(\"lending_uuid\", \"98e408b9-2d78-4a2a-b738-42ae07610410\");\n        params.put(\"include_entries\", \"true\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .get()\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
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
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)