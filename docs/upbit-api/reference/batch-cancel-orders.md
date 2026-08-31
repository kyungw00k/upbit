---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 주문 일괄 취소 접수

조건을 지정하여 해당 조건을 만족하는 최대 300개의 주문을 일괄 취소합니다.

### 주문 상태에 따른 일괄 취소 가능 여부

일괄 취소 API로는 **오직 체결 대기(WAIT) 상태의 주문만 취소**할 수 있습니다. 단, 취소 처리 중 주문 체결이 발생할 수 있으므로 취소요청을 보내는 시점의 주문 잔량과 일괄 취소가 완료된 후 취소된 주문의 주문 잔량은 다를 수 있습니다.

예약 주문(WATCH) 상태의 주문은 일괄 주문 취소로 취소할 수 없으므로 취소를 원하시는 경우 [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order) 또는 [id로 주문 목록 취소 접수](https://docs.upbit.com/kr/reference/cancel-orders-by-ids) API를 사용해주세요.

<br />

### 취소 대상 지정

일괄 취소 대상 주문을 지정하기 위해 다음 조건들을 사용할 수 있습니다.

* **주문 방향**: 매수 또는 매도 주문만 선택적으로 취소할 수 있습니다.
* **페어**: 최대 20개의 취소 대상 페어를 지정하거나, 취소 제외 페어를 지정할 수 있습니다.
* **마켓(호가 자산)**: 특정 마켓의 주문만 일괄 취소할 수 있습니다.
* **취소 주문 수와 정렬 옵션**: 최신 순 또는 오래된 순으로 지정한 수 만큼의 주문을 일괄 취소할 수 있습니다.

취소 대상 페어 또는 대상 마켓에 포함된 페어 중 취소 제외 페어(excluded\_pairs)에 포함된 페어의 경우, **취소 제외 페어 지정을 높은 우선순위로 적용**하여 취소 대상에서 제외됩니다.

<br />

<div className="callout-section callout-section--danger">
  <div className="callout-title">
    <i className="fa-solid fa-circle-exclamation"></i> 주문 취소 요청이 거절되는 경우
    </div>
    다음과 같은 사유로 주문 취소 요청이 거절될 수 있습니다.
    <li>이미 전량 체결 완료되어 취소가 불가한 주문</li>
    <li>이미 취소가 완료된 주문</li>
    <li>해당 주문의 페어가 리브랜딩 등의 이유로 서비스를 일시 정지한 경우</li>
</div>

  <div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> 이 API는 쿼리 파라미터 형식만 지원합니다.
      </div>
    모든 요청 파라미터는 쿼리 파라미터 형식으로만 요청할 수 있습니다. Request Body를 통한 form 형식 또는 json 형식의 파라미터 요청을 지원하지 않습니다. 
    <br /><br />
    [예시] /v1/orders/open?cancel_side=bid
</div>

<div className="callout-section">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  페어(pairs)과 마켓(quote_currencies) 조건은 동시에 사용할 수 없습니다.
      </div>
      둘 중 최대 하나의 파라미터만 포함할 수 있으며, 두 파라미터 모두 지정하지 않는 경우 전체 마켓을 대상으로 일괄 취소 주문이 선택됩니다.
  </div>

<div className="accordion-changelog">
    <input type="checkbox" id="api-changelog" />
    <label for="api-changelog">
        <div className="APISectionHeader-heading4MUMLbp4_nLs">Revision History <i className="fa-solid fa-angle-right"></i> </div>
    </label>

    <div className="accordion-changelog-content">
        



<table className="custom-table">
    <thead>
        <tr>
            <th>반영 버전</th>
            <th>반영 일자</th>
            <th>변경 사항</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td className="code-col">-</td>
            <td>2024-12-11</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_delete_orders"> '주문 일괄 취소' 기능 신규 지원</a></td>
        </tr>
    </tbody>
</table>




    </div>
</div>

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

최대 2 초당 1회 호출할 수 있습니다. 포켓 단위로 측정되며\[Exchange 주문 일괄 취소 그룹]내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>
  <br />
  <div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>
  <div className="box-rate-limit">
    <a href="auth">인증</a>이 필요한 API로, [주문하기] 권한이 설정된 API Key를 사용해야 합니다. <br />
    권한 오류(out_of_scope) 오류가 발생한다면, <a href="https://upbit.com/mypage/open_api_management">API Key 관리 메뉴</a>에서 권한 설정을 확인해주세요.
  </div>

<br />

# OpenAPI definition

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "EXCHANGE API",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.upbit.com"
    }
  ],
  "components": {
    "schemas": {
      "BatchCancelResult": {
        "type": "object",
        "required": [
          "success",
          "failed"
        ],
        "properties": {
          "success": {
            "description": "성공적으로 취소된 주문 정보",
            "type": "object",
            "required": [
              "count",
              "orders"
            ],
            "properties": {
              "count": {
                "type": "integer",
                "description": "성공적으로 취소된 주문 수",
                "example": 5
              },
              "orders": {
                "description": "성공적으로 취소된 주문 목록",
                "type": "array",
                "items": {
                  "type": "object",
                  "required": [
                    "uuid",
                    "market"
                  ],
                  "properties": {
                    "uuid": {
                      "type": "string",
                      "description": "주문의 유일 식별자",
                      "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                    },
                    "market": {
                      "type": "string",
                      "description": "페어(거래쌍)의 코드\n",
                      "example": "KRW-BTC"
                    },
                    "identifier": {
                      "type": "string",
                      "description": "주문 생성시 클라이언트가 지정한 주문 식별자.\n* identifier 필드는 2024년 10월 18일 이후 생성된 주문에 대해서만 제공됩니다.\n",
                      "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                    }
                  }
                }
              }
            }
          },
          "failed": {
            "description": "취소 실패한 주문 정보",
            "type": "object",
            "required": [
              "count",
              "orders"
            ],
            "properties": {
              "count": {
                "type": "integer",
                "description": "취소 실패한 주문 수",
                "example": 5
              },
              "orders": {
                "description": "취소 실패한 주문 목록",
                "type": "array",
                "items": {
                  "type": "object",
                  "required": [
                    "uuid",
                    "market"
                  ],
                  "properties": {
                    "uuid": {
                      "type": "string",
                      "description": "주문의 유일 식별자",
                      "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                    },
                    "market": {
                      "type": "string",
                      "description": "페어(거래쌍)의 코드\n",
                      "example": "KRW-BTC"
                    },
                    "identifier": {
                      "type": "string",
                      "description": "주문 생성시 클라이언트가 지정한 주문 식별자.\n* identifier 필드는 2024년 10월 18일 이후 생성된 주문에 대해서만 제공됩니다.\n",
                      "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                    }
                  }
                }
              }
            }
          }
        }
      },
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
      }
    }
  },
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
    "/v1/orders/open": {
      "delete": {
        "operationId": "batch-cancel-orders",
        "summary": "주문 일괄 취소 접수",
        "description": "조건을 지정하여 해당 조건을 만족하는 최대 300개의 주문을 일괄 취소합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request DELETE \\\n  --url 'https://api.upbit.com/v1/orders/open?cancel_side=all&excluded_pairs=KRW-ETH,BTC-XRP&quote_currencies=KRW,BTC' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/orders/open\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"cancel_side\": \"all\",\n  \"excluded_pairs\": \"KRW-ETH,BTC-XRP\",\n  \"quote_currencies\": \"KRW,BTC\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.delete(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/orders/open\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  cancel_side: \"all\",\n  excluded_pairs: \"KRW-ETH,BTC-XRP\",\n  quote_currencies: \"KRW,BTC\",\n};\n\nconst queryString = decodeURIComponent(new URLSearchParams(params).toString());\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"DELETE\",\n  url: `${baseURL}${path}?${queryString}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.net.URLEncoder;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class DeleteOrders {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/orders/open\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"cancel_side\", \"all\");\n        params.put(\"excluded_pairs\", \"KRW-ETH,BTC-XRP\");\n        params.put(\"quote_currencies\", \"KRW,BTC\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .delete()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "parameters": [
          {
            "name": "quote_currencies",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "주문 취소 대상 마켓 목록.\n디지털 자산 구매에 사용되는 통화(KRW,BTC,USDT)로 취소 대상 주문을 한정할 수 있는 필터 파라미터입니다.\n\n[예시] \"KRW\"로 지정시 원화마켓의 모든 미체결 주문 취소\n",
              "example": "KRW"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "KRW,BTC"
              }
            }
          },
          {
            "name": "cancel_side",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "취소하고자 하는 주문의 매수/매도 방향 구분.\n사용 가능한 값은 \"all\"(전체), \"ask\"(매도), \"bid\"(매수)입니다.\n",
              "enum": [
                "bid",
                "ask",
                "all"
              ],
              "default": "all"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "all"
              }
            }
          },
          {
            "name": "count",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "취소할 주문의 최대 개수.\n최대값은 300이며 미지정시 기본값은 20입니다.\n",
              "example": 300,
              "default": 20
            },
            "allowReserved": true
          },
          {
            "name": "order_by",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "enum": [
                "asc",
                "desc"
              ],
              "description": "취소할 주문의 정렬 방식. \"desc\"(최신 주문 순) 또는 \"asc\"(오래된 주문 순). 기본값은 \"desc\"입니다.",
              "example": "desc",
              "default": "desc"
            },
            "allowReserved": true
          },
          {
            "name": "pairs",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "페어(거래쌍) 목록.\n지정한 페어의 미체결 주문만 취소합니다. 최대 20개, 쉼표(,)로 구분합니다.\n\n[예시] pairs=KRW-BTC,KRW-ETH\n",
              "example": "KRW-BTC,KRW-ETH"
            },
            "allowReserved": true
          },
          {
            "name": "excluded_pairs",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "제외할 페어(거래쌍) 목록.\n지정한 페어의 주문을 제외한 모든 주문을 취소합니다. 최대 20개, 쉼표(,)로 구분합니다.\n\n[예시] excluded_pairs=KRW-BTC,KRW-ETH\n",
              "example": "KRW-BTC,KRW-ETH"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "KRW-ETH,BTC-XRP"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Object of batch cancellation result",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/BatchCancelResult"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "success": {
                        "count": 2,
                        "orders": [
                          {
                            "uuid": "bbbb8e07-1689-4769-af3e-a117016623f8",
                            "market": "KRW-ETH"
                          },
                          {
                            "uuid": "4312ba49-5f1a-4a01-9f3b-2d2bce17267e",
                            "market": "KRW-ETH"
                          }
                        ]
                      },
                      "failed": {
                        "count": 1,
                        "orders": [
                          {
                            "uuid": "bdb49a54-de36-4eb4-a963-9c8d4337a9da",
                            "market": "BTC-XRP"
                          }
                        ]
                      }
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
                  "body parameter not supported error": {
                    "value": {
                      "error": {
                        "name": "QUERY_PARAMETER_NOT_SUPPORTED",
                        "message": "이 API는 쿼리 파라미터 형식의 요청만 지원합니다. 요청 본문에 값을 포함할 수 없습니다."
                      }
                    }
                  },
                  "invalid parameter error": {
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
          "401": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "unauthorized jwt error": {
                    "value": {
                      "error": {
                        "name": "invalid_query_payload",
                        "message": "Failed to verify the query of Jwt."
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
                  "not found order error": {
                    "value": {
                      "error": {
                        "name": "order_not_found",
                        "message": "주문을 찾지 못했습니다."
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
  }
}
```

# Sibling pages

* [페어별 주문 가능 정보 조회](https://docs.upbit.com/kr/reference/available-order-information.md)
* [주문 생성](https://docs.upbit.com/kr/reference/new-order.md)
* [주문 생성 테스트](https://docs.upbit.com/kr/reference/order-test.md)
* [개별 주문 조회](https://docs.upbit.com/kr/reference/get-order.md)
* [id로 주문 목록 조회](https://docs.upbit.com/kr/reference/list-orders-by-ids.md)
* [체결 대기 주문 목록 조회](https://docs.upbit.com/kr/reference/list-open-orders.md)
* [종료 주문 목록 조회](https://docs.upbit.com/kr/reference/list-closed-orders.md)
* [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order.md)
* [id로 주문 목록 취소 접수](https://docs.upbit.com/kr/reference/cancel-orders-by-ids.md)
* [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order.md)