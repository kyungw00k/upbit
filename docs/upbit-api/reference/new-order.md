---
updatedAt: 2026-08-21T12:05:08.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 주문 생성

특정 페어를 매수/매도하기 위한 주문을 생성합니다.

### 주문 유형(ord\_type)

사용 가능한 주문 유형은 다음과 같습니다.

#### 지정가 주문

지정가 주문은 사용자가 직접 설정한 매수/매도 단가, 또는 더 유리한 가격에 호가가 도달한 경우에만 체결되는 주문 유형입니다. 체결 단가의 상한/하한을 통제할 수 있지만, 시장 가격이 지정한 단가에 도달하지 않을 수 있으므로 체결을 보장할 수 없습니다.

<div className="accordion">
  <input type="checkbox" id="limit-bid" />
  <label for="limit-bid">
    <i className="fa-solid fa-circle-info"></i>
    지정가 매수/매도 주문 생성 요청 파라미터 사용 예시
  </label>

  <div className="accordion-content">

아래 표를 참고하여 지정가 매수/매도 주문 생성 요청 시 사용 가능한 파라미터를 쉽게 확인할 수 있습니다. 각 파라미터에 대한 상세한 설명은 하단 Request Body를 참고해주세요.

<table className="custom-table">
<thead>
  <tr>
    <th>파라미터</th>
    <th>필수 여부</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col">market</td>
    <td>Required</td>
    <td>페어 코드. <code>KRW-BTC</code> 형식으로 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">side</td>
    <td>Required</td>
    <td>매수시 <code>bid</code>, 매도시 <code>ask</code>로 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">ord_type</td>
    <td>Required</td>
    <td><code>limit</code></td>
  </tr>
  <tr>
    <td className="code-col">volume</td>
    <td>Required</td>
    <td>주문 수량. 0.1 입력시 지정가로 0.1개의 자산을 매수/매도합니다.</td>
  </tr>
  <tr>
    <td className="code-col">price</td>
    <td>Required</td>
    <td>호가 자산 기준 주문 단가. 예를 들어, KRW-BTC 페어에서 BTC 1개당 1억원(KRW)으로 매수/매도하는 경우 100000000을 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">time_in_force</td>
    <td>Optional</td>
    <td><code>ioc</code>, <code>fok</code>, <code>post_only</code>. post_only 옵션은 smp_type 옵션과 함께 사용할 수 없습니다.</td>
  </tr>
  <tr>
    <td className="code-col">smp_type</td>
    <td>Optional</td>
    <td>자전거래 체결 방지 옵션. <code>cancel_maker</code>,<code>cancel_taker</code>,<code>reduce</code></td>
  </tr>
  <tr>
    <td className="code-col">identifier</td>
    <td>Optional</td>
    <td>조회, 삭제시 사용할 수 있는 사용자 지정 주문 ID.</td>
  </tr>
</tbody>

</table>

  </div>
</div>

#### 시장가 주문

시장가 주문은 현재 시장에서 가장 유리한 가격으로 즉시 체결되는 주문 유형입니다. 빠른 체결이 보장되지만, 시장 상황에 따라 체결 가격이 변동될 수 있습니다.

<div className="accordion">
  <input type="checkbox" id="market-bid" />
  <label for="market-bid">
    <i className="fa-solid fa-circle-info"></i>
    시장가 매수 주문 생성 요청 파라미터 사용 예시
  </label>

  <div className="accordion-content">

아래 표를 참고하여 시장가 매수 주문 생성 요청 시 사용 가능한 파라미터를 쉽게 확인할 수 있습니다. <code>volume</code> 파라미터를 입력하지 않습니다(value에 null로 입력 또는 key 자체를 제외). 각 파라미터에 대한 상세한 설명은 하단 Request Body를 참고해주세요.

<table className="custom-table">
<thead>
  <tr>
    <th>파라미터</th>
    <th>필수 여부</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col">market</td>
    <td>Required</td>
    <td>페어 코드. <code>KRW-BTC</code> 형식으로 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">side</td>
    <td>Required</td>
    <td><code>bid</code></td>
  </tr>
  <tr>
    <td className="code-col">ord_type</td>
    <td>Required</td>
    <td><code>price</code></td>
  </tr>
  <tr>
    <td className="code-col">price</td>
    <td>Required</td>
    <td>호가 자산 기준 주문 총액. 예를 들어, KRW-BTC 페어에서 100000000을 입력하는 경우 시장가로 1억 원어치의 BTC 수량이 매수됩니다.</td>
  </tr>
  <tr>
    <td className="code-col">smp_type</td>
    <td>Optional</td>
    <td>자전거래 체결 방지 옵션. <code>cancel_maker</code>,<code>cancel_taker</code>,<code>reduce</code></td>
  </tr>
  <tr>
    <td className="code-col">identifier</td>
    <td>Optional</td>
    <td>조회, 삭제시 사용할 수 있는 사용자 지정 주문 ID.</td>
  </tr>
</tbody>

</table>

  </div>
</div>

<div className="accordion">
  <input type="checkbox" id="market-ask" />
  <label for="market-ask">
    <i className="fa-solid fa-circle-info"></i>
    시장가 매도 주문 생성 요청 파라미터 사용 예시
  </label>

  <div className="accordion-content">

