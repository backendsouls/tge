# IDLE Progression Subsystem Plan

To solve the problem of resource regeneration (e.g., regaining Spell Slots after a rest, or restoring Magoi/Ki) and continuous progression without requiring constant user input, we will implement an **IDLE Subsystem**.

This subsystem will allow characters to passively gain points, power, levels, and items based on elapsed time or specific offline triggers.

## 1. Subsystem Architecture
The IDLE subsystem will track "offline progress" or background ticks.

### New Data Structure: `IdleState`
Attach an `IdleState` block to the `Character` document (in MongoDB):
```json
{
  "idle_state": {
    "last_sync_timestamp": 1718224500,
    "active_activity": "meditating_in_spirit_cave",
    "rates": {
      "stamina_regen_per_hour": 10,
      "mana_regen_per_hour": 50,
      "cultivation_points_per_hour": 100
    }
  }
}
```

## 2. Trigger Mechanics (Just-In-Time Calculation)
We do not need to run a background CRON job that constantly updates the database every second (which is highly unperformant). Instead, we use **Just-In-Time (JIT) Calculation**.

Whenever a character is loaded from the database (via `CharacterService.Character(ctx, name)`), the service intercepts the payload and calculates the IDLE gains:
1. `DeltaTime = CurrentTimestamp - IdleState.LastSyncTimestamp`
2. Apply Regeneration: `CurrentMana += (DeltaTime in hours) * ManaRegenRate`
3. Apply Progression: `CultivationPoints += (DeltaTime in hours) * CultivationRate`
4. Update `LastSyncTimestamp` to `CurrentTimestamp` and save back to the DB before returning the character to the user.

## 3. IDLE Activities (Assignments)
Characters can be assigned to different background activities that alter their `Rates`:
- **"Resting"**: Drastically increases `stamina_regen` and `mana_regen`. Automatically replenishes D&D "Spell Slots" if `DeltaTime > 8 hours` (Long Rest).
- **"Secluded Cultivation"**: Zero stamina regen, but high `cultivation_points_per_hour` to automatically advance through Xianxia realms.
- **"Gathering/Scavenging"**: Adds a % chance per hour to spawn specific items (e.g., herbs, beast cores) directly into the character's `Inventory`.

## 4. Implementation Steps
1. **Domain Update**: Add `IdleState` struct to `internal/core/domain/character/character.go`.
2. **Repository Update**: Ensure MongoDB persistence handles the new `idle_state` JSON block.
3. **Service Logic**: 
   - Create `internal/core/service/idle_service.go`.
   - Implement `CalculateOfflineGains(char *Character)` to perform the math.
   - Inject this calculation into the `FindByName` retrieval flow so it triggers invisibly to the user.
4. **CLI Commands**: Add commands to assign activities, e.g., `tge character idle <name> --activity="rest"`.
