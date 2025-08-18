package common

import "encoding/json"

func HasHexPrefix(s string) bool {
	return len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X")
}

func TrimHexPrefix(s string) string {
	if len(s) >= 2 && (s[0:2] == "0x" || s[0:2] == "0X") {
		return s[2:]
	}
	return s
}

func JSON(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}
