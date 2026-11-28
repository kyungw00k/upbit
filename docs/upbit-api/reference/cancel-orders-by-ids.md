---
updatedAt: 2026-11-25T17:47:26.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# id로 주문 목록 취소 접수

UUID 또는 Identifier 목록으로 취소 대상 주문을 지정 취소합니다. 한 번의 요청으로 최대 20개의 주문을 취소할 수 있습니다.

<Callout icon="far fa-octagon-exclamation" theme="error">
  #### 주문 취소 요청이 거절되는 경우

  다음과 같은 사유로 주문 취소 요청이 거절될 수 있습니다.

  * 이미 전량 체결 완료되어 취소가 불가한 주문
  * 이미 취소가 완료된 주문
  * 해당 주문의 페어가 리브랜딩 등의 이유로 서비스를 일시 정지한 경우
</Callout>

<Callout icon="far fa-circle-info" theme="info">
  #### **이 API는 쿼리 파라미터 형식만 지원합니다.**

  모든 요청 파라미터는 쿼리 파라미터 형식으로만 요청할 수 있습니다. Request Body를 통한 form 형식 또는 json 형식의 파라미터 요청을 지원하지 않습니다.

  \[예시] /v1/orders/uuids?uuids\[]=1234567890\&uuids\[]=1234567891
</Callout>

<Callout icon="far fa-circle-info" theme="info">
  #### 취소 시 uuid 또는 identifier 중 하나는 반드시 포함해야 합니다.

  두 파라미터 모두 선택(Optional) 파라미터이지만, 취소하고자 하는 주문 지정을 위해 반드시 하나의 파라미터는 포함해야 합니다. uuid와 identifier를 모두 사용하는 경우, uuid를 기준으로 취소됩니다.
</Callout>

***

#### Rate Limit

초당 최대 30회 호출할 수 있습니다.

