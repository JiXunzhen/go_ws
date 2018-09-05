package source

import (
	"context"

	"github.com/JiXunzhen/go_ws/novel_fetcher/base"
	"github.com/JiXunzhen/go_ws/novel_fetcher/sources/biquge"
)

// SearcherNames ...
var SearcherNames = []string{
	"笔趣阁",
}

// SearcherMap ...
var SearcherMap = map[string]base.Searcher{}

// Init ...
func Init(ctx context.Context) error {
	SearcherMap["笔趣阁"] = biquge.NewSearcher(ctx)
	return nil
}
