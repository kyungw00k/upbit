---
updatedAt: 2026-05-04T05:48:42.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# WebSocket으로 실시간 캔들 차트 생성

업비트 시세 WebSocket 실시간 캔들 스트림을 구독하여 캔들 차트를 생성합니다.

[block:html]
{
  "html": "<meta name=\"robots\" content=\"noindex\">"
}
[/block]

## 시작하기

시세 데이터로 가격 및 거래의 변동 추이를 관찰하기 위해서는 주기적, 지속적인 데이터 수집이 필요합니다. 반복적으로 REST API를 호출하는 대신, WebSocket 연결을 통해 실시간 데이터를 스트림 형태로 구독할 수 있습니다. 이 튜토리얼에서는 WebSocket으로 실시간 캔들 데이터를 구독하고 시각화 차트를 생성합니다.

업비트 기본적인 WebSocket 사용 방법과 언어별 예제는 [WebSocket 사용 및 에러 안내](https://docs.upbit.com/kr/reference/websocket-guide) 에서 확인하실 수 있습니다.

<br />

## 개발 환경 및 사용 라이브러리 안내

본 튜토리얼은 Python 언어를 사용하는 예제로 진행됩니다. [개발 환경 설정 가이드](https://docs.upbit.com/kr/docs/dev-environment)를 참고하여 Python 개발 환경을 설정해주세요. WebSocket 수신을 위해  `websocket-client` 라이브러리를, 차트 시각화를 위해 `pandas` , `matplotlib` 라이브러리를 사용합니다. 가상 환경에서 아래 명령어를 실행하여 해당 라이브러리를 설치해주세요.

```
pip install websocket-client matplotlib pandas
```

<br />

## WebSocket으로 실시간 캔들 스트림 구독

먼저, python으로 WebSocket 실시간 캔들 스트림을 구독하는 코드를 작성해보겠습니다. run\_ws() 함수가 실행되면 WebSocketApp을 생성하여 업비트 WebSocket 서버에 연결합니다. threading 기능을 사용하여 별도의 스레드에서 WebSocket을 연결하도록 하였습니다. on\_open 함수를 정의하여 WebSocket이 성공적으로 연결된 후 `candle.1s` 타입의 실시간 스트림을 구독 요청하는 메세지를 전송합니다. 캔들 스트림 구독 요청 메세지의 명세는 [캔들 (Candle)](https://docs.upbit.com/kr/reference/websocket-candle) 문서에서 확인할 수 있습니다. 본 예제에서는 KRW-BTC 페어를 구독합니다.

```python
import websocket
import threading
import uuid
import json
import time

def on_open(ws):
    print("WebSocket connected")
    subscribe_message = [
        {"ticket": str(uuid.uuid4())},
        {"type": "candle.1s", "codes": ["KRW-BTC"]}
    ]
    ws.send(json.dumps(subscribe_message))
    
def on_message(ws, message):
    try:
        obj = json.loads(message)
        print(obj)
    except Exception:
        print("Received(raw):", message)
        return
    
def on_error(err):
        print("Error:", err)

def on_close():
        print("Closed")

def run_ws():
    ws = websocket.WebSocketApp(
        "wss://api.upbit.com/websocket/v1",
        on_open=on_open,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close
    )
    ws.run_forever(ping_interval=30, ping_timeout=10, reconnect=2)

ws_thread = threading.Thread(target=run_ws, daemon=True)
ws_thread.start()

try:
    while True:
        time.sleep(1)
except KeyboardInterrupt:
    pass
```

<br />

구독 메세지가 전송되면 초캔들 데이터의 수신이 시작됩니다. 메세지 수신 시 차트 시각화에 필요한 데이터를 추출하고 시각화 데이터에 반영하기 위해 on\_message(ws, message) 함수를 아래와 같이 새로 정의한 뒤, 연결 요청에 포함합니다. 해당 함수는 메세지를 수신할 때마다 실행됩니다.

```python
import websocket
import threading
import uuid
import json
import time
import pandas as pd

candle_data_list = []
data_lock = threading.Lock()

def on_open(ws):
    print("WebSocket connected")
    subscribe_message = [
        {"ticket": str(uuid.uuid4())},
        {"type": "candle.1s", "codes": ["KRW-BTC"]}
    ]
    ws.send(json.dumps(subscribe_message))
    
def on_message(ws, message):
    global candle_data_list
    data = json.loads(message)

    if data.get('type') == 'candle.1s':
        candle_time = data['candle_date_time_kst']
        opening_price = float(data['opening_price'])
        high_price = float(data['high_price'])
        low_price = float(data['low_price'])
        close_price = float(data['trade_price'])
        volume = float(data['candle_acc_trade_volume'])

        with data_lock:
            timestamp = pd.to_datetime(candle_time)
            candle_data_list.append({
                'timestamp': timestamp,
                'open': opening_price,
                'high': high_price,
                'low': low_price,
                'close': close_price,
                'volume': volume
            })
            if len(candle_data_list) > 50:
                candle_data_list = candle_data_list[-50:]
    
def on_error(err):
        print("Error:", err)

def on_close():
        print("Closed")

def run_ws():
    ws = websocket.WebSocketApp(
        "wss://api.upbit.com/websocket/v1",
        on_open=on_open,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close
    )
    ws.run_forever(ping_interval=30, ping_timeout=10, reconnect=2)

ws_thread = threading.Thread(target=run_ws, daemon=True)
ws_thread.start()

try:
    while True:
        time.sleep(1)
except KeyboardInterrupt:
    pass
```

위 코드 예제에서는 차트에 반영할 데이터를 저장할 `candle_data_list` 를 정의하였습니다. 그런 다음, 매 메세지를 수신할 때마다 json 형식의 메세지로부터 캔들의 주요 데이터 필드들을 추출하여 해당 리스트에 넣도록 구현하였습니다.

이제 해당 데이터를 기반으로 차트를 그리는 부분을 추가합니다. update\_chart() 함수는 candle\_data\_list에 저장된 데이터를 읽어온 뒤 해당 데이터를 기반으로 차트를 갱신합니다. 해당 함수를 포함한 전체 예제 코드는 아래와 같습니다.

<br />

## 전체 예제 코드

```python
import websocket
import threading
import uuid
import json
import time
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation

candle_data_list = []
data_lock = threading.Lock()

def on_open(ws):
    print("WebSocket connected")
    subscribe_message = [
        {"ticket": str(uuid.uuid4())},
        {"type": "candle.1s", "codes": ["KRW-BTC"]}
    ]
    ws.send(json.dumps(subscribe_message))
    
def on_message(ws, message):
    global candle_data_list
    data = json.loads(message)

    if data.get('type') == 'candle.1s':
        candle_time = data['candle_date_time_kst']
        opening_price = float(data['opening_price'])
        high_price = float(data['high_price'])
        low_price = float(data['low_price'])
        close_price = float(data['trade_price'])
        volume = float(data['candle_acc_trade_volume'])

        with data_lock:
            timestamp = pd.to_datetime(candle_time)
            candle_data_list.append({
                'timestamp': timestamp,
                'open': opening_price,
                'high': high_price,
                'low': low_price,
                'close': close_price,
                'volume': volume
            })
            if len(candle_data_list) > 50:
                candle_data_list = candle_data_list[-50:]
    
def on_error(err):
        print("Error:", err)

def on_close():
        print("Closed")

def run_ws():
    ws = websocket.WebSocketApp(
        "wss://api.upbit.com/websocket/v1",
        on_open=on_open,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close
    )
    ws.run_forever(ping_interval=30, ping_timeout=10, reconnect=2)

def update_chart(frame):
    with data_lock:
        if len(candle_data_list) < 2:
            return

        data_copy = candle_data_list.copy()
    
    ax1.clear()
    ax2.clear()
    
    timestamps = [d['timestamp'] for d in data_copy]
    opens = [d['open'] for d in data_copy]
    highs = [d['high'] for d in data_copy]
    lows = [d['low'] for d in data_copy]
    closes = [d['close'] for d in data_copy]
    volumes = [d['volume'] for d in data_copy]
    
    for i, (ts, o, h, l, c) in enumerate(zip(timestamps, opens, highs, lows, closes)):
        color = 'red' if c >= o else 'blue'
        ax1.plot([i, i], [o, c], color=color, linewidth=12)
        ax1.plot([i, i], [l, h], color=color, linewidth=2)
    
    ax1.set_title('Upbit KRW-BTC Real-time 1s Candlestick Chart', fontsize=14)
    ax1.set_ylabel('Price (KRW)', fontsize=12)
    ax1.grid(True, alpha=0.3)
    
    if len(timestamps) > 0:
        step = max(1, len(timestamps) // 10) 
        ax1.set_xticks(range(0, len(timestamps), step))
        ax1.set_xticklabels([ts.strftime('%H:%M:%S') for ts in timestamps[::step]], 
                            rotation=45, fontsize=10)
    
    ax2.bar(range(len(volumes)), volumes, color='gray', alpha=0.6)
    ax2.set_ylabel('Volume', fontsize=12)
    ax2.set_xlabel('Time', fontsize=12)
    ax2.grid(True, alpha=0.3)
    
    if len(timestamps) > 0:
        ax2.set_xticks(range(0, len(timestamps), step))
        ax2.set_xticklabels([ts.strftime('%H:%M:%S') for ts in timestamps[::step]], 
                            rotation=45, fontsize=10)
    
    plt.tight_layout()

ws_thread = threading.Thread(target=run_ws, daemon=True)
ws_thread.start()

fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10), gridspec_kw={'height_ratios': [3, 1]})

ani = FuncAnimation(fig, update_chart, interval=1000, cache_frame_data=False)
plt.show()

```

위 파일을 저장하고 실행하면 다음과 같은 순서로 프로그램이 동작합니다.

1. 새로운 스레드를 생성하여 업비트 public WebSocket 서버에 연결합니다. KRW-BTC 페어의 실시간 초 캔들 데이터 구독을 요청하는 메세지를 전송합니다.
2. 초 캔들 데이터를 수신할 때 마다 캔들 데이터 저장 리스트에 저장합니다. (최대 50개의 데이터 저장)
3. 메인 스레드에서 차트 애니메이션을 실행합니다. 실시간으로 수신된 캔들 데이터를 1초에 한 번 조회하여 차트에 반영합니다.

<br />

생성되는 초 캔들 차트 예시는 아래와 같습니다. 캔들 데이터는 해당 시간 구간 내에서 체결이 발생하는 경우에만 수신되므로 1초 내 거래가 없거나 단일 체결가로만 체결된 경우 캔들 차트가 반영되지 않을 수 있으며 이는 정상 동작입니다.

[block:image]
{
  "images": [
    {
      "image": [
        "https://files.readme.io/33ba975adf8137a2b491b05d6661e467f7734a1e6a7addb936396504300f1d34-_2025-08-11__6.29.09.png",
        "",
        ""
      ],
      "align": "center"
    }
  ]
}
[/block]

<br />

## 마치며

본 튜토리얼에서는 WebSocket을 이용하여 실시간 캔들 데이터를 수신하고, 해당 데이터로 실시간 차트를 생성하는 예제를 함께 구현해보았습니다. 예제를 변형하여 데이터를 파일로 저장하거나 실시간 분석 기능을 구현해 투자 전략 도출에 활용할 수 있습니다.

# Sibling pages

* [REST API를 이용한 시세 조회](https://docs.upbit.com/kr/docs/upbit-quotation-restful-api.md)
* [캔들 API로 지정한 기간의 시세 데이터를 CSV 파일로 출력하기](https://docs.upbit.com/kr/docs/how-to-download-candle-data.md)