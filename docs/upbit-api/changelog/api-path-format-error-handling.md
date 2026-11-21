---
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# [중요] API Path 형식 및 오류 처리 방식 점검 안내

안녕하세요. 업비트 개발자센터입니다.

업비트 API 요청 시 개발자 문서에 안내된 Path 형식을 준수해 주시기 바랍니다. **지정된 형식과 다른 Path로 요청하는 경우 오류가 반환될 예정**입니다.&#x20;

또한 현재 `error.message`를 기준으로 API 오류를 처리하고 계신 경우, `error.name`**을 기준으로 처리하도록 변경**해 주시기 바랍니다.

API 이용에 영향이 발생하지 않도록 적용일 이전에 요청 Path 및 오류 응답 처리 방식을 점검하고 필요한 경우 수정해 주시길 바랍니다.

## 1.적용 일시

* 적용 예정일: 2026-11-16 (월)

## 2.주요 내용

**1) API 요청 Path 형식 준수**

API 요청 Path는 다음 형식을 준수해야 합니다.

* Path 시작 및 중간: 경로 구분자로 / 1개만 허용
* Path 끝: / 사용 불가

**정상 요청 Path**

`/v1/orders`&#x20;

**오류 반환 대상 Path**

```text 잘못된 Path
//v1/orders
/v1//orders
/v1/orders//
```

적용일 이후 위와 같이 불필요한 슬래시(`/`)가 포함된 Path로 요청하는 경우 오류가 반환됩니다.

**2) API 오류 응답 처리 기준**

`error.message`는 안내 문구가 변경될 수 있으므로, 현재 `error.message`를 기준으로 API 오류를 식별하거나 처리하고 계신 경우 `error.name`을 기준으로 오류를 식별하고 처리해 주시기 바랍니다.

```json 오류 응답 예시
{
  "error": {
    "name": "ERROR_CODE", 
    "message": "ERROR_MESSAGE" → 비권장
  }
}
```

## 3.적용 전 확인 사항

현재 사용 중인 API 연동 환경에서 다음 사항을 확인해 주세요.

* Path 시작과 중간에 슬래시(`/`)가 중복되거나 마지막에 포함되어 있지 않은지
* BASE\_URL과 PATH 결합 시 슬래시(`/`)가 중복되지 않는지
* API 오류 응답을 `error.name` 기준으로 처리하고 있는지

예를 들어 다음과 같이 URL을 구성하면 슬래시가 중복될 수 있습니다.

| 구분        | 슬래시 중복 (수정 필요)                     | 정상 구성                             |
| --------- | ---------------------------------- | --------------------------------- |
| BASE\_URL | `https://api.upbit.com/`           | `https://api.upbit.com`           |
| PATH      | `/v1/orders`                       | `/v1/orders`                      |
| 요청 URL    | `https://api.upbit.com//v1/orders` | `https://api.upbit.com/v1/orders` |

## 4.유의사항

* 적용일 이후 Path 형식 미준수로 인해 API 요청이 실패할 수 있으므로 사전에 점검 및 테스트를 진행해 주세요.
* 적용 일정은 운영 상황에 따라 변경될 수 있으며, 변경 시 본 공지를 통해 다시 안내드리겠습니다.

원활한 API 이용을 위해 적용일 이전에 요청 Path와 오류 응답 처리 로직을 확인해 주시기 바랍니다.

***

문의 사항은 <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.
업비트 개발자센터 드림