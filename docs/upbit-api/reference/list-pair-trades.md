---
updatedAt: 2026-07-28T06:00:55.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 페어 체결 이력 조회

지정한 페어의 최근 체결 목록을 조회합니다.

[block:html]
{
  "html": "<div class=\"accordion-changelog\">\n    <input type=\"checkbox\" id=\"api-changelog\">\n    <label for=\"api-changelog\">\n        <div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Revision History <i class=\"fa-solid fa-angle-right\"></i> </div>\n    </label>\n\n    <div class=\"accordion-changelog-content\">\n        <table class=\"custom-table\">\n            <thead>\n                <tr>\n                    <th>반영 버전</th>\n                    <th>반영 일자</th>\n                    <th>변경 사항</th>\n                </tr>\n            </thead>\n            <tbody>\n                <tr>\n                    <td class=\"code-col\">-</td>\n                    <td>2020-07-17</td>\n                    <td><a href=\"https://docs.upbit.com/kr/changelog/update_0717\">조회 기간 7일로 확대 지원</a></td>\n                </tr>\n            </tbody>\n        </table>\n    </div>\n</div>\n\n<div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Rate Limit</div>\n<div class=\"box-rate-limit\">\n  초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 [체결 그룹] 내에서 요청 가능 횟수를 공유합니다.\n</div>"
}
[/block]

# OpenAPI definition

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