---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 코인빌리기 가능 정보 조회

코인빌리기에 이용할 수 있는 자산과 담보 조건, 계정별 이용 한도를 조회합니다.

코인빌리기를 신청하기 전에 확인해야 하는 빌릴 자산 정책, 담보자산 정책과 계정별 이용 한도를 한 번에 조회합니다. 응답은 다음 항목으로 구성됩니다.

<Table>
  <thead>
    <tr>
      <th style={{ width: "41.4%" }}>
        주요 항목
      </th>

      <th style={{ width: "58.5%" }}>
        관련 응답 필드
      </th>
    </tr>
  </thead>

  <tbody>
    <tr>
      <td>
        빌릴 수 있는 자산과<br />신규·추가신청 가능 여부
      </td>

      <td>
        `lending_assets[].currency`, `lending_assets[].borrow_enabled`, `lending_assets[].borrow_more_enabled`
      </td>
    </tr>

    <tr>
      <td>
        신청 가능한 LTV 범위와<br />강제상환 기준 LTV
      </td>

      <td>
        `lending_assets[].min_ltv`, `lending_assets[].max_ltv`, `lending_assets[].liquidation_ltv`
      </td>
    </tr>

    <tr>
      <td>
        이용수수료율,<br />강제상환수수료율,<br />담보상환수수료율
      </td>

      <td>
        `lending_assets[].borrowing_fee_rate`, `lending_assets[].liquidation_fee_rate`, `lending_assets[].collateral_repayment_fee_rate`
      </td>
    </tr>

    <tr>
      <td>
        자산별 최소 신청 금액, 최대 이용 건수, 신청 가능 수량,<br />이용기간
      </td>

      <td>
        `lending_assets[].min_fiat`, `lending_assets[].max_lending_count`, `lending_assets[].available_budget_amount`, `lending_assets[].term`
      </td>
    </tr>

    <tr>
      <td>
        원화(KRW)를 포함한 담보 가능 자산,<br />최소 담보 평가 금액,<br />담보 인정 비율
      </td>

      <td>
        `collateral_assets[].currency`, `collateral_assets[].min_fiat`, `collateral_assets[].factor`
      </td>
    </tr>

    <tr>
      <td>
        계정의 대여한도,<br />담보 한도,<br />이용 건수
      </td>

      <td>
        `account_limits`
      </td>
    </tr>

    <tr>
      <td>
        1건에 설정할 수 있는 담보자산 종류 수
      </td>

      <td>
        `max_collateral_count`
      </td>
    </tr>
  </tbody>
</Table>

`lending_assets`와 `collateral_assets`에는 현재 지원 중인 자산만 포함되며, 페이지네이션 없이 전체 목록을 반환합니다.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### **포켓의 실제 잔고는 반영하지 않습니다.**&#x20;

  실제 신청 가능 여부와 예상 LTV는 코인빌리기 신청 테스트로 확인해 주세요. 담보로 맡길 자산의 보유 수량은 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/get-balance">포켓 잔고 조회</Anchor>API로 확인할 수 있습니다.
</Callout>

#### 신청 전 확인 순서

1. 이 API로 빌릴 자산의 신청 가능 여부(`borrow_enabled`)와 정책, 계정의 남은 대여한도(`account_limits.remaining_lending_fiat`)를 확인합니다.
2. 빌릴 수량과 담보를 정한 뒤 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-test">코인빌리기 신청 테스트</Anchor>로 신청 가능 여부와 예상 LTV(`estimated_ltv`)를 확인합니다.
3. <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/new-lending">코인빌리기 신청</Anchor>을 요청합니다.

#### 한도와 신청 가능 수량

* 대여한도(`account_limits.max_lending_fiat`)는 계정에 적용되는 한도로, 자산이나 포켓과 관계없이 계정 레벨에 따라 정해집니다.
* 최대 이용 건수는 계정 단위(`account_limits.max_lending_count`)와 자산 단위(`lending_assets[].max_lending_count`)로 각각 적용됩니다. 현재 이용 중인 건수는 `account_limits.ongoing_lending_count`로 확인합니다.
* `available_budget_amount`는 자산별 전체 공급 수량에서 빌리기 중인 수량을 뺀 값으로, 계정 한도가 아닙니다. 해당 값은 준실시간으로 갱신되므로 조회 시점과 실제 신청 시점의 수량이 다를 수 있습니다. 신청 시 공급 수량이 부족하면 `insufficient_lending_budget` 에러가 반환됩니다.

