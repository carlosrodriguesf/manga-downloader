package mangahostedbridge

import (
	"github.com/carlosrodriguesf/manga-downloader/pkg/core"
	"html"
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

func (m Manga) Chapters() (chapters []core.Chapter, err error) {
	h, err := core.GetHtmlFromURL(m.url)
	if err != nil {
		return
	}
	chapters, err = m.getChaptersFromRegex(chapterLinkRegexShort, h)
	if err != nil {
		return
	}
	if len(chapters) == 0 {
		chapters, err = m.getChaptersFromRegex(chapterLinkRegexLarge, h)
		if err != nil {
			return
		}
	}
	chapters = core.ReverseChapters(chapters)
	return
}

func (m Manga) parseChaptersSubmatch(matches [][]string) (chapters []core.Chapter) {
	for _, mt := range matches {
		l := len(mt)
		chapters = append(chapters, NewChapter(&m, html.UnescapeString(mt[l-1]), html.UnescapeString(mt[l-2])))
	}
	return
}

func (m Manga) getChaptersFromRegex(regex, html string) (chapters []core.Chapter, err error) {
	matches, err := core.GetStringsSubmatchFromRegex(regex, html)
	if err != nil {
		return
	}
	chapters = m.parseChaptersSubmatch(matches)
	return
}
