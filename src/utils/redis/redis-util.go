package redisUtil

import (
	"strings"
)

var redisSpecialChars = []string{
	"-", "[", "]", "{", "}", "(", ")", "*", "~", ":", "\"", "'", "|", "&", "!", "<", ">", ".",
}

func PrepareParam(input string) string {
	escaped := input
	for _, char := range redisSpecialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

func PrepareParamForEval(input string) string {
	replacer := strings.NewReplacer(
		`"`, `\\\"`,
		`'`, `\\\'`,
		`\\`, `\\\\`,
	)
	escaped := replacer.Replace(input)
	for _, char := range redisSpecialChars {
		if char == `"` || char == `'` || char == `\\` {
			continue
		}
		escaped = strings.ReplaceAll(escaped, char, `\\`+char)
	}
	return escaped
}
