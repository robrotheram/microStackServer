package main

import (
	"flag"
	"fmt"
	"microStackServer/app"
	"microStackServer/app/Background"
	"microStackServer/app/Config"
)

var (
	sourceAddr = flag.String("source", "", "source address")
	targeAddr  = flag.String("target", "", "target address")
)

func main() {
	fmt.Println("=========== STARTING ============")
	Config.ReadConfig("config.json")
	Config.LibConnect()
	Background.StartResourceCollector()
	app := &app.App{}
	app.Run()
}

/*
func main() {
	flag.Parse()
	if *sourceAddr == "" || *targeAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	http.Handle("/websockify", websocket.Server{Handshake: bootHandshake, Handler: handleWss})
	err := http.ListenAndServe(*sourceAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}*/
