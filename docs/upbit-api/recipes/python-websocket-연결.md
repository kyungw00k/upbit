---
updatedAt: 2026-01-07T10:53:34.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Python] Websocket 연결

```python Python
import jwt  # PyJWT
import uuid
import websocket  # websocket-client

def on_message(ws, message):
    """
    message 형식
    """
    data = message.decode('utf-8')
    print(data)


def on_connect(ws):
    """
    연결 및 연결 후 Request 전송
    """
    print("connected!")
    # Request after connection
    ws.send('[{"ticket":"test example"},{"type":"myAsset"}]')


def on_error(ws, err):
    """
    에러 발생 시 처리
    """
    print(err)


def on_close(ws, status_code, msg):
    """
    연결 종료
    """
    print("closed!")

payload = {
    'access_key': "<YOUR_ACCESS_KEY>",
    'nonce': str(uuid.uuid4()),
}

jwt_token = jwt.encode(payload, "<YOUR_SECRET_KEY>");
authorization_token = f'Bearer {jwt_token}'
headers = {"Authorization": authorization_token}

ws_app = websocket.WebSocketApp("wss://api.upbit.com/websocket/v1/private",
                                header=headers,
                                on_message=on_message,
                                on_open=on_connect,
                                on_error=on_error,
                                on_close=on_close)
ws_app.run_forever(ping_interval=30, ping_timeout=10, reconnect=2) 
```

```json Response Example
connected!
{"type":"myAsset","asset_uuid":"<my_asset_uuid>","assets":[{"currency":"KRW","balance":20517.4157543899035,"locked":15007.5}],"asset_timestamp":1753769511603,"timestamp":1753769511612,"stream_type":"REALTIME"}
{"type":"myAsset","asset_uuid":"<my_asset_uuid>","assets":[{"currency":"XRP","balance":1.15580212,"locked":0}],"asset_timestamp":1753769511603,"timestamp":1753769511613,"stream_type":"REALTIME"}

```

# 유틸 라이브러리 Import

<!-- python@1-3 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우 `pip install <module name>` 명령어를 실행해 설치할 수 있습니다. 

websocket의 경우, `pip install websocket-client` 명령어를 실행해 websocket 라이브러리를 설치해야 합니다.

# 메시지 수신 시 처리 방식 정의

<!-- python@5-10 -->

WebSocket으로 전달된 실시간 스트림 데이터를 UTF-8 디코딩된 형태로 출력하는 예제입니다.

# 요청 이벤트 정의

<!-- python@13-19 -->

Websocket이 연결된 후 수행할 동작을 정의합니다. 이 가이드에서는 Websocket 연결 이후, 사용자의 주문 및 체결 데이터를 수신하기 위한 요청 메세지를 전송합니다.

# 에러 이벤트 정의

<!-- python@22-26 -->

Websocket 연결 중 에러 발생 시 처리 방법을 정의합니다. 이 가이드에서는 에러 발생 시 사용자에게 전체 에러 메세지를 노출합니다.

# 연결 종료 이벤트 정의

<!-- python@29-33 -->

Websocket 연결 종료 후 실행할 동작을 정의합니다. 이 가이드에서는 연결 종료 시 “closed!” 라는 문자열을 출력합니다.

# Websocket 정의 및 실행

<!-- python@40-50 -->

Websocket 요청 시 필요한 payload와 JWT를 생성하고 이를 헤더에 추가해 Websocket을 연결합니다. 이 때, 앞서 정의한 함수를 인자로 전달하여 상황 별로 실행할 로직을 적용합니다.