package utils

import (
	"strings"

	"github.com/agelinazf/egb"
)

func StringToIntArray(str string) []int {
	if str == "" {
		return make([]int, 0)
	}
	intarr := make([]int, 0)
	if strings.Contains(str, ",") {
		result, _ := egb.StringToIntArray(str)
		return result
	} else {
		tempint := egb.StringToInt(str)
		intarr = append(intarr, tempint)
		return intarr
	}
}

func StringToStringArray(str string) []string {
	if str == "" {
		return make([]string, 0)
	}
	intarr := make([]string, 0)
	if strings.Contains(str, ",") {
		return strings.Split(str, ",")
	} else {
		intarr = append(intarr, str)
		return intarr
	}
}

func StringArrayToString(strArr []string) string {
	length := len(strArr)
	if length > 1 {
		return strings.Join(strArr, ",")
	}
	if length == 1 {
		return strArr[0]
	}
	return ""
}
