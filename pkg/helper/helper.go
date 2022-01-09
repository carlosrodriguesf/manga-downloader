package helper

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

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
