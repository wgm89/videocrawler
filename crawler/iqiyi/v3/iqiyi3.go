package iqiyiv3

import (
	"fmt"
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
	"github.com/satori/go.uuid"
)

type Iqiyi3 struct {
	*iqiyi.IqiyiT
}

func (this *Iqiyi3) GetVideoInfo(url string) (crawler.VideoDetail, error) {
	this.GetVideoId(url)

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.Title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	if this.VideoId == "" || this.Tvid == "" {
		return crawler.VideoDetail{}, nil
	}
	vData := this.getVps()
	this.parseData(vData)
	vDetail := crawler.VideoDetail{
		Title:   this.Title,
		Streams: this.Streams,
	}
	return vDetail, nil
}

func (this *Iqiyi3) getVps() []byte {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	u := uuid.NewV4()
	macid := strings.Replace(u.String(), "-", "", -1)

	ts := strconv.FormatInt(currentTime, 10)
	host := "http://cache.video.qiyi.com"
	params := "/vps?tvid=" + this.Tvid + "&vid=" + this.VideoId + "&v=0&qypid=" +
		this.Tvid + "_12&src=01012001010000000000&t=" + ts + "&k_tag=1&k_uid=" +
		macid + "&rs=1"
	vf := this.getVf(params)
	reqUrl := host + params + "&vf=" + vf
	body, _ := this.Client.GetUrlContent(reqUrl)
	return body
}

func (this *Iqiyi3) getVf(urlParams string) string {
	sufix := ""
	var v4, v8 int
	for j := 0; j < 8; j++ {
		for k := 0; k < 4; k++ {
			v4 = 13 * (66*k + 27*j) % 35
			if v4 >= 10 {
				v8 = v4 + 88
			} else {
				v8 = v4 + 49
			}
			sufix += string(v8)
		}
	}
	urlParams += sufix
	return util.Md5(urlParams)
}

func (this *Iqiyi3) parseData(vData []byte) error {
	code, _ := jsonparser.GetString(vData, "code")
	if code != "A00000" {
		return fmt.Errorf("code error")
	}
	urlPrefix, _ := jsonparser.GetString(vData, "data", "vp", "du")
	stream, _, _, _ := jsonparser.Get(vData, "data", "vp", "tkl", "[0]", "vs")
	jsonparser.ArrayEach(stream, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		bid, _ := jsonparser.GetInt(value, "bid")
		streamId := iqiyi.Vd2Id[bid]

		c := color.New(color.FgCyan)
		c.Println(streamId)

		profile := iqiyi.Id2Profile[streamId]
		flvs := make([]crawler.FlvInfo, 0, 5)

		jsonparser.ArrayEach(value, func(val []byte, dataType jsonparser.ValueType, offset int, err error) {
			l, _ := jsonparser.GetString(val, "l")
			realUrl := this.getRealUrl(urlPrefix + l)
			fmt.Println(realUrl)
			flvs = append(flvs, crawler.FlvInfo{
				Src:  realUrl,
				Size: 0,
			})
		}, "fs")

		fmt.Println("")
		fmt.Println("========================================")

		this.Streams[streamId] = crawler.StreamInfo{
			VideoProfile: profile,
			Container:    "m3u8",
			Src:          "",
			Width:        0,
			Height:       0,
			Flv:          flvs,
		}
	})
	return nil
}

func (this *Iqiyi3) getRealUrl(reqUrl string) string {
	body, err := this.Client.GetUrlContent(reqUrl)
	if err != nil {
		return ""
	}
	l, _ := jsonparser.GetString(body, "l")
	return l
}

func (this *Iqiyi3) GetVideoId(url string) (string, string) {
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

func NewIqiyi3() *Iqiyi3 {
	fmt.Println("init iqiyi###############################################")
	iqy := &Iqiyi3{
		&iqiyi.IqiyiT{
			&crawler.CrawlerNet{},
			clients.NewClient("http://iqiyi.com", map[string]string{}),
			make(crawler.StreamSet), "", "", "",
		},
	}
	return iqy
}
