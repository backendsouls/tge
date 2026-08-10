-- +goose Up
-- Cultivation realms (ordered by tier) and their ordered sub-levels.
CREATE TABLE IF NOT EXISTS realms (
	name                TEXT    PRIMARY KEY,
	tier                INTEGER NOT NULL DEFAULT 0,
	power_multiplier    REAL    NOT NULL,
	power_adder         REAL    NOT NULL,
	lifespan_multiplier REAL    NOT NULL,
	lifespan_adder      REAL    NOT NULL,
	max_levels          INTEGER NOT NULL DEFAULT 0,
	main_max_levels     INTEGER NOT NULL DEFAULT 0
);

-- Ordered levels (sub-stages) within a realm, unique by number per realm.
CREATE TABLE IF NOT EXISTS realm_levels (
	realm               TEXT    NOT NULL,
	number              INTEGER NOT NULL,
	name                TEXT    NOT NULL,
	breakthrough_points INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (realm, number)
);

-- +goose Down
DROP TABLE IF EXISTS realm_levels;
DROP TABLE IF EXISTS realms;
