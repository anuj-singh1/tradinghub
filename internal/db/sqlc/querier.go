package db

import "context"

type Querier interface {
	CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error)
	GetTokenById(ctx context.Context, id int) (Token, error)
	GetLastToken(ctx context.Context) (Token, error)
	ListTokens(ctx context.Context) ([]Token, error)
}

//var _ Querier = (*Queries)(nil)
