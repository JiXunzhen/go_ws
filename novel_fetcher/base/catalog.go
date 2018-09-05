package base

import (
	"context"
)

// Cataloger ...
type Cataloger interface {
	Count() int
	GetBookName() string
	List() []Sectioner

	Get(context.Context, int) (Sectioner, error)
	Flush(context.Context, bool) error
	Save(ctx context.Context, start, end int, source NovelSource) error
	LoadFromSource(context.Context, NovelSource) error
	Load(ctx context.Context, start, end int) error
}

const (
	// SaveRetry 保存可重试次数
	SaveRetry = 3
	// LoadRetry 预加载可重试次数
	LoadRetry = 3
)

type nilCatalog struct{}

func (c *nilCatalog) Count() int                                                         { return 0 }
func (c *nilCatalog) Get(context.Context, int) (Sectioner, error)                        { return NilSection, nil }
func (c *nilCatalog) Flush(context.Context, bool) error                                  { return nil }
func (c *nilCatalog) Save(ctx context.Context, start, end int, source NovelSource) error { return nil }
func (c *nilCatalog) LoadFromSource(context.Context, NovelSource) error                  { return nil }
func (c *nilCatalog) Load(ctx context.Context, start, end int) error                     { return nil }
func (c *nilCatalog) GetBookName() string                                                { return "😛空目录😱" }
func (c *nilCatalog) List() []Sectioner                                                  { return []Sectioner{} }

// NilCatalog ...
var NilCatalog Cataloger = &nilCatalog{}
