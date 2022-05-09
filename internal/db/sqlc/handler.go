package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TokenDB interface {
	Querier
	GetAccessTokens(ctx context.Context) ([]Token, error)
	GetLastAccessToken(ctx context.Context) (Token, error)
	GetAccessTokenById(ctx context.Context, id int) (Token, error)
	SetAccessToken(ctx context.Context, arg CreateTokenParams) (Token, error)
}

type SQLTokenDb struct {
	db *sql.DB
	*Queries
}

func NewTokenDb(db *sql.DB) TokenDB {
	return &SQLTokenDb{
		db:      db,
		Queries: New(db),
	}
}

func (tokenSql *SQLTokenDb) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := tokenSql.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (tokenSql *SQLTokenDb) SetAccessToken(ctx context.Context, arg CreateTokenParams) (Token, error) {
	var result Token
	err := tokenSql.execTx(ctx, func(q *Queries) error {
		var err error
		result, err = q.CreateToken(ctx, arg)

		if err != nil {
			return err
		}
		return err
	})

	return result, err
}

func (tokenSql *SQLTokenDb) GetAccessTokens(ctx context.Context) ([]Token, error) {
	return tokenSql.Queries.ListTokens(ctx)
}

func (tokenSql *SQLTokenDb) GetAccessTokenById(ctx context.Context, id int) (Token, error) {
	return tokenSql.Queries.GetTokenById(ctx, id)
}

func (tokenSql *SQLTokenDb) GetLastAccessToken(ctx context.Context) (Token, error) {
	return tokenSql.Queries.GetLastToken(ctx)
}
