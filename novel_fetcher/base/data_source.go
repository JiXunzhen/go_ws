package base

import "context"

// NovelSource ...
type NovelSource interface {
	Save(ctx context.Context, bookName string, sections []Sectioner) (failIndexes []int)
	Load(ctx context.Context, bookName string) ([]string, error)
}
