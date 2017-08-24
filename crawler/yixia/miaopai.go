package yixia

import (
	"fmt"
	"regexp"

	"videocrawler/clients"
	"videocrawler/common/util"
	"videocrawler/crawler"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

type Miaopai struct {
	*crawler.CrawlerNet
	Client  *clients.CrClient
	streams crawler.StreamSet
	videoId string
	title   string
	pageUrl string
}

func (this *Miaopai) Init(param crawler.CrawlerInitParam) {

}

func (this *Miaopai) GetVideoInfo(url string) (crawler.VideoDetail, error) {
	this.GetVideoId(url)

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	if this.videoId == "" {
		return crawler.VideoDetail{}, nil
	}
	params := util.HttpBuildQuery(map[string]string{
		"scid":     this.videoId,
		"vend":     "miaopai",
		"fillType": "259",
	})
	reqUrl := fmt.Sprintf("http://api.miaopai.com/m/v2_channel.json?" + params)
	body, err := this.Client.GetUrlContent(reqUrl)
	if err != nil {
		return crawler.VideoDetail{}, nil
	}
	flvSrc, _ := jsonparser.GetString(body, "result", "stream", "base")
	this.title, _ = jsonparser.GetString(body, "result", "ext", "t")

	fmt.Println(flvSrc)

	flvs := []crawler.FlvInfo{}
	flvs = append(flvs, crawler.FlvInfo{
		Src:  flvSrc,
		Size: 0,
	})
	streamInfo := crawler.StreamInfo{
		VideoProfile: "",
		Container:    "",
		Src:          "",
		Width:        0,
		Height:       0,
		Flv:          flvs,
	}
	this.streams["normal"] = streamInfo

	vDetail := crawler.VideoDetail{
		Title:   this.title,
		Streams: this.streams,
	}
	return vDetail, nil
}

func (this *Miaopai) GetVideoId(url string) string {
	body, err := this.Client.GetUrlContent(url)
	if err != nil {
		return ""
	}
	r, _ := regexp.Compile(`"scid":"(.+?)"`)
	matches := r.FindStringSubmatch(string(body))

	if matches != nil {
		this.videoId = matches[1]
	} else {
		r, _ = regexp.Compile(`scid\s=\s"(.+?)"`)
		matches = r.FindStringSubmatch(string(body))
		if matches != nil {
			this.videoId = matches[1]
		}
	}
	r, _ = regexp.Compile("<title>([^<]+)")
	matches = r.FindStringSubmatch(string(body))
	this.title = matches[1]

	return this.videoId
}

func New() crawler.Crawler {
	fmt.Println("init miaopai###############################################")
	mp := &Miaopai{
		&crawler.CrawlerNet{},
		clients.NewClient("https://miaopai.com", map[string]string{}),
		make(crawler.StreamSet), "", "", "",
	}
	return mp
}
