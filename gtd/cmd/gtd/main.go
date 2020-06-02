package main

import (
	"flag"
	"log"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Usage = Usage
	flag.Parse()

	REPL(flag.Args())
}
