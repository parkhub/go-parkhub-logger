package log

import (
	"regexp"
	"strings"
	"testing"
	"unicode"
)

type stringFunc func(s string) string

const testWhitespaceString = "\n  \t some message \n\t    "
const testRuneString = "LMeXcqgt3MEJBferDhkCvWB6UNvPwZ7kLCpQEwrcC3RTpYxMoGDtBk87MUbAdP4gfLnXxLkka9QzRRTb2J2CGtuoFcTx9MidhPxD"

var leftRegexp = regexp.MustCompile(`^(\s)+`)
var rightRegexp = regexp.MustCompile(`(\s)+$`)

func leadingWhitespaceRegex(s string) string {
	var trimmedLeft string
	left := leftRegexp.FindAllStringSubmatch(s, -1)
	if len(left) > 0 {
		trimmedLeft = left[0][0]
	}
	return trimmedLeft
}

func trailingWhitespaceRegex(s string) string {
	var trimmedRight string
	right := rightRegexp.FindAllStringSubmatch(s, -1)
	if len(right) > 0 {
		trimmedRight = right[0][0]
	}
	return trimmedRight
}

func leadingWhitespaceScan(s string) string {
	var b strings.Builder
	for i, _ := range s {
		if !unicode.IsSpace(rune(s[i])) {
			return b.String()
		}
		b.WriteString(string(rune(s[i])))
	}
	return b.String()
}

func trailingWhitespaceScan(s string) string {
	var rev strings.Builder
	for i := len(s); i > 0; i-- {
		if !unicode.IsSpace(rune(s[i-1])) {
			return reverseStringSwap(rev.String())
		}
		rev.WriteString(string(rune(s[i-1])))
	}
	return rev.String()
}

func reverseStringSwap(s string) string {
	r := []rune(s)
	for i, j := 0, 0; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func BenchmarkStringFuncs(b *testing.B) {
	b.Run("leading whitespace", func(b *testing.B) {
		b.Run("regex", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = leadingWhitespaceRegex(testWhitespaceString)
			}
		})

		b.Run("scan", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = leadingWhitespaceScan(testWhitespaceString)
			}
		})
	})

	b.Run("trailing whitespace", func(b *testing.B) {
		b.Run("regex", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = trailingWhitespaceRegex(testWhitespaceString)
			}
		})

		b.Run("scan", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = trailingWhitespaceScan(testWhitespaceString)
			}
		})
	})
}

func BenchmarkBuildString(b *testing.B) {
	b.Run("by index and casting", func(b *testing.B) {
		b.Run("WriteString(string(rune(s[i])))", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var b strings.Builder
				for i, _ := range testRuneString {
					_ = unicode.IsSpace(rune(testRuneString[i]))
					b.WriteString(string(rune(testRuneString[i])))
				}
			}
		})
		b.Run("WriteString(string(s[i]))", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var b strings.Builder
				for i, _ := range testRuneString {
					_ = unicode.IsSpace(rune(testRuneString[i]))
					b.WriteString(string(testRuneString[i]))
				}
			}
		})
		b.Run("Write([]byte{s[i]})", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var b strings.Builder
				for i, _ := range testRuneString {
					_ = unicode.IsSpace(rune(testRuneString[i]))
					b.Write([]byte{testRuneString[i]})
				}
			}
		})
	})

	b.Run("by range value", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var b strings.Builder
			for _, r := range testRuneString {
				b.WriteString(string(r))
			}
		}
	})
}
