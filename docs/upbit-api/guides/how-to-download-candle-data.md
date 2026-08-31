---
updatedAt: 2026-05-04T05:48:42.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 캔들 API로 지정한 기간의 시세 데이터를 CSV 파일로 출력하기

업비트의 캔들 API로 지정한 기간의 캔들 데이터를 조회하여 CSV 파일로 출력하는 방법을 알아봅니다.

# 개요

본 가이드는 업비트 분 캔들 조회 API를 활용해 지정한 기간의 시세 데이터를 CSV 파일로 저장합니다. 각 레코드는 UTC 기준 캔들 시각과 함께 시가, 최고가, 최저가, 종가, 누적 거래 금액, 누적 거래 수량을 포함합니다. 캔들의 집계 간격(분)은 업비트에서 지원하는unit 값으로 설정할 수 있으며, 실행 후 결과는 아래와 같이 CSV형식으로 확인할 수 있습니다.

[block:image]
{
  "images": [
    {
      "image": [
        "https://files.readme.io/edb97ecc15e7f5f74efe541fa67de13157a88a185ee464b1e8b66653c469232d-candle_data_csv.png",
        "",
        ""
      ],
      "align": "center"
    }
  ]
}
[/block]

***

# 캔들 API로 지정한 기간의 시세 데이터를 CSV 파일로 출력하기

원활한 튜토리얼 진행을 위해 Python 3.7 버전 이상을 사용하시기 바랍니다. 업비트의 캔들 API는 호출 당 최대 200개의 캔들 데이터를 반환합니다. 사용자가 조회하려는 기간의 캔들 데이터가 200개를 초과할 경우 호출 기간을 계산해 캔들 API를 여러 번 호출하는 로직 추가가 필요합니다.

1. CSV로 출력할 캔들의 기간과 분 캔들 생성 시간 단위(이하 unit)를 입력합니다. 분 캔들 API는 unit 간격으로 생성된 캔들 데이터의 배열을 반환합니다.
2. 캔들 조회 시작 시각, unit, 그리고 응답으로 받을 캔들 데이터의 개수(이하 count)로 분 캔들 API의 호출 간격을 계산합니다.
3. 계산한 간격에 따라 분 캔들 API를 호출합니다.
4. 응답으로 받은 캔들 데이터 중 사용자가 지정한 기간 내에 생성된 캔들 데이터를 CSV 파일에 기록합니다.

[block:html]
{
  "html": "<div class=\"callout-section\">\n    <div class=\"callout-title\">\n      <i class=\"fa-solid fa-circle-exclamation\"></i> <b>Note</b>\n  </div>\n  캔들 데이터는 체결된 거래 데이터를 기반으로 생성되기 때문에 실제 거래가 체결되지 않은 경우 해당 시간대에는 캔들 데이터가 존재하지 않습니다.<br><br>\n\n 따라서, 캔들 API의 응답만을 CSV 파일에 기록하는 경우 CSV에 입력된 캔들 데이터의 간격이 일정하지 않을 수 있습니다. 본 가이드에서는 캔들 데이터가 없는 경우 빈 행을 CSV 파일에 기록합니다. 이러한 로직을 통해 사용자는 거래가 체결된 시간과 체결되지 않은 시간을 확인할 수 있습니다.<br><br>\n    마지막 캔들 행은 사용자가 지정한 종료 시각 “to”(exclusive)이 캔들 간격(unit)과 정렬되지 않은 경우, 종료 시각 이전에 가장 최근에 생성된 캔들 시각으로 기록됩니다.<br><br>\nunit=15일 때 캔들은 00:00, 00:15, 00:30 … 시각에 생성됩니다. 종료 시각이 2025-09-10T00:21:00처럼 15분 단위가 아니면, CSV의 마지막 행은 가장 최근 캔들 시각인 2025-09-10T00:15:00 의 데이터가 됩니다.\n\t</div>\n</div>"
}
[/block]

