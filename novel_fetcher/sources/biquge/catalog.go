package biquge

import (
	"context"
	"errors"
	"net/http"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/utils"
	"github.com/PuerkitoBio/goquery"
)

// Catalog ...
type Catalog struct {
	request  *http.Request
	sections []base.Sectioner
	bookName string
	client   *http.Client
	preload  int
	loaded   int
}

// Count ...
func (c *Catalog) Count(_ context.Context) int {
	return len(c.sections)
}

// Get ...
func (c *Catalog) Get(ctx context.Context, index int) (base.Sectioner, error) {
	if len(c.sections) <= index {
		return nil, errors.New("out of sections range")
	}
	return c.sections[index], nil
}

// Flush ...
func (c *Catalog) Flush(ctx context.Context, withoutCache bool) error {
	// post
	resp, err := c.client.Do(c.request)
	if err != nil {
		return err
	}

	// get doc
	doc, err := goquery.NewDocumentFromResponse(resp)
	// 转字符
	html, err := doc.Html()
	if err != nil {
		return err
	}
	doc.SetHtml(utils.Convert(html, utils.GbkToUtf8))

	// refresh section
	if withoutCache {
		c.sections = []base.Sectioner{}
	}
	startIndex := len(c.sections)

	doc.Find("#list").Children().Children().Each(func(i int, s *goquery.Selection) {
		if i < startIndex {
			return
		}
		if s.Children() == nil {
			return
		}
		sectionNode := s.Children()
		html, err := sectionNode.Html()
		if err != nil {
			// log
			return
		}
		url, exist := sectionNode.Attr("href")
		if !exist {
			// log
			return
		}
		c.sections = append(c.sections, &Section{
			name:     html,
			url:      url,
			bookName: c.bookName,
			catalog:  c,
			index:    i,
		})
	})
	for i, j := startIndex, startIndex+1; j < len(c.sections); i, j = j, j+1 {
		c.sections[i].SetNext(ctx, c.sections[j])
		c.sections[j].SetPre(ctx, c.sections[i])
	}
	return nil
}

// Save ...
func (c *Catalog) Save(ctx context.Context, start, end int) error {
	return nil
}

// SetPreLoad ...
func (c *Catalog) SetPreLoad(_ context.Context, preload int) {
	c.preload = preload
}

// NewCatalog ...
func NewCatalog(ctx context.Context, req *http.Request, bookName string) (base.Cataloger, error) {
	c := &Catalog{
		request:  req,
		bookName: bookName,
		client:   &http.Client{},
		preload:  base.DefaultPreload,
		loaded:   0,
	}
	err := c.Flush(ctx, false)
	return c, err
}
