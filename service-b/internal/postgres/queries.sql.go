// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: queries.sql

package postgres

import (
	"context"

	"github.com/google/uuid"
)

const selectX = `-- name: SelectX :many
SELECT id, name
FROM x
`

func (q *Queries) SelectX(ctx context.Context) ([]X, error) {
	rows, err := q.db.QueryContext(ctx, selectX)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []X
	for rows.Next() {
		var i X
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertX = `-- name: insertX :exec
INSERT INTO x (id, name)
VALUES ($1, $2)
`

type insertXParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) insertX(ctx context.Context, arg insertXParams) error {
	_, err := q.db.ExecContext(ctx, insertX, arg.ID, arg.Name)
	return err
}
