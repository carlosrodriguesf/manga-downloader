package prompt

import (
	"errors"
	"github.com/carlosrodriguesf/manga-downloader/pkg/helper"
	"strconv"
)

func parseOptionRangeTerm(term string) (int, int, error) {
	matches, err := helper.GetStringSubmatchFromRegex(rangeRegex, term)
	if err != nil {
		return 0, 0, err
	}

	if len(matches) == 0 {
		return 0, 0, errors.New("invalid range")
	}

	start, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, err
	}
	if matches[3] == "" {
		matches[3] = matches[1]
	}

	end, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}
