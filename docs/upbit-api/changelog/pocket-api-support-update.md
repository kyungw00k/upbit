Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내] SDK, CLI, Skills '포켓 API' 지원 업데이트 안내

안녕하세요. 업비트 개발자센터입니다.

업비트 SDK, CLI 및 Skills가 업데이트되었습니다.

이번 업데이트를 통해 기존에 제공 중인 포켓(Pocket) API를 Python, TypeScript, Go SDK와 CLI, Skills에서도 사용할 수 있도록 지원 범위를 확대했습니다.

## 1. 적용 일시

* 적용일: 2026-09-01 (화)

## 2. 주요 업데이트 내용

### 업데이트 버전

| 구분                                                                                                            | 버전       |
| ------------------------------------------------------------------------------------------------------------- | -------- |
| <Anchor target="_blank" href="https://github.com/upbit-official/upbit-sdk-python">Python SDK</Anchor>         | `v1.0.0` |
| <Anchor target="_blank" href="https://github.com/upbit-official/upbit-sdk-typescript">TypeScript SDK</Anchor> | `v1.0.0` |
| <Anchor target="_blank" href="https://github.com/upbit-official/upbit-sdk-go">Go SDK</Anchor>                 | `v1.0.0` |
| <Anchor target="_blank" href="https://github.com/upbit-official/upbit-cli">Upbit CLI</Anchor>                 | `v1.0.0` |
| <Anchor target="_blank" href="https://github.com/upbit-official/upbit-agent-skills">Upbit Skills</Anchor>     | `v1.0.0` |

> 포켓 API 기능을 사용하려면 최신 버전으로 업데이트해 주세요.

### 포켓 API 지원 추가

Python, TypeScript, Go SDK와 CLI, Skills에서 <Anchor target="_blank" href="https://docs.upbit.com/kr/reference/pocket-overview">포켓 API</Anchor>를 지원합니다.

이번 업데이트를 통해 다음 포켓 API를 사용할 수 있습니다.

* 포켓 정보 조회
* 포켓별 API Key 목록 조회
* 서브포켓 잔고 조회
* 메인포켓 자산 이전
* 메인포켓 자산 이전 목록 조회
* 서브포켓 자산 이전
* 서브포켓 자산 이전 목록 조회

CLI에서는 `pockets` 리소스를 통해 포켓 관련 기능을 사용할 수 있습니다.

```bash
# 서브포켓 잔고 조회
upbit pockets retrieve-balance \
  --uuid "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
```

Upbit Skills에서도 포켓 조회 및 자산 이전 등 포켓 관련 작업을 자연어로 요청할 수 있습니다.

> “서브포켓 잔고 확인해줘”

자세한 사용 방법은 아래 개발자센터 문서를 참고해 주세요.

* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk)
* [Upbit CLI](https://docs.upbit.com/kr/docs/upbit-cli)
* [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills)
* [포켓 API](https://docs.upbit.com/kr/reference/pocket-overview)

***

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

감사합니다.

업비트 개발자센터 드림