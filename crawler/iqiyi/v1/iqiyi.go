package iqiyiv1

import (
	"fmt"
	_ "os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"videocrawler/clients"
	"videocrawler/common/util"
	"videocrawler/crawler"
	"videocrawler/crawler/iqiyi"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

type Iqiyi struct {
	*iqiyi.IqiyiT
}

func (this *Iqiyi) Init(param crawler.CrawlerInitParam) {

}

func (this *Iqiyi) GetVideoInfo(url string) (crawler.VideoDetail, error) {
	this.GetVideoId(url)

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.Title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	if this.VideoId == "" || this.Tvid == "" {
		return crawler.VideoDetail{}, nil
	}

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	sc := util.Md5(strconv.FormatInt(currentTime, 10) + iqiyi.KEY + this.VideoId)
	requestUrl := fmt.Sprintf("http://cache.m.iqiyi.com/tmts/%s/%s/?t=%d&sc=%s&src=%s",
		this.Tvid,
		this.VideoId,
		currentTime,
		sc,
		iqiyi.SRC,
	)
	body, err := this.Client.GetUrlContent(requestUrl)
	if err != nil {
		return crawler.VideoDetail{}, err
	}
	this.GetSegListUrl(body)

	for k, v := range this.Streams {
		fmt.Println(k)
		fmt.Printf("分辨率:%s\n", v.VideoProfile)
		fmt.Printf("宽:%d, 高:%d\n", v.Width, v.Height)
		fmt.Printf("m3u8:%s\n", v.Src)
		c = color.New(color.FgCyan).Add(color.Underline)
		c.Println("=========================================================")
	}
	vDetail := crawler.VideoDetail{
		Title:   this.Title,
		Streams: this.Streams,
	}
	return vDetail, nil
}

func (this *Iqiyi) GetSegListUrl(body []byte) error {
	code, _ := jsonparser.GetString(body, "code")

	if code != "A00000" {
		return nil
	}
	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		vd, _ := jsonparser.GetInt(value, "vd")
		src, _ := jsonparser.GetString(value, "m3u")
		screenSize, _ := jsonparser.GetString(value, "screenSize")
		if _, ok := iqiyi.Vd2Id[vd]; !ok {
			return
		}
		streamId, _ := iqiyi.Vd2Id[vd]
		profile, _ := iqiyi.Id2Profile[streamId]
		size := strings.Split(screenSize, "x")

		width, _ := strconv.ParseInt(size[0], 10, 64)
		height, _ := strconv.ParseInt(size[1], 10, 64)

		this.Streams[streamId] = crawler.StreamInfo{
			VideoProfile: profile,
			Container:    "m3u8",
			Src:          src,
			Width:        width,
			Height:       height,
		}
	}, "data", "vidl")
	return nil
}

func (this *Iqiyi) GetVideoId(url string) (string, string) {
	body, err := this.Client.GetUrlContent(url)
	if err != nil {
		return "", ""
	}
	r, _ := regexp.Compile("data-player-videoid=\"(.+?)\"")
	matches := r.FindStringSubmatch(string(body))

	if matches != nil {
		this.VideoId = matches[1]
	}
	r, _ = regexp.Compile("data-player-tvid=\"(.+?)\"")
	matches = r.FindStringSubmatch(string(body))
	if matches != nil {
		this.Tvid = matches[1]
	}
	r, _ = regexp.Compile("<title>([^<]+)")
	matches = r.FindStringSubmatch(string(body))
	this.Title = matches[1]

	return this.VideoId, this.Tvid
}

func New() crawler.Crawler {
	fmt.Println("init iqiyi###############################################")
	iqy := &Iqiyi{
		&iqiyi.IqiyiT{
			&crawler.CrawlerNet{},
			clients.NewClient("http://iqiyi.com", map[string]string{}),
			make(crawler.StreamSet), "", "", "",
		},
	}
	return iqy
}