아래 표를 참고하여 시장가 매도 주문 생성 요청 시 사용 가능한 파라미터를 쉽게 확인할 수 있습니다. <code>price</code> 파라미터를 입력하지 않습니다(value에 null로 입력 또는 key 자체를 제외). 각 파라미터에 대한 상세한 설명은 하단 Request Body를 참고해주세요.

<table className="custom-table">
<thead>
  <tr>
    <th>파라미터</th>
    <th>필수 여부</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col">market</td>
    <td>Required</td>
    <td>페어 코드. <code>KRW-BTC</code> 형식으로 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">side</td>
    <td>Required</td>
    <td><code>ask</code></td>
  </tr>
  <tr>
    <td className="code-col">ord_type</td>
    <td>Required</td>
    <td><code>market</code></td>
  </tr>
  <tr>
    <td className="code-col">volume</td>
    <td>Required</td>
    <td>매도 주문 수량. 예를 들어, KRW-BTC 페어에서 0.1을 입력하는 경우 시장가로 0.1개의 BTC 수량이 매도됩니다.</td>
  </tr>
  <tr>
    <td className="code-col">smp_type</td>
    <td>Optional</td>
    <td>자전거래 체결 방지 옵션. <code>cancel_maker</code>,<code>cancel_taker</code>,<code>reduce</code></td>
  </tr>
  <tr>
    <td className="code-col">identifier</td>
    <td>Optional</td>
    <td>조회, 삭제시 사용할 수 있는 사용자 지정 주문 ID.</td>
  </tr>
</tbody>

</table>

  </div>
</div>

#### 최유리 지정가 주문

최유리 지정가 주문은 현재 시장에서 가장 유리한 상대 호가를 가격으로 하는 주문 유형입니다. 전량 체결을 항상 보장할 수는 없으나, 빠르게 유리한 가격으로 호가창에 진입하고 싶은 경우 유용합니다.

<div className="accordion">
  <input type="checkbox" id="best-bid" />
  <label for="best-bid">
    <i className="fa-solid fa-circle-info"></i>
    최유리지정가 매수 주문 생성 요청 파라미터 사용 예시
  </label>

  <div className="accordion-content">

아래 표를 참고하여 최유리지정가 매수 주문 생성 요청 시 사용 가능한 파라미터를 쉽게 확인할 수 있습니다. <code>volume</code> 파라미터를 입력하지 않습니다(value에 null로 입력 또는 key 자체를 제외). 각 파라미터에 대한 상세한 설명은 하단 Request Body를 참고해주세요.

<table className="custom-table">
<thead>
  <tr>
    <th>파라미터</th>
    <th>필수 여부</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col">market</td>
    <td>Required</td>
    <td>페어 코드. <code>KRW-BTC</code> 형식으로 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">side</td>
    <td>Required</td>
    <td><code>bid</code></td>
  </tr>
  <tr>
    <td className="code-col">ord_type</td>
    <td>Required</td>
    <td><code>best</code></td>
  </tr>
  <tr>
    <td className="code-col">price</td>
    <td>Required</td>
    <td>호가 자산 기준 주문 총액. 최유리 호가로 주문 총액에 해당하는 수량을 매수하는 주문이 생성됩니다. 예를 들어, KRW-BTC 페어에서 100000000을 입력하는 경우 최유리 호가로 1억 원어치의 BTC 수량을 매수하는 주문이 생성됩니다.</td>
  </tr>
  <tr>
    <td className="code-col">time_in_force</td>
    <td><b>Required</b></td>
    <td><code>ioc</code>,<code>fok</code></td>
  </tr>
  <tr>
    <td className="code-col">smp_type</td>
    <td>Optional</td>
    <td>자전거래 체결 방지 옵션. <code>cancel_maker</code>,<code>cancel_taker</code>,<code>reduce</code></td>
  </tr>
  <tr>
    <td className="code-col">identifier</td>
    <td>Optional</td>
    <td>조회, 삭제시 사용할 수 있는 사용자 지정 주문 ID.</td>
  </tr>
</tbody>

</table>

  </div>
</div>

<div className="accordion">
  <input type="checkbox" id="best-ask" />
  <label for="best-ask">
    <i className="fa-solid fa-circle-info"></i>
    최유리지정가 매도 주문 생성 요청 파라미터 사용 예시
  </label>

  <div className="accordion-content">

아래 표를 참고하여 최유리지정가 매도 주문 생성 요청 시 사용 가능한 파라미터를 쉽게 확인할 수 있습니다. <code>price</code> 파라미터를 입력하지 않습니다(value에 null로 입력 또는 key 자체를 제외). 각 파라미터에 대한 상세한 설명은 하단 Request Body를 참고해주세요.

