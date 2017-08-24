package comm

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/user"
	"testing"
	"time"
)

func Test_cookie(t *testing.T) {
	jar, _ := cookiejar.New(nil)
	u, e := url.Parse("http://www.baidu.com")
	if e != nil {
		panic(e)
	}

	ti := time.Unix(1600196526, 0)
	cl := make([]*http.Cookie, 0)
	cl = append(cl, &http.Cookie{
		Name:     "name",
		Value:    "value",
		Path:     "/",
		Domain:   "baidu.com",
		Expires:  ti,
		Secure:   false,
		HttpOnly: false,
	})
	cl = append(cl, &http.Cookie{
		Name:     "hi",
		Value:    "abc",
		Path:     "/",
		Domain:   "baidu.com",
		Expires:  ti,
		Secure:   false,
		HttpOnly: false,
	})

	jar.SetCookies(u, cl)

	usr, _ := user.Current()
	fmt.Println(usr)

	c := New()
	c.ExportCookie(usr.HomeDir+"/.crawler/cookies", "http://www.baidu.com", jar)

	//jar, _ = cookiejar.New(nil)
	//c.LoadCookie(usr.HomeDir+"/.crawler/cookies", "http://www.baidu.com", jar)
	fmt.Println("===================")
	//fmt.Println(jar.Cookies(u))
}
