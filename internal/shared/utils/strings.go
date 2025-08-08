package utils

import (
	"regexp"
	"strings"
	"sync"
)

var (
	// Pre-compile regex for better performance
	camelCaseRegex *regexp.Regexp
	regexOnce      sync.Once
)

// ToSnakeCase converts camelCase or PascalCase strings to snake_case
// Examples: PersonID -> person_id, UserName -> user_name, firstName -> first_name
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	// Initialize regex only once for better performance
	regexOnce.Do(func() {
		camelCaseRegex = regexp.MustCompile("([a-z0-9])([A-Z])")
	})

	// Insert underscore before uppercase letters
	snake := camelCaseRegex.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamelCase converts snake_case strings to camelCase
// Examples: person_id -> personId, user_name -> userName
func ToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	var result strings.Builder
	result.WriteString(parts[0]) // First part stays lowercase

	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result.WriteString(strings.ToUpper(string(parts[i][0])))
			if len(parts[i]) > 1 {
				result.WriteString(parts[i][1:])
			}
		}
	}

	return result.String()
}

// ToPascalCase converts snake_case strings to PascalCase
// Examples: person_id -> PersonId, user_name -> UserName
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.Split(s, "_")
	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])))
			if len(part) > 1 {
				result.WriteString(part[1:])
			}
		}
	}

	return result.String()
}
