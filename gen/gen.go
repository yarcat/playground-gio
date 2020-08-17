package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var img = flag.String("img", "", "image  to encode")

const template = `package res

const %sPNG = %q`

func main() {
	flag.Parse()
	if *img == "" {
		fmt.Println("-img is missing")
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(*img)
	if err != nil {
		fmt.Println("unable to read file:", err)
		os.Exit(1)
	}

	name := strings.SplitN(path.Base(*img), ".", 2)[0]
	fmt.Printf(template, strings.Title(strings.ToLower(name)), b)
}
