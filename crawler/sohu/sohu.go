package sohu

import (
	"errors"
	"fmt"
	"math/rand"
	netUrl "net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"videocrawler/clients"
	"videocrawler/crawler"

	"code.google.com/p/mahonia"
	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

type vidMap map[string]string

type Sohu struct {
	*crawler.CrawlerNet
	Client  *clients.CrClient
	videoId string
	title   string
	streams crawler.StreamSet
	pageUrl string
	isMy    bool
}

var shSteamTypes = map[string]string{
	"nor": "普通",
	"hig": "高清",
	"sup": "超清",
}

func (this *Sohu) Init(param crawler.CrawlerInitParam) {

}

func (this *Sohu) GetVideoInfo(pageUrl string) (crawler.VideoDetail, error) {

	this.GetVideoId(pageUrl)

	if this.videoId == "" {
		return crawler.VideoDetail{}, errors.New("not found vid")
	}
	this.pageUrl = pageUrl

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	matched, _ := regexp.MatchString("http://tv.sohu.com/*", pageUrl)
	this.isMy = !matched
	if matched {
		this.parseHotRes()
	} else {
		this.parseMyRes()
	}

	vDetail := crawler.VideoDetail{
		Title:   this.title,
		Streams: this.streams,
	}
	return vDetail, nil
}

func (this *Sohu) parseHotRes() {
	var exists error

	fmt.Println("搜狐视频")

	body := this.getVideoRes()
	_, exists = jsonparser.GetString(body, "allot")
	if exists != nil {
		for _, qtype := range []string{"oriVid", "superVid",
			"highVid", "norVid", "relativeId"} {
			data, _, _, exi := jsonparser.Get(body, "data")
			var hqvid int64
			if exi != nil {
				hqvid, _ = jsonparser.GetInt(body, qtype)
			} else {
				hqvid, _ = jsonparser.GetInt(data, qtype)
			}
			hqvidStr := strconv.FormatInt(hqvid, 10)
			if hqvid != 0 && hqvidStr != this.videoId {
				this.videoId = hqvidStr
				body = this.getVideoRes()
				_, exi = jsonparser.GetString(body, "allot")
				if exi == nil {
					break
				}
			}
		}
	}

	this.parseBody(body)
}

func (this *Sohu) getVideoRes() []byte {
	requestUrl := fmt.Sprintf("http://hot.vrs.sohu.com/vrs_flash.action?vid=%s", this.videoId)
	fmt.Println(requestUrl)
	content, err := this.Client.GetUrlContent(requestUrl)
	if err != nil {
		return nil
	}
	return content
}

func (this *Sohu) parseMyRes() string {
	requestUrl := fmt.Sprintf("http://my.tv.sohu.com/play/videonew.do?vid=%s&referer=http://my.tv.sohu.com", this.videoId)
	fmt.Println("自媒体")
	fmt.Println(requestUrl)
	content, err := this.Client.GetUrlContent(requestUrl)
	if err != nil {
		return ""
	}

	this.parseBody(content)

	return requestUrl
}

func (this *Sohu) parseBody(body []byte) {
	var sus, clips, cks []string
	var allot string

	data, _, _, exists := jsonparser.Get(body, "data")
	if exists != nil {
		data = body
	}

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		sus = append(sus, string(value))
	}, "su")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		cks = append(cks, string(value))
	}, "ck")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		clips = append(clips, string(value))
	}, "clipsURL")

	allot, _ = jsonparser.GetString(body, "allot")
	tvid, _ := jsonparser.GetString(body, "tvid")

	flvs := make([]crawler.FlvInfo, 0, 10)
	for k, v := range clips {
		urlInfo, _ := netUrl.Parse(string(v))
		clipPath := urlInfo.Path
		realUrl := this.getRealUrl(allot, this.videoId, tvid, sus[k], clipPath, cks[k])
		flvs = append(flvs, crawler.FlvInfo{
			Src:  realUrl,
			Size: 0,
		})
		fmt.Printf("第%d段视频地址:%s\n", k, realUrl)

		c := color.New(color.FgCyan).Add(color.Underline)
		c.Println("=========================================================")
	}

	width, _ := jsonparser.GetInt(body, "width")
	height, _ := jsonparser.GetInt(body, "height")

	var streamId string
	if !this.isMy {
		vids := make(vidMap)
		vids = this.getVidMap(data)
		streamId = vids[this.videoId]
	} else {
		streamId = "none"
	}

	this.streams[streamId] = crawler.StreamInfo{
		VideoProfile: streamId,
		Container:    "m3u8",
		Src:          "",
		Width:        width,
		Height:       height,
		Flv:          flvs,
	}
}

func (this *Sohu) changeCharset(body string) string {
	decoder := mahonia.NewDecoder("gbk")
	content := decoder.ConvertString(body)
	return content
}

func (this *Sohu) getVidMap(data []byte) vidMap {
	vids := make(vidMap)
	jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		kStr := string(key)
		if ok := strings.HasSuffix(kStr, "Vid"); ok && kStr != "tvid" {
			vids[string(value)] = string(key)
		}
		return nil
	})
	return vids
}

func (this *Sohu) getRealUrl(host string, vid string, tvid string, new string,
	clipUrl string, ck string) string {
	ts := strconv.FormatInt(int64(time.Now().Unix()), 10)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := fmt.Sprintf("%f", r.Float64())
	requestUrl := "http://" + host + "/?prot=9&prod=flash&pt=1&file=" + clipUrl +
		"&new=" + new + "&key=" + ck + "&vid=" + vid + "&uid=" + ts +
		"&t=" + t + "&rb=1"
	content, err := this.Client.GetUrlContent(requestUrl)
	if err != nil {
		return ""
	}
	realUrl, _ := jsonparser.GetString(content, "url")
	return realUrl
}

func (this *Sohu) GetVideoId(pageUrl string) string {
	body, err := this.Client.GetUrlContent(pageUrl)
	if err != nil {
		return ""
	}
	content := this.changeCharset(string(body))

	r, _ := regexp.Compile(`\Wvid\s*?[\:=]\s*?["\'](\d+)["\']`)
	matches := r.FindStringSubmatch(content)

	if matches != nil {
		this.videoId = matches[1]
	}

	r, _ = regexp.Compile("<title>([^<]+)")
	matches = r.FindStringSubmatch(content)
	this.title = matches[1]

	return this.videoId
}

func New() crawler.Crawler {
	sh := &Sohu{
		&crawler.CrawlerNet{},
		clients.NewClient("http://tv.sohu.com", map[string]string{}),
		"", "", make(crawler.StreamSet), "",
		false,
	}
	return sh
}
