package mangahostedbridge

import (
	"github.com/carlosrodriguesf/manga-downloader/core"
	"github.com/carlosrodriguesf/manga-downloader/mangadownloader"
	"html"
)

const (
	defaultHostUrl  = `https://mangahosted.com/find/`
	searchLinkRegex = `entry-title">\s*<a\s*href="(.*)?"\s*title="(.*)"`
)

type Host struct {
}

func NewHost() Host {
	return Host{}
}

func (h Host) Search(term string) (list []mangadownloader.Manga, err error) {
	body, err := core.GetHtmlFromURL(defaultHostUrl + term)
	mangas, err := core.GetStringsSubmatchFromRegex(searchLinkRegex, body)
	if err != nil {
		return
	}
	for _, m := range mangas {
		l := len(m)
		list = append(list, NewManga(html.UnescapeString(m[l-1]), html.UnescapeString(m[l-2])))
	}
	return
}
