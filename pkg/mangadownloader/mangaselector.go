package mangadownloader

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/core"
	"os"
	"strings"
)

type MangaSelector struct {
	host Host
}

func NewMangaSelector(host Host) MangaSelector {
	return MangaSelector{host: host}
}

func (ms MangaSelector) Select() (manga Manga, err error) {
	term := ms.getTerm()

	list, err := ms.host.Search(term)
	if err != nil {
		return
	}

	manga = ms.selectFromList(list)
	return
}

func (MangaSelector) getTerm() string {
	term := core.PromptString("Search term for manga (Ex.: one piece):")
	term = strings.ReplaceAll(term, " ", "+")
	term = strings.ToLower(term)
	return term
}

func (MangaSelector) selectFromList(list []Manga) Manga {
	listLen := len(list)

	template := "%s\n\t %d) %s"
	message := "Select the manga you want to download:"
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
