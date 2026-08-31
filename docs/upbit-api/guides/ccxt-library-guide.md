---
updatedAt: 2026-05-04T07:01:20.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# CCXT 라이브러리 연동 안내

CCXT(CryptoCurrency eXchange Trading Library)는 전 세계 디지털 자산 거래소의 API와 각 거래소에서 지원하는 알고리즘 트레이딩, 전략 백테스팅 등의 기능을 다양한 프로그래밍 언어로 사용할 수 있게 도와주는 오픈 소스 라이브러리입니다.

## 시작하기

CCXT(CryptoCurrency eXchange Trading Library)는 100개 이상의 글로벌 디지털 자산 거래소의 서로 다른 API를 단일(Unified) 인터페이스 연동을 통해 모두 사용할 수 있도록 지원하는 거래소 연동 라이브러리입니다. 다양한 프로그래밍 언어를 지원하며, 직접 API 연동을 위한 모든 부분을 구현하지 않고도 적은 코드만으로 API 연동이 가능하여 편리합니다. 이 가이드에서는 Python과 Node.js 환경에서 CCXT 라이브러리를 사용해 업비트 API를 호출하는 방법을 안내합니다.

<br />

## CCXT 공식 문서

CCXT는 공식 문서를 제공합니다. 아래 링크를 클릭해 공식 웹사이트로 이동할 수 있습니다.

* **CCXT 업비트 API 안내**: <https://docs.ccxt.com/#/exchanges/upbit>
  * CCXT 개발자 문서에 작성된 업비트 API 소개 페이지 입니다. CCXT를 통해 사용할 수 있는 업비트 API와 사용 방법을 확인할 수 있습니다.
* **CCXT 소스 코드**: <https://github.com/ccxt/ccxt>
  * CCXT의 소스 코드를 clone할 수 있는 공식 Github 웹사이트입니다. 전체 커밋 이력과 수정 내역 등을 확인할 수 있습니다.
* **CCXT Docs**: <https://docs.ccxt.com/>
  * CCXT 개발자 문서가 작성된 공식 웹사이트입니다. CCXT가 지원하는 거래소와 API의 종류와 연동 방법 등을 확인할 수 있습니다.

<br />

## Python 연동 가이드

* 최소 버전: Python 3.7.0 +

**가상 환경 구축 및 CCXT 라이브러리 설치**

Python은 프로젝트별 가상 환경을 통한 개발 환경 구축을 권장하고 있습니다. 가상 환경 설정을 사용하는 경우 프로젝트 간 개발 환경이 독립적으로 생성되므로 패키지 충돌이나 버전 문제 등을 효과적으로 해결할 수 있습니다. 가상 환경 설정을 통한 CCXT 라이브러리 실행 순서는 다음과 같습니다.

1. **프로젝트 디렉토리 및 파일 생성**

터미널에 아래와 같이 입력하여 `ccxt_project` 디렉토리를 생성하고 해당 디렉토리 내에 `ccxt_upbit.py` 파일을 생성합니다.

```shell
mkdir ccxt_project
cd ccxt_project

touch ccxt_upbit.py
```

<br />

2. **가상 환경 생성**

터미널에 다음 명령어를 실행해 가상 환경을 생성합니다.

```shell
python3 -m venv .venv
```

<br />

3. **가상 환경 활성화**

운영체제별 가상 환경을 활성화 명령어는 다음과 같습니다. 터미널에 다음 명령어를 실행해 가상 환경을 활성화합니다.

* Linux/MacOS: `source .venv/bin/activate`
* Windows: `.venv\Scripts\activate`

정상적으로 가상 환경이 활성화된 경우, 터미널 프롬프트 앞에 가상 환경 이름이 표시됩니다.

```
(.venv) user@computer:~/project$
```

<br />

4. **CCXT 라이브러리 다운로드**

가상 환경에서 CCXT 라이브러리를 다운로드합니다. 이 패키지는 해당 가상 환경 내에서만 사용이 가능합니다. 가상 환경이 활성화된 터미널에 다음 명령어를 입력해 실행합니다.

```shell
pip install ccxt
```

<br />

5. **CCXT 인스턴스 설정**

CCXT는 거래소별 설정 및 인증 정보를 포함한 상태를 유지할 수 있도록 설계되어 있습니다. 따라서 아래와 같이 인스턴스를 생성하여 사용하는 방식이 일반적입니다. 인스턴스 생성 시 아래와 같이 **Access Key**와 **Secret Key**를 입력해 인증 토큰을 관리할 수 있습니다. 이 외에도 요청 속도 제한, 프록시 등 다양한 옵션 설정을 지원합니다.

```python
# ccxt_upbit.py
import ccxt

access_key = "<YOUR_ACCESS_KEY>"
secret_key = "<YOUR_SECRET_KEY>"

client = ccxt.upbit({
        "apiKey": access_key,
        "secret": secret_key
    })
```

