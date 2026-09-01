---
updatedAt: 2026-07-29T11:04:10.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 페어 목록 조회

업비트에서 지원하는 모든 페어 목록을 조회합니다.

<HTMLBlock>{`
<div class="accordion-changelog">
    <input type="checkbox" id="api-changelog">
    <label for="api-changelog">
        <div class="APISectionHeader-heading4MUMLbp4_nLs">Revision History <i class="fa-solid fa-angle-right"></i> </div>
    </label>

    <div class="accordion-changelog-content">
        <table class="custom-table">
            <thead>
                <tr>
                    <th>반영 버전</th>
                    <th>반영 일자</th>
                    <th>변경 사항</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td class="code-col">-</td>
                    <td>2024-11-20</td>
                    <td><a href="https://docs.upbit.com/kr/changelog/mytrade_market_warning_deprecated_11_20"> market_event 필드 신규 지원,<br>market_warning 필드 필수 여부 변경 </a></td>
              	</tr>
								<tr>
                    <td class="code-col">-</td>
                    <td>2024-02-22</td>
                    <td><a href="https://docs.upbit.com/kr/changelog/voc_update"> 페어별 시장경보 조회 지원</a></td>
                </tr>
								<tr>
                    <td class="code-col">-</td>
                    <td>-</td>
                    <td><a href="https://docs.upbit.com/kr/changelog/open-api-개선사항-안내-투자-유의-종목-필드-추가"> is_details 파라미터 지원</a></td>
                </tr>
            </tbody>
        </table>
    </div>
</div>

<div class="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>
<div class="box-rate-limit">
  초당 최대 10회 호출할 수 있습니다. IP 단위로 측정되며 [마켓 그룹] 내에서 요청 가능 횟수를 공유합니다.
</div>
`}</HTMLBlock>

<br />

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
    "/v1/market/all": {
      "get": {
        "operationId": "list-trading-pairs",
        "summary": "페어 목록 조회",
        "description": "업비트에서 지원하는 모든 페어 목록을 조회합니다.",
        "tags": [
          "Quotation"
        ],
        "parameters": [
          {
            "name": "is_details",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean",
              "description": "상세 정보를 포함한 조회 여부.\n\ntrue로 지정하여 유의종목 지정 여부, 주의종목 지정 여부와 같은 상세 정보를 응답에 포함할 수 있습니다. 기본값은 false입니다.\n",
              "example": true
            },
            "examples": {
              "default": {
                "value": false
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "List of trading pairs",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Market"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "market": "KRW-BTC",
                        "korean_name": "비트코인",
                        "english_name": "Bitcoin",
                        "market_event": {
                          "warning": false,
                          "caution": {
                            "PRICE_FLUCTUATIONS": false,
                            "TRADING_VOLUME_SOARING": false,
                            "DEPOSIT_AMOUNT_SOARING": false,
                            "GLOBAL_PRICE_DIFFERENCES": false,
                            "CONCENTRATION_OF_SMALL_ACCOUNTS": false
                          }
                        }
                      },
                      {
                        "market": "KRW-ETH",
                        "korean_name": "이더리움",
                        "english_name": "Ethereum",
                        "market_event": {
                          "warning": true,
                          "caution": {
                            "PRICE_FLUCTUATIONS": false,
                            "TRADING_VOLUME_SOARING": false,
                            "DEPOSIT_AMOUNT_SOARING": false,
                            "GLOBAL_PRICE_DIFFERENCES": false,
                            "CONCENTRATION_OF_SMALL_ACCOUNTS": false
                          }
                        }
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
      "Market": {
        "type": "object",
        "required": [
          "market",
          "korean_name",
          "english_name"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n\n[예시] \"KRW-BTC\"\n",
            "example": "KRW-BTC"
          },
          "korean_name": {
            "type": "string",
            "description": "해당 디지털 자산의 한글명",
            "example": "비트코인"
          },
          "english_name": {
            "type": "string",
            "description": "해당 디지털 자산의 영문명",
            "example": "Bitcoin"
          },
          "market_event": {
            "description": "종목 경보 정보",
            "type": "object",
            "properties": {
              "warning": {
                "type": "boolean",
                "description": "유의 종목 여부.\n업비트의 시장경보 시스템에 따라 해당 페어가 유의 종목으로 지정되었는지 여부를 나타냅니다.\n",
                "example": false
              },
              "caution": {
                "description": "주의 종목 여부.\n\n주의 종목으로 지정된 경우, 아래의 세부 경보 유형 중 하나 이상에 해당될 수 있습니다.\n",
                "type": "object",
                "properties": {
                  "PRICE_FLUCTUATIONS": {
                    "type": "boolean",
                    "description": "가격 급등락 경보",
                    "example": false
                  },
                  "TRADING_VOLUME_SOARING": {
                    "type": "boolean",
                    "description": "거래량 급증 경보",
                    "example": false
                  },
                  "DEPOSIT_AMOUNT_SOARING": {
                    "type": "boolean",
                    "description": "입금량 급증 경보",
                    "example": false
                  },
                  "GLOBAL_PRICE_DIFFERENCES": {
                    "type": "boolean",
                    "description": "국내외 가격 차이 경보",
                    "example": false
                  },
                  "CONCENTRATION_OF_SMALL_ACCOUNTS": {
                    "type": "boolean",
                    "description": "소수 계정 집중 거래 경보",
                    "example": false
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
```