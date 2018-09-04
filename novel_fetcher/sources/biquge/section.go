package biquge

import (
	"context"
	"fmt"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/PuerkitoBio/goquery"
)

// Section ...
type Section struct {
	url      string
	name     string
	bookName string
	catalog  base.Cataloger
	index    int
	pre      base.Sectioner
	next     base.Sectioner
	doc      *goquery.Document
	text     string
}

// Text ...
func (s *Section) Text(context.Context) (string, error) {
	if s.doc != nil {
		return s.text, nil
	}
	doc, err := goquery.NewDocument(s.url)
	if err != nil {
		return "", err
	}
	s.doc = doc
	fmt.Println(doc.Html())
	return "", nil
}

// GetPre ...
func (s *Section) GetPre(context.Context) (base.Sectioner, error) {
	return s.pre, nil
}

// GetNext ...
func (s *Section) GetNext(context.Context) (base.Sectioner, error) {
	return s.next, nil
}

// SetPre ...
func (s *Section) SetPre(ctx context.Context, pre base.Sectioner) {
	s.pre = pre
}

// SetNext ...
func (s *Section) SetNext(ctx context.Context, next base.Sectioner) {
	s.next = next
}

// GetCatalog ...
func (s *Section) GetCatalog(context.Context) (base.Cataloger, error) {
	return s.catalog, nil
}

func (s *Section) String() string {
	return fmt.Sprintf("&biquge.Section[bookName: %s, name: %s, url: %s, index: %d]", s.bookName, s.name, s.url, s.index)
}
