package db

import (
	"context"
	"database/sql"
)

type TokenDB interface {
	GetAccessToken(ctx context.Context, arg CreateJobParams) (Result, error)
	SetAccessToken(ctx context.Context) (Result, error)
}

type CreateJobParams struct {
	AccessToken		string			`json:"access_token"`
	CreatedAt  		sql.NullTime	`json:"created_at"`
}
type SQLTokenDb struct {
	db *sql.DB
	*Queries
}


type Result struct {
	Token
}