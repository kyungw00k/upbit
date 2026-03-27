# 매수/매도 가이드

Upbit CLI의 매수(`buy`)와 매도(`sell`) 명령 상세 가이드입니다.

## 주문 유형

### 지정가 주문 (Limit)

지정한 가격에 주문을 넣습니다. 해당 가격에 도달해야 체결됩니다.

```bash
upbit buy KRW-BTC -p 50000000 -V 0.001    # 5천만원에 0.001 BTC 매수
upbit sell KRW-BTC -p 55000000 -V 0.001   # 5500만원에 0.001 BTC 매도
```

### 시장가 주문 (Market)

현재 호가에 즉시 체결됩니다.

```bash
upbit buy KRW-BTC -t 100000               # 10만원어치 시장가 매수
upbit sell KRW-BTC -V 0.001               # 0.001 BTC 시장가 매도
```

- 매수: `--total`(-t)로 총액 지정
- 매도: `--volume`(-V)으로 수량 지정

### 최유리 지정가 주문 (Best)

상대방 최우선 호가에 지정가 주문을 넣습니다. `--best` 플래그를 사용합니다.

```bash
upbit buy KRW-BTC -V 0.001 --best         # 최유리 지정가 매수
upbit sell KRW-BTC -V 0.001 --best        # 최유리 지정가 매도
```

### 예약-지정가 주문 (Reserved)

감시가(`--watch`)에 도달하면 지정가 주문이 자동 발동됩니다.

```bash
upbit buy KRW-BTC --watch 49000000 -p 49500000 -V 0.001   # 4900만원 도달 시 4950만원에 매수
upbit sell KRW-BTC --watch 55000000 -p 54500000 -V 0.001  # 5500만원 도달 시 5450만원에 매도
```

## Time in Force (TIF)

`--tif` 플래그로 주문의 체결 조건을 지정합니다.

| TIF | 설명 | 동작 |
|-----|------|------|
| *(기본)* | GTC (Good Till Cancel) | 체결되거나 취소할 때까지 대기 |
| `ioc` | Immediate or Cancel | 즉시 체결 가능한 수량만 체결, 나머지 자동 취소 |
| `fok` | Fill or Kill | 전량 즉시 체결 가능하면 체결, 아니면 전량 취소 |
| `post_only` | Maker Only | 메이커 주문만 허용, 테이커로 즉시 체결되면 자동 취소 |

### 주문 유형별 TIF 조합

| 주문 유형 | 기본(GTC) | IOC | FOK | post_only |
|-----------|:---------:|:---:|:---:|:---------:|
| 지정가 (limit) | O | O | O | O |
| 시장가 (price/market) | O | - | - | - |
| 최유리 지정가 (best) | O | O | O | - |
| 예약-지정가 (watch) | O | - | - | - |

### 사용 예시

```bash
upbit buy KRW-BTC -p 50000000 -V 0.001 --tif ioc       # 지정가 IOC
upbit buy KRW-BTC -p 50000000 -V 0.001 --tif fok       # 지정가 FOK
upbit buy KRW-BTC -p 50000000 -V 0.001 --tif post_only # 메이커 전용
upbit buy KRW-BTC -V 0.001 --best --tif ioc             # 최유리 IOC
upbit sell KRW-BTC -V 0.001 --best --tif fok            # 최유리 FOK
```

## 퍼센트 주문

`--volume`(-V)과 `--total`(-t)에 퍼센트를 지정할 수 있습니다.

### 매수

```bash
upbit buy KRW-BTC -p 50000000 -V 50%      # KRW 잔고의 50%로 지정가 매수
upbit buy KRW-BTC -t 100%                 # KRW 잔고 전액 시장가 매수
upbit buy KRW-BTC -t 25%                  # KRW 잔고의 25% 시장가 매수
```

- 매수 시 **수수료가 자동 차감**됩니다
- 계산식: `사용 가능 금액 × 비율 ÷ (1 + 수수료율)`
- 잔고와 수수료율은 `GetOrderChance` API로 조회합니다

### 매도

```bash
upbit sell KRW-BTC -V 100%                # BTC 전량 시장가 매도
upbit sell KRW-BTC -p 55000000 -V 50%     # BTC의 50% 지정가 매도
```

- 매도 시 수수료는 체결 후 차감되므로 별도 계산 불필요
- 계산식: `보유 수량 × 비율`

## 가격 키워드

`--price`(-p)에 키워드를 사용하면 현재 시세를 자동 조회하여 지정가로 변환합니다.

| 키워드 | 설명 | 대응 값 |
|--------|------|---------|
| `now` | 현재가 (마지막 체결가) | `trade_price` |
| `open` | 금일 시가 | `opening_price` |
| `low` | 금일 저가 | `low_price` |
| `high` | 금일 고가 | `high_price` |

### 사용 예시

```bash
upbit buy KRW-BTC -p now -V 50%           # 현재가로 잔고 50% 지정가 매수
upbit buy KRW-BTC -p low -V 0.001         # 금일 저가로 지정가 매수
upbit sell KRW-BTC -p high -V 100%        # 금일 고가로 전량 지정가 매도
upbit sell KRW-BTC -p open -V 0.5         # 시가로 0.5 BTC 매도
```

가격 키워드는 퍼센트 주문과 함께 사용할 수 있습니다. 가격이 먼저 해석된 후 퍼센트가 계산됩니다.

## 수수료

| 주문 유형 | 수수료율 |
|-----------|---------|
| 지정가 / 시장가 / 최유리 지정가 | 0.05% |
| 예약-지정가 | 0.139% (부가세 포함) |

- 퍼센트 매수 시 수수료는 자동으로 차감되어 잔액 부족이 발생하지 않습니다
- 매도 수수료는 체결 후 차감됩니다

## 부가 기능

### 호가 단위 자동 보정

지정가 주문 시 입력한 가격이 호가 단위에 맞지 않으면 자동으로 보정됩니다.

- 매수: 내림 (매수자에게 유리)
- 매도: 올림 (매도자에게 유리)
- 예약 주문의 감시가도 자동 보정됩니다

### 테스트 주문

`--test` 플래그로 실제 체결 없이 주문을 시뮬레이션합니다.

```bash
upbit buy KRW-BTC -p 50000000 -V 0.001 --test
```

### 확인 프롬프트 스킵

`--force`(-f) 플래그로 확인 프롬프트를 건너뜁니다.

```bash
upbit buy KRW-BTC -t 100000 --force
```

## 주문 유형 요약

| 조합 | 주문 유형 | 수수료 |
|------|-----------|--------|
| `-p` + `-V` | 지정가 | 0.05% |
| `-t` | 시장가 매수 | 0.05% |
| `-V` (매도) | 시장가 매도 | 0.05% |
| `--best` + `-V` | 최유리 지정가 | 0.05% |
| `--watch` + `-p` + `-V` | 예약-지정가 | 0.139% |
