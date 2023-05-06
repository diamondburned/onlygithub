PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA strict = ON;

CREATE TABLE IF NOT EXISTS oauth_tokens (
	token TEXT PRIMARY KEY, -- our locally generated token
	provider TEXT NOT NULL,
	access_token TEXT NOT NULL,
	token_type TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	expires_in TIMESTAMP NOT NULL,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY, -- github user ID
	name TEXT NOT NULL, -- login
	email TEXT,
	nickname TEXT,
	avatar_url TEXT,
	joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tiers (
	id TEXT PRIMARY KEY, -- github tier ID
	name TEXT NOT NULL,
	price INTEGER NOT NULL, -- cents
	description TEXT
);

CREATE TABLE IF NOT EXISTS user_tiers (
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	tier_id TEXT REFERENCES tiers(id) ON DELETE SET NULL,
	price INTEGER NOT NULL, -- cents
	is_one_time BOOLEAN NOT NULL DEFAULT FALSE,
	is_custom_amount BOOLEAN NOT NULL DEFAULT FALSE,
	UNIQUE (user_id, tier_id)
);

CREATE TABLE IF NOT EXISTS assets (
	id TEXT PRIMARY KEY, -- xid
	type TEXT CHECK (type IN ('post', 'image')),
	data BLOB NOT NULL, -- JSON for post, blob data for image
	visibility TEXT CHECK (visibility IN ('', 'sponsor', 'public', 'private')) NOT NULL,
	minimum_cost INTEGER NOT NULL, -- cents
	last_updated TIMESTAMP
);
