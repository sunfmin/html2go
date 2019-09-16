package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sunfmin/html2go/parse"
)

var pkg = flag.String("pkg", "", "generated htmlgo pkg name")

func main() {
	flag.Parse()

	fmt.Println(parse.GenerateHTMLGo(*pkg, os.Stdin))
}
