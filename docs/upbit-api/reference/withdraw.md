---
updatedAt: 2026-07-29T11:04:00.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 디지털 자산 출금 요청

디지털 자산 출금을 요청합니다.

### API Key 출금 안심차단 비활성화

2026년 7월 9일부터 출금 권한이 포함된 API Key는 발급 시 **출금 안심차단이 자동으로 적용**됩니다. 출금 안심차단이 적용된 API Key로는 **출금 API를 호출할 수 없습니다.**

출금 API를 사용하려면 업비트 모바일 앱의 **\[더보기 > 인증/보안 > Open API 관리 > 활성]** 페이지에서 해당 API Key를 선택한 후, **출금 안심차단 토글을 비활성화해 주세요.**

### 출금 허용 주소 등록

API를 통한 디지털 자산 출금시 출금 대상 주소를 **반드시 사전에 출금 허용 주소로 등록**해야 합니다. 출금 허용 주소는 업비트 홈페이지를 통해 등록할 수 있으며 API를 통한 등록은 불가합니다. 출금 허용 주소 등록 방법은 [출금 허용 주소 등록하기](https://docs.upbit.com/kr/docs/open-api-withdraw_access_register) 페이지를 참고해주세요.

### 트랜잭션 유형 - 일반 출금(default)과 바로 출금(internal)

디지털 자산을 개인 지갑 주소 또는 타 거래소 지갑 주소로 자산을 출금하고자 하는 경우 일반 출금으로 진행합니다. 일반 출금은 실제 블록체인 트랜잭션으로 자산 송금을 진행함에 따라 트랜잭션 확정 시간을 포함한 처리 시간과 출금 수수료가 발생합니다.

업비트 계정 간 디지털 자산 송금 시 바로 출금을 사용할 수 있습니다. 바로 출금은 블록체인을 사용하지 않는 자산 송금 방식으로, 약 1분 내외로 출금이 반영되며 수수료가 발생하지 않습니다. **바로 출금시 반드시 업비트 회원의 지갑 주소로 출금**해야 하며, 업비트 회원 주소가 아닌 주소로 출금을 요청하는 경우 출금이 정상적으로 수행되지 않습니다.

일반 출금과 바로 출금에 관한 자세한 사항은 [업비트 고객센터 - 자주하는 질문](https://support.upbit.com/hc/ko/articles/34786146778009-%EC%9D%BC%EB%B0%98-%EC%B6%9C%EA%B8%88%EA%B3%BC-%EB%B0%94%EB%A1%9C-%EC%B6%9C%EA%B8%88%EC%9D%98-%EC%B0%A8%EC%9D%B4%EC%A0%90%EC%9D%B4-%EB%AC%B4%EC%97%87%EC%9D%B8%EA%B0%80%EC%9A%94)를 참고하시기 바랍니다.

  <div className="callout-section">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i>  네트워크 타입("net_type")은 필수 파라미터입니다.
      </div>
      <b>네트워크 타입("net_type")</b>은 디지털 자산 입출금시 실제 자산이 이동되는 블록체인 네트워크(대상 체인)를 지정하기 위한 식별자 필드(예: BTC)입니다. 디지털 자산 출금 시 필수 파라미터로, 정상적인 입출금 진행을 위해 정확한 식별자 값을 사용해야 합니다.
      디지털 자산 출금 API 호출 시 사전에 출금 <a href="list-withdrawal-addresses">허용 주소 목록 조회 API</a>를 호출하여 응답으로부터 정확한 네트워크 타입 값을 참조하여 사용하시기 바랍니다. 
  </div>

<div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> POST API에 대한 Form 방식 요청은 2022년 3월 1일부로 지원이 종료되었습니다.
      </div>
    Form 방식 지원 종료에 따라 Urlencoded Form 방식으로 전송하는 POST 요청에 대한 정상적인 동작을 보장하지 않습니다. <b>반드시 JSON 형식으로 요청 본문(Body)을 전송</b>해주시기 바랍니다.
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
            <td>2023-05-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/net_type">네트워크 타입(net_type) 필드 추가</a></td>
      </tr>
tr>
            <td className="code-col">-</td>
            <td>2020-05-29</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_feature_0529"><code>transaction_type</code> 필드 추가</a></td>
        </tr>
tr>
            <td className="code-col">-</td>
            <td>2019-04-23</td>
            <td><a href="https://docs.upbit.com/kr/changelog/open-api-%EB%B3%80%EA%B2%BD%EC%82%AC%ED%95%AD-%EC%95%88%EB%82%B4-%EC%A0%81%EC%9A%A9-%EC%9D%BC%EC%9E%90-423-1100">바로 출금 기능 지원</a></td>
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
    <a href="auth">인증</a>이 필요한 API로, [출금하기] 권한이 설정된 API Key를 사용해야 합니다. <br />
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
      "Withdraw": {
        "type": "object",
        "required": [
          "type",
          "uuid",
          "currency",
          "net_type",
          "txid",
          "state",
          "created_at",
          "amount",
          "fee",
          "transaction_type",
          "is_cancelable"
        ],
        "properties": {
          "type": {
            "type": "string",
            "description": "입출금 종류",
            "example": "withdraw",
            "default": "withdraw"
          },
          "uuid": {
            "type": "string",
            "description": "출금의 유일 식별자",
            "example": "9f432943-54e0-40b7-825f-b6fec8b42b79"
          },
          "currency": {
            "type": "string",
            "description": "조회하고자 하는 통화 코드",
            "example": "BTC"
          },
          "net_type": {
            "type": "string",
            "nullable": true,
            "description": "출금 네트워크 유형.\n\n업비트에서 사용하는 블록체인 네트워크 구분자입니다.\n\n[예시] \"ETH\", \"TRX\", \"SOL\"\n",
            "example": "BTC"
          },
          "txid": {
            "type": "string",
            "nullable": true,
            "description": "트랜잭션 ID",
            "example": 1234567890
          },
          "state": {
            "type": "string",
            "enum": [
              "WAITING",
              "PROCESSING",
              "DONE",
              "FAILED",
              "CANCELLED",
              "REJECTED"
            ],
            "description": "출금 처리 상태\n  - `WAITING`: 대기 중\n  - `PROCESSING`: 처리 중\n  - `DONE`: 완료\n  - `FAILED`: 실패\n  - `CANCELLED`: 취소됨\n  - `REJECTED`: 거절됨\n",
            "example": "WAITING"
          },
          "created_at": {
            "type": "string",
            "description": "출금 생성 시간",
            "example": "2024-01-01T00:00:00"
          },
          "done_at": {
            "type": "string",
            "nullable": true,
            "description": "출금 완료 시간. 출금이 완료되지 않은 경우 null로 반환됩니다.",
            "example": "2024-01-01T00:00:00"
          },
          "amount": {
            "type": "string",
            "format": "decimal",
            "description": "출금하고자 하는 자산의 수량",
            "example": "0.01"
          },
          "fee": {
            "type": "string",
            "format": "decimal",
            "description": "출금 수수료",
            "example": 0
          },
          "transaction_type": {
            "type": "string",
            "enum": [
              "default",
              "internal"
            ],
            "description": "출금 유형.\n\n* `default`: 일반출금\n* `internal`: 바로출금\n",
            "example": "default",
            "default": "default"
          },
          "is_cancelable": {
            "type": "boolean",
            "nullable": true,
            "description": "출금 취소 가능 여부",
            "example": false
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
    "/v1/withdraws/coin": {
      "post": {
        "operationId": "withdraw",
        "summary": "디지털 자산 출금 요청",
        "description": "디지털 자산 출금을 요청합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request POST \\\n  --url 'https://api.upbit.com/v1/withdraws/coin' \\\n  --header 'Authorization: Bearer {JWT_TOKEN}' \\\n  --header 'Accept: application/json' \\\n  --header 'Content-Type: application/json' \\\n  --data '\n{\n\"currency\": \"BTC\",\n\"net_type\": \"BTC\",\n\"amount\": \"0.01\",\n\"address\": \"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa\"\n}\n'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/withdraws/coin\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n  \"currency\": \"BTC\",\n  \"net_type\": \"BTC\",\n  \"amount\": \"0.01\",\n  \"address\": \"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa\",\n}\n\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n  \"access_key\": ACCESS_KEY,\n  \"nonce\": str(uuid.uuid4()),\n  \"query_hash\": query_hash,\n  \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n  \"Authorization\": f\"Bearer {jwt_token}\",\n  \"Accept\": \"application/json\",\n  \"Content-Type\": \"application/json\",\n}\n\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/withdraws/coin\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  currency: \"BTC\",\n  net_type: \"BTC\",\n  amount: \"0.01\",\n  address: \"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"POST\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n    \"Content-Type\": \"application/json\",\n  },\n  data: params,\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\nimport com.google.gson.Gson;\n\npublic class WithdrawCoin {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/withdraws/coin\";\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"currency\", \"BTC\");\n        params.put(\"net_type\", \"BTC\");\n        params.put(\"amount\", \"0.01\");\n        params.put(\"address\", \"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        String jsonBody = new Gson().toJson(params);\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .post(RequestBody.create(jsonBody, okhttp3.MediaType.parse(\"application/json; charset=utf-8\")))\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "requestBody": {
          "description": "디지털 자산 출금",
          "content": {
            "application/json": {
              "examples": {
                "default": {
                  "value": {
                    "currency": "BTC",
                    "net_type": "BTC",
                    "amount": "0.01",
                    "address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
                  }
                }
              },
              "schema": {
                "type": "object",
                "required": [
                  "currency",
                  "net_type",
                  "amount",
                  "address"
                ],
                "properties": {
                  "currency": {
                    "type": "string",
                    "description": "출금하고자 하는 디지털 자산의 통화 코드",
                    "example": "BTC"
                  },
                  "net_type": {
                    "type": "string",
                    "description": "출금 주소 등록 후 출금 허용 주소 목록 조회 API를 호출하여 응답에서 각 주소로의 출금시 사용 가능한 \"net_type\" 값을 확인할 수 있습니다.\n",
                    "example": "BTC"
                  },
                  "amount": {
                    "type": "string",
                    "format": "decimal",
                    "description": "출금하고자 하는 자산의 수량.\n\n숫자 형식의 String으로 입력합니다.\n",
                    "example": "0.01"
                  },
                  "address": {
                    "type": "string",
                    "description": "디지털 자산 출금 시 수신 주소.\n\n출금 가능 주소 목록에 등록된 주소만 사용 가능합니다.\n",
                    "example": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
                  },
                  "secondary_address": {
                    "type": "string",
                    "nullable": true,
                    "description": "2차 출금 주소.\n일부 디지털 자산의 경우 Destination Tag, Memo, 또는 Message와 같은 2차 주소를 포함합니다. 수신 주소 정보에 2차 주소가 포함되어 있다면 이 필드를 반드시 포함하여 출금을 요청해야 합니다.\n"
                  },
                  "transaction_type": {
                    "type": "string",
                    "enum": [
                      "default",
                      "internal"
                    ],
                    "description": "출금 유형.\n\n* `default`: 일반출금\n* `internal`: 바로출금\n",
                    "example": "default",
                    "default": "default"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Object of digital asset withdrawal",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Withdraw"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "type": "withdraw",
                      "uuid": "9f432943-54e0-40b7-825f-b6fec8b42b79",
                      "currency": "BTC",
                      "net_type": "BTC",
                      "txid": "ebe6937b-130e-4066-8ac6-4b0e67f28adc",
                      "state": "PROCESSING",
                      "created_at": "2018-04-13T11:24:01+09:00",
                      "amount": "0.01",
                      "fee": "0.0",
                      "transaction_type": "default",
                      "is_cancelable": false
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

* [출금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-withdrawal-information.md)
* [출금 허용 주소 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawal-addresses.md)
* [원화 출금 요청](https://docs.upbit.com/kr/reference/withdraw-krw.md)
* [개별 출금 조회](https://docs.upbit.com/kr/reference/get-withdrawal.md)
* [출금 목록 조회](https://docs.upbit.com/kr/reference/list-withdrawals.md)
* [디지털 자산 출금 취소 요청](https://docs.upbit.com/kr/reference/cancel-withdrawal.md)