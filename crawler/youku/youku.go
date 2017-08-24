/**
 * author: wangguangmao
 */

package youku

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"videocrawler/clients"
	"videocrawler/common/comm"
	"videocrawler/common/util"
	"videocrawler/crawler"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

var streamTypes = map[string]map[string]string{
	"mp4hd3": {"id": "mp4hd3", "aliasOf": "hd3"},
	"hd3":    {"id": "hd3", "container": "flv", "videoProfile": "1080P"},
	"mp4hd2": {"id": "mp4hd2", "aliasOf": "hd2"},
	"hd2":    {"id": "hd2", "container": "flv", "videoProfile": "超清"},
	"mp4hd":  {"id": "mp4hd", "aliasOf": "mp4"},
	"mp4":    {"id": "mp4", "container": "mp4", "videoProfile": "高清"},
	"flvhd":  {"id": "flvhd", "container": "flv", "videoProfile": "标清"},
	"flv":    {"id": "flv", "container": "flv", "videoProfile": "标清"},
	"3gphd":  {"id": "3gphd", "container": "3gp", "videoProfile": "标清（3GP）"},
}

type Youku struct {
	*crawler.CrawlerNet
	Client     *clients.CrClient
	utid       string
	videoId    string
	streams    crawler.StreamSet
	retry      int
	videoTitle string
}

func (this *Youku) Init(crawler.CrawlerInitParam) {

}

func (this *Youku) GetVideoInfo(videoUrl string) (crawler.VideoDetail, error) {

	fmt.Println(this.getYsuid())
	cookies := make([]comm.CookieItem, 0, 5)
	cookies = append(cookies, comm.CookieItem{
		Name:  "__ysuid",
		Value: this.getYsuid(),
	})
	cookies = append(cookies, comm.CookieItem{
		Name:  "xreferrer",
		Value: "http://www.youku.com",
	})
	this.Client.SetCookies(cookies)

	videoId := this.extractId(videoUrl)
	var err error
	if videoId == "" {
		err = errors.New("url error, not found video id")
		return crawler.VideoDetail{}, err
	}
	this.Client.SetHeaders(map[string]string{
		"Referer": videoUrl,
	})
	ts := int64(time.Now().Unix())
	params := map[string]string{
		"vid":       videoId,
		"ccode":     "0401",
		"client_ip": "192.168.1.1",
		"utid":      this.getCna(),
		"client_ts": strconv.FormatInt(ts, 10),
	}
	paramStr := util.HttpBuildQuery(params)
	videoInfoUrl := "https://ups.youku.com/ups/get.json?" + paramStr

	body, err := this.Client.GetUrlContent(videoInfoUrl)
	if err != nil {
		return crawler.VideoDetail{}, err
	}
	if respContentErr, err := jsonparser.GetString(body, "data", "error"); err == nil {
		code, _ := jsonparser.GetInt([]byte(respContentErr), "code")
		errStr, _ := jsonparser.GetString([]byte(respContentErr), "note")

		if code == -6004 {
			if this.retry == 0 {
				this.utid, _ = url.QueryUnescape(this.utid)
				this.retry = 1
				return this.GetVideoInfo(videoUrl)
			} else if this.retry == 1 {
				this.utid = this.getCna()
				this.retry = 2
				return this.GetVideoInfo(videoUrl)
			}
			return crawler.VideoDetail{}, errors.New(errStr)
		} else if code == -3307 {
			return crawler.VideoDetail{}, errors.New(errStr)
		} else if code == -2004 {
			return crawler.VideoDetail{}, errors.New(errStr)
		}
	} else {
		data, _, _, _ := jsonparser.Get(body, "data")
		this.parseRes(data)
	}
	vDetail := crawler.VideoDetail{
		Title:   this.videoTitle,
		Streams: this.streams,
	}
	return vDetail, nil
}

func (this *Youku) extractId(url string) (vid string) {
	vid = ""
	r, _ := regexp.Compile("id_(.+).*\\.html")
	matches := r.FindStringSubmatch(url)
	if matches != nil {
		vid = matches[1]
	}
	return
}

func (this *Youku) parseRes(data []byte) {
	this.videoTitle, _ = jsonparser.GetString(data, "video", "title")

	d := color.New(color.FgBlue, color.Bold)
	d.Printf("视频标题: %s \n", this.videoTitle)

	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("=========================================================")

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		streamType, _ := jsonparser.GetString(value, "stream_type")
		milliSeconds, _ := jsonparser.GetInt(value, "milliseconds_video")
		streamSize, _ := jsonparser.GetInt(value, "size")
		sInfo, _ := streamTypes[streamType]

		if alias, ok := sInfo["aliasOf"]; ok {
			streamType = alias
			sInfo = streamTypes[alias]
		}
		container := sInfo["container"]
		src, _ := jsonparser.GetString(value, "m3u8_url")
		videoProfile, _ := sInfo["videoProfile"]
		width, _ := jsonparser.GetInt(value, "width")
		height, _ := jsonparser.GetInt(value, "height")

		this.streams[streamType] = crawler.StreamInfo{
			VideoProfile: videoProfile,
			Container:    container,
			Src:          src,
			Width:        width,
			Height:       height,
			Flv:          this.getFlv(value),
		}

		fmt.Println("视频类型", streamType)
		fmt.Println("视频时长", this.millisecondsToTime(milliSeconds))
		fmt.Printf("视频大小(M)%.2f\n", float64(streamSize)/float64(2<<20))
		fmt.Println("m3u8地址:", src)

		c := color.New(color.FgCyan).Add(color.Underline)

		c.Println("=========================================================")

	}, "stream")

}

func (this *Youku) getFlv(data []byte) []crawler.FlvInfo {
	var segs []crawler.FlvInfo
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		cdnUrl, _ := jsonparser.GetString(value, "cdn_url")
		sizeByte, _, _, _ := jsonparser.Get(value, "size")
		sizeInt, err := strconv.ParseInt(string(sizeByte), 10, 64)
		segs = append(segs, crawler.FlvInfo{Src: cdnUrl, Size: sizeInt})
	}, "segs")
	fmt.Printf("flv数量:%d\n", len(segs))
	return segs
}

func (this *Youku) millisecondsToTime(milliseconds int64) string {
	seconds := milliseconds / 1000
	h := seconds / 3600
	s := seconds % 3600
	m := (s / 60)
	s = s % 60
	return fmt.Sprintf("%d小时, %d分钟, %d秒", h, m, s)
}

func (this *Youku) getCna() string {
	cnaUrl := "http://log.mmstat.com/eg.js"
	body, err := this.Client.GetUrlContent(cnaUrl)
	if err != nil {
		return ""
	}
	r, _ := regexp.Compile("Etag=\"(.+?)\"")
	matches := r.FindStringSubmatch(string(body))
	if matches != nil {
		return matches[1]
	}
	return ""
}

func (this *Youku) getYsuid() string {
	currentTime := time.Now().UnixNano() / int64(time.Second)
	return fmt.Sprintf("%d%s", currentTime, util.RandString(3))
}

func New() crawler.Crawler {
	fmt.Println("init youku################################################")
	cl := clients.NewClient("https://youku.com", map[string]string{})
	youku := &Youku{
		&crawler.CrawlerNet{},
		cl,
		url.QueryEscape("onBdERfZriwCAW+uM3cVByOa"),
		"", make(crawler.StreamSet), 0, "",
	}

	return youku
}
