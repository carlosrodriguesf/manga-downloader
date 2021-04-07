package core

type PageDownloadProgress func(page, downloaded, total int)

type Manga interface {
	Title() string
	Chapters() ([]Chapter, error)
}

type Chapter interface {
	Title() string
	Manga() Manga
	Download(chapterDir string, event PageDownloadProgress) error
}

type Host interface {
	Search(term string) ([]Manga, error)
}
