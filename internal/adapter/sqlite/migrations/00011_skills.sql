-- +goose Up
-- Skills: learned, trainable abilities (name + description).
CREATE TABLE IF NOT EXISTS skills (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS skills;
