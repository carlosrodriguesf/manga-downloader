package app

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/host/mangahosted"
	"github.com/carlosrodriguesf/manga-downloader/pkg/mobiqueue"
	"github.com/carlosrodriguesf/manga-downloader/pkg/model"
	"github.com/carlosrodriguesf/manga-downloader/pkg/prompt"
	"io"
	"os"
)

const chapterRange = "([0-9]+)(-([0-9]+))?"

type App struct {
	host       model.Host
	libraryDir string

	term     string
	manga    model.Manga
	chapters []model.Chapter
}

func New(host model.Host) *App {
	return &App{
		host:       host,
		libraryDir: "/home/carlosrodrigues/Documents/Mangas",
	}
}

func (a *App) Run() error {
	//a.enterTheTerm()
	a.term = "One piece"

	//a.selectManga()
	a.manga = mangahosted.NewManga("One Piece (BR)", "https://mangahosted.com/manga/one-piece-br-mh19204")

	err := a.selectChapters()
	if err != nil {
		return err
	}

	return a.download()
}

func (a *App) enterTheTerm() {
	a.term = prompt.String("Search term for manga (Ex.: one piece):")
}

func (a *App) selectManga() error {
	mangas, err := a.host.Search("One piece")
	if err != nil {
		return err
	}

	opt, _, err := prompt.Option("Select the manga you want to download:", optionsFromMangaList(mangas), false)
	if err != nil {
		return err
	}

	fmt.Println(opt)
	return nil
}

func (a *App) selectChapters() error {
	chapters, err := a.manga.Chapters()
	if err != nil {
		return err
	}

	start, end, err := prompt.Option("Select the manga you want to downloadChunked:", optionsFromChapterList(chapters), true)
	if err != nil {
		return err
	}

	a.chapters = chapters[start-1 : end]
	return nil
}

func (a *App) download() error {
	queue := mobiqueue.New()

	mangaDir := a.mangaDir()
	for _, chapter := range a.chapters {
		path, err := a.downloadChapter(chapter, mangaDir)
		if err != nil {
			return err
		}

		queue.Add(path)
	}

	return nil
}

func (a *App) downloadChapter(chapter model.Chapter, mangaDir string) (string, error) {
	template := "\nDownloading\n\tManga: %s\n\tChapter: %s\n\tDir: %s\n"
	fmt.Printf(template, a.manga.Title(), chapter.Title(), mangaDir)

	chapterDir := fmt.Sprintf("%s/%s", mangaDir, chapter.TitleSimplified())
	if err := os.MkdirAll(chapterDir, os.ModePerm); err != nil {
		return "", err
	}

	createPage := func(name string) (io.WriteCloser, error) {
		f, err := os.Create(chapterDir + "/" + name)
		if err != nil {
			return nil, err
		}

		return f, nil
	}
	err := chapter.Download(createPage, func(downloaded, total int) {
		fmt.Printf("\r\tDownloaded: %d/%d", downloaded, total)
	})
	if err != nil {
		return "", err
	}

	fmt.Println("")
	return chapterDir, nil
}

func (a App) mangaDir() string {
	return fmt.Sprintf("%s/%s", a.libraryDir, a.manga.Title())
}
