-- name: CreateVideo :one
INSERT INTO Videos (link, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetVideoByID :one
SELECT * FROM Videos
WHERE id = $1 LIMIT 1;

-- name: ListVideos :many
SELECT * FROM Videos
ORDER BY id
LIMIT $2
OFFSET $1;