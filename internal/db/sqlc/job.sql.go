package db

import (
	"context"
)

const createToken = `-- name: CreateJob :one
INSERT INTO token (
  access_token
) VALUES (
  $1
) RETURNING id, access_token, created_at
`

type CreateTokenParams struct {
	AccessToken string `json:"access_token"`
}

func (q *Queries) CreateToken(ctx context.Context, arg CreateTokenParams) (Token, error) {
	row := q.db.QueryRowContext(ctx, createToken,
		arg.AccessToken,
	)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.AccessToken,
		&i.CreatedAt,
	)
	return i, err
}

const getTokenById = `-- name: GetJobById :one
SELECT id, access_token, created_at FROM token
WHERE id = $1
`

func (q *Queries) GetTokenById(ctx context.Context, id int) (Token, error) {
	row := q.db.QueryRowContext(ctx, getTokenById, id)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.AccessToken,
		&i.CreatedAt,
	)
	return i, err
}

const getLastToken = `-- name: GetJobById :one
SELECT id, access_token, created_at FROM token
ORDER BY created_at desc LIMIT 1 
`

func (q *Queries) GetLastToken(ctx context.Context) (Token, error) {
	row := q.db.QueryRowContext(ctx, getLastToken)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.AccessToken,
		&i.CreatedAt,
	)
	return i, err
}

const listTokens = `-- name: GetJobById :one
SELECT id, access_token, created_at FROM token
WHERE id = $1
`

func (q *Queries) ListTokens(ctx context.Context) ([]Token, error) {
	rows, err := q.db.QueryContext(ctx, listTokens)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Token
	for rows.Next() {
		var i Token
		if err = rows.Scan(&i.ID, &i.AccessToken, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}
