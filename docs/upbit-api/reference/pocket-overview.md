---
updatedAt: 2026-08-25T12:21:10.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 포켓(Pocket)

'포켓(Pocket)' 기능을 통해 더욱 효율적이고 체계적인 자산 관리가 가능해졌습니다. 포켓의 기본 개념과 API Key 발급 시 설정할 수 있는 권한에 대해 안내해 드립니다.

### '포켓(Pocket)'이란?

포켓(Pocket)은 업비트 계정 안에서 자산을 목적별로 나누어 투자하고 관리할 수 있는 기능입니다. 기본으로 제공되는 ‘메인포켓’ 외에 서브포켓을 추가로 만들어 이용할 수 있습니다.

* 메인포켓(Main Pocket) : 업비트 계정을 만들 때 기본으로 제공되는 기본 포켓입니다. 외부 입출금을 포함한 계정의 기본 거래를 할 수 있고, 계정 내 전체 포켓을 관리할 수 있는 마스터 권한을 받을 수 있습니다.
* 서브포켓(Sub Pocket): 특정 목적에 따라 사용자가 추가로 만드는 포켓입니다. API Key에 부여된 권한 안에서 독립적으로 동작하며, 외부 입출금은 제한됩니다.

<br />

### 기능별 API Key 권한 및 지원 포켓 안내

API Key 발급 시 선택할 수 있는 권한과 각 권한의 주요 기능, 그리고 권한별 지원 포켓은 아래 표에서 확인할 수 있습니다. 특히 포켓 기능과 관련해 새롭게 추가된 포켓관리, 자산이전 권한은 설정 가능한 포켓이 다르므로 유의해 주시기 바랍니다.<br /> <Anchor target="_blank" href="https://upbit.com/mypage/open_api_management">포켓 별 API Key 발급 받으러 가기→</Anchor>

<style>
  .api-permission-table {
    width: 100%;
    table-layout: fixed;
    border-collapse: separate;
    border-spacing: 0 6px;
    font-size: 14px;
    color: #1f2937;
    background: #ffffff;
    border: 0 !important;
  }
  .api-permission-table th,
  .api-permission-table td {
    border: 0 !important;
    background: #ffffff;
  }
  .api-permission-table th {
    padding: 10px 14px;
    text-align: left;
    font-weight: 700;
    color: #111827;
    white-space: nowrap;
  }
  .api-permission-table td {
    padding: 10px 14px;
    vertical-align: middle;
  }
  .permission-col {
    width: 22%;
    min-width: 120px;
    word-break: keep-all;
  }
  .permission-name {
    font-weight: 500;
    color: #1f2937;
    word-break: keep-all;
  }
  .permission-name.pocket {
    font-weight: 700;
    color: #0062df;
  }
  .feature-col {
    width: 54%;
    line-height: 1.65;
    word-break: keep-all;
  }
  .pocket-col {
    width: 12%;
    text-align: center !important;
    white-space: nowrap;
  }
  .check-box {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 17px;
    height: 17px;
    border-radius: 5px;
    border: 0;
    background: #eef2f7;
    color: transparent;
    font-size: 12px;
    font-weight: 700;
    line-height: 1;
  }
  .check-box.checked {
    background: #0062df;
    color: #ffffff;
  }
</style>

<table className="api-permission-table">
  <colgroup>
    <col className="permission-col" />
    <col className="feature-col" />
    <col className="pocket-col" />
    <col className="pocket-col" />
  </colgroup>
  <thead>
    <tr>
      <th>권한명</th>
      <th>주요기능</th>
      <th className="pocket-col">메인포켓</th>
      <th className="pocket-col">서브포켓</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td className="permission-name pocket">포켓관리</td>
      <td className="feature-col">모든 포켓의 자산 조회, 자산 이전 실행 및 내역 조회, 포켓 정보 조회</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box"></span></td>
    </tr>
    <tr>
      <td className="permission-name pocket">자산이전</td>
      <td className="feature-col">해당 포켓 기준의 포켓 간 자산 이전 실행 및 이전 내역 조회</td>
      <td className="pocket-col"><span className="check-box"></span></td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
    </tr>
    <tr>
      <td className="permission-name">자산조회</td>
      <td className="feature-col">해당 포켓의 보유자산 및 잔고 조회</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
    </tr>
    <tr>
      <td className="permission-name">주문조회</td>
      <td className="feature-col">해당 포켓의 체결 및 미체결 주문 내역 조회</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
    </tr>
    <tr>
      <td className="permission-name">주문하기</td>
      <td className="feature-col">해당 포켓의 매수/매도 주문 접수 및 취소</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
    </tr>
    <tr>
      <td className="permission-name">출금조회</td>
      <td className="feature-col">디지털 자산 및 원화 출금 내역 조회</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box"></span></td>
    </tr>
    <tr>
      <td className="permission-name">출금하기</td>
      <td className="feature-col">디지털 자산 및 원화 출금 요청</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box"></span></td>
    </tr>
    <tr>
      <td className="permission-name">입금조회</td>
      <td className="feature-col">디지털 자산 및 원화 입금 내역 조회</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box"></span></td>
    </tr>
    <tr>
      <td className="permission-name">입금하기</td>
      <td className="feature-col">원화 입금 요청 및 트래블룰 검증 등</td>
      <td className="pocket-col"><span className="check-box checked">✓</span></td>
      <td className="pocket-col"><span className="check-box"></span></td>
    </tr>
  </tbody>
</table>

<br />

<Callout icon="📘" theme="info">
  거래 및 자산 관리(Exchange) API의 Rate Limit은 포켓(Pocket) 단위로 적용됩니다.<br />동일 포켓에서 발급된 API Key는 요청 한도를 공유하며, 서로 다른 포켓은 각각 독립적인 요청 한도를 가집니다.<br />자세한 내용은 [요청 수 제한(Rate Limits)](/reference/rate-limits) 문서를 참고하세요.
</Callout>

<br />

# Sub pages

* [포켓 정보 조회](https://docs.upbit.com/kr/reference/list-pockets.md)
* [포켓별 API Key 목록 조회](https://docs.upbit.com/kr/reference/list-pocket-api-keys.md)
* [서브포켓 잔고 조회](https://docs.upbit.com/kr/reference/get-sub-pocket-balance.md)
* [메인포켓 자산 이전](https://docs.upbit.com/kr/reference/universal-transfer.md)
* [메인포켓 자산 이전 목록 조회](https://docs.upbit.com/kr/reference/list-universal-transfers.md)
* [서브포켓 자산 이전](https://docs.upbit.com/kr/reference/transfer.md)
* [서브포켓 자산 이전 목록 조회](https://docs.upbit.com/kr/reference/list-transfers.md)

# Sibling pages

* [자산(Asset)](https://docs.upbit.com/kr/reference/자산asset.md)
* [주문(Order)](https://docs.upbit.com/kr/reference/exchange.md)
* [출금(Withdrawal)](https://docs.upbit.com/kr/reference/출금withdrawal.md)
* [입금(Deposit)](https://docs.upbit.com/kr/reference/입금deposit.md)
* [서비스 정보(Service)](https://docs.upbit.com/kr/reference/서비스-정보service.md)