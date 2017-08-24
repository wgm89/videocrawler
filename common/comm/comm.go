package comm

import (
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"videocrawler/common/randomAgent"
	"videocrawler/common/util"

	"golang.org/x/net/publicsuffix"
)

type CookieItem struct {
	Name  string
	Value string
}

type Comm struct{}

func (this *Comm) FakeIp() string {
	var ipSegs []string
	r := rand.New(rand.NewSource(int64(time.Now().Unix())))
	for i := 0; i < 4; i++ {
		ipSegs = append(ipSegs, strconv.Itoa(r.Intn(255)))
	}
	return strings.Join(ipSegs, ".")
}

func (this *Comm) SetHeader(req *http.Request,
	headers map[string]string) map[string]string {

	reqHeaders := map[string]string{
		"Accept-Encoding": "gzip, deflate, sdch",
		"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4",
	}
	if clientIp, ok := headers["CLIENT-IP"]; !ok {
		clientIp = this.FakeIp()
		reqHeaders["CLIENT-IP"] = clientIp
		reqHeaders["X-FORWARDED-FOR"] = clientIp
	}
	if _, ok := headers["User-Agent"]; !ok {
		reqHeaders["User-Agent"] = randomAgent.GetRandomAgent()
	}
	for k, v := range headers {
		reqHeaders[k] = v
	}
	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}
	return reqHeaders
}

func (this *Comm) Unzip(resp *http.Response) io.ReadCloser {
	var reader io.ReadCloser
	var err error
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil
		}
	} else {
		reader = resp.Body
	}
	return reader
}

func (this *Comm) GetMasterDomain(domain string) string {
	if this.IsIp(domain) {
		return domain
	}
	if strings.Contains(domain, "http") {
		u, err := url.Parse(domain)
		if err != nil {
			panic(err)
		}
		domain = u.Host
	}
	suffix, ok := publicsuffix.PublicSuffix(domain)
	if !ok {
		return domain
	}
	i := len(domain) - len(suffix)
	if i <= 0 || domain[i-1] != '.' {
		return domain
	}
	lastIndex := strings.LastIndex(domain[:i-1], ".")
	return domain[lastIndex+1:]
}

func (this *Comm) ExportCookie(dir string, domain string, jar *cookiejar.Jar) error {
	u, _ := url.Parse(domain)
	cookies := jar.Cookies(u)
	u.Host = this.GetMasterDomain(domain)
	cookieFile := path.Join(dir, u.Host+"_cookie.txt")
	f, err := os.OpenFile(cookieFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0700)
	if err != nil {
		panic(err)
	}
	f.Truncate(0)
	f.Seek(0, 0)
	defer f.Close()
	var cookieList = make([]string, 0, 10)
	expire := (time.Now().UnixNano() / int64(time.Second)) + 30*24*60*60
	expireStr := strconv.FormatInt(expire, 10)
	for _, v := range cookies {
		cookieStr := "." + u.Host + "\t" + "FALSE" + "\t" + "/" + "\t" + "FALSE" +
			"\t" + expireStr + "\t" + v.Name + "\t" + v.Value + "\n"
		cookieList = append(cookieList, cookieStr)
	}
	for _, cookieStr := range cookieList {
		if _, err := f.WriteString(cookieStr); err != nil {
			return err
		}
	}
	return nil
}

func (this *Comm) LoadCookie(dir string, domain string, jar *cookiejar.Jar) error {
	u, _ := url.Parse(domain)
	md := this.GetMasterDomain(domain)
	u.Host = md
	cookieFile := path.Join(dir, u.Host+"_cookie.txt")
	expire := (time.Now().UnixNano() / int64(time.Second)) + 30*24*60*60
	ti := time.Unix(expire, 0)
	cl := make([]*http.Cookie, 0)

	util.ReadLine(cookieFile, func(line string) {
		lineS := strings.Split(line, "\t")
		if len(lineS) != 7 {
			return
		}
		cl = append(cl, &http.Cookie{
			Name:     lineS[5],
			Value:    lineS[6],
			Path:     "/",
			Domain:   u.Host,
			Expires:  ti,
			Secure:   false,
			HttpOnly: false,
		})
	})
	jar.SetCookies(u, cl)
	return nil
}

func (this *Comm) LoadCookieFromList(domain string, jar *cookiejar.Jar, cookies []CookieItem) error {
	u, _ := url.Parse(domain)
	md := this.GetMasterDomain(domain)
	u.Host = md
	expire := (time.Now().UnixNano() / int64(time.Second)) + 30*24*60*60
	ti := time.Unix(expire, 0)
	cl := make([]*http.Cookie, 0)

	for _, ci := range cookies {
		fmt.Println(ci)
		fmt.Println(u.Host)
		cl = append(cl, &http.Cookie{
			Name:     ci.Name,
			Value:    ci.Value,
			Path:     "/",
			Domain:   u.Host,
			Expires:  ti,
			Secure:   false,
			HttpOnly: false,
		})
	}
	jar.SetCookies(u, cl)
	return nil

}

func (this *Comm) IsIp(host string) bool {
	return net.ParseIP(host) != nil
}

func New() *Comm {
	return &Comm{}
}
