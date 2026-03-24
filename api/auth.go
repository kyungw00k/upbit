package api

import (
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken JWT HS512 인증 토큰 생성
// query: GET/DELETE 쿼리 파라미터 또는 POST 본문에서 추출된 key=value 맵
func GenerateToken(accessKey, secretKey string, query map[string]string) (string, error) {
	claims := jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	}

	if len(query) > 0 {
		// Upbit 인증 문서: URL 인코딩 되지 않은 쿼리 문자열을 기준으로 Hash 생성
		queryString := buildQueryHashString(query)
		hash := sha512.Sum512([]byte(queryString))
		claims["query_hash"] = hex.EncodeToString(hash[:])
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

// GenerateTokenFromBody POST 본문(JSON 구조체의 필드)으로부터 토큰 생성
// Upbit POST 요청 시 body 파라미터를 query string 형태로 해시해야 함
func GenerateTokenFromBody(accessKey, secretKey string, bodyParams map[string]string) (string, error) {
	return GenerateToken(accessKey, secretKey, bodyParams)
}

// GenerateTokenFromRawQuery 배열 쿼리 등 raw 쿼리 문자열로부터 토큰 생성
// map[string]string으로 표현할 수 없는 배열 파라미터 (uuids[]=a&uuids[]=b) 지원
func GenerateTokenFromRawQuery(accessKey, secretKey string, rawQuery string) (string, error) {
	claims := jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	}

	if rawQuery != "" {
		hash := sha512.Sum512([]byte(rawQuery))
		claims["query_hash"] = hex.EncodeToString(hash[:])
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

// buildQueryHashString 해시 생성용 비인코딩 쿼리 문자열 생성
// Upbit 인증 문서: "URL 인코딩 되지 않은 쿼리 문자열을 기준으로 Hash 값을 생성해야 합니다."
// 결정적 동작을 위해 키를 정렬함 (map은 순서가 없으므로)
func buildQueryHashString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(params))
	for _, k := range keys {
		// URL 인코딩 없이 그대로 연결 (해시용)
		pairs = append(pairs, k+"="+params[k])
	}
	return strings.Join(pairs, "&")
}

// BuildQueryString 쿼리 맵을 URL 인코딩된 쿼리 문자열로 변환
// HTTP 요청 URL 구성 시 사용 (URL 인코딩 적용)
func BuildQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(params))
	for _, k := range keys {
		pairs = append(pairs, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	return strings.Join(pairs, "&")
}
