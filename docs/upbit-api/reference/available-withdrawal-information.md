---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 출금 가능 정보 조회

지정한 통화에 대한 출금 가능 정보를 조회합니다. 해당 통화의 출금 정책과 사용자 잔고를 확인할 수 있습니다.

통화 출금 가능 정보는 다음과 같은 주요 항목을 포함합니다.

<table className="custom-table">
    <thead>
      <tr>
        <th>주요 항목</th>
        <th>관련 주요 응답 필드</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td className="code-col"><b>통화 정보:</b><br />출금 수수료, 출금 및 지갑 상태</td>
        <td><code>currency.withdraw_fee</code>,
            <br /><code>currency.wallet_state</code>,
            <br /><code>currency.wallet_support</code></td>
      </tr>
      <tr>
        <td className="code-col"><b>통화 잔고</b></td>
        <td><code>account.balance</code>,<code>account.locked</code>,
            <br /><code>account.avg_buy_price</code></td>
      </tr>
      <tr>
        <td className="code-col"><b>출금 한도:</b><br />1회/일일/잔여한도</td>
        <td><code>withdraw_limit.onetime</code>,
            <br /><code>withdraw_limit.daily</code>,
            <br /><code>withdraw_limit.remaining_daily</code>,
            <br /><code>withdraw_limit.minimum</code>,
            <br /><code>withdraw_limit.can_withdraw</code></td>
      </tr>
      <tr>
        <td className="code-col"><b>출금 관련 계정 정보:</b><br />수수료 레벨, 인증 여부</td>
        <td><code>member_level.fee_level</code>,
            <br /><code>member_level.bank_account_verified</code>,
            <br /><code>member_level.two_factor_auth_verified</code>,
            <br /><code>member_level.locked</code>,
            <br /><code>member_level.wallet_locked</code></td>
      </tr>
    </tbody>
  </table>

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
            <td>2023-05-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/net_type">네트워크 타입(net_type) 필드 추가</a></td>
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
    <a href="auth">인증</a>이 필요한 API로, [출금조회] 권한이 설정된 API Key를 사용해야 합니다. <br />
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
      "WithdrawChance": {
        "type": "object",
        "required": [
          "member_level",
          "currency",
          "account",
          "withdraw_limit"
        ],
        "properties": {
          "member_level": {
            "description": "사용자 보안 등급 정보",
            "type": "object",
            "required": [
              "security_level",
              "fee_level",
              "email_verified",
              "identity_auth_verified",
              "bank_account_verified",
              "two_factor_auth_verified",
              "locked",
              "wallet_locked"
            ],
            "properties": {
              "security_level": {
                "type": "integer",
                "description": "보안 등급",
                "example": 4
              },
              "fee_level": {
                "type": "integer",
                "description": "수수료 등급",
                "example": 0
              },
              "email_verified": {
                "type": "boolean",
                "description": "이메일 인증 여부",
                "example": true
              },
              "identity_auth_verified": {
                "type": "boolean",
                "description": "실명 인증 여부",
                "example": true
              },
              "bank_account_verified": {
                "type": "boolean",
                "description": "계좌 인증 여부",
                "example": true
              },
              "two_factor_auth_verified": {
                "type": "boolean",
                "description": "2FA 활성화 여부",
                "example": true
              },
              "locked": {
                "type": "boolean",
                "description": "계정 보호 상태",
                "example": false
              },
              "wallet_locked": {
                "type": "boolean",
                "description": "출금 보호 상태",
                "example": false
              }
            }
          },
          "currency": {
            "description": "통화 정보",
            "type": "object",
            "required": [
              "code",
              "withdraw_fee",
              "is_coin",
              "wallet_state",
              "wallet_support"
            ],
            "properties": {
              "code": {
                "type": "string",
                "description": "통화 코드",
                "example": "BTC"
              },
              "withdraw_fee": {
                "type": "string",
                "format": "decimal",
                "description": "출금 수수료",
                "example": 0
              },
              "is_coin": {
                "type": "boolean",
                "description": "디지털 자산 여부",
                "example": true
              },
              "wallet_state": {
                "type": "string",
                "enum": [
                  "working",
                  "unsupported"
                ],
                "description": "자산별 입출금 지원 이력 여부. 현시점의 입출금 가능 여부는 wallet_support 필드를 참조하시기 바랍니다.\n  - `working`: 입출금 지원한 적 있음\n  - `unsupported`: 입출금 미지원\n",
                "example": "working"
              },
              "wallet_support": {
                "description": "해당 통화의 입출금 가능 여부. 입금가능한 상태인 경우 'deposit', 출금가능한 상태인 경우 'withdraw'가 추가되며, 빈 리스트일 경우 입출금 모두 불가한 상태입니다.",
                "type": "array",
                "enum": [
                  "deposit",
                  "withdraw"
                ],
                "items": {
                  "type": "string"
                },
                "example": [
                  "deposit",
                  "withdraw"
                ]
              }
            }
          },
          "account": {
            "description": "자산 잔고 정보",
            "type": "object",
            "required": [
              "currency",
              "balance",
              "locked",
              "avg_buy_price",
              "avg_buy_price_modified",
              "unit_currency"
            ],
            "properties": {
              "currency": {
                "type": "string",
                "description": "조회하고자 하는 통화 코드",
                "example": "BTC"
              },
              "balance": {
                "type": "string",
                "format": "decimal",
                "description": "주문 가능 수량 또는 금액.\n\n디지털 자산의 경우 수량, 법정 통화(KRW)의 경우 금액입니다.\n",
                "example": 1000000
              },
              "locked": {
                "type": "string",
                "format": "decimal",
                "description": "출금이나 주문 등에 잠겨 있는 잔액",
                "example": 0
              },
              "avg_buy_price": {
                "type": "string",
                "format": "decimal",
                "description": "매수 평균가",
                "example": 0
              },
              "avg_buy_price_modified": {
                "type": "boolean",
                "description": "매수 평균가 수정 여부",
                "example": false
              },
              "unit_currency": {
                "type": "string",
                "description": "평균가 기준 통화.\n\n[예시] KRW, BTC, USDT\n",
                "example": "KRW"
              }
            }
          },
          "withdraw_limit": {
            "description": "출금 제약 조건",
            "type": "object",
            "required": [
              "currency",
              "onetime",
              "daily",
              "remaining_daily",
              "remaining_daily_fiat",
              "fiat_currency"
            ],
            "properties": {
              "currency": {
                "type": "string",
                "description": "조회하고자 하는 통화 코드",
                "example": "BTC"
              },
              "onetime": {
                "type": "string",
                "deprecated": true,
                "format": "decimal",
                "description": "1회 출금 한도 (deprecated)",
                "example": 0.4
              },
              "daily": {
                "type": "string",
                "nullable": true,
                "deprecated": true,
                "format": "decimal",
                "description": "1일 출금 한도 (deprecated)",
                "example": 500000000
              },
              "remaining_daily": {
                "type": "string",
                "deprecated": true,
                "format": "decimal",
                "description": "1일 잔여 출금 한도 (deprecated)",
                "example": 0
              },
              "remaining_daily_fiat": {
                "type": "string",
                "format": "decimal",
                "description": "통합 1일 잔여 출금 한도",
                "example": 5000000000
              },
              "fiat_currency": {
                "type": "string",
                "description": "기준 법정 통화",
                "example": "KRW"
              },
              "minimum": {
                "type": "string",
                "format": "decimal",
                "description": "최소 출금 금액 또는 수량",
                "example": 300000
              },
              "fixed": {
                "type": "integer",
                "description": "출금 금액/수량 소수점 자리 수",
                "example": 6
              },
              "withdraw_delayed_fiat": {
                "type": "string",
                "format": "decimal",
                "description": "출금 지연 제도로 인해 출금이 제한된 금액",
                "example": 0
              },
              "can_withdraw": {
                "type": "boolean",
                "description": "출금 지원 여부\n* 현재 출금 '가능'여부를 확인하시고 싶은 경우 currency.wallet_support에 'withdraw'의 존재 여부로 확인하실 수 있습니다.\n",
                "example": true
              },
              "remaining_daily_krw": {
                "deprecated": true,
                "type": "string",
                "format": "decimal",
                "description": "통합 1일 잔여 출금 한도 (KRW 기준) (deprecated)",
                "example": 5000000000
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
    "/v1/withdraws/chance": {
      "get": {
        "operationId": "available-withdrawal-information",
        "summary": "출금 가능 정보 조회",
        "description": "지정한 통화에 대한 출금 가능 정보를 조회합니다. 해당 통화의 출금 정책과 사용자 잔고를 확인할 수 있습니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n    --url 'https://api.upbit.com/v1/withdraws/chance?currency=BTC&net_type=BTC' \\\n    --header 'Authorization: Bearer {JWT_TOKEN}' \\\n    --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/withdraws/chance\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"currency\": \"BTC\",\n  \"net_type\": \"BTC\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n}\n\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers, params=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/withdraws/chance\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  currency: \"BTC\",\n  net_type: \"BTC\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"GET\",\n  url: `${baseURL}${path}?${queryString}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class AvailableWithdraw {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/withdraws/chance\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"currency\", \"BTC\");\n        params.put(\"net_type\", \"BTC\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .get()\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "parameters": [
          {
            "name": "currency",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "출금 가능 정보를 조회하고자 하는 통화 코드",
              "example": "BTC"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "BTC"
              }
            }
          },
          {
            "name": "net_type",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "description": "디지털 자산 입출금에 사용되는 블록체인 네트워크 식별자.\n디지털 자산인 경우 필수(required) 필드입니다.\n",
              "example": "BTC"
            },
            "allowReserved": true,
            "examples": {
              "default": {
                "value": "BTC"
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Object of withdrawal policy",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/WithdrawChance"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "member_level": {
                        "security_level": 4,
                        "fee_level": 0,
                        "email_verified": true,
                        "identity_auth_verified": true,
                        "bank_account_verified": true,
                        "two_factor_auth_verified": true,
                        "locked": false,
                        "wallet_locked": false
                      },
                      "currency": {
                        "code": "BTC",
                        "withdraw_fee": "0.0008",
                        "is_coin": true,
                        "wallet_state": "working",
                        "wallet_support": [
                          "deposit",
                          "withdraw"
                        ]
                      },
                      "account": {
                        "currency": "BTC",
                        "balance": "0.0",
                        "locked": "0.0",
                        "avg_buy_price": "145115000",
                        "avg_buy_price_modified": false,
                        "unit_currency": "KRW"
                      },
                      "withdraw_limit": {
                        "currency": "BTC",
                        "onetime": "50.0",
                        "remaining_daily": "0.0",
                        "remaining_daily_fiat": "5000000000.0",
                        "fiat_currency": "KRW",
                        "minimum": "0.00001",
                        "fixed": 8,
                        "withdraw_delayed_fiat": "227.0",
                        "can_withdraw": true,
                        "remaining_daily_krw": "5000000000.0"
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
                  "invalid enum error": {
                    "value": {
                      "error": {
                        "name": "validation_error",
                        "message": "\"field name\" does not have a valid value"
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
          }
        }
      }
    }
  }
}
```

# Sibling pages

* [출금 허용 주소 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawal-addresses.md)
* [디지털 자산 출금 요청](https://docs.upbit.com/kr/reference/withdraw.md)
* [원화 출금 요청](https://docs.upbit.com/kr/reference/withdraw-krw.md)
* [개별 출금 조회](https://docs.upbit.com/kr/reference/get-withdrawal.md)
* [출금 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawals.md)
* [디지털 자산 출금 취소 요청](https://docs.upbit.com/kr/reference/cancel-withdrawal.md)