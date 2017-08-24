package youku

import "testing"

func Test_getvideo(t *testing.T) {
	youku := New()
	_, _ = youku.GetVideoInfo("http://v.youku.com/v_show/id_XMjkxMjgwNDc0NA==.html")
	//_, _ = youku.GetVideoInfo("http://v.youku.com/v_show/id_XMzc0OTM4NDg0.html?spm=a2h0k.8191407.0.0&from=s1.8-1-1.2")
}
