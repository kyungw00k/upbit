Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내] Open API 요청 헤더 중복 전송 방지 안내

안녕하세요. 업비트 개발자센터입니다.

`Content-Type`, `User-Agent`, `Cache-Control` 등의 헤더가 중복으로 전송되지 않도록 요청 환경을 점검해 주시기 바랍니다. 직접 중복으로 지정하지 않았더라도 사용 중인 라이브러리, 프레임워크 또는 공통 헤더 설정에 의해 의도하지 않게 발생할 수 있습니다.

향후 시스템 변경 이후 동일한 헤더가 중복으로 전송되면 여러 값이 하나로 병합될 수 있습니다. 특히 Content-Type 헤더가 병합되는 경우 올바르지 않은 값으로 인식되어 요청이 정상적으로 처리되지 않을 수 있습니다.

### 반영 예정일

* 2026-07-31 (금)

### 확인 및 조치 사항

Open API 요청에 포함되는 각 헤더가 한 번만 전송되도록 요청 환경을 점검해 주세요.
특히 Content-Type 헤더는 아래와 같이 한 번만 지정하시기 바랍니다.

```text
Content-Type: application/json
```

<br />다음 항목을 함께 확인해 주세요.

* 개별 요청과 공통 요청 설정에 동일한 헤더가 각각 추가되어 있지 않은지
* HTTP 클라이언트 또는 라이브러리의 기본 헤더와 직접 지정한 헤더가 중복되지 않는지

Open API를 직접 연동한 경우 요청 생성 로직과 공통 헤더 설정을 확인하시기 바랍니다.

향후 중복된 Content-Type 헤더가 포함된 요청에는 400 오류가 반환될 수 있습니다. 안정적인 Open API 이용을 위해 반영 예정일까지 요청 환경을 점검하시기 바랍니다.

관련 문의는 <open-api@upbit.com>으로 보내주시기 바랍니다.

<br />

감사합니다.