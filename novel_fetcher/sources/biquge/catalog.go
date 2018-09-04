package biquge

import (
	"context"
	"errors"
	"fmt"
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
		uri, exist := sectionNode.Attr("href")
		if !exist {
			// log
			return
		}
		c.sections = append(c.sections, NewFromURL(ctx, html, uri, c.bookName, c, i))
	})
	for i, j := startIndex, startIndex+1; j < len(c.sections); i, j = j, j+1 {
		c.sections[i].SetNext(ctx, c.sections[j])
		c.sections[j].SetPre(ctx, c.sections[i])
	}
	return nil
}

// Save ...
func (c *Catalog) Save(ctx context.Context, start, end int, source base.NovelSource) error {
	remaining := c.sections[start:end]
	fails := []int{}
	for i := 0; i < base.SaveRetry; i++ {
		fails = source.Save(ctx, c.bookName, remaining)
		if len(fails) == 0 {
			break
		}
		remaining = make([]base.Sectioner, 0, len(fails))
		for _, fail := range fails {
			remaining = append(remaining, c.sections[fail])
		}
	}

	if len(fails) != 0 {
		return fmt.Errorf("sections not saved: %v", fails)

	}

	return nil
}

// LoadFromSource 必须在Flush后使用
func (c *Catalog) LoadFromSource(ctx context.Context, source base.NovelSource) error {
	bodies, err := source.Load(ctx, c.bookName)
	if err != nil {
		return err
	}
	for _, body := range bodies {
		section, err := NewFromBody(ctx, body)
		if err != nil {
			// NOTE log
			continue
		}
		c.UpdateSection(ctx, section)
	}
	return nil
}

// UpdateSection ...
func (c *Catalog) UpdateSection(ctx context.Context, section base.Sectioner) error {
	index := section.GetIndex(ctx)
	if index >= len(c.sections) {
		return errors.New("out of section range")
	}

	section.SetCatalog(ctx, c)
	if index > 0 {
		section.SetPre(ctx, c.sections[index-1])
		c.sections[index-1].SetNext(ctx, section)
	}
	if index < len(c.sections)-1 {
		section.SetNext(ctx, c.sections[index+1])
		c.sections[index+1].SetPre(ctx, section)
	}
	return nil
}

// NewCatalog ...
func NewCatalog(ctx context.Context, req *http.Request, bookName string) (base.Cataloger, error) {
	c := &Catalog{
		request:  req,
		bookName: bookName,
		client:   &http.Client{},
	}
	err := c.Flush(ctx, false)
	return c, err
}
