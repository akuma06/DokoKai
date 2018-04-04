package templates

import (
	"github.com/CloudyKit/jet"
	"net/url"
	"html/template"
	"math"
	"strconv"
)

// Navigation is used to display navigation links to pages on list view
type Navigation struct {
	TotalItem      int
	MaxItemPerPage int // FIXME: shouldn't this be in SearchForm?
	CurrentPage    int
	Route          string
}

// templateFunctions : Functions accessible in templates by {{ $.Function }}
func templateFunctions(vars jet.VarMap) jet.VarMap {
	vars.Set("genNav", genNav)
	return vars
}

func genNav(nav Navigation, currentURL *url.URL, pagesSelectable int) template.HTML {
	var ret = ""
	if nav.TotalItem > 0 {
		maxPages := math.Ceil(float64(nav.TotalItem) / float64(nav.MaxItemPerPage))

		href :=  ""
		display := " style=\"display:none;\""
		if nav.CurrentPage-1 > 0 {
			display = ""
			href = " href=\"" + "/" + nav.Route + "/1" + "?" + currentURL.RawQuery + "\""
		}
		ret = ret + "<a class=\"page-prev\"" + display + href + " aria-label=\"Previous\"><span aria-hidden=\"true\">&laquo;</span></a>"

		startValue := 1
		if nav.CurrentPage > pagesSelectable/2 {
			startValue = (int(math.Min((float64(nav.CurrentPage)+math.Floor(float64(pagesSelectable)/2)), maxPages)) - pagesSelectable + 1)
		}
		if startValue < 1 {
			startValue = 1
		}
		endValue := (startValue + pagesSelectable - 1)
		if endValue > int(maxPages) {
			endValue = int(maxPages)
		}
		for i := startValue; i <= endValue; i++ {
			pageNum := strconv.Itoa(i)
			url := "/" + nav.Route + "/" + pageNum
			ret = ret + "<a aria-label=\"Page " + strconv.Itoa(i) + "\" href=\"" + url + "?" + currentURL.RawQuery + "\">" + "<span"
			if i == nav.CurrentPage {
				ret = ret + " class=\"active\""
			}
			ret = ret + ">" + strconv.Itoa(i) + "</span></a>"
		}

		href = ""
		display = " style=\"display:none;\""
		if nav.CurrentPage < int(maxPages) {
			display = ""
			href = " href=\"" + "/" + nav.Route + "/" + strconv.Itoa(int(maxPages)) + "?" + currentURL.RawQuery + "\""
		}
		ret = ret + "<a class=\"page-next\"" + display + href +" aria-label=\"Next\"><span aria-hidden=\"true\">&raquo;</span></a>"

		itemsThisPageStart := nav.MaxItemPerPage*(nav.CurrentPage-1) + 1
		itemsThisPageEnd := nav.MaxItemPerPage * nav.CurrentPage
		if nav.TotalItem < itemsThisPageEnd {
			itemsThisPageEnd = nav.TotalItem
		}
		ret = ret + "<p>" + strconv.Itoa(itemsThisPageStart) + "-" + strconv.Itoa(itemsThisPageEnd) + "/" + strconv.Itoa(nav.TotalItem) + "</p>"
	}
	return template.HTML(ret)
}