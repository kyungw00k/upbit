---
updatedAt: 2026-07-08T04:49:23.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Upbit Strategy Toolkit 관련 문의

### Upbit Strategy Toolkit으로 실제 매매가 되나요?

되지 않습니다. 백테스트(과거 데이터 시뮬레이션) 전용이며, 실시간 주문·계좌 변경은 지원하지 않습니다. 실거래가 필요하면 [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills)를 사용합니다.

<br />

### 설치하려면 무엇이 필요한가요?

AI 코딩 에이전트(Claude Code·Cursor·Codex 중 하나)와 uv가 필요합니다. 이후 `npx skills add upbit-official/upbit-strategy-toolkit -s '*' -a claude-code -a codex -a cursor -y --copy`로 설치합니다.

<br />

### 코딩을 전혀 몰라도 사용할 수 있나요?

네. 설치 이후에는 코딩 없이 대화(자연어)로만 사용합니다. 전략 생성·백테스트·결과 비교 모두 에이전트에게 말로 요청하면 됩니다. 다만 최초 1회 설치는 터미널에 명령을 붙여넣는 과정이 필요합니다. 시작 가이드의 순서를 그대로 따라 하면 됩니다.

<br />

### uv, npx 같은 명령이 낯섭니다. 꼭 필요한가요?

uv는 툴킷을 구동하는 파이썬 실행 도구, npx는 에이전트에 스킬을 설치하는 명령입니다. 설치는 한 번만 하면 됩니다.

<br />

### 설치했는데 에이전트가 스킬을 인식하지 못하거나 명령이 동작하지 않습니다

`/setup` 스킬로 환경을 진단할 수 있습니다. 에이전트에게 `/setup`을 요청하면 `uv` 설치 여부, 스킬 인식 상태 등을 점검하고 필요한 조치를 안내합니다. 그래도 해결되지 않으면 에이전트를 재시작하거나 설치 명령을 다시 실행합니다.

<br />

### 전략 아이디어가 없어도 사용할 수 있나요?

이 도구는 이용자가 정의한 매매 규칙을 과거 데이터로 계산·검증하는 도구이며, 수익이 나는 전략을 추천하거나 종목·매매 시점을 제시하지 않습니다. 검증할 규칙(매수·매도 조건 등)은 이용자가 먼저 정해야 합니다. 규칙을 구체화하는 과정에서 에이전트에게 아이디어를 정리해달라고 요청할 수 있으나, 에이전트가 제시하거나 정리한 내용은 투자 권유가 아니며, 그 내용이 의도와 맞는지는 백테스트 전에 반드시 직접 확인해야 합니다.

<br />

### 전략을 어떻게 요청해야 잘 만들어지나요?

"언제 사고(매수), 언제 팔지(매도)"를 중심으로, 가능하면 지표·기간·손절/익절 기준을 함께 말하면 더 정확히 정리됩니다. 전문 용어를 몰라도 "하락장에서 단기 반등을 노리는 규칙"처럼 말해도 되고, 에이전트가 이를 전략 JSON으로 구성합니다. 구성된 전략은 저장 전 내용을 함께 확인하고, 조건을 바꿔가며 반복해 다듬는 것을 권장합니다.

<br />

### 요청한 기간과 백테스트 결과의 기간이 다릅니다

거래가 없어 캔들이 존재하지 않는 구간이 있으면, 백테스트는 데이터가 실제로 존재하는 기간에 대해서만 수행됩니다. 이때 결과의 Period에는 요청 기간이 아니라 실제 분석된 기간(min\~max) 이 표시됩니다. 거래가 드문 마켓에서는 요청 기간과 실제 기간의 차이가 클 수 있으므로, 결과의 Period 표기를 함께 확인합니다.

<br />

### "캔들 데이터가 부족하다"는 오류가 납니다

해당 구간에 거래 가능한 캔들이 없으면, 오류와 함께 백테스트가 중단됩니다.<Anchor target="_blank" href="https://www.upbit.com/historical_data/main"> 업비트 Historical Market Data </Anchor>에서 캔들이 제공되는 기간을 확인하여 시작일을 데이터가 있는 구간으로 조정합니다.

<br />

### 수수료와 초기 자본은 어떻게 설정되나요?

초기 자본 기본값은 KRW: 1,000,000 / BTC: 0.01 / USDT: 1,000 이며, 수수료율은 거래소·마켓별 내장 규칙에서 자동 적용됩니다(예: 업비트 한국 KRW 0.05%, BTC 0.25%). CLI에서는 `--capital`, `--fee-rate`로 재정의할 수 있습니다.

<br />

### 백테스트 수익률이 좋으면 실제로도 수익이 나나요?

아닙니다. 백테스트 결과는 과거 데이터 기반 참고 수치이며 미래 수익을 보장하지 않습니다. 또한 슬리피지·스프레드·시장 충격은 반영되지 않아 실제 매매와 차이가 있을 수 있고, 과거에 과도하게 최적화된 전략은 실제 시장에서 작동하지 않을 수 있습니다.

<br />

### 데이터는 어디서 오고, 매번 새로 받나요?

<Anchor target="_blank" href="https://www.upbit.com/historical_data/main">업비트 Historical Market Data</Anchor>에서 캔들 데이터를 가져와 로컬에 캐시합니다. 최신 데이터로 다시 받으려면 CLI에서 `--force-refresh` 옵션을 사용합니다.

# Sibling pages

* [API 공통 문의](https://docs.upbit.com/kr/docs/faq-api.md)
* [주문 관련 문의](https://docs.upbit.com/kr/docs/faq-order.md)
* [입출금 관련 문의](https://docs.upbit.com/kr/docs/faq-withdraw-deposit.md)