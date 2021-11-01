package main

import (
	"log"

	"github.com/qtumproject/janus/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	cli.Run()
}
