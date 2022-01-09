package mangahosted

import (
	"github.com/carlosrodriguesf/manga-downloader/pkg/helper"
	"github.com/carlosrodriguesf/manga-downloader/pkg/model"
	"html"
	"sort"
)

const (
	chapterLinkRegexShort = `capitulo.*?Ler\s+Online\s+-\s+(.*?)['"]\s+href=['"](.*?)['"]`
	chapterLinkRegexLarge = `<a\s+href=['"](.*?)['"]\s+title=['"]Ler\s+Online\s+-\s+(.*?)\s+\[\]`
)

type Manga struct {
	title string
	url   string
}

func NewManga(title string, url string) Manga {
	return Manga{title: title, url: url}
}

func (m Manga) Title() string {
	return m.title
}

func (m Manga) Chapters() ([]model.Chapter, error) {
	body, err := helper.GetHtmlFromURL(m.url)
	if err != nil {
		return nil, err
	}
	chapters, err := m.getChaptersFromHTML(body)
	if err != nil {
		return nil, err
	}
	sort.Slice(chapters, func(i, j int) bool {
		return i > j
	})
	return chapters, nil
}

func (m Manga) getChaptersFromHTML(body string) ([]model.Chapter, error) {
	chapters, err := m.getChaptersFromRegex(chapterLinkRegexShort, body)
	if err != nil {
		return nil, err
	}
	if len(chapters) > 0 {
		return chapters, nil
	}
	return m.getChaptersFromRegex(chapterLinkRegexLarge, body)
}

func (m Manga) getChaptersFromRegex(regex, html string) (chapters []model.Chapter, err error) {
	matches, err := helper.GetStringsSubmatchFromRegex(regex, html)
	if err != nil {
		return
	}
	chapters = m.parseChaptersSubmatch(matches)
	return
}

func (m Manga) parseChaptersSubmatch(matches [][]string) (chapters []model.Chapter) {
	for _, mt := range matches {
		l := len(mt)
		chapters = append(chapters, NewChapter(html.UnescapeString(mt[l-1]), html.UnescapeString(mt[l-2])))
	}
	return
}
