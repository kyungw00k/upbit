---
updatedAt: 2026-05-07T12:48:19.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Go

Go 환경에서 Upbit SDK를 사용하여 Upbit API를 호출하기 위한 개발 환경 설정 방법을 안내합니다.

## 시작하기

Upbit Go SDK는 Go 애플리케이션에서 Upbit REST API를 사용할 수 있도록 제공되는 공식 라이브러리입니다.

요청 파라미터와 응답 필드에 대한 타입 정의를 제공하며, API 호출에 필요한 인증 설정, 요청 구성, 응답 처리 등을 Go 코드에서 일관된 방식으로 사용할 수 있습니다.

<br />

## SDK 공식 문서

Upbit Go SDK 관련 문서는 아래에서 확인할 수 있습니다.

* [시작 가이드](https://github.com/upbit-official/upbit-sdk-go/blob/main/README_KR.md) — Upbit Go SDK의 설치 및 기본 사용 방법
* [Upbit SDK API Reference](https://github.com/upbit-official/upbit-sdk-go/blob/main/api.md) — SDK가 지원하는 전체 API 목록
* [SDK 예제 코드](https://github.com/upbit-official/upbit-sdk-go/tree/main/examples) — 예제 코드 및 상세 사용법

<br />

## 사전 준비

Upbit Go SDK를 사용하려면 다음이 필요합니다.

* Go 1.25 이상
* Upbit Open API Key
  * 인증이 필요한 API를 호출하려면 Access Key와 Secret Key가 필요합니다.

<br />

## Go 설치 확인

터미널에서 다음 명령어를 실행하여 Go가 설치되어 있는지 확인합니다.

```bash
go version
```

Go가 설치되어 있지 않거나 버전이 낮은 경우 Go 공식 문서를 참고하여 설치합니다.

```text
https://go.dev/doc/install
```

<br />

## 프로젝트 생성

프로젝트 디렉터리를 생성하고 Go 모듈을 초기화합니다.

```bash
mkdir <project-name>
cd <project-name>
go mod init <module-path>
```

예시:

```bash
mkdir upbit-go-example
cd upbit-go-example
go mod init github.com/username/upbit-go-example
```

<br />

## SDK 설치

Upbit Go SDK를 설치합니다.

```bash
go get github.com/upbit-official/upbit-sdk-go
```

특정 버전으로 고정하려면 다음과 같이 설치합니다.

```bash
go get github.com/upbit-official/upbit-sdk-go@v0.9.0
```

<br />

## 인증 정보 설정

인증이 필요한 API를 사용하려면 API Key를 환경 변수로 설정합니다.

### macOS / Linux

```bash
export UPBIT_ACCESS_KEY=<your-access-key>
export UPBIT_SECRET_KEY=<your-secret-key>
```

### Windows PowerShell

```powershell
$env:UPBIT_ACCESS_KEY="your-access-key"
$env:UPBIT_SECRET_KEY="your-secret-key"
```

### Windows CMD

```bash
set UPBIT_ACCESS_KEY=your-access-key
set UPBIT_SECRET_KEY=your-secret-key
```

<br />

## 기본 사용 예제

다음 예제는 Upbit Go SDK를 사용하여 계좌 정보를 조회하는 코드입니다.

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/upbit-official/upbit-sdk-go"
	"github.com/upbit-official/upbit-sdk-go/option"
)

func main() {
	client := upbit.NewClient(
		option.WithAccessKey(os.Getenv("UPBIT_ACCESS_KEY")),
		option.WithSecretKey(os.Getenv("UPBIT_SECRET_KEY")),
	)

	accounts, err := client.Accounts.List(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%+v\n", accounts)
}
```

작성한 코드를 `main.go`로 저장한 뒤 실행합니다.

```bash
go run main.go
```

정상적으로 실행되면 계좌 정보가 출력됩니다.

<br />

## 예제 코드 실행

Upbit Go SDK는 주요 기능을 시나리오 형태로 구성한 예제 코드를 제공합니다.

각 예제는 독립 실행 가능한 `package main` 프로그램이며, 다음 형식으로 실행할 수 있습니다.

```bash
go run examples/<file>.go
```

### 인증이 필요 없는 예제

다음 예제는 API Key 없이 실행할 수 있습니다.

| 예제              | 설명                       |
| --------------- | ------------------------ |
| `quotation.go`  | 시세, 캔들, 체결, 호가 조회        |
| `indicators.go` | Quotation 데이터를 활용한 지표 계산 |

```bash
go run examples/quotation.go
go run examples/indicators.go
```

<br />

### 인증이 필요한 예제

다음 예제는 Upbit API Key가 필요합니다.

| 예제                   | 설명                          |
| -------------------- | --------------------------- |
| `orders.go`          | 주문 생성, 조회, 취소 흐름            |
| `orders_validate.go` | 주문 생성 테스트 API를 활용한 주문 유형 검증 |
| `deposits.go`        | 입금 주소 및 입금 내역 관리            |
| `withdrawals.go`     | 출금 가능 정보 및 출금 흐름            |
| `dca.go`             | 정기적 시장가 매수 예제               |
| `tp_sl.go`           | 익절/손절 자동 매도 예제              |

```go
UPBIT_ACCESS_KEY=<key>
UPBIT_SECRET_KEY=<secret>
go run examples/orders.go
```

<br />

## Dry run 모드

일부 예제는 기본적으로 Dry run 모드로 동작합니다.

Dry run 모드에서는 조회 중심으로 동작하며, 실제 주문이나 출금과 같은 쓰기 작업은 실행하지 않습니다.

실제 실행이 필요한 경우 `DRY_RUN=false`를 설정합니다.

```bash
DRY_RUN=false UPBIT_ACCESS_KEY=<key>
UPBIT_SECRET_KEY=<secret> 
go run examples/orders.go
```

실제 실행 시 주문, 출금 등 자산 변경이 발생할 수 있으므로 실행 전 내용을 반드시 확인하세요.

<br />

## 주요 기능

Upbit Go SDK는 다음과 같은 기능을 제공합니다.

* 시세 및 마켓 데이터 조회
* 계좌 및 자산 조회
* 주문 생성, 조회, 취소
* 입금 주소 및 입금 내역 조회
* 출금 가능 정보 및 출금 내역 조회
* 페이지네이션 목록 조회
* 요청 옵션 설정
* 에러 처리 및 디버깅

<br />

## 고급 기능

필요에 따라 다음 기능을 사용할 수 있습니다.

### 페이지네이션

목록 조회 API에서는 자동 페이지 순회를 사용할 수 있습니다.

```go
iter := client.Orders.ListOpenAutoPaging(context.TODO(), upbit.OrderListOpenParams{})

for iter.Next() {
	order := iter.Current()
	fmt.Printf("%+v\n", order)
}

if err := iter.Err(); err != nil {
	panic(err.Error())
}
```

### 요청 옵션

요청별로 헤더, 타임아웃, 재시도 설정 등을 지정할 수 있습니다.

```go
client.Accounts.List(
	context.TODO(),
	option.WithRequestTimeout(20*time.Second),
)
```

### 에러 처리

API 호출 중 오류가 발생하면 `upbit.Error` 타입으로 상태 코드와 응답 정보를 확인할 수 있습니다.

```go
accounts, err := client.Accounts.List(context.TODO())
if err != nil {
	var apierr *upbit.Error
	if errors.As(err, &apierr) {
		fmt.Println(apierr.StatusCode)
	}
	panic(err.Error())
}

fmt.Printf("%+v\n", accounts)
```

### 원시 응답 데이터 접근

응답 헤더나 상태 코드 등 원시 HTTP 응답 정보가 필요한 경우 사용할 수 있습니다.

```go
var response *http.Response

accounts, err := client.Accounts.List(
	context.TODO(),
	option.WithResponseInto(&response),
)
if err != nil {
	panic(err.Error())
}

fmt.Println(response.StatusCode)
fmt.Printf("%+v\n", accounts)
```

# Sibling pages

* [Python ](https://docs.upbit.com/kr/docs/python-sdk.md)
* [TypeScript](https://docs.upbit.com/kr/docs/typescript-sdk.md)