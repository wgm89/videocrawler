package iqiyiv1

import "testing"

func Test_iqiyi(t *testing.T) {
	iqiyi := New()
	//iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr7w54oo.html")
	//iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr80lv9c.html")

	//iqiyi.Login("", "")
	//iqiyi.ExportCookie()

	iqiyi.GetVideoInfo("http://www.iqiyi.com/v_19rr7s1zis.html")

}
