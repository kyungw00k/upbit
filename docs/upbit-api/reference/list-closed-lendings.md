---
updatedAt: 2026-11-24T10:46:46.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# 종료된 코인빌리기 목록 조회

종료된 코인빌리기 목록과 종료 정보를 조회합니다.

해당 포켓에서 신청한 코인빌리기 중 상환 또는 청산이 완료되어 종료된 포지션 목록을 조회합니다. `borrowed_currency`, `start_time`, `end_time`을 지정하여 빌린 자산과 조회 기간별로 조회할 수 있습니다.

#### 조회 대상

상환 또는 청산이 완료된 코인빌리기 포지션을 조회합니다.

| 포지션 상태       | 종료 유형(`closure_type`)   | 설명               |
| ------------ | ----------------------- | ---------------- |
| `repaid`     | `repay`                 | 빌린 자산으로 전액 상환 완료 |
| `liquidated` | `repay_with_collateral` | 담보로 전액 상환 완료     |
| `liquidated` | `liquidation`           | 강제상환 완료          |

담보상환과 강제상환은 모두 `liquidated` 상태로 종료되며, `closure_type`으로 구분합니다. 강제상환 사유는 `liquidation_reason`에서 확인할 수 있습니다.

#### 종료된 포지션의 수량 및 미수금

| 필드                | 설명                                                                                                                                         |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `borrowed_amount` | 포지션 종료 직전 남아 있던 빌린 수량입니다.<br /><br />전액 상환 후에도 0이 아닌 값이 반환될 수 있으며, 현재 갚아야 할 빌린 수량을 의미하지 않습니다.                                              |
| `unpaid_amount`   | 강제상환 정산 후 갚지 못한 빌린 수량과 미수취 수수료의 합계입니다.<br /><br />빌린 자산(`borrowed_currency`) 기준 수량으로, 미수금이 있는 경우 숫자 형식의 문자열로 반환하고, 미수금이 없으면 `null`을 반환합니다. |

* 미수금 추심 내역은 개별 코인빌리기 조회에서 신청번호(`lending_number`)를 확인한 후 업비트 웹·앱에서 확인해 주세요.
* 종료되지 않은 포지션은 진행 중 코인빌리기 목록 조회에서 확인할 수 있습니다.

#### 조회 기준 및 제약 사항

1. 조회 기간과 정렬은 포지션의 종료 시각(`closed_at`)이 아닌 **생성 시각(**`created_at`**)을 기준**으로 합니다.
2. 별도의 조회 조건을 지정하지 않으면 요청 시점 기준 최근 30일 이내에 생성된 종료 포지션을 생성 시각(`created_at`) 기준 최신순으로 20건 조회합니다. `limit`을 지정하면 한 번에 최대 100건까지 조회할 수 있습니다.
3. `start_time`과 `end_time`으로 조회 기간을 지정할 수 있습니다.
   * 두 파라미터를 모두 생략하면 요청 시점 이전 30일간 생성된 종료 포지션을 조회합니다.
   * `start_time`만 지정하면 해당 시각부터 30일간 생성된 종료 포지션을 조회합니다.
   * `end_time`만 지정하면 해당 시각 이전 30일간 생성된 종료 포지션을 조회합니다.
   * 두 파라미터를 모두 지정하면 해당 기간에 생성된 종료 포지션을 조회합니다.
4. 조회 기간은 최대 30일까지 지정할 수 있습니다. 30일을 초과하면 `invalid_parameter` 에러가 반환되므로, 더 긴 기간은 30일 이하로 나누어 조회해 주세요.
5. 조회 시각은 ISO 8601 형식(`2025-06-24T13:56:53+09:00`) 또는 밀리초 단위 타임스탬프(`1750741013000`)로 입력할 수 있습니다.
   * ISO 8601 형식은 URL 인코딩해 전달해 주세요. 오프셋의 `+`는 `%2B`로 인코딩합니다.
   * 같은 시작 시각을 두 형식으로 전달하는 예시입니다. ISO 8601: `?start_time=2025-06-24T13%3A56%3A53%2B09%3A00`, 밀리초 타임스탬프: `?start_time=1750741013000`
   * 입력한 시각의 소수 초는 버리며, `created_at`이 시작 초부터 종료 초까지인 포지션을 조회합니다. 예를 들어 `end_time`이 `13:56:53.001`이면 `13:56:53`에 생성된 포지션이 모두 포함됩니다.
   * 최대 조회 기간(30일)도 소수 초를 버린 시각을 기준으로 계산합니다.
