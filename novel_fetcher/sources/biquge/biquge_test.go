package biquge

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	name := "大王饶命"
	s := NewSearcher(nil)
	catalog, err := s.Search(nil, name)
	assert.NoError(t, err)
	fmt.Println(err)
	if err == nil {
		fmt.Println(catalog.Count(nil))
		section, err := catalog.Get(nil, 200)
		assert.NoError(t, err)
		fmt.Println(section)
		fmt.Println(section.GetNext(nil))
		fmt.Println(section.GetPre(nil))
	}
}
