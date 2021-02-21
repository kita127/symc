package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kita127/symc"
)

func main() {
	input, _ := ioutil.ReadAll(os.Stdin)
	module := symc.ParseModule(string(input))
	fmt.Println(module.PrettyString())
	for _, s := range module.Statements {
		if i, ok := s.(*symc.InvalidStatement); ok {
			fmt.Printf("remain:%v\n", i.Remain[:10])
		}
	}
}
