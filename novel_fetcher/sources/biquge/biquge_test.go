package biquge

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestSearch(t *testing.T) {
	name := "大王饶命"
	s := NewSearcher(ctx)
	catalog, err := s.Search(ctx, name)
	assert.NoError(t, err)
	fmt.Println(err)
	if err == nil {
		fmt.Println(catalog.Count())
		section, err := catalog.Get(ctx, 1110)
		assert.NoError(t, err)
		fmt.Println(section.GetText(ctx))
		fmt.Println(section.GetNext())
		fmt.Println(section.GetPre())
	}
}

func TestSection(t *testing.T) {
	section := &Section{
		Name:     "测试章节",
		uri:      "/18_18727/8697302.html",
		BookName: "大王饶命",
		catalog:  nil,
		Index:    199,
	}
	b1 := time.Now()
	fmt.Println(section.GetText(ctx))
	b2 := time.Now()
	body, err := section.GetBody(ctx)
	assert.NoError(t, err)
	loadedSection, err := NewFromBody(ctx, body)
	actualText, err := loadedSection.GetText(ctx)
	assert.NoError(t, err)
	expectText, err := section.GetText(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectText, actualText)
	b3 := time.Now()
	fmt.Println(b2.Sub(b1), b3.Sub(b2))
}
