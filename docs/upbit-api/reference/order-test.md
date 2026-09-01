---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 주문 생성 테스트

실제 주문을 생성하지 않고 주문 요청 형식과 주문 가능 여부를 검증합니다.

주문 생성 테스트 API는 실제 주문을 생성하지 않고, 주문 요청 형식과 주문 가능 상태를 사전에 검증할 수 있는 API입니다.

* 실제 주문 API를 호출하기 전 이 API를 통해 거래 수수료 없이 주문 생성 가능 여부를 테스트하고, 발생 가능한 오류를 사전에 확인할 수 있습니다.
* 특히 **거래소 점검 직후 특정 페어의 주문 가능 상태를 확인하는 용도**로 활용할 수 있습니다.
  * 정상 응답이 반환된 경우, 주문 요청의 형식이 올바르고 해당 페어가 현재 주문 가능한 상태임을 의미합니다.
  * 만약 `market_offline` 오류가 반환된다면 해당 페어는 아직 주문이 불가능한 상태임을 의미합니다.
* 테스트 주문은 실제 주문과 동일한 검증 과정을 거치지만, 주문이 실제로 생성되지는 않습니다. **따라서 응답으로 반환된 UUID (혹은 identifier)는 주문 조회나 취소 요청에 사용할 수 없습니다.**

<br />

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
            <td className="code-col">v1.6.1</td>
            <td>2025-10-27</td>
            <td><a href="https://docs.upbit.com/kr/changelog/order-test">주문 테스트 기능 신규 지원</a></td>
      	</tr>


    </tbody>
</table>




    </div>
