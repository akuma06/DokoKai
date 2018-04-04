package templates

// SearchForm struct used to display the search form
type SearchForm struct {
	Category         string
	ShowItemsPerPage bool
	SearchURL string
}

// NewSearchForm return a searchForm struct with
// Some Default Values to ease things out
func NewSearchForm() SearchForm {
	return SearchForm{
		Category:         "_",
		ShowItemsPerPage: true,
		SearchURL:        "/search",
	}
}
