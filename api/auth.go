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

// GenerateToken generates a JWT HS512 authentication token.
// query: key=value map extracted from GET/DELETE query parameters or POST body.
// See https://docs.upbit.com/reference/%EC%9D%B8%EC%A6%9D for authentication docs.
func GenerateToken(accessKey, secretKey string, query map[string]string) (string, error) {
	claims := jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	}

	if len(query) > 0 {
		// Per Upbit auth docs: hash is computed from the unencoded query string.
		queryString := buildQueryHashString(query)
		hash := sha512.Sum512([]byte(queryString))
		claims["query_hash"] = hex.EncodeToString(hash[:])
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

// GenerateTokenFromBody generates a token from POST body fields (JSON struct fields).
// Upbit POST requests require the body parameters to be hashed as a query string.
func GenerateTokenFromBody(accessKey, secretKey string, bodyParams map[string]string) (string, error) {
	return GenerateToken(accessKey, secretKey, bodyParams)
}

// GenerateTokenFromRawQuery generates a token from a raw query string.
// Supports array parameters (e.g. uuids[]=a&uuids[]=b) that cannot be expressed as map[string]string.
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

// buildQueryHashString builds an unencoded query string for hash generation.
// Per Upbit auth docs: "Hash must be computed from the URL-unencoded query string."
// Keys are sorted for deterministic output (map iteration order is undefined).
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
		// Concatenate without URL encoding (for hashing).
		pairs = append(pairs, k+"="+params[k])
	}
	return strings.Join(pairs, "&")
}

// BuildQueryString converts a parameter map to a URL-encoded query string.
// Used when constructing HTTP request URLs (applies URL encoding).
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
