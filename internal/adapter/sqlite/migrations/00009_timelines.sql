-- +goose Up
-- Each location owns one timeline (keyed by owner kind + key) with ordered events.
CREATE TABLE IF NOT EXISTS timelines (
	owner_kind TEXT NOT NULL,
	owner_key  TEXT NOT NULL,
	name       TEXT NOT NULL,
	PRIMARY KEY (owner_kind, owner_key)
);

CREATE TABLE IF NOT EXISTS timeline_events (
	owner_kind  TEXT    NOT NULL,
	owner_key   TEXT    NOT NULL,
	ord         INTEGER NOT NULL,
	description TEXT    NOT NULL,
	PRIMARY KEY (owner_kind, owner_key, ord)
);

-- +goose Down
DROP TABLE IF EXISTS timeline_events;
DROP TABLE IF EXISTS timelines;
