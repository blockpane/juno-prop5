package main

import (
	"flag"
	jsb "github.com/blockpane/juno-set-bump"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&jsb.GRPCHost, "g", "127.0.0.1:9090", "grpc endpoint")
	flag.Parse()

	if jsb.GRPCHost == "" {
		flag.PrintDefaults()
		log.Fatal("grpc endpoint cannot be empty")
	}

	jsb.Run()
}
