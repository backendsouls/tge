-- +goose Up
-- Quests and their ordered objectives (unique by order within a quest).
CREATE TABLE IF NOT EXISTS quests (
	name        TEXT PRIMARY KEY,
	description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS quest_objectives (
	quest       TEXT    NOT NULL,
	ord         INTEGER NOT NULL,
	description TEXT    NOT NULL,
	PRIMARY KEY (quest, ord)
);

-- +goose Down
DROP TABLE IF EXISTS quest_objectives;
DROP TABLE IF EXISTS quests;
