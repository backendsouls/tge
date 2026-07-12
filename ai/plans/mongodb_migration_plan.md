# MongoDB Migration & Implementation Plan

Based on the decision to resolve the JSON and DAG structural complexity by migrating from SQLite to NoSQL (MongoDB), this plan outlines the architecture for the new persistence layer.

MongoDB is uniquely suited for this project because it natively stores JSON-like documents (BSON). This means our complex, highly-nested `PowerSystem` DAGs and our multi-dimensional `MechanicState` can be stored exactly as they look in Go memory, without complex relational table joins.

## 1. Dependency Updates
- Add the official MongoDB Go Driver to the project:
  ```bash
  go get go.mongodb.org/mongo-driver/mongo
  ```

## 2. Collection Design (Document Schemas)
Instead of 10+ disjoint relational tables, we will condense the data into clear, self-contained collections.

### Collection: `power_systems`
Each document represents a complete power universe.
```json
{
  "_id": "Nasuverse",
  "type": "Magic",
  "nodes": {
    "node_magic_circuits": {
      "name": "Magic Circuits",
      "category": "Foundation",
      "tags": ["internal", "mana_generator"],
      "stat_vector": {"Quality": "B", "Quantity": "A"},
      "parents": [],
      "drawbacks": ["Extreme pain upon activation"]
    },
    "node_tracing": {
      "name": "Tracing",
      "parents": ["node_magic_circuits"],
      "material_req": {"mana": 50}
    }
  }
}
```
**Advantage**: We load the entire DAG into memory with a single `db.collection.findOne()` call. No complex `JOIN` or recursive CTE logic is required.

### Collection: `characters`
Each document holds the complete state of a character. We completely eliminate the need for `character_items`, `character_stats`, `character_cultivations`, etc.
```json
{
  "_id": "Shirou",
  "type": "MainCharacter",
  "gender": "Male",
  "species": ["Human"],
  "stats": { "STR": 15, "AGI": 20 /* ... */ },
  "inventory": [
    {"item": "Health Potion", "quantity": 5}
  ],
  "mechanic_state": {
    "base_power": 100.0,
    "tier": 2,
    "is_awakened": false,
    "alignment": 10.5,
    "energy_pools": {"mana": 200, "stamina": 50},
    "spell_slots": {"1": 4, "2": 2}
  },
  "unlocked_nodes": [
    {
      "system": "Nasuverse",
      "node_id": "node_magic_circuits",
      "level": 5,
      "progress": 100.0
    }
  ]
}
```
**Advantage**: MongoDB supports atomic operations on single documents. When a character trains or unlocks a node, we can atomically push to `unlocked_nodes` and update `mechanic_state` simultaneously without needing complex multi-table SQL transactions.

## 3. Adapter Layer Implementation
- Create a new package: `internal/adapter/mongodb/`
- Implement the existing ports (`CharacterRepository`, `PowerSystemRepository`, etc.) using MongoDB.
  - `mongodb.NewCharacterRepository(client *mongo.Client, dbName string)`
  - Replace SQL `INSERT ON CONFLICT` with MongoDB `UpdateOne` using `$set` and `$upsert: true`.
  - Replace SQL `JOIN` queries with simple single-document fetches.

## 4. Configuration and Bootstrapping
- Update `internal/config/config.go` to accept `MongoURI` and `MongoDBName` instead of `SQLiteDSN`.
- Modify the `main.go` / `server` initialization to connect to the MongoDB cluster, run `client.Ping()`, and inject the MongoDB repository implementations into the `CharacterService`.

## 5. Solving Question 2 (State De-synchronization)
With MongoDB, we don't have to worry about data desync across multiple tables. Because `MechanicState` and `UnlockedNodes` exist inside the *exact same character document*, we can easily fetch the document, dynamically compute the `Tier` based on the `UnlockedNodes` array in Go memory, and save it back if caching is needed. 

## Next Steps
To begin this migration:
1. We will initialize the MongoDB driver in `go.mod`.
2. We will create the `internal/adapter/mongodb` directory and scaffold the repository implementations.
3. We will swap the DI bindings in the main server setup.
