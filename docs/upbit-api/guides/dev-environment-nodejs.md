---
updatedAt: 2026-05-04T05:42:01.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# Node.js

Node.js 환경에서 Upbit Open API를 연동하기 위한 개발 환경 설정 방법을 안내합니다. 

## macOS 환경 설정

MacOS에서 Node.js를 설치하는 방법입니다. Homebrew라는 MacOS용 소프트웨어 패키지 관리자를 사용해 명령어로 간편하게 설치할 수 있습니다.

1. **Homebrew 설치**

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

설치를 완료하고 터미널에서 다음 명령어를 실행해 Homebrew의 버전을 확인합니다. 정상적으로 설치를 완료한 경우, Homebrew의 버전을 확인할 수 있습니다.

```shell
brew -v

Homebrew 4.5.3
```

<br />

2. **NVM 설치**

```shell
brew install nvm
```

NVM 설치 후, 쉘 프로필에 환경 변수를 추가합니다. 터미널에서 다음 명령어를 실행합니다.

* zsh 사용자

```shell
echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.zshrc
echo '[ -s "$(brew --prefix nvm)/nvm.sh" ] && source "$(brew --prefix nvm)/nvm.sh"' >> ~/.zshrc
```

* bash 사용자

```shell
echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.bash_profile
echo '[ -s "$(brew --prefix nvm)/nvm.sh" ] && source "$(brew --prefix nvm)/nvm.sh"' >> ~/.bash_profile
```

쉘 프로필을 업데이트하고 터미널에서 다음 명령어를 실행해 설정을 반영합니다.

* zsh 사용자

```shell
source ~/.zshrc
```

* bash 사용자

```shell
source ~/.bash_profile
```

설정을 반영하고 터미널에서 다음 명령어를 실행해 NVM의 버전을 확인합니다. 정상적으로 설치를 완료한 경우 NVM의 버전을 확인할 수 있습니다.

```shell
nvm --version

0.40.1
```

<br />

3. **Node.js 설치**

```shell
nvm install <version>

nvm install --lts
```

Node.js를 설치한 후 터미널에서 다음 명령어를 실행해 Node.js의 버전을 확인합니다. 정상적으로 설치를 완료한 경우 Node.js의 버전을 확인할 수 있습니다.

```shell
node -v

v22.14.0
```

<br />

## Windows 환경 설정

1. **Node.js 공식 웹사이트에서 설치 파일 다운로드**

Windows 운영체제에서 Node.js를 사용하기 위해서는 Node.js 공식 웹사이트에서 제공하는 설치 파일을 다운로드 받아야 합니다. 아래 링크를 클릭해 공식 웹사이트를 방문하고 설치 파일을 다운로드하세요. 설치 과정에서 **\[Add to PATH]** 옵션이 기본적으로 설정되어 있습니다. 이 옵션을 사용하는 경우 별도의 환경 변수 설정 없이 Node.js를 바로 사용할 수 있습니다.

* [Node.js 다운로드 바로가기](https://nodejs.org/ko/download)

Node.js를 설치한 후 터미널에서 다음 명령어를 실행해 Node.js의 버전을 확인합니다. 정상적으로 설치를 완료한 경우 Node.js의 버전을 확인할 수 있습니다.

```shell
node -v

v22.14.0
```

<br />

### HTTP 클라이언트 라이브러리 안내

Node.js 환경에서 REST API와 WebSocket을 호출하기 위해 사용할 수 있는 대표적인 라이브러리입니다.

#### REST API - `Axios` 라이브러리

Axios는 가장 널리 쓰이는 Promise 기반의 HTTP 클라이언트 라이브러리입니다. RESTful API 호출에 최적화되어 있으며 간결한 문법과 다양한 환경(브라우저, Node.js) 지원이 특징입니다. 다음과 같이 설치합니다.

```shell
npm install axios
```

REST API 호출 예제는 아래와 같습니다.

```javascript
const axios = require('axios');

axios.get('https://api.upbit.com/v1/ticker', {
  params: { markets: 'KRW-BTC' },
  headers: { 'accept': 'application/json' }
})
.then(response => {
  console.log(response.data[0].trade_price);
})
.catch(error => {
  console.error(error);
});
```

<br />

#### WebSocket - `ws` 라이브러리

ws는 Node.js에서 가장 널리 쓰이는 WebSocket 클라이언트/서버 구현 라이브러리입니다. 실시간 데이터 제공 API에 적합하며 이벤트 기반으로 메시지를 주고받을 수 있습니다. 다음과 같이 설치합니다.

```shell
npm install ws
```

WebSocket 사용 예시는 아래와 같습니다.

```javascript
const WebSocket = require('ws');

const ws = new WebSocket('wss://api.upbit.com/websocket/v1', {
  headers: {
    'accept': 'application/json'
  }
});

ws.on('open', () => {
  const subscribeMessage = [
    { ticket: 'test' },
    { type: 'ticker', codes: ['KRW-BTC'] }
  ];
  ws.send(JSON.stringify(subscribeMessage));
});

ws.on('message', (data) => {
  console.log('Received:', data.toString());
});

ws.on('close', () => {
  console.log('WebSocket connection closed');
});

ws.on('error', (error) => {
  console.error('WebSocket error:', error);
});
```

# Sibling pages

* [Python](https://docs.upbit.com/kr/docs/dev-environment-python.md)
* [Java](https://docs.upbit.com/kr/docs/dev-environment-java.md)