<br />

## 전체 코드

특정 기간의 캔들 데이터를 조회해 CSV 파일로 출력하는 전체 코드는 다음과 같습니다. 이 코드를 실행할 경우 현재 경로에 data 디렉토리를 생성하고 해당 디렉토리 내에 csv 파일을 생성합니다.

```python
import requests
from datetime import datetime, timedelta, timezone
from pathlib import Path
from time import sleep
import csv


class MinCandleCsvExporter():
    def __init__(self):
        self.csv_header = [
            "candle_date_time_utc", "trading_pair", "unit",
            "opening_price", "high_price", "low_price", "trade_price",
            "candle_acc_trade_volume", "candle_acc_trade_price"
        ]
        self.unit_limit = [
            1, 3, 5, 10, 15, 30, 60, 240
        ]

    def _get_candles_minutes(self, trading_pair, to_dt, unit):
        if unit not in self.unit_limit:
            raise ValueError("unit must be in 1, 3, 5, 10, 15, 30, 60, 240")
        path = f"https://api.upbit.com/v1/candles/minutes/{unit}"
        params = {
            "market": trading_pair,
            "to": to_dt,
            "count": 200
        }
        try:
            res = requests.get(path, params=params)
            res.raise_for_status()
            data = res.json()
        except Exception as e:
            raise RuntimeError(f"Error: {e}")

        if isinstance(data, list):
            return data[::-1]
        else:
            raise RuntimeError(f"Unexpected API response format: {data}")

    def _string_to_datetime(self, zulu_str):
        if zulu_str.endswith("Z"):
            zulu_str = zulu_str[:-1]
        return datetime.fromisoformat(zulu_str).replace(tzinfo=timezone.utc)

    def _datetime_to_string(self, dt):
        if dt.tzinfo is None:
            dt = dt.replace(tzinfo=timezone.utc)
        else:
            dt = dt.astimezone(timezone.utc)

        return dt.strftime("%Y-%m-%dT%H:%M:%SZ")

    def _floor_time_by_unit(self, dt, unit_min):
        if dt.tzinfo is None:
            dt = dt.replace(tzinfo=timezone.utc)
        else:
            dt = dt.astimezone(timezone.utc)

        if unit_min >= 60:
            total_minutes = dt.hour * 60 + dt.minute
            floored_total_minutes = (total_minutes // unit_min) * unit_min
            floored_hour = floored_total_minutes // 60
            floored_minute = floored_total_minutes % 60
            floored = dt.replace(
                hour=floored_hour, minute=floored_minute, second=0, microsecond=0)
        else:
            minute = (dt.minute // unit_min) * unit_min
            floored = dt.replace(minute=minute, second=0, microsecond=0)

        return floored.strftime("%Y-%m-%dT%H:%M:%SZ")

    def _get_csv_writer(self, csv_path):
        exists = csv_path.exists()
        csv_path.parent.mkdir(parents=True, exist_ok=True)
        f = open(csv_path, "a", newline="", encoding="utf-8")
        writer = csv.writer(f)
        if not exists:
            writer.writerow(self.csv_header)
        return f, writer

    def _write_empty_row(self, writer, ts_utc, trading_pair, unit):
        writer.writerow([
            self._datetime_to_string(ts_utc),
            trading_pair, unit, "", "", "", "", "", ""
        ])

    def _write_candle_row(self, writer, candle):
        ts_utc = candle.get("candle_date_time_utc")
        trading_pair = candle.get("market", "")
        unit = candle.get("unit", "")
        opening = candle.get("opening_price", "")
        high = candle.get("high_price", "")
        low = candle.get("low_price", "")
        close = candle.get("trade_price", "")
        acc_vol = candle.get("candle_acc_trade_volume", "")
        acc_val = candle.get("candle_acc_trade_price", "")
        if isinstance(ts_utc, str) and not ts_utc.endswith("Z"):
            ts_utc = ts_utc + "Z"
        writer.writerow([
            ts_utc, trading_pair, unit, opening, high, low, close, acc_vol, acc_val
        ])

    def _iter_write_range(self, trading_pair, unit, start_dt, end_dt, writer):
        from_dt = self._string_to_datetime(start_dt)
        to_dt = self._string_to_datetime(end_dt)
        unit_dt = timedelta(minutes=unit)
        interval = timedelta(minutes=unit * 200)
        to_params_dt = from_dt + interval
        last_dt = to_dt + interval

        while to_params_dt <= last_dt:
            print(f"querying {from_dt} to {to_params_dt} Candle Data...")
            candles = self._get_candles_minutes(
                trading_pair, self._datetime_to_string(to_params_dt), unit)
            expected_candle_dt = from_dt
            candle_idx = 0

            if not candles or len(candles) == 0:
                raise LookupError(
                    f"No candle data found for {trading_pair} from {self._datetime_to_string(from_dt)} to {self._datetime_to_string(to_params_dt)}. "
                    f"Please check the trading_pair and time range.")

            batch_end_dt = min(to_dt, to_params_dt - unit_dt)

            while expected_candle_dt <= batch_end_dt and candle_idx < len(candles):
                prev_expected_dt = expected_candle_dt

                current_candle_dt = self._string_to_datetime(
                    candles[candle_idx].get("candle_date_time_utc"))

                if current_candle_dt == expected_candle_dt:
                    self._write_candle_row(writer, candles[candle_idx])
                    candle_idx = candle_idx + 1
                    expected_candle_dt += unit_dt

                elif current_candle_dt < expected_candle_dt:
                    candle_idx = candle_idx + 1

                else:
                    self._write_empty_row(
                        writer, expected_candle_dt, trading_pair, unit)
                    expected_candle_dt += unit_dt

                if expected_candle_dt == prev_expected_dt and candle_idx >= len(candles):
                    self._write_empty_row(
                        writer, expected_candle_dt, trading_pair, unit)
                    expected_candle_dt += unit_dt

            while expected_candle_dt <= batch_end_dt:
                self._write_empty_row(
                    writer, expected_candle_dt, trading_pair, unit)
                expected_candle_dt += unit_dt

            from_dt = expected_candle_dt
            to_params_dt = from_dt + interval

            sleep(0.5)

    def run_export_minutes_csv(self, trading_pair, unit, from_dt_str, to_dt_str, csv_file_path):
        from_dt = self._string_to_datetime(from_dt_str)
        to_dt = self._string_to_datetime(to_dt_str)
        floored_from_dt = self._floor_time_by_unit(from_dt, unit)
        floored_to_dt = self._floor_time_by_unit(to_dt, unit)

        if from_dt > to_dt:
            raise ValueError("from_dt must be before to_dt")

        if to_dt > datetime.now(timezone.utc):
            raise ValueError("to_dt must be before current time")

        csv_path = Path(csv_file_path)
        f, writer = self._get_csv_writer(csv_path)
        try:
            print(
                f"{trading_pair}_{from_dt}_{to_dt}_{unit}m_candles.csv is creating...")
            self._iter_write_range(
                trading_pair, unit, floored_from_dt, floored_to_dt, writer)
            print(
                f"{trading_pair}_{from_dt}_{to_dt}_{unit}m_candles.csv has been created.")
        finally:
            f.close()


if __name__ == "__main__":

    exporter = MinCandleCsvExporter()

    trading_pair = "KRW-BTC"
    from_dt = "2025-09-01T00:00:00Z"  # UTC
    to_dt = "2025-09-10T10:00:00Z"  # UTC
    unit = 10
    csv_file = f"./data/{trading_pair}_{from_dt}_{to_dt}_{unit}m_candles.csv"

    exporter.run_export_minutes_csv(
        trading_pair, unit, from_dt, to_dt, csv_file_path=csv_file)
```

