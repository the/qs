package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/the/qs/internal/process"
)

func main() {
	var in = os.Stdin
	var out = os.Stdout

	var filename string
	var output string
	var jsonproc bool
	var help bool

	flag.BoolVar(&help, "h", false, "show help")
	flag.BoolVar(&jsonproc, "json", false, "output JSON")
	flag.StringVar(&filename, "f", "", "read from file")
	flag.StringVar(&output, "o", "", "write to file")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		in = file
	}

	if output != "" {
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		out = file
	}

	if jsonproc {
		process.JSON(in, out, flag.Args())
	} else {
		process.Highlight(in, out, flag.Args())
	}
}
