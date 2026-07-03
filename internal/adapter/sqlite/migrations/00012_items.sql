-- +goose Up
-- Items: things that can be held in an inventory (name + description).
CREATE TABLE IF NOT EXISTS items (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS items;
