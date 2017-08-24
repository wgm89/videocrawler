package iqiyi

import (
	"videocrawler/clients"
	"videocrawler/crawler"
)

const SRC = "76f90cbd92f94a2e925d83e8ccd22cb7"
const KEY = "d5fb4bd9d50c4be6948c97edd7254b0e"

var Ids []string = []string{"4k", "BD", "TD", "HD", "SD", "LD"}
var Vd2Id map[int64]string = map[int64]string{
	10: "4k", 19: "4k", 5: "BD", 18: "BD",
	21: "HD", 2: "HD", 4: "TD", 17: "TD",
	96: "LD", 1: "SD",
}
var Id2Profile map[string]string = map[string]string{
	"4k": "4k", "BD": "1080p", "TD": "720p",
	"HD": "540p", "SD": "360p", "LD": "210p",
}

type IqiyiT struct {
	*crawler.CrawlerNet
	Client  *clients.CrClient
	Streams crawler.StreamSet
	VideoId string
	Tvid    string
	Title   string
}