아래 링크를 클릭해 CCXT를 통해 사용할 수 있는 업비트 API를 확인할 수 있습니다.

* **CCXT 업비트 API 안내**: <https://docs.ccxt.com/#/exchanges/upbit>

<br />

6. **인스턴스 설정 확인**

인스턴스가 정상적으로 설정되었는지 확인하기 위해 다음과 같이 `ccxt_upbit.py` 파일의 코드를 작성하고 해당 파일을 실행합니다. 정상적으로 설정된 경우 API의 응답을 반환하고 정상적으로 설정되지 않은 경우, 에러를 반환합니다.

```python Python - API Key 설정 API 호출
# ccxt_upbit.py
import ccxt

access_key = "<YOUR_ACCESS_KEY>"
secret_key = "<YOUR_SECRET_KEY>"

def list_balances():
client = ccxt.upbit({
    "apiKey": access_key,
    "secret": secret_key,
})

    balance = client.fetchBalance()
    print(balance)

if __name__ == "__main__":
    list_balances()

```
```python Python - API Key 비설정 API 호출
# ccxt_upbit.py
import ccxt

def list_markets(symbols: list[str]):
    client = ccxt.upbit()

    market_list = client.fetchMarkets()
    return [item for item in market_list if item.get('symbol') in symbols]

if __name__ == "__main__":
    symbols = ["BTC/KRW", "ETH/KRW", "SOL/KRW"]
    markets = list_markets(symbols)
    print(markets)
```

<br />

CCXT를 통해 실행한 업비트 API 호출 결과를 다음과 같이 확인할 수 있습니다.

```json API Key 설정 API 호출 결과
{
  "info": [
    {
      "currency": "BTC",
      "balance": "0.00000104",
      "locked": "0",
      "avg_buy_price": "160210789.42247014",
      "avg_buy_price_modified": False,
      "unit_currency": "KRW"
    },
...
  ],
  "timestamp": None,
  "datetime": None,
  "BTC": {
    "free": 1.04e-6,
    "used": 0.0,
    "total": 1.04e-6
  },
  "total": {
    "BTC": 1.04e-6,
  }
}

```
```json API Key 비설정 API 호출 결과
[
  {
    "id": "KRW-ETH",
    "lowercaseId": None,
    "symbol": "ETH/KRW",
    "base": "ETH",
    "quote": "KRW",
    "settle": None,
    "baseId": "ETH",
    "quoteId": "KRW",
    "settleId": None,
    "type": "spot",
    "spot": True,
    "margin": False,
    "swap": False,
    "future": False,
    "option": False,
    "index": False,
    "active": True,
    "contract": False,
    "linear": None,
    "inverse": None,
    "subType": None,
    "taker": 0.0005,
    "maker": 0.0005,
    "contractSize": None,
    "expiry": None,
    "expiryDatetime": None,
    "strike": None,
    "optionType": None,
    "precision": {
      "price": 1e-8,
      "amount": 1e-8
    },
    "limits": {
      "leverage": {
        "min": None,
        "max": None
      },
      "amount": {
        "min": None,
        "max": None
      },
      "price": {
        "min": None,
        "max": None
      },
      "cost": {
        "min": None,
        "max": None
      }
    },
    "marginModes": {
      "cross": None,
      "isolated": None
    },
    "created": None,
    "info": {
      "market": "KRW-ETH",
      "korean_name": "이더리움",
      "english_name": "Ethereum"
    }
  },
  {
    "id": "KRW-SOL",
    "lowercaseId": None,
    "symbol": "SOL/KRW",
    "base": "SOL",
    "quote": "KRW",
    "settle": None,
    "baseId": "SOL",
    "quoteId": "KRW",
    "settleId": None,
    "type": "spot",
    "spot": True,
    "margin": False,
    "swap": False,
    "future": False,
    "option": False,
    "index": False,
    "active": True,
    "contract": False,
    "linear": None,
    "inverse": None,
    "subType": None,
    "taker": 0.0005,
    "maker": 0.0005,
    "contractSize": None,
    "expiry": None,
    "expiryDatetime": None,
    "strike": None,
    "optionType": None,
    "precision": {
      "price": 1e-8,
      "amount": 1e-8
    },
    "limits": {
      "leverage": {
        "min": None,
        "max": None
      },
      "amount": {
        "min": None,
        "max": None
      },
      "price": {
        "min": None,
        "max": None
      },
      "cost": {
        "min": None,
        "max": None
      }
    },
    "marginModes": {
      "cross": None,
      "isolated": None
    },
    "created": None,
    "info": {
      "market": "KRW-SOL",
      "korean_name": "솔라나",
      "english_name": "Solana"
    }
  },
  {
    "id": "KRW-BTC",
    "lowercaseId": None,
    "symbol": "BTC/KRW",
    "base": "BTC",
    "quote": "KRW",
    "settle": None,
    "baseId": "BTC",
    "quoteId": "KRW",
    "settleId": None,
    "type": "spot",
    "spot": True,
    "margin": False,
    "swap": False,
    "future": False,
    "option": False,
    "index": False,
    "active": True,
    "contract": False,
    "linear": None,
    "inverse": None,
    "subType": None,
    "taker": 0.0005,
    "maker": 0.0005,
    "contractSize": None,
    "expiry": None,
    "expiryDatetime": None,
    "strike": None,
    "optionType": None,
    "precision": {
      "price": 1e-8,
      "amount": 1e-8
    },
    "limits": {
      "leverage": {
        "min": None,
        "max": None
      },
      "amount": {
        "min": None,
        "max": None
      },
      "price": {
        "min": None,
        "max": None
      },
      "cost": {
        "min": None,
        "max": None
      }
    },
    "marginModes": {
      "cross": None,
      "isolated": None
    },
    "created": None,
    "info": {
      "market": "KRW-BTC",
      "korean_name": "비트코인",
      "english_name": "Bitcoin"
    }
  }
]

```

