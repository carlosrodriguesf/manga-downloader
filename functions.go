package main

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

func promptString(message string) string {
	reader := bufio.NewReader(os.Stdin)

	template := "%s\n\n-> "
	fmt.Printf(template, message)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func promptInt(message string) int {
	input := promptString(message)
	inputInt, _ := strconv.Atoi(input)
	return inputInt
}

func getStringsSubmatchFromRegex(strRegex, text string) (matches [][]string, err error) {
	exp, err := regexp.Compile(strRegex)
	if err != nil {
		return
	}
	matches = exp.FindAllStringSubmatch(text, -1)
	return
}

func getStringSubmatchFromRegex(strRegex, text string) (matches []string, err error) {
	matchArray, err := getStringsSubmatchFromRegex(strRegex, text)
	if err != nil || len(matchArray) == 0 {
		return
	}
	matches = matchArray[0]
	return
}

func getHtmlFromURL(url string) (html string, err error) {
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

func downloadPicture(dir, url string) (err error) {
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

func downloadPictureAsync(wg *sync.WaitGroup, cErr chan error, dir, url string, onDone func()) {
	err := downloadPicture(dir, url)
	if err != nil {
		cErr <- err
	} else {
		onDone()
	}
	wg.Done()
}

func display(str string) {
	fmt.Printf("\033[0;0H")
	fmt.Println(str)
}
