package cultivation

// CultivationPath is a cultivation aspect a character develops. Paths are modeled as
// entities rather than enum values so each can accrue its own attributes and
// items over time without forcing changes on the cultivations that reference
// them (Open/Closed Principle).
type CultivationPath interface {
	// Name is the human-readable identity of the path.
	Name() string
}

// BodyCultivation is the physical cultivation path: tempering flesh, bones and blood.
// Future items (body-tempering techniques, resources, ...) will hang off it.
type BodyCultivation struct{}

// Name implements CultivationPath.
func (BodyCultivation) Name() string { return "Body" }

// SpiritCultivation is the energy cultivation path: refining qi and spiritual power.
type SpiritCultivation struct{}

// Name implements CultivationPath.
func (SpiritCultivation) Name() string { return "Spirit" }

// SoulCultivation is the consciousness cultivation path: strengthening the soul and will.
type SoulCultivation struct{}

// Name implements CultivationPath.
func (SoulCultivation) Name() string { return "Soul" }