***

## 세부 가이드

### MinCandleCsvExporter 모듈 구현

사용자가 입력한 기간 동안 캔들 데이터를 조회하고 CSV 파일에 입력하는 메서드를 구현할 클래스를 생성합니다.

<br />

### 초기 설정

CSV 파일 내 작성할 헤더의 값 그리고 업비트 분 캔들 API에서 허용하는 unit을 설정합니다. MinCandleCsvExporter 모듈로 인스턴스를 생성할 때 마다 다음과 같이 기본 헤더와 unit의 단위가 설정됩니다.

```python
class MinCandleCsvExporter():
    def __init__(self):
        self.csv_header = [
            "candle_date_time_utc", "trading_pair", "unit",
            "opening_price", "high_price", "low_price", "trade_price",
            "candle_acc_trade_volume", "candle_acc_trade_price"
        ]
        self.unit_limit = [
            1, 3, 5, 10, 15, 30, 60, 240
        ]
```

***

### CSV 파일 메서드

**캔들 데이터 추출**

MinCandleCsvExporter 클래스에서 실제 사용자가 호출하는 메서드 입니다. 조회할 캔들 데이터의 페어와 unit, 조회 시작 시각, 종료 시각, csv 파일의 경로를 입력받아 캔들 API 조회 및 캔들 데이터를 CSV 파일에 입력하는 동작을 수행합니다.

