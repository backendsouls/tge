-- +goose Up
-- Novels led by one main character, with ordered volumes and their chapters.
CREATE TABLE IF NOT EXISTS novels (
	title          TEXT PRIMARY KEY,
	main_character TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS volumes (
	novel  TEXT    NOT NULL,
	number INTEGER NOT NULL,
	title  TEXT    NOT NULL,
	PRIMARY KEY (novel, number)
);

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
