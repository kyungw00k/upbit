---
updatedAt: 2026-07-28T05:32:06.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 일(Day) 캔들 조회

일 단위 캔들 목록을 조회합니다.

### 종가 환산 통화(converting\_price\_unit)

'converting\_price\_unit' 파라미터를 지정하여 원화 마켓이 아닌 다른 마켓(예시: BTC마켓)의 일 캔들 조회시 종가('trade\_price')를 원화로 환산하여 조회할 수 있습니다. 파라미터를 원화('KRW')로 지정하는 경우 converted\_trade\_price 필드에 원화로 환산된 종가 정보가 반환됩니다. 현재 원화 환산만을 지원하고 있으며, 추후 지원 통화는 추가될 수 있습니다.

[block:html]
{
  "html": "<div class=\"callout-section callout-section--danger\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i> 캔들은 해당 시간대에 체결이 발생한 경우에만 생성됩니다.\n      </div>\n    해당 캔들의 시작 시각부터 종료 시각까지 체결이 발생하지 않은 경우 캔들이 생성되지 않으며, 응답에도 포함되지 않습니다. \n    예를 들어, candle_date_time이 2024-08-31T00:00:00인 일 캔들의 경우 2024-08-31T00:00:00(이상)부터 2024-09-01T00:00:00(미만)까지 체결이 발생하지 않은 경우 생성되지 않습니다.\n  </div>"
}
[/block]

