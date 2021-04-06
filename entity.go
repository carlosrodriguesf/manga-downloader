package main

type Manga struct {
	Title string
	Url   string
}

type Chapter struct {
	Title string
	Url   string
}

type Host interface {
	Search(term string) ([]Manga, error)
	GetChapters(manga Manga) ([]Chapter, error)
	GetChapterPages(chapter Chapter) ([]string, error)
}