6. 목록은 커서 기반 페이지네이션 방식으로 조회합니다.
   * 최초 요청에서는 `cursor`를 지정하지 않습니다.
   * 다음 페이지가 있는 경우 응답의 `next_cursor`를 다음 요청의 `cursor`에 입력합니다. `next_cursor`가 `null`이면 마지막 페이지입니다.
   * 다음 페이지를 조회할 때는 `borrowed_currency`, `start_time`, `end_time` 등 다른 조회 조건을 이전 요청과 동일하게 유지해야 합니다.
7. 여러 페이지에 걸쳐 조회하는 경우 동일한 시점의 데이터가 보장되지 않습니다. 조회 중 포지션 정보가 변경될 수 있습니다.
8. 조회 조건에 맞는 포지션이 없는 경우 다음과 같이 빈 목록을 반환합니다.

```json
{
  "lendings": [],
  "next_cursor": null
}
```

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

  코인빌리기 권한은 API Key 발급 시 이용 자격(고객 확인, 약관 동의 등)을 확인한 뒤 설정할 수 있습니다. DAXA 온라인 교육을 이수하지 않아도 권한 설정과 조회 API 이용은 가능합니다. 이후 동의를 철회하거나 이용 자격이 바뀌면 공통 에러가 반환됩니다.
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
      "ClosedLendings": {
        "type": "object",
        "required": [
          "lendings",
          "next_cursor"
        ],
        "properties": {
          "lendings": {
            "type": "array",
            "items": {
              "type": "object",
              "required": [
                "uuid",
                "state",
                "closure_type",
                "liquidation_reason",
                "liquidation_ltv",
                "borrowed_currency",
                "borrowed_amount",
                "unpaid_amount",
                "collaterals",
                "final_fee",
                "created_at",
                "expire_at",
                "updated_at",
                "closed_at"
              ],
              "properties": {
                "uuid": {
                  "type": "string",
                  "description": "포지션의 UUID"
                },
                "state": {
                  "type": "string",
                  "enum": [
                    "repaid",
                    "liquidated"
                  ],
                  "description": "포지션 상태\n\n* `repaid`: 정상 상환 완료\n* `liquidated`: 청산 또는 담보상환 완료"
                },
                "closure_type": {
                  "type": "string",
                  "enum": [
                    "repay",
                    "repay_with_collateral",
                    "liquidation"
                  ],
                  "description": "포지션의 종료 유형\n\n* `repay`: 빌린 자산으로 전액 상환\n* `repay_with_collateral`: 담보로 전액 상환\n* `liquidation`: 강제상환"
                },
                "liquidation_reason": {
                  "type": "string",
                  "enum": [
                    "ltv_exceeded",
                    "expired",
                    "scheduled",
                    "account_status",
                    null
                  ],
                  "description": "강제상환 사유\n\n* `ltv_exceeded`: 강제상환 기준 LTV 도달\n* `expired`: 이용기간 만료\n* `scheduled`: 거래지원 종료 등 자산 이벤트에 따른 예약 청산\n* `account_status`: 회원 상태, 고객 확인 등 계정 관련 사유",
                  "nullable": true
                },
                "liquidation_ltv": {
                  "type": "string",
                  "description": "해당 포지션에 적용된 강제상환 기준 LTV"
                },
                "borrowed_currency": {
                  "type": "string",
                  "description": "빌린 자산의 자산 코드"
                },
                "borrowed_amount": {
                  "type": "string",
                  "description": "포지션 종료 직전에 남아 있던 빌린 수량. 전액 상환 후에도 0이 아닌 값이 반환될 수 있으며, 현재 갚아야 할 빌린 수량을 의미하지 않습니다."
                },
                "unpaid_amount": {
                  "type": "string",
                  "description": "강제상환 정산 후 갚지 못한 빌린 수량과 미수취 수수료의 합계. 빌린 자산(`borrowed_currency`) 기준 수량을 숫자 형식의 문자열로 반환하며, 미수금이 없으면 `null`을 반환합니다.",
                  "nullable": true
                },
                "collaterals": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "required": [
                      "currency",
                      "amount",
                      "factor"
                    ],
                    "properties": {
                      "currency": {
                        "type": "string",
                        "description": "담보자산의 자산 코드"
                      },
                      "amount": {
                        "type": "string",
                        "description": "포지션 종료 직전의 담보 수량"
                      },
                      "factor": {
                        "type": "string",
                        "description": "해당 포지션에 적용된 담보 인정 비율"
                      }
                    }
                  },
                  "description": "담보자산 목록"
                },
                "final_fee": {
                  "type": "object",
                  "description": "포지션 종료 후 확정된 최종 수수료 정보",
                  "nullable": true,
                  "properties": {
                    "currency": {
                      "type": "string",
                      "description": "최종 수수료의 자산 코드"
                    },
                    "amount": {
                      "type": "string",
                      "description": "빌린 자산 기준 최종 수수료 합계"
                    }
                  }
                },
                "created_at": {
                  "type": "string",
                  "description": "빌리기 시작 시각 (ISO 8601)"
                },
                "expire_at": {
                  "type": "string",
                  "description": "최초 약정된 만기일 (ISO 8601)"
                },
                "updated_at": {
                  "type": "string",
                  "description": "포지션 정보가 마지막으로 변경된 시각 (ISO 8601)"
                },
                "closed_at": {
                  "type": "string",
                  "description": "포지션의 상환 또는 청산이 완료된 시각 (ISO 8601)"
                }
              }
            },
            "description": "조회 조건을 만족하는 종료된 포지션 목록"
          },
          "next_cursor": {
            "type": "string",
            "description": "다음 페이지 조회에 사용할 커서",
            "nullable": true
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
    "/v1/lendings/closed": {
      "get": {
        "operationId": "list-closed-lendings",
        "summary": "종료된 코인빌리기 목록 조회",
        "description": "종료된 코인빌리기 목록과 종료 정보를 조회합니다.",
        "tags": [
          "Exchange"
        ],
        "parameters": [
          {
            "name": "borrowed_currency",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회할 빌린 자산의 자산 코드"
            }
          },
          {
            "name": "start_time",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회 시작 시각. ISO 8601 형식(`2025-06-24T13:56:53+09:00`) 또는 밀리초 단위 타임스탬프(`1750741013000`). ISO 8601 형식은 URL 인코딩이 필요하며(`+` → `%2B`), 소수 초는 버립니다."
            }
          },
          {
            "name": "end_time",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "조회 종료 시각. 형식은 `start_time`과 같으며, 소수 초를 버린 뒤 해당 초에 생성된 포지션까지 포함합니다.\n\n* 둘 다 미입력: 요청 시점 이전 30일간 생성된 종료 포지션\n* `start_time`만 입력: `start_time` 이후 30일간\n* `end_time`만 입력: `end_time` 이전 30일간\n* 둘 다 입력: `start_time`부터 `end_time`까지(30일 이하). 30일을 넘으면 `invalid_parameter`"
            }
          },
          {
            "name": "limit",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "description": "조회할 건수",
              "default": 20,
              "minimum": 1,
              "maximum": 100
            }
          },
          {
            "name": "cursor",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "description": "다음 페이지 조회의 기준이 되는 포지션의 UUID"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ClosedLendings"
                },
                "examples": {
                  "Example": {
                    "description": "정상 응답",
                    "value": {
                      "lendings": [
                        {
                          "uuid": "7d38eec5-8e15-48d0-999f-81333b76387b",
                          "state": "liquidated",
                          "closure_type": "liquidation",
                          "liquidation_reason": "ltv_exceeded",
                          "liquidation_ltv": "0.9500",
                          "borrowed_currency": "USDT",
                          "borrowed_amount": "50000.00000000",
                          "unpaid_amount": "137.45000000",
                          "collaterals": [
                            {
                              "currency": "BTC",
                              "amount": "0.58000000",
                              "factor": "0.90"
                            },
                            {
                              "currency": "ETH",
                              "amount": "3.20000000",
                              "factor": "0.90"
                            }
                          ],
                          "final_fee": {
                            "currency": "USDT",
                            "amount": "575.00000000"
                          },
                          "created_at": "2026-09-04T10:15:32+09:00",
                          "expire_at": "2026-10-04T10:15:32+09:00",
                          "updated_at": "2026-09-07T11:42:18+09:00",
                          "closed_at": "2026-09-07T11:42:18+09:00"
                        },
                        {
                          "uuid": "942662c8-4587-489d-b775-6a58aa9ce47f",
                          "state": "liquidated",
                          "closure_type": "repay_with_collateral",
                          "liquidation_reason": null,
                          "liquidation_ltv": "0.9500",
                          "borrowed_currency": "USDT",
                          "borrowed_amount": "35000.00000000",
                          "unpaid_amount": null,
                          "collaterals": [
                            {
                              "currency": "BTC",
                              "amount": "0.41000000",
                              "factor": "0.90"
                            },
                            {
                              "currency": "ETH",
                              "amount": "2.50000000",
                              "factor": "0.90"
                            }
                          ],
                          "final_fee": {
                            "currency": "USDT",
                            "amount": "542.50000000"
                          },
                          "created_at": "2026-08-25T13:42:17+09:00",
                          "expire_at": "2026-09-24T13:42:17+09:00",
                          "updated_at": "2026-09-05T11:08:32+09:00",
                          "closed_at": "2026-09-05T11:08:32+09:00"
                        },
                        {
                          "uuid": "3f75bf24-bc12-4dc4-a87c-77ac8a5ab251",
                          "state": "repaid",
                          "closure_type": "repay",
                          "liquidation_reason": null,
                          "liquidation_ltv": "0.9500",
                          "borrowed_currency": "USDT",
                          "borrowed_amount": "20000.00000000",
                          "unpaid_amount": null,
                          "collaterals": [
                            {
                              "currency": "BTC",
                              "amount": "0.32000000",
                              "factor": "0.90"
                            }
                          ],
                          "final_fee": {
                            "currency": "USDT",
                            "amount": "300.00000000"
                          },
                          "created_at": "2026-08-09T09:30:00+09:00",
                          "expire_at": "2026-09-08T09:30:00+09:00",
                          "updated_at": "2026-09-08T09:28:41+09:00",
                          "closed_at": "2026-09-08T09:28:41+09:00"
                        }
                      ],
                      "next_cursor": "3f75bf24-bc12-4dc4-a87c-77ac8a5ab251"
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
                  "invalid_parameter": {
                    "description": "* limit 범위 초과 (100 초과)\n* limit 0 또는 음수\n* limit에 정수형이 아님\n* borrowed_currency - 잘못된 형식\n* 허용되지 않는 자산 코드 - 빈 문자열, 소문자 - 존재하지 않는 코드 등\n* start_time\n* end_time 포맷 위반\n* start_time이 end_time보다 늦음\n* 날짜 범위 30일 초과\n* cursor 형식 오류 (uuid)",
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
              "code": "curl --request GET \\\n  --url 'https://api.upbit.com/v1/lendings/closed?borrowed_currency=USDT&limit=3' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt",
              "code": "import hashlib\nimport os\nimport uuid\nfrom urllib.parse import unquote, urlencode\n\nimport jwt  # PyJWT\nimport requests\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/lendings/closed\"\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\"borrowed_currency\": \"USDT\", \"limit\": 3}\nquery_string = unquote(urlencode(params, doseq=True))\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": hashlib.sha512(query_string.encode(\"utf-8\")).hexdigest(),\n    \"query_hash_alg\": \"SHA512\",\n}\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS512\")\n\nheaders = {\"Authorization\": f\"Bearer {jwt_token}\", \"Accept\": \"application/json\"}\nres = requests.get(f\"{BASE_URL}{PATH}?{query_string}\", headers=headers)\nprint(res.status_code, res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\n\nconst BASE_URL = \"https://api.upbit.com\";\nconst PATH = \"/v1/lendings/closed\";\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = { borrowed_currency: \"USDT\", limit: 3 };\nconst queryString = Object.entries(params)\n  .map(([key, value]) => `${key}=${value}`)\n  .join(\"&\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: crypto.createHash(\"sha512\").update(queryString, \"utf-8\").digest(\"hex\"),\n  query_hash_alg: \"SHA512\",\n};\nconst jwtToken = sign(payload, SECRET_KEY, { algorithm: \"HS512\" });\n\naxios\n  .get(`${BASE_URL}${PATH}?${queryString}`, {\n    headers: { Authorization: `Bearer ${jwtToken}`, Accept: \"application/json\" },\n  })\n  .then((res) => console.log(res.data))\n  .catch((err) => console.error(err.response ? err.response.data : err.message));\n"
            },
            {
              "language": "java",
              "code": "import com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.util.LinkedHashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.Response;\n\npublic class ListClosedLendings {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/lendings/closed\";\n\n    public static void main(String[] args) throws Exception {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new LinkedHashMap<>();\n        params.put(\"borrowed_currency\", \"USDT\");\n        params.put(\"limit\", \"3\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + e.getValue())\n            .collect(Collectors.joining(\"&\"));\n\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", sha512(queryString))\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(Algorithm.HMAC512(secretKey));\n\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH + \"?\" + queryString)\n            .addHeader(\"Authorization\", \"Bearer \" + jwtToken)\n            .addHeader(\"Accept\", \"application/json\")\n            .get()\n            .build();\n\n        try (Response response = new OkHttpClient().newCall(request).execute()) {\n            System.out.println(response.code() + \" \" + Objects.requireNonNull(response.body()).string());\n        }\n    }\n\n    private static String sha512(String input) throws Exception {\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest(input.getBytes(StandardCharsets.UTF_8))) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        return sb.toString();\n    }\n}\n"
            }
          ]
        }
      }
    }
  }
}
```

# 같은 수준의 페이지

* [코인빌리기 가능 정보 조회](https://docs.upbit.com/kr/reference/available-lending-information.md)
* [코인빌리기 신청 테스트](https://docs.upbit.com/kr/reference/lending-test.md)
* [코인빌리기 신청](https://docs.upbit.com/kr/reference/new-lending.md)
* [코인빌리기 요청 상태 조회](https://docs.upbit.com/kr/reference/get-lending-request.md)
* [개별 코인빌리기 조회](https://docs.upbit.com/kr/reference/get-lending.md)
* [코인빌리기 추가 신청](https://docs.upbit.com/kr/reference/borrow-more.md)
* [담보 추가](https://docs.upbit.com/kr/reference/adjust-collateral.md)
* [코인빌리기 상환](https://docs.upbit.com/kr/reference/repay-lending.md)
* [담보로 상환](https://docs.upbit.com/kr/reference/repay-lending-with-collateral.md)
* [진행 중 코인빌리기 목록 조회](https://docs.upbit.com/kr/reference/list-ongoing-lendings.md)
* [코인빌리기 거래 내역 조회](https://docs.upbit.com/kr/reference/list-lending-transactions.md)