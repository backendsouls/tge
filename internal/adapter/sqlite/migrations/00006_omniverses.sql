-- +goose Up
-- Omniverses and their member multiverses. A multiverse belongs to at most one
-- omniverse.
CREATE TABLE IF NOT EXISTS omniverses (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS omniverse_multiverses (
	omniverse  TEXT NOT NULL,
	multiverse TEXT NOT NULL UNIQUE,
	PRIMARY KEY (omniverse, multiverse)
);

-- +goose Down
DROP TABLE IF EXISTS omniverse_multiverses;
DROP TABLE IF EXISTS omniverses;
