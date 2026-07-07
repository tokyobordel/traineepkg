package strings

import (
	"regexp"
	"strings"
	"unicode"

	"traineesheep/imageservice/pkg/errors"
)

func PascalToUpperSnakeCase(input string) (string, errors.DomainError) {
	if input == "" {
		return "", errors.NewInvalidParametersError("input", input, "входная строка не может быть пустой")
	}

	if !isPascalCase(input) {
		return "", errors.NewInvalidParametersError("input", input, "строка не соответствует формату PascalCase")
	}

	var result strings.Builder

	for i, r := range input {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToUpper(r))
	}

	return result.String(), nil
}

func isPascalCase(s string) bool {
	if s == "" {
		return false
	}

	if !unicode.IsUpper(rune(s[0])) {
		return false
	}

	pascalRegex := regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`)
	return pascalRegex.MatchString(s)
}
