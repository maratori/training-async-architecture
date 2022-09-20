-- name: SelectX :many
SELECT *
FROM x;

-- name: insertX :exec
INSERT INTO x (id, name)
VALUES ($1, $2);
