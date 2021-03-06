package utils

import (
	"os"
	"strings"
	"unicode"
)

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

func StripUnderscore(s string) string {
	s = strings.Replace(s, "_", "", -1)
	return s
}

func StripLineBreakChar(s string) string {
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\r", "", -1)
	return s
}

func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return false
}

// CreateDir 创建文件夹
func CreateDir(path string) error {
	if !FileIsExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
