package env

import (
	"fmt"

	"github.com/subosito/gotenv"
	"os"
	"os/user"
	"videocrawler/common/util"
)

var (
	ConfigPath string = ""
	HomePath   string = ""
	CookieDir  string = ""
)

func init() {
	fmt.Println("init env=================================")
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	HomePath := usr.HomeDir + "/"
	ConfigPath := HomePath + ".crawler/"

	err = util.CreateDir(ConfigPath)
	if err != nil {
		panic(err)
	}

	CookieDir := ConfigPath + "cookies/"
	err = util.CreateDir(CookieDir)
	if err != nil {
		panic(err)
	}

	envFile := ConfigPath + ".env.production"
	_, err = os.Stat(envFile)
	if err == nil {
		gotenv.Load(envFile)
	} else {
		fmt.Println("crawler config not exists")
	}
}

func LoadEnv() {

}
