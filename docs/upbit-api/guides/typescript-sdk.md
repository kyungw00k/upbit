---
updatedAt: 2026-05-05T09:58:02.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# TypeScript

 TypeScript 환경에서 Upbit SDK 라이브러리를 사용해 업비트 API를 호출하는 방법을 안내합니다.

## 시작하기

Upbit TypeScript SDK는 TypeScript 및 Node.js 애플리케이션에서 Upbit API를 사용할 수 있도록 제공되는 공식 라이브러리입니다.

요청 파라미터와 응답 데이터에 대한 타입 정의를 제공하며, API 호출에 필요한 인증 설정, 요청 구성, 응답 처리 등을 TypeScript 코드에서 일관된 방식으로 사용할 수 있습니다.

이 가이드에서는 TypeScript 환경에서 Upbit SDK를 사용하여 API를 호출하는 방법을 안내합니다.

<br />

## SDK 공식 문서

Upbit TypeScript SDK는 공식 문서를 제공합니다. 아래 링크를 클릭해 공식 웹사이트로 이동할 수 있습니다.

* [시작 가이드](https://github.com/upbit-official/upbit-sdk-typescript/blob/main/README_KR.md) — Upbit TypeScript SDK의 설치 및 기본 사용 방법
* [Upbit SDK API Reference](https://github.com/upbit-official/upbit-sdk-typescript/blob/main/api.md) — SDK가 지원하는 전체 API 목록
* [SDK 예제 코드](https://github.com/upbit-official/upbit-sdk-typescript/tree/main/examples) — 예제 코드 및 상세 사용법

<br />

## TypeScript 연동 가이드

* 최소 버전: TypeScript 4.9 이상
* Node.js 22 LTS 이상 권장

### **프로젝트 초기화 및 SDK 설치**

TypeScript는 프로젝트별 환경 구성을 권장합니다. 다음과 같이 프로젝트를 생성하고 SDK를 설치합니다.

#### 1. 프로젝트 디렉토리 생성

```bash
mkdir upbit_sdk_project
cd upbit_sdk_project

npm init -y
```

#### 2. TypeScript 및 SDK 설치

```bash
npm install @upbit-official/upbit-sdk
npm install -D typescript ts-node @types/node
```

#### 3. TypeScript 설정 파일 생성

```bash
npx tsc --init
```

<br />

## 클라이언트 인스턴스 설정

Upbit SDK는 인증 정보와 환경 설정을 포함한 상태를 유지할 수 있도록 설계되어 있습니다. 따라서 아래와 같이 클라이언트 인스턴스를 생성하여 사용하는 방식이 일반적입니다.

> **환경변수를 통한 인증**: `UPBIT_ACCESS_KEY`와 `UPBIT_SECRET_KEY` 환경변수를 설정하면 클라이언트 생성 시 키를 직접 전달하지 않아도 자동으로 인증 정보를 읽어옵니다. 소스코드에 API Key를 직접 입력하는 것은 보안상 권장하지 않으며, 환경변수 또는 `.env` 파일을 사용하세요.

```typescript
import Upbit from '@upbit-official/upbit-sdk';

const client = new Upbit({
  accessKey: process.env['UPBIT_ACCESS_KEY'],
  secretKey: process.env['UPBIT_SECRET_KEY'],
});

```

아래 링크를 클릭해 Upbit SDK를 통해 사용할 수 있는 업비트 API를 확인할 수 있습니다.

* **Upbit SDK API Reference**: <https://github.com/upbit-official/upbit-sdk-typescript/blob/main/api.md>

<br />

### 인스턴스 설정 확인

인스턴스가 정상적으로 설정되었는지 확인하기 위해 다음과 같이 코드를 작성하고 실행합니다. 정상적으로 설정된 경우 API의 응답을 반환하고, 설정이 올바르지 않은 경우 에러를 반환합니다.

#### **인증이 필요한 API 호출**

```typescript
import Upbit from '@upbit-official/upbit-sdk';

async function getBalance(): Promise<void> {
  const client = new Upbit({
    accessKey: process.env['UPBIT_ACCESS_KEY'],
    secretKey: process.env['UPBIT_SECRET_KEY'],
  });
  const result = await client.accounts.list();
  console.log(result);
}

async function main(): Promise<void> {
  await getBalance();
}

main().catch(console.error);
```

#### **인증 없이 API 호출**

```typescript
import Upbit from '@upbit-official/upbit-sdk';

async function listTradingPairsDeafult(): Promise<void> {
  const client = new Upbit();
  const result = await client.tradingPairs.list({
    is_details: false,
  });
  console.log(result);
}

async function main(): Promise<void> {
  await listTradingPairsDeafult();
}

main().catch(console.error);
```

<br />

## 에러 핸들링

API 호출 시 인증 실패, 요청 초과 등의 에러가 발생할 수 있습니다. 다음과 같이 예외를 처리할 수 있습니다.

```typescript
const accounts = await client.accounts.list().catch(async (err) => {
  if (err instanceof Upbit.APIError) {
    console.log(err.status);  // 404
    console.log(err.name);    // "not_found"  (Upbit API error name, or null)
    console.log(err.message); // "404 [not_found] no Route matched"
    console.log(err.headers); // {server: 'nginx', ...}
  } else {
    throw err;
  }
});
```

### 에러 타입

| 상태 코드  | 오류 타입                      |
| ------ | -------------------------- |
| 400    | `BadRequestError`          |
| 401    | `AuthenticationError`      |
| 403    | `PermissionDeniedError`    |
| 404    | `NotFoundError`            |
| 418    | `RateLimitPenaltyError`    |
| 422    | `UnprocessableEntityError` |
| 429    | `RateLimitError`           |
| > =500 | `InternalServerError`      |
| N/A    | `APIConnectionError`       |

<br />

## 추가 기능

### 재시도 설정

일부 오류에 대해서는 기본적으로 최대 2회까지 자동 재시도됩니다.

```typescript
//모든 요청의 기본값 설정:
const client = new Upbit({
  maxRetries: 0, //기본값은 2
});

//또는 요청별 설정:
await client.accounts.list({
  maxRetries: 5,
});

```

<br />

### 타임아웃 설정

기본 타임아웃은 60초이며, 다음과 같이 변경할 수 있습니다.

```typescript TypeScript
//모든 요청의 기본값 설정:
const client = new Upbit({
  timeout: 20 * 1000, //20초 (기본값은 1분)
});

//요청별 설정:
const client.accounts.list({
	timeout: 5 * 1000,
}); 


```

<br />

### 자동 페이지네이션

Upbit API의 목록 조회 메서드는 페이지네이션을 지원합니다. `for await ... of`구문을 사용하면 모든 페이지의 항목을 순회할 수 있습니다.

```typescript TypeScript
async function fetchAllOrders(params) {
	const allOrders = [];
	// 필요에 따라 자동으로 다음 페이지를 가져옵니다.
	for await (const order of client.orders.listOpen()) {
		allOrders.push(order);
	}
	return allOrders;
}

```

한 번에 한 페이지씩 요청할 수도 있습니다.

```typescript TypeScript
let page = await client.orders.listOpen();
for (const order of page.items) {
	console.log(order);
}

// 수동 페이지네이션을 위한 편의 메서드:
while (page.hasNextPage()) {
	page = await page.getNextPage();
	// ...
}
```

# Sibling pages

* [Python ](https://docs.upbit.com/kr/docs/python-sdk.md)
* [Go](https://docs.upbit.com/kr/docs/go-sdk.md)