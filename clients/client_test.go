package clients

import (
	"fmt"
	"testing"
)

func Test_client(t *testing.T) {
	c := NewClient("http://weibo.com", map[string]string{})
	fmt.Println(c)
}
