package main

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/mangadownloader"
	"github.com/carlosrodriguesf/manga-downloader/mangahostedbridge"
	"os"
)

func main() {
	host := mangahostedbridge.NewHost()
	downloader := mangadownloader.NewDownloader(host)
	if err := downloader.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
