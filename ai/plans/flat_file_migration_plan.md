# Exhaustive Hexagonal TDD/DDD Migration Plan

This document is the absolute blueprint for migrating the `tge-go` codebase to a Flat File (JSON) storage engine while implementing the Universal Power System (DAG) framework. 

Every change is strictly governed by **Hexagonal Architecture (Ports and Adapters)**, **Domain-Driven Design (DDD)**, and **Test-Driven Development (TDD)**. All implementations will adhere to **SOLID** principles.

---

## 1. Architectural Map (Hexagonal / Ports & Adapters)
- **Domain (`internal/core/domain`)**: The core rules of the power system. Completely isolated from I/O.
- **Ports (`internal/core/port`)**: Interfaces defining what the Application expects from the outside world.
- **Adapters (`internal/adapter/file`)**: Driven adapters that implement the Ports by writing JSON to disk.
- **Application Service (`internal/core/service`)**: Orchestrates Use Cases by calling the Domain and driving the Ports.

---

## 2. Exhaustive File & Test Plan

### Phase A: The Domain Layer (Core Business Rules)
*Principles: Single Responsibility (Domain logic only). Open/Closed (Extensible traits).*

1. **`internal/core/domain/progression/power_node.go`** (Replaces `power.go`)
   - *Logic*: Define `PowerNode` struct (Parents, Siblings, Exclusives, Tags, StatVector, MaterialReq).
   - **TDD `power_node_test.go`**:
     - *Test 1*: Assert that adding a mutual exclusive tag prevents instantiation.
     - *Test 2*: Assert that `StatVector` correctly initializes default keys.

2. **`internal/core/domain/progression/mechanic_state.go`** (New)
   - *Logic*: Define `MechanicState` (BasePower, Tier, Awakenings, Alignments, EnergyPools, SpellSlots).
   - **TDD `mechanic_state_test.go`**:
     - *Test 1*: Assert `EnergyPools` correctly caps at a defined maximum.
     - *Test 2*: Assert `Alignment` cannot exceed domain boundaries (-100 to 100).

3. **`internal/core/domain/progression/power_system.go`**
   - *Logic*: Refactor to hold `Nodes map[string]*PowerNode` (O(1) graph traversal).
   - **TDD `power_system_test.go`**:
     - *Test 1*: Graph Validation—Assert that adding a cyclic dependency (A -> B -> A) returns a domain error.
     - *Test 2*: Assert `Names()` correctly traverses the DAG topologically.

4. **`internal/core/domain/character/character.go`**
   - *Logic*: Embed `MechanicState` and a record of Unlocked Nodes. Add `CalculateTotalPower()` method.
   - **TDD `character_test.go`**:
     - *Test 1*: Assert `CalculateTotalPower()` correctly sums `MechanicState.BasePower` and all active node vectors.

---

### Phase B: The Port Layer (Dependency Inversion)
*Principles: Interface Segregation (Small, focused interfaces). Dependency Inversion (Service depends on Interfaces).*

1. **`internal/core/port/power_system.go`**
   - *Refactor*: Remove SQLite-specific granularity.
   - *Interface*: `Save(ctx context.Context, ps progression.PowerSystem) error`, `FindByName(...)`, `List(...)`.
2. **`internal/core/port/character.go`**
   - *Refactor*: Remove `SaveCultivation`, `SaveSuperPower`, `UpdatePowerValue`. Because Flat Files serialize the aggregate root (the entire character), we only need `Save(ctx, char)`.
   - *Interface*: `Save(ctx, char)`, `FindByName(...)`, `List(...)`.

*(No tests written here, as these are pure interfaces).*

---

### Phase C: The Adapter Layer (JSON Persistence)
*Principles: Liskov Substitution (Adapters fulfill the exact contract of the Port).*

1. **`internal/adapter/file/power_system_repository.go`**
   - *Logic*: Implements `port.PowerSystemRepository`. Uses `os.WriteFile` and `encoding/json`.
   - **TDD `power_system_repository_test.go`**:
     - *Setup*: `t.TempDir()` to ensure complete isolation.
     - *Test 1*: Assert `Save()` marshals a complex DAG into `<name>.json` without losing relational data.
     - *Test 2*: Assert `FindByName()` reads the JSON and reconstructs the Go pointers correctly.
     - *Test 3*: Assert `FindByName()` returns `port.ErrPowerSystemNotFound` if the file does not exist.
     - *Test 4*: Assert reading a corrupted file returns a wrapped format error.

2. **`internal/adapter/file/character_repository.go`**
   - *Logic*: Implements `port.CharacterRepository`.
   - **TDD `character_repository_test.go`**:
     - *Setup*: `t.TempDir()`.
     - *Test 1*: Assert `Save()` correctly serializes the deeply nested `MechanicState` and `Inventory`.
     - *Test 2*: Assert `List()` reads the entire `data/characters/` directory, ignores non-json files, and returns a slice of `Character` structs.

---

### Phase D: The Application Service (Use Cases)
*Principles: Orchestration over computation. Delegates logic to Domain and persistence to Adapters.*

1. **`internal/core/service/character_service.go`**
   - *Logic*: Replace `Cultivate()` with generalized `TrainNode()`. Ensure it calls `CalculateTotalPower()` before calling `chars.Save()`.
   - **TDD `character_service_test.go`**:
     - *Setup*: Use mock implementations (or the file adapter pointed to a TempDir).
     - *Test 1*: Assert `TrainNode()` fails if the Character's inventory does not meet the Domain's `MaterialReq`.
     - *Test 2*: Assert `TrainNode()` applies post-combat Hooks (e.g., triggering a Zenkai boost logic if simulated HP was low).
     - *Test 3*: Assert Service correctly rolls back or handles errors if the driven Port returns an I/O error.

---

## 3. Execution Checklist
By enforcing TDD, we will execute the plan in this strict order:
1. Write the tests in `internal/core/domain/*_test.go`. Watch them fail.
2. Write the domain logic to make tests pass.
3. Update `internal/core/port/` interfaces.
4. Write the tests in `internal/adapter/file/*_test.go`. Watch them fail.
5. Write the file I/O logic to make tests pass.
6. Write the tests in `internal/core/service/*_test.go`. Watch them fail.
7. Refactor the service logic to make tests pass.
8. Wire the new adapter in `cmd/main.go` and delete `internal/adapter/sqlite`.
