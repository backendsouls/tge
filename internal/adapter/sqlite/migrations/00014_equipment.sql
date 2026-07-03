-- +goose Up
-- Equipment: an equippable item with a slot and a stat bonus block.
CREATE TABLE IF NOT EXISTS equipment (
	name  TEXT PRIMARY KEY,
	slot  TEXT    NOT NULL,
	str   INTEGER NOT NULL,
	agi   INTEGER NOT NULL,
	intel INTEGER NOT NULL,
	vit   INTEGER NOT NULL,
	dex   INTEGER NOT NULL,
	wis   INTEGER NOT NULL,
	cha   INTEGER NOT NULL,
	luk   INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS equipment;
