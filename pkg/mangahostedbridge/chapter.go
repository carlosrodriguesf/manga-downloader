package mangahostedbridge

import (
	"errors"
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/core"
	"strings"
	"sync"
)

const (
	imgLinkRegex1 = `img_\d+['"]\s+src=['"](.*?)['"]`
	imgLinkRegex2 = `url['"]:['"](.*?)['"]\}`
	numberRegex   = `[\w\W]+ #(\d+).*`
)

type Chapter struct {
	title           string
	url             string
	number          string
	titleSimplified string
	manga           *Manga
}

func NewChapter(manga *Manga, title, url string) Chapter {
	number := extractChapterTitleNumber(title)
	titleSimplified := fmt.Sprintf("Chapter %s", number)
	return Chapter{
		title:           title,
		url:             url,
		number:          number,
		titleSimplified: titleSimplified,
		manga:           manga,
	}
}

func (c Chapter) Manga() core.Manga {
	return *c.manga
}

func (c Chapter) Title() string {
	return c.title
}

func (c Chapter) Number() string {
	return c.number
}

func (c Chapter) TitleSimplified() string {
	return c.titleSimplified
}

func (c Chapter) Download(chapterDir string, pageDownloadCallback core.PageDownloadProgress) error {
	pages, err := c.pages()
	if err != nil {
		return err
	}

	pagesCount := len(pages)
	pagesDownloaded := 0
	pageDownloadCallback(0, 0, pagesCount)

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
			pagesDownloaded++
			pageDownloadCallback(i+1, pagesDownloaded, pagesCount)
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

func extractChapterTitleNumber(title string) string {
	str, err := core.GetStringSubmatchFromRegex(numberRegex, title)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if len(str) == 0 {
		return ""
	}
	return str[1]
}