```python
    def run_export_minutes_csv(self, trading_pair, unit, from_dt_str, to_dt_str, csv_file_path):
        from_dt = self._string_to_datetime(from_dt_str)
        to_dt = self._string_to_datetime(to_dt_str)
        floored_from_dt = self._floor_time_by_unit(from_dt, unit)
        floored_to_dt = self._floor_time_by_unit(to_dt, unit)

        if from_dt > to_dt:
            raise ValueError("from_dt must be before to_dt")

        if to_dt > datetime.now(timezone.utc):
            raise ValueError("to_dt must be before current time")

        csv_path = Path(csv_file_path)
        f, writer = self._get_csv_writer(csv_path)
        try:
            print(
                f"{trading_pair}_{from_dt}_{to_dt}_{unit}m_candles.csv is creating...")
            self._iter_write_range(
                trading_pair, unit, floored_from_dt, floored_to_dt, writer)
            print(
                f"{trading_pair}_{from_dt}_{to_dt}_{unit}m_candles.csv has been created.")
        finally:
            f.close()
```

**CSV 파일 입력 객체 생성**

사용자가 입력한 CSV 파일의 경로를 확인해 CSV 파일 존재 여부를 확인하는 메서드입니다. CSV 파일이 있을 경우 해당 파일을 열고 파일 객체와 writer 객체를 반환합니다. 이 객체들을 사용해 해당 CSV 파일에 데이터를 입력할 수 있습니다. 파일이 존재하지 않는 경우, 상위 폴더를 만들고 CSV 파일을 생성해 파일 객체와 writer 객체를 반환합니다.

```python
    def _get_csv_writer(self, csv_path):
        exists = csv_path.exists()
        csv_path.parent.mkdir(parents=True, exist_ok=True)
        f = open(csv_path, "a", newline="", encoding="utf-8")
        writer = csv.writer(f)
        if not exists:
            writer.writerow(self.csv_header)
        return f, writer
```

***

### API 호출 메서드

**분 캔들 API 호출**

업비트의 분 캔들 API를 호출하는 메서드입니다. 사용자가 입력한 페어의 캔들을 조회하며 unit 간격으로 생성된 캔들 데이터의 배열을 반환합니다. 업비트의 캔들 API가 지원하는 to 파라미터를 사용해 캔들 데이터의 조회 구간을 설정할 수 있으며 이를 바탕으로 캔들 API 호출 간격을 계산합니다.

