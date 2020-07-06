package server

import (
	"net/url"
	"strconv"
	"strings"
	"fmt"
	"html/template"
)

type Pagination struct {
	ItemOffset int
	ItemLimit int
	ItemTotal int
	CurrentPage int
	TotalPages int

	queryStringFormat string
}

type Page struct {
	Index int
	Current bool
	Offset int
	Limit int
	QueryString template.URL

	Prev bool
	Next bool
}

const defaultQueryLimit = 24		// default number of items to return, if no limit is specified; 24 is evenly divisible by the default of 4 items displayed per row
const maxQueryLimit = 1000			// maximum number of items to return, to prevent users from doing dumb things

func queryStringFormat(o url.Values) string {
	n := url.Values{}
	for k, v := range o {
		if k != "offset" && k != "limit" {
			for _, vv := range v {
				n.Set(k, vv)
			}
		}
	}
	f := strings.Replace(n.Encode(),"%","%%", -1)
	if len(f) != 0 {
		f += "&"
	}
	f += "offset=%d&limit=%d"
	return f
}

func NewPagination(v url.Values, totalItems int) *Pagination {

	p := &Pagination{
		queryStringFormat: queryStringFormat(v),
		ItemTotal: totalItems,
	}

	p.ItemOffset, _ = strconv.Atoi(v.Get("offset"))
	if p.ItemOffset < 0 {
		p.ItemOffset = 0
	}

	p.ItemLimit, _ = strconv.Atoi(v.Get("limit"))
	if p.ItemLimit < 1 {
		p.ItemLimit = defaultQueryLimit
	}
	if p.ItemLimit > maxQueryLimit {
		p.ItemLimit = maxQueryLimit
	}

	p.TotalPages = totalItems / p.ItemLimit
	if totalItems % p.ItemLimit != 0 {
		p.TotalPages++
	}

	p.CurrentPage = (p.ItemOffset+1) / p.ItemLimit
	if (p.ItemOffset+1) % p.ItemLimit != 0 {
		p.CurrentPage++
	}

	return p
}

func (p *Pagination) Pages() []Page {
	pages := make([]Page,0,p.TotalPages)

	if p.CurrentPage != 1 {
		offset := p.ItemOffset - p.ItemLimit
		if offset < 0 {
			offset = 0
		}
		pages = append(pages,Page{
			Prev: true,
			Current: false,
			Offset: offset,
			Limit: p.ItemLimit,
			QueryString: template.URL(fmt.Sprintf(p.queryStringFormat,offset,p.ItemLimit)),
		})
	}

	for idx := 0; idx<p.TotalPages; idx++ {
		pages = append(pages,Page{
			Index: idx+1,
			Current: idx+1 == p.CurrentPage,
			Offset: idx*p.ItemLimit,
			Limit: p.ItemLimit,
			QueryString: template.URL(fmt.Sprintf(p.queryStringFormat,idx*p.ItemLimit,p.ItemLimit)),
		})
	}

	if p.CurrentPage != p.TotalPages {
		offset := p.ItemOffset + p.ItemLimit
		pages = append(pages,Page{
			Next: true,
			Current: false,
			Offset: offset,
			Limit: p.ItemLimit,
			QueryString: template.URL(fmt.Sprintf(p.queryStringFormat,offset,p.ItemLimit)),
		})
	}



	return pages
}
