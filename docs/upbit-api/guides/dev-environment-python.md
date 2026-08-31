---
updatedAt: 2026-05-04T05:41:32.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Python

Python 환경에서 Upbit Open API를 연동하기 위한 개발 환경 설정 방법을 안내합니다.

## macOS 환경 설정

macOS에서 Python을 설치하는 방법입니다. Homebrew라는 macOS용 소프트웨어 패키지 관리자를 사용해 명령어로 간편하게 설치할 수 있습니다.

<br />

1. **Homebrew 설치**

터미널에서 다음 명령어를 실행해 Homebrew를 설치합니다.

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

설치를 완료하고 터미널에서 다음 명령어를 실행해 Homebrew의 버전을 확인합니다.

```shell
brew -v

Homebrew 4.6.1
```

<br />

2. **Python 설치**

터미널에서 다음 명령어를 실행해 Python을 설치합니다. Homebrew를 사용해 Python을 설치합니다.

```shell
brew install python
```

설치를 완료하고 터미널에서 다음 명령어를 실행해 Python의 버전을 확인합니다. 정상적으로 설치를 완료한 경우 Python의 버전을 확인할 수 있습니다.

```shell
python3 --version

Python 3.13.2
```

<br />

3. **가상 환경 설정**

프로젝트별 의존성 충돌 방지를 위해 가상 환경을 사용합니다.

```shell
python3 -m venv .venv

source .venv/bin/activate
```

정상적으로 가상 환경이 활성화된 경우, 프롬프트에 (venv)가 표시됩니다.

```shell
(.venv) user@computer:~/project$
```

가상 환경 내에서 개발 환경에 따라 필요한 package를 자유롭게 설치할 수 있으며, 이 때 설치된 package들은 글로벌 환경에 변경을 주지 않으므로 편리하게 프로젝트별 의존성 관리를 수행할 수 있습니다.

개발을 완료한 후 터미널에서 다음 명령어를 실행해 가상 환경을 비활성화할 수 있습니다.

```shell
deactivate
```

<br />

## Windows 환경 설정

1. **Python 공식 웹사이트에서 설치 파일 다운로드**

Windows 운영체제에서 Python을 사용하기 위해서는 Python 공식 웹사이트에서 제공하는 설치 파일을 다운로드해야 합니다. 아래 링크를 클릭해 공식 웹사이트를 방문하고 설치 파일을 다운로드하세요. 설치 과정에서 **\[Add Python to PATH]** 옵션을 선택하면 별도의 환경 변수 설정 없이 Python을 바로 사용할 수 있습니다.

* [Python 다운로드 바로가기](https://www.python.org/downloads/)

<br />

2. **가상 환경 설정**

Python은 프로젝트간 충돌을 방지하고 의존성을 관리하기 위해 가상 환경을 구성해 개발 환경을 관리합니다. Python을 설치한 후 터미널에 다음 명령어를 실행해 가상 환경을 구성하고 활성화할 수 있습니다.

```shell
python -m venv .venv

.venv\Scripts\activate
```

개발을 완료한 후 터미널에서 다음 명령어를 실행해 가상 환경을 비활성화할 수 있습니다

```shell
deactivate
```

<br />

### HTTP 클라이언트 라이브러리 안내

REST API와 WebSocket을 호출하기 위해 사용할 수 있는 대표적인 HTTP 클라이언트 라이브러리를 소개합니다.

#### REST API - `requests` 라이브러리

1. **설치**

```shell
pip install requests
```

<br />

2. **기본 사용법**

```python
import requests

url = "https://api.upbit.com/v1/ticker?markets=KRW-BTC"
response = requests.get(url)
data = response.json()
print(data[0]["trade_price"])  
```

* GET, POST 등 다양한 HTTP 요청을 쉽게 보낼 수 있습니다.

<br />

#### WebSockets - 'websocket-client', 'websockets' 라이브러리

각각 동기 방식의 WebSocket 연결, 비동기 방식의 WebSocket 연결을 위한 라이브러리입니다. 개발 환경에 맞게 선택하여 사용하실 수 있습니다.

1. **설치**

```shell
pip install websocket-client
pip install websockets
```

<br />

2. **기본 사용법**

```python
import websocket # websocket-client
import json

def on_message(ws, message):
    print("Received:", message)

ws = websocket.WebSocketApp(
    "wss://api.upbit.com/websocket/v1",
    on_message=on_message
)

subscribe_message = [
    {"ticket":"test"},
    {"type":"ticker","codes":["KRW-BTC"]}
]

def on_open(ws):
    ws.send(json.dumps(subscribe_message))

ws.on_open = on_open
ws.run_forever(ping_interval=30, ping_timeout=10, reconnect=2)
```

* 콜백 함수(on\_message)를 통해 서버로부터 받은 메시지를 처리할 수 있습니다.
* WebSocket 연결 관리 및 이벤트 처리를 위해 on\_open, on\_close, on\_error 등의 콜백 함수를 함께 구현할 수 있습니다.
  * on\_open: 서버 연결 직후 구독 메시지 전송 등 초기 작업 처리
  * on\_close: 연결 종료 시 리소스 정리, 재연결 로직 등 구현 가능
  * on\_error: 네트워크 오류 등 예외 상황 대응
* 인증이 필요한 /private WebSocket 연결 시에는 Header에 인증 정보를 포함해야 합니다. 아래 Recipe를 참고하여 구현할 수 있습니다.

[block:tutorial-tile]
{
  "backgroundColor": "#ffffff",
  "emoji": "🔌",
  "id": "691be0e426e5b5b0e24ced18",
  "link": "https://docs.upbit.com/v1.6.0/recipes/python-websocket-연결",
  "slug": "python-websocket-연결",
  "title": "[Python] Websocket 연결"
}
[/block]

# Sibling pages

* [Java](https://docs.upbit.com/kr/docs/dev-environment-java.md)
* [Node.js](https://docs.upbit.com/kr/docs/dev-environment-nodejs.md)