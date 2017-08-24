package manager

import (
	"fmt"
	"testing"
)

func Test_manage(t *testing.T) {
	r, e := GetVideoInfo("http://www.iqiyi.com/v_19rr7hix8s.html")
	fmt.Println(r)
	fmt.Println(e)
}
