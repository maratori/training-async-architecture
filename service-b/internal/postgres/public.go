package postgres

import "context"

func (q *Queries) InsertX(ctx context.Context, x X) error {
	return q.insertX(ctx, insertXParams(x))
}
