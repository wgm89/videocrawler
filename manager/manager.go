package manager

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"videocrawler/crawler"
	"videocrawler/crawler/iqiyi/v1"
	"videocrawler/crawler/letv"
	"videocrawler/crawler/qq"
	"videocrawler/crawler/sohu"
	"videocrawler/crawler/weibo"
	"videocrawler/crawler/yixia"
	"videocrawler/crawler/youku"
)

var crawlers = make(map[string]crawler.InsFunc)

func initCrawlers() {
	crawlers = map[string]crawler.InsFunc{
		"youku":   youku.New,
		"qq":      qq.New,
		"iqiyi":   iqiyiv1.New,
		"sohu":    sohu.New,
		"le":      letv.New,
		"miaopai": yixia.New,
		"weibo":   weibo.New,
	}
}

func GetVideoInfo(url string) (crawler.VideoDetail, error) {
	var (
		module string
		getIns crawler.InsFunc
		err    error
	)
	var vd = crawler.VideoDetail{}
	module, getIns, err = matchModule(url)
	fmt.Println(module)
	if err == nil {
		crawler := getIns()
		vd, err = crawler.GetVideoInfo(url)
	}
	return vd, err
}

func matchModule(videoUrl string) (string, crawler.InsFunc, error) {
	u, e := url.Parse(videoUrl)
	if e != nil {
		return "", nil, e
	}
	host := u.Host
	//TODO
	for k, v := range crawlers {
		if strings.Contains(host, k) {
			return k, v, nil
		}
	}
	return "", nil, errors.New("not found")
}

func init() {
	fmt.Println("init manager##################################################")
	initCrawlers()
}
