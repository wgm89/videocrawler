package qq

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"videocrawler/clients"
	"videocrawler/crawler"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

type Qq struct {
	*crawler.CrawlerNet
	Client  *clients.CrClient
	videoId string
	title   string
	streams crawler.StreamSet
	pageUrl string
}

func (this *Qq) Init(param crawler.CrawlerInitParam) {

}

func (this *Qq) GetVideoInfo(pageUrl string) (crawler.VideoDetail, error) {
	this.pageUrl = pageUrl
	this.getVideoId()
	this.parseRes()

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	for k, v := range this.streams {
		fmt.Println(k)
		for i := 0; i < len(v.Flv); i++ {
			fmt.Println(v.Flv[i].Src)
			c = color.New(color.FgCyan).Add(color.Underline)
			c.Println("=========================================================")
		}
	}
	vDetail := crawler.VideoDetail{
		Title:   this.title,
		Streams: this.streams,
	}
	return vDetail, nil
}

func (this *Qq) getVideoId() string {
	fmt.Println("page url:", this.pageUrl)
	body, err := this.Client.GetUrlContent(this.pageUrl)
	if err != nil {
		return ""
	}
	content := string(body)

	r, _ := regexp.Compile(`vid\s*:\s*"\s*([^"]+)`)
	matches := r.FindStringSubmatch(content)
	videoId := matches[1]
	this.videoId = videoId

	r, _ = regexp.Compile("<title>([^<]+)")
	matches = r.FindStringSubmatch(content)
	this.title = matches[1]

	return ""
}

func (this *Qq) parseRes() error {
	requestUrl := "http://vv.video.qq.com/getinfo?otype=json&appver=3%2E2%2E19%2E333&platform=11&defnpayver=1&vid=" + this.videoId
	body, err := this.Client.GetUrlContent(requestUrl)
	if err != nil {
		return err
	}
	content := strings.Replace(string(body), "QZOutputJson=", "", 1)
	content = strings.TrimRight(content, ";")
	body = []byte(content)

	partsVid, _ := jsonparser.GetString(body, "vl", "vi", "[0]", "vid")
	partsTi, _ := jsonparser.GetString(body, "vl", "vi", "[0]", "ti")
	this.title = partsTi
	partsPrefix, _ := jsonparser.GetString(body, "vl", "vi", "[0]", "ul", "ui", "[0]", "url")
	partsPrefix = strings.TrimRight(partsPrefix, "/")

	var firstFi []byte
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		sl, _ := jsonparser.GetInt(value, "sl")
		if sl != 0 {
			firstFi = value
		}
	}, "fl", "fi")
	if len(firstFi) == 0 {
		firstFi, _, _, _ = jsonparser.Get(body, "fl", "fi", "[0]")
	}
	firstFi, _, _, _ = jsonparser.Get(body, "fl", "fi", "[2]")

	formatId, _ := jsonparser.GetInt(firstFi, "id")
	formatSl, _ := jsonparser.GetInt(firstFi, "sl")
	streamId, _ := jsonparser.GetString(firstFi, "name")
	profile, _ := jsonparser.GetString(firstFi, "cname")

	var mp4Urls []string
	if formatSl == 0 {
		for i := 1; i < 100; i++ {
			filename := fmt.Sprintf("%s.p%d.%d.mp4", this.videoId, formatId%10000, i)
			keyApi := fmt.Sprintf("http://vv.video.qq.com/getkey?"+
				"otype=json&platform=11&format=%d&vid=%s&filename=%s",
				formatId, partsVid, filename)
			body, _ := this.Client.GetUrlContent(keyApi)
			content := strings.Replace(string(body), "QZOutputJson=", "", 1)
			content = strings.TrimRight(content, ";")
			body = []byte(content)
			vkey, err := jsonparser.GetString(body, "key")
			if err != nil {
				break
			}
			url := fmt.Sprintf("%s/%s?vkey=%s", partsPrefix, filename, vkey)
			mp4Urls = append(mp4Urls, url)
		}
	} else {
		fvKey, _ := jsonparser.GetString(body, "vl", "vi", "[0]", "fvkey")
		mp4, _, _, err := jsonparser.Get(body, "vl", "vi", "[0]", "cl", "ci")
		if err != nil {
			mp4, _, _, _ = jsonparser.Get(body, "vl", "vi", "[0]", "fn")
		} else {
			keyId, _ := jsonparser.GetString(mp4, "[0]", "keyid")
			oldId := strings.Split(keyId, ".")[1]
			oldIdInt, _ := strconv.ParseInt(oldId, 10, 64)
			newId := "p" + fmt.Sprintf("%d", (oldIdInt%1000))
			mp4 = []byte(strings.Replace(keyId, oldId, newId, -1) + ".mp4")
		}
		mp4Url := fmt.Sprintf("%s/%s?vkey=%s", partsPrefix, string(mp4), fvKey)
		mp4Urls = append(mp4Urls, mp4Url)
	}
	var flvs []crawler.FlvInfo
	for i := 0; i < len(mp4Urls); i++ {
		flvs = append(flvs, crawler.FlvInfo{
			Src:  mp4Urls[i],
			Size: 0,
		})
	}
	/*
		if len(mp4Urls) > 0 {
			flvs = append(flvs, FlvInfo{
				Src:  mp4Urls[0],
				Size: 0,
			})
		}
	*/
	width, _ := jsonparser.GetInt(body, "vl", "vi", "[0]", "vw")
	height, _ := jsonparser.GetInt(body, "vl", "vi", "[0]", "vh")

	this.streams[streamId] = crawler.StreamInfo{
		VideoProfile: profile,
		Container:    "",
		Src:          "",
		Width:        width,
		Height:       height,
		Flv:          flvs,
	}
	return nil
}

func New() crawler.Crawler {
	fmt.Println("init qq#########################################")
	qq := &Qq{
		&crawler.CrawlerNet{},
		clients.NewClient("https://v.qq.com", map[string]string{}),
		"", "", make(crawler.StreamSet), "",
	}
	return qq
}
