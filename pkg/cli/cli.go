package cli

import (
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/core"
	"os"
)

type Cli struct {
	host       core.Host
	manga      core.Manga
	chapters   []core.Chapter
	libraryDir string
}

func NewCli(host core.Host) Cli {
	pwd, _ := os.Getwd()
	return Cli{host: host, libraryDir: pwd}
}

func (d *Cli) Run() error {
	if err := d.selectManga(); err != nil {
		return err
	}
	if err := d.selectChapters(); err != nil {
		return err
	}
	return d.download()
}

func (d *Cli) selectManga() (err error) {
	d.manga, err = NewMangaSelector(d.host).Select()
	return
}

func (d *Cli) selectChapters() (err error) {
	d.chapters, err = NewChapterSelector(d.manga).Select()
	return
}

func (d Cli) mangaDir() string {
	return d.libraryDir + "/" + d.manga.Title()
}

func (d Cli) download() error {
	queue := core.NewMobiGeneratorQueue()
	for _, c := range d.chapters {
		chapterDir, err := d.downloadChapter(c)
		if err != nil {
			return err
		}
		queue.Add(chapterDir)
	}

	fmt.Println("Waiting mobi files generation finish.")
	queue.Wait()
	return queue.Err()
}

func (d Cli) downloadChapter(c core.Chapter) (string, error) {
	chapterDir := d.mangaDir() + "/" + c.Title()
	template := "\nDownloading\n\tManga: %s\n\tChapter: %s\n\tDir: %s\n"
	fmt.Printf(template, d.manga.Title(), c.Title(), chapterDir)

	fmt.Printf("\tDownloaded: 0/0")
	err := c.Download(chapterDir, func(_, downloaded, total int) {
		fmt.Printf("\r\tDownloaded: %d/%d", downloaded, total)
	})

	fmt.Print("\n")
	return chapterDir, err
}
