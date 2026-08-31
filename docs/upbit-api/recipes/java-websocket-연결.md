---
updatedAt: 2026-01-07T10:53:35.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [Java] Websocket 연결

```java Java
package com.upbit.openapi.test;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import java.util.UUID;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;
import okio.ByteString;
import org.jetbrains.annotations.NotNull;

public class AuthWebSocket {

    /**
     * JWT Token 생성 메서드 예시
     */ 
    private static String createJwt(String accessKey, String secretKey) {
        try {
            Algorithm algorithm = Algorithm.HMAC256(secretKey);
            return JWT.create().withClaim("access_key", accessKey)
                .withClaim("nonce", UUID.randomUUID().toString()).sign(algorithm);
        } catch (Exception e) {
            throw new RuntimeException("JWT token generation failed", e);
        }
    }

    static WebSocketListener createWebSocketListener() {
        return new WebSocketListener() {
		    @Override
		    public void onOpen(@NotNull WebSocket webSocket, @NotNull Response response) {
		        System.out.println("connected!");
		
		        // 연결 이후 데이터 구독 요청
		        String subscribeMessage = "[{\"ticket\":\"test example\"},{\"type\":\"ticker\",\"codes\":[\"KRW-BTC\"]}]";
		        webSocket.send(subscribeMessage);
		    }
		
		    @Override
		    public void onMessage(@NotNull WebSocket webSocket, @NotNull String text) {
		        System.out.println("Received: " + text);
		    }
		
		    @Override
		    public void onMessage(@NotNull WebSocket webSocket, @NotNull ByteString bytes) {
		        System.out.println("Received: " + bytes.utf8());
		    }
		
		    @Override
		    public void onClosing(@NotNull WebSocket webSocket, int code, @NotNull String reason) {
		        System.out.println("Connection closing: " + reason);
		        webSocket.close(code, reason);
		    }
		
		    @Override
		    public void onClosed(@NotNull WebSocket webSocket, int code, @NotNull String reason) {
		        System.out.println("Connection closed: " + reason);
		    }
		
		    @Override
		    public void onFailure(@NotNull WebSocket webSocket, @NotNull Throwable t, Response response) {
		        System.err.println("Error occurred: " + t.getMessage());
		        if (response != null) {
		            System.err.println("Response: " + response);
		        }
		    };
    }

    public static void main(String[] args){

        String accessKey = "<YOUR_ACCESS_KEY>";
        String secretKey = "<YOUR_SECRET_KEY>"; // 실제로는 안전하게 로드하거나 주입하세요.
				OkHttpClient httpClient = new OkHttpClient.Builder()
            .pingInterval(30, TimeUnit.SECONDS) // ping interval 30초
            .connectTimeout(10, TimeUnit.SECONDS) // connect timeout
            .writeTimeout(10, TimeUnit.SECONDS) // write timeout
          .build();


        try {
            String jwtToken = createJwt(accessKey, secretKey);
          
						// /websocket/v1 endpoint로 요청시 인증 헤더는 필수가 아닙니다.
            // 인증 예제를 안내하고 데이터 수신을 확인하기 위해 임의로 시세 데이터를 구독하는 예제를 제공합니다.
            Request request = new Request.Builder()
                .url("wss://api.upbit.com/websocket/v1")
                .addHeader("Authorization", "Bearer " + jwtToken)
                .build();

            httpClient.newWebSocket(request, createWebSocketListener());
        } catch (Exception e) {
            throw new RuntimeException("Failed to connect to private WebSocket", e);
        }
    }

}

```

```json Response Example
connected!
{"type":"ticker","code":"KRW-BTC","opening_price":158219000.00000000,"high_price":159278000.00000000,"low_price":157588000.00000000,"trade_price":158201000.00000000,"prev_closing_price":158200000.00000000,"acc_trade_price":95078733694.3792000000000000,"change":"RISE","change_price":1000.00000000,"signed_change_price":1000.00000000,"change_rate":0.0000063211,"signed_change_rate":0.0000063211,"ask_bid":"BID","trade_volume":0.00006321,"acc_trade_volume":600.02894052,"trade_date":"20250926","trade_time":"082645","trade_timestamp":1758875205145,"acc_ask_volume":341.67985038,"acc_bid_volume":258.34909014,"highest_52_week_price":169900000.00000000,"highest_52_week_date":"2025-08-14","lowest_52_week_price":80635000.00000000,"lowest_52_week_date":"2024-10-10","market_state":"ACTIVE","is_trading_suspended":false,"delisting_date":null,"market_warning":"NONE","timestamp":1758875205276,"acc_trade_price_24h":331187169240.21526,"acc_trade_volume_24h":2083.82485091,"stream_type":"SNAPSHOT"}
{"type":"ticker","code":"KRW-BTC","opening_price":158219000.00000000,"high_price":159278000.00000000,"low_price":157588000.00000000,"trade_price":158200000.00000000,"prev_closing_price":158200000.00000000,"acc_trade_price":95088751392.9099200000000000,"change":"EVEN","change_price":0E-8,"signed_change_price":0E-8,"change_rate":0,"signed_change_rate":0,"ask_bid":"ASK","trade_volume":0.00011188,"acc_trade_volume":600.09226312,"trade_date":"20250926","trade_time":"082646","trade_timestamp":1758875206936,"acc_ask_volume":341.67996226,"acc_bid_volume":258.41230086,"highest_52_week_price":169900000.00000000,"highest_52_week_date":"2025-08-14","lowest_52_week_price":80635000.00000000,"lowest_52_week_date":"2024-10-10","market_state":"ACTIVE","is_trading_suspended":false,"delisting_date":null,"market_warning":"NONE","timestamp":1758875206974,"acc_trade_price_24h":331188513477.17233,"acc_trade_volume_24h":2083.83396972,"stream_type":"REALTIME"}
closed!
```

# 유틸 라이브러리 Import

<!-- java@1-12 -->

기능 구현을 위해 필요한 라이브러리를 import 합니다.

# JWT 생성

<!-- java@16-27 -->

Private Websocket 연결 시 인증을 위해 전달할 JWT를 생성하는 메서드를 구현합니다.

# 이벤트 리스너 객체 정의

<!-- java@29-68 -->

Websocket 연결 후 발생할 이벤트를 정의합니다. 

1. Websocket 연결이 성공한 경우 실시간 스트림 구독 요청 메세지를 전송합니다.
2. 수신한 메세지를 표준 출력으로 출력합니다.
3. Websocket 연결 종료 시 메시지와 종료 사유를 노출합니다.
4. 에러 발생 시 에러 메시지를 노출합니다.

# Websocket 연결

<!-- java@70-92 -->

Websocket 요청 시 필요한 JWT를 생성하고 인증 헤더로 추가한 후 Websocket을 연결합니다. 이 때, 앞서 정의한 이벤트 리스너 객체를 인자로 전달하여 동작 별로 실행할 로직을 적용합니다.