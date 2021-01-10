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
}
