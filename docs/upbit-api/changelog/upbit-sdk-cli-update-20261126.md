---
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# [업데이트] Upbit SDK 및 CLI 업데이트 안내

안녕하세요. 업비트 개발자센터입니다.

2026년 11월 25일, 업비트 공식 SDK 및 CLI의 신규 버전이 배포됩니다.

이번 업데이트에서는 WebSocket 연결 및 이벤트 처리, 인증 및 응답 처리 기능을 개선하고 일부 오류를 수정합니다.

### 1. 업데이트 일정

* **업데이트 일자:** 2026년 11월 26일 (목)

### 2. 업데이트 대상 및 주요 변경 사항

| 프로젝트                                                                     | 업데이트 버전  | 주요 변경 사항                                                |
| ------------------------------------------------------------------------ | -------- | ------------------------------------------------------- |
| [Python SDK](https://github.com/upbit-official/upbit-sdk-python)         | `v1.2.0` | WebSocket 재연결 및 이벤트 처리 개선, 인증 및 응답 모델 개선                |
| [TypeScript SDK](https://github.com/upbit-official/upbit-sdk-typescript) | `v1.2.0` | WebSocket 기능 개선, 최소 지원 TypeScript 버전 상향, 일부 응답 필드 타입 정정 |
| [Go SDK](https://github.com/upbit-official/upbit-sdk-go)                 | `v1.0.1` | 인증 및 연결 처리 관련 오류 수정                                     |
| [Upbit CLI](https://github.com/upbit-official/upbit-cli)                 | `v1.0.1` | 일부 오류 수정                                                |

### 3. 주요 업데이트 내용

#### WebSocket 연결 및 이벤트 처리 개선

Python 및 TypeScript SDK의 WebSocket 관련 기능을 개선합니다.

* **Python SDK:** WebSocket 재연결 및 이벤트 처리 개선
* **TypeScript SDK:** WebSocket 관련 기능 개선

#### 인증 및 응답 처리 개선

* **Python SDK:** 인증 처리 및 응답 모델 개선
* **Go SDK:** 인증 및 연결 처리 관련 오류 수정

#### TypeScript SDK 호환성 변경

TypeScript SDK v1.2.0부터 최소 지원 TypeScript 버전이 상향되며, 일부 응답 필드의 반환 타입이 정정됩니다.

| 구분       | 기능                                       | 변경 대상                  | 변경 전     | 변경 후     |
| -------- | ---------------------------------------- | ---------------------- | -------- | -------- |
| 최소 지원 버전 | 개발 환경                                    | TypeScript             | `4.9`    | `5.8`    |
| 응답 타입 정정 | 입출금 서비스 상태 조회 (`GET /v1/status/wallet`)  | `block_height`         | `string` | `number` |
| 응답 타입 정정 | 출금 가능 정보 조회 (`GET /v1/withdraws/chance`) | `withdraw_limit.fixed` | `string` | `number` |

TypeScript SDK를 이용 중인 경우, 업데이트 전 프로젝트의 TypeScript 버전과 해당 응답 필드를 사용하는 코드의 호환성을 확인해 주시기 바랍니다.

### 4. 이용 안내

* 업데이트된 SDK 및 CLI는 각 프로젝트의 공식 GitHub 저장소에서 확인하실 수 있습니다.
* 기존 SDK 및 CLI를 이용 중인 경우, 프로젝트별 변경 사항과 애플리케이션의 호환성을 확인한 후 업데이트해 주시기 바랍니다.

***

자세한 내용은 [업비트 공식 GitHub](https://github.com/upbit-official)를 참고해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림