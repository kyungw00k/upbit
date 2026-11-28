---
updatedAt: 2026-11-25T17:14:52.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 페어 체결 이력 조회

지정한 페어의 최근 체결 목록을 조회합니다.

#### Rate Limit

초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 \[<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits#:~:text=%EB%8A%98%EB%A6%B4%20%EC%88%98%20%EC%97%86%EC%8A%B5%EB%8B%88%EB%8B%A4.-,Quotation%20REST%20API,-Rate%20Limit%20%EA%B7%B8%EB%A3%B9">`trade`그룹</Anchor>] 내에서 요청 가능 횟수를 공유합니다.

<br />

# OpenAPI 정의

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "QUOTATION API",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.upbit.com"
    }
  ],
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
    "/v1/trades/ticks": {
      "get": {
        "operationId": "list-pair-trades",
        "summary": "페어 체결 이력 조회",
        "description": "지정한 페어의 최근 체결 목록을 조회합니다.",
        "tags": [
          "Quotation"
        ],
        "parameters": [
          {
            "name": "market",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "조회하고자 하는 페어(거래쌍)",
              "example": "KRW-BTC"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "KRW-BTC"
              }
            }
          },
          {
            "name": "to",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회 대상 일자 내 조회 기간의 종료 시각(UTC).\n지정한 조회 일자 내에서 특정 시간대의 체결 내역을 조회하고자 하는 경우 선택적으로 사용할 수 있는 파라미터입니다.\n\nHHmmss 또는 HH:mm:ss 형식의 시간 포맷으로 입력합니다. 체결 목록이 지정한 시간부터 시간 역순으로 반환됩니다.\n\n[예시] days_ago=1&to=130000 으로 조회하는 경우 1일 전(어제) UTC 기준 오후 1시 정각 이전 체결 이력을 최신순으로 반환\n",
              "example": "134501"
            },
            "allowReserved": true
          },
          {
            "name": "count",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "조회하고자 하는 체결 내역의 개수.\n\n최대 500개 조회를 지원하며 기본 값은 1입니다.\n",
              "example": 500,
              "default": 1
            },
            "allowReserved": true
          },
          {
            "name": "cursor",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "Pagination을 위한 조회 범위 지정용 커서.\n응답에 포함된 체결의 \"sequential_id\" 값을 이 필드에 입력하여 해당 체결 직전 데이터부터 \"count\"개의 이전 체결 이력을 이어서 조회할 수 있습니다.\n",
              "example": 17510687993360000
            },
            "allowReserved": true
          },
          {
            "name": "days_ago",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "체결 내역 조회 대상 일자와 요청 시점과의 일 단위 offset.\n체결 내역은 체결 일을 지정하여 조회해야 하며 최대 7일의 조회 기간을 지원합니다. 일자 구분은 UTC를 기준으로 합니다.\n\n1 이상 7 이하의 정수형으로 입력합니다. 빈 값으로 입력할 경우, 요청 일자에 발생한 체결 목록을 반환하며, 7을 입력하면 7일 이전에 발생한 체결 목록을 시간 역순(최신 체결 순)으로 반환합니다.\n",
              "example": 1
            },
            "allowReserved": true
          }
        ],
        "responses": {
          "200": {
            "description": "List of trades",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Trade"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "trade_date_utc": "2025-06-27",
                        "trade_time_utc": "23:59:59",
                        "timestamp": 1751068799336,
                        "trade_price": 147058000,
                        "trade_volume": 0.00006043,
                        "prev_closing_price": 146852000,
                        "change_price": 206000,
                        "ask_bid": "BID",
                        "sequential_id": 17510687993360000
                      }
                    ]
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
                  "parameter range error": {
                    "value": {
                      "error": {
                        "name": 400,
                        "message": "Invalid parameter. Check the given value!"
                      }
                    }
                  },
                  "parameter type error": {
                    "value": {
                      "error": {
                        "name": 400,
                        "message": "Type mismatch error. Check the parameters type!"
                      }
                    }
                  },
                  "missing required parameter error": {
                    "value": {
                      "error": {
                        "name": 400,
                        "message": "Missing request parameter error. Check the required parameters!"
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
                  "not found error": {
                    "value": {
                      "error": {
                        "name": 404,
                        "message": "Code not found"
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
  },
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
      "Trade": {
        "type": "object",
        "required": [
          "market",
          "trade_date_utc",
          "trade_time_utc",
          "timestamp",
          "trade_price",
          "trade_volume",
          "prev_closing_price",
          "change_price",
          "ask_bid",
          "sequential_id"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n\n[예시] \"KRW-BTC\"\n",
            "example": "KRW-BTC"
          },
          "trade_date_utc": {
            "type": "string",
            "description": "체결 일자 (UTC 기준)\n\n[형식] yyyy-MM-dd\n",
            "example": "2025-06-17"
          },
          "trade_time_utc": {
            "type": "string",
            "description": "체결 시각 (UTC 기준)\n\n[형식] HH:mm:ss\n",
            "example": "16:59:59"
          },
          "timestamp": {
            "type": "integer",
            "format": "int64",
            "description": "체결 시각의 밀리초단위 타임스탬프",
            "example": 1750749886904
          },
          "trade_price": {
            "type": "number",
            "format": "double",
            "description": "최근 체결 가격",
            "example": 142340000
          },
          "trade_volume": {
            "type": "number",
            "format": "double",
            "description": "최근 거래 수량",
            "example": 100000000
          },
          "prev_closing_price": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 (UTC 0시 기준)",
            "example": 140908000
          },
          "change_price": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 대비 가격 변화.\n\"trade_price\" - \"prev_closing_price\"로 계산되며, 현재 종가가 전일 종가보다 얼마나 상승 또는 하락했는지를 나타냅니다.\n  - 양수(+): 현재 종가가 전일 종가보다 상승한 경우\n  - 음수(-): 현재 종가가 전일 종가보다 하락한 경우\n  - 0: 전일 종가와 동일하여 가격 변화가 없는 경우\n",
            "example": 100000000
          },
          "ask_bid": {
            "type": "string",
            "enum": [
              "ASK",
              "BID"
            ],
            "description": "매수/매도 주문 구분",
            "example": "ASK"
          },
          "sequential_id": {
            "type": "integer",
            "description": "체결의 유일 식별자.\n\n해당 필드는 체결 순서를 보장하지 않습니다.\n",
            "example": 17501795992650000
          }
        }
      }
    }
  }
}
```