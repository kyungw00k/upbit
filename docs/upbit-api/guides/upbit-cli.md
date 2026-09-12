---
updatedAt: 2026-09-01T06:32:23.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Upbit CLI

터미널 환경에서 업비트 API를 호출하는 공식 CLI 도구의 개념과 사용 방법을 안내합니다.

## Upbit CLI란 무엇인가요?

CLI(Command Line Interface)는 텍스트 기반의 명령줄 인터페이스로, 터미널에서 명령어를 입력하여 기능을 실행하는 방식입니다.<br />Upbit CLI는 터미널 환경에서 Upbit API를 호출할 수 있도록 제공되는 공식 명령줄 도구입니다.

Upbit CLI를 사용하면 별도의 코드 작성 없이 업비트 개발자센터에서 지원하는 시세 조회, 계좌 조회, 주문, 입출금, 포켓 조회 및 자산 이전 등의 작업을 수행할 수 있습니다. 명령어 기반으로 API를 실행할 수 있어 반복적인 작업을 간단하게 처리할 수 있습니다.

<br />

## 언제 사용하나요?

다음과 같은 경우 Upbit CLI 사용을 권장합니다.

* 터미널에서 업비트 API를 빠르게 호출하고 결과를 확인하고 싶은 경우
* 반복적인 API 호출을 스크립트로 구성하려는 경우
* API 응답을 검증하거나 문제를 디버깅하려는 경우

<br />

## 사전 준비

Upbit CLI를 사용하려면 다음이 필요합니다.

* Node.js 및 npm 또는 Go 환경
* 업비트 Open API Key
  * 인증이 필요한 API를 사용하려면 Access Key와 Secret Key가 필요합니다.
  * 인증이 필요한 API 목록은 [Upbit API 이용 준비](https://docs.upbit.com/kr/docs/api-setup) 문서에서 확인 가능합니다.

<br />

## 어떻게 시작하나요?

### npm으로 설치

```Text bash
npm install -g @upbit-official/upbit-cli
```

### Go로 설치

```
go install github.com/upbit-official/upbit-cli/cmd/upbit@latest
```

설치 후 바이너리는 Go bin 디렉토리에 생성됩니다.

설치가 완료되면 upbit 실행 파일이 Go에서 사용되는 bin폴더에 생성됩니다. 해당 경로가 PATH에 포함되어 있으면 터미널에서 upbit 명령을 바로 실행할 수 있습니다.

### 설치 확인

```Text bash
upbit --version
```

Windows에서도 동일하게 설치할 수 있습니다. 단, 환경변수 설정 방식은 PowerShell 또는 CMD에 따라 다릅니다.

<br />

### 인증 정보 설정

인증이 필요한 API를 사용하려면 API Key를 환경 변수로 설정(권장)합니다.

#### macOS / Linux

```bash bash
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

#### Windows PowerShell

```powershell
$env:UPBIT_ACCESS_KEY="your-access-key"
$env:UPBIT_SECRET_KEY="your-secret-key"
```

#### Windows CMD

```Text CMD
set UPBIT_ACCESS_KEY=your-access-key
set UPBIT_SECRET_KEY=your-secret-key
```

명령 실행 시 플래그로 API Key를 직접 전달할 수도 있습니다.

```Text bash
upbit accounts list \
  --access-key "$UPBIT_ACCESS_KEY" \
  --secret-key "$UPBIT_SECRET_KEY"
```

<br />

## 정상 동작 확인은 어떻게 하나요?

먼저 공개 API를 호출해 CLI가 정상적으로 설치되었는지 확인합니다.

```Text bash
upbit trading-pairs list --is-details=false
```

정상적으로 실행되면 마켓 목록이 반환됩니다.

이후 인증이 필요한 API를 호출해 인증 정보가 올바른지 확인합니다.

```Text bash
upbit accounts list
```

정상적으로 설정된 경우 보유 자산 목록을 반환합니다.

<br />

## 어떤 작업을 할 수 있나요?

#### 시세 조회

```Text bash
upbit tickers list-by-trading-pairs --markets "KRW-BTC"
```

#### 마켓 정보 조회

```
upbit trading-pairs list
```

#### 주문 및 계좌 조회

```
upbit orders list-open --market KRW-BTC
upbit accounts list
```

#### 포켓 조회 및 관리

> 포켓 자산 조회는 메인포켓 API Key로만 이용할 수 있습니다. 서브포켓 API Key로 요청하는 경우 out\_of\_scope 오류가 발생합니다.

```bash
upbit pockets retrieve-balance 
  --uuid "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
```

#### 입출금 조회

```
upbit deposits list
upbit withdraws list
```

<br />

## 출력 형식과 데이터 처리는 어떻게 하나요?

### 출력 형식 설정

```Text bash
upbit accounts list --format json
upbit trading-pairs list --format yaml
```

지원 형식: `auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`

자동화 스크립트에서는 일반적으로 `json` 형식을 권장합니다.

### 데이터 필터링

`--transform` 옵션과 *\*GJSON* 구문을 사용하여 필요한 데이터만 추출할 수 있습니다.

*\*GJSON : JSON 데이터에서 원하는 값을 경로(path)로 지정해서 추출하는 문법.*

```Text bash
# 마켓 코드만 추출
upbit trading-pairs list --transform "#.market"

# BTC 현재가만 추출(GJSON)
upbit tickers list-by-trading-pairs \
  --markets "KRW-BTC" \
  --transform "0.trade_price"
```

<br />

## 페이지네이션은 어떻게 동작하나요?

목록 조회 시 `--max-items` 옵션으로 가져올 항목 수를 지정할 수 있습니다.

```Text bash
upbit orders list-open --max-items 20
upbit deposits list --max-items 20
```

<br />

## 에러 발생 시 어떻게 확인하나요?

`--debug` 옵션을 사용하면 HTTP 요청 및 응답 내용을 확인할 수 있습니다.

```Text bash
upbit accounts list --debug
```

<br />

## 참고

더 많은 명령과 예제는 아래 문서를 참고하세요.

* [Upbit CLI ReadMe](https://github.com/upbit-official/upbit-cli)
* [Upbit CLI 예제 코드](https://github.com/upbit-official/upbit-cli/tree/main/examples)

# Sibling pages

* [개요](https://docs.upbit.com/kr/docs/upbit-ai-agent-overview.md)
* [Upbit SDK](https://docs.upbit.com/kr/docs/upbit-sdk.md)
* [Upbit Skills](https://docs.upbit.com/kr/docs/upbit-agent-skills.md)
* [Upbit Strategy Toolkit](https://docs.upbit.com/kr/docs/upbit-strategy-toolkit.md)
* [CCXT 라이브러리 연동 안내](https://docs.upbit.com/kr/docs/ccxt-library-guide.md)
* [업비트 어시스턴트 이용 안내](https://docs.upbit.com/kr/docs/assistant-guide.md)