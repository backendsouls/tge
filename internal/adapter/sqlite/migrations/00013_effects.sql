-- +goose Up
-- Effects: modifiers/conditions with a kind (Buff, Debuff, Status).
CREATE TABLE IF NOT EXISTS effects (
	name        TEXT PRIMARY KEY,
	kind        TEXT NOT NULL,
	description TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS effects;
