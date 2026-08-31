---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 종료 주문 목록 조회

종료 주문(Closed Order) 목록을 조회합니다.

포켓의 종료 주문 목록을 조회할 수 있습니다. 종료 주문은 전량 체결 주문과 취소 주문을 포함합니다. 필터 파라미터인 페어와 주문 상태(체결 완료 또는 주문 취소), 조회 기간을 선택적으로 입력하여 해당 조건을 만족하는 종료 주문들을 조회할 수 있습니다. 조회 기간을 정하는 경우 최대 7일 구간(window)을 조회할 수 있습니다.

<div className="callout-section">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  state와 states[]는 동시에 사용할 수 없습니다.
      </div>
      주문 상태 필터링 시 조건에 따라 적절한 파라미터를 하나만 사용해주세요.
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
            <td className="code-col">v1.5.8</td>
            <td className="code-col">2025-07-02</td>
            <td><a href="https://docs.upbit.com/kr/changelog/smp"> 자전거래 체결 방지(SMP) 기능<br />신규 지원에 따른 필드 추가<br /><code>smp_type</code>,<code>prevented_volume</code>,<code>prevented_locked</code></a></td>
      	</tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-12-18</td>
            <td><a href="https://docs.upbit.com/kr/changelog/closed_order_support_timestamp"> <code>start_time</code>, <code>end_time</code> 파라미터 timestamp 형식 추가 지원</a></td>
        </tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-12-04</td>
            <td><a href="https://docs.upbit.com/kr/changelog/myorder_identifier"> <code>identifier</code> 필드 신규 지원</a></td>
     	  </tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-10-02</td>
            <td><a href="https://docs.upbit.com/kr/changelog/extend_time_window"> 주문 조회 가능 범위 7일로 확대 지원</a></td>
        </tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-06-27</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_get_orders"> '종료 주문 목록 조회' 신규 기능 지원</a></td>
        </tr>
    </tbody>
</table>




  </div>

