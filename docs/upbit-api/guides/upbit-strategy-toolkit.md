---
updatedAt: 2026-07-08T04:48:59.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Upbit Strategy Toolkit

Upbit Strategy Toolkit은 AI 에이전트와 대화하며 매매 전략을 설계하고, 업비트 캔들 데이터로 백테스트하는 도구입니다. 자연어로 매매 아이디어를 전달하면 에이전트가 이를 전략으로 정리하고, 과거 데이터로 시뮬레이션한 결과를 제공합니다.

> 본 도구는 **백테스트(과거 데이터 시뮬레이션) 전용**이며, 실시간 주문이나 실제 자금 이동은 발생하지 않습니다. 실시간 주문·계좌 변경·자동 매매가 필요한 경우 별도의 방법(<Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-agent-skills">Upbit Skills</Anchor> 등)을 이용하길 권장드립니다.

> 본 도구는 **베타(Beta)**&#xB85C; 제공됩니다. 정식 출시 전까지 기능·사양 및 전략 JSON 스키마가 사전 예고 없이 변경될 수 있으며, 이용 중 일부 동작이 원활하지 않을 수 있는 점 양해 부탁드립니다.

## 주요 용어

* **백테스트(Backtest)**: 매매 규칙을 과거 시세에 적용해, 해당 규칙으로 거래했다면 결과가 어땠을지 계산하는 과정입니다.
* **전략(Strategy)**: 매수·매도 시점을 정의한 매매 규칙으로, JSON 형식의 파일로 저장됩니다.
* **에이전트(Agent)**: Claude Code·Cursor 등 AI 코딩 도구를 의미하며, 본 도구는 그 위에서 대화 방식으로 동작합니다.

## 1. 사전 준비

본 툴킷을 실행하려면 `uv`와 `Node.js` / `npm`이 필요합니다.

* `uv`는 Python 프로그램을 설치하고 실행하는 도구이며, 본 툴킷 구동에 사용됩니다.
* `npm`은 Node.js와 함께 설치되는 패키지 관리 도구입니다.

아래 명령어는 최초 1회만 실행하면 됩니다. 이미 설치되어 있다면 해당 단계는 생략하셔도 됩니다.

> 터미널이 처음이신가요? 터미널은 명령어를 입력하는 검은 창입니다. 아래 방법으로 열 수 있습니다.
>
> * macOS: Spotlight(⌘ + Space)에서 "터미널" 또는 "Terminal" 검색 후 실행
> * Windows: 시작 메뉴에서 "PowerShell" 검색 후 실행
>
> 이후에는 아래 명령어를 한 줄씩 복사해 붙여넣고 Enter를 누르면 됩니다.

### 1-1. uv 설치

운영체제에 맞는 명령어를 터미널에 입력해 설치합니다.

#### macOS / Linux

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

#### Windows PowerShell

```powershell
powershell -ExecutionPolicy ByPass -c "irm https://astral.sh/uv/install.ps1 | iex"
```

설치가 완료되면 터미널을 새로 열고 아래 명령어로 정상 설치 여부를 확인합니다.

```bash
uv --version
```

버전 번호가 출력되면 정상적으로 설치된 것입니다.

### 1-2. Node.js / npm 설치

`npm`은 Node.js 설치 시 함께 설치됩니다. 가장 간단한 방법은 공식 인스톨러를 내려받는 것입니다.

[nodejs.org/en/download](https://nodejs.org/en/download) 에 접속해 **LTS** 버전 인스톨러를 내려받아 실행합니다.

* macOS: `.pkg` 파일을 실행하고 안내에 따라 설치합니다.
* Windows: `.msi` 파일을 실행하고 안내에 따라 설치합니다. (설치 시 npm이 함께 포함됩니다.)

<Accordion title="설치 명령어" icon="fa-info-circle">

아래 명령어를 사용하여 운영체제에 맞게 Node.js LTS(Long Term Support) 버전을 설치합니다.

`npm`은 Node.js 설치 시 함께 설치되므로 별도로 설치하지 않아도 됩니다.

> **유의사항**
>
> macOS / Linux 설치 명령어는 외부 스크립트를 다운로드하여 실행하는 방식입니다. 실행 전 스크립트 내용을 확인하고, 회사 또는 기관 장비에서는 내부 보안 정책에 따라 진행하시기 바랍니다.

---

### macOS / Linux

macOS / Linux에서는 `nvm`을 설치한 뒤, `nvm`을 통해 Node.js LTS 버전을 설치합니다.

터미널에서 아래 명령어를 실행합니다.

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/master/install.sh | bash && \
export NVM_DIR="$HOME/.nvm" && \
[ -s "$NVM_DIR/nvm.sh" ] && \
. "$NVM_DIR/nvm.sh" && \
nvm install --lts
```

설치가 완료되면 아래 명령어로 설치 여부를 확인합니다.

```bash
node --version
npm --version
```

버전 번호가 출력되면 정상적으로 설치된 것입니다.

예시:

```bash
v22.11.0
10.9.0
```

---

### Windows PowerShell

Windows에서는 `winget`을 사용하여 Node.js LTS 버전을 설치합니다.

PowerShell에서 아래 명령어를 실행합니다.

```powershell
winget install OpenJS.NodeJS.LTS
```

설치가 완료되면 PowerShell을 새로 열고 아래 명령어로 설치 여부를 확인합니다.

```powershell
node --version
npm --version
```

버전 번호가 출력되면 정상적으로 설치된 것입니다.

예시:

```powershell
v22.11.0
10.9.0
```

---

### macOS / Linux 삭제

macOS / Linux 설치 명령어는 `nvm`을 통해 Node.js를 설치합니다.  
따라서 삭제 시에는 `nvm` 설치 디렉터리와 셸 설정 파일에 추가된 nvm 관련 설정을 함께 정리해야 합니다.

먼저 nvm 설치 디렉터리를 삭제합니다.

```bash
nvm_dir="${NVM_DIR:-$HOME/.nvm}"

nvm unload 2>/dev/null || true

rm -rf "$nvm_dir"
```

그다음 `.bashrc`, `.zshrc`, `.profile` 등에 추가된 nvm 관련 줄을 삭제합니다.

```bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"
```

삭제가 완료되면 터미널을 새로 열고 아래 명령어로 삭제 여부를 확인합니다.

```bash
node --version
npm --version
```

정상적으로 삭제된 경우 `node` 또는 `npm` 명령어를 찾을 수 없다는 메시지가 출력됩니다.

자세한 내용은 [nvm 공식 문서](https://github.com/nvm-sh/nvm#uninstalling--removal)를 참고하세요.

---

### Windows PowerShell 삭제

Windows에서 `winget`으로 설치한 Node.js LTS를 삭제하려면 PowerShell에서 아래 명령어를 실행합니다.

```powershell
winget uninstall OpenJS.NodeJS.LTS
```

삭제가 완료되면 PowerShell을 새로 열고 아래 명령어로 삭제 여부를 확인합니다.

```powershell
node --version
npm --version
```

정상적으로 삭제된 경우 `node` 또는 `npm` 명령어를 찾을 수 없다는 메시지가 출력됩니다.

예시:

```powershell
node : The term 'node' is not recognized as the name of a cmdlet, function, script file, or operable program.
```

또는 아래와 유사한 메시지가 출력될 수 있습니다.

```powershell
'node' is not recognized as an internal or external command,
operable program or batch file.
```

</Accordion>

설치가 완료되면 터미널을 새로 열고 아래 명령어로 정상 설치 여부를 확인합니다.

```bash
node --version
npm --version
```

이미 `uv`, `node`, `npm`이 설치되어 있는 경우에는 별도로 재설치하지 않아도 됩니다. 터미널에서 아래 명령어를 실행했을 때 버전 번호가 출력되면 바로 사용 가능합니다.

```bash
uv --version
node --version
npm --version
```

<br />

## 2. 설치 및 실행

에이전트(Claude Code / Cursor / Codex)에 아래 명령으로 설치한 뒤, 대화하듯 전략을 요청합니다.&#x20;

```bash
npx skills add upbit-official/upbit-strategy-toolkit -s '*' -a claude-code -a codex -a cursor -y --copy
```

설치가 완료되면 에이전트에 다음과 같이 요청합니다.

```text
사용자:   "RSI 과매도에 진입하고 볼린저 중간선에서 청산하는 전략, 2024년으로 백테스트해줘"
에이전트: 전략 정리 → 백테스트 실행 → 결과 출력
```

아래와 같은 결과가 출력됩니다. 각 숫자의 의미는 [4. 결과 지표 안내](#4-결과-지표-안내)에서 설명합니다.

```text
Period       2024-01-01 ~ 2024-12-31
Total Return +18.32%     ← 전략의 수익률
Benchmark    +62.45%     ← 단순 보유(buy & hold) 시 수익률
MDD          -7.21%      ← 최대 낙폭 (절대값이 작을수록 안정적)
Sharpe        1.34       ← 위험 대비 수익 (통상 1 이상이면 양호)
Win Rate      60%
```

### 두 가지 설치 모드

설치 방식은 두 가지입니다. 대부분의 경우 위의 `npx skills add`(Agent Skills 모드)로 충분하며, 도구의 내부 코드를 직접 확인하거나 수정하려는 경우에만 Clone 모드를 사용합니다. 두 모드는 전략·결과 파일이 저장되는 위치만 다릅니다. 설정 가이드는 업비트 공식 Github의 <Anchor target="_blank" href="https://github.com/upbit-official/upbit-strategy-toolkit">Upbit Strategy Toolkit 리포지토리(Repository)</Anchor>에서 확인 가능합니다.

#### Agent Skills 모드 (권장)

리포지토리 클론 없이 현재 프로젝트에 에이전트 스킬만 설치합니다. 출력 파일은 홈 폴더 아래 `$HOME/.upbit-strategy-toolkit`에 저장되므로, 작업 중인 프로젝트 폴더에 영향을 주지 않습니다.

```bash
npx skills add upbit-official/upbit-strategy-toolkit -s '*' -a claude-code -a codex -a cursor -y --copy
```

```text
$HOME/.upbit-strategy-toolkit/
├── strategies/   # 저장된 전략 JSON
├── reports/      # 백테스트 결과 CSV
└── cache/
    └── upbit/    # 캔들 데이터 캐시
```

#### Clone 모드

리포지토리를 내려받아(클론) 그 폴더 안에서 실행합니다. 출력 파일도 그 폴더 안에 저장됩니다.

```bash
git clone https://github.com/upbit-official/upbit-strategy-toolkit.git
cd upbit-strategy-toolkit
uv sync
uv run upbit-strategy-toolkit --help
```

```text
upbit-strategy-toolkit/
├── strategies/   # 저장된 전략 JSON
├── reports/      # 백테스트 결과 CSV
└── cache/
    └── upbit/    # 캔들 데이터 캐시
```

어떤 모드든 결과물의 종류는 같습니다. 만든 전략은 `strategies/`, 백테스트 리포트는 `reports/`에 쌓입니다.

## 3. 기본 사용 방법

사용 흐름은 **① 전략 생성 → ② 백테스트 → ③ 결과 검토 및 수정**이며, 세 단계를 반복하며 전략을 다듬는 것이 핵심입니다. 모든 단계는 대화로 진행됩니다.

**① 전략 생성** — 매매 아이디어를 전달하면, 에이전트가 매수·매도·손절(stop-loss)·익절(take-profit) 조건을 정리해 전략 파일로 생성합니다. 전문 용어를 사용하지 않아도 다음과 같이 요청할 수 있습니다.

```text
"하락장에서도 단기 반등을 잡는 전략을 만들어줘"
"RSI 과매도 진입, 볼린저 밴드 중간선 청산 전략"
"SMA 5/20 골든크로스, 손절 3% 익절 8% 전략"
```

**② 백테스트 실행** — 전략을 적용할 과거 기간을 지정하면, 해당 기간의 실제 시세로 시뮬레이션합니다.

```text
"방금 만든 전략을 2024년 한 해 백테스트해줘"
```

**③ 결과 검토 및 수정** — 결과에 따라 조건을 바꿔 다시 실행합니다. 값 하나만 변경해 비교하는 방식으로 전략을 빠르게 검증할 수 있습니다.

```text
"진입 조건을 RSI 25 이하로 낮추고 다시 백테스트해줘"
"손절을 3%로 줄이면 결과가 어떻게 달라지는지 비교해줘"
```

이렇게 만든 전략은 `strategies/` 폴더에, 백테스트 결과 CSV는 `reports/` 폴더에 자동으로 저장됩니다.

## 4. 결과 지표 안내

백테스트를 실행하면 아래와 같은 요약이 출력됩니다. 항목이 많지만, 처음에는 굵게 표시한 세 가지(**Total Return · MDD · Benchmark**)를 중심으로 보면 됩니다.

```text
Period       2024-01-01 ~ 2024-12-31 (trading bars: 366, warmup bars: 30)
Benchmark    +62.45%
Total Return +18.32%
CAGR         +18.11%
MDD          -7.21%
Sharpe       1.34  (portfolio / full equity curve)
Sharpe       1.11  (trades / position holding periods only)
Trades       23  Win Rate 60% (before fees)
Profit Factor  2.14 (before fees)
SL 4 / TP 9 / Sell 8 / Final bar 2
Total Fees   12,345
```

각 지표의 의미는 다음과 같습니다.

| 지표                 | 의미                                  | 참고                          |
| ------------------ | ----------------------------------- | --------------------------- |
| **Total Return**   | 전략의 총 수익률 = (최종자산 / 초기자산 − 1)       | 전략의 최종 손익                   |
| **Benchmark**      | 같은 기간 단순 보유(buy & hold) 시 수익률       | 전략이 이 값을 넘어야 유의미함           |
| **MDD**            | 최대 낙폭(고점 대비 최대 하락 폭)                | 절대값이 작을수록 안정적, 클수록 중간 손실이 큼 |
| CAGR               | 연평균 복리 수익률                          | 기간이 다른 전략을 비교할 때            |
| Sharpe (portfolio) | 전체 자산 흐름 기준, 위험 대비 수익 (현금 보유 구간 포함) | 통상 1 이상 양호, 2 이상 우수         |
| Sharpe (trades)    | 개별 거래 기준, 위험 대비 수익 (포지션 보유 구간만 계산)  | 거래 단위 효율 비교용                |
| Win Rate           | 수익을 낸 거래의 비율 (수수료 전)                | 승률이 높아도 큰 손실 한 번으로 손실 가능    |
| Profit Factor      | 총이익 ÷ 총손실 (수수료 전)                   | 1보다 크면 이익이 손실보다 큼           |

> **결과 읽는 법**
>
> 위 예시는 **+18.32%** 수익을 냈으나 같은 기간 단순 보유(**+62.45%**)보다 낮은 결과입니다.<br />즉 수익이 발생했다는 것이 곧 좋은 전략을 의미하지는 않습니다.
>
> 항상 Benchmark와 비교하고, MDD(최대 손실 폭)를 함께 확인합니다. Sharpe가 두 줄로 출력되는 것은 계산 기준이 다르기 때문입니다. `portfolio`는 전체 자산 흐름, `trades`는 개별 거래 기준입니다.

체결·수수료가 결과에 어떻게 반영되는지 등 백테스트의 정확한 계산 방식은 레퍼런스의 "백테스트 해석"에서 확인할 수 있습니다.

### 결과 분석이 어렵다면

지표가 많아 결과를 판단하기 어렵다면, 여러 전략을 백테스트한 뒤 그 결과를 모아 AI 에이전트에게 비교·분석을 맡길 수 있습니다.

1. 조건을 바꿔가며 여러 전략을 백테스트합니다. 결과 CSV는 `reports/` 폴더에 순차적으로 저장됩니다.
2. 저장된 리포트를 하나의 전략집(예: `strategy-report.md`)으로 정리합니다. 파일마다 전략 설정과 성과 요약(Total Return·Benchmark·MDD·Win Rate 등)을 함께 담습니다.
3. 정리한 전략집을 에이전트에 전달하고 분석을 요청합니다.

예시 프롬프트:

```text
"reports/ 폴더의 백테스트 결과들을 비교해서 전략집 strategy-report.md로 정리해줘.
각 전략의 설정, Total Return, Benchmark, MDD, Win Rate를 표로 만들고,
벤치마크 대비 초과 성과 기준으로 정렬해줘."

"strategy-report.md를 읽고, 어떤 전략이 하락장에서 손실을 잘 방어했는지,
어떤 전략이 과최적화로 의심되는지 분석해줘. 지표만 보고 단정하지 말고
거래 횟수와 MDD도 함께 근거로 들어줘."
```

> 이때 에이전트의 분석은 참고 의견이며, 투자 판단의 근거가 되지 않습니다. AI는 결과를 잘못 해석할 수 있으므로 원본 리포트 수치를 함께 확인하세요.

***

### 유의사항

* **백테스트 결과는 투자 조언이 아닙니다.** 수익률·승률·Profit Factor·Sharpe는 과거 데이터 기반 참고 수치이며, 미래 수익을 보장하지 않습니다.
* **전략 JSON만 지원합니다.** 임의 Python 코드 실행은 지원하지 않습니다. 새로운 전략 표현이 필요한 경우 Pydantic 모델과 명시적 스키마로 정의합니다.
* **실전 매매는 지원하지 않습니다.** 실시간 주문·계좌 변경·자동 매매가 필요한 경우 `upbit-agent-skills`를 사용합니다.
* 모든 계산은 IEEE 754 배정밀도(float64)를 사용하며, 실거래 주문 실행이나 정산 목적이 아닌 백테스트 참고용 근삿값입니다.

<Accordion title="책임 한정 및 이용 주의사항" icon="fa-info-circle">
  - Upbit Strategy Toolkit(이하 '본 툴킷')은 이용자가 사용하는 제3자 제공 AI 서비스(이하 '이용자 AI')를 활용하여 기술지표를 기반으로 매매 전략을 설계하고 백테스트를 할 수 있도록 지원하는 도구입니다.

- 본 툴킷은 이용자 AI를 활용한 투자 전략 수립 및 백테스트 수행 시 편의를 제공하기 위한 보조 도구에 불과합니다. 본 툴킷은 법적·재정적·투자적 자문 또는 조언을 목적으로 하지 않으며, 회사는 어떠한 경우에도 특정 종목, 거래 시점, 투자 방식, 투자 전략 등을 추천·권유하지 않습니다.

- 이용자는 본 툴킷을 참고 목적으로만 활용하여야 하며, 본 툴킷을 활용하여 이용자 AI가 생성·제공하는 결과물, 답변, 전략, 백테스트 결과, 전략 개선 제안 또는 그 밖의 산출물(이하 총칭하여 'AI 산출물')을 자신의 책임과 판단하에 독립적으로 검토·검증하여야 합니다. 이용자는 AI 산출물이 실제 시장 상황, 거래 환경, 수수료, 슬리피지, 유동성, 체결 가능성, 가격 변동성 등 다양한 요인을 충분히 반영하지 못할 수 있음을 충분히 인지하고 본 툴킷을 이용하여야 하며, 회사는 과거 데이터를 기초로 한 백테스트 결과가 장래 또는 실제 거래에서 동일하거나 유사한 투자 성과로 이어질 것임을 보장하지 않습니다. 이용자는 AI 산출물을 투자 판단의 유일한 또는 주요한 근거로 삼지 말아야 하며, 이용자 AI를 활용한 거래·투자 행위와 관련된 모든 책임은 이용자에게 귀속됩니다.

- 회사는 본 툴킷을 활용한 AI 산출물의 정확성, 진실성, 완전성, 무결성, 최신성 등을 보증하지 않으며, AI 산출물에는 회사의 어떠한 의사나 입장도 반영되지 않습니다.

- 본 툴킷은 실전 매매(Live Trading)와 직접 연동되지 않으며, 이용자 AI를 통한 실전 매매를 희망하는 경우 별도의 방법(Upbit Agent Skills 등)을 이용하여야 합니다. 이러한 회사의 입장에도 불구하고, 이용자가 본 툴킷을 임의로 수정·변경·우회하거나 제3자가 제공하는 프로그램, 플러그인, 서비스 등과 결합하여 실전 매매, 자동 매매 또는 이에 준하는 행위를 하는 경우, 그로 인하여 발생하는 모든 결과와 책임은 이용자에게 귀속됩니다.

- 회사는 이용자 AI 등 제3자가 제공하는 서비스의 가용성, 보안성, 적법성, 정확성 및 성능 등에 대하여 책임을 부담하지 않습니다. 이용자는 제3자 서비스 이용에 따른 약관, 정책 및 위험을 스스로 확인하여야 합니다. 회사는 본 툴킷을 검증되지 않은 제3자 서비스와 연동 또는 조합하여 사용하는 것을 권장하지 않으며, 이용자의 제3자 서비스 사용으로 발생한 보안 사고(API Key 노출 등을 포함)나 금전적 손해 등에 대하여 책임을 부담하지 않습니다.

- 이용자는 본 툴킷을 활용함에 있어 관련 이용약관(업비트 이용약관, Open API 이용약관 등)을 준수하여야 합니다.

- 회사는 이용자 간 또는 이용자와 제3자 간에 본 툴킷과 관련하여 발생한 분쟁에 관여할 법적 의무가 없으며, 이와 관련하여 어떠한 책임도 지지 않습니다.

- 본 툴킷은 베타 버전으로 제공됩니다. 회사의 사정에 의하여 별도의 사전 고지 없이 제공 중단되거나 수정·변경될 수 있습니다.

- 본 툴킷 및 이에 포함된 문서, 코드, 예시, 가이드, 명령어 등 관련 리소스에 관한 권리 일체는 회사 또는 회사가 지정한 정당한 권리자에게 귀속됩니다. 본 툴킷 및 관련 리소스의 복제·배포·변경 및 영리 목적 이용은 관련 저작권 및 약관의 적용을 받으며, 본 툴킷을 활용한 프로그램 및 관련 리소스 등의 상업적 배포·판매·이용은 허용되지 않습니다.
</Accordion>

***

### 다음으로

* 전략을 직접 편집하거나 CLI를 쓰고 싶다면 → <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/upbit-strategy-toolkit-reference">**Upbit Strategy Toolkit — 레퍼런스**</Anchor>

- 설치 문제 해결이나 스킬·도구 제거 방법이 필요하다면 → <Anchor target="_blank" href="https://docs.upbit.com/kr/docs/faq-upbit-strategy-toolkit">**Upbit Strategy Toolkit — FAQ**</Anchor>

<br />

# Sub pages

* [기술 레퍼런스](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit-reference.md)
* [지표 정의](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit-indicators.md)

# Sibling pages

* [개요](https://docs.upbit.com/kr/docs/upbit-ai-agent-overview.md)
* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk.md)
* [Upbit CLI](https://docs.upbit.com/kr/docs/upbit-cli.md)
* [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills.md)
* [CCXT 라이브러리 연동 안내](https://docs.upbit.com/kr/docs/ccxt-library-guide.md)
* [업비트 어시스턴트 이용 안내](https://docs.upbit.com/kr/docs/assistant-guide.md)