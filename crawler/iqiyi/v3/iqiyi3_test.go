package iqiyiv3

import "testing"

func Test_iqiyi3(t *testing.T) {
	iqiyi := NewIqiyi3()
	//iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr7w54oo.html")
	//iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr80lv9c.html")

	//iqiyi.Login("", "")
	//iqiyi.ExportCookie()

	iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr7hix8s.html")

}
