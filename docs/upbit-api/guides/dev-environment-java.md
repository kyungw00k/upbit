---
updatedAt: 2026-05-04T05:41:42.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Java

Java 환경에서 Upbit Open API를 연동하기 위한 개발 환경 설정 방법을 안내합니다. 

## IDE(IntelliJ IDEA) 환경 설정

IntelliJ IDEA는 JetBrains에서 개발한 통합 개발 환경으로 Java와 Kotlin 개발 환경에서 널리 사용됩니다. 이 가이드에서는 IntelliJ IDEA를 사용해 Java 개발 환경을 구축합니다.

1. **IntelliJ IDEA 설치**\
   JetBrains 공식 웹사이트에서 IntelliJ IDEA를 설치할 수 있습니다. JetBrains은 Community Edition(무료), Ultimate(유료) 등 2가지 버전의 IntelliJ IDEA를 제공합니다. 사용자의 목적에 맞는 버전을 선택해 설치합니다.
   * [IntelliJ IDEA 다운로드 바로가기](http://jetbrains.com/idea/download)

2. **JDK 설치**\
   Java 개발 환경을 구축하기 위해 JDK(Java Development Kit)를 설치합니다. Amazon, IBM, Microsoft 등 다양한 기관에서 OpenJDK를 제공하고 있으며 사용 환경, 요구 사항에 따라 가장 적합한 배포판을 선택해 설치합니다.

3. **IntelliJ IDEA 설정**
   * IntelliJ IDEA를 실행하고 기존 프로젝트를 열거나 새로운 프로젝트를 생성합니다.
   * 상단 메뉴에서 File > Project Structure를 클릭합니다.
   * Project Structure 창의 좌측 메뉴에서 Platform Settings > SDKs를 선택합니다.
   * 설치한 JDK가 목록에 없다면 상단의 + 버튼을 클릭하고 \[Add JDK] 버튼을 클릭합니다.
   * JDK 설치 경로를 선택합니다. 운영체제별 일반적인 JDK 설치 경로는 다음과 같습니다.
     * macOS: `/Library/Java/JavaVirtualMachines/<jdk-version>/Contents/Home`
     * Windows: `C:\Program Files\Java\<jdk-version>`\
       올바르게 경로를 선택하면 IntelliJ IDEA가 자동으로 JDK를 인식해 추가합니다.
   * Project Structure 창 좌측 메뉴에서 Project를 선택합니다.
   * 오른쪽 화면의 Project SDK 드롭다운 메뉴에서 추가한 JDK를 선택합니다.
   * Language level을 프로젝트에 맞는 Java 버전으로 선택합니다.
   * 설정 변경 완료 후 우측 하단의 \[OK] 버튼 또는 \[Apply] 버튼을 클릭해 설정을 저장하고 창을 닫습니다.
     <br />

4. **Java 환경 설정 테스트**\
   JDK 설정을 확인하기 위해 간단한 Java 클래스 코드를 작성해 테스트합니다.

```java
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, Java!");
    }
}
```

실행 후 다음과 같은 결과가 출력되면 JDK 설정이 올바르게 완료된 것을 확인할 수 있습니다.

```
Hello, Java!
```

<br />

### HTTP 클라이언트 라이브러리 안내

REST API와 WebSocket을 Java에서 호출하기 위한 대표적인 HTTP 클라이언트 라이브러리들을 소개합니다. 개발 환경 및 Framework에 따라 적합한 라이브러리를 선택하여 사용하십시오.

#### 라이브러리별 특징 및 설치 방법

1. **OkHttp**

OkHttp는 널리 사용되는 Java/Android HTTP 클라이언트로, REST API와 WebSocket 통신 모두 지원합니다. 경량이면서 비동기 및 동기 요청 모두 지원하고, 다양한 커스텀 설정에 용이합니다. 아래와 같이 설치 의존성을 설정하거나 직접 다운로드 받을 수 있습니다.

```Text Gradle
implementation 'com.squareup.okhttp3:okhttp:{version}'
```
```Text Maven
<dependency>
  <groupId>com.squareup.okhttp3</groupId>
  <artifactId>okhttp</artifactId>
  <version>{version}</version>
</dependency>
```

REST API 호출 및 WebSocket 연결 예제는 아래와 같습니다.

```java REST API
OkHttpClient client = new OkHttpClient();

Request request = new Request.Builder()
  	.url("https://api.upbit.com/v1/ticker?markets=KRW-BTC")
		.addHeader("accept", "application/json")
    .build();

try (Response response = client.newCall(request).execute()) {
    System.out.println(response.body().string());
}
```
```java WebSocket
import okhttp3.*;

OkHttpClient client = new OkHttpClient();

Request request = new Request.Builder()
    .url("wss://api.upbit.com/websocket/v1")
    .build();

WebSocketListener listener = new WebSocketListener() {
    @Override
    public void onOpen(WebSocket webSocket, Response response) {
        String subscribeMessage = "[{\"ticket\":\"test\"},{\"type\":\"ticker\",\"codes\":[\"KRW-BTC\"]}]";
        webSocket.send(subscribeMessage);
    }
    @Override
    public void onMessage(WebSocket webSocket, String text) {
        System.out.println("Received: " + text);
    }
    @Override
    public void onClosed(WebSocket webSocket, int code, String reason) {
        System.out.println("Closed: " + reason);
    }
    @Override
    public void onFailure(WebSocket webSocket, Throwable t, Response response) {
        t.printStackTrace();
    }
};
client.newWebSocket(request, listener);
client.dispatcher().executorService().shutdown();
```

<br />

2. **Spring WebClient**

Spring WebClient는 Spring 5부터 제공되는 Reactive HTTP/WebSocket 클라이언트입니다. 비동기 및 스트림 기반 처리에 최적화되어 있어 대용량, 고성능 API 연동에 적합합니다. Spring Boot 2.0 이상을 사용하는 경우 REST, WebSocket 모두 손쉽게 사용할 수 있습니다. 아래와 같이 설치 의존성을 설정하거나 직접 다운로드 받을 수 있습니다.

```Text Gradle
implementation 'org.springframework.boot:spring-boot-starter-webflux'
```
```Text Maven
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-webflux</artifactId>
</dependency>
```

REST API 호출 및 WebSocket 연결 예제는 아래와 같습니다

```java REST API
import org.springframework.web.reactive.function.client.WebClient;

WebClient client = WebClient.create();
String response = client.get()
  	.uri("https://api.upbit.com/v1/ticker?markets=KRW-BTC")
		.header("accept", "application/json") 
    .retrieve()
    .bodyToMono(String.class)
    .block();
System.out.println(response);
```
```java WebSocket
import org.springframework.web.reactive.socket.client.ReactorNettyWebSocketClient;
import org.springframework.web.reactive.socket.WebSocketHandler;
import org.springframework.web.reactive.socket.WebSocketSession;
import reactor.core.publisher.Mono;

import java.net.URI;

ReactorNettyWebSocketClient client = new ReactorNettyWebSocketClient();
client.execute(
        URI.create("wss://api.upbit.com/websocket/v1"),
        session -> session.send(
                Mono.just(session.textMessage(
                        "[{\"ticket\":\"test\"},{\"type\":\"ticker\",\"codes\":[\"KRW-BTC\"]}]"
                ))
            ).thenMany(session.receive()
                .map(msg -> {
                    System.out.println("Received: " + msg.getPayloadAsText());
                    return msg;
                })
            ).then()
).block();
```

<br />

3. **Retrofit**

Retrofit은 OkHttp 기반의 REST API 클라이언트로, 인터페이스로 API 명세를 정의한 뒤 gson 라이브러리와 함께 사용하여 DTO 자동 매핑, 비동기 호출 등을 쉽게 구현할 수 있습니다. 응답을 객체로 역직렬화하거나, API 명세 관리가 중요한 프로젝트에 유용합니다. WebSocket은 직접 지원하지 않으나 OkHttp를 기반으로 하는 라이브러리이므로 OkHttp의 WebSocket 기능을 활용할 수 있습니다.

```Text Gradle
implementation 'com.squareup.retrofit2:retrofit:2.11.0'
implementation 'com.squareup.retrofit2:converter-gson:2.11.0'
```
```Text Maven
<dependency>
  <groupId>com.squareup.retrofit2</groupId>
  <artifactId>retrofit</artifactId>
  <version>2.11.0</version>
</dependency>
<dependency>
  <groupId>com.squareup.retrofit2</groupId>
  <artifactId>converter-gson</artifactId>
  <version>2.11.0</version>
</dependency>
```

REST API 호출 예제는 아래와 같습니다.

```java
import retrofit2.Call;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.http.GET;
import retrofit2.http.Query;

import java.util.List;
import java.util.Map;

public interface UpbitService {
  	@GET("v1/ticker")
		@Headers("accept: application/json")
    Call<List<Map<String, Object>>> getTicker(@Query("markets") String markets);
}

Retrofit retrofit = new Retrofit.Builder()
    .baseUrl("https://api.upbit.com/")
    .addConverterFactory(GsonConverterFactory.create())
    .build();

UpbitService service = retrofit.create(UpbitService.class);
Call<List<Map<String, Object>>> call = service.getTicker("KRW-BTC");
List<Map<String, Object>> response = call.execute().body();
System.out.println(response);
```

<br />

4. **Java 표준 HttpClient (Java 11+)**

Java 11 이상에서 내장되는 표준 HTTP 클라이언트로, 별도의 외부 라이브러리 설치 없이 REST API 연동이 가능합니다. HTTP/2, 비동기, 동기 모두 지원하며 간단한 REST 연동 시 간편하게 사용할 수 있습니다. WebSocket도 지원하나 사용법은 상대적으로 복잡합니다.

```java REST API
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

HttpClient client = HttpClient.newHttpClient();
HttpRequest request = HttpRequest.newBuilder()
  	.uri(URI.create("https://api.upbit.com/v1/ticker?markets=KRW-BTC"))
		.header("accept", "application/json")
    .build();
HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
System.out.println(response.body());
```
```java WebSocket
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.WebSocket;
import java.util.concurrent.CompletionStage;

HttpClient client = HttpClient.newHttpClient();
WebSocket webSocket = client.newWebSocketBuilder()
    .buildAsync(URI.create("wss://api.upbit.com/websocket/v1"), new WebSocket.Listener() {
        @Override
        public void onOpen(WebSocket webSocket) {
            String subscribeMessage = "[{\"ticket\":\"test\"},{\"type\":\"ticker\",\"codes\":[\"KRW-BTC\"]}]";
            webSocket.sendText(subscribeMessage, true);
            WebSocket.Listener.super.onOpen(webSocket);
        }
        @Override
        public CompletionStage<?> onText(WebSocket webSocket, CharSequence data, boolean last) {
            System.out.println("Received: " + data);
            return WebSocket.Listener.super.onText(webSocket, data, last);
        }
        @Override
        public void onError(WebSocket webSocket, Throwable error) {
            error.printStackTrace();
        }
    }).join();
```

# Sibling pages

* [Python](https://docs.upbit.com/kr/docs/dev-environment-python.md)
* [Node.js](https://docs.upbit.com/kr/docs/dev-environment-nodejs.md)