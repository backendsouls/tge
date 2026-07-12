-- +goose Up
ALTER TABLE power_systems ADD COLUMN kind TEXT NOT NULL DEFAULT 'Cultivation';

CREATE TABLE character_superpowers (
	character           TEXT    NOT NULL,
	system              TEXT    NOT NULL,
	path                TEXT    NOT NULL,
	tier                INTEGER NOT NULL,
	PRIMARY KEY (character, system, path),
	FOREIGN KEY (character) REFERENCES characters(name) ON DELETE CASCADE,
	FOREIGN KEY (system) REFERENCES power_systems(name) ON DELETE CASCADE,
	FOREIGN KEY (system, path) REFERENCES powers(system, name) ON DELETE CASCADE
);

-- +goose Down
-- +goose NO TRANSACTION
DROP TABLE character_superpowers;

-- Drop the kind column via table recreation for broader SQLite compatibility
-- Ensure foreign keys are disabled during the swap
PRAGMA foreign_keys = OFF;
BEGIN;
CREATE TABLE power_systems_dg_tmp (name TEXT PRIMARY KEY);
INSERT INTO power_systems_dg_tmp(name) SELECT name FROM power_systems;
DROP TABLE power_systems;
ALTER TABLE power_systems_dg_tmp RENAME TO power_systems;
COMMIT;
PRAGMA foreign_keys = ON;