```python
    def _get_candles_minutes(self, trading_pair, to_dt, unit):
        if unit not in self.unit_limit:
            raise ValueError("unit must be in 1, 3, 5, 10, 15, 30, 60, 240")
        path = f"https://api.upbit.com/v1/candles/minutes/{unit}"
        params = {
            "market": trading_pair,
            "to": to_dt,
            "count": 200
        }
        try:
            res = requests.get(path, params=params)
            res.raise_for_status()
            data = res.json()
        except Exception as e:
            raise RuntimeError(f"Error: {e}")

        if isinstance(data, list):
            return data[::-1]
        else:
            raise RuntimeError(f"Unexpected API response format: {data}")
```

**기간별 분 캔들 API 호출 및 CSV 파일 출력**

조회 시작, 종료 시각과 unit을 기반으로 분 캔들 API를 호출하여 CSV 파일에 기록하는 메서드입니다. 캔들 데이터 존재 시 시간 순서대로 기록하고, 누락된 시간대는 빈 값으로 채워 연속적인 시계열 데이터를 보장합니다. API 호출 간 0.5초 대기로 안정성을 확보합니다.

```python
    def _iter_write_range(self, trading_pair, unit, start_dt, end_dt, writer):
        from_dt = self._string_to_datetime(start_dt)
        to_dt = self._string_to_datetime(end_dt)
        unit_dt = timedelta(minutes=unit)
        interval = timedelta(minutes=unit * 200)
        to_params_dt = from_dt + interval
        last_dt = to_dt + interval

        while to_params_dt <= last_dt:
            print(f"querying {from_dt} to {to_params_dt} Candle Data...")
            candles = self._get_candles_minutes(
                trading_pair, self._datetime_to_string(to_params_dt), unit)
            expected_candle_dt = from_dt
            candle_idx = 0

            if not candles or len(candles) == 0:
                raise LookupError(
                    f"No candle data found for {trading_pair} from {self._datetime_to_string(from_dt)} to {self._datetime_to_string(to_params_dt)}. "
                    f"Please check the trading_pair and time range.")

            batch_end_dt = min(to_dt, to_params_dt - unit_dt)

            while expected_candle_dt <= batch_end_dt and candle_idx < len(candles):
                prev_expected_dt = expected_candle_dt

                current_candle_dt = self._string_to_datetime(
                    candles[candle_idx].get("candle_date_time_utc"))

                if current_candle_dt == expected_candle_dt:
                    self._write_candle_row(writer, candles[candle_idx])
                    candle_idx = candle_idx + 1
                    expected_candle_dt += unit_dt

                elif current_candle_dt < expected_candle_dt:
                    candle_idx = candle_idx + 1

                else:
                    self._write_empty_row(
                        writer, expected_candle_dt, trading_pair, unit)
                    expected_candle_dt += unit_dt

                if expected_candle_dt == prev_expected_dt and candle_idx >= len(candles):
                    self._write_empty_row(
                        writer, expected_candle_dt, trading_pair, unit)
                    expected_candle_dt += unit_dt

            while expected_candle_dt <= batch_end_dt:
                self._write_empty_row(
                    writer, expected_candle_dt, trading_pair, unit)
                expected_candle_dt += unit_dt

            from_dt = expected_candle_dt
            to_params_dt = from_dt + interval

            sleep(0.5)
```

<br />

***

### 캔들 데이터 출력 메서드

**실제 캔들 데이터 출력**

실제 캔들 데이터의 값을 추출해 CSV 파일에 입력하는 메서드입니다. 캔들 데이터에서 생성 시각, 페어, unit, 시가, 최고가, 최저가, 종가, 누적 거래 수량, 누적 거래 금액 등을 추출해 CSV 파일에 입력합니다.

