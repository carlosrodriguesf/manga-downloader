package main

import (
	"fmt"
	"os"
)

func main() {
	pwd, _ := os.Getwd()

	name := promptString("Search manga:")
	host := GetDefaultHost()
	manga, err := host.SearchAndSelect(name)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = manga.Download(pwd)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("End")
}
