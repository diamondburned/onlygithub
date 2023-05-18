PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA strict = ON;

-- MIGRATE --

CREATE TABLE IF NOT EXISTS oauth_tokens (
	token TEXT PRIMARY KEY, -- our locally generated token
	provider TEXT NOT NULL,
	access_token TEXT NOT NULL,
	token_type TEXT NOT NULL,
	refresh_token TEXT NOT NULL,
	expires_in TIMESTAMP NOT NULL,
	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY, -- github user ID
	username TEXT UNIQUE NOT NULL, -- login
	email TEXT NOT NULL,
	real_name TEXT NOT NULL,
	pronouns TEXT NOT NULL,
	avatar_url TEXT NOT NULL,
	joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	is_owner BOOLEAN NOT NULL DEFAULT FALSE,
	user_config BLOB NOT NULL DEFAULT 'null', -- JSON
	site_config BLOB
);

CREATE TABLE IF NOT EXISTS tiers (
	id TEXT PRIMARY KEY, -- github tier ID
	name TEXT NOT NULL,
	price INTEGER NOT NULL, -- cents
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tiers (
	user_id TEXT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	tier_id TEXT REFERENCES tiers(id) ON DELETE SET NULL,
	price INTEGER NOT NULL, -- cents
	is_one_time BOOLEAN NOT NULL DEFAULT FALSE,
	is_custom_amount BOOLEAN NOT NULL DEFAULT FALSE,
	started_at TIMESTAMP NOT NULL,
	renewed_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS assets (
	id TEXT PRIMARY KEY, -- xid
	type TEXT CHECK (type IN ('post', 'image')),
	data BLOB NOT NULL, -- JSON for post, blob data for image
	visibility TEXT NOT NULL CHECK (visibility IN ('', 'sponsor', 'public', 'private')) NOT NULL,
	minimum_cost INTEGER NOT NULL, -- cents
	last_updated TIMESTAMP
);

-- MIGRATE --

ALTER TABLE assets ADD COLUMN filename TEXT NOT NULL DEFAULT '';
