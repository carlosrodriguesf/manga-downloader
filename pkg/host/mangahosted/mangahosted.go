package mangahosted

import (
	"github.com/carlosrodriguesf/manga-downloader/pkg/helper"
	"github.com/carlosrodriguesf/manga-downloader/pkg/model"
	"strings"
)

const (
	defaultHostUrl  = `https://mangahosted.com/find/`
	searchLinkRegex = `entry-title">\s*<a\s*href="(.*)?"\s*title="(.*)"`
)

type MangaHosted struct {
}

func New() MangaHosted {
	return MangaHosted{}
}

func (h MangaHosted) Search(term string) ([]model.Manga, error) {
	term = strings.ReplaceAll(term, " ", "+")
	term = strings.ToLower(term)
	body, err := helper.GetHtmlFromURL(defaultHostUrl + term)
	mangas, err := helper.GetStringsSubmatchFromRegex(searchLinkRegex, body)
	if err != nil {
		return nil, err
	}

	mangaList := make([]model.Manga, len(mangas))
	for i, m := range mangas {
		j := len(m) - 1
		mangaList[i] = NewManga(m[j], m[j-1])
	}
	return mangaList, nil
}