#### 담보 한도

* 담보 한도(`account_limits.max_collateral_fiat`)는 계정에서 담보로 설정할 수 있는 원화(KRW) 및 지원되는 디지털 자산의 총 평가금액 한도입니다. 담보 한도 계산에는 담보 인정 비율(factor)을 적용하지 않습니다.
* `account_limits.used_collaterals`는 이용 중·청산 중·이용 제한 상태의 포지션을 포함하여 아직 반환되지 않은 담보의 자산별 수량을 반환합니다. 원화(KRW)는 금액을 그대로 반영하고, 디지털 자산은 평가가격을 적용하여 남은 담보 한도를 계산할 수 있습니다.
* 디지털 자산의 가격 변동과 업비트 기준 평가가격에 따라 실제 남은 담보 한도는 달라질 수 있습니다. 자세한 가격 기준은 [코인빌리기 가이드 FAQ](https://www.upbit.com/lending/guide?tab=faq)를 참고해 주세요.

#### 수수료

* 이용수수료는 일 단위 이용수수료율(`borrowing_fee_rate`)에 **일 최대 수수료 수량을 적용해** 신청 다음 날부터 하루 단위로 계산하며, 당일 상환해도 1일치가 적용됩니다.
* 강제상환 또는 담보상환 시에는 다음과 같이 계산한 수수료가 추가로 적용됩니다.
  * **강제상환수수료** = 강제상환 시점의 남은 빌린 수량 × 강제상환수수료율(`liquidation_fee_rate`)
  * **담보상환수수료** = 담보상환 처리 직전 남은 빌린 수량 × 담보상환수수료율(`collateral_repayment_fee_rate`)
* 이용수수료는 이용 중에 차감하지 않고 최종 상환 시 정산합니다. 단, 담보상환하거나 강제상환할 때는 반대매매로 확보한 자산 또는 남은 담보 등으로 정산할 수 있습니다

#### **수량 입력 규칙**

* 수량 파라미터는 숫자 형식의 문자열로 입력하며, 소수점 이하 최대 8자리까지 허용합니다. 소수점 이하 자릿수가 8자리를 초과하면 에러가 반환됩니다.
* 자산 수량은 소수점 아래 최대 20자리까지 반환될 수 있습니다. `entries.fiat_value`은 소수점 아래 8자리로 반올림하여 제공합니다. 필드의 용도에 따라 소수점 자릿수가 다를 수 있습니다.

***

#### 에러 안내

코인빌리기 API 이용 중 에러가 발생하면 반환된 에러 코드와 발생 사유를 확인하고, 해당 조건에 맞게 조치한 뒤 다시 요청해 주세요.

발생 가능한 에러와 상세 응답 명세는 [코인빌리기 API 안내](https://docs.upbit.com/kr/reference/lending-overview)의 **에러 안내** 항목을 참고해 주세요.

#### Rate Limit

초당 최대 30회 호출할 수 있습니다.&#x20;

포켓 단위로 측정되며 \[`default` 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

#### API Key Permission

[인증](https://docs.upbit.com/kr/reference/auth)이 필요한 API로, 메인포켓 API Key만 사용할 수 있습니다. \[코인빌리기] 권한이 설정된 API Key를 사용해야 합니다. 권한 오류(`out_of_scope`)가 발생하면 <Anchor target="_blank" href="https://www.upbit.com/mypage/open_api_management">API Key 관리</Anchor> 메뉴에서 권한 설정을 확인해 주세요.

<Callout icon="far fa-triangle-exclamation" theme="warn">
  #### 코인빌리기 API 이용 권한

  코인빌리기 권한은 API Key 발급 시 이용 자격(고객 확인, 약관 동의 등)을 확인한 뒤 설정할 수 있습니다. DAXA 온라인 교육을 이수하지 않아도 권한 설정과 조회 API 이용은 가능합니다. 이후 동의를 철회하거나 이용 자격이 바뀌면 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/lending-overview">공통 에러</Anchor>가 반환됩니다.
</Callout>

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
      "LendingChance": {
        "type": "object",
        "required": [
          "lending_assets",
          "collateral_assets",
          "account_limits",
          "max_collateral_count"
        ],
        "properties": {
          "lending_assets": {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "currency",
                "borrow_enabled",
                "borrow_more_enabled",
                "min_ltv",
                "max_ltv",
                "liquidation_ltv",
                "borrowing_fee_rate",
                "liquidation_fee_rate",
                "collateral_repayment_fee_rate",
                "min_fiat",
                "max_lending_count",
                "available_budget_amount",
                "term"
              ],
              "properties": {
                "currency": {
                  "type": "string",
                  "description": "빌릴 수 있는 자산의 자산 코드"
                },
                "borrow_enabled": {
                  "type": "boolean",
                  "description": "해당 자산의 신규 신청 가능 여부\n\n* `true`: 신규 신청 가능 상태\n* `false`: 신규 신청 불가 상태"
                },
                "borrow_more_enabled": {
                  "type": "boolean",
                  "description": "해당 자산의 추가신청 가능 여부\n\n* `true`: 추가신청 가능 상태\n* `false`: 추가신청 불가 상태"
                },
                "min_ltv": {
                  "type": "string",
                  "description": "신청 가능한 최소 LTV"
                },
                "max_ltv": {
                  "type": "string",
                  "description": "신청 가능한 최대 LTV"
                },
                "liquidation_ltv": {
                  "type": "string",
                  "description": "강제상환 기준 LTV"
                },
                "borrowing_fee_rate": {
                  "type": "string",
                  "description": "신청 전 기준 현재 정책의 일 단위 이용수수료율. 신청 후 포지션에 적용된 요율은 개별 코인빌리기 조회에서 확인합니다."
                },
                "liquidation_fee_rate": {
                  "type": "string",
                  "description": "신청 전 기준 현재 정책의 강제상환수수료율. 강제상환 시 추가로 부과됩니다. 신청 후 포지션에 적용된 요율은 개별 코인빌리기 조회에서 확인합니다."
                },
                "collateral_repayment_fee_rate": {
                  "type": "string",
                  "description": "신청 전 기준 현재 정책의 담보상환수수료율. 담보상환 시 추가로 부과됩니다. 신청 후 포지션에 적용된 요율은 개별 코인빌리기 조회에서 확인합니다."
                },
                "min_fiat": {
                  "type": "string",
                  "description": "최소 신청 금액 (단위: KRW)"
                },
                "max_lending_count": {
                  "type": "integer",
                  "description": "해당 자산으로 동시에 이용할 수 있는 최대 건수"
                },
                "available_budget_amount": {
                  "type": "string",
                  "description": "해당 자산의 신청 가능한 최대 수량. 조회 시점 기준 값이며 준실시간으로 갱신되므로, 실제 신청 시점의 수량과 다를 수 있습니다. 신청 시 수량이 소진되면 `insufficient_lending_budget` 에러가 반환됩니다."
                },
                "term": {
                  "type": "string",
                  "description": "이용기간 (단위: 일)"
                }
              }
            },
            "description": "빌릴 수 있는 자산별 정책과 현재 신청 가능 정보"
          },
          "collateral_assets": {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "currency",
                "min_fiat",
                "factor"
              ],
              "properties": {
                "currency": {
                  "type": "string",
                  "description": "담보로 설정할 수 있는 자산의 자산 코드. 원화(KRW)도 담보로 설정할 수 있습니다."
                },
                "min_fiat": {
                  "type": "string",
                  "description": "해당 자산을 담보로 설정할 때 필요한 최소 담보 평가 금액 (단위: KRW)"
                },
                "factor": {
                  "type": "string",
                  "description": "담보 인정 비율"
                }
              }
            },
            "description": "담보로 사용할 수 있는 자산별 정책"
          },
          "account_limits": {
            "type": "object",
            "description": "API Key 소유 계정에 적용되는 코인빌리기 이용 한도",
            "required": [
              "max_lending_fiat",
              "remaining_lending_fiat",
              "max_lending_count",
              "ongoing_lending_count",
              "max_collateral_fiat",
              "used_collaterals"
            ],
            "properties": {
              "max_lending_fiat": {
                "type": "string",
                "description": "계정의 최대 대여한도 (단위: KRW)"
              },
              "remaining_lending_fiat": {
                "type": "string",
                "description": "계정의 남은 대여한도 (단위: KRW)"
              },
              "max_lending_count": {
                "type": "integer",
                "description": "계정이 동시에 이용할 수 있는 최대 빌리기 건수"
              },
              "ongoing_lending_count": {
                "type": "integer",
                "description": "계정에서 현재 이용 중인 빌리기 건수. 이용 제한(`suspended`)과 청산·담보상환 진행 중(`liquidating`) 포지션도 포함합니다."
              },
              "max_collateral_fiat": {
                "type": "string",
                "description": "계정에서 담보로 설정할 수 있는 자산의 총 평가 금액 한도 (단위: KRW)"
              },
              "used_collaterals": {
                "type": "array",
                "items": {
                  "type": "object",
                  "required": [
                    "currency",
                    "amount"
                  ],
                  "properties": {
                    "currency": {
                      "type": "string",
                      "description": "담보로 사용 중인 자산의 자산 코드"
                    },
                    "amount": {
                      "type": "string",
                      "description": "아직 반환되지 않은 담보 수량. 이용 제한(`suspended`)과 청산·담보상환 진행 중(`liquidating`) 포지션의 담보도 포함합니다."
                    }
                  }
                },
                "description": "현재 담보로 사용 중인 자산별 수량"
              }
            }
          },
          "max_collateral_count": {
            "type": "integer",
            "description": "코인빌리기 1건에 담보로 설정할 수 있는 최대 자산 종류 수"
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
    "/v1/lendings/chance": {
      "get": {
        "operationId": "available-lending-information",
        "summary": "코인빌리기 가능 정보 조회",
        "description": "코인빌리기에 이용할 수 있는 자산과 담보 조건, 계정별 이용 한도를 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/LendingChance"
                },
                "examples": {
                  "Example": {
                    "description": "정상 응답",
                    "value": {
                      "lending_assets": [
                        {
                          "currency": "BTC",
                          "borrow_enabled": true,
                          "borrow_more_enabled": true,
                          "min_ltv": "0.00001",
                          "max_ltv": "0.85",
                          "liquidation_ltv": "0.95",
                          "borrowing_fee_rate": "0.0005",
                          "liquidation_fee_rate": "0.01",
                          "collateral_repayment_fee_rate": "0.01",
                          "min_fiat": "10000",
                          "max_lending_count": 1,
                          "available_budget_amount": "2.73184200",
                          "term": "30"
                        },
                        {
                          "currency": "ETH",
                          "borrow_enabled": true,
                          "borrow_more_enabled": true,
                          "min_ltv": "0.00001",
                          "max_ltv": "0.85",
                          "liquidation_ltv": "0.95",
                          "borrowing_fee_rate": "0.0005",
                          "liquidation_fee_rate": "0.01",
                          "collateral_repayment_fee_rate": "0.01",
                          "min_fiat": "10000",
                          "max_lending_count": 1,
                          "available_budget_amount": "184.52000000",
                          "term": "30"
                        }
                      ],
                      "collateral_assets": [
                        {
                          "currency": "KRW",
                          "min_fiat": "20000",
                          "factor": "1.0"
                        },
                        {
                          "currency": "ETH",
                          "min_fiat": "20000",
                          "factor": "0.90"
                        }
                      ],
                      "account_limits": {
                        "max_lending_fiat": "100000000",
                        "remaining_lending_fiat": "100000000",
                        "max_lending_count": 1,
                        "ongoing_lending_count": 1,
                        "max_collateral_fiat": "300000000",
                        "used_collaterals": [
                          {
                            "currency": "KRW",
                            "amount": "1500000"
                          },
                          {
                            "currency": "ETH",
                            "amount": "0.20000000"
                          }
                        ]
                      },
                      "max_collateral_count": 4
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
                  "out_of_scope": {
                    "description": "코인빌리기(`lending`) 권한이 없는 API Key로 호출",
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
        },
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/lendings/chance' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import os\nimport uuid\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/chance\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\npayload = {\"access_key\": ACCESS_KEY, \"nonce\": str(uuid.uuid4())}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"}\nres = requests.get(f\"{BASE_URL}{PATH}\", headers=headers)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/chance\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst payload = { access_key: ACCESS_KEY, nonce: uuidv4() };\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .get(`${BASE_URL}${PATH}`, {\n    headers: { Authorization: `Bearer ${jwtToken}`, Accept: \"application/json\" },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.util.Objects;\nimport java.util.UUID;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class AvailableLendingInformation {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/chance\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .get()\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        }
      }
    }
  }
}
```

# 같은 수준의 페이지

* [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test.md)
* [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending.md)
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request.md)
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [종료된 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-closed-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)