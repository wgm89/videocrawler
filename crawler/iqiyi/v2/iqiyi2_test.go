package iqiyiv2

import "testing"

func Test_iqiyi2(t *testing.T) {
	iqiyi := NewIqiyi2()
	iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr7w54oo.html")
	//iqiyi.GetVideoInfo("http://www.iqiyi.com/a_19rrh9vl6t.html")
}