<br />

6. **가상 환경 비활성화**

프로그램을 종료한 뒤 아래 명령어를 입력해 가상 환경을 비활성화할 수 있습니다. 가상 환경 비활성화 명령어는 운영체제에 상관없이 동일합니다.

```shell
deactivate
```

<br />

## Node.js 연동 가이드

* 최소 버전: Node v7.6 +

**NPM을 통한 CCXT 라이브러리 다운로드**

Node.js는 NPM(Node Package Manager)라는 Node.js 환경의 패키지 관리자를 사용해 간편하게 CCXT 라이브러리를 다운로드할 수 있습니다. NPM을 사용해 CCXT 라이브러리를 다운로드하는 방법은 다음과 같습니다.

1. **프로젝트 디렉토리 및 파일 생성**

터미널에 아래와 같이 입력하여 `ccxt_project` 디렉토리를 생성하고 해당 디렉토리 내에 `ccxt.js` 파일을 생성합니다.

```shell
mkdir ccxt_project
cd ccxt_project

touch ccxt.js
```

<br />

2. **CCXT 라이브러리 다운로드**

NPM을 사용해 CCXT 라이브러리를 다운로드합니다. 터미널에 다음 명령어를 입력하고 실행합니다.

```shell
npm install ccxt
```

<br />

3. **CCXT 인스턴스 설정**

CCXT는 거래소별 설정 및 인증 정보를 포함한 상태를 유지할 수 있도록 설계되어 있습니다. 따라서 아래와 같이 인스턴스를 생성하여 사용하는 방식이 일반적입니다. 인스턴스 생성 시 아래와 같이 **Access Key**와 **Secret Key**를 입력해 인증 토큰을 관리할 수 있습니다. 이 외에도 요청 속도 제한, 프록시 설정 등 다양한 설정값을 추가해 관리할 수 있습니다.

```node
// ccxt.js
const ccxt = require('ccxt');

accessKey = "<YOUR_ACCESS_KEY>";
secretKey = "<YOUR_SECRET_KEY>";

const client = new ccxt.upbit({
    apiKey: accessKey,
    secret: secretKey,
});

```

아래 링크를 클릭해 CCXT를 통해 사용할 수 있는 업비트 API를 확인할 수 있습니다.

* **CCXT 업비트 API 안내**: <https://docs.ccxt.com/#/exchanges/upbit>

<br />

4. **인스턴스 설정 확인**

인스턴스가 정상적으로 설정되었는지 확인하기 위해 다음과 같이 ccxt.js 파일의 코드를 작성하고 해당 파일을 실행합니다. 정상적으로 설정된 경우 API의 응답을 반환하고 정상적으로 설정되지 않은 경우, 에러를 반환합니다.

```node Node - API Key 설정 API 호출
const ccxt = require('ccxt');

const accessKey = "<YOUR_ACCESS_KEY>";
const secretKey = "<YOUR_SECRET_KEY>";

const client = new ccxt.upbit({
    apiKey: accessKey,
    secret: secretKey,
});
async function listBalances() {
    const balance = await client.fetchBalance();
    console.log(balance);
}

listBalances();
```
```node Node - API Key 비설정 API 호출
// ccxt.js
const ccxt = require('ccxt');

const upbit = new ccxt.upbit();
const symbols = ["BTC/KRW", "ETH/KRW", "SOL/KRW"];

async function listMarkets(symbols) {
    const markets = await upbit.fetchMarkets();
    return markets.filter(market => symbols.includes(market.symbol));
}

listMarkets(symbols).then(markets => {
    console.log(markets);
});

```

<br />

CCXT를 통해 실행한 업비트 API 호출 결과를 다음과 같이 확인할 수 있습니다.

