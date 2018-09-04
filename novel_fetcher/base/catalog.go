package base

import "context"

// Cataloger ...
type Cataloger interface {
	Count(context.Context) int
	Get(context.Context, int) (Sectioner, error)
	Flush(context.Context, bool) error
	Save(ctx context.Context, start, end int, source NovelSource) error
	LoadFromSource(context.Context, NovelSource) error
}

const (
	// DefaultPreload 默认预加载
	DefaultPreload = 50

	// SaveRetry 保存可重试次数
	SaveRetry = 3
)
