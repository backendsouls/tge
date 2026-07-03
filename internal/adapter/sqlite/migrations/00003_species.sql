-- +goose Up
-- Species with base status values and a default gender.
CREATE TABLE IF NOT EXISTS species (
	name           TEXT PRIMARY KEY,
	power          REAL NOT NULL,
	lifespan       INTEGER NOT NULL,
	default_gender TEXT NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE IF EXISTS species;
