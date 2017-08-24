package iqiyi

import (
	"fmt"
	"net/url"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	_ "time"

	"github.com/buger/jsonparser"
)

/*
func (this *IqiyiT) Rsa(data string) string {
	var N float64 = 0xab86b6371b5318aaa1d3c9e612a9f1264f372323c8c0f19875b5fc3b3fd3afcc1e5bec527aa94bfa85bffc157e4245aebda05389a5357b75115ac94f074aefcd
	var e float64 = 65537
	return this.ohdaveRsaEncrypt(data, e, N)
}

func (this *IqiyiT) ohdaveRsaEncrypt(data string, exponent float64, modulus float64) string {

	dataHex := hex.EncodeToString([]byte(data))

	dataInt64, _ := strconv.ParseUint(dataHex, 16, 64)
	dataF := float64(dataInt64)
	encryptData := math.Mod(math.Pow(dataF, exponent), modulus)
	fmt.Println(encryptData)
	return fmt.Sprintf("%x", encryptData)
}
*/

func (this *IqiyiT) getFormData(username, password string) string {
	_, filename, _, _ := runtime.Caller(1)
	pyScript := path.Join(path.Dir(filename), "encrypt_pwd.py")
	out, _ := exec.Command(pyScript, "-u"+username, "-p"+password).Output()
	fmt.Println(out)
	return string(out)
}

func (this *IqiyiT) decodeCode(code string) string {
	r, _ := regexp.Compile(`}\('(.+)',(\d+),(\d+),'([^']+)'\.split\('\|'\)`)
	matches := r.FindStringSubmatch(code)
	obfucastedCode, base, count, symbolStr := matches[1], matches[2], matches[3], matches[4]

	baseInt, _ := strconv.ParseInt(base, 10, 64)
	countInt, _ := strconv.ParseInt(count, 10, 64)
	symbols := strings.Split(symbolStr, "|")
	symbolTable := make(map[string]string)

	for countInt > 0 {
		countInt -= 1
		baseNcount := this.encodeBaseN(countInt, baseInt)

		if countInt < int64(len(symbols)) {
			symbolTable[baseNcount] = symbols[countInt]
		} else {
			symbolTable[baseNcount] = baseNcount
		}
	}
	fmt.Println(obfucastedCode)
	r, _ = regexp.Compile(`\b(\w+)\b`)
	matches = r.FindStringSubmatch(obfucastedCode)
	for _, m := range matches[1:] {
		obfucastedCode = strings.Replace(obfucastedCode, m, symbolTable[m], -1)
	}
	//return re.sub( r'\b(\w+)\b', lambda mobj: symbol_table[mobj.group(0)],obfucasted_code)
	return obfucastedCode
}

func (this *IqiyiT) encodeBaseN(num, n int64) string {
	fullTable := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if num == 0 {
		return string(fullTable[0])
	}
	var ret string = ""
	for num > 0 {
		ret = string(fullTable[num%n]) + ret
		num = num / n
	}
	return ret
}

func (this *IqiyiT) Login(username string, password string) error {
	formDataStr := this.getFormData(username, password)
	formDataByte := []byte(formDataStr)

	target, _ := jsonparser.GetString(formDataByte, "target")
	bird_t, _ := jsonparser.GetString(formDataByte, "bird_t")
	sign, _ := jsonparser.GetString(formDataByte, "sign")
	token, _ := jsonparser.GetString(formDataByte, "token")
	bird_src, _ := jsonparser.GetString(formDataByte, "bird_src")
	server, _ := jsonparser.GetString(formDataByte, "server")

	formData := map[string]string{
		"target":   target,
		"bird_t":   bird_t,
		"sign":     sign,
		"token":    token,
		"bird_src": bird_src,
		"server":   server,
	}
	uv := url.Values{}
	for k, v := range formData {
		uv.Add(k, v)
	}
	urlParams := uv.Encode()
	reqUrl := "http://kylin.iqiyi.com/validate?" + urlParams
	fmt.Println(reqUrl)
	content, err := this.Client.GetUrlContent(reqUrl)
	if err != nil {
		return err
	}
	fmt.Println("empty")
	fmt.Println(string(content))
	return nil
}
