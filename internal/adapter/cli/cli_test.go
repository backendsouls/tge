package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"tge/internal/adapter/cli"
	"tge/internal/core/domain/character"
	"tge/internal/core/domain/cosmology"
	"tge/internal/core/domain/novel"
	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
)

// fakeService is an in-memory port.RealmService double for driving the CLI
// without any storage or process.
type fakeService struct {
	realms []progression.Realm
}

func (f *fakeService) AddRealm(_ context.Context, cfg progression.RealmConfig) (progression.Realm, error) {
	r, err := progression.NewRealm(cfg)
	if err != nil {
		return progression.Realm{}, err
	}
	f.realms = append(f.realms, r)
	return r, nil
}

func (f *fakeService) ListRealms(context.Context) ([]progression.Realm, error) {
	return f.realms, nil
}

func (f *fakeService) GetRealm(_ context.Context, name string) (progression.Realm, error) {
	for _, r := range f.realms {
		if r.Name == name {
			return r, nil
		}
	}
	return progression.Realm{}, port.ErrRealmNotFound
}

func (f *fakeService) AddLevel(_ context.Context, in port.AddLevelInput) (progression.Realm, error) {
	for i := range f.realms {
		if f.realms[i].Name == in.Realm {
			if err := f.realms[i].AddLevel(in.Number, in.Name, in.BreakthroughPoints, in.BottleneckPoints); err != nil {
				return progression.Realm{}, err
			}
			return f.realms[i], nil
		}
	}
	return progression.Realm{}, port.ErrRealmNotFound
}

// fakeCharacterService is a port.CharacterService double for driving the CLI.
type fakeCharacterService struct {
	created    *character.Character
	createErr  error
	main       character.Character
	mainErr    error
	byName     map[string]character.Character
	list       []character.Character
	listErr    error
	cultivated *port.CultivateInput
	cultErr    error
	trained    *port.TrainInput
	trainErr   error
}

func (f *fakeCharacterService) CreateCharacter(_ context.Context, in port.CreateCharacterInput) (character.Character, error) {
	if f.createErr != nil {
		return character.Character{}, f.createErr
	}
	c := character.Character{
		Name:   in.Name,
		Type:   character.CharacterType(in.Type),
		Gender: character.Gender(in.Gender),
	}
	for _, s := range in.Systems {
		c.PowerSystems = append(c.PowerSystems, progression.PowerSystem{Name: s})
	}
	f.created = &c
	return c, nil
}

func (f *fakeCharacterService) MainCharacter(context.Context) (character.Character, error) {
	return f.main, f.mainErr
}

func (f *fakeCharacterService) Character(_ context.Context, name string) (character.Character, error) {
	c, ok := f.byName[name]
	if !ok {
		return character.Character{}, port.ErrCharacterNotFound
	}
	return c, nil
}

func (f *fakeCharacterService) ListCharacters(context.Context) ([]character.Character, error) {
	return f.list, f.listErr
}

func (f *fakeCharacterService) GiveItem(context.Context, port.GiveItemInput) (character.Character, error) {
	return character.Character{}, nil
}

func (f *fakeCharacterService) Cultivate(_ context.Context, in port.CultivateInput) (character.Character, error) {
	f.cultivated = &in
	return f.main, f.cultErr
}

func (f *fakeCharacterService) Train(_ context.Context, in port.TrainInput) (character.Character, error) {
	f.trained = &in
	return f.main, f.trainErr
}

// fakePowerSystemService is a port.PowerSystemService double for driving the CLI.
type fakePowerSystemService struct {
	createErr   error
	addPowerErr error
	added       []port.AddPowerInput
}

func (f *fakePowerSystemService) CreateSystem(_ context.Context, name string) (progression.PowerSystem, error) {
	if f.createErr != nil {
		return progression.PowerSystem{}, f.createErr
	}
	return progression.PowerSystem{Name: name}, nil
}

func (f *fakePowerSystemService) AddPower(_ context.Context, in port.AddPowerInput) (progression.PowerSystem, error) {
	if f.addPowerErr != nil {
		return progression.PowerSystem{}, f.addPowerErr
	}
	f.added = append(f.added, in)
	return progression.PowerSystem{Name: in.System}, nil
}

func (f *fakePowerSystemService) GetSystem(context.Context, string) (progression.PowerSystem, error) {
	return progression.PowerSystem{}, nil
}

func (f *fakePowerSystemService) ListSystems(context.Context) ([]progression.PowerSystem, error) {
	return nil, nil
}

