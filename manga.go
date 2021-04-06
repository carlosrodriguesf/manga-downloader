package main

import (
	"errors"
	"fmt"
	"html"
	"strconv"
)

const chapterRange = "([0-9]+)(-([0-9]+))?"

type Manga struct {
	Title    string
	url      string
	chapters []MangaChapter
}

func NewManga(url, title string) Manga {
	return Manga{
		url:   url,
		Title: title,
	}
}

func (m Manga) parseChaptersSubmatch(matches [][]string) (chapters []MangaChapter) {
	for _, mt := range matches {
		l := len(mt)
		chapters = append(chapters, NewMangaChapter(&m, html.UnescapeString(mt[l-1]), html.UnescapeString(mt[l-2])))
	}
	return
}

func (m Manga) getChaptersFromRegex(regex, html string) (chapters []MangaChapter, err error) {
	matches, err := getStringsSubmatchFromRegex(regex, html)
	if err != nil {
		return
	}
	chapters = m.parseChaptersSubmatch(matches)
	return
}

func (m *Manga) loadChapters() (err error) {
	html, err := getHtmlFromURL(m.url)
	if err != nil {
		return
	}

	chapters, err := m.getChaptersFromRegex(chapterLinkRegexShort, html)
	if err != nil {
		return
	}
	if len(chapters) == 0 {
		chapters, err = m.getChaptersFromRegex(chapterLinkRegexLarge, html)
		if err != nil {
			return
		}
	}

	m.chapters = chapters
	return
}

func (m Manga) parseDownloadRangeTerm(term string) (start, end int, err error) {
	matches, err := getStringSubmatchFromRegex(chapterRange, term)
	if err != nil {
		return
	}

	if len(matches) == 0 {
		err = errors.New("Invalid range")
		return
	}

	start, err = strconv.Atoi(matches[1])
	if err != nil {
		return
	}
	if matches[3] == "" {
		matches[3] = matches[1]
	}
	end, err = strconv.Atoi(matches[3])
	return
}

func (m Manga) reversedChapters() (chapters []MangaChapter) {
	for i := len(m.chapters) - 1; i >= 0; i-- {
		chapters = append(chapters, m.chapters[i])
	}
	return
}

func (m Manga) download(mangaDir string, start, end int) (err error) {
	chapters := m.reversedChapters()[start-1 : end]
	for _, c := range chapters {
		err = c.Download(mangaDir)
		if err != nil {
			return err
		}
	}
	return
}

func (m Manga) Download(dir string) (err error) {
	mangaDir := dir + "/" + m.Title
	if len(m.chapters) == 0 {
		err = m.loadChapters()
		if err != nil {
			return
		}
	}

	template := "%s\n\t%d | %s"
	message := ""
	for i, c := range m.chapters {
		message = fmt.Sprintf(template, message, i+1, c.Title)
	}
	message += "\nSelecione os capítulos:"
	rangeTerm := promptString(message)
	if rangeTerm == "all" {
		err = m.download(mangaDir, 1, len(m.chapters))
		return
	}

	start, end, err := m.parseDownloadRangeTerm(rangeTerm)
	if err != nil {
		return
	}
	err = m.download(mangaDir, start, end)
	return
}
