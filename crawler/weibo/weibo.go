package weibo

import (
	"fmt"
	"regexp"

	"videocrawler/clients"
	"videocrawler/crawler"

	"github.com/fatih/color"
)

type Weibo struct {
	*crawler.CrawlerNet
	Client  *clients.CrClient
	streams crawler.StreamSet
	videoId string
	title   string
	pageUrl string
}

func (this *Weibo) Init(param crawler.CrawlerInitParam) {

}

func (this *Weibo) GetVideoInfo(url string) (crawler.VideoDetail, error) {
	this.pageUrl = url
	mp4 := this.GetVideoId(url)

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")
	flvs := []crawler.FlvInfo{}
	flvs = append(flvs, crawler.FlvInfo{
		Src:  mp4,
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

func (this *Weibo) GetVideoId(url string) string {
	var mp4 string
	body, err := this.Client.GetUrlContent(url)
	if err != nil {
		return ""
	}

	r, _ := regexp.Compile(`"stream_url":\s"(.+)"`)
	matches := r.FindStringSubmatch(string(body))
	if matches != nil {
		mp4 = matches[1]
	}

	r, _ = regexp.Compile(`"status_title":\s"(.+?)"`)
	matches = r.FindStringSubmatch(string(body))
	if matches != nil {
		this.title = matches[1]
	}

	return mp4
}

func New() crawler.Crawler {
	fmt.Println("init weibo###############################################")
	mp := &Weibo{
		&crawler.CrawlerNet{},
		clients.NewClient("http://weibo.com", map[string]string{}),
		make(crawler.StreamSet), "", "", "",
	}
	return mp
}
