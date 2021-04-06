package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type MangaChapter struct {
	Title string
	url   string
	manga *Manga
}

func NewMangaChapter(manga *Manga, title, url string) MangaChapter {
	return MangaChapter{
		Title: title,
		url:   url,
		manga: manga,
	}
}

func (mc MangaChapter) getPages() (pages []string, err error) {
	html, err := getHtmlFromURL(mc.url)
	matches, err := getStringsSubmatchFromRegex(imgLinkRegex1, html)
	if err != nil {
		return
	}
	if len(matches) == 0 {
		matches, err = getStringsSubmatchFromRegex(imgLinkRegex2, html)
		if err != nil {
			return
		}
	}
	for _, m := range matches {
		pages = append(pages, m[1])
	}
	return
}

func (mc MangaChapter) Download(mangaDir string) (err error) {
	chapterDir := mangaDir + "/" + mc.Title
	template := "\nDownloading\n\tManga: %s\n\tChapter: %s\n\tDir: %s\n"
	fmt.Printf(template, mc.manga.Title, mc.Title, chapterDir)
	pages, err := mc.getPages()
	if err != nil {
		return
	}

	totalPages := len(pages)
	fmt.Printf("\tDownloaded: 0/%d", totalPages)

	cErr := make(chan error)
	defer close(cErr)

	var errStr []string
	go (func() {
		for err := range cErr {
			errStr = append(errStr, err.Error())
		}
	})()

	var finishedPages int

	var wg sync.WaitGroup
	for _, p := range pages {
		wg.Add(1)
		go downloadPictureAsync(&wg, cErr, chapterDir, p, func() {
			finishedPages++
			fmt.Printf("\r\tDownloaded: %d/%d", finishedPages, totalPages)
		})
	}
	wg.Wait()

	// clean console
	fmt.Print("\n")

	if len(errStr) > 0 {
		err = errors.New(strings.Join(errStr, "\n"))
	}
	return
}

func (mc MangaChapter) GenerateMOBI(mangaDir string) {
	fmt.Println("Generating mobi file")
	chapterDir := mangaDir + "/" + mc.Title
	exec.Command("kcc-c2e", "-f", "MOBI", "-t", mc.Title, chapterDir)

}
