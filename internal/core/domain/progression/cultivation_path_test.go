package progression

import "testing"

func TestPathNames(t *testing.T) {
	cases := []struct {
		path CultivationPath
		want string
	}{
		{BodyCultivation{}, "Body"},
		{SpiritCultivation{}, "Spirit"},
		{SoulCultivation{}, "Soul"},
	}
	for _, c := range cases {
		if got := c.path.Name(); got != c.want {
			t.Errorf("%T.Name() = %q, want %q", c.path, got, c.want)
		}
	}
}
