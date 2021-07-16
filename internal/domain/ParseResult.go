package domain

type PageResponse struct {
	PageUrl string
	Data    map[string]string
}

type ParseResult struct {
	SiteId            uint16
	PagesParseResults []PageResponse
}
