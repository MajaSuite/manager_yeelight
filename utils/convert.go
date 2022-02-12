package utils

import (
	"strconv"
	"strings"
)

func ConvertHex(v string) uint32 {
	res, _ := strconv.ParseUint(strings.Replace(v, "0x", "", -1), 16, 64)
	return uint32(res)
}

func ConvertBool(v string) bool {
	if v == "off" {
		return false
	}
	return true
}

func ConverArray(v string) []string {
	return strings.Split(v, " ")
}

func ConvertInt(v string) int {
	if i, err := strconv.Atoi(v); err != nil {
		return 0
	} else {
		return i
	}
}

func ConvertToString(msg map[string]interface{}, k string) string {
	switch msg[k].(type) {
	case string:
		return msg[k].(string)
	}
	return ""
}

func ConvertToInt(msg map[string]interface{}, k string) int {
	switch msg[k].(type) {
	case float64:
		return int(msg[k].(float64))
	}
	return 0
}

func ConvertToBool(msg map[string]interface{}, k string) bool {
	switch msg[k].(type) {
	case bool:
		return msg[k].(bool)
	}
	return false
}
