package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/handlename/go-mysql-slowlog-parser"
)

func main() {
	opts := parseOpts()
	parser := slowlog.NewParser()

	stylist, err := selectStylist(opts.style)

	if err != nil {
		log.Panic(err)
	}

	for parsed := range parser.Parse(os.Stdin) {
		fmt.Println(stylist(parsed))
	}
}

type stylist func(slowlog.Parsed) string

func selectStylist(style string) (s stylist, err error) {
	switch style {
	default:
		return nil, errors.New("invalid style")
	case "ltsv":
		return func(p slowlog.Parsed) string {
			return p.AsLTSV()
		}, nil
	case "json":
		return func(p slowlog.Parsed) string {
			return p.AsJSON()
		}, nil
	}
}

type opts struct {
	style string
}

func parseOpts() opts {
	var (
		style string
	)

	flag.StringVar(&style, "style", "ltsv", "output style. \"ltsv\" or \"json\"")
	flag.Parse()

	return opts{style: style}
}
