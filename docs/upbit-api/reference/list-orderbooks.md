---
updatedAt: 2026-07-28T06:06:55.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 호가 조회

지정한 종목들의 실시간 호가(Orderbook) 정보를 조회합니다.

### 호가 모아보기 (level)

원화마켓(KRW)에서만 지원하는 기능으로, 지정한 단위로 ask/bid price와 size를 모아(group) 조회할 수 있습니다. 숫자 형식의 String으로 요청합니다.

\[예시]: KRW-BTC 종목에 대해 level=100000으로 요청시 10만원(KRW) 단위로 ask/bid price가 반환되며, 각 금액대 내에 포진된 매수/매도 주문량의 합산이 size로 반환됩니다.

종목별 호가 단위에 따라 지원하는 모아보기 단위가 다릅니다. 지원하는 모아보기 단위 정보는 [마켓별 주문 정책](https://docs.upbit.com/kr/docs/faq-market-policy) 문서 또는 [호가 정책 조회](https://docs.upbit.com/kr/reference/list-orderbook-instruments) API 응답을 참고하여 사용해주시기 바랍니다. 미지원 단위를 지정하여 요청하는 경우 빈 배열이 반환되므로 호출 전 지원하는 단위를 반드시 확인해주시기 바랍니다.

[block:html]
{
  "html": "<div class=\"accordion-changelog\">\n    <input type=\"checkbox\" id=\"api-changelog\">\n    <label for=\"api-changelog\">\n        <div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Revision History <i class=\"fa-solid fa-angle-right\"></i> </div>\n    </label>\n\n    <div class=\"accordion-changelog-content\">\n        <table class=\"custom-table\">\n            <thead>\n                <tr>\n                    <th>반영 버전</th>\n                    <th>반영 일자</th>\n                    <th>변경 사항</th>\n                </tr>\n            </thead>\n            <tbody>\n                <tr>\n                    <td class=\"code-col\">v1.5.8</td>\n                    <td>2025-07-02</td>\n                  <td><a href=\"https://docs.upbit.com/kr/changelog/rest_orderbook_unit_options\"><code>count</code> 파라미터 신규 지원, 최대 30호가 지원</a></td>\n              \t</tr>\n\t\t\t\t\t\t\t\t<tr>\n                    <td class=\"code-col\">v1.4.4</td>\n                    <td>2024-01-22</td>\n                  <td><a href=\"https://docs.upbit.com/kr/changelog/tick_unit_change\">호가 모아보기 기능 신규 지원(원화 마켓)</a></td>\n                </tr>\n            </tbody>\n        </table>\n    </div>\n</div>\n\n<div class=\"APISectionHeader-heading4MUMLbp4_nLs\">Rate Limit</div>\n<div class=\"box-rate-limit\">\n  초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 [호가 그룹] 내에서 요청 가능 횟수를 공유합니다.\n</div>"
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
    "/v1/orderbook": {
      "get": {
        "operationId": "list-orderbooks",
        "summary": "호가 조회",
        "description": "지정한 종목들의 실시간 호가(Orderbook) 정보를 조회합니다.",
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
          },
          {
            "name": "level",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "호가 모아보기 단위.\n원화마켓(KRW)에서만 지원하는 기능으로, 지정한 단위로 ask/bid price와 size를 모아(group) 조회할 수 있습니다. 숫자 형식의 String으로 요청합니다.\n\n0 또는 1 이상의 모아보기 단위인 경우 소수점을 포함하지 않은 정수형 문자열로, 1 미만의 소수점 단위 모아보기 단위인 경우 double형 문자열로 요청합니다. 미지정시 기본값은 0입니다.\n",
              "example": "1",
              "default": "0"
            },
            "allowReserved": true
          },
          {
            "name": "count",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "조회하고자 하는 호가 쌍의 개수.\n\n최고 매수 호가 - 최저 매도 호가의 쌍을 기준으로, 지정한 개수 만큼의 호가 쌍 정보가 반환됩니다. 미지정시 기본값은 30입니다.\n",
              "example": 10,
              "default": 30
            },
            "allowReserved": true
          }
        ],
        "responses": {
          "200": {
            "description": "List of orderbook",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Orderbook"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "timestamp": 1751606867762,
                        "total_ask_size": 10.37591054,
                        "total_bid_size": 9.49577219,
                        "orderbook_units": [
                          {
                            "ask_price": 148520000,
                            "bid_price": 148490000,
                            "ask_size": 0.0134662,
                            "bid_size": 0.04296774
                          },
                          {
                            "ask_price": 148590000,
                            "bid_price": 148480000,
                            "ask_size": 0.00075646,
                            "bid_size": 0.00336561
                          },
                          {
                            "ask_price": 148610000,
                            "bid_price": 148470000,
                            "ask_size": 0.34336092,
                            "bid_size": 0.05256162
                          },
                          {
                            "ask_price": 148620000,
                            "bid_price": 148460000,
                            "ask_size": 0.32169944,
                            "bid_size": 0.24641946
                          },
                          {
                            "ask_price": 148630000,
                            "bid_price": 148450000,
                            "ask_size": 0.00092837,
                            "bid_size": 0.04871907
                          },
                          {
                            "ask_price": 148640000,
                            "bid_price": 148440000,
                            "ask_size": 0.10529204,
                            "bid_size": 0.02047981
                          },
                          {
                            "ask_price": 148670000,
                            "bid_price": 148430000,
                            "ask_size": 0.0227922,
                            "bid_size": 0.01414705
                          },
                          {
                            "ask_price": 148690000,
                            "bid_price": 148420000,
                            "ask_size": 0.0801863,
                            "bid_size": 0.09912504
                          },
                          {
                            "ask_price": 148700000,
                            "bid_price": 148410000,
                            "ask_size": 1.51712604,
                            "bid_size": 0.07080781
                          },
                          {
                            "ask_price": 148710000,
                            "bid_price": 148400000,
                            "ask_size": 0.0348708,
                            "bid_size": 0.17449241
                          },
                          {
                            "ask_price": 148720000,
                            "bid_price": 148390000,
                            "ask_size": 0.13259751,
                            "bid_size": 0.08526862
                          },
                          {
                            "ask_price": 148730000,
                            "bid_price": 148380000,
                            "ask_size": 0.00067236,
                            "bid_size": 0.20565151
                          },
                          {
                            "ask_price": 148740000,
                            "bid_price": 148370000,
                            "ask_size": 0.02441535,
                            "bid_size": 0.63581042
                          },
                          {
                            "ask_price": 148750000,
                            "bid_price": 148360000,
                            "ask_size": 0.20022932,
                            "bid_size": 0.0000674
                          },
                          {
                            "ask_price": 148760000,
                            "bid_price": 148350000,
                            "ask_size": 2.15099186,
                            "bid_size": 0.00808661
                          },
                          {
                            "ask_price": 148770000,
                            "bid_price": 148340000,
                            "ask_size": 0.44039467,
                            "bid_size": 0.46319466
                          },
                          {
                            "ask_price": 148780000,
                            "bid_price": 148330000,
                            "ask_size": 0.13921986,
                            "bid_size": 2.14774855
                          },
                          {
                            "ask_price": 148790000,
                            "bid_price": 148320000,
                            "ask_size": 1.36430777,
                            "bid_size": 0.21213562
                          },
                          {
                            "ask_price": 148800000,
                            "bid_price": 148310000,
                            "ask_size": 0.15378602,
                            "bid_size": 0.29317524
                          },
                          {
                            "ask_price": 148810000,
                            "bid_price": 148300000,
                            "ask_size": 0.15674438,
                            "bid_size": 1.27687911
                          },
                          {
                            "ask_price": 148820000,
                            "bid_price": 148290000,
                            "ask_size": 0.61109469,
                            "bid_size": 0.69204352
                          },
                          {
                            "ask_price": 148830000,
                            "bid_price": 148280000,
                            "ask_size": 0.01054949,
                            "bid_size": 0.48827162
                          },
                          {
                            "ask_price": 148840000,
                            "bid_price": 148270000,
                            "ask_size": 0.25879902,
                            "bid_size": 0.03685685
                          },
                          {
                            "ask_price": 148850000,
                            "bid_price": 148260000,
                            "ask_size": 0.20287088,
                            "bid_size": 0.32931252
                          },
                          {
                            "ask_price": 148860000,
                            "bid_price": 148250000,
                            "ask_size": 0.03347512,
                            "bid_size": 0.3973801
                          },
                          {
                            "ask_price": 148870000,
                            "bid_price": 148240000,
                            "ask_size": 0.6276239,
                            "bid_size": 0.09035009
                          },
                          {
                            "ask_price": 148880000,
                            "bid_price": 148230000,
                            "ask_size": 1.20602569,
                            "bid_size": 0.09511834
                          },
                          {
                            "ask_price": 148890000,
                            "bid_price": 148220000,
                            "ask_size": 0.13032101,
                            "bid_size": 0.14136818
                          },
                          {
                            "ask_price": 148900000,
                            "bid_price": 148210000,
                            "ask_size": 0.01613179,
                            "bid_size": 0.01734704
                          },
                          {
                            "ask_price": 148910000,
                            "bid_price": 148200000,
                            "ask_size": 0.07518108,
                            "bid_size": 1.10662057
                          }
                        ],
                        "level": 10000
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
      "Orderbook": {
        "type": "object",
        "required": [
          "market",
          "timestamp",
          "total_ask_size",
          "total_bid_size",
          "orderbook_units",
          "level"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n\n[예시] \"KRW-BTC\"\n",
            "example": "KRW-BTC"
          },
          "timestamp": {
            "type": "integer",
            "format": "int64",
            "description": "조회 요청 시각의 타임스탬프 (ms)",
            "example": 1750749886904
          },
          "total_ask_size": {
            "type": "number",
            "format": "double",
            "description": "현재 호가의 전체 매도 잔량 합계",
            "example": 10.64580882
          },
          "total_bid_size": {
            "type": "number",
            "format": "double",
            "description": "현재 호가의 전체 매수 잔량 합계",
            "example": 8.26105616
          },
          "orderbook_units": {
            "description": "리스트에는 호가 정보가 들어가며 차례대로 1호가, 2호가 ... 30호가의 정보를 담고 있습니다.",
            "type": "array",
            "items": {
              "type": "object",
              "description": "호가 단위 목록",
              "required": [
                "ask_price",
                "bid_price",
                "ask_size",
                "bid_size"
              ],
              "properties": {
                "ask_price": {
                  "type": "number",
                  "format": "double",
                  "description": "매도 호가",
                  "example": 146030000
                },
                "bid_price": {
                  "type": "number",
                  "format": "double",
                  "description": "매수 호가",
                  "example": 146000000
                },
                "ask_size": {
                  "type": "number",
                  "format": "double",
                  "description": "매도 잔량",
                  "example": 0.74547181
                },
                "bid_size": {
                  "type": "number",
                  "format": "double",
                  "description": "매수 잔량",
                  "example": 0.20338834
                }
              }
            }
          },
          "level": {
            "type": "number",
            "format": "double",
            "description": "해당 호가가 적용된 가격 단위",
            "example": 10000,
            "default": 0
          }
        }
      }
    }
  }
}
```

# Sibling pages

* [호가 정책 조회](https://docs.upbit.com/kr/reference/list-orderbook-instruments.md)