포켓 단위로 측정되며, \[<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits#:~:text=%ED%98%B8%EA%B0%80%20%EC%A0%95%EC%B1%85%20%EC%A1%B0%ED%9A%8C-,Exchange%20REST%20API,-Rate%20Limit%20%EA%B7%B8%EB%A3%B9">`default` 그룹</Anchor>]에 속한 API는 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보를 위해 제한 또는 조정될 수 있습니다.

#### **API Key Permission**<br />

<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/auth">인증</Anchor>이 필요한 API로, \[주문하기] 권한이 설정된 API Key를 사용해야 합니다.

권한 오류(out\_of\_scope)가 발생하면[ API Key 관리 메뉴](https://www.upbit.com/mypage/open_api_management?tab=api_key)에서 권한 설정을 확인해 주세요.

# OpenAPI 정의

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "EXCHANGE API",
    "version": "1.0.1"
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
      "node",
      "java"
    ]
  },
  "paths": {
    "/v1/orders/uuids": {
      "delete": {
        "operationId": "cancel-orders-by-ids",
        "summary": "id로 주문 목록 취소 접수",
        "description": "UUID 또는 Identifier 목록으로 취소 대상 주문을 지정 취소합니다. 한 번의 요청으로 최대 20개의 주문을 취소할 수 있습니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request DELETE \\\n    --url 'https://api.upbit.com/v1/orders/uuids?uuids[]=bbbb8e07-1689-4769-af3e-a117016623f8&uuids[]=4312ba49-5f1a-4a01-9f3b-2d2bce17267e&uuids[]=bdb49a54-de36-4eb4-a963-9c8d4337a9da' \\\n    --header 'Authorization: Bearer {JWT_TOKEN}' \\\n    --header 'accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/orders/uuids\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n    \"uuids[]\": [\"bbbb8e07-1689-4769-af3e-a117016623f8\", \"4312ba49-5f1a-4a01-9f3b-2d2bce17267e\", \"bdb49a54-de36-4eb4-a963-9c8d4337a9da\"],\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": query_hash,\n    \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n}\n\nres = requests.delete(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/orders/uuids\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst arrayParams = {\n  \"uuids[]\": [\n  \"bbbb8e07-1689-4769-af3e-a117016623f8\",\n  \"4312ba49-5f1a-4a01-9f3b-2d2bce17267e\",\n  \"bdb49a54-de36-4eb4-a963-9c8d4337a9da\",\n  ],\n};\n\nconst arrayParamsString = Object.entries(arrayParams)\n  .flatMap(([key, values]) => values.map((value) => `${key}=${value}`))\n  .join(\"&\");\n\nconst queryString = arrayParamsString;\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"DELETE\",\n  url: `${baseURL}${path}?${queryString}`,\n  headers: {\n  Authorization: `Bearer ${jwtToken}`,\n  Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n  console.log(response.data);\n  })\n  .catch((error) => {\n  console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.List;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class DeleteOrdersById {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/orders/uuids\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, List<String>> params = new HashMap<>();\n        params.put(\"uuids[]\", List.of(\"bbbb8e07-1689-4769-af3e-a117016623f8\", \"4312ba49-5f1a-4a01-9f3b-2d2bce17267e\", \"bdb49a54-de36-4eb4-a963-9c8d4337a9da\"));\n\n        String queryString = params.entrySet().stream()\n            .flatMap(e -> e.getValue().stream().map(v -> e.getKey() + \"=\" + v))\n            .collect(Collectors.joining(\"&\"));\n\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "parameters": [
          {
            "name": "uuids[]",
            "in": "query",
            "required": false,
            "schema": {
              "type": "array",
              "description": "취소하고자 하는 주문의 UUID 목록.\n\n취소 가능한 최대 주문 개수는 20개입니다. 2개 이상의 UUID로 조회하는 경우 쿼리 파라미터를 다음과 같이 요청합니다.\n\n[예시] uuids[]=uuid1&uuids[]=uuid2…\n",
              "items": {
                "type": "string",
                "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
              }
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": [
                  "bbbb8e07-1689-4769-af3e-a117016623f8",
                  "4312ba49-5f1a-4a01-9f3b-2d2bce17267e",
                  "bdb49a54-de36-4eb4-a963-9c8d4337a9da"
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
              "description": "취소하고자 하는 주문의 클라이언트 지정 식별자 목록.\n\n취소 가능한 최대 주문 개수는 20개입니다. 2개 이상의 identifiers로 조회하는 경우 쿼리 파라미터를 다음과 같이 요청합니다.\n\n[예시] identifiers[]=id1&identifiers[]=id2…\n",
              "items": {
                "type": "string",
                "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
              }
            },
            "allowReserved": true
          }
        ],
        "responses": {
          "200": {
            "description": "Object of batch cancel result",
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
                  "QUERY_PARAMETER_NOT_SUPPORTED": {
                    "value": {
                      "error": {
                        "name": "QUERY_PARAMETER_NOT_SUPPORTED",
                        "message": "이 API는 쿼리 파라미터 형식의 요청만 지원합니다. 요청 본문에 값을 포함할 수 없습니다."
                      }
                    }
                  },
                  "bad_request": {
                    "value": {
                      "error": {
                        "name": "bad_request",
                        "message": "잘못된 요청"
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

# 같은 수준의 페이지

* [페어별 주문 가능 정보 조회](https://docs.upbit.com/kr/reference/available-order-information.md)
* [주문 생성](https://docs.upbit.com/kr/reference/new-order.md)
* [주문 생성 테스트](https://docs.upbit.com/kr/reference/order-test.md)
* [개별 주문 조회](https://docs.upbit.com/kr/reference/get-order.md)
* [id로 주문 목록 조회](https://docs.upbit.com/kr/reference/list-orders-by-ids.md)
* [체결 대기 주문 목록 조회](https://docs.upbit.com/kr/reference/list-open-orders.md)
* [종료 주문 목록 조회](https://docs.upbit.com/kr/reference/list-closed-orders.md)
* [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order.md)
* [주문 일괄 취소 접수](https://docs.upbit.com/kr/reference/batch-cancel-orders.md)
* [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order.md)