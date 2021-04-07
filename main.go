package main

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/cli"
	"github.com/carlosrodriguesf/manga-downloader/pkg/mangahostedbridge"
	"os"
)

func main() {
	host := mangahostedbridge.NewHost()
	downloader := cli.NewCli(host)
	if err := downloader.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
