package base

import "context"

// Sectioner ...
type Sectioner interface {
	GetPre() Sectioner
	GetNext() Sectioner
	GetCatalog() Cataloger
	GetIndex() int
	GetName() string

	GetText(context.Context) (string, error)
	GetBody(context.Context) (string, error)

	Load(context.Context) error
	SetPre(context.Context, Sectioner)
	SetNext(context.Context, Sectioner)
	SetCatalog(ctx context.Context, catalog Cataloger)
}

type nilSection struct{}

func (s *nilSection) GetPre() Sectioner     { return NilSection }
func (s *nilSection) GetNext() Sectioner    { return NilSection }
func (s *nilSection) GetCatalog() Cataloger { return NilCatalog }
func (s *nilSection) GetIndex() int         { return -1 }
func (s *nilSection) GetName() string       { return "⭐️空章节⭐️" }

func (s *nilSection) GetText(context.Context) (string, error) { return "", nil }
func (s *nilSection) GetBody(context.Context) (string, error) { return "", nil }

func (s *nilSection) Load(context.Context) error                        { return nil }
func (s *nilSection) SetPre(context.Context, Sectioner)                 {}
func (s *nilSection) SetNext(context.Context, Sectioner)                {}
func (s *nilSection) SetCatalog(ctx context.Context, catalog Cataloger) {}

// NilSection 空章节
var NilSection Sectioner = &nilSection{}
