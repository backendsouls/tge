-- +goose Up
-- Power systems and their tree of powers (adjacency list; parent NULL = root).
CREATE TABLE IF NOT EXISTS power_systems (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS powers (
	system TEXT NOT NULL,
	name   TEXT NOT NULL,
	parent TEXT,
	PRIMARY KEY (system, name)
);

-- +goose Down
DROP TABLE IF EXISTS powers;
DROP TABLE IF EXISTS power_systems;
