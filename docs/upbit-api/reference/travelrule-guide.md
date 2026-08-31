---
updatedAt: 2026-07-28T12:53:26.000Z
---

Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# 트래블룰 검증

타 거래소로부터 업비트로 입금시 트래블룰 준수를 위해 계정주 확인을 수행합니다.

### 계정주 확인 서비스란?

트래블룰(Travel Rule) 시행 의무에 따라, 가상자산사업자(VASP)간 입출금시 입금 계정주(업비트 계정)와 출금 계정주(해외 거래소 계정)의 동일 여부를 반드시 확인해야 합니다. 업비트는 VerifyVASP(베리파이바스프) 프로토콜을 기반으로 트래블룰 준수를 위한 계정주 확인 서비스를 제공하고 있습니다. VerifyVASP 프로토콜에 따라 업비트 입금건은 양쪽 계정주의 성명(한글 또는 영문) 및 생년월일 정보의 일치 여부를 확인한 뒤 정상 반영됩니다. 계정주 확인 서비스와 관련된 자세한 사항은 <Anchor target="_blank" href="https://support.upbit.com/hc/ko/articles/5127318107033-%EA%B3%84%EC%A0%95%EC%A3%BC-%ED%99%95%EC%9D%B8-%EC%84%9C%EB%B9%84%EC%8A%A4%EA%B0%80-%EB%AC%B4%EC%97%87%EC%9D%B8%EA%B0%80%EC%9A%94">업비트 고객센터 > 자주 묻는 질문 > 트래블룰</Anchor> 문서를 확인해주세요.

### 계정주 확인 서비스 실행이 필요한 경우

입금 상태 조회 시 상태가 `TRAVEL_RULE_SUSPECTED`로 반환되는 경우, 해당 입금건에 대한 트래블룰 검증이 수행되지 않아 입금 반영이 지연되고 있는 상태입니다. 해당 입금 UUID 또는 TxID로 트래블룰 검증 요청 API를 호출하여 계정주 확인을 완료하여 입금 반영을 재개할 수 있습니다.

### Open API를 통한 계정주 확인 방법

Open API를 통해 해외 거래소로부터 출금한 자산의 업비트 입금 요청 이후 계정주 확인을 자동화 할 수 있습니다. 계정주 확인 자동화는 다음과 같은 순서로 진행합니다.

1. <Anchor target="_blank" href="ref:list-travelrule-vasps">트래블룰 지원 거래소 목록 조회</Anchor> API를 호출하여 지원 거래소 목록으로부터 출금을 실행하는 해외 거래소의 UUID를 확인합니다.
2. 입금 UUID를 알고 있는 경우, <Anchor target="_blank" href="ref:verify-travelrule-by-uuid">입금 UUID로 트래블룰 검증 요청</Anchor> API를 호출하여 계정주 확인을 진행할 수 있습니다. 입금 UUID와 거래소 UUID를 파라미터로 입력합니다.
3. 트랜잭션 ID(TxID)를 알고 있는 경우 <Anchor target="_blank" href="ref:verify-travelrule-by-txid">입금 TxID로 트래블룰 검증 요청</Anchor> API를 호출하여 계정주 확인을 진행할 수 있습니다. 트랜잭션 ID(TxID)와 거래소 UUID를 파라미터로 입력합니다.

<br />

# Sub pages

* [계정주 확인 서비스 지원 거래소 목록 조회](https://docs.upbit.com/kr/reference/list-travelrule-vasps.md)
* [입금 UUID로 계정주 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-uuid.md)
* [입금 TxID로 계정주 검증 요청](https://docs.upbit.com/kr/reference/verify-travelrule-by-txid.md)

# Sibling pages

* [디지털 자산 입금 가능 정보 조회](https://docs.upbit.com/kr/reference/available-deposit-information.md)
* [입금 주소 생성 요청](https://docs.upbit.com/kr/reference/create-deposit-address.md)
* [개별 입금 주소 조회](https://docs.upbit.com/kr/reference/get-deposit-address.md)
* [입금 주소 목록 조회](https://docs.upbit.com/kr/reference/list-deposit-addresses.md)
* [원화 입금](https://docs.upbit.com/kr/reference/deposit-krw.md)
* [개별 입금 조회](https://docs.upbit.com/kr/reference/get-deposit.md)
* [입금 목록 조회](https://docs.upbit.com/kr/reference/list-deposits.md)