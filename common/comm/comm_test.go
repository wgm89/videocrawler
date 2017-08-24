package comm

import "testing"
import "fmt"

var c *Comm

func Test_FakeIp(t *testing.T) {
	c = &Comm{}
	fmt.Println(c.FakeIp())
	fmt.Println(c.GetMasterDomain("www.baidu.com"))
}
