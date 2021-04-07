package cli

import (
	"errors"
	"fmt"
	"github.com/carlosrodriguesf/manga-downloader/pkg/core"
	"strconv"
)

const chapterRange = "([0-9]+)(-([0-9]+))?"

type ChapterSelector struct {
	manga core.Manga
}

func NewChapterSelector(manga core.Manga) ChapterSelector {
	return ChapterSelector{manga: manga}
}

func (cs ChapterSelector) Select() (chapters []core.Chapter, err error) {
	chapters, err = cs.manga.Chapters()

	template := "%s\n\t%d | %s"
	message := ""
	for i, c := range chapters {
		message = fmt.Sprintf(template, message, i+1, c.Title())
	}
	message += "\n\nSelect chapters: \n\n\tall - To download all chapters\n\tnumber - To download specifc chapter (Ex.: 941)\n\trange - To donwload a interval of chapters (Ex.: 100-400)"
	rangeTerm := core.PromptString(message)
	if rangeTerm == "all" {
		return
	}
	start, end, err := parseDownloadRangeTerm(rangeTerm)
	if err != nil {
		return
	}
	chapters = chapters[start-1 : end]
	return
}

func parseDownloadRangeTerm(term string) (start, end int, err error) {
	matches, err := core.GetStringSubmatchFromRegex(chapterRange, term)
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