<table className="custom-table">
<thead>
  <tr>
    <th>파라미터</th>
    <th>필수 여부</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col">market</td>
    <td>Required</td>
    <td>페어 코드. <code>KRW-BTC</code> 형식으로 입력합니다.</td>
  </tr>
  <tr>
    <td className="code-col">side</td>
    <td>Required</td>
    <td><code>ask</code></td>
  </tr>
  <tr>
    <td className="code-col">ord_type</td>
    <td>Required</td>
    <td><code>best</code></td>
  </tr>
  <tr>
    <td className="code-col">volume</td>
    <td>Required</td>
    <td>매도 주문 수량. 예를 들어, KRW-BTC 페어에서 0.1을 입력하는 경우 최유리 호가로 0.1개의 BTC 수량을 매도하는 주문이 생성됩니다.</td>
  </tr>
  <tr>
    <td className="code-col">time_in_force</td>
    <td><b>Required</b></td>
    <td><code>ioc</code>,<code>fok</code></td>
  </tr>
  <tr>
    <td className="code-col">smp_type</td>
    <td>Optional</td>
    <td>자전거래 체결 방지 옵션. <code>cancel_maker</code>,<code>cancel_taker</code>,<code>reduce</code></td>
  </tr>
  <tr>
    <td className="code-col">identifier</td>
    <td>Optional</td>
    <td>조회, 삭제시 사용할 수 있는 사용자 지정 주문 ID.</td>
  </tr>
</tbody>

</table>

  </div>
</div>

<br />

### 주문 체결 조건(time\_in\_force)

주문 옵션으로, 주문 생성 시점의 체결 상황에 따른 주문 처리 방식을 지정할 수 있습니다.

<div className="accordion">
    <input type="checkbox" id="time_in_force" />
    <label for="time_in_force">
      <i className="fa-solid fa-circle-info"></i>
      사용 가능한 주문 체결 조건(time_in_force) 옵션
    </label>

<div className="accordion-content">








<table className="custom-table">
<thead>
  <tr>
    <th>옵션</th>
    <th>파라미터 값</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col"><b>IOC(Immediate or Cancel)</b></td>
    <td><code>ioc</code></td>
    <td>지정가 조건으로 즉시 체결 가능한 수량만 부분 체결하고, 잔여 수량은 취소합니다. <b>지정가 주문과 최유리 지정가 주문 에서 사용 가능</b>한 옵션입니다.</td>
  </tr>
  <tr>
    <td className="code-col"><b>FOK(Fill or Kill)</b></td>
    <td><code>fok</code></td>
    <td>지정가 조건으로 주문량 전량 체결 가능할 때만 주문을 실행하고, 아닌 경우 전량 주문 취소합니다. <b>지정가 주문과 최유리 지정가 주문 에서 사용 가능</b>한 옵션입니다.</td>
  </tr>
  <tr>
    <td className="code-col"><b>Post Only</b></td>
    <td><code>post_only</code></td>
    <td>
      지정가 조건으로 부분 또는 전체에 대해 즉시 체결 가능한 상황인 경우 주문을 실행하지 않고 취소합니다. 
      즉, 메이커(maker)주문으로 생성될 수 있는 상황에서만 주문이 생성되며 테이커(taker) 주문으로 체결되는 것을 방지합니다. 
      <b>지정가 주문(ord_type이 limit)에서만 사용 가능</b>한 옵션입니다. <b>자전 거래 체결 방지 옵션과 함께 사용할 수 없습니다.</b>
    </td>
  </tr>
</tbody>

</table>








  </div>
</div>

<br />

### 자전거래 체결 방지 옵션(SMP, Self-Matching Prevention)

`smp_type` 파라미터를 설정하여 자전거래 체결 방지 옵션을 원하는 모드로 활성화할 수 있습니다. 자전거래 체결 방지 기능과 관련된 자세한 사항은 [자전거래 체결 방지(Self-Match Prevention, SMP)](https://docs.upbit.com/kr/docs/smp)  페이지를 참고하시기 바랍니다.

<div className="accordion">
    <input type="checkbox" id="smp_type" />
    <label for="smp_type">
      <i className="fa-solid fa-circle-info"></i>
      사용 가능한 자전 거래 체결 방지(SMP) 옵션
    </label>

  <div className="accordion-content">
			<li>메이커(maker) 주문과 테이커(taker) 주문에 설정된 SMP 모드가 서로 상이한 경우 테이커 주문 모드에 따라 동작합니다.</li>
			<li>주문 생성 시 설정한 SMP 모드에 따라 기존 주문 또는 신규 주문의 전체 또는 부분 취소되는 경우 취소된 주문 수량과 금액은 주문 생성 응답의 "prevented_volume"필드와 "prevented_locked" 필드로 반환됩니다.</li>
			<br />









<table className="custom-table">
<thead>
  <tr>
    <th>옵션</th>
    <th>파라미터 값</th>
    <th>설명</th>
  </tr>
</thead>
<tbody>
  <tr>
    <td className="code-col"><b>메이커 주문 취소</b></td>
    <td><code>cancel_maker</code></td>
    <td>메이커 주문을 취소합니다. 즉, 새로운 주문 생성 시 자전 거래 조건이 성립하는 경우 이전에 생성한 주문을 취소하여 체결을 방지합니다.</td>
  </tr>
  <tr>
    <td className="code-col"><b>테이커 주문 취소</b></td>
    <td><code>cancel_taker</code></td>
    <td>테이커 주문을 취소합니다. 즉, 새로운 주문 생성 시 자전 거래 조건이 성립하는 경우 새롭게 생성한 주문을 취소하여 체결을 방지합니다.</td>
  </tr>
  <tr>
    <td className="code-col"><b>주문 수량 조정</b></td>
    <td><code>reduce</code></td>
    <td>
      새로운 주문 생성 시 자전 거래 조건이 성립하는 경우 기존 주문과 신규 주문의 주문 수량을 줄여 체결을 방지합니다. 잔량이 0인 경우 주문을 취소합니다.
    </td>
  </tr>
</tbody>

</table>









  </div>
</div>

<br />

### 체결 대기 중 자산 잠금

주문 생성시 해당 주문에 사용되는 호가 자산(매수 주문의 경우) 또는 기준 자산(매도 주문의 경우)이 즉시 잠금(locked) 상태로 전환되며, 다른 용도로 사용할 수 없게 됩니다. 이는 사용자의 잔고가 주문 체결 시점에도 유효하도록 보장하기 위한 동작이며 <Anchor target="_blank" href="ref:get-balance">포켓 잔고 조회</Anchor>API 를 호출하여 잠금 자산 현황을 확인할 수 있습니다. 자산 잠금은 아래 조건 중 하나가 충족될 때까지 유지됩니다.

\- 주문이 전량 체결되는 경우
\- 사용자 요청으로 주문이 취소되는 경우
\- `time_in_force` 조건에 따라 주문이 만료되는 경우

예시 ) KRW-BTC 마켓에서 지정가 매수 주문을 생성할 경우, 지정한 KRW 금액이 체결 전까지 잠금 상태로 유지됩니다.

