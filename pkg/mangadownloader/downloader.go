package mangadownloader

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Downloader struct {
	host       Host
	manga      Manga
	chapters   []Chapter
	libraryDir string
}

func NewDownloader(host Host) Downloader {
	pwd, _ := os.Getwd()
	return Downloader{host: host, libraryDir: pwd}
}

func (d *Downloader) Run() error {
	if err := d.selectManga(); err != nil {
		return err
	}
	if err := d.selectChapters(); err != nil {
		return err
	}
	return d.download()
}

func (d *Downloader) selectManga() (err error) {
	d.manga, err = NewMangaSelector(d.host).Select()
	return
}

func (d *Downloader) selectChapters() (err error) {
	d.chapters, err = NewChapterSelector(d.manga).Select()
	return
}

func (d Downloader) mangaDir() string {
	return d.libraryDir + "/" + d.manga.Title()
}

func (d Downloader) download() error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	cErr := make(chan error)
	defer close(cErr)

	var mobiGenerationErrors []string
	go (func() {
		for err := range cErr {
			mobiGenerationErrors = append(mobiGenerationErrors, err.Error())
		}
	})()

	for _, c := range d.chapters {
		chapterDir, err := d.downloadChapter(c)
		if err != nil {
			return err
		}

		wg.Add(1)
		go d.generateMOBI(&wg, &mu, cErr, chapterDir)
	}

	fmt.Println("Waiting mobi files generation finish.")
	wg.Wait()

	if len(mobiGenerationErrors) > 0 {
		return errors.New(strings.Join(mobiGenerationErrors, "\n\n"))
	}
	return nil
}

func (d Downloader) downloadChapter(c Chapter) (string, error) {
	chapterDir := d.mangaDir() + "/" + c.Title()
	template := "\nDownloading\n\tManga: %s\n\tChapter: %s\n\tDir: %s\n"
	fmt.Printf(template, d.manga.Title(), c.Title(), chapterDir)

	fmt.Printf("\tDownloaded: 0/0")
	finished := 0
	err := c.Download(chapterDir, func(_, total int) {
		finished++
		fmt.Printf("\r\tDownloaded: %d/%d", finished, total)
	})

	fmt.Print("\n")
	return chapterDir, err
}

func (d Downloader) generateMOBI(wg *sync.WaitGroup, mu *sync.Mutex, cErr chan error, chapterDir string) {
	mu.Lock()

	cmdOutput := &bytes.Buffer{}
	cmd := exec.Command("kcc-c2e", "-f", "MOBI", chapterDir)
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		template := "%s\n\t%s\n\t%s\n\n======%s====="
		cErr <- fmt.Errorf(template, chapterDir, err.Error(), fmt.Sprint("kcc-c2e", "-f", "MOBI", chapterDir), string(cmdOutput.Bytes()))
	}

	mu.Unlock()
	wg.Done()
}
