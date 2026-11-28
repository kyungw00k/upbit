---
updatedAt: 2026-11-25T17:07:57.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 초(Second) 캔들 조회

초 단위 캔들 목록을 조회합니다.

<Callout icon="far fa-circle-info" theme="info">
  #### 초 캔들(초봉) 조회 지원 기간은 최대 3개월입니다.

  초 단위 캔들 조회 API는 요청 시점을 기준으로 최근 3개월 이내의 데이터만 제공합니다. 조회 가능 기간을 초과한 경우, 응답이 빈 리스트로 반환되거나, 요청한 개수(count)만큼 반환되지 않을 수 있습니다. 조회 가능 기간은 to 파라미터를 활용하여 확인하시기 바랍니다.
</Callout>

<Callout icon="far fa-circle-info" theme="info">
  #### 캔들은 해당 시간대에 체결이 발생한 경우에만 생성됩니다.

  해당 캔들의 시작 시각부터 종료 시각까지 체결이 발생하지 않은 경우 캔들이 생성되지 않으며, 응답에도 포함되지 않습니다. 예를 들어, candle_date_time이 2024-08-31T22:25:00인 초 캔들의 경우 22:25:00(이상)부터 22:25:01(미만)까지 체결이 발생하지 않은 경우 생성되지 않습니다.
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
    "/v1/candles/seconds": {
      "get": {
        "operationId": "list-candles-seconds",
        "summary": "초(Second) 캔들 조회",
        "description": "초 단위 캔들 목록을 조회합니다.",
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
              "description": "조회 기간의 종료 시각.\n지정한 시각 이전 캔들을 조회합니다. 미지정시 요청 시각을 기준으로 최근 캔들이 조회됩니다.\n\nISO 8601 형식의 datetime으로 다음과 같이 요청 할 수 있습니다. 실제 요청 시에는 공백 및 특수문자가 정상적으로 처리되도록 URL 인코딩을 수행해야 합니다.\n\n[예시] 2025-06-24T04:56:53Z\n2025-06-24 04:56:53\n2025-06-24T13:56:53+09:00\n\n초 캔들은 요청 시점으로부터 최대 3개월 이전 데이터까지의 조회만 지원하므로, 3개월 이전 시각을 지정하는 경우 응답이 빈 배열로 반환됩니다.\n",
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
                    "$ref": "#/components/schemas/Candle"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "candle_date_time_utc": "2025-06-30T23:59:59",
                        "candle_date_time_kst": "2025-07-01T08:59:59",
                        "opening_price": 145794000,
                        "high_price": 145794000,
                        "low_price": 145759000,
                        "trade_price": 145759000,
                        "timestamp": 1751327999833,
                        "candle_acc_trade_price": 185042,
                        "candle_acc_trade_volume": 0.0013
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
      }
    }
  }
}
```

# 같은 수준의 페이지

* [분(Minute) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-minutes.md)
* [일(Day) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-days.md)
* [주(Week) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-weeks.md)
* [월(Month) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-months.md)
* [연(Year) 캔들 조회](https://docs.upbit.com/kr/reference/list-candles-years.md)