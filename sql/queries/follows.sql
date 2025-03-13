-- name: CreateFeedFollow :one
WITH feed_follows_insertion AS (
    INSERT INTO feed_follows (created_at, updated_at, feed_id, user_id)
    VALUES (
        $1,
        $2,
        $3,
        $4
    )
    RETURNING *
)
SELECT 
    feed_follows_insertion.created_at, 
    feed_follows_insertion.updated_at, 
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows_insertion
INNER JOIN feeds
ON feed_follows_insertion.feed_id = feeds.id
INNER JOIN users
ON feed_follows_insertion.user_id = users.id
WHERE feed_follows_insertion.feed_id = feeds.id and feed_follows_insertion.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT users.name AS user_name, feeds.name as feed_name
FROM feed_follows
INNER JOIN users
ON feed_follows.user_id = users.id
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollowForUser :one
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2
RETURNING *;
