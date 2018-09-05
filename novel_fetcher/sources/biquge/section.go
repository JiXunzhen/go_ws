package biquge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/utils"
	"github.com/PuerkitoBio/goquery"
)

// Section ...
type Section struct {
	Name     string
	BookName string
	Text     string
	Index    int

	catalog base.Cataloger
	pre     base.Sectioner
	next    base.Sectioner
	uri     string
	loaded  bool
}

// GetPre ...
func (s *Section) GetPre() base.Sectioner {
	return s.pre
}

// GetNext ...
func (s *Section) GetNext() base.Sectioner {
	return s.next
}

// GetCatalog ...
func (s *Section) GetCatalog() base.Cataloger {
	return s.catalog
}

// GetIndex ...
func (s *Section) GetIndex() int {
	return s.Index
}

// GetName ...
func (s *Section) GetName() string {
	return s.Name
}

// GetText ...
func (s *Section) GetText(ctx context.Context) (string, error) {
	if s.loaded {
		return s.Text, nil
	}
	err := s.Load(ctx)
	return s.Text, err
}

// GetBody ...
func (s *Section) GetBody(ctx context.Context) (string, error) {
	err := s.Load(ctx)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(s)
	return string(b), err
}

// Load ...
func (s *Section) Load(context.Context) error {
	if s.loaded {
		return nil
	}

	// get doc
	doc, err := goquery.NewDocument(utils.GenerateURL(Domain, s.uri))
	if err != nil {
		return err
	}

	// generate text
	//var textBuffer bytes.Buffer
	cont := doc.Find("#content")
	if cont == nil {
		return errors.New("section doc have no content")
	}
	text, err := cont.Html()
	if err != nil {
		return err
	}

	s.Text = s.handlerText(text)
	s.loaded = true
	return nil
}

func (s *Section) handlerText(rawText string) string {
	text := utils.Convert(rawText, utils.GbkToUtf8)
	// 笔趣阁拉下来的文章，每段开头都有四个聴 要干掉 顺便干掉<br/>
	re := regexp.MustCompile(`聽聽聽聽|<br/>`)
	text = re.ReplaceAllString(text, "")
	return text
}

// SetPre ...
func (s *Section) SetPre(ctx context.Context, pre base.Sectioner) {
	s.pre = pre
}

// SetNext ...
func (s *Section) SetNext(ctx context.Context, next base.Sectioner) {
	s.next = next
}

// SetCatalog ...
func (s *Section) SetCatalog(ctx context.Context, catalog base.Cataloger) {
	s.catalog = catalog
}

func (s *Section) String() string {
	return fmt.Sprintf("&biquge.Section[bookName: %s, name: %s, url: %s, index: %d]", s.BookName, s.Name, s.uri, s.Index)
}

// NewFromURL ...
func NewFromURL(_ context.Context, name, uri, bookName string, catalog base.Cataloger, index int) base.Sectioner {
	return &Section{
		Name:     name,
		BookName: bookName,
		Index:    index,
		uri:      uri,
		catalog:  catalog,
		loaded:   false,
	}
}

// NewFromBody ...
func NewFromBody(ctx context.Context, body string) (base.Sectioner, error) {
	var s Section
	s.loaded = true
	err := json.Unmarshal([]byte(body), &s)
	return &s, err
}
