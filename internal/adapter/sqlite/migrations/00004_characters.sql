-- +goose Up
-- Characters and their owned data: stats, power-system membership, inventory and
-- cultivation state.
CREATE TABLE IF NOT EXISTS characters (
	name       TEXT    PRIMARY KEY,
	type       TEXT    NOT NULL,
	gender     TEXT    NOT NULL,
	species    TEXT    NOT NULL,
	power      TEXT    NOT NULL,
	age        INTEGER NOT NULL,
	lifespan   INTEGER NOT NULL,
	class      TEXT    NOT NULL DEFAULT '',
	profession TEXT    NOT NULL DEFAULT ''
);

-- One stat block per character.
CREATE TABLE IF NOT EXISTS character_stats (
	character TEXT    PRIMARY KEY,
	str       INTEGER NOT NULL,
	agi       INTEGER NOT NULL,
	intel     INTEGER NOT NULL,
	vit       INTEGER NOT NULL,
	dex       INTEGER NOT NULL,
	wis       INTEGER NOT NULL,
	cha       INTEGER NOT NULL,
	luk       INTEGER NOT NULL
);

-- A character belongs to one or more power systems (many-to-many by name).
CREATE TABLE IF NOT EXISTS character_systems (
	character TEXT NOT NULL,
	system    TEXT NOT NULL,
	PRIMARY KEY (character, system)
);

-- A character's carried items (inventory), unique by item per character.
CREATE TABLE IF NOT EXISTS character_items (
	character TEXT    NOT NULL,
	item      TEXT    NOT NULL,
	quantity  INTEGER NOT NULL,
	PRIMARY KEY (character, item)
);

-- A character's cultivation state: one row per (system, path) node, anchored to a
-- realm + level, with the two-phase breakthrough/bottleneck progress and the
-- level's thresholds captured for a self-contained status view.
CREATE TABLE IF NOT EXISTS character_cultivations (
	character           TEXT    NOT NULL,
	system              TEXT    NOT NULL,
	path                TEXT    NOT NULL,
	realm               TEXT    NOT NULL,
	level_number        INTEGER NOT NULL,
	level_name          TEXT    NOT NULL,
	breakthrough_points INTEGER NOT NULL DEFAULT 0,
	bottleneck_points   INTEGER NOT NULL DEFAULT 0,
	points              INTEGER NOT NULL DEFAULT 0,
	bottleneck          INTEGER NOT NULL DEFAULT 0,
	progress            REAL    NOT NULL DEFAULT 0,
	PRIMARY KEY (character, system, path)
);

-- +goose Down
DROP TABLE IF EXISTS character_cultivations;
DROP TABLE IF EXISTS character_items;
DROP TABLE IF EXISTS character_systems;
DROP TABLE IF EXISTS character_stats;
DROP TABLE IF EXISTS characters;
