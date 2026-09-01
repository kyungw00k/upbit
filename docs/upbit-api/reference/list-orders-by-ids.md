---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# id로 주문 목록 조회

UUID 또는 Identifier 목록으로 주문을 조회합니다.

조회하고자 하는 주문들의 UUID 또는 Identifier를 알고 있는 경우 이 API를 호출하여 주문 목록을 조회할 수 있습니다. 필터 파라미터인 페어를 선택적으로 입력하여 주문 중 해당 페어에 해당하는 주문 목록만 조회할 수 있습니다.

주문 식별자 없이 특정 페어의 전체 주문 내역을 조회하고자 하시는 경우, 또는 시간을 조회 기준으로 적용해야 하는 경우 [체결 대기 주문 목록 조회](https://docs.upbit.com/kr/reference/list-open-orders) 또는 [종료 주문 목록 조회](https://docs.upbit.com/kr/reference/list-closed-orders) API를 사용할 수 있습니다.

<div className="callout-section">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  조회 요청 시 uuid[] 또는 identifier[]를 반드시 포함해야 합니다.
      </div>
    두 파라미터 모두 선택(Optional) 파라미터이지만, 조회하고자 하는 주문 지정을 위해 반드시 하나의 파라미터를 포함해야 합니다. 두 파라미터를 동시에 사용할 수 없습니다.
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
            <td>2025-07-02</td>
            <td><a href="https://docs.upbit.com/kr/changelog/smp"> 자전거래 체결 방지(SMP) 기능<br />신규 지원에 따른 필드 추가<br /><code>smp_type</code>,<code>prevented_volume</code>,<code>prevented_locked</code></a></td>
      	</tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-12-04</td>
            <td><a href="https://docs.upbit.com/kr/changelog/myorder_identifier"> <code>identifier</code> 필드 신규 지원</a></td>
      	</tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-06-27</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_get_orders"> '지정 주문 목록 조회' 기능 신규 지원</a></td>
        </tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-04-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_ord_type_expand"> 최유리지정가 주문 유형 신규 지원<br />주문 옵션(time_in_force) 추가 지원</a></td>
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
    "/v1/orders/uuids": {
      "get": {
        "operationId": "list-orders-by-ids",
        "summary": "id로 주문 목록 조회",
        "description": "UUID 또는 Identifier 목록으로 주문을 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/orders/uuids?uuids[]=5d303952-8be9-41e6-915b-121a90026248&uuids[]=3944c2c1-bd8c-441a-aa25-2370d08217a9&uuids[]=3b67e543-8ad3-48d0-8451-0dad315cae73' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/orders/uuids\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n    \"uuids[]\": [\"5d303952-8be9-41e6-915b-121a90026248\",\"3944c2c1-bd8c-441a-aa25-2370d08217a9\",\"5b95451b-971e-4e76-8f61-5ff441f078d5\",\"3b67e543-8ad3-48d0-8451-0dad315cae73\"],\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/orders/uuids\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst arrayParams = {\n  \"uuids[]\": [\n    \"5d303952-8be9-41e6-915b-121a90026248\",\n    \"3944c2c1-bd8c-441a-aa25-2370d08217a9\",\n    \"5b95451b-971e-4e76-8f61-5ff441f078d5\",\n    \"3b67e543-8ad3-48d0-8451-0dad315cae73\",\n  ],\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}?${queryString}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListOrders {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/orders/uuids\";\n\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, List<String>> params = new HashMap<>();\n        params.put(\"uuids[]\", List.of(\"5d303952-8be9-41e6-915b-121a90026248\",\"3944c2c1-bd8c-441a-aa25-2370d08217a9\",\"5b95451b-971e-4e76-8f61-5ff441f078d5\",\"3b67e543-8ad3-48d0-8451-0dad315cae73\"));\n\n        String queryString = params.entrySet().stream()\n            .flatMap(e -> e.getValue().stream().map(v -> e.getKey() + \"=\" + v))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
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
            "allowReserved": true
          },
          {
            "name": "uuids[]",
            "in": "query",
            "required": false,
            "schema": {
              "type": "array",
              "description": "조회하고자 하는 주문의 UUID 목록.\n\n조회 가능한 최대 주문 개수는 100개입니다. 2개 이상의 UUID로 조회하는 경우 쿼리 파라미터를 다음과 같이 요청합니다.\n\n[예시] uuids[]=uuid1&uuids[]=uuid2…\n",
              "items": {
                "type": "string",
                "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
              }
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": [
                  "5d303952-8be9-41e6-915b-121a90026248",
                  "3944c2c1-bd8c-441a-aa25-2370d08217a9",
                  "5b95451b-971e-4e76-8f61-5ff441f078d5",
                  "3b67e543-8ad3-48d0-8451-0dad315cae73"
                ]
              }
            }
          },
          {
            "name": "identifiers[]",
            "in": "query",
            "required": false,
            "schema": {
              "type": "array",
              "description": "조회하고자 하는 주문의 클라이언트 지정 식별자 목록.\n\n조회 가능한 최대 주문 개수는 100개입니다. 2개 이상의 identifiers로 조회하는 경우 쿼리 파라미터를 다음과 같이 요청합니다.\n\n[예시] identifiers[]=id1&identifiers[]=id2…\n",
              "items": {
                "type": "string",
                "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
              }
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
              "description": "결과 정렬 방식.\n주문 생성 시각을 기준으로 \"desc\"(내림차순, 최신 주문 순) 또는 \"asc\"(오름차순, 오래된 주문 순)으로 정렬합니다. 기본값은 \"desc\"입니다.\n",
              "example": "desc",
              "default": "desc"
            },
            "allowReserved": true
          }
        ],
        "responses": {
          "200": {
            "description": "List of orders",
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
                        "uuid": "5d303952-8be9-41e6-915b-121a90026248",
                        "side": "bid",
                        "ord_type": "limit",
                        "price": "155772000",
                        "state": "wait",
                        "market": "KRW-BTC",
                        "created_at": "2025-08-09T16:52:49+09:00",
                        "volume": "0.0001",
                        "remaining_volume": "0.0001",
                        "prevented_volume": "0",
                        "reserved_fee": "7.7886",
                        "remaining_fee": "7.7886",
                        "paid_fee": "0",
                        "locked": "15584.9886",
                        "prevented_locked": "0",
                        "executed_volume": "0",
                        "trades_count": 0
                      },
                      {
                        "uuid": "3b67e543-8ad3-48d0-8451-0dad315cae73",
                        "side": "ask",
                        "ord_type": "market",
                        "state": "done",
                        "market": "KRW-USDT",
                        "created_at": "2025-08-09T16:44:00+09:00",
                        "volume": "5.377594",
                        "remaining_volume": "0",
                        "prevented_volume": "0",
                        "reserved_fee": "0",
                        "remaining_fee": "0",
                        "paid_fee": "3.697095875",
                        "locked": "0",
                        "prevented_locked": "0",
                        "executed_volume": "5.377594",
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
* [체결 대기 주문 목록 조회](https://docs.upbit.com/kr/reference/list-open-orders.md)
* [종료 주문 목록 조회](https://docs.upbit.com/kr/reference/list-closed-orders.md)
* [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order.md)
* [id로 주문 목록 취소 접수](https://docs.upbit.com/kr/reference/cancel-orders-by-ids.md)
* [주문 일괄 취소 접수](https://docs.upbit.com/kr/reference/batch-cancel-orders.md)
* [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order.md)