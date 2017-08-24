package letv

import (
	"crypto/sha1"
	"fmt"
	_ "math"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"videocrawler/clients"
	"videocrawler/crawler"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

type Letv struct {
	*crawler.CrawlerNet
	Client  *clients.CrClient
	videoId string
	title   string
	streams crawler.StreamSet
	pageUrl string
}

func (this *Letv) Init(param crawler.CrawlerInitParam) {

}

func (this *Letv) GetVideoInfo(pageUrl string) (crawler.VideoDetail, error) {

	this.pageUrl = pageUrl
	var matched bool
	if matched, _ = regexp.MatchString("http://yuntv.letv.com/*", pageUrl); matched {
		fmt.Println("yuntv")
	} else if matched = strings.Contains(pageUrl, "sports.le.com"); matched {
		fmt.Println("sports")
	} else {
		this.parseTvData()
	}
	vDetail := crawler.VideoDetail{
		Title:   this.title,
		Streams: this.streams,
	}
	return vDetail, nil
}

func (this *Letv) parseTvData() error {
	fmt.Println("tv")
	this.getVideoId()

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	key := this.calcTimeKey()
	requestUrl := fmt.Sprintf("http://player-pc.le.com/mms/out/video/playJson?"+
		"id=%s&platid=1&splatid=101&format=1&tkey=%d"+
		"&domain=www.le.com&region=cn&source=1000&accesyx=1",
		this.videoId, key)
	body, err := this.Client.GetUrlContent(requestUrl)
	if err != nil {
		return err
	}
	info, _, _, _ := jsonparser.Get(body, "msgs")

	dispatch := make(map[string][]byte)
	jsonparser.ObjectEach(info, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		dispatch[string(key)] = value
		return nil
	}, "playurl", "dispatch")

	var dispatchInfo []byte
	var streamId string
	for streamId, dispatchInfo = range dispatch {
		break
	}
	fmt.Println("streamId:", streamId)

	domains := make([]string, 0, 3)
	jsonparser.ArrayEach(info, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		domains = append(domains, string(value))
	}, "playurl", "domain")

	dispatchSrcs := make([]string, 0, 2)
	jsonparser.ArrayEach(dispatchInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		dispatchSrcs = append(dispatchSrcs, string(value))
	})

	playUrl := domains[0] + dispatchSrcs[0]
	fmt.Println(playUrl)
	h := sha1.New()
	h.Write([]byte(playUrl))
	bs := h.Sum(nil)
	uuid := fmt.Sprintf("%x_0", bs)
	//ext := path.Ext(dispatchSrcs[1])
	playUrl = strings.Replace(playUrl, "tss=0", "tss=ios", -1)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := fmt.Sprintf("%f", r.Float64())
	playUrl += fmt.Sprintf("&m3v=1&termid=1&format=1&hwtype=un&ostype=MacOS10.12.4&p1=1&p2=10&p3=-&expect=3&tn=%s&vid=%s&uuid=%s&sign=letv", t, this.videoId, uuid)

	var body2, m3Content []byte
	body2, err = this.Client.GetUrlContent(playUrl)

	m3Url, _ := jsonparser.GetString(body2, "location")

	m3Url += "&r=" + fmt.Sprintf("%d", int64(time.Now().UnixNano()/1000000)) + "&appid=500"
	m3Content, err = this.Client.GetUrlContent(m3Url)
	m3Data := this.decodeData(m3Content)

	_ = m3Data

	return nil
}

func (this *Letv) getVideoId() string {
	body, err := this.Client.GetUrlContent(this.pageUrl)
	if err != nil {
		return ""
	}
	content := string(body)
	var r *regexp.Regexp

	r, _ = regexp.Compile(`vplay/(\d+).html`)
	matches := r.FindStringSubmatch(this.pageUrl)
	if matches != nil && len(matches) == 2 && matches[1] != "" {
		this.videoId = matches[1]
	} else {
		r, _ = regexp.Compile(`\Wvid\s*?[\:=]\s*?["\']?(\d+)["\']?`)
		matches := r.FindStringSubmatch(content)
		if matches != nil && len(matches) == 2 {
			this.videoId = matches[1]
		}
	}

	r, _ = regexp.Compile("<title>([^<]+)")
	matches = r.FindStringSubmatch(content)
	this.title = matches[1]

	return this.videoId
}

func (this *Letv) decodeData(data []byte) string {
	version := data[0:5]
	if strings.ToLower(string(version)) == "vc_01" {
		fmt.Println("yes")
		loc2 := data[5:]
		length := len(loc2)
		loc4 := make([]byte, 2*length)
		for i := 0; i < length; i++ {
			loc4[2*i] = loc2[i] >> 4
			loc4[2*i+1] = loc2[i] & 15
		}
		loc6 := append(loc4[len(loc4)-11:], loc4[:len(loc4)-11]...)
		loc7 := make([]byte, length)
		for i := 0; i < length; i++ {
			loc7[i] = (loc6[2*i] << 4) + loc6[2*i+1]
		}
		var decodeStr string = ""
		for i := 0; i < len(loc7); i++ {
			decodeStr += string(loc7[i])
		}
		return decodeStr
	} else {
		return string(data)
	}
}

func (this *Letv) calcTimeKey() int64 {
	var magic int64 = 185025305
	t := int64(time.Now().Unix())
	rBits := uint64(magic % 17)
	ror := ((t & ((1 << 32) - 1)) >> (rBits % 32)) | (t << (32 - (rBits % 32)) & ((1 << 32) - 1))
	return ror ^ magic
}

func New() crawler.Crawler {
	fmt.Println("init letv################################################")
	lt := &Letv{
		&crawler.CrawlerNet{},
		clients.NewClient("http://www.le.com", map[string]string{}),
		"", "", make(crawler.StreamSet), "",
	}
	return lt
}
