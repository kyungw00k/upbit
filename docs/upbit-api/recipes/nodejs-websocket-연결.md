---
updatedAt: 2026-01-07T10:53:35.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Node.js] Websocket 연결

```javascript JavaScript
const jwt = require("jsonwebtoken");
const {v4: uuidv4} = require('uuid');
const WebSocket = require("ws");

const payload = {
    access_key: "<YOUR_ACCESS_KEY>", 
    nonce: uuidv4(),
};

/**
 * JWT 생성
 */
const jwtToken = jwt.sign(payload, "<YOUR_SECRET_KEY>");

/**
 * 설정된 endpoint에 웹소켓 연경
 */
const ws = new WebSocket("wss://api.upbit.com/websocket/v1/private", {
    headers: {
        authorization: `Bearer ${jwtToken}`
    }
});

/**
 * 요청 이벤트 정의, 웹소켓 연결 후 요청을 전송
 */
ws.on("open", () => {
    console.log("connected!");
    ws.send('[{"ticket":"test example"},{"type":"myAsset"}]');

});

/**
 * 에러 이벤트 정의
 */
ws.on("error", console.error);

/**
 * 사용자게 보여줄 메시지 포맷 정의
 */
ws.on("message", (data) => console.log(data.toString()));

/**
 * 연결 종료 이벤트 정의
 */
ws.on("close", () => console.log("closed!"));
```

```json Response Example
connected!
{"type":"myAsset","asset_uuid":"<my_asset_uuid>","assets":[{"currency":"KRW","balance":20517.4157543899035,"locked":15007.5}],"asset_timestamp":1753769511603,"timestamp":1753769511612,"stream_type":"REALTIME"}
{"type":"myAsset","asset_uuid":"<my_asset_uuid>","assets":[{"currency":"XRP","balance":1.15580212,"locked":0}],"asset_timestamp":1753769511603,"timestamp":1753769511613,"stream_type":"REALTIME"}
closed!
```

# 유틸 라이브러리 Import

<!-- javascript@1-3 -->

기능 구현을 위해 필요한 모듈을 import 합니다. 별도의 설치가 필요한 모듈의 경우 `pip install <module name>` 명령어를 실행해 설치할 수 있습니다. 

websocket의 경우, npm install ws 명령어를 실행해 websocket 패키지를 설치해야 합니다.

# Websocket 정의 및 실행

<!-- javascript@5-22 -->

Websocket 요청 시 필요한 payload와 JWT를 생성하고 이를 헤더에 추가해 Websocket을 연결합니다.

# 요청 이벤트 정의

<!-- javascript@24-31 -->

Websocket이 연결된 후 수행할 동작을 정의합니다. 이 가이드에서는 Websocket 연결 이후, 사용자의 주문 및 체결 데이터를 수신하기 위한 요청 메세지를 전송합니다.

# 에러 이벤트 정의

<!-- javascript@33-36 -->

Websocket 연결 중 에러 발생 시 처리 방법을 정의합니다. 이 가이드에서는 에러 발생 시 사용자에게 전체 에러 데이터를 노출합니다.

# 메시지 이벤트  정의

<!-- javascript@38-41 -->

WebSocket으로 전달된 실시간 스트림 데이터를 콘솔에 출력하는 예제입니다.

# 연결 종료 이벤트 정의

<!-- javascript@43-46 -->

Websocket 연결 종료 후 실행할 동작을 정의합니다. 이 가이드에서는 연결 종료 시 “closed!” 라는 문자열을 출력합니다.