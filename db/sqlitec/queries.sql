-- name: SaveToken :exec
INSERT INTO oauth_tokens (token, provider, access_token, token_type, refresh_token, expires_in)
	VALUES (?, ?, ?, ?, ?, ?);

-- name: RestoreToken :one
SELECT * FROM oauth_tokens WHERE token = ? AND provider = ?;

-- name: DeleteToken :exec
DELETE FROM oauth_tokens WHERE token = ? AND provider = ?;

-- omitted: User :one

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
REPLACE INTO user_tiers (user_id, tier_id, price, is_one_time, is_custom_amount, started_at, renewed_at)
	VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: DeleteUserTier :exec
DELETE FROM user_tiers WHERE user_id = ?;

-- name: ImageAsset :one
SELECT filename, visibility, minimum_cost, last_updated FROM assets WHERE id = ?;

-- name: ImageData :one
SELECT data FROM assets WHERE id = ?;

-- name: DeleteImageAsset :exec
DELETE FROM assets WHERE id = ?;

-- name: CreateImageAsset :one
INSERT INTO assets (id, data, filename, preview_url, width, height, visibility, minimum_cost, last_updated)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime())
	RETURNING last_updated;

-- name: SetAssetVisibility :exec
UPDATE assets SET visibility = ? WHERE id = ?;

-- name: Post :one
SELECT * FROM posts WHERE id = ?;

-- name: Posts :many
SELECT * FROM posts ORDER BY id LIMIT 15;

-- name: PostsBefore :many
SELECT * FROM posts WHERE id < ? ORDER BY id LIMIT 15;

-- name: CreatePost :one
INSERT INTO posts (id, data, visibility, minimum_cost, last_updated)
	VALUES (?, ?, ?, ?, datetime())
	RETURNING last_updated;

-- name: SetPostVisibility :exec
UPDATE posts SET visibility = ? WHERE id = ?;

-- name: BindAssetToPost :exec
INSERT INTO asset_refs (post_id, asset_id) VALUES (?, ?);

-- name: PinAssetRef :exec
REPLACE INTO asset_refs (post_id, asset_id) VALUES (NULL, ?);

-- name: PostImages :many
SELECT assets.id, assets.preview_url, assets.width, assets.height
	FROM asset_refs
	INNER JOIN assets ON assets.id = asset_refs.asset_id
	WHERE asset_refs.post_id = ?;

-- -- name: UnusedImageAssetIDs :many
-- SELECT assets.id
-- FROM ASSETS
-- LEFT JOIN posts ON posts.image_asset_id = assets.id
-- WHERE posts.id IS NULL AND assets.type = 'image';

-- name: UserConfig :one
SELECT user_config FROM users WHERE id = ?;

-- name: SetUserConfig :exec
UPDATE users SET user_config = ? WHERE id = ?;

-- name: SiteConfig :one
SELECT site_config FROM users WHERE is_owner = TRUE;

-- name: SetSiteConfig :exec
UPDATE users SET site_config = ? WHERE is_owner = TRUE;

-- name: Tiers :many
SELECT * FROM tiers;

-- name: DeleteAllTiers :exec
DELETE FROM tiers;

-- name: CreateTier :exec
INSERT INTO tiers (id, name, price, description)
	VALUES (?, ?, ?, ?);
