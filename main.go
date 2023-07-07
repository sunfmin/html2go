package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sunfmin/html2go/parse"
)

var pkg = flag.String("pkg", "", "generated htmlgo pkg name")
var childrenMode = flag.Bool("c", false, "children mode")

func main() {
	flag.Parse()

	fmt.Println(parse.GenerateHTMLGo(*pkg, *childrenMode, os.Stdin))
}
