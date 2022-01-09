package main

import (
	"github.com/carlosrodriguesf/manga-downloader/pkg/app"
	"github.com/carlosrodriguesf/manga-downloader/pkg/host/mangahosted"
	"log"
	"os"
)

func main() {
	host := mangahosted.New()
	err := app.New(host).Run()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
