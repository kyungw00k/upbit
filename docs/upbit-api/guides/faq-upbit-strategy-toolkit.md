---
updatedAt: 2026-11-25T17:00:05.000Z
agentTools:
  siteIndex: https://docs.upbit.com/llms.txt
  projectIndex: https://docs.upbit.com/kr/llms.txt
---

# Upbit Strategy Toolkit 관련 문의

<Accordion title="Upbit Strategy Toolkit으로 실제 매매가 가능한가요?" icon="fa-info-circle">
  아니요. Upbit Strategy Toolkit은 과거 데이터를 활용한 백테스트(시뮬레이션) 전용 도구로, 실시간 주문이나 계좌 변경 기능은 지원하지 않습니다.

  실제 매매 기능이 필요한 경우 [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills)를 참고해 주세요.
</Accordion>

<Accordion title="설치하려면 무엇이 필요한가요?" icon="fa-info-circle">
  AI 코딩 에이전트(Claude Code, Cursor, Codex 중 하나)와 Python 실행 도구인 `uv`가 필요합니다.

  필요한 환경을 준비한 후, 터미널에서 다음 명령어를 실행하여 설치할 수 있습니다.

  ```bash
  npx skills add upbit-official/upbit-strategy-toolkit -s '*' -a claude-code -a codex -a cursor -y --copy
  ```

  자세한 설치 방법은 시작 가이드를 참고해 주세요.
</Accordion>

<Accordion title="코딩을 몰라도 사용할 수 있나요?" icon="fa-info-circle">
  네. 설치를 완료한 후에는 별도의 코드를 작성하지 않고 자연어로 에이전트에게 요청하여 사용할 수 있습니다.

  전략 생성, 백테스트 실행, 결과 비교 등을 대화로 요청할 수 있습니다. 다만, 최초 설치 시에는 터미널에서 명령어를 실행하는 과정이 필요합니다.

  자세한 설치 절차는 시작 가이드를 참고해 주세요.
</Accordion>

<Accordion title="uv와 npx는 무엇이며, 꼭 필요한가요?" icon="fa-info-circle">
  `uv`는 툴킷 실행에 필요한 Python 환경을 관리하는 도구이며, `npx`는 AI 코딩 에이전트에 스킬을 설치할 때 사용하는 명령어입니다.

  최초 설치 과정에서 필요하며, 설치가 완료된 후에는 일반적인 전략 생성 및 백테스트 작업을 자연어로 요청할 수 있습니다.
</Accordion>

<Accordion title="설치 후 에이전트가 스킬을 인식하지 못하거나 명령이 동작하지 않으면 어떻게 하나요?" icon="fa-info-circle">
  `/setup` 스킬을 이용하여 설치 환경을 점검할 수 있습니다.

  에이전트에 `/setup` 실행을 요청하면 `uv` 설치 여부, 스킬 인식 상태 등을 점검하고 필요한 조치를 안내합니다.

  문제가 지속된다면 에이전트를 재시작하거나 설치 명령어를 다시 실행해 주세요.
</Accordion>

<Accordion title="전략 아이디어가 없어도 사용할 수 있나요?" icon="fa-info-circle">
  Upbit Strategy Toolkit은 이용자가 정의한 매매 규칙을 과거 데이터로 계산하고 검증하는 도구입니다. 수익이 나는 전략을 추천하거나 특정 종목 및 매매 시점을 제시하지 않습니다.

  백테스트를 실행하려면 매수·매도 조건 등 검증할 규칙을 먼저 정해야 합니다. 규칙을 구체화하는 과정에서 에이전트에게 아이디어 정리를 요청할 수 있습니다.

  다만, 에이전트가 제시하거나 정리한 내용은 투자 권유가 아니며, 백테스트를 실행하기 전에 해당 규칙이 이용자의 의도와 일치하는지 직접 확인해야 합니다.
</Accordion>

