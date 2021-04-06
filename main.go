package main

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/mangadownloader"
	"github.com/carlosrodriguesf/manga-downloader/pkg/mangahostedbridge"
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
