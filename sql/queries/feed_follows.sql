-- name: CreateFeedFollow :one

INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetFeedFollowsByUserId :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: GetFeedFollowsByUserApiKey :many
SELECT * FROM feed_follows WHERE user_id = (SELECT id from users WHERE api_key = $1);

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_follows.id =$1 AND user_id =(SELECT users.id FROM users WHERE api_key =$2);
