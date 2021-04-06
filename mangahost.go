package main

import (
	"fmt"
	"html"
	"os"
	"strings"
)

const (
	defaultHostUrl  = `https://mangahosted.com/find/`
	searchLinkRegex = `entry-title">\s*<a\s*href="(.*)?"\s*title="(.*)"`
)

type MangaHost struct {
	url string
}

func NewMangaHost(url string) MangaHost {
	return MangaHost{
		url: url,
	}
}

func GetDefaultHost() MangaHost {
	return NewMangaHost(defaultHostUrl)
}

func (mh MangaHost) Search(term string) (list []Manga, err error) {
	body, err := getHtmlFromURL(mh.url + term)
	mangas, err := getStringsSubmatchFromRegex(searchLinkRegex, body)
	if err != nil {
		return
	}
	for _, m := range mangas {
		l := len(m)
		list = append(list, NewManga(html.UnescapeString(m[l-2]), html.UnescapeString(m[l-1])))
	}
	return
}

func (mh MangaHost) SelectMangaFromList(list []Manga) Manga {
	template := "%s\n\t %d) %s"
	message := "Selecione o manga que você deseja baixar:"
	for i, m := range list {
		message = fmt.Sprintf(template, message, i+1, m.Title)
	}
	option := promptInt(message)
	listLen := len(list)
	if option > listLen || option < 1 {
		fmt.Printf("Invalid option: %d\n", option)
		os.Exit(1)
	}
	return list[option-1]
}

func (mh MangaHost) SearchAndSelect(term string) (manga Manga, err error) {
	term = strings.ReplaceAll(term, " ", "+")
	term = strings.ToLower(term)
	list, err := mh.Search(term)
	if err != nil {
		return
	}
	manga = mh.SelectMangaFromList(list)
	return
}
