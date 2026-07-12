### Literature Review: Power Systems and Tech Tree Mechanics in Fandoms and Game Design

**Research Question**: How do renowned fictional power systems and gaming tech trees structure progression, balance power, and manage complexity, and how can these mechanics be adapted for programmatic implementation?
**Search Parameters**: Fandom databases (fandom.com) + "Lord of the Mysteries Pathways", "Hunter x Hunter Nen", "Path of Exile Passive Skill Tree". Scope limited to systemic mechanics, progression structures, and balancing rules.

**Thematic Synthesis**:

- **Theme 1: Rigid Sequential Progression with Behavioral Constraints**
  In *Lord of the Mysteries*, the "Pathway" system forces individuals (Beyonders) through a strict 10-tier sequence (Sequence 9 down to 0). While progression appears linearly deterministic, it is heavily gated by two mechanics: resource scarcity (potion formulas) and psychological stability (the "Acting Method"). Beyonders must act out the thematic role of their sequence to safely digest power. Failure to do so leads to "loss of control" (corruption and madness). 
  *Applicability*: Introduces the concept of "Alignment" or "Behavioral States" as a prerequisite for safely advancing tiers, penalizing players/characters who rush progression without fulfilling systemic behavioral requirements.

- **Theme 2: Spatial Progression and Keystone Drawbacks**
  *Path of Exile* (PoE) utilizes a massive, unified web of over 1,300 nodes. Progression is spatially limited; players can only unlock nodes adjacent to their currently allocated path. Furthermore, the system introduces "Keystones"—game-changing nodes that offer immense, build-defining buffs but strictly enforce severe mechanical drawbacks (e.g., massive damage but complete inability to evade). 
  *Applicability*: Progression mapped as a traversable graph (similar to a `SuperPower` tree) where major nodes carry forced negative modifiers (Drawbacks) to perfectly balance astronomical scaling.

- **Theme 3: Affinities, Limitations, and Self-Imposed Vows**
  The *Nen* system from *Hunter x Hunter* combines biological affinity (Enhancement, Transmutation, Conjuration, etc.) with intense personalization. A standout mechanic is "Limitations and Vows." Users can exponentially increase their aura output or specific abilities by setting severe self-imposed rules (e.g., "I will only use this technique against a specific group, otherwise I will die").
  *Applicability*: Allows for dynamic power multiplication based on strict condition checking (e.g., temporal, target-specific, or situational modifiers) rather than static numeric progression.

**Research Gaps for Implementation**:
1. **Dynamic Risk-Reward Enforcers**: Current implementations of progression (such as standard Cultivation interfaces) rely on linear numeric costs (e.g., breakthrough points). There is a significant gap in modeling *negative consequences* for failed or rushed progression. We lack a systemic `Affliction` or `Corruption` state that degrades a character's stats if they violate their power system's behavioral vows.
2. **Topological Web Pathing (Non-Hierarchical)**: Current standard tech trees strictly enforce a directed parent-child hierarchy. A gap exists in supporting undirected, web-based matrices (like Path of Exile or Final Fantasy X's Sphere Grid) where players can bridge across entirely different branches dynamically if nodes share geographical proximity.

**Annotated Bibliography**:
- *Lord of the Mysteries Wiki (Fandom)* - Contributes the framework for the "Acting Method" and potion-based sequence tiers. Excellent model for high-risk, high-reward linear progression scaling.
- *Hunterpedia: Nen (Fandom)* - Provides the foundation for "Vows and Limitations", demonstrating how to balance overpowered abilities via strict conditional programming.
- *Path of Exile Official / Community Guides* - Highlights the architectural design of a unified, attribute-aligned web matrix, demonstrating how spatial start locations and Keystone drawbacks naturally enforce diverse character builds.
