package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const rangeRegex = "([0-9]+)(-([0-9]+))?"

func String(message string) string {
	reader := bufio.NewReader(os.Stdin)

	template := "%s\n\n-> "
	fmt.Printf(template, message)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func Int(message string) int {
	input := String(message)
	inputInt, _ := strconv.Atoi(input)
	return inputInt
}

func Option(message string, options []string, multi bool) (int, int, error) {
	listLen := len(options)

	template := "%s\n\t %d) %s"
	for i, option := range options {
		message = fmt.Sprintf(template, message, i+1, option)
	}

	if multi {
		opts := strings.Join([]string{
			"all - To download all chapters",
			"number - To downloadChunked specifc chapter (Ex.: 941)",
			"range - To donwload a interval of chapters (Ex.: 100-400)",
		}, "\n\t")
		message = fmt.Sprintf("%s\n\n\t%s", message, opts)
	}

	if multi {
		term := String(message)
		if term == "all" {
			return 0, len(options), nil
		}

		start, end, err := parseOptionRangeTerm(term)
		if err != nil {
			return 0, 0, err
		}
		if start > listLen || start < 1 {
			return 0, 0, fmt.Errorf("Invalid option: %d\n", start)
		}
		if multi && end > listLen || end < 1 {
			return 0, 0, fmt.Errorf("Invalid option: %d\n", end)
		}
		return start, end, nil
	}

	option := Int(message)
	if option > listLen || option < 1 {
		return 0, 0, fmt.Errorf("Invalid option: %d\n", option)
	}
	return option - 1, 0, nil
}