<br />

### 주문 가격 단위와 최소 주문 가능 금액

마켓(호가 자산)과 기준 자산 단가에 따라 주문 시 사용 가능한 주문 가격 단위와 최소 주문 금액이 상이합니다. 마켓별 호가 정책은 아래 가이드를 참고하시기 바랍니다.

\- [원화(KRW) 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/krw-market-info)
\- [BTC 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/btc-market-info)
\- [USDT 마켓 주문 가격 단위 / 최소 주문 가능 금액](https://docs.upbit.com/kr/docs/usdt-market-info)

<br />

### 클라이언트 주문 식별자(identifier)

주문 생성 시 업비트 시스템에서 주문을 고유하게 식별하기 위한 UUID와는 별도로, 사용자 클라이언트에서 각 주문을 식별하기 위해 직접 지정하는 값입니다. 사용자가 자체적인 주문 ID 체계를 통해 주문을 관리(조회 및 취소)하고자 할 때 활용됩니다.

identifier는 사용자 계정 내 전체 주문 기준으로 고유해야 하며, 이미 사용된 값은 해당 주문의 생성 또는 체결 여부와 관계없이 재사용할 수 없습니다. 또한 identifier의 최대 길이는 64자로 제한됩니다.

<br />

<div className="callout-section callout-section--danger">
    <div className="callout-title">
      <i className="fa-solid fa-circle-exclamation"></i> POST API에 대한 Form 방식 요청은 2022년 3월 1일부로 지원이 종료되었습니다.
      </div>
    Form 방식 지원 종료에 따라 Urlencoded Form 방식으로 전송하는 POST 요청에 대한 정상적인 동작을 보장하지 않습니다. <b>반드시 JSON 형식으로 요청 본문(Body)을 전송</b>해주시기 바랍니다.
  </div>

<div className="accordion-changelog">
    <input type="checkbox" id="api-changelog" />
    <label for="api-changelog">
        <div className="APISectionHeader-heading4MUMLbp4_nLs">Revision History <i className="fa-solid fa-angle-right"></i> </div>
    </label>

    <div className="accordion-changelog-content">
        













<table className="custom-table">
    <thead>
        <tr>
            <th>반영 버전</th>
            <th>반영 일자</th>
            <th>변경 사항</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td className="code-col">v1.5.8</td>
            <td>2025-07-07</td>
            <td><a href="https://docs.upbit.com/kr/changelog/post_only"><code>Post Only</code> 주문 옵션 신규 지원</a></td>
      	</tr>
<tr>
            <td className="code-col">v1.5.8</td>
            <td>2025-07-02</td>
            <td><a href="https://docs.upbit.com/kr/changelog/smp"> 자전거래 체결 방지(SMP) 기능<br />신규 지원에 따른 파라미터 추가<br /><code>smp_type</code>,<code>prevented_volume</code>,<code>prevented_locked</code></a></td>
        </tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-12-04</td>
            <td><a href="https://docs.upbit.com/kr/changelog/myorder_identifier"> <code>identifier</code> 필드 신규 지원</a></td>
     	 	</tr>
<tr>
            <td className="code-col">-</td>
            <td>2024-04-22</td>
            <td><a href="https://docs.upbit.com/kr/changelog/new_ord_type_expand"> 최유리지정가 주문 유형 신규 지원<br />주문 옵션(time_in_force) 추가 지원</a></td>
        </tr>

    </tbody>
</table>














    </div>
</div>

<div className="APISectionHeader-heading4MUMLbp4_nLs">Rate Limit</div>

<div className="box-rate-limit">

초당 최대 12회 호출할 수 있습니다. 포켓 단위로 측정되며 \[주문 생성 그룹] 내에서 요청 가능 횟수를 공유합니다. Rate Limit은 요청 처리량을 보장하는 기준이 아니며, 트래픽 상황이나 서비스 안정성 확보 필요에 따라 제한 또는 조정될 수 있습니다.

</div>
  <br />
  <div className="APISectionHeader-heading4MUMLbp4_nLs">API Key Permission</div>
  <div className="box-rate-limit">
    <a href="auth">인증</a>이 필요한 API로, [주문하기] 권한이 설정된 API Key를 사용해야 합니다. <br />
    권한 오류(out_of_scope) 오류가 발생한다면, <a href="https://upbit.com/mypage/open_api_management">API Key 관리 메뉴</a>에서 권한 설정을 확인해주세요.
  </div>

# OpenAPI definition

```json
{
  "openapi": "3.0.2",
  "info": {
    "title": "EXCHANGE API",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "https://api.upbit.com"
    }
  ],
  "components": {
    "schemas": {
      "SideEnum": {
        "type": "string",
        "enum": [
          "ask",
          "bid"
        ],
        "description": "주문 방향(매수/매도)"
      },
      "OrdTypeEnum": {
        "type": "string",
        "enum": [
          "limit",
          "price",
          "market",
          "best"
        ],
        "description": "주문 유형.\n\n* `limit`: 지정가 주문\n* `price`: 시장가 주문(매수)\n* `market`: 시장가 주문(매도)\n* `best`: 최유리 주문\n"
      },
      "TimeInForceEnum": {
        "type": "string",
        "enum": [
          "fok",
          "ioc",
          "post_only"
        ],
        "description": "주문 체결 조건.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: 메이커 전용 주문\n"
      },
      "SmpTypeEnum": {
        "type": "string",
        "enum": [
          "cancel_maker",
          "cancel_taker",
          "reduce"
        ],
        "description": "자전거래 체결 방지(SMP) 모드.\n\n* `cancel_maker`: 메이커 주문 취소\n* `cancel_taker`: 테이커 주문 취소\n* `reduce`: 수량 감소\n"
      },
      "Order": {
        "type": "object",
        "required": [
          "market",
          "uuid",
          "side",
          "ord_type",
          "state",
          "created_at",
          "remaining_volume",
          "executed_volume",
          "reserved_fee",
          "remaining_fee",
          "paid_fee",
          "locked",
          "trades_count",
          "prevented_volume",
          "prevented_locked"
        ],
        "properties": {
          "market": {
            "type": "string",
            "description": "페어(거래쌍)의 코드\n",
            "example": "KRW-BTC"
          },
          "uuid": {
            "type": "string",
            "description": "주문의 유일 식별자",
            "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
          },
          "side": {
            "type": "string",
            "enum": [
              "ask",
              "bid"
            ],
            "description": "주문 방향(매수/매도)",
            "example": "ask"
          },
          "ord_type": {
            "type": "string",
            "enum": [
              "limit",
              "price",
              "market",
              "best"
            ],
            "description": "주문 유형.",
            "example": "limit"
          },
          "price": {
            "type": "string",
            "format": "decimal",
            "description": "주문 단가 또는 총액\n\n지정가 주문의 경우 단가, 시장가 매수 주문의 경우 매수 총액입니다.\n",
            "example": 1000
          },
          "state": {
            "type": "string",
            "enum": [
              "wait",
              "watch",
              "done",
              "cancel"
            ],
            "description": "주문 상태\n\n- `wait`: 체결 대기\n- `watch`: 예약 주문 대기\n- `done`: 체결 완료\n- `cancel`: 주문 취소\n",
            "example": "wait"
          },
          "created_at": {
            "type": "string",
            "description": "주문 생성 시각 (KST 기준)\n\n[형식] yyyy-MM-ddTHH:mm:ss+09:00\n",
            "example": "2025-06-25T15:42:25+09:00"
          },
          "volume": {
            "type": "string",
            "format": "decimal",
            "description": "주문 요청 수량",
            "example": 10
          },
          "remaining_volume": {
            "type": "string",
            "format": "decimal",
            "description": "체결 후 남은 주문 양",
            "example": 8
          },
          "executed_volume": {
            "type": "string",
            "format": "decimal",
            "description": "체결된 양",
            "example": 2
          },
          "reserved_fee": {
            "type": "string",
            "format": "decimal",
            "description": "수수료로 예약된 비용",
            "example": 5
          },
          "remaining_fee": {
            "type": "string",
            "format": "decimal",
            "description": "남은 수수료",
            "example": 5
          },
          "paid_fee": {
            "type": "string",
            "format": "decimal",
            "description": "사용된 수수료",
            "example": 0
          },
          "locked": {
            "type": "string",
            "format": "decimal",
            "description": "거래에 사용 중인 비용",
            "example": 0
          },
          "trades_count": {
            "type": "integer",
            "description": "해당 주문에 대한 체결 건수",
            "example": 1
          },
          "time_in_force": {
            "type": "string",
            "enum": [
              "fok",
              "ioc",
              "post_only"
            ],
            "description": "주문 체결 옵션",
            "example": "ioc"
          },
          "identifier": {
            "type": "string",
            "description": "주문 생성시 클라이언트가 지정한 주문 식별자.\n* identifier 필드는 2024년 10월 18일 이후 생성된 주문에 대해서만 제공됩니다.\n",
            "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
          },
          "smp_type": {
            "type": "string",
            "enum": [
              "reduce",
              "cancel_maker",
              "cancel_taker"
            ],
            "description": "자전거래 체결 방지(Self-Match Prevention) 모드",
            "example": "cancel_maker"
          },
          "prevented_volume": {
            "type": "string",
            "format": "decimal",
            "description": "자전거래 방지로 인해 취소된 수량.\n\n동일 사용자의 주문 간 체결이 발생하지 않도록 설정(SMP)에 따라 취소된 주문 수량입니다.\n",
            "example": 2
          },
          "prevented_locked": {
            "type": "string",
            "format": "decimal",
            "description": "자전거래 방지로 인해 해제된 자산.\n자전거래 체결 방지 설정으로 인해 취소된 주문의 잔여 자산입니다.\n  - 매수 주문의 경우: 취소된 금액\n  - 매도 주문의 경우: 취소된 수량\n",
            "example": 2000
          },
          "trades": {
            "type": "array",
            "description": "체결 목록 (개별 주문 조회 시 반환)",
            "items": {
              "type": "object",
              "properties": {
                "market": {
                  "type": "string",
                  "description": "페어(거래쌍)의 코드",
                  "example": "KRW-BTC"
                },
                "uuid": {
                  "type": "string",
                  "description": "체결의 유일 식별자",
                  "example": "795dff29-bba6-49b2-baab-63473ab7931c"
                },
                "price": {
                  "type": "string",
                  "format": "decimal",
                  "description": "체결 단가",
                  "example": "1375"
                },
                "volume": {
                  "type": "string",
                  "format": "decimal",
                  "description": "체결 수량",
                  "example": "5.377594"
                },
                "funds": {
                  "type": "string",
                  "format": "decimal",
                  "description": "체결 금액",
                  "example": "7394.19175"
                },
                "trend": {
                  "type": "string",
                  "description": "체결 방향 (up/down)",
                  "example": "down"
                },
                "created_at": {
                  "type": "string",
                  "description": "체결 시각",
                  "example": "2025-08-09T16:44:00.597751+09:00"
                },
                "side": {
                  "type": "string",
                  "description": "주문 방향",
                  "example": "ask"
                }
              }
            }
          }
        }
      },
      "Error": {
        "type": "object",
        "properties": {
          "error": {
            "type": "object",
            "required": [
              "name",
              "message"
            ],
            "properties": {
              "name": {
                "type": "string",
                "description": "에러명"
              },
              "message": {
                "type": "string",
                "description": "에러 메세지"
              }
            }
          }
        }
      }
    }
  },
  "x-readme": {
    "explorer-enabled": false,
    "samples-languages": [
      "shell",
      "python",
      "java",
      "node"
    ]
  },
  "paths": {
    "/v1/orders": {
      "post": {
        "operationId": "new-order",
        "summary": "주문 생성",
        "description": "특정 페어를 매수/매도하기 위한 주문을 생성합니다.",
        "tags": [
          "Exchange"
        ],
        "x-readme": {
          "code-samples": [
            {
              "language": "curl",
              "code": "curl --request POST \\\n--url 'https://api.upbit.com/v1/orders' \\\n--header 'Authorization: Bearer {JWT_TOKEN}' \\\n--header 'Accept: application/json' \\\n--header 'Content-Type: application/json' \\\n--data '\n  {\n  \"market\":\"KRW-BTC\",\n  \"side\":\"bid\",\n  \"volume\":\"1\",\n  \"price\":\"140000000\",\n  \"ord_type\":\"limit\"\n  }\n'\n"
            },
            {
              "language": "python",
              "install": "pip install requests pyjwt python-dotenv",
              "code": "import os\nimport uuid\nimport hashlib\nimport jwt\nimport requests\nfrom urllib.parse import unquote, urlencode\nfrom dotenv import load_dotenv\n\nload_dotenv()\n\nBASE_URL = \"https://api.upbit.com\"\nPATH = \"/v1/orders\"\n\nACCESS_KEY = os.environ[\"UPBIT_OPEN_API_ACCESS_KEY\"]\nSECRET_KEY = os.environ[\"UPBIT_OPEN_API_SECRET_KEY\"]\n\nparams = {\n    \"market\": \"KRW-BTC\",\n    \"side\": \"bid\",\n    \"volume\": \"1\",\n    \"price\": \"140000000\",\n    \"ord_type\": \"limit\",\n}\nquery_string = unquote(urlencode(params, doseq=True)).encode(\"utf-8\")\n\nm = hashlib.sha512()\nm.update(query_string)\nquery_hash = m.hexdigest()\n\npayload = {\n    \"access_key\": ACCESS_KEY,\n    \"nonce\": str(uuid.uuid4()),\n    \"query_hash\": query_hash,\n    \"query_hash_alg\": \"SHA512\",\n}\n\njwt_token = jwt.encode(payload, SECRET_KEY, algorithm=\"HS256\")\n\nheaders = {\n    \"Authorization\": f\"Bearer {jwt_token}\",\n    \"Accept\": \"application/json\",\n}\n\nres = requests.post(f\"{BASE_URL}{PATH}\", headers=headers, json=params)\nprint(res.json())\n"
            },
            {
              "language": "node",
              "name": "Axios",
              "install": "npm install axios jsonwebtoken uuid",
              "code": "const axios = require(\"axios\");\nconst crypto = require(\"crypto\");\nconst { sign } = require(\"jsonwebtoken\");\nconst { v4: uuidv4 } = require(\"uuid\");\nrequire(\"dotenv\").config();\n\nconst baseURL = \"https://api.upbit.com\";\nconst path = \"/v1/orders\";\n\nconst ACCESS_KEY = process.env.UPBIT_OPEN_API_ACCESS_KEY;\nconst SECRET_KEY = process.env.UPBIT_OPEN_API_SECRET_KEY;\n\nconst params = {\n  market: \"KRW-BTC\",\n  side: \"bid\",\n  volume: \"1\",\n  price: \"140000000\",\n  ord_type: \"limit\",\n};\n\nconst queryString = new URLSearchParams(params).toString();\n\nconst queryHash = crypto\n  .createHash(\"sha512\")\n  .update(queryString, \"utf-8\")\n  .digest(\"hex\");\n\nconst payload = {\n  access_key: ACCESS_KEY,\n  nonce: uuidv4(),\n  query_hash: queryHash,\n  query_hash_alg: \"SHA512\",\n};\n\nconst jwtToken = sign(payload, SECRET_KEY);\n\nconst options = {\n  method: \"POST\",\n  url: `${baseURL}${path}`,\n  headers: {\n    Authorization: `Bearer ${jwtToken}`,\n    Accept: \"application/json\",\n  },\n  data: params,\n};\n\naxios\n  .request(options)\n  .then((response) => {\n    console.log(response.data);\n  })\n  .catch((error) => {\n    console.error(error.response ? error.response.data : error.message);\n  });\n"
            },
            {
              "language": "java",
              "code": "package main;\n\nimport com.auth0.jwt.JWT;\nimport com.auth0.jwt.algorithms.Algorithm;\nimport java.io.IOException;\nimport java.nio.charset.StandardCharsets;\nimport java.security.MessageDigest;\nimport java.security.NoSuchAlgorithmException;\nimport java.util.HashMap;\nimport java.util.Map;\nimport java.util.Objects;\nimport java.util.UUID;\nimport java.util.stream.Collectors;\nimport okhttp3.OkHttpClient;\nimport okhttp3.Request;\nimport okhttp3.RequestBody;\nimport okhttp3.Response;\nimport com.google.gson.Gson;\n\npublic class CreateOrder {\n    private static final String BASE_URL = \"https://api.upbit.com\";\n    private static final String PATH = \"/v1/orders\";\n\n\n    public static void main(String[] args) throws NoSuchAlgorithmException, IOException {\n        String accessKey = System.getenv(\"UPBIT_OPEN_API_ACCESS_KEY\");\n        String secretKey = System.getenv(\"UPBIT_OPEN_API_SECRET_KEY\");\n\n        Map<String, String> params = new HashMap<>();\n        params.put(\"market\", \"KRW-BTC\");\n        params.put(\"side\", \"bid\");\n        params.put(\"volume\", \"1\");\n        params.put(\"price\", \"140000000\");\n        params.put(\"ord_type\", \"limit\");\n        String queryString = params.entrySet().stream()\n            .map(e -> e.getKey() + \"=\" + String.valueOf(e.getValue()))\n            .collect(Collectors.joining(\"&\"));\n\n        MessageDigest md = MessageDigest.getInstance(\"SHA-512\");\n        md.update(queryString.getBytes(StandardCharsets.UTF_8));\n        StringBuilder sb = new StringBuilder();\n        for (byte b : md.digest()) {\n            sb.append(String.format(\"%02x\", b));\n        }\n        String queryHash = sb.toString();\n\n        Algorithm algorithm = Algorithm.HMAC512(secretKey.getBytes(StandardCharsets.UTF_8));\n        String jwtToken = JWT.create()\n            .withClaim(\"access_key\", accessKey)\n            .withClaim(\"nonce\", UUID.randomUUID().toString())\n            .withClaim(\"query_hash\", queryHash)\n            .withClaim(\"query_hash_alg\", \"SHA512\")\n            .sign(algorithm);\n\n        String authHeader = \"Bearer \" + jwtToken;\n\n        String jsonBody = new Gson().toJson(params);\n        OkHttpClient client = new OkHttpClient();\n        Request request = new Request.Builder()\n            .url(BASE_URL + PATH)\n            .post(RequestBody.create(jsonBody, okhttp3.MediaType.parse(\"application/json; charset=utf-8\")))\n            .addHeader(\"Content-Type\", \"application/json\")\n            .addHeader(\"Authorization\", authHeader)\n            .build();\n\n        try (Response response = client.newCall(request).execute()) {\n            System.out.println(response.code());\n            System.out.println(Objects.requireNonNull(response.body()).string());\n        }\n    }\n}\n"
            }
          ]
        },
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "examples": {
                "default": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "bid",
                    "volume": "1",
                    "price": "14000000",
                    "ord_type": "limit"
                  }
                },
                "market_sell": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "ask",
                    "volume": "0.001",
                    "ord_type": "market"
                  }
                },
                "market_buy": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "bid",
                    "price": "10000",
                    "ord_type": "price"
                  }
                },
                "best_buy": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "bid",
                    "price": "10000",
                    "ord_type": "best",
                    "time_in_force": "ioc"
                  }
                },
                "best_sell": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "ask",
                    "volume": "0.001",
                    "ord_type": "best",
                    "time_in_force": "ioc"
                  }
                },
                "limit_ioc": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "bid",
                    "volume": "0.001",
                    "price": "10000000",
                    "ord_type": "limit",
                    "time_in_force": "ioc"
                  }
                },
                "limit_smp_reduce": {
                  "value": {
                    "market": "KRW-BTC",
                    "side": "bid",
                    "volume": "0.001",
                    "price": "10000000",
                    "ord_type": "limit",
                    "smp_type": "reduce"
                  }
                }
              },
              "schema": {
                "type": "object",
                "required": [
                  "market",
                  "side",
                  "ord_type"
                ],
                "properties": {
                  "market": {
                    "type": "string",
                    "description": "주문을 생성하고자 하는 대상 페어(거래쌍). (필수)",
                    "example": "KRW-BTC"
                  },
                  "side": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/SideEnum"
                      }
                    ],
                    "description": "주문 방향(매수/매도). (필수)"
                  },
                  "volume": {
                    "type": "string",
                    "format": "decimal",
                    "description": "주문 수량",
                    "example": "0.01"
                  },
                  "price": {
                    "type": "string",
                    "format": "decimal",
                    "description": "주문 단가 또는 총액",
                    "example": "1000"
                  },
                  "ord_type": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/OrdTypeEnum"
                      }
                    ],
                    "description": "주문 유형. (필수)\n\n* `limit`: 지정가 주문\n* `price`: 시장가 주문(매수)\n* `market`: 시장가 주문(매도)\n* `best`: 최유리 주문\n"
                  },
                  "identifier": {
                    "type": "string",
                    "description": "클라이언트 지정 주문 식별자",
                    "example": "9ca023a5-851b-4fec-9f0a-48cd83c2eaae"
                  },
                  "time_in_force": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/TimeInForceEnum"
                      }
                    ],
                    "description": "주문 체결 조건.\n\n* `fok`: Fill or Kill\n* `ioc`: Immediate or Cancel\n* `post_only`: 메이커 전용 주문\n"
                  },
                  "smp_type": {
                    "allOf": [
                      {
                        "$ref": "#/components/schemas/SmpTypeEnum"
                      }
                    ],
                    "description": "자전거래 체결 방지(SMP) 모드.\n\n* `cancel_maker`: 메이커 주문 취소\n* `cancel_taker`: 테이커 주문 취소\n* `reduce`: 수량 감소\n"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Object of created order",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Order"
                },
                "examples": {
                  "Successful Example": {
                    "value": {
                      "market": "KRW-BTC",
                      "uuid": "cdd92199-2897-4e14-9448-f923320408ad",
                      "side": "ask",
                      "ord_type": "limit",
                      "price": "140000000",
                      "state": "wait",
                      "created_at": "2025-07-04T15:00:00+09:00",
                      "volume": "1.0",
                      "remaining_volume": "1.0",
                      "reserved_fee": "70000.0",
                      "remaining_fee": "70000.0",
                      "paid_fee": "0.0",
                      "locked": "0.0",
                      "executed_volume": "0.0",
                      "prevented_volume": "0",
                      "prevented_locked": "0",
                      "trades_count": 0
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "invalid parameter error": {
                    "value": {
                      "error": {
                        "name": "invalid_parameter",
                        "message": "잘못된 파라미터"
                      }
                    }
                  },
                  "invalid time in force error": {
                    "value": {
                      "error": {
                        "name": "invalid_time_in_force",
                        "message": "주문조건을 다시 확인해 주세요."
                      }
                    }
                  },
                  "over krw funds bid error": {
                    "value": {
                      "error": {
                        "name": "over_krw_funds_bid",
                        "message": "최대 주문 가능 금액 1000000000 KRW 보다 작은 주문을 입력해 주세요."
                      }
                    }
                  },
                  "insufficient funds bid error": {
                    "value": {
                      "error": {
                        "name": "insufficient_funds_bid",
                        "message": "주문 가능한 금액(KRW)이 부족합니다."
                      }
                    }
                  },
                  "not found market error": {
                    "value": {
                      "error": {
                        "name": "notfoundmarket",
                        "message": "다음 마켓에 없는 종목입니다: KRW-BTCs"
                      }
                    }
                  }
                }
              }
            }
          },
          "403": {
            "description": "error object",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Error"
                },
                "examples": {
                  "market offline error": {
                    "value": {
                      "error": {
                        "name": "market_offline",
                        "message": "시스템 점검 중입니다. 점검이 완료된 후에 다시 시도해주세요."
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

# Sibling pages

* [페어별 주문 가능 정보 조회](https://docs.upbit.com/kr/reference/available-order-information.md)
* [주문 생성 테스트](https://docs.upbit.com/kr/reference/order-test.md)
* [개별 주문 조회](https://docs.upbit.com/kr/reference/get-order.md)
* [id로 주문 목록 조회](https://docs.upbit.com/kr/reference/list-orders-by-ids.md)
* [체결 대기 주문 목록 조회](https://docs.upbit.com/kr/reference/list-open-orders.md)
* [종료 주문 목록 조회](https://docs.upbit.com/kr/reference/list-closed-orders.md)
* [개별 주문 취소 접수](https://docs.upbit.com/kr/reference/cancel-order.md)
* [id로 주문 목록 취소 접수](https://docs.upbit.com/kr/reference/cancel-orders-by-ids.md)
* [주문 일괄 취소 접수](https://docs.upbit.com/kr/reference/batch-cancel-orders.md)
* [취소 후 재주문](https://docs.upbit.com/kr/reference/cancel-and-new-order.md)