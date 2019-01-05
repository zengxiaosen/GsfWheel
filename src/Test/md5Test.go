package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {

	fmt.Printf("%s\n", str2Md5("1111"))
	a := ""
	fmt.Println(len(a))
}

func str2Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
