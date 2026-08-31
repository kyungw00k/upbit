---
updatedAt: 2026-07-30T02:37:44.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 포켓별 API Key 목록 조회

메인포켓의 API Key 권한을 사용하여 생성된 포켓별 API Key 목록과 정보를 조회하는 API입니다.

### 포켓별 API Key 조회 흐름

조회하고자하는 포켓 UUID는 <Anchor target="_blank" href="ref:list-pockets">포켓 정보 조회 API</Anchor>를 통해 확인할 수 있습니다.

```mermaid
flowchart LR
    A[포켓 정보 조회<br/>GET /v1/pockets]
    A --> B[포켓 UUID 확인]

    B --> C[포켓별 API Key 목록 조회<br/>GET /v1/pockets/api_keys]

    C --> D[포켓 정보]

    D --> E[UUID]
    D --> F[API Key 목록]

    F --> G[key #1]
    F --> H[key #2]
    F --> I[...]

    G --> J["access_key
permissions
allowed_ips
created_at
expired_at"]
```

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 \[Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>

<br />

<div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>

<div className="box-rate-limit">
  <a href="auth">인증</a>이 필요한 API로, 메인포켓만 사용 가능합니다.
  <br />

메인포켓 키발급 페이지에서 \[포켓관리] 권한이 설정된 API Key를 사용해야 합니다. <br />

권한 오류(out\_of\_scope) 오류가 발생한다면, <a href="https://upbit.com/mypage/open_api_management">API Key 관리 메뉴</a>에서 권한 설정을 확인해주세요.

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
      "PocketApiKeyPermission": {
        "type": "string",
        "enum": [
          "view_account",
          "view_orders",
          "make_orders",
          "view_withdrawals",
          "withdraw",
          "view_deposits",
          "deposit",
          "transfer",
          "manage_pockets"
        ],
        "description": "API Key에 부여된 접근 권한.\n\n* `view_account`: 자산/계좌 조회\n* `view_orders`: 주문 내역 조회\n* `make_orders`: 주문 생성 및 체결\n* `view_withdrawals`: 출금 내역 조회\n* `withdraw`: 디지털 자산/원화 출금\n* `view_deposits`: 입금 내역 조회\n* `deposit`: 입금 주소 생성 및 입금\n* `transfer`: 포켓 간 자산 이체\n* `manage_pockets`: 포켓 생성 및 관리\n",
        "example": "transfer"
      },
      "PocketApiKey": {
        "type": "object",
        "required": [
          "access_key",
          "permissions",
          "allowed_ips",
          "created_at",
          "expired_at"
        ],
        "properties": {
          "access_key": {
            "type": "string",
            "description": "API Key의 Access Key",
            "example": "8f3a2c1d9b7e6a5c4d3f2e1a0b9c8d7e6f5a4b3c"
          },
          "permissions": {
            "type": "array",
            "description": "API Key에 부여된 접근 권한 목록",
            "items": {
              "$ref": "#/components/schemas/PocketApiKeyPermission"
            },
            "example": [
              "deposit",
              "transfer",
              "manage_pockets"
            ]
          },
          "allowed_ips": {
            "type": "array",
            "description": "API 요청이 허용된 IP 주소 목록",
            "items": {
              "type": "string"
            },
            "example": {
              "allowed_ips": [
                "203.0.113.10",
                "198.51.100.25"
              ]
            }
          },
          "created_at": {
            "type": "string",
            "description": "API Key가 생성된 일시.\n\n[형식] yyyy-MM-dd'T'HH:mm:ss+09:00\n",
            "example": "2026-05-28T15:17:51+09:00"
          },
          "expired_at": {
            "type": "string",
            "description": "API Key의 만료 일시.\n\n[형식] yyyy-MM-dd'T'HH:mm:ss+09:00\n",
            "example": "2027-05-27T15:17:51+09:00"
          }
        }
      },
      "PocketApiKeyGroup": {
        "type": "object",
        "required": [
          "uuid",
          "keys"
        ],
        "properties": {
          "uuid": {
            "type": "string",
            "description": "포켓 UUID",
            "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
          },
          "keys": {
            "type": "array",
            "description": "해당 포켓에 발급된 API Key 목록",
            "items": {
              "$ref": "#/components/schemas/PocketApiKey"
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
    "/v1/pockets/api_keys": {
      "get": {
        "operationId": "list-pocket-api-keys",
        "summary": "포켓별 API Key 목록 조회",
        "description": "메인포켓의 API Key 권한을 사용하여 생성된 포켓별 API Key 목록과 정보를 조회하는 API입니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/pockets/api_keys?uuids[]=9ca023a5-851b-4fec-9f0a-48cd83c2eaae&include_expired=false' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/pockets/api_keys\"\nparams = [(\"uuids[]\", \"9ca023a5-851b-4fec-9f0a-48cd83c2eaae\"), (\"include_expired\", False)]\nheaders = {\n  \"Authorization\": \"Bearer {JWT_TOKEN}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios",
              "code": "const axios = require(\"axios\");\n\naxios.request({\n  method: \"GET\",\n  url: \"https://api.upbit.com/v1/pockets/api_keys\",\n  params: new URLSearchParams([\n    [\"uuids[]\", \"9ca023a5-851b-4fec-9f0a-48cd83c2eaae\"],\n    [\"include_expired\", \"false\"],\n  ]),\n  headers: {\n    Authorization: \"Bearer {JWT_TOKEN}\",\n    Accept: \"application/json\",\n  },\n}).then((response) => console.log(response.data));\n"
            },
            {
              "language": "java",
              "code": "import java.io.IOException;\nimport okhttp3.HttpUrl;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListPocketApiKeys {\n    public static void main(String[] args) throws IOException {\n        HttpUrl url = HttpUrl.parse(\"https://api.upbit.com/v1/pockets/api_keys\").newBuilder()\n            .addQueryParameter(\"uuids[]\", \"9ca023a5-851b-4fec-9f0a-48cd83c2eaae\")\n            .addQueryParameter(\"include_expired\", \"false\")\n            .build();\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(url)\n            .get()\n            .addHeader(\"Authorization\", \"Bearer {JWT_TOKEN}\")\n            .addHeader(\"Accept\", \"application/json\")\n            .build();\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.body().string());\n        }\n    }\n}\n"
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
              "description": "조회하고자 하는 포켓 UUID 배열.\n\n[예시] uuids[]=uuid1&uuids[]=uuid2\n",
              "items": {
                "type": "string"
              }
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": [
                  "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                ]
              }
            }
          },
          {
            "name": "include_expired",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean",
              "description": "만료된 API Key를 조회 결과에 포함할지 여부",
              "default": false
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": false
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "List of pocket API keys",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/PocketApiKeyGroup"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "uuid": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
                        "keys": [
                          {
                            "access_key": "8f3a2c1d9b7e6a5c4d3f2e1a0b9c8d7e6f5a4b3c",
                            "permissions": [
                              "deposit",
                              "transfer",
                              "manage_pockets"
                            ],
                            "allowed_ips": [
                              "203.0.113.10",
                              "198.51.100.25"
                            ],
                            "created_at": "2026-05-28T15:17:51+09:00",
                            "expired_at": "2027-05-27T15:17:51+09:00"
                          }
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
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터 입니다."
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
                  "out of scope error": {
                    "value": {
                      "error": {
                        "name": "out_of_scope",
                        "message": "권한이 부족합니다."
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

* [포켓 정보 조회](https://docs.upbit.com/kr/reference/list-pockets.md)
* [서브포켓 잔고 조회](https://docs.upbit.com/kr/reference/get-sub-pocket-balance.md)
* [메인포켓 자산 이전](https://docs.upbit.com/kr/reference/universal-transfer.md)
* [메인포켓 자산 이전 목록 조회](https://docs.upbit.com/kr/reference/list-universal-transfers.md)
* [서브포켓 자산 이전](https://docs.upbit.com/kr/reference/transfer.md)
* [서브포켓 자산 이전 목록 조회](https://docs.upbit.com/kr/reference/list-transfers.md)