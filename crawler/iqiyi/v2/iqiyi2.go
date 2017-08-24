package iqiyiv2

import (
	"errors"
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
)

type Iqiyi2 struct {
	*iqiyi.IqiyiT
}

func (this *Iqiyi2) GetVideoInfo(url string) (crawler.VideoDetail, error) {
	this.GetVideoId(url)

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.Title)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")
	if this.VideoId == "" || this.Tvid == "" {
		return crawler.VideoDetail{}, errors.New("not found vid")
	}
	err := this.parseVideo()
	if err != nil {
		return crawler.VideoDetail{}, err
	}
	vDetail := crawler.VideoDetail{
		Title:   this.Title,
		Streams: this.Streams,
	}
	return vDetail, nil
}

func (this *Iqiyi2) parseVideo() error {
	body := this.getRawData()
	content := string(body)
	content = strings.TrimLeft(content, "var tvInfoJs=")
	body = []byte(content)
	code, _ := jsonparser.GetString(body, "code")
	if code != "A00000" {
		return errors.New("json data format error")
	}
	data, _, _, _ := jsonparser.Get(body, "data")
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		_, e := jsonparser.GetString(value, "m3utx")
		if e != nil {
			return
		}
		vd, _ := jsonparser.GetString(value, "vd")
		vdInt, err := strconv.ParseInt(vd, 10, 64)
		m3u8, _ := jsonparser.GetString(value, "m3utx")
		streamId := iqiyi.Vd2Id[vdInt]
		videoProfile := iqiyi.Id2Profile[streamId]

		fmt.Println(streamId)
		fmt.Println("m3u8: " + m3u8)

		c := color.New(color.FgCyan).Add(color.Underline)
		c.Println("=========================================================")

		this.Streams[streamId] = crawler.StreamInfo{
			VideoProfile: videoProfile,
			Container:    "m3u8",
			Src:          m3u8,
			Width:        0,
			Height:       0,
			Flv:          make([]crawler.FlvInfo, 0),
		}
	}, "vidl")
	return nil
}

func (this *Iqiyi2) getRawData() []byte {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	timeStr := fmt.Sprintf("%d", currentTime)
	key := "d5fb4bd9d50c4be6948c97edd7254b0e"
	sc := util.Md5(timeStr + key + this.Tvid)

	reqUrl := fmt.Sprintf("http://cache.m.iqiyi.com/jp/tmts/%s/%s/", this.Tvid, this.VideoId)
	var params = map[string]string{
		"tvid": this.Tvid,
		"vid":  this.VideoId,
		"src":  "76f90cbd92f94a2e925d83e8ccd22cb7",
		"sc":   sc,
		"t":    timeStr,
	}
	reqUrl = reqUrl + "?" + util.HttpBuildQuery(params)
	body, _ := this.Client.GetUrlContent(reqUrl)
	return body
}

func (this *Iqiyi2) GetVideoId(url string) (string, string) {
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

func NewIqiyi2() *Iqiyi2 {
	fmt.Println("init iqiyi###############################################")
	iqy := &Iqiyi2{
		&iqiyi.IqiyiT{
			&crawler.CrawlerNet{},
			clients.NewClient("http://iqiyi.com", map[string]string{}),
			make(crawler.StreamSet), "", "", "",
		},
	}
	return iqy
}
