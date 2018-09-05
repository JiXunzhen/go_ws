package biquge

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/utils"
)

const (
	// Domain ...
	Domain = "http://www.biquge.com.tw"

	// SearchURI ...
	SearchURI = "/modules/article/soshu.php"
)

// Entrance ...
type Entrance struct {
}

// Search ...
func (s *Entrance) Search(ctx context.Context, name string) (base.Cataloger, error) {
	formValues := url.Values{
		"searchkey": {utils.Convert(name, utils.Utf8ToGbk)},
	}
	request, err := http.NewRequest("POST", utils.GenerateURL(Domain, SearchURI), strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return NewCatalog(ctx, request, name)
}

// Name ...
func (s *Entrance) Name() string {
	return "笔趣阁"
}

// NewSearcher ...
func NewSearcher(_ context.Context) base.Searcher {
	return &Entrance{}
}
