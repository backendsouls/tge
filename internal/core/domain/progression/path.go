package progression

// Path is a cultivation aspect a character develops. Paths are modeled as
// entities rather than enum values so each can accrue its own attributes and
// items over time without forcing changes on the cultivations that reference
// them (Open/Closed Principle).
type Path interface {
	// Name is the human-readable identity of the path.
	Name() string
}

// Body is the physical cultivation path: tempering flesh, bones and blood.
// Future items (body-tempering techniques, resources, ...) will hang off it.
type Body struct{}

// Name implements Path.
func (Body) Name() string { return "Body" }

// Spirit is the energy cultivation path: refining qi and spiritual power.
type Spirit struct{}

// Name implements Path.
func (Spirit) Name() string { return "Spirit" }

// Soul is the consciousness cultivation path: strengthening the soul and will.
type Soul struct{}

// Name implements Path.
func (Soul) Name() string { return "Soul" }
