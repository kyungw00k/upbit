---
updatedAt: 2026-11-25T17:08:47.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 분(Minute) 캔들 조회

분 단위 캔들 목록을 조회합니다.

### 분캔들 조회 단위(Unit)

분 캔들 조회 API는 Unit Path 파라미터를 지원합니다. Unit은 분단위 내에서 캔들 하나당 조회 단위(캔들의 너비)를 지정하기 위한 파라미터 입니다. 현재 1, 3, 5, 10, 15, 30, 60, 240분 단위로 캔들 너비를 조정할 수 있으며, 예를 들어 `/v1/candles/minutes/5` 로 요청하는 경우 체결 정보를 5분 단위로 Group하여 캔들을 조회합니다.

<Callout icon="far fa-circle-info" theme="info">
  #### 캔들은 해당 시간대에 체결이 발생한 경우에만 생성됩니다.

  해당 캔들의 시작시각부터 종료시각까지 체결이 발생하지 않은 경우 캔들이 생성되지 않으며, 응답에도 포함되지 않습니다. 예를 들어, candle_date_time이 2024-08-31T22:25:00인 분 캔들의 경우 22:25:00(이상)부터 22:26:00(미만)까지 체결이 발생하지 않은 경우 생성되지 않습니다.
</Callout>

***

#### Rate Limit

초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 \[<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits#quotation-rest-api:~:text=%EB%8A%98%EB%A6%B4%20%EC%88%98%20%EC%97%86%EC%8A%B5%EB%8B%88%EB%8B%A4.-,Quotation%20REST%20API,-Rate%20Limit%20%EA%B7%B8%EB%A3%B9">`candle`그룹</Anchor>] 내에서 요청 가능 횟수를 공유합니다.

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
    "/v1/candles/minutes/{unit}": {
      "get": {
        "operationId": "list-candles-minutes",
        "summary": "분(Minute) 캔들 조회",
        "description": "분 단위 캔들 목록을 조회합니다.",
        "tags": [
          "Quotation"
        ],
        "x-readme": {
          "parameter-ordering": [
            "header",
            "path",
            "query",
            "body",
            "form",
            "cookie"
          ]
        },
        "parameters": [
          {
            "name": "unit",
            "in": "path",
            "required": true,
            "schema": {
              "type": "integer",
              "format": "int32",
              "description": "분(Minute) 캔들 단위\n\n캔들 단위를 지정하여 캔들 조회를 할 수 있습니다.\n\n최대 240분(4시간) 캔들까지 조회할 수 있습니다.\n",
              "example": 15,
              "enum": [
                1,
                3,
                5,
                10,
                15,
                30,
                60,
                240
              ]
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": 15
              }
            }
          },
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
                    "$ref": "#/components/schemas/CandleMinute"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "candle_date_time_utc": "2025-07-01T12:00:00",
                        "candle_date_time_kst": "2025-07-01T21:00:00",
                        "opening_price": 145831000,
                        "high_price": 145831000,
                        "low_price": 145752000,
                        "trade_price": 145759000,
                        "timestamp": 1751327999833,
                        "candle_acc_trade_price": 4022470467.03403,
                        "candle_acc_trade_volume": 27.58904602,
                        "unit": 1
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
      "CandleMinute": {
        "allOf": [
          {
            "$ref": "#/components/schemas/Candle"
          },
          {
            "type": "object",
            "required": [
              "unit"
            ],
            "properties": {
              "unit": {
                "type": "integer",
                "format": "int32",
                "enum": [
                  1,
                  3,
                  5,
                  10,
                  15,
                  30,
                  60,
                  240
                ],
                "description": "분(Minute) 캔들 단위",
                "example": 1
              }
            }
          }
        ]
      }
    }
  }
}
```

# 같은 수준의 페이지

* [초(Second) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-seconds.md)
* [일(Day) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-days.md)
* [주(Week) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-weeks.md)
* [월(Month) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-months.md)
* [연(Year) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-years.md)