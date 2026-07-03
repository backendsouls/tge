-- +goose Up
-- Professions: a character's vocation (name + description).
CREATE TABLE IF NOT EXISTS professions (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS professions;
