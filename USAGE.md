# TGE (The Grand Element) CLI Usage Guide

The `tge` command line interface is the primary tool for interacting with The Grand Element engine, allowing you to manage the engine's cosmos, power systems, characters, RPG mechanics, and more.

## General Usage

```bash
tge <command> [subcommand] [arguments...]
```

To see a list of available commands or get help for a specific command, you can use the `help` or `--help` flag:
```bash
tge help
tge <command> --help
```

## Available Commands

Here is an overview of the primary commands available:

### Core Progression & Systems
- **`realm`**: Manage cultivation realms, linear scaling models, and stage levels.
  - Subcommands: `add-level`, `create`, `list`, `show`.
- **`powersystem`**: Manage power systems modeled as trees of powers.
- **`species`**: Manage character species (with base stats, power values, and lifespan).

### Cosmology & Time
- **`omniverse`**, **`multiverse`**, **`universe`**, **`reality`**, **`realm`**: Manage the hierarchy of reality.
- **`timeline`**: Manage timelines and events belonging to any level of the cosmology.

### RPG Mechanics
- **`character`**: Manage characters, their attributes, inventory, and progression states.
- **`class`** / **`profession`**: Define classes and professions for characters.
- **`ability`** / **`skill`**: Configure abilities and skills that characters can learn.
- **`item`** / **`equipment`**: Create and configure items and equipment.
- **`effect`**: Define status effects.
- **`quest`**: Outline narrative or mechanical quests for characters.
- **`recipe`**: Setup crafting or alchemy recipes.

### Utility
- **`novel`**: Manage the story-level novel data.
- **`status`**: Output the full current status and progression of a character.

## Examples

### Character Status
View a character's current state, stats, cultivation progress, and unlocked systems:
```bash
tge status [character_id]
```

### Cultivation Realm Creation
Create a new cultivation realm and define its progression scales:
```bash
tge realm create --name "Foundation Establishment" --power-mult 10 --power-add 500 --lifespan-mult 3 --lifespan-add 500 --max-levels 9
```

### Adding Realm Levels
Add a stage to an existing realm:
```bash
tge realm add-level --realm "Foundation Establishment" --number 1 --name "First Level" --breakthrough 600
```

*Note: Breakthrough points define how many accumulated points are needed to break through to the next level.*
