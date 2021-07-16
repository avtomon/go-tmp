package service

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"parser/internal/domain"
	"time"
)

const ToManyRequestsResponseCode = 429

func ParseSite(urls *[]string, siteConfig *domain.SiteConfig) (*domain.ParseResult, *[]error) {
	var (
		errs []error
		maxRequestDelay = GetMaxRequestDelay(siteConfig.ParserMaxExecutionTime, uint16(len(*urls)))
		result = domain.ParseResult{SiteId: siteConfig.Id, PagesParseResults: []domain.PageResponse{}}
	)

	for _, catalogUrl := range *urls {
		parseResult, responseCode, err := parse(catalogUrl, siteConfig)
		if responseCode == ToManyRequestsResponseCode {
			*urls = append(*urls, catalogUrl)
		}

		errs = append(errs, err)
		result.PagesParseResults = append(result.PagesParseResults, domain.PageResponse{PageUrl: catalogUrl, Data: *parseResult})

		time.Sleep(time.Duration(GetRandomDelay(maxRequestDelay) * 1000) * time.Millisecond)
	}

	return &result, &errs
}

func parse(catalogUrl string, siteConfig *domain.SiteConfig) (*map[string]string, int, error) {
	var result map[string]string
	client := &http.Client{}
	request, err := http.NewRequest("GET", catalogUrl, nil)
	statusCode := 0
	if err != nil {
		return &result, statusCode, err
	}

	for name, value := range siteConfig.Headers {
		request.Header.Add(name, value)
	}

	request.Header.Add("Set-Cookie", siteConfig.Cookies)

	response, err := client.Do(request)

	if response != nil {
		statusCode = response.StatusCode
	}

	if err != nil {
		return &result, statusCode, err
	}

	if response.ContentLength == 0 {
		err = errors.New("response is empty")
		return &result, statusCode, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("response body reader close error: %v", err.Error())
		}
	}(response.Body)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return &result, statusCode, err
	}

	for name, selector := range siteConfig.Data {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			result[name] = s.Text()
		})
	}

	return &result, statusCode, nil
}
