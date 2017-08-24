package weibo

import (
	"testing"
)

func Test_weibo(t *testing.T) {
	wb := New()
	//wb.GetUrlContent("http://weibo.com")
	//wb.ExportCookie()
	//wb.GetVideoInfo("http://weibo.com/tv/v/FcX9lhsFQ?from=vhot")
	wb.GetVideoInfo("http://weibo.com/tv/v/Fd4hM958S?from=vhot")
}
