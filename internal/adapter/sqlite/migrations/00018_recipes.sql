-- +goose Up
-- Recipes and their input ingredients (item + quantity, unique per item).
CREATE TABLE IF NOT EXISTS recipes (
	name   TEXT PRIMARY KEY,
	output TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS recipe_inputs (
	recipe   TEXT    NOT NULL,
	item     TEXT    NOT NULL,
	quantity INTEGER NOT NULL,
	PRIMARY KEY (recipe, item)
);

-- +goose Down
DROP TABLE IF EXISTS recipe_inputs;
DROP TABLE IF EXISTS recipes;
