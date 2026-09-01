---
updatedAt: 2026-07-30T02:37:44.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 포켓 정보 조회

메인포켓 및 사용자 생성 포켓의 UUID와 이름을 조회합니다.

사용자가 보유한 메인 포켓 및 사용자 생성 포켓의 UUID와 이름 정보를 조회하는 API입니다. 조회된 포켓 UUID는 포켓별 API Key 목록 조회, 서브포켓 잔고 조회, 메인포켓·서브포켓 간 자산 이전 및 자산 이전 내역 조회 등의 기능을 수행하기 위한 식별자로 사용됩니다.

<br />

### 포켓 정보 활용 구조

```mermaid
flowchart LR
    A[포켓 정보 조회<br/>GET /v1/pockets]

    A --> B[메인포켓 UUID / 이름]
    A --> C[사용자 생성 포켓 UUID / 이름]

    B --> D[포켓별 API Key 목록 조회<br/>GET /v1/pockets/api_keys]
    B --> E[서브포켓 잔고 조회<br/>GET /v1/pockets/asset]
    B --> F[메인포켓 자산 이전<br/>POST /v1/pockets/transfers]
    B --> G[메인포켓 자산 이전 목록 조회<br/>GET /v1/pockets/transfers]

    C --> H[서브포켓 자산 이전<br/>POST /v1/pockets/transfers]
    C --> I[서브포켓 자산 이전 목록 조회<br/>GET /v1/pockets/transfers]

    style A fill:#eef,stroke:#446
    style B fill:#efe,stroke:#4a4
    style C fill:#efe,stroke:#4a4
```

<br />

<br />

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 30회 호출할 수 있습니다. 포켓 단위로 측정되며 \[Exchange 기본 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>

<br />

<div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>

<div className="box-rate-limit">
  <a href="auth">인증</a>이 필요한 API로, 메인포켓만 사용 가능합니다.
  <br />

메인포켓 키발급 페이지에서 \[포켓관리] 권한이 설정된 API Key를 사용해야 합니다.<br />

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
      "PocketType": {
        "type": "string",
        "enum": [
          "main",
          "user_spot_trading"
        ],
        "description": "포켓 종류.\n\n* `main`: 메인 포켓\n* `user_spot_trading`: 사용자 생성 포켓\n",
        "example": "user_spot_trading"
      },
      "Pocket": {
        "type": "object",
        "required": [
          "uuid",
          "name",
          "type"
        ],
        "properties": {
          "uuid": {
            "type": "string",
            "description": "포켓 UUID",
            "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
          },
          "name": {
            "type": "string",
            "description": "포켓 이름",
            "example": "단기 트레이딩 계좌"
          },
          "type": {
            "$ref": "#/components/schemas/PocketType"
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
    "/v1/pockets": {
      "get": {
        "operationId": "list-pockets",
        "summary": "포켓 정보 조회",
        "description": "메인포켓 및 사용자 생성 포켓의 UUID와 이름을 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/pockets' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/pockets\"\nheaders = {\n  \"Authorization\": \"Bearer {JWT_TOKEN}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios",
              "code": "const axios = require(\"axios\");\n\naxios.request({\n  method: \"GET\",\n  url: \"https://api.upbit.com/v1/pockets\",\n  headers: {\n    Authorization: \"Bearer {JWT_TOKEN}\",\n    Accept: \"application/json\",\n  },\n}).then((response) => console.log(response.data));\n"
            },
            {
              "language": "java",
              "code": "import java.io.IOException;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListPockets {\n    public static void main(String[] args) throws IOException {\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(\"https://api.upbit.com/v1/pockets\")\n            .get()\n            .addHeader(\"Authorization\", \"Bearer {JWT_TOKEN}\")\n            .addHeader(\"Accept\", \"application/json\")\n            .build();\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.body().string());\n        }\n    }\n}\n"
            }
          ]
        },
        "responses": {
          "200": {
            "description": "List of pockets",
            "content": {
              "application/json": {
                "schema": {
                  "type": "array",
                  "items": {
                    "$ref": "#/components/schemas/Pocket"
                  }
                },
                "examples": {
                  "Successful Example": {
                    "value": [
                      {
                        "uuid": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
                        "name": "메인포켓",
                        "type": "main"
                      },
                      {
                        "uuid": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae",
                        "name": "단기 트레이딩 계좌",
                        "type": "user_spot_trading"
                      }
                    ]
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

* [포켓별 API Key 목록 조회](https://docs.upbit.com/kr/reference/list-pocket-api-keys.md)
* [서브포켓 잔고 조회](https://docs.upbit.com/kr/reference/get-sub-pocket-balance.md)
* [메인포켓 자산 이전](https://docs.upbit.com/kr/reference/universal-transfer.md)
* [메인포켓 자산 이전 목록 조회](https://docs.upbit.com/kr/reference/list-universal-transfers.md)
* [서브포켓 자산 이전](https://docs.upbit.com/kr/reference/transfer.md)
* [서브포켓 자산 이전 목록 조회](https://docs.upbit.com/kr/reference/list-transfers.md)