[block:html]
{
  "html": "<div class=\"accordion-changelog\">\n    <input type=\"checkbox\" id=\"api-changelog\">\n    <label for=\"api-changelog\">\n        <div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Revision History <i class=\"fa-solid fa-angle-right\"></i> </div>\n    </label>\n\n    <div class=\"accordion-changelog-content\">\n        <table class=\"custom-table\">\n            <thead>\n                <tr>\n                    <th>반영 버전</th>\n                    <th>반영 일자</th>\n                    <th>변경 사항</th>\n                </tr>\n            </thead>\n            <tbody>\n                <tr>\n                    <td class=\"code-col\">-</td>\n                    <td>2018-06-21</td>\n                    <td>API 신규 지원</td>\n                </tr>\n            </tbody>\n        </table>\n    </div>\n</div>\n\n<div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Rate Limit</div>\n<div class=\"box-rate-limit\">\n  초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 [캔들 그룹] 내에서 요청 가능 횟수를 공유합니다.\n</div>"
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
    "/v1/candles/days": {
      "get": {
        "operationId": "list-candles-days",
        "summary": "일(Day) 캔들 조회",
        "description": "일 단위 캔들 목록을 조회합니다.",
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
              "description": "조회 기간의 종료 시각.\n지정한 시각 이전 캔들을 조회합니다. 미지정시 요청 시각을 기준으로 최근 캔들이 조회됩니다.\n\nISO 8601 형식의 datetime으로 아래와 같이 요청 할 수 있습니다.\n\n[예시]\n2025-06-24T04:56:53Z\n2025-06-24 04:56:53\n2025-06-24T13:56:53+09:00\n",
              "example": "2024-01-01T00:00:00"
            },
            "allowReserved": true
          },
          {
            "name": "count",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "조회하고자 하는 캔들의 개수.\n\n최대 200개의 캔들 조회를 지원하며, 기본값은 1입니다.\n",
              "example": 200,
              "default": 1
            },
            "allowReserved": true
          },
          {
            "name": "converting_price_unit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "종가 환산 통화.\n종가를 특정 통화로 환산하고자 하는 경우 선택적으로 지정할 수 있습니다. 사용시 응답에 \"converted_trade_price\" 필드가 추가로 반환됩니다.\n\n[예시] \"KRW\"로 지정시 원가로 환산된 종가가 반환됨.\n",
              "example": "KRW"
            },
            "allowReserved": true
          }
        ],
        "responses": {
          "200": {
            "description": "List of candles",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/CandleDay"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "candle_date_time_utc": "2025-06-30T00:00:00",
                        "candle_date_time_kst": "2025-06-30T09:00:00",
                        "opening_price": 147996000,
                        "high_price": 148480000,
                        "low_price": 145740000,
                        "trade_price": 145759000,
                        "timestamp": 1751327999833,
                        "candle_acc_trade_price": 138812096716.42776,
                        "candle_acc_trade_volume": 944.35761221,
                        "prev_closing_price": 147996000,
                        "change_price": -2237000,
                        "change_rate": -0.0151152734
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
      "Candle": {
        "type": "object",
        "required": [
          "market",
          "candle_date_time_utc",
          "candle_date_time_kst",
          "opening_price",
          "high_price",
          "low_price",
          "trade_price",
          "timestamp",
          "candle_acc_trade_price",
          "candle_acc_trade_volume"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n\n[예시] \"KRW-BTC\"\n",
            "example": "KRW-BTC"
          },
          "candle_date_time_utc": {
            "type": "string",
            "description": "캔들 구간의 시작 시각 (UTC 기준)\n\n[형식] yyyy-MM-dd'T'HH:mm:ss\n",
            "example": "2025-06-23T14:29:59"
          },
          "candle_date_time_kst": {
            "type": "string",
            "description": "캔들 구간의 시작 시각 (KST 기준)\n\n[형식] yyyy-MM-dd'T'HH:mm:ss\n",
            "example": "2025-06-23T23:29:59"
          },
          "opening_price": {
            "type": "number",
            "format": "double",
            "description": "시가.\n\n해당 캔들의 첫 거래 가격입니다.\n",
            "example": 142340000
          },
          "high_price": {
            "type": "number",
            "format": "double",
            "description": "고가.\n\n해당 캔들의 최고 거래 가격입니다.\n",
            "example": 142340000
          },
          "low_price": {
            "type": "number",
            "format": "double",
            "description": "저가.\n\n해당 캔들의 최저 거래 가격입니다.\n",
            "example": 142340000
          },
          "trade_price": {
            "type": "number",
            "format": "double",
            "description": "종가.\n\n해당 페어의 현재 가격입니다.\n",
            "example": 142340000
          },
          "timestamp": {
            "type": "integer",
            "format": "int64",
            "description": "해당 캔들의 마지막 틱이 저장된 시각의 타임스탬프 (ms)",
            "example": 1750688999123
          },
          "candle_acc_trade_price": {
            "type": "number",
            "format": "double",
            "description": "해당 캔들 동안의 누적 거래 금액",
            "example": 185042
          },
          "candle_acc_trade_volume": {
            "type": "number",
            "format": "double",
            "description": "해당 캔들 동안의 누적 거래된 디지털 자산의 수량",
            "example": 0.0013
          }
        }
      },
      "CandleDay": {
        "allOf": [
          {
            "$ref": "#/components/schemas/Candle"
          },
          {
            "type": "object",
            "required": [
              "prev_closing_price",
              "change_price",
              "change_rate"
            ],
            "properties": {
              "prev_closing_price": {
                "type": "number",
                "format": "double",
                "description": "전일 종가 (UTC 0시 기준)",
                "example": 147996000
              },
              "change_price": {
                "type": "number",
                "format": "double",
                "description": "전일 종가 대비 가격 변화 (trade_price - prev_closing_price)",
                "example": -2237000
              },
              "change_rate": {
                "type": "number",
                "format": "double",
                "description": "전일 종가 대비 가격 변화율 ((trade_price - prev_closing_price) / prev_closing_price)",
                "example": -0.0151152734
              },
              "converted_trade_price": {
                "type": "number",
                "format": "double",
                "description": "converting_price_unit 지정시 환산된 종가. 디지털 자산이 아닌 마켓에서는 null.",
                "example": 148601000
              }
            }
          }
        ]
      }
    }
  }
}
```

# Sibling pages

* [초(Second) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-seconds.md)
* [분(Minute) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-minutes.md)
* [주(Week) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-weeks.md)
* [월(Month) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-months.md)
* [연(Year) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-years.md)