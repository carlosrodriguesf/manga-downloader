package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func PromptString(message string) string {
	reader := bufio.NewReader(os.Stdin)

	template := "%s\n\n-> "
	fmt.Printf(template, message)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func PromptInt(message string) int {
	input := PromptString(message)
	inputInt, _ := strconv.Atoi(input)
	return inputInt
}

func GetStringsSubmatchFromRegex(strRegex, text string) (matches [][]string, err error) {
	exp, err := regexp.Compile(strRegex)
	if err != nil {
		return
	}
	matches = exp.FindAllStringSubmatch(text, -1)
	return
}

func GetStringSubmatchFromRegex(strRegex, text string) (matches []string, err error) {
	matchArray, err := GetStringsSubmatchFromRegex(strRegex, text)
	if err != nil || len(matchArray) == 0 {
		return
	}
	matches = matchArray[0]
	return
}

func GetHtmlFromURL(url string) (html string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html = string(body)
	return
}

func DownloadPicture(dir, url string) (err error) {
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}

	splitted := strings.Split(url, "/")
	fileName := splitted[len(splitted)-1]
	filePath := dir + "/" + fileName

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return
}

func DownloadPictureAsync(wg *sync.WaitGroup, cErr chan error, dir, url string, onDone func()) {
	err := DownloadPicture(dir, url)
	if err != nil {
		cErr <- err
	} else {
		onDone()
	}
	wg.Done()
}

func ReverseChapters(input []Chapter) (chapters []Chapter) {
	for i := len(input) - 1; i >= 0; i-- {
		chapters = append(chapters, input[i])
	}
	return
}

func ChapterChunkName(mangaName string, chapters []Chapter) string {
	return fmt.Sprintf("%s (%s-%s)", mangaName, chapters[0].Number(), chapters[len(chapters)-1].Number())
}
