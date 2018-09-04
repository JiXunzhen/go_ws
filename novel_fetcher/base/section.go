package base

import "context"

// Sectioner ...
type Sectioner interface {
	Text(context.Context) (string, error)
	GetPre(context.Context) (Sectioner, error)
	GetNext(context.Context) (Sectioner, error)
	SetPre(context.Context, Sectioner)
	SetNext(context.Context, Sectioner)
	GetCatalog(context.Context) (Cataloger, error)
}