```json API Key 설정 API 호출 결과
{
  "info": [
    {
      "currency": "BTC",
      "balance": "0.00000104",
      "locked": "0",
      "avg_buy_price": "160210789.42247014",
      "avg_buy_price_modified": False,
      "unit_currency": "KRW"
    },
...
  ],
  "timestamp": None,
  "datetime": None,
  "BTC": {
    "free": 1.04e-6,
    "used": 0.0,
    "total": 1.04e-6
  },
  "total": {
    "BTC": 1.04e-6,
  }
}

```
```json API Key 비설정 API 호출 결과
[
  {
    "id": "KRW-ETH",
    "lowercaseId": None,
    "symbol": "ETH/KRW",
    "base": "ETH",
    "quote": "KRW",
    "settle": None,
    "baseId": "ETH",
    "quoteId": "KRW",
    "settleId": None,
    "type": "spot",
    "spot": True,
    "margin": False,
    "swap": False,
    "future": False,
    "option": False,
    "index": False,
    "active": True,
    "contract": False,
    "linear": None,
    "inverse": None,
    "subType": None,
    "taker": 0.0005,
    "maker": 0.0005,
    "contractSize": None,
    "expiry": None,
    "expiryDatetime": None,
    "strike": None,
    "optionType": None,
    "precision": {
      "price": 1e-8,
      "amount": 1e-8
    },
    "limits": {
      "leverage": {
        "min": None,
        "max": None
      },
      "amount": {
        "min": None,
        "max": None
      },
      "price": {
        "min": None,
        "max": None
      },
      "cost": {
        "min": None,
        "max": None
      }
    },
    "marginModes": {
      "cross": None,
      "isolated": None
    },
    "created": None,
    "info": {
      "market": "KRW-ETH",
      "korean_name": "이더리움",
      "english_name": "Ethereum"
    }
  },
  {
    "id": "KRW-SOL",
    "lowercaseId": None,
    "symbol": "SOL/KRW",
    "base": "SOL",
    "quote": "KRW",
    "settle": None,
    "baseId": "SOL",
    "quoteId": "KRW",
    "settleId": None,
    "type": "spot",
    "spot": True,
    "margin": False,
    "swap": False,
    "future": False,
    "option": False,
    "index": False,
    "active": True,
    "contract": False,
    "linear": None,
    "inverse": None,
    "subType": None,
    "taker": 0.0005,
    "maker": 0.0005,
    "contractSize": None,
    "expiry": None,
    "expiryDatetime": None,
    "strike": None,
    "optionType": None,
    "precision": {
      "price": 1e-8,
      "amount": 1e-8
    },
    "limits": {
      "leverage": {
        "min": None,
        "max": None
      },
      "amount": {
        "min": None,
        "max": None
      },
      "price": {
        "min": None,
        "max": None
      },
      "cost": {
        "min": None,
        "max": None
      }
    },
    "marginModes": {
      "cross": None,
      "isolated": None
    },
    "created": None,
    "info": {
      "market": "KRW-SOL",
      "korean_name": "솔라나",
      "english_name": "Solana"
    }
  },
  {
    "id": "KRW-BTC",
    "lowercaseId": None,
    "symbol": "BTC/KRW",
    "base": "BTC",
    "quote": "KRW",
    "settle": None,
    "baseId": "BTC",
    "quoteId": "KRW",
    "settleId": None,
    "type": "spot",
    "spot": True,
    "margin": False,
    "swap": False,
    "future": False,
    "option": False,
    "index": False,
    "active": True,
    "contract": False,
    "linear": None,
    "inverse": None,
    "subType": None,
    "taker": 0.0005,
    "maker": 0.0005,
    "contractSize": None,
    "expiry": None,
    "expiryDatetime": None,
    "strike": None,
    "optionType": None,
    "precision": {
      "price": 1e-8,
      "amount": 1e-8
    },
    "limits": {
      "leverage": {
        "min": None,
        "max": None
      },
      "amount": {
        "min": None,
        "max": None
      },
      "price": {
        "min": None,
        "max": None
      },
      "cost": {
        "min": None,
        "max": None
      }
    },
    "marginModes": {
      "cross": None,
      "isolated": None
    },
    "created": None,
    "info": {
      "market": "KRW-BTC",
      "korean_name": "비트코인",
      "english_name": "Bitcoin"
    }
  }
]

```

# Sibling pages

* [개요](https://docs.upbit.com/kr/docs/upbit-ai-agent-overview.md)
* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk.md)
* [Upbit CLI](https://docs.upbit.com/kr/docs/upbit-cli.md)
* [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills.md)
* [Upbit Strategy Toolkit](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit.md)
* [업비트 어시스턴트 이용 안내](https://docs.upbit.com/kr/docs/assistant-guide.md)