package main

import _ "fmt"

import "videocrawler/crawler"

func main() {
	var youku *crawler.Youku = crawler.NewYouku()
	//fmt.Println(youku.GetCna())
	_, _ = youku.GetVideoInfo("http://v.youku.com/v_show/id_XMjg1Mzg0MzkyOA==.html?spm=a2hww.20023042.m_223465.5~5~5~5!2~5~5~A&f=50219377", 0)
	//fmt.Println(res)
}
