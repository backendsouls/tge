-- +goose Up
-- Multiverses and their member universes. A universe belongs to at most one
-- multiverse.
CREATE TABLE IF NOT EXISTS multiverses (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS multiverse_universes (
	multiverse TEXT NOT NULL,
	universe   TEXT NOT NULL UNIQUE,
	PRIMARY KEY (multiverse, universe)
);

-- +goose Down
DROP TABLE IF EXISTS multiverse_universes;
DROP TABLE IF EXISTS multiverses;
