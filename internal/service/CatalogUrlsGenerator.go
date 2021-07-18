package service

import (
	"regexp"
	"strconv"
)

func GenerateUrls(catalogUrlTemplate string, maxPage uint16) []string {
	re := regexp.MustCompile(`(\d+)`)
	var urls []string
	var i uint16
	for i = 1; i <= maxPage; i++ {
		urls = append(urls, re.ReplaceAllString(catalogUrlTemplate, strconv.Itoa(int(i))))
	}

	return urls
}
