package crawler

import (
	"net/http/cookiejar"
)

type FlvInfo struct {
	Src  string
	Size int64
}

type StreamInfo struct {
	VideoProfile string
	Container    string
	Src          string
	Width        int64
	Height       int64
	Flv          []FlvInfo
}

type StreamSet map[string]StreamInfo

type VideoDetail struct {
	Title   string
	Streams StreamSet
	Extra   interface{}
}

type InsFunc func() Crawler

type Crawler interface {
	GetVideoInfo(url string) (VideoDetail, error)
	Init(param CrawlerInitParam)
}

type CrawlerInitParam struct {
	Jar *cookiejar.Jar
}

type CrawlerNet struct {
}
