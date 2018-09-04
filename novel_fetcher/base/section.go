package base

import "context"

// Sectioner ...
type Sectioner interface {
	GetText(context.Context) (string, error)
	Load(context.Context) error
	GetPre(context.Context) Sectioner
	GetNext(context.Context) Sectioner
	SetPre(context.Context, Sectioner)
	SetNext(context.Context, Sectioner)
	GetCatalog(context.Context) Cataloger
	GetIndex(context.Context) int
	GetBody(context.Context) (string, error)
	SetCatalog(ctx context.Context, catalog Cataloger)
}
