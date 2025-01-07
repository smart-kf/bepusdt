package sign

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

func ValidateSign(data map[string]any, token string) bool {
	var signString string
	reqSign, ok := data["sign"]
	if ok {
		signString = reqSign.(string)
	}

	var jsonKeys = make([]string, 0, len(data))
	for k := range data {
		if k == "sign" {
			continue
		}
		jsonKeys = append(jsonKeys, k)
	}
	sort.Strings(jsonKeys)
	var builder strings.Builder
	for idx, k := range jsonKeys {
		builder.WriteString(fmt.Sprintf("%s=%v", k, data[k]))
		if idx != len(jsonKeys)-1 {
			builder.WriteString("&")
		}
	}
	builder.WriteString(token)
	h := sha256.New()
	h.Write([]byte(builder.String()))
	x := h.Sum(nil)
	return fmt.Sprintf("%x", x) == signString
}

func SetSign(data map[string]any, token string) error {
	var jsonKeys = make([]string, 0, len(data))
	for k := range data {
		if k == "sign" {
			continue
		}
		jsonKeys = append(jsonKeys, k)
	}
	sort.Strings(jsonKeys)
	var builder strings.Builder
	for idx, k := range jsonKeys {
		builder.WriteString(fmt.Sprintf("%s=%v", k, data[k]))
		if idx != len(jsonKeys)-1 {
			builder.WriteString("&")
		}
	}
	builder.WriteString(token)
	h := sha256.New()
	h.Write([]byte(builder.String()))
	x := h.Sum(nil)
	data["sign"] = fmt.Sprintf("%x", x)
	return nil
}
