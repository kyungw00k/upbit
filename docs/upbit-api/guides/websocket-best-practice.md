---
updatedAt: 2026-08-26T10:05:07.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# WebSocket 연동 Best Practice

업비트 WebSocket 연동 구현을 위한 가이드라인 문서로서 인증, 연결 관리, 요청 수 제한 등 실제 구현 시 참고해야 하는 구현 요구사항을 안내합니다.

## 시작하기

이 가이드는 단순한 WebSocket 연결 예제를 보다 안정적이고 실전적인 연동 예제로 발전시키기 위해 필요한 권장 고려사항과 모범 사례를 제시합니다. 구독하고자 하는 데이터의 종류에 따른 인증, 연결 관리 기법, 요청 수 제한 등 WebSocket 연동 시 고려해야 하는 구현 요구사항들을 적용해볼 수 있습니다.

<br />

## 업비트 WebSocket 연동시 고려해야 할 사항

<br />

### 최초 연결과 구독 요청

WebSocket연동을 통한 데이터 조회 과정은 크게 **(1)최초 연결 생성 단계**와 **(2)구독 요청 메세지 전송을 통한 데이터 수신** 단계의 두 단계로 나누어 진행됩니다.

최초 연결 생성 단계는 양방향 통신을 위한 통신 채널을 여는 단계로, 연결이 성공적으로 완료되면 클라이언트에서 서버로, 또는 서버에서 클라이언트로 데이터를 전송할 수 있게 됩니다. 클라이언트로부터 별도의 요청이 있기 전까지 기본적으로 업비트 WebSocket 서버는 데이터 전송 없이 타임아웃까지 대기하게 됩니다.

데이터를 구독하기 위해서는 클라이언트에서 WebSocket 서버로 데이터 구독 요청 메세지를 전송해야 합니다. 구독 요청 메세지는 같은 필드들을 포함합니다.

1. **요청 티켓(ticket)**: 요청을 구분하기 위한 Ticket ID.
2. **구독하고자 하는 데이터 유형(type)** : candle, ticker, orderbook, trade, announcement, myOrder(내 주문), myAsset(내 자산) 등 조회하고 싶은 데이터의 타입과 페어 목록을 명시합니다.

* **하나의 타입에 대해 여러 페어의 정보를 수신**할 수 있으며
* **동시에 여러 유형의 데이터를 구독**할 수 있습니다.<br />예를 들어, 아래와 같이 요청하는 경우 KRW-BTC, KRW-ETH 호가 데이터와 KRW-XRP 체결 데이터를 동시에 구독합니다.

```json
[{"ticket":"uuid_"},{"type":"orderbook", "codes":["KRW-BTC","KRW-ETH"]},{"type":"trade", "codes":["KRW-XRP"]}]
```

3\. **데이터 형식(format)**: 기본형(DEFAULT) 또는 축약형(SIMPLE) 형식으로 데이터를 구독할 수 있습니다. 축약형으로 구독 시 JSON 데이터의 각 key값이 짧게(예시: type > ty) 사용되어 데이터 크기를 줄일 수 있습니다.  또한 목록형(JSON\_LIST) 또는 축약 목록형(SIMPLE\_LIST)로 요청하여 응답을 JSON Object 형식이 아닌 JSON Array(List) 형식으로 수신할 수 있습니다.

동일한 데이터 항목에 대해, **구독 요청은 여러 번 전송할 필요 없이 1회 요청만으로도 서버로부터 지속적으로 스트림 데이터를 수신**할 수 있습니다. 또한 새로운 데이터를 구독하고 싶은 경우 새로 연결을 생성할 필요 없이 새로운 구독 메세지를 전송하여 **이전 구독을 중단하고 새로운 데이터 스트림 구독을 시작**할 수 있습니다.

<br />

### 스냅샷 데이터와 실시간 스트림 데이터

WebSocket을 통해 수신할 수 있는 데이터는 스냅샷 데이터와 실시간 스트림의 두가지 유형으로 구분됩니다. 사용자는 구독 요청 시 두 데이터를 모두 요청하거나, 스냅샷 데이터만을 요청할 수 있습니다.

* **스냅샷** 데이터란, 요청 시점의 정보를 1회 수신하는 방식입니다.
* **실시간 스트림**은 WebSocket 연결이 유지되는 동안 지속적으로 정보가 수신되는 방식으로, 데이터 항목에 따라 일정한 갱신 주기마다 또는 이벤트가 발생하는 시점에 정보를 수신할 수 있습니다.

<br />

### 연결 유지 및 재연결

업비트 WebSocket 서버는 120초 동안 데이터 송수신이 없으면 Idle Timeout으로 연결을 종료합니다. 데이터가 전송되지 않더라도  [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide) 문서의 안내와 같이 ping/pong 옵션, timeout 설정 또는 명시적 Ping 메시지 전송을 활용하여 연결을 유지할 수 있습니다. 또한 연결이 해제된 경우 적절한 연결 재시도를 통해 데이터 누락을 방지할 수 있습니다.

<br />

### 인증