// fakeNovelService is a port.NovelService double for driving the CLI.
type fakeNovelService struct {
	created   *novel.Novel
	createErr error
	volumes   []port.AddVolumeInput
	chapters  []port.AddChapterInput
	addErr    error
	getResult novel.Novel
	getErr    error
	list      []novel.Novel
	listErr   error
}

func (f *fakeNovelService) CreateNovel(_ context.Context, in port.CreateNovelInput) (novel.Novel, error) {
	if f.createErr != nil {
		return novel.Novel{}, f.createErr
	}
	n := novel.Novel{Title: in.Title, MainCharacter: in.MainCharacter}
	f.created = &n
	return n, nil
}

func (f *fakeNovelService) AddVolume(_ context.Context, in port.AddVolumeInput) (novel.Novel, error) {
	if f.addErr != nil {
		return novel.Novel{}, f.addErr
	}
	f.volumes = append(f.volumes, in)
	return novel.Novel{Title: in.Novel}, nil
}

func (f *fakeNovelService) AddChapter(_ context.Context, in port.AddChapterInput) (novel.Novel, error) {
	if f.addErr != nil {
		return novel.Novel{}, f.addErr
	}
	f.chapters = append(f.chapters, in)
	return novel.Novel{Title: in.Novel}, nil
}

func (f *fakeNovelService) GetNovel(context.Context, string) (novel.Novel, error) {
	return f.getResult, f.getErr
}

func (f *fakeNovelService) ListNovels(context.Context) ([]novel.Novel, error) {
	return f.list, f.listErr
}

// fakeUniverseService is a port.UniverseService double for driving the CLI.
type fakeUniverseService struct {
	createErr error
	added     []port.AddUniverseSystemInput
	addErr    error
	getResult cosmology.Universe
	getErr    error
	list      []cosmology.Universe
}

func (f *fakeUniverseService) CreateUniverse(_ context.Context, name string) (cosmology.Universe, error) {
	if f.createErr != nil {
		return cosmology.Universe{}, f.createErr
	}
	return cosmology.Universe{Name: name}, nil
}

func (f *fakeUniverseService) AddSystem(_ context.Context, in port.AddUniverseSystemInput) (cosmology.Universe, error) {
	if f.addErr != nil {
		return cosmology.Universe{}, f.addErr
	}
	f.added = append(f.added, in)
	return cosmology.Universe{Name: in.Universe}, nil
}

func (f *fakeUniverseService) GetUniverse(context.Context, string) (cosmology.Universe, error) {
	return f.getResult, f.getErr
}

func (f *fakeUniverseService) ListUniverses(context.Context) ([]cosmology.Universe, error) {
	return f.list, nil
}

func run(t *testing.T, svc *fakeService, args ...string) (int, string, string) {
	t.Helper()
	return runApp(t, svc, &fakePowerSystemService{}, &fakeUniverseService{}, &fakeCharacterService{}, &fakeNovelService{}, args...)
}

func runApp(t *testing.T, realms *fakeService, systems *fakePowerSystemService, universes *fakeUniverseService, chars *fakeCharacterService, novels *fakeNovelService, args ...string) (int, string, string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := cli.New(realms, systems, universes, &fakeMultiverseService{}, &fakeOmniverseService{}, &fakeRealityService{}, &fakeTimelineService{}, chars, novels, &fakeSpeciesService{}, cli.RPGServices{}, &out, &errOut).Run(context.Background(), args)
	return code, out.String(), errOut.String()
}

type fakeMultiverseService struct{}

func (f *fakeMultiverseService) CreateMultiverse(context.Context, string) (cosmology.Multiverse, error) {
	return cosmology.Multiverse{}, nil
}
func (f *fakeMultiverseService) GetMultiverse(context.Context, string) (cosmology.Multiverse, error) {
	return cosmology.Multiverse{}, nil
}
func (f *fakeMultiverseService) ListMultiverses(context.Context) ([]cosmology.Multiverse, error) {
	return nil, nil
}
func (f *fakeMultiverseService) AddUniverse(context.Context, port.AddMultiverseUniverseInput) (cosmology.Multiverse, error) {
	return cosmology.Multiverse{}, nil
}

type fakeOmniverseService struct{}

func (f *fakeOmniverseService) CreateOmniverse(context.Context, string) (cosmology.Omniverse, error) {
	return cosmology.Omniverse{}, nil
}
func (f *fakeOmniverseService) GetOmniverse(context.Context, string) (cosmology.Omniverse, error) {
	return cosmology.Omniverse{}, nil
}
func (f *fakeOmniverseService) ListOmniverses(context.Context) ([]cosmology.Omniverse, error) {
	return nil, nil
}
func (f *fakeOmniverseService) AddMultiverse(context.Context, port.AddOmniverseMultiverseInput) (cosmology.Omniverse, error) {
	return cosmology.Omniverse{}, nil
}

