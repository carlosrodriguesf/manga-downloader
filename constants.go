package main

const (
	searchUrl = `https://mangahosted.com/find/`

	chapterLinkRegexShort = `capitulo.*?Ler\s+Online\s+-\s+(.*?)['"]\s+href=['"](.*?)['"]`
	chapterLinkRegexLarge = `<a\s+href=['"](.*?)['"]\s+title=['"]Ler\s+Online\s+-\s+(.*?)\s+\[\]`
	imgLinkRegex1         = `img_\d+['"]\s+src=['"](.*?)['"]`
	imgLinkRegex2         = `url['"]:['"](.*?)['"]\}`
)
