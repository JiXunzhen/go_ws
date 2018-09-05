package base

import "context"

// NovelSource ...
type NovelSource interface {
	Save(ctx context.Context, bookName string, sections map[int]Sectioner) (failIndexes []int)
	Load(ctx context.Context, bookName string) ([]string, error)
}
