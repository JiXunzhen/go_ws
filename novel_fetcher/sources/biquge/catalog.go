package biquge

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/utils"
	"github.com/PuerkitoBio/goquery"
)

const (
	// LoadConcurrent 加载并发数
	LoadConcurrent = 5
)

// Catalog ...
type Catalog struct {
	request  *http.Request
	sections map[int]base.Sectioner
	bookName string
	client   *http.Client
	newest   int
}

// Count ...
func (c *Catalog) Count() int {
	return c.newest
}

// GetBookName ...
func (c *Catalog) GetBookName() string {
	return c.bookName
}

// List ...
func (c *Catalog) List() []base.Sectioner {
	sections := make([]base.Sectioner, 0, c.newest)
	for i := 0; i < c.newest; i++ {
		if s, ok := c.sections[i]; ok {
			sections = append(sections, s)
		}
	}
	return sections
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
		c.sections = map[int]base.Sectioner{}
	}

	doc.Find("#list").Children().Children().Each(func(i int, s *goquery.Selection) {
		if _, ok := c.sections[i]; ok {
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
		c.sections[i] = NewFromURL(ctx, html, uri, c.bookName, c, i)
		c.newest = i
	})
	c.flushLink(ctx)
	return nil
}

// Save ...
func (c *Catalog) Save(ctx context.Context, start, end int, source base.NovelSource) error {
	if end <= start {
		return nil
	}

	fails := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		fails = append(fails, i)
	}

	for i := 0; i < base.SaveRetry; i++ {
		remain := make(map[int]base.Sectioner, len(fails))
		for _, fail := range fails {
			remain[fail] = c.sections[fail]
		}
		fails = source.Save(ctx, c.bookName, remain)
		if len(fails) == 0 {
			break
		}
	}

	if len(fails) != 0 {
		return fmt.Errorf("sections not saved: %v", fails)
	}

	return nil
}

// Load ...
func (c *Catalog) Load(ctx context.Context, start, end int) error {
	fails := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		fails = append(fails, i)
	}
	for i := 0; i < base.LoadRetry; i++ {
		remain := make(map[int]base.Sectioner, len(fails))
		for _, fail := range fails {
			remain[fail] = c.sections[i]
		}
		fails = c.stepLoad(ctx, remain)
		if len(fails) == 0 {
			break
		}
	}
	if len(fails) > 0 {
		return fmt.Errorf("sections not loaded: %v", fails)
	}
	return nil
}

func (c *Catalog) stepLoad(ctx context.Context, sections map[int]base.Sectioner) (fails []int) {
	sectionChan := make(chan base.Sectioner)
	failChan := make(chan int)

	finished := &sync.WaitGroup{}
	finished.Add(LoadConcurrent)

	// go for work
	for i := 0; i < LoadConcurrent; i++ {
		go func(i int) {
			defer finished.Done()

			for section := range sectionChan {
				err := section.Load(ctx)
				if err != nil {
					// NOTE log err
					failChan <- section.GetIndex()
				}

				// sleep for 0.5s
				time.Sleep(500 * time.Millisecond)
			}
		}(i)
	}
	// fail collect
	go func() {
		for index := range failChan {
			fails = append(fails, index)
		}
	}()

	// send works
	for _, section := range sections {
		sectionChan <- section
	}
	close(sectionChan)

	// wait for work
	finished.Wait()

	// handle fails
	close(failChan)
	return
}

// LoadFromSource ...
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
	c.flushLink(ctx)
	return nil
}

// UpdateSection ...
func (c *Catalog) UpdateSection(ctx context.Context, section base.Sectioner) error {
	index := section.GetIndex()
	if index >= len(c.sections) {
		return errors.New("out of section range")
	}

	section.SetCatalog(ctx, c)
	if old, ok := c.sections[index]; ok {
		section.SetPre(ctx, old.GetPre())
		section.SetNext(ctx, old.GetNext())
	}
	if index > c.newest {
		c.newest = index
	}
	return nil
}

func (c *Catalog) flushLink(ctx context.Context) {
	var pre base.Sectioner
	for i := 0; i < len(c.sections); i++ {
		cur, ok := c.sections[i]
		if !ok {
			continue
		}
		if pre != nil {
			pre.SetNext(ctx, cur)
			cur.SetPre(ctx, pre)
		}
		pre = cur
	}
}

// NewCatalog ...
func NewCatalog(ctx context.Context, req *http.Request, bookName string) (base.Cataloger, error) {
	c := &Catalog{
		request:  req,
		bookName: bookName,
		client:   &http.Client{},
		sections: make(map[int]base.Sectioner),
	}
	err := c.Flush(ctx, false)
	return c, err
}
