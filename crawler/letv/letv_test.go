package letv

import (
	"testing"
)

func Test_letv(t *testing.T) {
	lt := New()
	lt.GetVideoInfo("http://www.le.com/ptv/vplay/30249668.html")
}
