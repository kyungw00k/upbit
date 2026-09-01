Fetch the complete documentation index at: https://docs.upbit.com/kr/llms.txt. Use this file to discover all available pages before exploring further. Append .md to any documentation page URL to get its markdown version.

# [안내]업비트 Open API Rate Limit 상향 안내

안녕하세요.  업비트 개발자센터입니다.

업비트 회원님들을 위해 업비트 Open API Rate Limit을 상향합니다.

기존 초당 8회였던 Open API Rate Limit을 초당 12회로 상향하여, 더욱 여유로운 환경에서 Open API를 이용하실 수 있도록 지원합니다.

## 적용 일시

* 2026-08-21(금) 21:00

## 변경 사항

* 대상 그룹 : [order](https://docs.upbit.com/kr/reference/rate-limits#rate-limit-%EA%B7%B8%EB%A3%B9%EB%B3%84-%EC%A0%95%EC%B1%85:~:text=%EB%B3%B4%EB%82%BC%20%EC%88%98%20%EC%9E%88%EC%8A%B5%EB%8B%88%EB%8B%A4.-,Rate%20Limit%20%EA%B7%B8%EB%A3%B9%EB%B3%84%20%EC%A0%95%EC%B1%85,-API%EB%8A%94%20Rate)
* 기존 : 초당 최대 8회
* 변경 : 초당 최대 12회

## 유의사항

* **서비스 안정성 및 운영 상황에 따라 Open API Rate Limit은 변경될 수 있으며, 변경 시 별도 공지사항을 통해 안내드릴 예정입니다.**
* **API Key 및 인증 정보는 외부에 노출되지 않도록 환경 변수 등을 활용하여 안전하게 관리해 주시기 바랍니다.**
* 발급된 API Key는 등록된 IP에서만 사용 가능하며, 발급 한도는 포켓별로 메인포켓 최대 10개, 서브포켓 최대 5개입니다.
* API Key는 1년간 유효하며 연장할 수 없습니다. 만료되거나 기능 변경이 필요할 경우 기존 Key를 삭제한 뒤 재발급해야 합니다.
* 과도한 요청 또는 비정상적인 트래픽 발생 시 서비스 이용이 제한될 수 있으므로, 업비트 Open API의 호출 제한(Rate Limit) 정책을 준수해 주시기 바랍니다.
* 각 API의 상세 측정 단위와 한도는 요청 수 제한(Rate Limits) 문서를 확인해 주시기 바랍니다.
* 분 단위 요청 제한은 더 이상 적용되지 않으며,<Anchor target="_blank" href="https://docs.upbit.com/kr/reference/rate-limits#%EC%9E%94%EC%97%AC-%EC%9A%94%EC%B2%AD-%EC%88%98-%ED%99%95%EC%9D%B8-%EB%B0%A9%EB%B2%95"> Remaining-Req 헤더의 min 값</Anchor>은 하위 호환성을 위해 유지되는 값으로 실제 Rate Limit 계산에는 사용되지 않습니다.

상향된 Open API Rate Limit을 통해 보다 원활하게 Open API를 이용하실 수 있으며,
앞으로도 더욱 편리하고 쾌적한 거래 환경 및 다양한 기능을 제공하기 위해 최선을 다하는 업비트가 되겠습니다.

감사합니다.