# Combat Engine Concepts (Deferred)

*Note: The implementation of a dedicated combat engine has been deferred. This document preserves the architectural ideas and design blueprints for future use.*

## The Problem with "God Functions"
If we implement a single `ExecuteAction(node, target, context)` function that handles every mechanic from every power system (e.g., Shinsu density, SCP reality bending, Logia intangibility), the function will quickly become an unmaintainable monolith of `if/else` checks.

## Proposed Solution: The Event / Hook Pattern
Instead of hardcoding tropes into the combat loop, the combat engine should act as an Event Bus. Mechanics (Traits, Vows, Environment) subscribe to these events and inject their own logic.

### 1. The Core Event Loop
When an action is taken, the engine fires sequential events:
1. `OnBeforeAction`: Traits can interrupt or prevent the action. (e.g., A "Silenced" status hook cancels verbal spellcasting).
2. `OnCalculateHit`: Modifies accuracy or enforces absolutes. (e.g., "Logia Intangibility" sets hit chance to 0% unless the attacker emits an "Armament Haki" tag).
3. `OnCalculateDamage`: Applies multipliers. (e.g., JJK "Binding Vows" multiply damage output by 1.5x based on a condition).
4. `OnAfterAction`: Triggers recoil or side effects. (e.g., MHA "Quirk Overuse" hook deals recoil damage to the caster).

### 2. Environmental Rules as Overrides
The arena/environment is just another entity that registers hooks.
- **Tower of God (Shinsu)**: Registers an `OnTurnStart` hook that deals massive damage or applies "Crushed" status to any character whose `Resistance` stat is lower than the environment's `AmbientDensity`.
- **The Matrix (Anomaly)**: A meta-trait that registers an `OnCalculateDamage` hook with supreme priority, directly editing the enemy's `Damage` struct to 0, completely bypassing standard math.

## Future Implementation Path
When combat is ready to be built:
1. Build an `EventBus` in `internal/core/domain/combat/bus.go`.
2. Update the `PowerNode` structure to hold an array of `Hooks` instead of hardcoded rules.
3. Write isolated hook functions (e.g., `LogiaHook`, `ZenkaiHook`) that are mapped to node tags.
