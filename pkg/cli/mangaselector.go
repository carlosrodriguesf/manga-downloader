package cli

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/core"
	"os"
)

type MangaSelector struct {
	host core.Host
}

func NewMangaSelector(host core.Host) MangaSelector {
	return MangaSelector{host: host}
}

func (ms MangaSelector) Select() (manga core.Manga, err error) {
	term := core.PromptString("Search term for manga (Ex.: one piece):")

	list, err := ms.host.Search(term)
	if err != nil {
		return
	}

	manga = ms.selectFromList(list)
	return
}

func (MangaSelector) selectFromList(list []core.Manga) core.Manga {
	listLen := len(list)

	template := "%s\n\t %d) %s"
	message := "Select the manga you want to downloadChunked:"
	for i, m := range list {
		message = fmt.Sprintf(template, message, i+1, m.Title())
	}

	option := core.PromptInt(message)
	if option > listLen || option < 1 {
		fmt.Printf("Invalid option: %d\n", option)
		os.Exit(1)
	}

	return list[option-1]
}