내 자산, 주문 및 체결 정보 구독은 \[wss\://api.upbit.com/websocket/v1 Endpoint]가 아닌 \[wss\://api.upbit.com/websocket/v1/private Endpoint] 연동을 통해서만 가능합니다. /private 채널 연결 요청시 API Key로 생성한 [인증](https://docs.upbit.com/kr/reference/auth) 토큰이 요청 헤더에 반드시 포함되어야 합니다.

<br />

### 요청 수 제한(Rate Limit) 정책 준수

REST API 뿐만 아니라 WebSocket도 연결 요청 및 메세지 전송에 대해 요청 수 제한 정책을 적용하고 있습니다. 따라서 Websocket 연동 구현 시 과도한 연결 요청 및 메세지 발신이 이루어지지 않도록 적절한 throttling 전략을 구현해야 합니다. 단, REST API와 달리 WebSocket은 한 번의 요청으로 지속적인 데이터 수신이 가능하므로 상대적으로 제한의 영향이 적습니다.

<br />

## Best Practice - Python 예제로 알아보는 구현 가이드라인

튜토리얼 및 예제 코드에서 볼 수 있는 간단한 호출 예제로 시작하여 실 연동 시 고려해야 하는 사항들을 구현에 적용하는 과정을 Python 예제를 통해 순차적으로 안내합니다.

<div className="callout-section">
  <div className="callout-title">
    <i className="fa-solid fa-circle-exclamation"></i> 가이드라인 예제 코드 
    </div>
    <li>본 가이드라인에서 사용된 Python 코드는 보다 다양한 버전의 Python 환경을 지원하기 위해 타입 힌트와 최신 문법을 적극적으로 사용하기 보다는 <b>호환성을 고려하여 구현</b>되었습니다. </li>
    <li>websocket-client 라이브러리에서 동기 방식의 연결 방식에 한정한 예제로, <b>절대적인 구현 제약사항이 아닌 하나의 구현 예시</b>로 작성되었음을 사용시 참고부탁드립니다. </li>
    <li>보다 간결한 예제 코드 작성을 위해 docstring 및 주석은 최소화하였으며, <b>주요 코드 설명은 코드 상단 문서 영역을 참조</b> 부탁드립니다.</li>
</div>

<br />

### 기본 WebSocket 연결 예제

[개발 환경 설정 가이드](https://docs.upbit.com/kr/docs/dev-environment)에서 안내한 websocket-client 라이브러리를 활용한 기본적인 업비트 WebSocket 연결 예제는 아래와 같습니다. 업비트 WebSocket 서버와 연결만 생성 및 종료하는 예제로, 구독 메세지 전송 및 인증과 같은 기능은 수행하지 않습니다.

* websocket-client의 WebSocketApp을 기반으로 이벤트 루프를 통해 WebSocket 서버로 연결하며, open, message, error, close 상황에 대한 콜백 함수를 작성하여 각 상황이 발생했을 때 수행할 동작을 정의할 수 있습니다. 이를 포함한 ThreadedWebSocketApp 클래스를 구현하여 WebSocket 연결 시 새로운 스레드를 실행하여 연결 및 데이터를 구독하고, 연결 해제할 수 있도록 하였습니다.
* 메인 함수에서는 ThreadedWebSocketApp 인스턴스를 생성하여 연결 스레드를 실행하고 종료합니다.
* 150초 이후 자동으로 **연결 종료하는 예제로 작성**되었습니다. 실제 사용시에는 해당 부분은 제외하여 연결을 유지해야 합니다.

```python
import threading
import websocket
import json
import time

class ThreadedWebSocketApp(threading.Thread):
    def __init__(self, url):
        # type: (str) -> None
        threading.Thread.__init__(self)
        self.daemon = True
        self.url = url
        self.ws_app = None
        self._stop_evt = threading.Event()

    @staticmethod
    def connect(url):
        # type: (str) -> "ThreadedWebSocketApp"
        t = ThreadedWebSocketApp(url)
        t.start()
        return t

    def run(self):
        self.ws_app = websocket.WebSocketApp(
            self.url,
            on_open=self._on_open,
            on_message=self._on_message,
            on_error=self._on_error,
            on_close=self._on_close
        )

        self.ws_app.run_forever()
        self.ws_app = None

    def close(self):
        self._stop_evt.set()
        try:
            if self.ws_app:
                self.ws_app.close()
        except Exception:
            pass

    def send_message(self, message):
        try:
            if self.ws_app and self.ws_app.sock and self.ws_app.sock.connected:
                self.ws_app.send(message)
        except Exception as e:
            self.on_error(e)

    def _on_open(self, ws):
        print("Opened")

    def _on_message(self, ws, data):
        try:
            obj = json.loads(data)
            print("Received(JSON):", obj)
        except Exception:
            print("Received(raw):", data)

    def _on_error(self, ws, err):
        print("Error:", err)

    def _on_close(self, ws, code, reason):
        print("Closed")


if __name__ == "__main__":
    ws = ThreadedWebSocketApp.connect(url="wss://api.upbit.com/websocket/v1")

    try:
        time.sleep(150)
    finally:
        ws.close()
        ws.join(timeout=3)
```

위 코드를 실행하면, 연결을 생성하고 아무 동작도 하지 않은채 150초간 대기하게 됩니다. 하지만 업비트 WebSocket 서버의 정책상 120초간 데이터 송수신이 없는 경우 연결이 해제되기 때문에 약 120초 경과 후 아래와 같은 에러 메세지와 함께 연결이 종료되는 것을 확인할 수 있습니다. (프로그램은 150초가 모두 경과한 후 종료됩니다.)

```shell
Opened
Error: Connection to remote host was lost.
Closed
```

<br />

### 연결 관리 - 연결 유지 및 재연결

myOrder 또는 myAsset과 같은 private 스트림 데이터의 경우 실제 주문 또는 자산의 변동이 발생할 때만 데이터가 수신되기 때문에, 위와 같은 기본 예제로는 충분히 대기하며 연결을 유지할 수 없습니다. 이를 방지하기 위해 아래와 같이 연결 유지 및 연결 재시도 코드를 추가하였습니다. 이러한 개선을 통해, WebSocketApp은 이제 설정된 ping\_interval(기본 값 30초)에 따라 일정 간격으로 ping 프레임을 전송하고 pong이 10초간 수신되지 않으면 연결을 종료합니다. 또한 서버 측 연결 해제가 발생한 경우 최대 재시도 횟수(기본 값 3회) 내에서 연결을 재시도 할 수 있습니다.

* ping, 재연결 관련 설정을 연결 생성 시 주입 (13-14 라인, 27-28 라인, 48-50라인)
* 최대 재시도 횟수 내에서 연결을 재시도합니다. 단, 사용자 요청에 의한 연결 종료의 경우 재시도 대상에서 제외합니다. (84-93 라인)

```python
import threading
import websocket
import json
import time

class ThreadedWebSocketApp(threading.Thread):
    def __init__(self, url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0):
        # type: (str, int, int, int, float) -> None
        threading.Thread.__init__(self)
        self.daemon = True
        self.url = url
        self.ping_interval = ping_interval
        self.ping_timeout = ping_timeout
        self.max_retries = max_retries
        self.retry_sleep = retry_sleep

        self.ws_app = None
        self._stop_evt = threading.Event()

    @staticmethod
    def connect(url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0):
        # type: (str, int, int, int, float) -> "ThreadedWebSocketApp"
        t = ThreadedWebSocketApp(
            url,
            ping_interval=ping_interval,
            ping_timeout=ping_timeout,
            max_retries=max_retries,
            retry_sleep=retry_sleep
        )
        t.start()
        return t

    def run(self):
        attempts = 0
        while not self._stop_evt.is_set():
            self._opened_once = False
            self.ws_app = websocket.WebSocketApp(
                self.url,
                on_open=self._on_open,
                on_message=self._on_message,
                on_error=self._on_error,
                on_close=self._on_close
            )
                
            self.ws_app.run_forever(
                ping_interval=self.ping_interval,
                ping_timeout=self.ping_timeout,
                reconnect=int(self.retry_sleep) if self.retry_sleep else None
            )
            self.ws_app = None

    def close(self):
        self._stop_evt.set()
        try:
            if self.ws_app:
                self.ws_app.close()
                try:
                    self.ws_app.keep_running = False
                except Exception:
                    pass
        except Exception:
            pass

    def send_message(self, message):
        try:
            if self.ws_app and self.ws_app.sock and self.ws_app.sock.connected:
                self.ws_app.send(message)
        except Exception as e:
            self.on_error(e)

    def _on_open(self, ws):
        self._attempts = 0
        self.on_open()

    def _on_message(self, ws, data):
        self.on_message(data)

    def _on_error(self, ws, err):
        self.on_error(err)

    def _on_close(self, ws, code, reason):
        if not self._stop_evt.is_set():
            self._attempts += 1
            if self._attempts <= self.max_retries:
                print("Try to reconnect:", self._attempts)
            if self._attempts > self.max_retries:
                try:
                    ws.keep_running = False
                except Exception:
                    pass
        self.on_close()

    def on_open(self):
        print("Opened")

    def on_message(self, data):
        try:
            obj = json.loads(data)
            print("Received(JSON):", obj)
        except Exception:
            print("Received(raw):", data)

    def on_error(self, err):
        print("Error:", err)

    def on_close(self):
        print("Closed")


if __name__ == "__main__":
    ws = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1",
        ping_interval=30,
        ping_timeout=10,
        max_retries=3,     
        retry_sleep=2.0    
    )

    try:
        time.sleep(150)
    finally:
        ws.close()
        ws.join(timeout=3)
```

위 코드 실행시 이제 120초 이후에도 연결이 해제되지 않고, 설정한바와 같이 150초 이후 연결을 클라이언트 요구로 정상 종료하는 것을 확인할 수 있습니다.

<br />

### 실시간 스트림 데이터 구독 요청 예제

위 예제는 단순히 업비트 WebSocket 서버로의 연결만 생성한 예제로, 실제 데이터를 구독하기 위해서는 구독 요청 메세지를 전송해야 합니다. [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide)를 참고하여 예제 코드를 아래와 같이 수정하였습니다.

* main 함수를 수정하여 연결 생성 이후 호가(orderbook) 데이터를 구독하기  위한 요청 메세지를 만들고, send\_message 메서드를 호출하여 전송합니다. 약 5초간 스트림을 수신하다가, 구독 데이터를 추가하기 위해 체결(trade) 요청 메세지를 추가한 뒤 다시 전송하여 두 종류의 데이터 스트림을 함께 수신할 수 있습니다. 약 5초간 수신을 지속하다가 close()를 호출하여 연결을 종료하는 예시입니다. (151-166 라인)
* 확인을 위해 \[on *message] 메서드를 수정하여 수신한 데이터의 type과 pair 코드를 기반으로  \[{type}*{pair}.jsonl] 파일에 수신한 json 데이터를 append 하도록 변경하였습니다. 코드를 실행하면, 실행한 디렉토리 하위에 \[KRW-BTC\_orderbook.jsonl], \[KRW-BTC\_trade.jsonl] 파일이 생성되며, 전송된 데이터가 기록되어 있는 것을 확인할 수 있습니다. (101-133 라인)

```python
import threading
import websocket
import json
import time
import uuid
import os

class ThreadedWebSocketApp(threading.Thread):
    def __init__(self, url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0):
        # type: (str, int, int, int, float) -> None
        threading.Thread.__init__(self)
        self.daemon = True
        self.url = url
        self.ping_interval = ping_interval
        self.ping_timeout = ping_timeout
        self.max_retries = max_retries
        self.retry_sleep = retry_sleep

        self.ws_app = None
        self._stop_evt = threading.Event()

    @staticmethod
    def connect(url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0):
        # type: (str, int, int, int, float) -> "ThreadedWebSocketApp"
        t = ThreadedWebSocketApp(
            url,
            ping_interval=ping_interval,
            ping_timeout=ping_timeout,
            max_retries=max_retries,
            retry_sleep=retry_sleep
        )
        t.start()
        return t

    def run(self):
        attempts = 0
        while not self._stop_evt.is_set():
            self._opened_once = False
            self.ws_app = websocket.WebSocketApp(
                self.url,
                on_open=self._on_open,
                on_message=self._on_message,
                on_error=self._on_error,
                on_close=self._on_close
            )
                
            self.ws_app.run_forever(
                ping_interval=self.ping_interval,
                ping_timeout=self.ping_timeout,
                reconnect=int(self.retry_sleep) if self.retry_sleep else None
            )
            self.ws_app = None

    def close(self):
        self._stop_evt.set()
        try:
            if self.ws_app:
                self.ws_app.close()
                try:
                    self.ws_app.keep_running = False
                except Exception:
                    pass
        except Exception:
            pass

    def send_message(self, message):
        try:
            if self.ws_app and self.ws_app.sock and self.ws_app.sock.connected:
                self.ws_app.send(message)
        except Exception as e:
            self.on_error(e)

    def _on_open(self, ws):
        self._attempts = 0
        print("Opened")

    def _on_message(self, ws, data):
        if isinstance(data, bytes):
            try:
                data = data.decode("utf-8")
            except Exception:
                print("Received(binary): <{} bytes>".format(len(data)))
                return

        try:
            obj = json.loads(data)
        except Exception:
            print("Received(raw):", data)
            self._append_jsonl("misc.jsonl", {"raw": data})
            return

        code = obj.get("code")
        mtype = obj.get("type")
        if code and mtype:
            path = "{0}_{1}.jsonl".format(code, mtype)
        else:
            path = "misc.jsonl"
        self._append_jsonl(path, obj)

    def _on_error(self, ws, err):
        print("Error:", err)

    def _on_close(self, ws, code, reason):
        if not self._stop_evt.is_set():
            self._attempts += 1
            if self._attempts <= self.max_retries:
                print("Try to reconnect:", self._attempts)
            if self._attempts > self.max_retries:
                try:
                    ws.keep_running = False
                except Exception:
                    pass
        print("Closed")

    def on_open(self):
        print("Opened")

    def _append_jsonl(self, path, obj):
        try:
            d = os.path.dirname(path)
            if d and not os.path.exists(d):
                os.makedirs(d)
            with open(path, "a", encoding="utf-8") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")
        except TypeError:
            with open(path, "a") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")

if __name__ == "__main__":
    ws = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1",
        ping_interval=30,
        ping_timeout=10,
        max_retries=3,     
        retry_sleep=2.0    
    )

    try:
        time.sleep(1)
        sub_orderbook = [
            {"ticket": str(uuid.uuid4())},
            {"type": "orderbook", "codes": ["KRW-BTC.1"]}
        ]
        ws.send_message(json.dumps(sub_orderbook))
        time.sleep(5)

        sub_trade = [
            {"ticket": str(uuid.uuid4())},
            {"type": "orderbook", "codes": ["KRW-BTC.1"]},
            {"type": "trade", "codes": ["KRW-BTC"]}
        ]
        ws.send_message(json.dumps(sub_trade))
        time.sleep(5)
    finally:
        ws.close()
        ws.join(timeout=3)
```

<br />

### 인증

인증이 필요한 구독 데이터(My Order 또는 My Asset) 요청 시 /private url로 WebSocket을 새로 연결하는 예제입니다.

* API Key의 Access Key, Secret Key를 설정할 수 있도록 지원 (25-26 라인, 39-40 라인, 178-179라인)
* 인증을 위한 JWT 토큰을 생성하고 인증 헤더로 추가하는 부분 추가 구현 (45-62 라인)
* /private 채널 연결 및 myOrder 타입 데이터 구독 추가 (176-180 라인, 191-195 라인)

```python
import threading
import websocket
import json
import time
import uuid
import os
import jwt

class ThreadedWebSocketApp(threading.Thread):
    def __init__(self, url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0, 
                access_key=None, secret_key=None):
        # type: (str, int, int, int, float, str, str) -> None
        threading.Thread.__init__(self)
        self.daemon = True
        self.url = url
        self.ping_interval = ping_interval
        self.ping_timeout = ping_timeout
        self.max_retries = max_retries
        self.retry_sleep = retry_sleep

        self.ws_app = None
        self._stop_evt = threading.Event()
        
        self.access_key = access_key
        self.secret_key = secret_key

    @staticmethod
    def connect(url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0,
                access_key=None, secret_key=None):
        # type: (str, int, int, int, float, str, str) -> "ThreadedWebSocketApp"
        t = ThreadedWebSocketApp(
            url,
            ping_interval=ping_interval,
            ping_timeout=ping_timeout,
            max_retries=max_retries,
            retry_sleep=retry_sleep,
            access_key=access_key,
            secret_key=secret_key,
        )
        t.start()
        return t
    
    def _create_jwt_token(self):
        # type: () -> str
        if not self.access_key or not self.secret_key:
            return None
        payload = {
            "access_key": self.access_key,
            "nonce": str(uuid.uuid4())
        }
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _build_headers(self):
        # type: () -> list
        headers = []
        token = self._create_jwt_token()
        if token:
            headers.append("Authorization: Bearer {0}".format(token))
        return headers

    def run(self):
        attempts = 0
        while not self._stop_evt.is_set():
            self._opened_once = False
            self.ws_app = websocket.WebSocketApp(
                self.url,
                header=self._build_headers(),
                on_open=self._on_open,
                on_message=self._on_message,
                on_error=self._on_error,
                on_close=self._on_close
            )
                
            self.ws_app.run_forever(
                ping_interval=self.ping_interval,
                ping_timeout=self.ping_timeout,
                reconnect=int(self.retry_sleep) if self.retry_sleep else None
            )

            self.ws_app = None

    def close(self):
        self._stop_evt.set()
        try:
            if self.ws_app:
                self.ws_app.close()
                try:
                    self.ws_app.keep_running = False
                except Exception:
                    pass
        except Exception:
            pass

    def send_message(self, message):
        try:
            if self.ws_app and self.ws_app.sock and self.ws_app.sock.connected:
                self.ws_app.send(message)
        except Exception as e:
            self.on_error(e)

    def _on_open(self, ws):
        self._attempts = 0
        print("Opened")

    def _on_message(self, ws, data):
        if isinstance(data, bytes):
            try:
                data = data.decode("utf-8")
            except Exception:
                print("Received(binary): <{} bytes>".format(len(data)))
                return

        try:
            obj = json.loads(data)
        except Exception:
            print("Received(raw):", data)
            self._append_jsonl("misc.jsonl", {"raw": data})
            return

        code = obj.get("code")
        mtype = obj.get("type")
        if code and mtype:
            path = "{0}_{1}.jsonl".format(code, mtype)
        else:
            path = "misc.jsonl"
        self._append_jsonl(path, obj)

    def _on_error(self, ws, err):
        print("Error:", err)

    def _on_close(self, ws, code, reason):
        if not self._stop_evt.is_set():
            self._attempts += 1
            if self._attempts <= self.max_retries:
                print("Try to reconnect:", self._attempts)
            if self._attempts > self.max_retries:
                try:
                    ws.keep_running = False
                except Exception:
                    pass
        print("Closed")

    def _append_jsonl(self, path, obj):
        try:
            d = os.path.dirname(path)
            if d and not os.path.exists(d):
                os.makedirs(d)
            with open(path, "a", encoding="utf-8") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")
        except TypeError:
            with open(path, "a") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")

if __name__ == "__main__":

    ws = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1",
    )
    
    ws_private = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1/private",
        access_key="YOUR_ACCESS_KEY",          
        secret_key="YOUR_SECRET_KEY"  
    )

    try:
        time.sleep(1)
        sub_orderbook = [
            {"ticket": str(uuid.uuid4())},
            {"type": "orderbook", "codes": ["KRW-BTC.1"]}
        ]
        ws.send_message(json.dumps(sub_orderbook))
        time.sleep(5)

        sub_trade = [
            {"ticket": str(uuid.uuid4())},
            {"type": "myOrder", "codes": ["KRW-BTC"]}
        ]
        ws_private.send_message(json.dumps(sub_trade))
        time.sleep(5)
        
    finally:
        ws.close()
        ws.join(timeout=3)
        ws_private.close()
        ws_private.join(timeout=3)
```

<br />

### 요청 수 제한(Rate Limit) 관련 처리

[요청 수 제한(Rate Limits)](https://docs.upbit.com/kr/reference/rate-limits) 페이지에 명시된 바와 같이, 업비트 WebSocket은 초당 최대 5회의 연결 요청 횟수 제한과, 초당 최대 5회/분당 최대 100회의 최대 메세지 전송 횟수 제한 정책을 적용하고 있습니다. 이를 위해 과도하게 구독 요청 메세지가 전송되지 않도록, WebSocket 메세지 전송을 위한 Rate Limiter를 구현하여 적용할 수 있습니다.

* \_FixedWindowLimiter 클래스를 생성하고, 초당 최대 5회 / 분당 최대 100회 이상 요청이 전송되지 않도록 구현하였습니다. acquire 메서드를 호출하면 메세지 전송 가능 횟수를 차감하거나, 최대 호출 횟수를 초과한 경우 다음 시간 window(초 또는 분)까지 대기합니다. (9-48라인)
* send\_message 메서드 내부에서 메세지 전송 전 self.\_send\_limiter.acquire() 를 호출하여 전송 가능 여부를 확인합니다. (144 라인)

```
import threading
import websocket
import json
import time
import uuid
import os
import jwt

class _FixedWindowLimiter(object):
    def __init__(self, per_sec=5, per_min=100):
        # type: (int, int) -> None
        self.per_sec = per_sec
        self.per_min = per_min
        self._sec_ts = 0
        self._min_ts = 0
        self._sec_used = 0
        self._min_used = 0

    def acquire(self):
        now = time.time()
        sec = int(now)
        minute = int(now // 60)

        if sec != self._sec_ts:
            self._sec_ts = sec
            self._sec_used = 0
        if minute != self._min_ts:
            self._min_ts = minute
            self._min_used = 0

        if self._sec_used >= self.per_sec:
            sleep_for = (self._sec_ts + 1) - now + 0.001
            if sleep_for > 0:
                time.sleep(sleep_for)
            now = time.time()
            self._sec_ts = int(now)
            self._sec_used = 0

        if self._min_used >= self.per_min:
            sleep_for = ((self._min_ts + 1) * 60) - now + 0.001
            if sleep_for > 0:
                time.sleep(sleep_for)
            now = time.time()
            self._min_ts = int(now // 60)
            self._min_used = 0

        self._sec_used += 1
        self._min_used += 1

class ThreadedWebSocketApp(threading.Thread):
    def __init__(self, url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0, 
                access_key=None, secret_key=None, send_per_sec=5, send_per_min=100):
        # type: (str, int, int, int, float, str, str, int, int) -> None
        threading.Thread.__init__(self)
        self.daemon = True
        self.url = url
        self.ping_interval = ping_interval
        self.ping_timeout = ping_timeout
        self.max_retries = max_retries
        self.retry_sleep = retry_sleep

        self.ws_app = None
        self._stop_evt = threading.Event()
        
        self.access_key = access_key
        self.secret_key = secret_key
        self._send_limiter = _FixedWindowLimiter(per_sec=send_per_sec, per_min=send_per_min)

    @staticmethod
    def connect(url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0,
                access_key=None, secret_key=None, send_per_sec=5, send_per_min=100):
        # type: (str, int, int, int, float, str, str, int, int) -> "ThreadedWebSocketApp"
        t = ThreadedWebSocketApp(
            url,
            ping_interval=ping_interval,
            ping_timeout=ping_timeout,
            max_retries=max_retries,
            retry_sleep=retry_sleep,
            access_key=access_key,
            secret_key=secret_key,
            send_per_sec=send_per_sec,
            send_per_min=send_per_min
        )
        t.start()
        return t
    
    def _create_jwt_token(self):
        # type: () -> str
        if not self.access_key or not self.secret_key:
            return None
        payload = {
            "access_key": self.access_key,
            "nonce": str(uuid.uuid4())
        }
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _build_headers(self):
        # type: () -> list
        headers = []
        token = self._create_jwt_token()
        if token:
            headers.append("Authorization: Bearer {0}".format(token))
        return headers

    def run(self):
        attempts = 0
        while not self._stop_evt.is_set():
            self._opened_once = False
            self.ws_app = websocket.WebSocketApp(
                self.url,
                header=self._build_headers(),
                on_open=self._on_open,
                on_message=self._on_message,
                on_error=self._on_error,
                on_close=self._on_close
            )
                
            self.ws_app.run_forever(
                ping_interval=self.ping_interval,
                ping_timeout=self.ping_timeout,
                reconnect=int(self.retry_sleep) if self.retry_sleep else None
            )

            self.ws_app = None

    def close(self):
        self._stop_evt.set()
        try:
            if self.ws_app:
                self.ws_app.close()
                try:
                    self.ws_app.keep_running = False
                except Exception:
                    pass
        except Exception:
            pass

    def send_message(self, message):
        try:
            if self.ws_app and self.ws_app.sock and self.ws_app.sock.connected:
                self._send_limiter.acquire()
                self.ws_app.send(message)
        except Exception as e:
            self.on_error(e)

    def _on_open(self, ws):
        self._attempts = 0
        print("Opened")

    def _on_message(self, ws, data):
        if isinstance(data, bytes):
            try:
                data = data.decode("utf-8")
            except Exception:
                print("Received(binary): <{} bytes>".format(len(data)))
                return

        try:
            obj = json.loads(data)
        except Exception:
            print("Received(raw):", data)
            self._append_jsonl("misc.jsonl", {"raw": data})
            return

        code = obj.get("code")
        mtype = obj.get("type")
        if code and mtype:
            path = "{0}_{1}.jsonl".format(code, mtype)
        else:
            path = "misc.jsonl"
        self._append_jsonl(path, obj)

    def _on_error(self, ws, err):
        print("Error:", err)

    def _on_close(self, ws, code, reason):
        if not self._stop_evt.is_set():
            self._attempts += 1
            if self._attempts <= self.max_retries:
                print("Try to reconnect:", self._attempts)
            if self._attempts > self.max_retries:
                try:
                    ws.keep_running = False
                except Exception:
                    pass
        print("Closed")

    def _append_jsonl(self, path, obj):
        try:
            d = os.path.dirname(path)
            if d and not os.path.exists(d):
                os.makedirs(d)
            with open(path, "a", encoding="utf-8") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")
        except TypeError:
            with open(path, "a") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")

if __name__ == "__main__":

    ws = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1",
    )
    
    ws_private = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1/private",
        access_key="YOUR_ACCESS_KEY",          
        secret_key="YOUR_SECRET_KEY"  
    )

    try:
        time.sleep(1)
        sub_orderbook = [
            {"ticket": str(uuid.uuid4())},
            {"type": "orderbook", "codes": ["KRW-BTC.1"]}
        ]
        ws.send_message(json.dumps(sub_orderbook))
        time.sleep(5)

        sub_trade = [
            {"ticket": str(uuid.uuid4())},
            {"type": "myOrder", "codes": ["KRW-BTC"]}
        ]
        ws_private.send_message(json.dumps(sub_trade))
        time.sleep(5)
        
    finally:
        ws.close()
        ws.join(timeout=3)
        ws_private.close()
        ws_private.join(timeout=3)
```

<br />

### SIMPLE 버전의 스트림 구독

마지막으로, SIMPLE 포맷의 데이터를 구독하기 위한 처리를 반영할 수 있습니다.

* on\_message 메서드의 데이터 파싱 부분에 type의 축약형인 ty, code의 축약형인 cd를 기반으로 데이터를 처리하는 로직을 추가하였습니다. (200-201 라인)
* 구독 메세지 내에 {"format": "SIMPLE"}을 추가하여 전송 데이터 포맷을 SIMPLE 타입으로 변경할 수 있습니다. (232 라인)

```python
import threading
import websocket
import json
import time
import uuid
import os
import jwt

class _FixedWindowLimiter(object):
    def __init__(self, per_sec=5, per_min=100):
        # type: (int, int) -> None
        self.per_sec = per_sec
        self.per_min = per_min
        self._sec_ts = 0
        self._min_ts = 0
        self._sec_used = 0
        self._min_used = 0

    def acquire(self):
        now = time.time()
        sec = int(now)
        minute = int(now // 60)

        if sec != self._sec_ts:
            self._sec_ts = sec
            self._sec_used = 0
        if minute != self._min_ts:
            self._min_ts = minute
            self._min_used = 0

        if self._sec_used >= self.per_sec:
            sleep_for = (self._sec_ts + 1) - now + 0.001
            if sleep_for > 0:
                time.sleep(sleep_for)
            now = time.time()
            self._sec_ts = int(now)
            self._sec_used = 0

        if self._min_used >= self.per_min:
            sleep_for = ((self._min_ts + 1) * 60) - now + 0.001
            if sleep_for > 0:
                time.sleep(sleep_for)
            now = time.time()
            self._min_ts = int(now // 60)
            self._min_used = 0

        self._sec_used += 1
        self._min_used += 1

class ThreadedWebSocketApp(threading.Thread):
    def __init__(self, url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0, 
                access_key=None, secret_key=None, send_per_sec=5, send_per_min=100):
        # type: (str, int, int, int, float, str, str, int, int) -> None
        threading.Thread.__init__(self)
        self.daemon = True
        self.url = url
        self.ping_interval = ping_interval
        self.ping_timeout = ping_timeout
        self.max_retries = max_retries
        self.retry_sleep = retry_sleep

        self.ws_app = None
        self._stop_evt = threading.Event()
        
        self.access_key = access_key
        self.secret_key = secret_key
        self._send_limiter = _FixedWindowLimiter(per_sec=send_per_sec, per_min=send_per_min)

    @staticmethod
    def connect(url, ping_interval=30, ping_timeout=10,
                max_retries=3, retry_sleep=2.0,
                access_key=None, secret_key=None, send_per_sec=5, send_per_min=100):
        # type: (str, int, int, int, float, str, str, int, int) -> "ThreadedWebSocketApp"
        t = ThreadedWebSocketApp(
            url,
            ping_interval=ping_interval,
            ping_timeout=ping_timeout,
            max_retries=max_retries,
            retry_sleep=retry_sleep,
            access_key=access_key,
            secret_key=secret_key,
            send_per_sec=send_per_sec,
            send_per_min=send_per_min
        )
        t.start()
        return t
    
    def _create_jwt_token(self):
        # type: () -> str
        if not self.access_key or not self.secret_key:
            return None
        payload = {
            "access_key": self.access_key,
            "nonce": str(uuid.uuid4())
        }
        token = jwt.encode(payload, self.secret_key, algorithm="HS512")
        return token if isinstance(token, str) else token.decode("utf-8")

    def _build_headers(self):
        # type: () -> list
        headers = []
        token = self._create_jwt_token()
        if token:
            headers.append("Authorization: Bearer {0}".format(token))
        return headers

    def run(self):
        attempts = 0
        while not self._stop_evt.is_set():
            self._opened_once = False
            self.ws_app = websocket.WebSocketApp(
                self.url,
                header=self._build_headers(),
                on_open=self._on_open,
                on_message=self._on_message,
                on_error=self._on_error,
                on_close=self._on_close
            )
                
            self.ws_app.run_forever(
                ping_interval=self.ping_interval,
                ping_timeout=self.ping_timeout,
                reconnect=int(self.retry_sleep) if self.retry_sleep else None
            )

            self.ws_app = None

    def close(self):
        self._stop_evt.set()
        try:
            if self.ws_app:
                self.ws_app.close()
                try:
                    self.ws_app.keep_running = False
                except Exception:
                    pass
        except Exception:
            pass

    def send_message(self, message):
        try:
            if self.ws_app and self.ws_app.sock and self.ws_app.sock.connected:
                self._send_limiter.acquire()
                self.ws_app.send(message)
        except Exception as e:
            self.on_error(e)

    def _on_open(self, ws):
        self._attempts = 0
        print("Opened")

    def _on_message(self, ws, data):
        if isinstance(data, bytes):
            try:
                data = data.decode("utf-8")
            except Exception:
                print("Received(binary): <{} bytes>".format(len(data)))
                return

        try:
            obj = json.loads(data)
        except Exception:
            print("Received(raw):", data)
            self._append_jsonl("misc.jsonl", {"raw": data})
            return

        code = obj.get("code")
        mtype = obj.get("type")
        if code and mtype:
            path = "{0}_{1}.jsonl".format(code, mtype)
        else:
            path = "misc.jsonl"
        self._append_jsonl(path, obj)

    def _on_error(self, ws, err):
        print("Error:", err)

    def _on_close(self, ws, code, reason):
        if not self._stop_evt.is_set():
            self._attempts += 1
            if self._attempts <= self.max_retries:
                print("Try to reconnect:", self._attempts)
            if self._attempts > self.max_retries:
                try:
                    ws.keep_running = False
                except Exception:
                    pass
        print("Closed")

    def _append_jsonl(self, path, obj):
        try:
            d = os.path.dirname(path)
            if d and not os.path.exists(d):
                os.makedirs(d)
            with open(path, "a", encoding="utf-8") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")
        except TypeError:
            with open(path, "a") as f:
                f.write(json.dumps(obj, ensure_ascii=False) + "\n")

if __name__ == "__main__":

    ws = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1",
    )
    
    ws_private = ThreadedWebSocketApp.connect(
        url="wss://api.upbit.com/websocket/v1/private",
        access_key="YOUR_ACCESS_KEY",          
        secret_key="YOUR_SECRET_KEY"  
    )

    try:
        time.sleep(1)
        sub_orderbook = [
            {"ticket": str(uuid.uuid4())},
            {"type": "orderbook", "codes": ["KRW-BTC.1"]},
            {"format": "SIMPLE"}
        ]
        ws.send_message(json.dumps(sub_orderbook))
        time.sleep(5)

        sub_trade = [
            {"ticket": str(uuid.uuid4())},
            {"type": "myOrder", "codes": ["KRW-BTC"]}
        ]
        ws_private.send_message(json.dumps(sub_trade))
        time.sleep(5)
        
    finally:
        ws.close()
        ws.join(timeout=3)
        ws_private.close()
        ws_private.join(timeout=3)
```

<br />

## 마치며

업비트 WebSocket의 운영 정책을 준수하면서 안정적인 연동환경을 구축하기 위한 Best Practice 구현을 간단한 동기 클라이언트 예제를 통해 살펴보았습니다. 본 가이드에서 설명한 내용 외에도 [REST API 연동 Best Practice](https://docs.upbit.com/kr/docs/rest-api-best-practice) 에서 소개된 것과 같이 Logging 등의 모듈을 추가하는 등의 개선을 반영할 수 있습니다. 아래 페이지 중 하나로 이동하여, WebSocket에 대한 더 많은 정보를 확인하세요.

# Sibling pages

* [REST API 연동 Best Practice](https://docs.upbit.com/kr/docs/rest-api-best-practice.md)