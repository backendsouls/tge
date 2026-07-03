-- +goose Up
-- Universes, the power systems that belong to them (a system is in one universe),
-- and their in-universe realms (locations; unrelated to the cultivation realms).
CREATE TABLE IF NOT EXISTS universes (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS universe_systems (
	universe TEXT NOT NULL,
	system   TEXT NOT NULL UNIQUE,
	PRIMARY KEY (universe, system)
);

CREATE TABLE IF NOT EXISTS universe_realms (
	universe TEXT NOT NULL,
	name     TEXT NOT NULL,
	PRIMARY KEY (universe, name)
);

-- +goose Down
DROP TABLE IF EXISTS universe_realms;
DROP TABLE IF EXISTS universe_systems;
DROP TABLE IF EXISTS universes;
