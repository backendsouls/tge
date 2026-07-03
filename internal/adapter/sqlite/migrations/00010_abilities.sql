-- +goose Up
-- Abilities: innate powers (name + description).
CREATE TABLE IF NOT EXISTS abilities (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS abilities;
