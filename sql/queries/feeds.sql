-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ResetFeeds :exec
DELETE FROM feeds WHERE id IS NOT NULL;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedFromURL :one
SELECT * FROM feeds WHERE url=$1;
