package progression

// SystemKind identifies the family a power system belongs to. Different kinds
// progress by different rules — a Cultivation system advances through realms and
// levels, a Magic system will advance by its own mechanics — so a character's
// per-node progress (PowerState) is kind-specific while the surrounding
// PowerSystem stays general.
type SystemKind string

const (
	// Cultivation is the realm/level progression kind.
	Cultivation SystemKind = "Cultivation"
	// Magic is a placeholder for a future, differently-progressing kind.
	Magic SystemKind = "Magic"
)

// Valid reports whether k is a known system kind.
func (k SystemKind) Valid() bool {
	switch k {
	case Cultivation, Magic:
		return true
	}
	return false
}

// PowerState is a character's progress at a single power node. Each SystemKind
// supplies its own implementation (CultivationState today, a Magic state later),
// which keeps a character's attained power (Character.Power) general over the
// kind of system rather than tied to cultivation.
type PowerState interface {
	// Kind reports which system kind this state belongs to.
	Kind() SystemKind
}
