package base

import "context"

// Cataloger ...
type Cataloger interface {
	Count(context.Context) int
	Get(context.Context, int) (Sectioner, error)
	Flush(context.Context, bool) error
	Save(ctx context.Context, start, end int) error
	SetPreLoad(context.Context, int)
}

const (
	// DefaultPreload 默认预加载
	DefaultPreload = 50
)
