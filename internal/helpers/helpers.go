package helpers

import (
	"strings"
	"unicode"
)

func IsValidInput(input string) bool {
	invalidInputs := map[string]bool{
		"error": true,
		"-Inf":  true,
		"+Inf":  true,
		"NaN":   true,
		"-NaN":  true,
	}

	return !invalidInputs[input]
}

func ReplaceEConstant(expression string) string {
	var result strings.Builder
	runes := []rune(expression)

	for i := 0; i < len(runes); i++ {
		if runes[i] == 'e' {
			if i > 0 && unicode.IsDigit(runes[i-1]) && i+1 < len(runes) && (runes[i+1] == '+' || runes[i+1] == '-') {
				result.WriteRune('e')
			} else {
				result.WriteString("2.71828182846")
			}
		} else {
			result.WriteRune(runes[i])
		}
	}

	return result.String()
}
