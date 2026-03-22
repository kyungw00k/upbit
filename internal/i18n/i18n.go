package i18n

import (
	"fmt"
	"os"
	"strings"
)

// Key is a message key type.
type Key string

var current map[Key]string
var isKorean bool

func init() {
	// POSIX locale 우선순위: LC_ALL > LC_MESSAGES > LANG
	lang := os.Getenv("LC_ALL")
	if lang == "" {
		lang = os.Getenv("LC_MESSAGES")
	}
	if lang == "" {
		lang = os.Getenv("LANG")
	}
	if strings.HasPrefix(lang, "ko") {
		current = ko
		isKorean = true
	} else {
		current = en
		isKorean = false
	}
}

// T returns the translated message for the given key.
func T(key Key) string {
	if msg, ok := current[key]; ok {
		return msg
	}
	return string(key)
}

// Tf returns the translated message with format arguments.
func Tf(key Key, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}

// IsKorean returns true if the current locale is Korean.
func IsKorean() bool {
	return isKorean
}
