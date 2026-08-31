---
updatedAt: 2026-07-28T06:01:29.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 페어 단위 현재가 조회

지정한 페어의 현재가를 조회합니다. 요청 시점 기준으로 해당 페어의 티커 스냅샷이 반환됩니다.

### 가격 변동 지표

페어 현재가 조회시 반환되는 change, change\_price, change\_rate, signed\_change\_price, signed\_change\_rate 필드는 가격 변동에 관련된 지표를 반환하는 필드들입니다. 해당 변동 지표들은 **전일 종가를 기준으로** 산출합니다.

[block:html]
{
  "html": "<div class=\"accordion-changelog\">\n    <input type=\"checkbox\" id=\"api-changelog\">\n    <label for=\"api-changelog\">\n        <div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Revision History <i class=\"fa-solid fa-angle-right\"></i> </div>\n    </label>\n\n    <div class=\"accordion-changelog-content\">\n        <table class=\"custom-table\">\n            <thead>\n                <tr>\n                    <th>반영 버전</th>\n                    <th>반영 일자</th>\n                    <th>변경 사항</th>\n                </tr>\n            </thead>\n            <tbody>\n                <tr>\n                    <td class=\"code-col\">-</td>\n                    <td>2018-06-21</td>\n                    <td>API 신규 지원</td>\n                </tr>\n            </tbody>\n        </table>\n    </div>\n</div>\n\n<div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Rate Limit</div>\n<div class=\"box-rate-limit\">\n  초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 [현재가 그룹] 내에서 요청 가능 횟수를 공유합니다.\n</div>"
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
    "/v1/ticker": {
      "get": {
        "operationId": "list-tickers",
        "summary": "페어 단위 현재가 조회",
        "description": "지정한 페어의 현재가를 조회합니다. 요청 시점 기준으로 해당 페어의 티커 스냅샷이 반환됩니다.",
        "tags": [
          "Quotation"
        ],
        "parameters": [
          {
            "name": "markets",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "조회하고자 하는 페어(거래쌍) 목록.\n\n2개 이상의 페어에 대해 조회하고자 하는 경우 쉼표(,)로 구분된 문자열 형식으로 요청합니다.\n\n[예시] KRW-BTC,KRW-ETH,BTC-ETH,BTC-XRP\n",
              "example": "KRW-BTC,KRW-ETH"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "KRW-BTC,KRW-ETH"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "List of tickers",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Ticker"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "trade_date": "20250704",
                        "trade_time": "051400",
                        "trade_date_kst": "20250704",
                        "trade_time_kst": "141400",
                        "trade_timestamp": 1751606040365,
                        "opening_price": 148737000,
                        "high_price": 149360000,
                        "low_price": 148288000,
                        "trade_price": 148601000,
                        "prev_closing_price": 148737000,
                        "change": "FALL",
                        "change_price": 136000,
                        "change_rate": 0.0009143656,
                        "signed_change_price": -136000,
                        "signed_change_rate": -0.0009143656,
                        "trade_volume": 0.00016823,
                        "acc_trade_price": 31615925234.05438,
                        "acc_trade_price_24h": 178448329314.96686,
                        "acc_trade_volume": 212.38911576,
                        "acc_trade_volume_24h": 1198.26954807,
                        "highest_52_week_price": 163325000,
                        "highest_52_week_date": "2025-01-20",
                        "lowest_52_week_price": 72100000,
                        "lowest_52_week_date": "2024-08-05",
                        "timestamp": 1751606040403
                      },
                      {
                        "market": "KRW-ETH",
                        "trade_date": "20250704",
                        "trade_time": "051400",
                        "trade_date_kst": "20250704",
                        "trade_time_kst": "141400",
                        "trade_timestamp": 1751606040327,
                        "opening_price": 3518000,
                        "high_price": 3542000,
                        "low_price": 3495000,
                        "trade_price": 3510000,
                        "prev_closing_price": 3517000,
                        "change": "FALL",
                        "change_price": 7000,
                        "change_rate": 0.0019903327,
                        "signed_change_price": -7000,
                        "signed_change_rate": -0.0019903327,
                        "trade_volume": 0.00712453,
                        "acc_trade_price": 20633572185.6844,
                        "acc_trade_price_24h": 111690309851.794,
                        "acc_trade_volume": 5871.4585559,
                        "acc_trade_volume_24h": 31689.60167858,
                        "highest_52_week_price": 5900000,
                        "highest_52_week_date": "2024-12-16",
                        "lowest_52_week_price": 2096000,
                        "lowest_52_week_date": "2025-04-09",
                        "timestamp": 1751606040373
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
      "Ticker": {
        "type": "object",
        "required": [
          "market",
          "trade_date",
          "trade_time",
          "trade_date_kst",
          "trade_time_kst",
          "trade_timestamp",
          "opening_price",
          "high_price",
          "low_price",
          "trade_price",
          "prev_closing_price",
          "change",
          "change_price",
          "change_rate",
          "signed_change_price",
          "signed_change_rate",
          "trade_volume",
          "acc_trade_price",
          "acc_trade_price_24h",
          "acc_trade_volume",
          "acc_trade_volume_24h",
          "highest_52_week_price",
          "highest_52_week_date",
          "lowest_52_week_price",
          "lowest_52_week_date",
          "timestamp"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n\n[예시] \"KRW-BTC\"\n",
            "example": "KRW-BTC"
          },
          "trade_date": {
            "type": "string",
            "description": "최근 체결 일자 (UTC 기준)\n\n[형식] yyyyMMdd\n",
            "example": 20250624
          },
          "trade_time": {
            "type": "string",
            "description": "최근 체결 시각 (UTC 기준)\n\n[형식] HHmmss\n",
            "example": 63514
          },
          "trade_date_kst": {
            "type": "string",
            "description": "최근 체결 일자 (KST 기준)\n\n[형식] yyyyMMdd\n",
            "example": 20250624
          },
          "trade_time_kst": {
            "type": "string",
            "description": "최근 체결 시각 (KST 기준)\n\n[형식] HHmmss\n",
            "example": 153514
          },
          "trade_timestamp": {
            "type": "integer",
            "format": "int64",
            "description": "체결 시각의 밀리초단위 타임스탬프",
            "example": 1750179599265
          },
          "opening_price": {
            "type": "number",
            "format": "double",
            "description": "시가.\n\n해당 페어의 첫 거래 가격입니다.\n",
            "example": 142340000
          },
          "high_price": {
            "type": "number",
            "format": "double",
            "description": "고가.\n\n해당 페어의 최고 거래 가격입니다.\n",
            "example": 142340000
          },
          "low_price": {
            "type": "number",
            "format": "double",
            "description": "저가.\n\n해당 페어의 최저 거래 가격입니다.\n",
            "example": 142340000
          },
          "trade_price": {
            "type": "number",
            "format": "double",
            "description": "종가.\n\n해당 페어의 현재 가격입니다.\n",
            "example": 142340000
          },
          "prev_closing_price": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 (UTC 0시 기준)",
            "example": 140908000
          },
          "change": {
            "type": "string",
            "enum": [
              "EVEN",
              "RISE",
              "FALL"
            ],
            "description": "가격 변동 상태\n  - `EVEN`: 보합\n  - `RISE`: 상승\n  - `FALL`: 하락\n",
            "example": "FALL"
          },
          "change_price": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 대비 가격 변화(절대값)\n\n\"trade_price\" - \"prev_closing_price\"로 계산됩니다.\n",
            "example": 100000000
          },
          "change_rate": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 대비 가격 변화 (절대값)\n\n(\"trade_price\" - \"prev_closing_price\") ÷ \"prev_closing_price\" 으로 계산됩니다.\n",
            "example": 0.0031057175
          },
          "signed_change_price": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 대비 가격 변화.\n\"trade_price\" - \"prev_closing_price\"로 계산되며, 현재 종가가 전일 종가보다 얼마나 상승 또는 하락했는지를 나타냅니다.\n  - 양수(+): 현재 종가가 전일 종가보다 상승한 경우\n  - 음수(-): 현재 종가가 전일 종가보다 하락한 경우\n",
            "example": -454000
          },
          "signed_change_rate": {
            "type": "number",
            "format": "double",
            "description": "전일 종가 대비 가격 변화율\n(\"trade_price\" - \"prev_closing_price\") ÷ \"prev_closing_price\" 으로 계산됩니다.\n  - 양수(+): 가격 상승\n  - 음수(-): 가격 하락\n\n  [예시] 0.015 = 1.5% 상승\n",
            "example": -0.0031057175
          },
          "trade_volume": {
            "type": "number",
            "format": "double",
            "description": "최근 거래 수량",
            "example": 0.0068
          },
          "acc_trade_price": {
            "type": "number",
            "format": "double",
            "description": "누적 거래 금액 (UTC 0시 기준)",
            "example": 72976655137.21983
          },
          "acc_trade_price_24h": {
            "type": "number",
            "format": "double",
            "description": "24시간 누적 거래 금액",
            "example": 344100708901.19037
          },
          "acc_trade_volume": {
            "type": "number",
            "format": "double",
            "description": "누적 거래량 (UTC 0시 기준)",
            "example": 501.38991748
          },
          "acc_trade_volume_24h": {
            "type": "number",
            "format": "double",
            "description": "24시간 누적 거래량",
            "example": 2389.39562434
          },
          "highest_52_week_price": {
            "type": "number",
            "format": "double",
            "description": "52주 신고가",
            "example": 163325000
          },
          "highest_52_week_date": {
            "type": "string",
            "description": "52주 신고가 달성일\n\n[형식] yyyy-MM-dd\n",
            "example": "2025-01-20"
          },
          "lowest_52_week_price": {
            "type": "number",
            "format": "double",
            "description": "52주 신저가",
            "example": 72100000
          },
          "lowest_52_week_date": {
            "type": "string",
            "description": "52주 신저가 달성일\n\n[형식] yyyy-MM-dd\n",
            "example": "2024-08-05"
          },
          "timestamp": {
            "type": "integer",
            "format": "int64",
            "description": "현재가 정보가 반영된 시각의 타임스탬프(ms)",
            "example": 1750746915062
          }
        }
      }
    }
  }
}
```

# Sibling pages

* [마켓 단위 현재가 조회](https://docs.upbit.com/kr/reference/list-quote-tickers.md)