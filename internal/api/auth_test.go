package api

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_EmptyQuery_NoQueryHash(t *testing.T) {
	token, err := GenerateToken("test-access", "test-secret", nil)
	if err != nil {
		t.Fatalf("GenerateToken 실패: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("토큰 파싱 실패: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims 타입 변환 실패")
	}

	if claims["access_key"] != "test-access" {
		t.Errorf("access_key 기대: test-access, 실제: %v", claims["access_key"])
	}

	if _, exists := claims["query_hash"]; exists {
		t.Error("빈 query일 때 query_hash가 존재하면 안 됨")
	}

	if _, exists := claims["query_hash_alg"]; exists {
		t.Error("빈 query일 때 query_hash_alg가 존재하면 안 됨")
	}

	if _, exists := claims["nonce"]; !exists {
		t.Error("nonce가 존재해야 함")
	}
}

func TestGenerateToken_WithQuery_HasQueryHash(t *testing.T) {
	query := map[string]string{
		"market": "KRW-BTC",
		"side":   "bid",
	}

	token, err := GenerateToken("test-access", "test-secret", query)
	if err != nil {
		t.Fatalf("GenerateToken 실패: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("토큰 파싱 실패: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims 타입 변환 실패")
	}

	if _, exists := claims["query_hash"]; !exists {
		t.Error("query가 있으면 query_hash가 포함되어야 함")
	}

	if claims["query_hash_alg"] != "SHA512" {
		t.Errorf("query_hash_alg 기대: SHA512, 실제: %v", claims["query_hash_alg"])
	}
}

func TestBuildQueryHashString_SortedKeys_NoURLEncoding(t *testing.T) {
	params := map[string]string{
		"market":   "KRW-BTC",
		"side":     "bid",
		"ord_type": "limit",
	}

	result := buildQueryHashString(params)

	// 키가 정렬되어야 함 (market < ord_type < side)
	expected := "market=KRW-BTC&ord_type=limit&side=bid"
	if result != expected {
		t.Errorf("기대: %s, 실제: %s", expected, result)
	}

	// URL 인코딩이 적용되지 않아야 함 (- 기호가 그대로)
	if strings.Contains(result, "%") {
		t.Error("URL 인코딩이 적용되면 안 됨")
	}
}

func TestBuildQueryHashString_Empty(t *testing.T) {
	result := buildQueryHashString(nil)
	if result != "" {
		t.Errorf("빈 맵일 때 빈 문자열 기대, 실제: %s", result)
	}
}

func TestGenerateTokenFromRawQuery_HashesRawQuery(t *testing.T) {
	rawQuery := "uuids[]=abc&uuids[]=def"

	token, err := GenerateTokenFromRawQuery("test-access", "test-secret", rawQuery)
	if err != nil {
		t.Fatalf("GenerateTokenFromRawQuery 실패: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("토큰 파싱 실패: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims 타입 변환 실패")
	}

	// 해시 검증
	expectedHash := sha512.Sum512([]byte(rawQuery))
	expectedHex := hex.EncodeToString(expectedHash[:])

	if claims["query_hash"] != expectedHex {
		t.Errorf("query_hash 불일치.\n기대: %s\n실제: %v", expectedHex, claims["query_hash"])
	}
}

func TestGenerateTokenFromRawQuery_EmptyQuery(t *testing.T) {
	token, err := GenerateTokenFromRawQuery("test-access", "test-secret", "")
	if err != nil {
		t.Fatalf("GenerateTokenFromRawQuery 실패: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("토큰 파싱 실패: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims 타입 변환 실패")
	}

	if _, exists := claims["query_hash"]; exists {
		t.Error("빈 rawQuery일 때 query_hash가 존재하면 안 됨")
	}
}
