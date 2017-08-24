package sohu

import (
	"testing"
)

func Test_sohu(t *testing.T) {
	sh := New()
	//sh.GetVideoInfo("http://tv.sohu.com/20170629/n600029461.shtml")
	//sh.GetVideoInfo("http://tv.sohu.com/20170630/n600029870.shtml")
	//sh.GetVideoInfo("http://tv.sohu.com/20170623/n600019963.shtml")
	sh.GetVideoInfo("http://my.tv.sohu.com/pl/9354876/90079223.shtml")
	//sh.GetVideoInfo("http://tv.sohu.com/20170216/n480929591.shtml")
}
