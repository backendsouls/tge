-- +goose Up
-- Initial schema. CREATE ... IF NOT EXISTS keeps this safe to apply to a database
-- that predates migrations (created by the old bootstrap): the tables are no-ops
-- and goose simply records version 1 as applied.

CREATE TABLE IF NOT EXISTS realms (
	name                TEXT    PRIMARY KEY,
	power_multiplier    REAL    NOT NULL,
	power_adder         REAL    NOT NULL,
	lifespan_multiplier REAL    NOT NULL,
	lifespan_adder      REAL    NOT NULL,
	bottleneck_points   INTEGER NOT NULL,
	max_levels          INTEGER NOT NULL DEFAULT 0,
	main_max_levels     INTEGER NOT NULL DEFAULT 0
);

-- Ordered levels (sub-stages) within a realm, unique by number per realm.
CREATE TABLE IF NOT EXISTS realm_levels (
	realm               TEXT    NOT NULL,
	number              INTEGER NOT NULL,
	name                TEXT    NOT NULL,
	breakthrough_points INTEGER NOT NULL DEFAULT 0,
	bottleneck_points   INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (realm, number)
);

CREATE TABLE IF NOT EXISTS power_systems (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS realities (
	name TEXT PRIMARY KEY
);

-- Each location owns exactly one timeline, keyed by kind + key.
CREATE TABLE IF NOT EXISTS timelines (
	owner_kind TEXT NOT NULL,
	owner_key  TEXT NOT NULL,
	name       TEXT NOT NULL,
	PRIMARY KEY (owner_kind, owner_key)
);

-- Ordered events on a timeline, unique by order within the owning location.
CREATE TABLE IF NOT EXISTS timeline_events (
	owner_kind  TEXT    NOT NULL,
	owner_key   TEXT    NOT NULL,
	ord         INTEGER NOT NULL,
	description TEXT    NOT NULL,
	PRIMARY KEY (owner_kind, owner_key, ord)
);

-- --- RPG entities (keyed by name) ---
CREATE TABLE IF NOT EXISTS abilities (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS skills (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS effects (
	name        TEXT PRIMARY KEY,
	kind        TEXT NOT NULL,
	description TEXT NOT NULL
);

-- Equipment carries a slot and a stat bonus block.
CREATE TABLE IF NOT EXISTS equipment (
	name  TEXT PRIMARY KEY,
	slot  TEXT    NOT NULL,
	str   INTEGER NOT NULL,
	agi   INTEGER NOT NULL,
	intel INTEGER NOT NULL,
	vit   INTEGER NOT NULL,
	dex   INTEGER NOT NULL,
	wis   INTEGER NOT NULL,
	cha   INTEGER NOT NULL,
	luk   INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS professions (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS classes (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS quests (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

-- Ordered objectives within a quest.
CREATE TABLE IF NOT EXISTS quest_objectives (
	quest       TEXT    NOT NULL,
	ord         INTEGER NOT NULL,
	description TEXT    NOT NULL,
	PRIMARY KEY (quest, ord)
);

CREATE TABLE IF NOT EXISTS recipes (
	name   TEXT PRIMARY KEY,
	output TEXT NOT NULL
);

-- Input ingredients of a recipe, unique by item within the recipe.
CREATE TABLE IF NOT EXISTS recipe_inputs (
	recipe   TEXT    NOT NULL,
	item     TEXT    NOT NULL,
	quantity INTEGER NOT NULL,
	PRIMARY KEY (recipe, item)
);

-- A character's carried items (inventory), unique by item per character.
CREATE TABLE IF NOT EXISTS character_items (
	character TEXT    NOT NULL,
	item      TEXT    NOT NULL,
	quantity  INTEGER NOT NULL,
	PRIMARY KEY (character, item)
);

CREATE TABLE IF NOT EXISTS reality_omniverses (
	reality   TEXT NOT NULL,
	omniverse TEXT NOT NULL UNIQUE,
	PRIMARY KEY (reality, omniverse)
);

CREATE TABLE IF NOT EXISTS omniverses (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS omniverse_multiverses (
	omniverse  TEXT NOT NULL,
	multiverse TEXT NOT NULL UNIQUE,
	PRIMARY KEY (omniverse, multiverse)
);

CREATE TABLE IF NOT EXISTS multiverses (
	name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS multiverse_universes (
	multiverse TEXT NOT NULL,
	universe   TEXT NOT NULL UNIQUE,
	PRIMARY KEY (multiverse, universe)
);

CREATE TABLE IF NOT EXISTS universes (
	name TEXT PRIMARY KEY
);

-- A power system belongs to at most one universe (system is unique).
CREATE TABLE IF NOT EXISTS universe_systems (
	universe TEXT NOT NULL,
	system   TEXT NOT NULL UNIQUE,
	PRIMARY KEY (universe, system)
);

-- In-universe realms (locations), unique within a universe.
CREATE TABLE IF NOT EXISTS universe_realms (
	universe TEXT NOT NULL,
	name     TEXT NOT NULL,
	PRIMARY KEY (universe, name)
);

-- Powers form a tree per system (adjacency list). parent is NULL for roots.
CREATE TABLE IF NOT EXISTS powers (
	system TEXT NOT NULL,
	name   TEXT NOT NULL,
	parent TEXT,
	PRIMARY KEY (system, name)
);

CREATE TABLE IF NOT EXISTS species (
	name           TEXT PRIMARY KEY,
	power          REAL NOT NULL,
	lifespan       INTEGER NOT NULL,
	default_gender TEXT NOT NULL DEFAULT ''
);

-- Characters are keyed by name. They are created mortal; power is stored as a
-- string for now (computed numerically later).
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

-- A character's stat block, one row per character.
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

-- A character's cultivation state: one row per (system, path) node it has
-- progressed, anchored to a realm + level. The persisted form of Character.Power.
CREATE TABLE IF NOT EXISTS character_cultivations (
	character    TEXT    NOT NULL,
	system       TEXT    NOT NULL,
	path         TEXT    NOT NULL,
	realm        TEXT    NOT NULL,
	level_number INTEGER NOT NULL,
	level_name   TEXT    NOT NULL,
	points       INTEGER NOT NULL DEFAULT 0,
	progress     REAL    NOT NULL DEFAULT 0,
	PRIMARY KEY (character, system, path)
);

-- Each novel is led by one main character (main_character is unique).
CREATE TABLE IF NOT EXISTS novels (
	title          TEXT PRIMARY KEY,
	main_character TEXT NOT NULL UNIQUE
);

-- Volumes belong to a novel, numbered uniquely within it.
CREATE TABLE IF NOT EXISTS volumes (
	novel  TEXT    NOT NULL,
	number INTEGER NOT NULL,
	title  TEXT    NOT NULL,
	PRIMARY KEY (novel, number)
);

-- Chapters belong to a volume, numbered uniquely within that volume.
CREATE TABLE IF NOT EXISTS chapters (
	novel   TEXT    NOT NULL,
	volume  INTEGER NOT NULL,
	number  INTEGER NOT NULL,
	title   TEXT    NOT NULL,
	PRIMARY KEY (novel, volume, number)
);

-- +goose Down
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS volumes;
DROP TABLE IF EXISTS novels;
DROP TABLE IF EXISTS character_cultivations;
DROP TABLE IF EXISTS character_systems;
DROP TABLE IF EXISTS character_stats;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS species;
DROP TABLE IF EXISTS powers;
DROP TABLE IF EXISTS universe_realms;
DROP TABLE IF EXISTS universe_systems;
DROP TABLE IF EXISTS universes;
DROP TABLE IF EXISTS multiverse_universes;
DROP TABLE IF EXISTS multiverses;
DROP TABLE IF EXISTS omniverse_multiverses;
DROP TABLE IF EXISTS omniverses;
DROP TABLE IF EXISTS reality_omniverses;
DROP TABLE IF EXISTS character_items;
DROP TABLE IF EXISTS recipe_inputs;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS quest_objectives;
DROP TABLE IF EXISTS quests;
DROP TABLE IF EXISTS classes;
DROP TABLE IF EXISTS professions;
DROP TABLE IF EXISTS equipment;
DROP TABLE IF EXISTS effects;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS abilities;
DROP TABLE IF EXISTS timeline_events;
DROP TABLE IF EXISTS timelines;
DROP TABLE IF EXISTS realities;
DROP TABLE IF EXISTS power_systems;
DROP TABLE IF EXISTS realm_levels;
DROP TABLE IF EXISTS realms;