type fakeRealityService struct{}

func (f *fakeRealityService) CreateReality(context.Context, string) (cosmology.Reality, error) {
	return cosmology.Reality{}, nil
}
func (f *fakeRealityService) GetReality(context.Context, string) (cosmology.Reality, error) {
	return cosmology.Reality{}, nil
}
func (f *fakeRealityService) ListRealities(context.Context) ([]cosmology.Reality, error) {
	return nil, nil
}
func (f *fakeRealityService) AddOmniverse(context.Context, port.AddRealityOmniverseInput) (cosmology.Reality, error) {
	return cosmology.Reality{}, nil
}

type fakeTimelineService struct{}

func (f *fakeTimelineService) GetTimeline(context.Context, port.LocationRef) (cosmology.Timeline, error) {
	return cosmology.Timeline{}, nil
}
func (f *fakeTimelineService) AddEvent(context.Context, port.AddTimelineEventInput) (cosmology.Timeline, error) {
	return cosmology.Timeline{}, nil
}

type fakeSpeciesService struct{}

func (f *fakeSpeciesService) CreateSpecies(context.Context, port.CreateSpeciesInput) (character.Species, error) {
	return character.Species{}, nil
}
func (f *fakeSpeciesService) Species(context.Context, string) (character.Species, error) {
	return character.Species{}, nil
}
func (f *fakeSpeciesService) ListSpecies(context.Context) ([]character.Species, error) {
	return nil, nil
}

func TestCLI_RealmAdd(t *testing.T) {
	svc := &fakeService{}
	code, _, errOut := run(t, svc,
		"realm", "add",
		"--name", "Qi Condensation",
		"--power-mult", "2", "--power-add", "10",
		"--lifespan-mult", "5", "--lifespan-add", "100",
		"--bottleneck", "250",
	)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	if len(svc.realms) != 1 || svc.realms[0].Name != "Qi Condensation" {
		t.Fatalf("realm not added: %+v", svc.realms)
	}
	if svc.realms[0].PowerMultiplier != 2 || svc.realms[0].BottleneckPoints != 250 {
		t.Errorf("flags not parsed correctly: %+v", svc.realms[0])
	}
}

func TestCLI_RealmAddInvalid(t *testing.T) {
	svc := &fakeService{}
	code, _, errOut := run(t, svc, "realm", "add", "--name", "")
	if code == 0 {
		t.Fatalf("expected non-zero exit for invalid realm")
	}
	if errOut == "" {
		t.Errorf("expected an error message on stderr")
	}
}

func TestCLI_RealmList(t *testing.T) {
	svc := &fakeService{}
	run(t, svc, "realm", "add", "--name", "Foundation")
	code, out, errOut := run(t, svc, "realm", "list")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	if !strings.Contains(out, "Foundation") {
		t.Errorf("list output missing realm name:\n%s", out)
	}
}

func TestCLI_PowerSystemAddAndAddPower(t *testing.T) {
	systems := &fakePowerSystemService{}
	code, out, errOut := runApp(t, &fakeService{}, systems, &fakeUniverseService{}, &fakeCharacterService{}, &fakeNovelService{},
		"powersystem", "add", "--name", "Universe A Cultivation")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	if !strings.Contains(out, "Universe A Cultivation") {
		t.Errorf("expected confirmation, got %q", out)
	}

	code, _, errOut = runApp(t, &fakeService{}, systems, &fakeUniverseService{}, &fakeCharacterService{}, &fakeNovelService{},
		"powersystem", "add-power", "--system", "Universe A Cultivation", "--name", "Body")
	if code != 0 {
		t.Fatalf("add-power exit = %d, stderr = %q", code, errOut)
	}
	if len(systems.added) != 1 || systems.added[0].Name != "Body" {
		t.Fatalf("power not added: %+v", systems.added)
	}
}

func TestCLI_CharacterCreate(t *testing.T) {
	chars := &fakeCharacterService{}
	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{},
		"character", "create",
		"--name", "Yun Hai",
		"--type", "Hero",
		"--gender", "Male",
		"--system", "Universe A Cultivation",
		"--system", "Universe B Sorcery",
	)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	if chars.created == nil || chars.created.Name != "Yun Hai" {
		t.Fatalf("character not created: %+v", chars.created)
	}
	if len(chars.created.PowerSystems) != 2 {
		t.Errorf("repeatable --system not parsed: %+v", chars.created.PowerSystems)
	}
	if !strings.Contains(out, "Hero") || !strings.Contains(out, "Yun Hai") {
		t.Errorf("expected confirmation output, got %q", out)
	}
}

