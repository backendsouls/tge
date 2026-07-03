-- +goose Up
-- Classes: a character's combat archetype (name + description).
CREATE TABLE IF NOT EXISTS classes (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS classes;
