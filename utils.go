package jsonschemaline

import (
	"strings"

	"goa.design/goa/v3/codegen"
)

// 封装 goa.design/goa/v3/codegen 方便后续可定制
func ToCamel(name string) string {
	return codegen.CamelCase(name, true, true)
}

func ToLowerCamel(name string) string {
	return codegen.CamelCase(name, false, true)
}

func SnakeCase(name string) string {
	return codegen.SnakeCase(name)
}

func Addslashes(str string) string {
	var tmpRune []rune
	for _, ch := range str {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// BaseName 获取最后.后的文本
func BaseName(fullname string) (baseName string) {
	baseName = fullname
	lastDotIndex := strings.LastIndex(baseName, ".")
	if lastDotIndex > -1 {
		baseName = baseName[lastDotIndex+1:]
	}
	return baseName
}

// Namespace 获取最后.前的文本
func Namespace(fullname string) (namespace string) {
	lastDotIndex := strings.LastIndex(fullname, ".")
	if lastDotIndex > -1 {
		namespace = fullname[:lastDotIndex]
	}
	return namespace
}