</div>

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 \[Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>

<br />
  <div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>
  <div className="box-rate-limit">
    <a href="auth">인증</a>이 필요한 API로, [주문조회] 권한이 설정된 API Key를 사용해야 합니다. <br />
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
      "ClosedOrderState": {
        "type": "string",
        "enum": [
          "done",
          "cancel"
        ],
        "description": "완료 주문 상태.\n\n* `done`: 완전 체결\n* `cancel`: 취소\n"
      },
      "Order": {
        "type": "object",
        "required": [
          "market",
          "uuid",
          "side",
          "ord_type",
          "state",
          "created_at",
          "remaining_volume",
          "executed_volume",
          "reserved_fee",
          "remaining_fee",
          "paid_fee",
          "locked",
          "trades_count",
          "prevented_volume",
          "prevented_locked"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n",
            "example": "KRW-BTC"
          },
          "uuid": {
            "type": "string",
            "description": "주문의 유일 식별자",
            "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
          },
          "side": {
            "type": "string",
            "enum": [
              "ask",
              "bid"
            ],
            "description": "주문 방향(매수/매도)",
            "example": "ask"
          },
          "ord_type": {
            "type": "string",
            "enum": [
              "limit",
              "price",
              "market",
              "best"
            ],
            "description": "주문 유형.",
            "example": "limit"
          },
          "price": {
            "type": "string",
            "format": "decimal",
            "description": "주문 단가 또는 총액\n\n지정가 주문의 경우 단가, 시장가 매수 주문의 경우 매수 총액입니다.\n",
            "example": 1000
          },
          "state": {
            "type": "string",
            "enum": [
              "wait",
              "watch",
              "done",
              "cancel"
            ],
            "description": "주문 상태\n\n- `wait`: 체결 대기\n- `watch`: 예약 주문 대기\n- `done`: 체결 완료\n- `cancel`: 주문 취소\n",
            "example": "wait"
          },
          "created_at": {
            "type": "string",
            "description": "주문 생성 시각 (KST 기준)\n\n[형식] yyyy-MM-ddTHH:mm:ss+09:00\n",
            "example": "2025-06-25T15:42:25+09:00"
          },
          "volume": {
            "type": "string",
            "format": "decimal",
            "description": "주문 요청 수량",
            "example": 10
          },
          "remaining_volume": {
            "type": "string",
            "format": "decimal",
            "description": "체결 후 남은 주문 양",
            "example": 8
          },
          "executed_volume": {
            "type": "string",
            "format": "decimal",
            "description": "체결된 양",
            "example": 2
          },
          "reserved_fee": {
            "type": "string",
            "format": "decimal",
            "description": "수수료로 예약된 비용",
            "example": 5
          },
          "remaining_fee": {
            "type": "string",
            "format": "decimal",
            "description": "남은 수수료",
            "example": 5
          },
          "paid_fee": {
            "type": "string",
            "format": "decimal",
            "description": "사용된 수수료",
            "example": 0
          },
          "locked": {
            "type": "string",
            "format": "decimal",
            "description": "거래에 사용 중인 비용",
            "example": 0
          },
          "trades_count": {
            "type": "integer",
            "description": "해당 주문에 대한 체결 건수",
            "example": 1
          },
          "time_in_force": {
            "type": "string",
            "enum": [
              "fok",
              "ioc",
              "post_only"
            ],
            "description": "주문 체결 옵션",
            "example": "ioc"
          },
          "identifier": {
            "type": "string",
            "description": "주문 생성시 클라이언트가 지정한 주문 식별자.\n* identifier 필드는 2024년 10월 18일 이후 생성된 주문에 대해서만 제공됩니다.\n",
            "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
          },
          "smp_type": {
            "type": "string",
            "enum": [
              "reduce",
              "cancel_maker",
              "cancel_taker"
            ],
            "description": "자전거래 체결 방지(Self-Match Prevention) 모드",
            "example": "cancel_maker"
          },
          "prevented_volume": {
            "type": "string",
            "format": "decimal",
            "description": "자전거래 방지로 인해 취소된 수량.\n\n동일 사용자의 주문 간 체결이 발생하지 않도록 설정(SMP)에 따라 취소된 주문 수량입니다.\n",
            "example": 2
          },
          "prevented_locked": {
            "type": "string",
            "format": "decimal",
            "description": "자전거래 방지로 인해 해제된 자산.\n자전거래 체결 방지 설정으로 인해 취소된 주문의 잔여 자산입니다.\n  - 매수 주문의 경우: 취소된 금액\n  - 매도 주문의 경우: 취소된 수량\n",
            "example": 2000
          },
          "trades": {
            "type": "array",
            "description": "체결 목록 (개별 주문 조회 시 반환)",
            "items": {
              "type": "object",
              "properties": {
                "market": {
                  "type": "string",
                  "description": "페어(거래쌍)의 코드",
                  "example": "KRW-BTC"
                },
                "uuid": {
                  "type": "string",
                  "description": "체결의 유일 식별자",
                  "example": "795dff29-bba6-49b2-baab-63473ab7931c"
                },
                "price": {
                  "type": "string",
                  "format": "decimal",
                  "description": "체결 단가",
                  "example": "1375"
                },
                "volume": {
                  "type": "string",
                  "format": "decimal",
                  "description": "체결 수량",
                  "example": "5.377594"
                },
                "funds": {
                  "type": "string",
                  "format": "decimal",
                  "description": "체결 금액",
                  "example": "7394.19175"
                },
                "trend": {
                  "type": "string",
                  "description": "체결 방향 (up/down)",
                  "example": "down"
                },
                "created_at": {
                  "type": "string",
                  "description": "체결 시각",
                  "example": "2025-08-09T16:44:00.597751+09:00"
                },
                "side": {
                  "type": "string",
                  "description": "주문 방향",
                  "example": "ask"
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
    "/v1/orders/closed": {
      "get": {
        "operationId": "list-closed-orders",
        "summary": "종료 주문 목록 조회",
        "description": "종료 주문(Closed Order) 목록을 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/orders/closed?market=KRW-BTC&states[]=done&states[]=cancel' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/orders/closed\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\"market\": \"KRW-BTC\", \"states[]\": [\"done\", \"cancel\"]}\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/orders/closed\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  market: \"KRW-BTC\",\n  \"states[]\": [\"done\", \"cancel\"],\n};\n\nconst encodedQueryString = new URLSearchParams(\n  Object.entries(params).flatMap(([key, value]) =>\n    Array.isArray(value) ? value.map((v) => [key, v]) : [[key, value]]\n  )\n).toString();\n\nconst queryString = decodeURIComponent(encodedQueryString);\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\nconsole.log(jwtToken);\n\nconst options = {\n  method: \"GET\",\n  url: `${BASE_URL}${PATH}?${encodedQueryString}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.List;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ClosedOrders {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/orders/closed\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, List<String>> params = new HashMap<>();\n        params.put(\"market\", \"KRW-BTC\");\n        params.put(\"states[]\", List.of(\"done\", \"cancel\"));\n        String queryString = params.entrySet().stream()\n            .flatMap(e -> e.getValue().stream().map(v -> e.getKey() + \"=\" + v))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "parameters": [
          {
            "name": "market",
            "in": "query",
            "required": false,
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
            "name": "state",
            "in": "query",
            "required": false,
            "schema": {
              "allOf": [
                {
                  "$ref": "#/components/schemas/ClosedOrderState"
                }
              ],
              "description": "주문의 상태.\n\n- 시장가 매수 주문은 체결 후 `cancel`, `done` 두 상태 모두 발생할 수 있습니다.\n- 체결 후 잔량이 발생하는 경우 `cancel`, 딱 맞아떨어지는 경우 `done`이 됩니다.\n"
            },
            "allowReserved": true
          },
          {
            "name": "states[]",
            "in": "query",
            "required": false,
            "schema": {
              "type": "array",
              "description": "주문의 상태 (배열 형식).\n사용 가능한 값은 \"done\"(전체 체결 완료), \"cancel\"(전체 또는 부분 취소)입니다. 미지정시 기본값은 모든 상태(done, cancel)입니다.\n\n[예시] states[]=done&states[]=cancel\n",
              "items": {
                "allOf": [
                  {
                    "$ref": "#/components/schemas/ClosedOrderState"
                  }
                ]
              },
              "default": [
                "done",
                "cancel"
              ]
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": [
                  "done",
                  "cancel"
                ]
              }
            }
          },
          {
            "name": "start_time",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회 기간의 시작 시각.\n지정한 시각으로부터 \"end_time\"까지 생성된 주문을 조회합니다. 최대 조회 범위는 7일입니다.\n\n* \"start_time\"만 입력하는 경우 해당 시각 기준으로 이후 7일을 조회합니다.\n* 두 필드 모두 미입력시 요청 시각 기준 이전 7일을 조회합니다.\n\n형식: ISO 8601 (2025-06-24T13:56:53+09:00) 또는 밀리초 타임스탬프 (1750741013000)\n",
              "example": "2024-12-09T13:56:53+09:00"
            }
          },
          {
            "name": "end_time",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회 기간의 종료 시각.\n\"start_time\"부터 이 시각까지 생성된 주문을 조회합니다. 최대 조회 범위는 7일입니다.\n\n* \"end_time\"만 입력하는 경우 해당 시각 기준으로 이전 7일을 조회합니다.\n* 두 필드 모두 미입력시 요청 시각 기준 이전 7일을 조회합니다.\n\n형식: ISO 8601 (2025-06-24T13:56:53+09:00) 또는 밀리초 타임스탬프 (1750741013000)\n",
              "example": "2024-12-09T13:56:53+09:00"
            }
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "요청 개수(default: 100, max: 1,000)",
              "example": 1000,
              "default": 100
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
              "description": "결과 정렬 방식. \"desc\"(최신 주문 순) 또는 \"asc\"(오래된 주문 순). 기본값은 \"desc\"입니다.",
              "example": "desc",
              "default": "desc"
            },
            "allowReserved": true
          }
        ],
        "responses": {
          "200": {
            "description": "List of closed orders",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Order"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "uuid": "d098ceaf-6811-4df8-97f2-b7e01aefc03f",
                        "side": "bid",
                        "ord_type": "limit",
                        "price": "140000000",
                        "state": "done",
                        "market": "KRW-BTC",
                        "created_at": "2025-07-04T15:00:00+09:00",
                        "volume": "1.0",
                        "remaining_volume": "0.0",
                        "executed_volume": "1.0",
                        "reserved_fee": "0.0",
                        "remaining_fee": "0.0",
                        "paid_fee": "70000.0",
                        "locked": "0.0",
                        "prevented_volume": "0",
                        "prevented_locked": "0",
                        "trades_count": 1
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
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
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
* [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order.md)
* [id로 주문 목록 취소 접수](https://docs.upbit.com/kr/reference/cancel-orders-by-ids.md)
* [주문 일괄 취소 접수](https://docs.upbit.com/kr/reference/batch-cancel-orders.md)
* [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order.md)