<Accordion title="매매 전략은 어떻게 요청하면 되나요?" icon="fa-info-circle">
  매수·매도 조건을 중심으로 사용하는 지표, 적용 기간, 손절·익절 기준 등을 구체적으로 설명하면 보다 정확하게 전략을 구성할 수 있습니다.

  전문 용어를 몰라도 "하락장에서 단기 반등을 노리는 규칙"과 같이 자연어로 요청할 수 있으며, 에이전트가 이를 전략 JSON 형식으로 구성합니다.

  생성된 전략은 저장하기 전에 내용을 확인하고, 필요한 조건을 수정하면서 반복적으로 검증하는 것을 권장합니다.
</Accordion>

<Accordion title="요청한 기간과 백테스트 결과의 기간이 다른 이유는 무엇인가요?" icon="fa-info-circle">
  요청한 기간에 거래가 발생하지 않아 캔들 데이터가 없는 구간이 포함된 경우, 백테스트는 실제 데이터가 존재하는 기간을 기준으로 수행됩니다.

  따라서 결과의 `Period`에는 요청한 전체 기간이 아닌 실제 분석된 데이터의 시작일과 종료일(min\~max)이 표시됩니다.

  거래가 적은 마켓에서는 요청 기간과 실제 분석 기간의 차이가 클 수 있으므로, 결과의 `Period`를 함께 확인해 주세요.
</Accordion>

<Accordion title="캔들 데이터가 부족하다는 오류가 발생하면 어떻게 하나요?" icon="fa-info-circle">
  요청한 기간에 백테스트를 수행할 수 있는 캔들 데이터가 충분하지 않으면 오류가 발생하고 백테스트가 중단될 수 있습니다.

  [업비트 Historical Market Data](https://www.upbit.com/historical_data/main)에서 해당 마켓의 데이터 제공 기간을 확인한 후, 백테스트 시작일을 데이터가 존재하는 기간으로 조정해 주세요.
</Accordion>

<Accordion title="초기 자본과 거래 수수료는 어떻게 설정되나요?" icon="fa-info-circle">
  백테스트의 초기 자본은 마켓별 기본값이 적용되며, 거래 수수료율은 거래소 및 마켓별로 설정된 내장 규칙에 따라 자동 적용됩니다.

  **초기 자본 기본값**

  - KRW 마켓: 1,000,000 KRW
  - BTC 마켓: 0.01 BTC
  - USDT 마켓: 1,000 USDT

  **거래 수수료율 예시**

  - 업비트 한국 KRW 마켓: 0.05%
  - 업비트 한국 BTC 마켓: 0.25%

  기본 설정을 변경하려면 CLI에서 `--capital` 또는 `--fee-rate` 옵션을 사용할 수 있습니다.
</Accordion>

<Accordion title="백테스트 수익률이 좋으면 실제 매매에서도 수익이 발생하나요?" icon="fa-info-circle">
  아니요. 백테스트 결과는 과거 데이터를 기반으로 계산한 참고 수치이며, 미래의 수익을 보장하지 않습니다.

  백테스트에는 슬리피지, 스프레드, 시장 충격 등이 반영되지 않아 실제 매매 결과와 차이가 발생할 수 있습니다.

  또한 과거 데이터에 지나치게 최적화된 전략은 실제 시장에서 동일한 성과를 내지 못할 수 있으므로, 결과를 해석할 때 이러한 한계를 고려해야 합니다.
</Accordion>

<Accordion title="백테스트 데이터는 어디에서 가져오며, 매번 새로 다운로드하나요?" icon="fa-info-circle">
  백테스트에 사용하는 캔들 데이터는 [업비트 Historical Market Data](https://www.upbit.com/historical_data/main)에서 가져오며, 다운로드한 데이터는 로컬 환경에 캐시됩니다.

  따라서 동일한 데이터를 매번 새로 다운로드할 필요는 없습니다.

  최신 데이터를 다시 다운로드하려면 CLI에서 `--force-refresh` 옵션을 사용해 주세요.
</Accordion>

# 같은 수준의 페이지

* [API 공통 문의](https://docs.upbit.com/kr/docs/faq-api.md)
* [주문 관련 문의](https://docs.upbit.com/kr/docs/faq-order.md)
* [입출금 관련 문의](https://docs.upbit.com/kr/docs/faq-withdraw-deposit.md)
* [코인빌리기 관련 문의](https://docs.upbit.com/kr/docs/lending-faq.md)