package domain

type SiteConfig struct {
	Id                      uint16
	Name                    string
	CatalogUrls             []string 		  `db:"catalog_urls"`
	Headers                 map[string]string
	PageCountSearchInterval uint16 			  `db:"page_count_search_interval"`
	ParserMaxExecutionTime  uint16			  `db:"parser_max_execution_time"`
	Cookies                 string
	Data                    map[string]string
	LastSearchPageCount		uint16			  `db:"last_search_page_count"`
}
