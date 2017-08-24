package crawler

import (
	"net/http"
	"net/http/cookiejar"
	"testing"

	"videocrawler/common/comm"
)

func Test_crawler(t *testing.T) {
	videoCookieJar, _ := cookiejar.New(nil)
	c := &CrawlerNet{
		Comm: &comm.Comm{},
		Client: &http.Client{
			Jar: videoCookieJar,
		},
		Headers: map[string]string{},
	}
	c.CrossdomainCheck("http://data.video.iqiyi.com/videos/v0/20170605/34/f7/5c2b746bd92b440d80562a967ddc4c87.ts?qdv=1&qypid=692992100_04022000001000000000_96&start=0&end=154633&hsize=8356&tag=0&v=&contentlength=86104&qd_uid=&qd_vip=0&qd_src=3_31_312&qd_tm=1499501811793&qd_ip=76c7d6ad&qd_p=76c7d6ad&qd_k=7eb580d76bfae71faec5b6ce095cecbc&qd_sc=5828dfac7490eacbd8a728a5f3bc56c3")
}
