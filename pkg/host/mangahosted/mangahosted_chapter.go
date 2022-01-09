package mangahosted

import (
	"context"
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/helper"
	"github.com/carlosrodriguesf/manga-downloader/pkg/model"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strings"
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
}

func NewChapter(title, url string) Chapter {
	number := extractChapterTitleNumber(title)
	titleSimplified := fmt.Sprintf("Chapter %s", number)
	return Chapter{
		title:           title,
		url:             url,
		number:          number,
		titleSimplified: titleSimplified,
	}
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

func (c Chapter) pages() (pages []string, err error) {
	html, err := helper.GetHtmlFromURL(c.url)
	matches, err := helper.GetStringsSubmatchFromRegex(imgLinkRegex1, html)
	if err != nil {
		return
	}
	if len(matches) == 0 {
		matches, err = helper.GetStringsSubmatchFromRegex(imgLinkRegex2, html)
		if err != nil {
			return
		}
	}
	for _, m := range matches {
		pages = append(pages, m[1])
	}
	return
}

func (c Chapter) Download(createPage model.CreatePage, callback model.PageDownloaded) error {
	pages, err := c.pages()
	if err != nil {
		return err
	}

	pagesCount := len(pages)
	callback(0, pagesCount)

	pagesDownloaded := 0
	errs, ctx := errgroup.WithContext(context.Background())
	for _, page := range pages {
		page := page
		errs.Go(func() error {
			err := downloadPage(ctx, page, createPage)
			if err != nil {
				return err
			}

			pagesDownloaded++
			callback(pagesDownloaded, pagesCount)
			return nil
		})
	}

	return errs.Wait()
}

func downloadPage(ctx context.Context, url string, createPage model.CreatePage) error {
	filename := filenameFromURL(url)
	dest, err := createPage(filename)
	if err != nil {
		return err
	}
	defer dest.Close()

	source, err := requestFile(ctx, url)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return err
	}

	return err
}

func extractChapterTitleNumber(title string) string {
	str, err := helper.GetStringSubmatchFromRegex(numberRegex, title)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if len(str) == 0 {
		return ""
	}
	return str[1]
}

func filenameFromURL(url string) string {
	splitted := strings.Split(url, "/")
	fileName := splitted[len(splitted)-1]
	return fileName
}

func requestFile(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
