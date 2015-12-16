package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/iambocai/mailer/g"
	"github.com/iambocai/mailer/http"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	go http.Start()

	select {}

}
