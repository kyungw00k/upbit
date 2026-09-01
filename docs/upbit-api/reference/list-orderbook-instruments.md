---
updatedAt: 2026-07-28T06:10:02.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 호가 정책 조회

지정한 페어들의 호가 단위(tick_size)와 호가 모아보기 단위(supported_levels) 정보를 조회합니다.

[block:html]
{
  "html": "<div class=\"accordion-changelog\">\n    <input type=\"checkbox\" id=\"api-changelog\">\n    <label for=\"api-changelog\">\n        <div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Revision History <i class=\"fa-solid fa-angle-right\"></i> </div>\n    </label>\n\n    <div class=\"accordion-changelog-content\">\n        <table class=\"custom-table\">\n            <thead>\n                <tr>\n                    <th>반영 버전</th>\n                    <th>반영 일자</th>\n                    <th>변경 사항</th>\n                </tr>\n            </thead>\n            <tbody>\n                <tr>\n                    <td class=\"code-col\">v1.5.9</td>\n                    <td>2025-07-31</td>\n                  <td><a href=\"https://docs.upbit.com/kr/changelog/orderbook-instruments\">'호가 정책 조회' 기능 추가</a></td>\n              \t</tr>\n            </tbody>\n        </table>\n    </div>\n</div>\n<div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Rate Limit</div>\n<div class=\"box-rate-limit\">\n  초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 [호가 그룹] 내에서 요청 가능 횟수를 공유합니다.\n</div>"
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
    "/v1/orderbook/instruments": {
      "get": {
        "operationId": "list-orderbook-instruments",
        "summary": "호가 정책 조회",
        "description": "지정한 페어들의 호가 단위(tick_size)와 호가 모아보기 단위(supported_levels) 정보를 조회합니다.",
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
              "description": "조회하고자 하는 페어(거래쌍) 목록.\n\n2개 이상의 페어에 대해 조회하고자 하는 경우 쉼표(,)로 구분된 문자열 형식으로 요청합니다.\n",
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
            "description": "List of orderbook instruments",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/OrderbookInstrument"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "quote_currency": "KRW",
                        "tick_size": "1000",
                        "supported_levels": [
                          "0",
                          "10000",
                          "100000",
                          "1000000",
                          "10000000",
                          "100000000"
                        ]
                      },
                      {
                        "market": "KRW-ETH",
                        "quote_currency": "KRW",
                        "tick_size": "1000",
                        "supported_levels": [
                          "0",
                          "10000",
                          "100000",
                          "1000000"
                        ]
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
                  "invalid parameter error": {
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
      "OrderbookInstrument": {
        "type": "object",
        "required": [
          "market",
          "quote_currency",
          "tick_size",
          "supported_levels"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n\n[예시] \"KRW-BTC\"\n",
            "example": "KRW-BTC"
          },
          "quote_currency": {
            "type": "string",
            "description": "해당 페어의 마켓 통화 코드(ex. KRW, BTC, USDT)",
            "example": "KRW"
          },
          "tick_size": {
            "type": "string",
            "description": "해당 페어에 적용되는 호가 단위",
            "example": "10000"
          },
          "supported_levels": {
            "type": "array",
            "description": "해당 페어에서 지원하는 호가 모아보기 단위.\n- 0: 기본 호가단위\n- 호가 모아보기 기능은 원화마켓(KRW)에서만 지원합니다.\n  (BTC, USDT 마켓의 경우 0만 존재)\n",
            "items": {
              "type": "string"
            },
            "example": [
              "0",
              "10000",
              "100000",
              "1000000",
              "10000000",
              "100000000"
            ]
          }
        }
      }
    }
  }
}
```

# Sibling pages

* [호가 조회](https://docs.upbit.com/kr/reference/list-orderbooks.md)