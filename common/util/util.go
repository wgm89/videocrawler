package util

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetParentDirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func Md5(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	cipherStr := md5Ctx.Sum(nil)
	sc := hex.EncodeToString(cipherStr)
	return sc
}

func HttpBuildQuery(params map[string]string) (param_str string) {
	uv := url.Values{}
	for k, v := range params {
		uv.Add(k, v)
	}
	return uv.Encode()
}
