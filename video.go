package videocrawler

import (
	"fmt"
	"net/url"
	"videocrawler/crawler"
	"videocrawler/env"
	"videocrawler/manager"
)

type VideoRes struct {
	Title       string
	M3          map[string]string
	Flv         []string
	DomainLimit bool
}

func GetVideoInfo(videoUrl string) (*VideoRes, error) {
	if e := checkUrl(videoUrl); e != nil {
		return nil, e
	}
	fmt.Println("get video info ...")
	videoInfo, err := manager.GetVideoInfo(videoUrl)
	if err != nil {
		return nil, err
	}
	title := videoInfo.Title
	if err != nil {
		return nil, err
	}
	//var quality string
	var stream crawler.StreamInfo
	var m3Urls = make(map[string]string)
	var flvUrls = make([]string, 0)
	for _, stream = range videoInfo.Streams {
		if stream.Src != "" {
			m3Urls[stream.VideoProfile] = stream.Src
		}
	}
	for _, stream = range videoInfo.Streams {
		if len(stream.Flv) != 0 {
			for _, v := range stream.Flv {
				flvUrls = append(flvUrls, v.Src)
			}
			break
		}
	}
	res := &VideoRes{
		Title:       title,
		M3:          m3Urls,
		Flv:         flvUrls,
		DomainLimit: false,
	}
	return res, nil
}

func checkUrl(videoUrl string) error {
	_, err := url.Parse(videoUrl)
	return err
}

func init() {
	env.LoadEnv()
}
