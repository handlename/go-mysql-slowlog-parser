package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/handlename/go-mysql-slowlog-parser"
)

func main() {
	opts := parseOpts()
	parser := slowlog.NewParser()

	loc, err := time.LoadLocation(opts.location)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	parser.Location = loc

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
		return nil, fmt.Errorf("invalid style\n")
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
	style    string
	location string
}

func parseOpts() opts {
	o := opts{}

	flag.StringVar(&o.style, "style", "ltsv", "output style. \"ltsv\" or \"json\"")
	flag.StringVar(&o.location, "location", "UTC", "location of slowlog's datetime")
	flag.Parse()

	return o
}
