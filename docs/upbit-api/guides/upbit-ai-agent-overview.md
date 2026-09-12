---
updatedAt: 2026-08-31T13:21:46.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 개요

Upbit AI Agent는 AI 기반 개발 환경 (Claude Code, Cursor, Codex 등)에서 Upbit API를 활용해 개발과 자동화 작업을 수행하도록 지원하는 문서 모음입니다.  시세 조회, 시장 분석, 주문 실행, 서비스 연동, 반복 작업 자동화, 전략 백테스트까지 하나의 흐름으로 다룰 수 있으며, 사용자는 목적에 따라 필요한 도구와 문서를 선택하면 됩니다.

## 전체 구조 한눈에 보기

모든 도구는 **Upbit API**를 기반으로 합니다. API에 접근하는 방식(코드 / 명령어 / 자연어)에 따라 사용하는 도구가 달라집니다.

```
                        Upbit API  (시세 · 주문 · 자산)
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
   코드로 연동            명령어로 실행           자연어로 실행
        │                    │              (AI 에이전트)
     Upbit SDK           Upbit CLI ───────────┐  │
   (Python·TS·Go)       (터미널 명령어)        │  │
                                        Upbit Agent Skills
                                        Upbit Strategy Toolkit
```

* **Upbit CLI**는 명령어로 직접 실행하는 도구이며, **Upbit Agent Skills**와 **Upbit Strategy Toolkit**은 그 위에서 AI 에이전트가 자연어 요청을 실제 작업으로 변환하도록 돕습니다.
* **Upbit SDK**는 애플리케이션 코드에 직접 연동하는 별도 경로입니다.

<br />

## 구성 요소

| 구성 요소                      | 역할                                                | 접근 방식 |
| :------------------------- | :------------------------------------------------ | :---- |
| **Upbit API**              | 시세 조회, 주문, 자산 관리 기능을 제공하는 기반 인터페이스                | —     |
| **Upbit SDK**              | Python, TypeScript, Go 환경에서 API를 연동하는 라이브러리       | 코드    |
| **Upbit CLI**              | 터미널에서 API 기능을 명령어로 실행·자동화하는 도구                    | 명령어   |
| **Upbit Agent Skills**     | AI 에이전트가 CLI를 통해 시세·주문·입출금을 자연어로 처리하도록 돕는 패키지     | 자연어   |
| **Upbit Strategy Toolkit** | AI 에이전트와 대화하며 매매 전략을 설계하고 과거 데이터로 백테스트하는 도구 (베타)  | 자연어   |
| **AI 개발 환경**               | Claude Code, Cursor, Codex 등에서 코드 작성과 실행을 지원하는 환경 | —     |

<br />

## 시작하기 전에: 공개 vs 인증

Upbit API는 두 종류의 엔드포인트로 나뉩니다. 어떤 작업을 하느냐에 따라 API 키 필요 여부가 달라집니다.

| 구분              | 하는 일                       | API 키 |
| :-------------- | :------------------------- | :---- |
| **공개(Public)**  | 시세, 호가, 캔들, 마켓 목록 조회       | 불필요   |
| **인증(Private)** | 포켓, 잔고 조회, 주문, 입출금 등 계좌 작업 | 필요    |

자세한 준비 과정은 <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/api-setup">Upbit API 이용 준비</Anchor>를 참고합니다.

<br />

## 목적별 시작 가이드

무엇을 하려는지에 따라 아래 문서부터 시작하세요.

| 목적                        | 시작 문서                                                                                                                | 적합한 작업                     |
| :------------------------ | :------------------------------------------------------------------------------------------------------------------- | :------------------------- |
| **처음이라 일단 한 번 호출해 보고 싶다** | <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-ai-agent-quick-start">빠른 시작</Anchor>              | 최소 설정으로 첫 API 호출           |
| 터미널에서 명령어로 작업             | <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-cli">Upbit CLI</Anchor>                           | 반복 실행, 자동화, 셸 스크립트, 응답 디버깅 |
| 애플리케이션·서비스에 코드로 연동        | <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-sdk">Upbit SDK</Anchor>                           | 백엔드 개발, 서비스 구현, 데이터 처리     |
| AI 에이전트로 시세·주문을 자연어 처리    | <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-agent-skills">Upbit Agent Skills</Anchor>         | 자연어 기반 조회·주문·입출금, 자동 매매 예제 |
| 매매 전략을 과거 데이터로 검증         | <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-strategy-toolkit">Upbit Strategy Toolkit</Anchor> | 자연어 전략 설계, 백테스트, 결과 비교·분석  |
| 여러 거래소를 함께 사용             | [CCXT 연동 가이드](https://docs.upbit.com/kr/docs/ccxt-library-guide)                                                     | 멀티 거래소 연동, 데이터 비교          |
| 저장소·예제 확인                 | <Anchor target="_blank" href="https://github.com/upbit-official">GitHub</Anchor>                                     | 코드 확인, 설치, 예제 참고           |

<br />

## 실전 매매와 백테스트의 차이

전략을 다루는 두 도구는 목적이 다릅니다. 혼동하지 않도록 구분해 사용하세요.

|           | Upbit Agent Skills | Upbit Strategy Toolkit |
| :-------- | :----------------- | :--------------------- |
| **목적**    | 실제 주문·입출금 실행       | 과거 데이터 기반 전략 검증        |
| **자산 이동** | 발생 (실전 매매)         | 없음 (시뮬레이션)             |

> **Upbit Strategy Toolkit은 백테스트 전용 도구입니다.** 실시간 주문·계좌 변경·자동 매매는 지원하지 않으며, 백테스트 결과는 과거 데이터 기반 참고 수치로 미래 수익을 보장하지 않습니다. 실시간 주문·계좌 변경·자동 매매가 필요한 경우 별도의 방법(Upbit Skills 등)을 이용하길 권장드립니다.

<br />

## AI 도구에서 참조하기

Upbit 개발자 문서는 AI 도구가 구조를 파악하고 참조할 수 있도록 `llms.txt` 형식을 제공합니다.

* **전체 문서 구조**: <Anchor target="_blank" href="https://docs.upbit.com/kr/llms.txt">`https://docs.upbit.com/llms.txt`</Anchor>
* **개별 문서**: URL 뒤에 `.md`를 붙이면 마크다운 형식으로 확인할 수 있습니다.

```
https://docs.upbit.com/kr/docs/getting-started
→ https://docs.upbit.com/kr/docs/getting-started.md
```

# Sibling pages

* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk.md)
* [Upbit CLI](https://docs.upbit.com/kr/docs/upbit-cli.md)
* [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills.md)
* [Upbit Strategy Toolkit](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit.md)
* [CCXT 라이브러리 연동 안내](https://docs.upbit.com/kr/docs/ccxt-library-guide.md)
* [업비트 어시스턴트 이용 안내](https://docs.upbit.com/kr/docs/assistant-guide.md)