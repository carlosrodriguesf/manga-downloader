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
	return fmt.Sprintf("%s/%s", d.libraryDir, d.manga.Title())
}

func (d Cli) download() error {
	queue := core.NewMobiGeneratorQueue()

	for _, c := range d.chapters {
		mangaDir := d.mangaDir()
		chapterDir := fmt.Sprintf("%s/%s", mangaDir, c.TitleSimplified())
		if err := d.downloadChapter(chapterDir, c); err != nil {
			fmt.Println("Waiting mobi files generation finish.")
			queue.Wait()
			return err
		}
		queue.Add(chapterDir)
	}

	fmt.Println("Waiting mobi files generation finish.")
	queue.Wait()
	return queue.Err()
}

func (d Cli) downloadChunked() error {
	queue := core.NewMobiGeneratorQueue()

	chunkSize := 30
	total := len(d.chapters)

	for start := 0; start < total; start += chunkSize {
		end := start + chunkSize
		if end > total {
			end = total
		}

		mangaDir := d.mangaDir()
		chunkDir, err := d.downloadChunk(mangaDir, start, end)

		if err != nil {
			fmt.Println("Waiting mobi files generation finish.")
			queue.Wait()
			return err
		}
		queue.Add(chunkDir)
	}

	fmt.Println("Waiting mobi files generation finish.")
	queue.Wait()
	return queue.Err()
}

func (d Cli) downloadChunk(mangaDir string, start, end int) (string, error) {
	chapters := d.chapters[start:end]
	chunkName := core.ChapterChunkName(d.manga.Title(), chapters)
	chunkDir := fmt.Sprintf("%s/%s", mangaDir, chunkName)
	for _, c := range chapters {
		chapterDir := fmt.Sprintf("%s/%s", chunkDir, c.TitleSimplified())
		err := d.downloadChapter(chapterDir, c)
		if err != nil {
			return chunkDir, err
		}
	}
	return chunkDir, nil
}

func (d Cli) downloadChapter(chapterDir string, c core.Chapter) error {
	template := "\nDownloading\n\tManga: %s\n\tChapter: %s\n\tDir: %s\n"
	fmt.Printf(template, d.manga.Title(), c.Title(), chapterDir)

	fmt.Printf("\tDownloaded: 0/0")
	err := c.Download(chapterDir, func(_, downloaded, total int) {
		fmt.Printf("\r\tDownloaded: %d/%d", downloaded, total)
	})

	fmt.Print("\n")
	return err
}
