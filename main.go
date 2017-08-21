package main

import (
	"flag"

	"github.com/swhsiang/gone/server"
)

func main() {
	host := flag.String("Host", "localhost", "Host")
	port := flag.String("Port", "8020", "Port")
	flag.Parse()

	s := server.NewServer(*host, *port)
	s.Run()
}
