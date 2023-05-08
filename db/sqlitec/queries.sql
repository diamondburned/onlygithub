-- name: SaveToken :exec
INSERT INTO oauth_tokens (token, provider, access_token, token_type, refresh_token, expires_in)
	VALUES (?, ?, ?, ?, ?, ?);

-- name: RestoreToken :one
SELECT * FROM oauth_tokens WHERE token = ? AND provider = ?;

-- name: User :one
SELECT
	users.id, users.username, users.email, users.real_name, users.pronouns, users.avatar_url, users.joined_at,
	user_tiers.price AS tier_price,
	user_tiers.is_one_time AS tier_is_one_time,
	user_tiers.is_custom_amount AS tier_is_custom_amount,
	user_tiers.started_at AS tier_started_at,
	user_tiers.renewed_at AS tier_renewed_at,
	tiers.id AS tier_id,
	tiers.name AS tier_name,
	tiers.description AS tier_description
FROM users AS users -- https://github.com/kyleconroy/sqlc/issues/2271
LEFT JOIN user_tiers ON users.id = user_tiers.user_id
LEFT JOIN tiers ON user_tiers.tier_id = tiers.id
WHERE users.id = ?;

-- name: Owner :one
SELECT * FROM users WHERE is_owner = TRUE;

-- name: MakeOwner :exec
UPDATE users SET is_owner = TRUE WHERE username = ?;

-- name: UpdateUser :exec
INSERT INTO users (id, username, email, real_name, pronouns, avatar_url)
	VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (id) DO UPDATE SET
	username = excluded.username,
	email = excluded.email,
	real_name = excluded.real_name,
	pronouns = excluded.pronouns,
	avatar_url = excluded.avatar_url;

-- name: SetTier :exec
REPLACE INTO tiers (id, name, price, description)
	VALUES (?, ?, ?, ?);

-- name: SetUserTier :exec
REPLACE INTO user_tiers (user_id, tier_id, price, is_one_time, is_custom_amount)
	VALUES (?, ?, ?, ?, ?);

-- name: DeleteUserTier :exec
DELETE FROM user_tiers WHERE user_id = ?;

-- name: ImageAsset :one
SELECT * FROM assets WHERE id = ? AND type = 'image';

-- name: PostAsset :one
SELECT * FROM assets WHERE id = ? AND type = 'post';

-- name: UserConfig :one
SELECT user_config FROM users WHERE id = ?;

-- name: SetUserConfig :exec
UPDATE users SET user_config = ? WHERE id = ?;

-- name: SiteConfig :one
SELECT site_config FROM users WHERE is_owner = TRUE;

-- name: SetSiteConfig :exec
UPDATE users SET site_config = ? WHERE is_owner = TRUE;
