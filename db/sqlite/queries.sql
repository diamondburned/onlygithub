-- name: SaveToken :exec
INSERT INTO oauth_tokens (token, provider, access_token, token_type, refresh_token, expires_in)
	VALUES (?, ?, ?, ?, ?, ?);

-- name: RestoreToken :one
SELECT * FROM oauth_tokens WHERE token = ? AND provider = ?;

-- name: User :one
SELECT * FROM users WHERE id = ?;

-- name: SetUser :exec
REPLACE INTO users (id, name, email, nickname)
	VALUES (?, ?, ?, ?);

-- name: SetTier :exec
REPLACE INTO tiers (id, name, price, description)
	VALUES (?, ?, ?, ?);

-- name: SetUserTier :exec
REPLACE INTO user_tiers (user_id, tier_id, price, is_one_time, is_custom_amount)
	VALUES (?, ?, ?, ?, ?);

-- name: ImageAsset :one
SELECT * FROM assets WHERE id = ? AND type = 'image';

-- name: PostAsset :one
SELECT * FROM assets WHERE id = ? AND type = 'post';
