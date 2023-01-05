package main

import (
	"flag"

	"utm5-sber-osmp/internal/service"
)

func main() {
	fCfgPath := flag.String("c", "config.toml", "path to conf file")
	flag.Parse()
	config := LoadConfig(*fCfgPath)

	initInterrupt(service.Run(*config))
}
