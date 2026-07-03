-- +goose Up
-- Realities (the outermost "Box") and their member omniverses. An omniverse
-- belongs to at most one reality.
CREATE TABLE IF NOT EXISTS realities (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS reality_omniverses (
	reality   TEXT NOT NULL,
	omniverse TEXT NOT NULL UNIQUE,
	PRIMARY KEY (reality, omniverse)
);

-- +goose Down
DROP TABLE IF EXISTS reality_omniverses;
DROP TABLE IF EXISTS realities;
