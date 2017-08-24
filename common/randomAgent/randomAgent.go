package randomAgent

import (
	"math/rand"
	"time"

	"github.com/buger/jsonparser"
)

var jsonContent []byte
var agents []string

func GetRandomAgent() string {
	if len(agents) == 0 {
		data, err := Asset("agent.json")
		if err != nil {
			panic(err)
		}
		jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			agents = append(agents, string(value))
		})
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	agent := agents[r.Int31n(int32(len(agents)))]
	return agent
}
