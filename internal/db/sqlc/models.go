package db

import (
	"database/sql"
)

type Token struct{
	ID				int64			`json:"id"`
	AccessToken		string			`json:"access_token"`
	CreatedAt  		sql.NullTime	`json:"created_at"`
}

