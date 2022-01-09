package app

import "github.com/carlosrodriguesf/manga-downloader/pkg/model"

func optionsFromMangaList(list []model.Manga) []string {
	options := make([]string, len(list))
	for i, manga := range list {
		options[i] = manga.Title()
	}
	return options
}

func optionsFromChapterList(list []model.Chapter) []string {
	options := make([]string, len(list))
	for i, chapter := range list {
		options[i] = chapter.Title()
	}
	return options
}