```python
    def _write_candle_row(self, writer, candle):
        ts_utc = candle.get("candle_date_time_utc", "")
        trading_pair = candle.get("market", "")
        unit = candle.get("unit", "")
        opening = candle.get("opening_price", "")
        high = candle.get("high_price", "")
        low = candle.get("low_price", "")
        close = candle.get("trade_price", "")
        acc_vol = candle.get("candle_acc_trade_volume", "")
        acc_val = candle.get("candle_acc_trade_price", "")
        if isinstance(ts_utc, str) and not ts_utc.endswith("Z"):
            ts_utc = ts_utc + "Z"
        writer.writerow([
            ts_utc, trading_pair, unit, opening, high, low, close, acc_vol, acc_val
        ])
```

**빈 캔들 데이터 출력**

빈 캔들 데이터를 CSV 파일에 입력하는 메서드입니다. 생성된 캔들 데이터가 없는 경우, 이 메서드를 호출해 CSV 파일에 빈 값을 입력합니다.

```python
    def _write_empty_row(self, writer, ts_utc, trading_pair, unit):
        writer.writerow([
            self._datetime_to_string(ts_utc),
            trading_pair, unit, "", "", "", "", "", ""
        ])
```

***

### 시간 변환 메서드

**문자열을 datetime 객체로 변환**

문자열을 datetime 객체로 변환하는 메서드입니다. 사용자가 문자열로 입력한 시간을 datetime 객체로 변환해 시간의 대소 비교, 시간 연산 등에 사용합니다.

```python
    def _string_to_datetime(self, zulu_str):
        if zulu_str.endswith("Z"):
            zulu_str = zulu_str[:-1]
        return datetime.fromisoformat(zulu_str).replace(tzinfo=timezone.utc)
```

**datetime 객체를 문자열로 변환**

datetime 객체를 문자열로 변환하는 메서드입니다. 업비트의 캔들 API는 시간 관련 파라미터를 문자열로 입력받기 때문에 시간 연산을 마친 후 해당 메서드로 datetime 객체를 다시 문자열로 변환해 API를 호출해야 합니다.

```python
    def _datetime_to_string(self, dt):
        if dt.tzinfo is None:
            dt = dt.replace(tzinfo=timezone.utc)
        else:
            dt = dt.astimezone(timezone.utc)

        return dt.strftime("%Y-%m-%dT%H:%M:%SZ")
```

**캔들 생성 단위에 따른 시간 계산**

사용자 입력 시간을 캔들 unit 단위로 버림 처리하여 유효한 캔들 시간으로 변환하는 메서드입니다. 예를 들어 unit이 5분이고 파라미터로 입력한 시각이 2025-09-10T10:22:00인 경우, 해당 시각으로 조회할 수 있는 마지막 캔들 데이터의 생성 시각은 2025-09-10T10:20:00입니다. 이 메서드를 사용해 캔들 API 호출 시 더욱 정확한 캔들 데이터를 조회할 수 있습니다.

```python
    def _floor_time_by_unit(self, dt, unit_min):
        if dt.tzinfo is None:
            dt = dt.replace(tzinfo=timezone.utc)
        else:
            dt = dt.astimezone(timezone.utc)

        if unit_min >= 60:
            total_minutes = dt.hour * 60 + dt.minute
            floored_total_minutes = (total_minutes // unit_min) * unit_min
            floored_hour = floored_total_minutes // 60
            floored_minute = floored_total_minutes % 60
            floored = dt.replace(
                hour=floored_hour, minute=floored_minute, second=0, microsecond=0)
        else:
            minute = (dt.minute // unit_min) * unit_min
            floored = dt.replace(minute=minute, second=0, microsecond=0)

        return floored.strftime("%Y-%m-%dT%H:%M:%SZ")
```

# Sibling pages

* [REST API를 이용한 시세 조회](https://docs.upbit.com/kr/docs/upbit-quotation-restful-api.md)
* [WebSocket으로 실시간 캔들 차트 생성](https://docs.upbit.com/kr/docs/how-to-subscribe-quotation-stream.md)