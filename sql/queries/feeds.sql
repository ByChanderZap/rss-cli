-- name: CreateFeed :one
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

-- name: GetFeedsPopulated :many
SELECT f.id, f.created_at, f.updated_at, f.url, f.name as feedName, u.name as userName
FROM feeds f
INNER JOIN users u ON u.id = f.user_id;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;
