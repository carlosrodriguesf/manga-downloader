package model

import "io"

type PageDownloaded func(downloaded, total int)
type CreatePage func(name string) (io.WriteCloser, error)

type Host interface {
	Search(term string) ([]Manga, error)
}

type Manga interface {
	Title() string
	Chapters() ([]Chapter, error)
}

type Chapter interface {
	Title() string
	TitleSimplified() string
	Number() string
	Download(createPage CreatePage, callback PageDownloaded) error
}