</div>

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 8회 호출할 수 있습니다. 포켓 단위로 측정되며 \[주문 테스트 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

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
      "SideEnum": {
        "type": "string",
        "enum": [
          "ask",
          "bid"
        ],
        "description": "주문 방향(매수/매도)"
      },
      "OrdTypeEnum": {
        "type": "string",
        "enum": [
          "limit",
          "price",
          "market",
          "best"
        ],
        "description": "주문 유형.\n\n* `limit`: 지정가 주문\n* `price`: 시장가 주문(매수)\n* `market`: 시장가 주문(매도)\n* `best`: 최유리 주문\n"
      },
      "TimeInForceEnum": {
        "type": "string",
        "enum": [
          "fok",
          "ioc",
          "post_only"
        ],
        "description": "주문 체결 조건.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: 메이커 전용 주문\n"
      },
      "SmpTypeEnum": {
        "type": "string",
        "enum": [
          "cancel_maker",
          "cancel_taker",
          "reduce"
        ],
        "description": "자전거래 체결 방지(SMP) 모드.\n\n* `cancel_maker`: 메이커 주문 취소\n* `cancel_taker`: 테이커 주문 취소\n* `reduce`: 수량 감소\n"
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
    "/v1/orders/test": {
      "post": {
        "operationId": "order-test",
        "summary": "주문 생성 테스트",
        "description": "실제 주문을 생성하지 않고 주문 요청 형식과 주문 가능 여부를 검증합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request POST \\\n--url 'https://api.upbit.com/v1/orders/test' \\\n--header 'Authorization: Bearer {JWT_TOKEN}' \\\n--header 'Accept: application/json' \\\n--header 'Content-Type: application/json' \\\n--data '\n  {\n  \"market\":\"KRW-BTC\",\n  \"side\":\"bid\",\n  \"volume\":\"1\",\n  \"price\":\"140000000\",\n  \"ord_type\":\"limit\"\n  }\n'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/orders/test\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n    \"market\": \"KRW-BTC\",\n    \"side\": \"bid\",\n    \"volume\": \"1\",\n    \"price\": \"140000000\",\n    \"ord_type\": \"limit\",\n}\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": query_hash,\n    \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n}\n\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/orders/test\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  market: \"KRW-BTC\",\n  side: \"bid\",\n  volume: \"1\",\n  price: \"140000000\",\n  ord_type: \"limit\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"POST\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n  data: params,\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\nimport com.google.gson.Gson;\n\npublic class CreateOrder {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/orders/test\";\n\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"market\", \"KRW-BTC\");\n        params.put(\"side\", \"bid\");\n        params.put(\"volume\", \"1\");\n        params.put(\"price\", \"140000000\");\n        params.put(\"ord_type\", \"limit\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        String jsonBody = new Gson().toJson(params);\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .post(RequestBody.create(jsonBody, okhttp3.MediaType.parse(\"application/json; charset=utf-8\")))\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "examples": {
                "default": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "bid",
                    "volume": "1",
                    "price": "14000000",
                    "ord_type": "limit"
                  }
                }
              },
              "schema": {
                "type": "object",
                "required": [
                  "market",
                  "side",
                  "ord_type"
                ],
                "properties": {
                  "market": {
                    "type": "string",
                    "description": "주문을 생성하고자 하는 대상 페어(거래쌍). (필수)",
                    "example": "KRW-BTC"
                  },
                  "side": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/SideEnum"
                      }
                    ],
                    "description": "주문 방향(매수/매도). (필수)"
                  },
                  "volume": {
                    "type": "string",
                    "format": "decimal",
                    "description": "주문 수량",
                    "example": "0.01"
                  },
                  "price": {
                    "type": "string",
                    "format": "decimal",
                    "description": "주문 단가 또는 총액",
                    "example": "1000"
                  },
                  "ord_type": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/OrdTypeEnum"
                      }
                    ],
                    "description": "주문 유형. (필수)\n\n* `limit`: 지정가 주문\n* `price`: 시장가 주문(매수)\n* `market`: 시장가 주문(매도)\n* `best`: 최유리 주문\n"
                  },
                  "identifier": {
                    "type": "string",
                    "description": "클라이언트 지정 주문 식별자",
                    "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                  },
                  "time_in_force": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/TimeInForceEnum"
                      }
                    ],
                    "description": "주문 체결 조건.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: 메이커 전용 주문\n"
                  },
                  "smp_type": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/SmpTypeEnum"
                      }
                    ],
                    "description": "자전거래 체결 방지(SMP) 모드.\n\n* `cancel_maker`: 메이커 주문 취소\n* `cancel_taker`: 테이커 주문 취소\n* `reduce`: 수량 감소\n"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Object of created order (dry-run)",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Order"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "market": "KRW-BTC",
                      "uuid": "cdd92199-2897-4e14-9448-f923320408ad",
                      "side": "ask",
                      "ord_type": "limit",
                      "price": "140000000",
                      "state": "wait",
                      "created_at": "2025-07-04T15:00:00+09:00",
                      "volume": "1.0",
                      "remaining_volume": "1.0",
                      "reserved_fee": "70000.0",
                      "remaining_fee": "70000.0",
                      "paid_fee": "0.0",
                      "locked": "0.0",
                      "executed_volume": "0.0",
                      "prevented_volume": "0",
                      "prevented_locked": "0",
                      "trades_count": 0
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
                  "invalid parameter error": {
                    "value": {
                      "error": {
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
                      }
                    }
                  },
                  "invalid time in force error": {
                    "value": {
                      "error": {
                        "name": "invalid_time_in_force",
                        "message": "주문조건을 다시 확인해 주세요."
                      }
                    }
                  },
                  "over krw funds bid error": {
                    "value": {
                      "error": {
                        "name": "over_krw_funds_bid",
                        "message": "최대 주문 가능 금액 1000000000 KRW 보다 작은 주문을 입력해 주세요."
                      }
                    }
                  },
                  "insufficient funds bid error": {
                    "value": {
                      "error": {
                        "name": "insufficient_funds_bid",
                        "message": "주문 가능한 금액(KRW)이 부족합니다."
                      }
                    }
                  },
                  "not found market error": {
                    "value": {
                      "error": {
                        "name": "notfoundmarket",
                        "message": "다음 마켓에 없는 종목입니다: KRW-BTCs"
                      }
                    }
                  }
                }
              }
            }
          },
          "403": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "market offline error": {
                    "value": {
                      "error": {
                        "name": "market_offline",
                        "message": "시스템 점검 중입니다. 점검이 완료된 후에 다시 시도해주세요."
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
* [개별 주문 조회](https://docs.upbit.com/kr/reference/get-order.md)
* [id로 주문 목록 조회](https://docs.upbit.com/kr/reference/list-orders-by-ids.md)
* [체결 대기 주문 목록 조회](https://docs.upbit.com/kr/reference/list-open-orders.md)
* [종료 주문 목록 조회](https://docs.upbit.com/kr/reference/list-closed-orders.md)
* [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order.md)
* [id로 주문 목록 취소 접수](https://docs.upbit.com/kr/reference/cancel-orders-by-ids.md)
* [주문 일괄 취소 접수](https://docs.upbit.com/kr/reference/batch-cancel-orders.md)
* [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order.md)