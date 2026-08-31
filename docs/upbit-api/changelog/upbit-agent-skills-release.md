Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내] Upbit Skills 출시

안녕하세요. 업비트 개발자센터입니다.

업비트 Open API를 AI 에이전트 환경에서 보다 일관되고 편리하게 활용하실 수 있도록 **Upbit Skills**를 새롭게 출시했습니다.

Upbit Skills는 AI 에이전트가 `upbit` CLI 사용 방법과 업비트 Open API의 주요 규칙을 이해하고, 시세 조회, 잔고 확인, 주문, 입출금 등 다양한 작업을 자연어 요청 기반으로 보조할 수 있도록 지원합니다.

<br />

## 1. 이용 안내

GitHub와 개발자센터를 통해 설치 방법과 상세 가이드를 확인하실 수 있습니다.

* 출시 일자: 2026년 5월 7일
* Upbit 공식 GitHub: <https://github.com/upbit-official>
  * [Upbit Agent Skills](https://github.com/upbit-official/upbit-agent-skills)
  * [Upbit CLI](https://github.com/upbit-official/upbit-cli)
* [Upbit Agent Skills가이드](https://docs.upbit.com/kr/docs/upbit-agent-skills)

<br />

### 설치 방법

Upbit Skills 아래 명령으로 설치할 수 있습니다.

```bash bash
npx skills add upbit-official/upbit-agent-skills
```

전역 설치가 필요한 경우 `-g` 옵션을 사용할 수 있습니다.

```bash bash
npx skills add upbit-official/upbit-agent-skills -g
```

Upbit Agent Skills는 `upbit` CLI와 함께 동작합니다.

```bash bash
npm install -g @upbit-official/upbit-cli
```

설치 후 아래 명령으로 정상 설치 여부를 확인할 수 있습니다.

```bash bash
upbit --version
```

<br />

## 2. 자연어 사용 예시

설치 후 Skills를 지원하는 AI 에이전트에서 자연어로 요청할 수 있습니다.

예시:

* KRW-BTC 현재가 조회해줘.
* 내 업비트 잔고를 확인해줘.
* KRW-BTC 주문 가능 정보를 확인해줘.
* BTC 1만 원 시장가 매수 테스트 명령을 만들어줘.
* USDT 출금 가능한 네트워크를 확인해줘.
* 계정주 확인이 필요한 입금 내역을 확인해줘.

AI 에이전트는 요청 내용을 바탕으로 적절한 `upbit` CLI 명령을 제안하거나 실행을 보조합니다.

<br />

## 3. 주요 기능

| 구분      | 설명                        |
| ------- | ------------------------- |
| 시세 조회   | 현재가, 호가, 체결, 캔들, 마켓 목록 조회 |
| 계정 조회   | 잔고, 보유 자산 조회              |
| 주문 보조   | 주문 생성, 조회, 취소 및 주문 테스트    |
| 입출금 보조  | 입금 주소 조회, 출금 정보 확인        |
| 트래블룰 검증 | 계정주 확인이 필요한 입금건에 대한 검증    |
| 자동화 지원  | CLI 기반 명령 구성 및 실행 보조      |

<br />

## 4. 활용 이점

* 자연어 기반으로 업비트 Open API 작업 수행 가능
* CLI 명령 구성을 자동화하여 개발 및 테스트 효율 향상
* 업비트 마켓 표기 및 주문 규칙을 일관되게 적용
* 실제 자산에 영향을 줄 수 있는 작업 전 확인 절차를 통해 안전한 실행 지원

<br />

## 5. 유의사항

* 계정 조회, 주문, 입출금 등 일부 기능은 API 키 발급 및 권한 설정이 필요합니다.
* API Key 및 Secret Key는 환경 변수 등을 활용하여 외부에 노출되지 않도록 안전하게 관리해 주시기 바랍니다.
* 주문, 출금 등 실제 자산에 영향을 줄 수 있는 작업은 실행 전 반드시 내용을 확인해 주시기 바랍니다.
* 실제 주문 전에는 테스트 명령 또는 Dry run 예제를 먼저 사용하는 것을 권장드립니다.
* 과도한 요청 또는 비정상적인 트래픽 발생 시 서비스 이용이 제한될 수 있으므로, 업비트 Open API의 Rate Limit 정책을 준수해 주시기 바랍니다.
* Upbit Agent Skills및 관련 도구는 지속적으로 업데이트될 예정입니다.

<br />

문의 사항이 있으실 경우, <open-api@upbit.com>으로 연락해 주시기 바랍니다.

항상 더 나은 서비스를 제공하기 위해 노력하겠습니다.

감사합니다.

업비트 개발자센터 드림