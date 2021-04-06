package mangahostedbridge

import (
	"errors"
	"github.com/carlosrodriguesf/manga-downloader/core"
	"github.com/carlosrodriguesf/manga-downloader/mangadownloader"
	"strings"
	"sync"
)

const (
	imgLinkRegex1 = `img_\d+['"]\s+src=['"](.*?)['"]`
	imgLinkRegex2 = `url['"]:['"](.*?)['"]\}`
)

type Chapter struct {
	title string
	url   string
	manga *Manga
}

func NewChapter(manga *Manga, title, url string) Chapter {
	return Chapter{
		title: title,
		url:   url,
		manga: manga,
	}
}

func (c Chapter) Manga() mangadownloader.Manga {
	return *c.manga
}

func (c Chapter) Title() string {
	return c.title
}

func (c Chapter) Download(chapterDir string, pageDownloadCallback mangadownloader.PageDownloadProgress) error {
	pages, err := c.pages()
	if err != nil {
		return err
	}

	pagesCount := len(pages)
	pageDownloadCallback(0, pagesCount)

	cErr := make(chan error)
	defer close(cErr)

	var errStr []string
	go (func() {
		for err := range cErr {
			errStr = append(errStr, err.Error())
		}
	})()

	var wg sync.WaitGroup
	for i, p := range pages {
		i := i
		wg.Add(1)
		go core.DownloadPictureAsync(&wg, cErr, chapterDir, p, func() {
			pageDownloadCallback(i+1, pagesCount)
		})
	}
	wg.Wait()

	if len(errStr) > 0 {
		return errors.New(strings.Join(errStr, "\n"))
	}
	return nil
}

func (c Chapter) pages() (pages []string, err error) {
	html, err := core.GetHtmlFromURL(c.url)
	matches, err := core.GetStringsSubmatchFromRegex(imgLinkRegex1, html)
	if err != nil {
		return
	}
	if len(matches) == 0 {
		matches, err = core.GetStringsSubmatchFromRegex(imgLinkRegex2, html)
		if err != nil {
			return
		}
	}
	for _, m := range matches {
		pages = append(pages, m[1])
	}
	return
}
