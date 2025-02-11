package utils

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func Sha1Object(v interface{}) string {
	x, _ := json.Marshal(v)
	a := sha1.New()
	a.Write(x)
	res := a.Sum(nil)
	return fmt.Sprintf("%x", res)
}
