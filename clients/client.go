package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"videocrawler/common/comm"
	"videocrawler/common/util"
	"videocrawler/env"
)

type Entry struct {
	jar    *cookiejar.Jar
	domain string
	ready  chan struct{}
}

type Jars struct {
	JarTable map[string]*Entry
	mu       *sync.Mutex
	comm     *comm.Comm
}

func (this *Jars) GetJar(domain string) *Entry {
	this.mu.Lock()
	u, err := url.Parse(domain)
	if err != nil {
		return nil
	}
	host := u.Host
	e := this.JarTable[host]
	if e == nil {
		e = &Entry{
			ready: make(chan struct{}),
		}
		this.JarTable[host] = e
		this.mu.Unlock()
		e.domain = domain
		e.jar, err = cookiejar.New(nil)
		this.LoadCookie(e)
		close(e.ready)
	} else {
		this.mu.Unlock()
		<-e.ready
	}
	return e
}

func (this *Jars) LoadCookie(e *Entry) {
	fmt.Println(e.domain)
	this.comm.LoadCookie(env.CookieDir, e.domain, e.jar)
}

func (this *Jars) ExportCookie(e *Entry) {
	this.comm.ExportCookie(env.CookieDir, e.domain, e.jar)
}

type CrClient struct {
	Comm     *comm.Comm
	JarEntry *Entry
	Client   *http.Client
	Headers  map[string]string
	Domain   string
}

func (this *CrClient) SetHeaders(headers map[string]string) {
	for key, value := range headers {
		this.Headers[key] = value
	}
}

func (this *CrClient) SetCookies(cookies []comm.CookieItem) {
	this.Comm.LoadCookieFromList(this.JarEntry.domain, this.JarEntry.jar, cookies)
}

func (this *CrClient) GetUrlContent(requestUrl string) ([]byte, error) {
	var err error
	var resp *http.Response
	var req *http.Request

	req, err = http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	this.Headers = this.Comm.SetHeader(req, this.Headers)

	resp, err = this.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := this.Comm.Unzip(resp)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	//jars.ExportCookie(this.JarEntry)
	return body, nil
}

func (this *CrClient) CrossdomainCheck(pageUrl string) bool {
	u, e := url.Parse(pageUrl)
	if e != nil {
		return false
	}
	fileUrl := u.Scheme + u.Host + "crossdomain.xml"
	fmt.Println(fileUrl)
	body, _ := this.GetUrlContent(pageUrl)
	if body == nil {
		return false
	}
	content := string(body)
	r, e := regexp.Compile(`.+allow-access-from domain="\*".+`)
	if r.MatchString(content) {
		return true
	}
	return false
}

func (this *CrClient) ConvertContentEncode(content string) string {
	var encode string
	r, _ := regexp.Compile(`content="text/html;\s*?charset=(.+?)"`)
	matches := r.FindStringSubmatch(content)
	if matches != nil {
		encode = strings.ToUpper(matches[1])
		if encode != "" && encode != "UTF-8" {
			content = util.CovertEncode(content, encode)
		}
	}
	return content
}

var jars *Jars

func init() {
	jars = &Jars{
		JarTable: make(map[string]*Entry),
		mu:       new(sync.Mutex),
		comm:     comm.New(),
	}
}

func NewClient(domain string, headers map[string]string) *CrClient {
	jarEntry := jars.GetJar(domain)
	return &CrClient{
		Comm:     comm.New(),
		JarEntry: jarEntry,
		Client: &http.Client{
			Jar: jarEntry.jar,
		},
		Headers: headers,
		Domain:  domain,
	}
}
