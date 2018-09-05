package base

import "context"

// Searcher ...
type Searcher interface {
	Search(ctx context.Context, name string) (Cataloger, error)
	Name() string
}