func TestCLI_Status(t *testing.T) {
	ch := character.Character{
		Name:         "Lin Feng",
		Type:         character.MainCharacter,
		Gender:       character.Female,
		PowerValue:   "Mortal",
		PowerSystems: []progression.PowerSystem{{Name: "Universe A Cultivation"}},
		Mortal:       character.Mortal{Age: 16, Lifespan: 80},
	}

	chars := &fakeCharacterService{main: ch}
	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{}, "status")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	for _, want := range []string{"Lin Feng", "Female", "Mortal", "Universe A Cultivation", "Age 16/80"} {
		if !strings.Contains(out, want) {
			t.Errorf("status output missing %q:\n%s", want, out)
		}
	}
}

func TestCLI_StatusPowerTree(t *testing.T) {
	cult := func(realm, level string) progression.CultivationState {
		return progression.CultivationState{
			Realm: progression.Realm{Name: realm},
			Level: progression.Level{Name: level},
		}
	}
	system := progression.PowerSystem{
		Name: "Origin Cultivation",
		Kind: progression.Cultivation,
		Powers: []progression.Power{
			{Name: "Spirit", State: cult("Qi Condensation Realm", "Seventh Level")},
			{Name: "Body", State: cult("Pink Muscle Realm", "Condensation Level")},
			{Name: "Soul", State: cult("Adult Soul Realm", "Initial Stage Level")},
		},
	}
	ch := character.Character{
		Name:    "Lin Feng",
		Type:    character.MainCharacter,
		Gender:  character.Female,
		Species: []character.Species{{Name: "Human"}},
		Power:   []progression.PowerSystem{system},
		Mortal:  character.Mortal{Age: 16, Lifespan: 80},
	}

	chars := &fakeCharacterService{main: ch}
	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{}, "status")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	t.Logf("status output:\n%s", out)
	for _, want := range []string{
		"Power:",
		"- Cultivation:",
		"- Spirit:",
		"- Realm: Qi Condensation Realm",
		"- Level: Seventh Level",
		"- Body:",
		"- Realm: Pink Muscle Realm",
		"- Soul:",
		"- Level: Initial Stage Level",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("status power tree missing %q:\n%s", want, out)
		}
	}
}

func TestCLI_CharacterList(t *testing.T) {
	chars := &fakeCharacterService{list: []character.Character{
		{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female, PowerValue: "Mortal", PowerSystems: []progression.PowerSystem{{Name: "Universe A Cultivation"}}},
		{Name: "Yun Hai", Type: character.Hero, Gender: character.Male, PowerValue: "Mortal"},
	}}
	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{}, "character", "list")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	for _, want := range []string{"Lin Feng", "MainCharacter", "Yun Hai", "Hero"} {
		if !strings.Contains(out, want) {
			t.Errorf("list output missing %q:\n%s", want, out)
		}
	}
}

func TestCLI_StatusByName(t *testing.T) {
	hero := character.Character{
		Name:         "Yun Hai",
		Type:         character.Hero,
		Gender:       character.Male,
		Species:      []character.Species{{Name: "Human"}},
		PowerValue:   "Mortal",
		PowerSystems: []progression.PowerSystem{{Name: "Universe A Cultivation"}},
		Mortal:       character.Mortal{Age: 16, Lifespan: 80},
	}
	chars := &fakeCharacterService{
		main:   character.Character{Name: "Lin Feng", Type: character.MainCharacter, Gender: character.Female},
		byName: map[string]character.Character{"Yun Hai": hero},
	}

	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{}, "status", "--name", "Yun Hai")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	if !strings.Contains(out, "Yun Hai") || !strings.Contains(out, "Human") {
		t.Errorf("expected Yun Hai's status, got:\n%s", out)
	}
	if strings.Contains(out, "Lin Feng") {
		t.Errorf("should not show the main character when --name is given:\n%s", out)
	}
}

func TestCLI_StatusByNameNotFound(t *testing.T) {
	chars := &fakeCharacterService{byName: map[string]character.Character{}}
	code, _, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{}, "status", "--name", "Ghost")
	if code == 0 {
		t.Fatalf("expected non-zero exit for unknown character")
	}
	if !strings.Contains(errOut, "Ghost") {
		t.Errorf("expected error to mention the name, got %q", errOut)
	}
}

func TestCLI_StatusNoMainCharacter(t *testing.T) {
	chars := &fakeCharacterService{mainErr: port.ErrCharacterNotFound}
	code, _, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, chars, &fakeNovelService{}, "status")
	if code == 0 {
		t.Fatalf("expected non-zero exit when no main character exists")
	}
	if !strings.Contains(errOut, "no main character") {
		t.Errorf("expected guidance, got %q", errOut)
	}
}

func TestCLI_NovelAddVolumeAndChapter(t *testing.T) {
	novels := &fakeNovelService{}
	if code, _, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, &fakeCharacterService{}, novels,
		"novel", "add-volume", "--novel", "Ascension", "--number", "1", "--title", "Beginnings"); code != 0 {
		t.Fatalf("add-volume exit, stderr = %q", errOut)
	}
	if len(novels.volumes) != 1 || novels.volumes[0].Number != 1 || novels.volumes[0].Title != "Beginnings" {
		t.Fatalf("volume not forwarded: %+v", novels.volumes)
	}

	if code, _, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, &fakeCharacterService{}, novels,
		"novel", "add-chapter", "--novel", "Ascension", "--volume", "1", "--number", "3", "--title", "First Steps"); code != 0 {
		t.Fatalf("add-chapter exit, stderr = %q", errOut)
	}
	if len(novels.chapters) != 1 || novels.chapters[0].Volume != 1 || novels.chapters[0].Number != 3 {
		t.Fatalf("chapter not forwarded: %+v", novels.chapters)
	}
}

func TestCLI_NovelShow(t *testing.T) {
	n := novel.Novel{Title: "Ascension", MainCharacter: "Lin Feng", Volumes: []novel.Volume{
		{Number: 1, Title: "Beginnings", Chapters: []novel.Chapter{{Number: 1, Title: "A Mortal's Dream"}}},
		{Number: 2}, // no title -> renders "Volume 2"
	}}
	novels := &fakeNovelService{getResult: n}
	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, &fakeUniverseService{}, &fakeCharacterService{}, novels,
		"novel", "show", "--title", "Ascension")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	for _, want := range []string{"Ascension", "Lin Feng", "Volume 1: Beginnings", "Chapter 1: A Mortal's Dream", "Volume 2"} {
		if !strings.Contains(out, want) {
			t.Errorf("show output missing %q:\n%s", want, out)
		}
	}
}

func TestCLI_UniverseAddAndAddSystem(t *testing.T) {
	universes := &fakeUniverseService{}
	if code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, universes, &fakeCharacterService{}, &fakeNovelService{},
		"universe", "add", "--name", "Universe A"); code != 0 || !strings.Contains(out, "Universe A") {
		t.Fatalf("add exit = %d, out = %q, stderr = %q", code, out, errOut)
	}
	if code, _, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, universes, &fakeCharacterService{}, &fakeNovelService{},
		"universe", "add-system", "--universe", "Universe A", "--system", "Cultivation"); code != 0 {
		t.Fatalf("add-system exit, stderr = %q", errOut)
	}
	if len(universes.added) != 1 || universes.added[0].System != "Cultivation" {
		t.Fatalf("system not forwarded: %+v", universes.added)
	}
}

func TestCLI_UniverseShow(t *testing.T) {
	u := cosmology.Universe{Name: "Universe A", Systems: []progression.PowerSystem{{Name: "Cultivation"}, {Name: "Sorcery"}}}
	universes := &fakeUniverseService{getResult: u}
	code, out, errOut := runApp(t, &fakeService{}, &fakePowerSystemService{}, universes, &fakeCharacterService{}, &fakeNovelService{},
		"universe", "show", "--name", "Universe A")
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, errOut)
	}
	for _, want := range []string{"Universe A", "Cultivation", "Sorcery"} {
		if !strings.Contains(out, want) {
			t.Errorf("show output missing %q:\n%s", want, out)
		}
	}
}

func TestCLI_UnknownCommand(t *testing.T) {
	code, _, errOut := run(t, &fakeService{}, "bogus")
	if code == 0 {
		t.Fatalf("expected non-zero exit for unknown command")
	}
	if !strings.Contains(errOut, "bogus") {
		t.Errorf("stderr should mention the unknown command:\n%s", errOut)
	}
}

func (f *fakeUniverseService) AddRealms(ctx context.Context, in port.AddRealmsInput) (cosmology.Universe, error) {
	return cosmology.Universe{}, nil
}
