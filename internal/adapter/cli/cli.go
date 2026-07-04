// Package cli is the inbound (driving) adapter: it translates command-line
// invocations into calls on the application's driving ports (realm, power system
// and character services) and renders the results.
//
// It depends only on those ports, so the whole CLI can be exercised with fake
// services and in-memory buffers, with no process or storage involved.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"tge/internal/core/domain/character"
	"tge/internal/core/domain/progression"
	"tge/internal/core/port"
	"tge/internal/logger"
)

// App dispatches CLI commands to the application services.
type App struct {
	realms       port.RealmService
	powerSystems port.PowerSystemService
	universes    port.UniverseService
	multiverses  port.MultiverseService
	omniverses   port.OmniverseService
	realities    port.RealityService
	timelines    port.TimelineService
	characters   port.CharacterService
	novels       port.NovelService
	species      port.SpeciesService
	rpg          RPGServices
	out          io.Writer
	err          io.Writer
}

// RPGServices bundles the driving ports for the RPG bounded context so the App
// constructor stays readable.
type RPGServices struct {
	Abilities   port.AbilityService
	Skills      port.SkillService
	Items       port.ItemService
	Effects     port.EffectService
	Equipment   port.EquipmentService
	Professions port.ProfessionService
	Classes     port.ClassService
	Quests      port.QuestService
	Recipes     port.RecipeService
}

// New builds an App. Services are injected as interfaces and writers are
// injected so output can be captured in tests.
func New(realms port.RealmService, powerSystems port.PowerSystemService, universes port.UniverseService, multiverses port.MultiverseService, omniverses port.OmniverseService, realities port.RealityService, timelines port.TimelineService, characters port.CharacterService, novels port.NovelService, species port.SpeciesService, rpg RPGServices, out, err io.Writer) *App {
	return &App{realms: realms, powerSystems: powerSystems, universes: universes, multiverses: multiverses, omniverses: omniverses, realities: realities, timelines: timelines, characters: characters, novels: novels, species: species, rpg: rpg, out: out, err: err}
}

// Run executes the command described by args (excluding the program name) and
// returns a process exit code: 0 on success, non-zero on failure.
func (a *App) Run(ctx context.Context, args []string) int {
	if len(args) == 0 {
		a.usage()
		return 2
	}
	switch args[0] {
	case "realm":
		return a.runRealm(ctx, args[1:])
	case "powersystem":
		return a.runPowerSystem(ctx, args[1:])
	case "universe":
		return a.runUniverse(ctx, args[1:])
	case "multiverse":
		return a.runMultiverse(ctx, args[1:])
	case "omniverse":
		return a.runOmniverse(ctx, args[1:])
	case "reality":
		return a.runReality(ctx, args[1:])
	case "timeline":
		return a.runTimeline(ctx, args[1:])
	case "ability":
		return a.runAbility(ctx, args[1:])
	case "skill":
		return a.runSkill(ctx, args[1:])
	case "item":
		return a.runItem(ctx, args[1:])
	case "effect":
		return a.runEffect(ctx, args[1:])
	case "equipment":
		return a.runEquipment(ctx, args[1:])
	case "profession":
		return a.runProfession(ctx, args[1:])
	case "class":
		return a.runClass(ctx, args[1:])
	case "quest":
		return a.runQuest(ctx, args[1:])
	case "recipe":
		return a.runRecipe(ctx, args[1:])
	case "character":
		return a.runCharacter(ctx, args[1:])
	case "novel":
		return a.runNovel(ctx, args[1:])
	case "species":
		return a.runSpecies(ctx, args[1:])
	case "status":
		return a.status(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown command %q\n", args[0])
		a.usage()
		return 2
	}
}

func (a *App) runRealm(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge realm <add|add-level|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.realmAdd(ctx, args[1:])
	case "add-level":
		return a.realmAddLevel(ctx, args[1:])
	case "list":
		return a.realmList(ctx)
	case "show":
		return a.realmShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown realm subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) realmAddLevel(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("realm add-level", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddLevelInput
	fs.StringVar(&in.Realm, "realm", "", "realm name")
	fs.IntVar(&in.Number, "number", 0, "level number (e.g. 1 for the First Level)")
	fs.StringVar(&in.Name, "name", "", "level name (e.g. \"First Level\")")
	fs.IntVar(&in.BreakthroughPoints, "breakthrough", 0, "breakthrough points to advance to the next level")
	fs.IntVar(&in.BottleneckPoints, "bottleneck", 0, "bottleneck points lost when a breakthrough fails")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.realms.AddLevel(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added level %d (%q) to realm %q\n", in.Number, in.Name, in.Realm)
	return 0
}

func (a *App) realmShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("realm show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "realm name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	r, err := a.realms.GetRealm(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "%s\n", r.Name)
	_, _ = fmt.Fprintf(a.out, "  power: %gx+%g, lifespan: %gx+%g\n",
		r.PowerMultiplier, r.PowerAdder, r.LifespanMultiplier, r.LifespanAdder)
	caps := fmt.Sprintf("normal %s, main %s", levelCap(r.MaxLevelsFor(false)), levelCap(r.MaxLevelsFor(true)))
	if len(r.Levels) == 0 {
		_, _ = fmt.Fprintf(a.out, "  Levels: none (%s)\n", caps)
		return 0
	}
	_, _ = fmt.Fprintf(a.out, "  Levels (%d defined; %s):\n", len(r.Levels), caps)
	for _, l := range r.Levels {
		_, _ = fmt.Fprintf(a.out, "    %d. %s (breakthrough %d, bottleneck %d)\n",
			l.Number, l.Name, l.BreakthroughPoints, l.BottleneckPoints)
	}
	return 0
}

// levelCap renders a level cap, where 0 means unlimited.
func levelCap(n int) string {
	if n == 0 {
		return "unlimited"
	}
	return fmt.Sprintf("%d", n)
}

func (a *App) realmAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("realm add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var cfg progression.RealmConfig
	fs.StringVar(&cfg.Name, "name", "", "realm name")
	fs.IntVar(&cfg.Tier, "tier", 0, "realm tier ordering (1 = lowest; 0 = unordered)")
	fs.Float64Var(&cfg.PowerMultiplier, "power-mult", 0, "power multiplier (a in ax+b)")
	fs.Float64Var(&cfg.PowerAdder, "power-add", 0, "power adder (b in ax+b)")
	fs.Float64Var(&cfg.LifespanMultiplier, "lifespan-mult", 0, "lifespan multiplier (a in ax+b)")
	fs.Float64Var(&cfg.LifespanAdder, "lifespan-add", 0, "lifespan adder (b in ax+b)")
	fs.IntVar(&cfg.BottleneckPoints, "bottleneck", 0, "realm-wide bottleneck points")
	fs.IntVar(&cfg.MaxLevels, "max-levels", 0, "max levels a normal character may reach (0 = unlimited)")
	fs.IntVar(&cfg.MainCharacterMaxLevels, "max-levels-main", 0, "max levels the main character may reach (0 = same as --max-levels)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	r, err := a.realms.AddRealm(ctx, cfg)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added realm %q\n", r.Name)
	return 0
}

func (a *App) realmList(ctx context.Context) int {
	realms, err := a.realms.ListRealms(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(realms) == 0 {
		_, _ = fmt.Fprintln(a.out, "no realms")
		return 0
	}

	tw := tabwriter.NewWriter(a.out, 0, 2, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "NAME\tPOWER(ax+b)\tLIFESPAN(ax+b)\tBOTTLENECK\tLEVELS\tMAX(normal/main)")
	for _, r := range realms {
		_, _ = fmt.Fprintf(tw, "%s\t%gx+%g\t%gx+%g\t%d\t%d\t%s/%s\n",
			r.Name,
			r.PowerMultiplier, r.PowerAdder,
			r.LifespanMultiplier, r.LifespanAdder,
			r.BottleneckPoints,
			len(r.Levels),
			levelCap(r.MaxLevelsFor(false)), levelCap(r.MaxLevelsFor(true)),
		)
	}
	tw.Flush()
	return 0
}

// stringSlice is a repeatable string flag (e.g. --system A --system B).
type stringSlice []string

func (s *stringSlice) String() string     { return strings.Join(*s, ",") }
func (s *stringSlice) Set(v string) error { *s = append(*s, v); return nil }

func (a *App) runCharacter(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge character <create|list|give-item|cultivate|train>")
		return 2
	}
	switch args[0] {
	case "create":
		return a.characterCreate(ctx, args[1:])
	case "list":
		return a.characterList(ctx)
	case "give-item":
		return a.characterGiveItem(ctx, args[1:])
	case "cultivate":
		return a.characterCultivate(ctx, args[1:])
	case "train":
		return a.characterTrain(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown character subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) characterGiveItem(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character give-item", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.GiveItemInput
	fs.StringVar(&in.Character, "name", "", "character name")
	fs.StringVar(&in.Item, "item", "", "item name")
	fs.IntVar(&in.Quantity, "quantity", 1, "quantity to give")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.characters.GiveItem(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "gave %d x %q to %q\n", in.Quantity, in.Item, in.Character)
	return 0
}

func (a *App) characterCultivate(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character cultivate", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CultivateInput
	var realmName string
	level := fs.Int("level", 1, "level number within the realm")
	fs.StringVar(&in.Character, "name", "", "character name")
	fs.StringVar(&in.System, "system", "", "power system (defaults to the character's first system)")
	fs.StringVar(&in.Path, "path", "", "power/path node (defaults to the system name)")
	fs.StringVar(&realmName, "realm", "", "cultivation realm (defaults to the first listed realm)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	realm, err := a.resolveRealm(ctx, realmName)
	if err != nil {
		return a.fail(err)
	}
	lvl, ok := levelByNumber(realm, *level)
	if !ok {
		return a.fail(fmt.Errorf("realm %q has no level %d", realm.Name, *level))
	}
	in.Realm = realm.Name
	in.LevelNumber = lvl.Number
	in.LevelName = lvl.Name
	in.BreakthroughPoints = lvl.BreakthroughPoints
	in.BottleneckPoints = lvl.BottleneckPoints

	c, err := a.characters.Cultivate(ctx, in)
	if err != nil {
		return a.fail(err)
	}
	logger.System("Congratulation to %s on stepping into %s, %s!", c.Name, realm.Name, lvl.Name)
	_, _ = fmt.Fprintf(a.out, "%s now cultivating in %s, %s\n", c.Name, realm.Name, lvl.Name)
	return 0
}

func (a *App) characterTrain(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character train", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.TrainInput
	fs.StringVar(&in.Character, "name", "", "character name")
	fs.StringVar(&in.System, "system", "", "power system (defaults to the character's first system)")
	fs.StringVar(&in.Path, "path", "", "power/path node (defaults to the system name)")
	fs.IntVar(&in.Points, "points", 0, "cultivation points to add")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	realms, err := a.realms.ListRealms(ctx)
	if err != nil {
		return a.fail(err)
	}
	in.Realms = realms

	c, err := a.characters.Train(ctx, in)
	if err != nil {
		return a.fail(err)
	}

	// Resolve the effective node (mirroring the service defaults) to report progress.
	system := in.System
	if system == "" && len(c.PowerSystems) > 0 {
		system = c.PowerSystems[0].Name
	}
	path := in.Path
	if path == "" {
		path = system
	}
	if cs, ok := findCultivation(c, system, path); ok {
		logger.System("Training complete! %s gained %d cultivation points.", c.Name, in.Points)
		_, _ = fmt.Fprintf(a.out, "%s trained +%d — %s, %s (breakthrough %d/%d, bottleneck %d/%d)\n",
			c.Name, in.Points, cs.Realm.Name, cs.Level.Name,
			cs.Points, cs.Level.BreakthroughPoints, cs.Bottleneck, cs.Level.BottleneckPoints)
	} else {
		logger.System("Training complete! %s gained %d cultivation points.", c.Name, in.Points)
		_, _ = fmt.Fprintf(a.out, "%s trained +%d points\n", c.Name, in.Points)
	}
	return 0
}

// findCultivation returns the CultivationState at a character's (system, path)
// node. Empty system/path match the first system and its node respectively.
func findCultivation(c character.Character, system, path string) (progression.CultivationState, bool) {
	for _, ps := range c.Power {
		if system != "" && ps.Name != system {
			continue
		}
		for _, p := range ps.Powers {
			if path != "" && p.Name != path {
				continue
			}
			if cs, ok := p.State.(progression.CultivationState); ok {
				return cs, true
			}
		}
	}
	return progression.CultivationState{}, false
}

// resolveRealm returns the named realm, or the first listed realm when name is
// empty. (The realm list is ordered by name; there is no tier ordering yet.)
func (a *App) resolveRealm(ctx context.Context, name string) (progression.Realm, error) {
	if name != "" {
		return a.realms.GetRealm(ctx, name)
	}
	realms, err := a.realms.ListRealms(ctx)
	if err != nil {
		return progression.Realm{}, err
	}
	if len(realms) == 0 {
		return progression.Realm{}, fmt.Errorf("no realms exist; add one with: tge realm add ...")
	}
	return realms[0], nil
}

// levelByNumber finds a realm's level by its number.
func levelByNumber(r progression.Realm, number int) (progression.Level, bool) {
	for _, l := range r.Levels {
		if l.Number == number {
			return l, true
		}
	}
	return progression.Level{}, false
}

func (a *App) characterCreate(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CreateCharacterInput
	var systems stringSlice
	fs.StringVar(&in.Name, "name", "", "character name")
	fs.StringVar(&in.Type, "type", "MainCharacter", "character type (MainCharacter|SideCharacter|SupportCharacter|Hero|Heroine)")
	fs.StringVar(&in.Gender, "gender", "", "character gender (Male|Female); a main character defaults to Male")
	fs.StringVar(&in.Species, "species", "", "character species; a main character defaults to the Human base")
	fs.StringVar(&in.Class, "class", "", "character class (optional; must exist)")
	fs.StringVar(&in.Profession, "profession", "", "character profession (optional; must exist)")
	fs.Var(&systems, "system", "power system the character belongs to (repeatable)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	in.Systems = systems

	c, err := a.characters.CreateCharacter(ctx, in)
	if err != nil {
		return a.fail(err)
	}
	logger.System("Host has bound to new character: %s!", c.Name)
	_, _ = fmt.Fprintf(a.out, "created %s %q (%s, %s)\n", c.Type, c.Name, c.Gender, speciesName(c.Species))
	return 0
}

func (a *App) characterList(ctx context.Context) int {
	chars, err := a.characters.ListCharacters(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(chars) == 0 {
		_, _ = fmt.Fprintln(a.out, "no characters")
		return 0
	}
	tw := tabwriter.NewWriter(a.out, 0, 2, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "NAME\tTYPE\tGENDER\tSPECIES\tPOWER\tSYSTEMS")
	for _, c := range chars {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			c.Name, c.Type, c.Gender, speciesName(c.Species), c.PowerValue, strings.Join(systemNames(c.PowerSystems), ", "))
	}
	tw.Flush()
	return 0
}

func (a *App) status(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "character to show (defaults to the main character)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	c, err := a.statusTarget(ctx, *name)
	if err != nil {
		if errors.Is(err, port.ErrCharacterNotFound) {
			if *name != "" {
				_, _ = fmt.Fprintf(a.err, "no character named %q\n", *name)
			} else {
				_, _ = fmt.Fprintln(a.err, "no main character yet; create one with: tge character create --name ... --type MainCharacter --gender ... --system ...")
			}
			return 1
		}
		return a.fail(err)
	}

	_, _ = fmt.Fprintf(a.out, "%s  (%s, %s)\n", c.Name, c.Gender, speciesName(c.Species))
	_, _ = fmt.Fprintf(a.out, "  power: %s\n", c.PowerValue)
	_, _ = fmt.Fprintf(a.out, "  Age %d/%d\n", c.Mortal.Age, c.Mortal.Lifespan)
	if c.Class.Name != "" {
		_, _ = fmt.Fprintf(a.out, "  Class: %s\n", c.Class.Name)
	}
	if c.Profession.Name != "" {
		_, _ = fmt.Fprintf(a.out, "  Profession: %s\n", c.Profession.Name)
	}
	_, _ = fmt.Fprintf(a.out, "  Stats: %s\n", formatStats(c.Stats))
	systems := systemNames(c.PowerSystems)
	if len(systems) == 0 {
		_, _ = fmt.Fprintln(a.out, "  Systems: none")
	} else {
		_, _ = fmt.Fprintf(a.out, "  Systems: %s\n", strings.Join(systems, ", "))
	}
	if len(c.Power) == 0 {
		_, _ = fmt.Fprintln(a.out, "  Power: none")
	} else {
		_, _ = fmt.Fprintln(a.out, "  Power:")
		for _, ps := range c.Power {
			label := string(ps.Kind)
			if label == "" {
				label = ps.Name
			}
			_, _ = fmt.Fprintf(a.out, "    - %s:\n", label)
			for _, p := range ps.Powers {
				writePowerState(a.out, p, 3)
			}
		}
	}
	if len(c.Inventory.Items) == 0 {
		_, _ = fmt.Fprintln(a.out, "  Inventory: empty")
	} else {
		_, _ = fmt.Fprintln(a.out, "  Inventory:")
		for _, st := range c.Inventory.Items {
			_, _ = fmt.Fprintf(a.out, "    %d x %s\n", st.Quantity, st.Item)
		}
	}
	return 0
}

// statusTarget returns the named character, or the main character if name is empty.
func (a *App) statusTarget(ctx context.Context, name string) (character.Character, error) {
	if name == "" {
		return a.characters.MainCharacter(ctx)
	}
	return a.characters.Character(ctx, name)
}

func systemNames(systems []progression.PowerSystem) []string {
	names := make([]string, len(systems))
	for i, s := range systems {
		names[i] = s.Name
	}
	return names
}

// speciesName renders a character's species as a comma-separated list of names.
func speciesName(species []character.Species) string {
	names := make([]string, len(species))
	for i, s := range species {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}

func (a *App) runPowerSystem(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge powersystem <add|add-power|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.powerSystemAdd(ctx, args[1:])
	case "add-power":
		return a.powerSystemAddPower(ctx, args[1:])
	case "list":
		return a.powerSystemList(ctx)
	case "show":
		return a.powerSystemShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown powersystem subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) powerSystemAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("powersystem add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "power system name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ps, err := a.powerSystems.CreateSystem(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added power system %q\n", ps.Name)
	return 0
}

func (a *App) powerSystemAddPower(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("powersystem add-power", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddPowerInput
	fs.StringVar(&in.System, "system", "", "power system to add to")
	fs.StringVar(&in.Name, "name", "", "power name")
	fs.StringVar(&in.Parent, "parent", "", "parent power name (optional; empty = root)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.powerSystems.AddPower(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added power %q to %q\n", in.Name, in.System)
	return 0
}

func (a *App) powerSystemList(ctx context.Context) int {
	systems, err := a.powerSystems.ListSystems(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(systems) == 0 {
		_, _ = fmt.Fprintln(a.out, "no power systems")
		return 0
	}
	for _, ps := range systems {
		_, _ = fmt.Fprintf(a.out, "%s (%d powers)\n", ps.Name, len(ps.Names()))
	}
	return 0
}

func (a *App) powerSystemShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("powersystem show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "power system name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	ps, err := a.powerSystems.GetSystem(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, ps.Name)
	for _, p := range ps.Powers {
		writePowerTree(a.out, p, 1)
	}
	return 0
}

// writePowerTree renders a power and its children with indentation.
func writePowerTree(w io.Writer, p progression.Power, depth int) {
	_, _ = fmt.Fprintf(w, "%s- %s\n", strings.Repeat("  ", depth), p.Name)
	for _, c := range p.Children {
		writePowerTree(w, c, depth+1)
	}
}

// writePowerState renders a character's progressed power node and its children:
// the node name, then its cultivation Realm/Level when the node carries a
// CultivationState, recursing into sub-powers.
func writePowerState(w io.Writer, p progression.Power, depth int) {
	indent := strings.Repeat("  ", depth)
	_, _ = fmt.Fprintf(w, "%s- %s:\n", indent, p.Name)
	if cs, ok := p.State.(progression.CultivationState); ok {
		_, _ = fmt.Fprintf(w, "%s  - Realm: %s\n", indent, cs.Realm.Name)
		_, _ = fmt.Fprintf(w, "%s  - Level: %s\n", indent, cs.Level.Name)
		_, _ = fmt.Fprintf(w, "%s  - Breakthrough: %d/%d\n", indent, cs.Points, cs.Level.BreakthroughPoints)
		_, _ = fmt.Fprintf(w, "%s  - Bottleneck: %d/%d\n", indent, cs.Bottleneck, cs.Level.BottleneckPoints)
	}
	for _, c := range p.Children {
		writePowerState(w, c, depth+1)
	}
}

func (a *App) runUniverse(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge universe <add|add-system|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.universeAdd(ctx, args[1:])
	case "add-system":
		return a.universeAddSystem(ctx, args[1:])
	case "list":
		return a.universeList(ctx)
	case "show":
		return a.universeShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown universe subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) universeAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("universe add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "universe name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	u, err := a.universes.CreateUniverse(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added universe %q\n", u.Name)
	return 0
}

func (a *App) universeAddSystem(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("universe add-system", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddUniverseSystemInput
	fs.StringVar(&in.Universe, "universe", "", "universe name")
	fs.StringVar(&in.System, "system", "", "power system to add")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.universes.AddSystem(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added power system %q to universe %q\n", in.System, in.Universe)
	return 0
}

func (a *App) universeList(ctx context.Context) int {
	universes, err := a.universes.ListUniverses(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(universes) == 0 {
		_, _ = fmt.Fprintln(a.out, "no universes")
		return 0
	}
	for _, u := range universes {
		_, _ = fmt.Fprintf(a.out, "%s (%d power systems)\n", u.Name, len(u.Systems))
	}
	return 0
}

func (a *App) universeShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("universe show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "universe name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	u, err := a.universes.GetUniverse(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, u.Name)
	for _, s := range u.Systems {
		_, _ = fmt.Fprintf(a.out, "  - %s\n", s.Name)
	}
	return 0
}

func (a *App) runNovel(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge novel <add|list>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.novelAdd(ctx, args[1:])
	case "add-volume":
		return a.novelAddVolume(ctx, args[1:])
	case "add-chapter":
		return a.novelAddChapter(ctx, args[1:])
	case "list":
		return a.novelList(ctx)
	case "show":
		return a.novelShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown novel subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) novelAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("novel add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CreateNovelInput
	fs.StringVar(&in.Title, "title", "", "novel title")
	fs.StringVar(&in.MainCharacter, "main-character", "", "name of the main character")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	n, err := a.novels.CreateNovel(ctx, in)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added novel %q (main character %q)\n", n.Title, n.MainCharacter)
	return 0
}

func (a *App) novelList(ctx context.Context) int {
	novels, err := a.novels.ListNovels(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(novels) == 0 {
		_, _ = fmt.Fprintln(a.out, "no novels")
		return 0
	}
	tw := tabwriter.NewWriter(a.out, 0, 2, 2, ' ', 0)
	for _, n := range novels {
		_, _ = fmt.Fprintf(tw, "%s\t-> %s\n", n.Title, a.describeCharacter(ctx, n.MainCharacter))
	}
	tw.Flush()
	return 0
}

func (a *App) novelAddVolume(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("novel add-volume", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddVolumeInput
	fs.StringVar(&in.Novel, "novel", "", "novel title")
	fs.IntVar(&in.Number, "number", 0, "volume number")
	fs.StringVar(&in.Title, "title", "", "volume title (optional)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.novels.AddVolume(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added volume %d to %q\n", in.Number, in.Novel)
	return 0
}

func (a *App) novelAddChapter(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("novel add-chapter", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddChapterInput
	fs.StringVar(&in.Novel, "novel", "", "novel title")
	fs.IntVar(&in.Volume, "volume", 0, "volume number")
	fs.IntVar(&in.Number, "number", 0, "chapter number")
	fs.StringVar(&in.Title, "title", "", "chapter title (optional)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.novels.AddChapter(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added chapter %d to volume %d of %q\n", in.Number, in.Volume, in.Novel)
	return 0
}

func (a *App) novelShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("novel show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	title := fs.String("title", "", "novel title")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	n, err := a.novels.GetNovel(ctx, *title)
	if err != nil {
		return a.fail(err)
	}

	_, _ = fmt.Fprintf(a.out, "%s  (main character: %s)\n", n.Title, n.MainCharacter)
	for _, v := range n.Volumes {
		_, _ = fmt.Fprintf(a.out, "  %s\n", heading("Volume", v.Number, v.Title))
		for _, c := range v.Chapters {
			_, _ = fmt.Fprintf(a.out, "    %s\n", heading("Chapter", c.Number, c.Title))
		}
	}
	return 0
}

// heading renders "<kind> <number>" with an optional ": <title>" suffix.
func heading(kind string, number int, title string) string {
	if title == "" {
		return fmt.Sprintf("%s %d", kind, number)
	}
	return fmt.Sprintf("%s %d: %s", kind, number, title)
}

// describeCharacter renders "Name (Gender)" for display, falling back to just the
// name if the character can't be loaded.
func (a *App) describeCharacter(ctx context.Context, name string) string {
	c, err := a.characters.Character(ctx, name)
	if err != nil {
		return name
	}
	return fmt.Sprintf("%s (%s)", c.Name, c.Gender)
}

func (a *App) usage() {
	_, _ = fmt.Fprintln(a.err, "usage: tge <realm|powersystem|universe|multiverse|omniverse|reality|timeline|character|novel|species|ability|skill|item|effect|equipment|profession|class|quest|recipe|status> ...")
}

func (a *App) runMultiverse(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge multiverse <add|add-universe|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.multiverseAdd(ctx, args[1:])
	case "add-universe":
		return a.multiverseAddUniverse(ctx, args[1:])
	case "list":
		return a.multiverseList(ctx)
	case "show":
		return a.multiverseShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown multiverse subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) multiverseAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("multiverse add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "multiverse name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	m, err := a.multiverses.CreateMultiverse(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added multiverse %q\n", m.Name)
	return 0
}

func (a *App) multiverseAddUniverse(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("multiverse add-universe", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddMultiverseUniverseInput
	fs.StringVar(&in.Multiverse, "multiverse", "", "multiverse name")
	fs.StringVar(&in.Universe, "universe", "", "universe to add")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.multiverses.AddUniverse(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added universe %q to multiverse %q\n", in.Universe, in.Multiverse)
	return 0
}

func (a *App) multiverseList(ctx context.Context) int {
	multiverses, err := a.multiverses.ListMultiverses(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(multiverses) == 0 {
		_, _ = fmt.Fprintln(a.out, "no multiverses")
		return 0
	}
	for _, m := range multiverses {
		_, _ = fmt.Fprintf(a.out, "%s (%d universes)\n", m.Name, len(m.Universes))
	}
	return 0
}

func (a *App) multiverseShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("multiverse show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "multiverse name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	m, err := a.multiverses.GetMultiverse(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, m.Name)
	for _, u := range m.Universes {
		_, _ = fmt.Fprintf(a.out, "  - %s\n", u.Name)
	}
	return 0
}

func (a *App) runOmniverse(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge omniverse <add|add-multiverse|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.omniverseAdd(ctx, args[1:])
	case "add-multiverse":
		return a.omniverseAddMultiverse(ctx, args[1:])
	case "list":
		return a.omniverseList(ctx)
	case "show":
		return a.omniverseShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown omniverse subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) omniverseAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("omniverse add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "omniverse name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	o, err := a.omniverses.CreateOmniverse(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added omniverse %q\n", o.Name)
	return 0
}

func (a *App) omniverseAddMultiverse(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("omniverse add-multiverse", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddOmniverseMultiverseInput
	fs.StringVar(&in.Omniverse, "omniverse", "", "omniverse name")
	fs.StringVar(&in.Multiverse, "multiverse", "", "multiverse to add")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.omniverses.AddMultiverse(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added multiverse %q to omniverse %q\n", in.Multiverse, in.Omniverse)
	return 0
}

func (a *App) omniverseList(ctx context.Context) int {
	omniverses, err := a.omniverses.ListOmniverses(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(omniverses) == 0 {
		_, _ = fmt.Fprintln(a.out, "no omniverses")
		return 0
	}
	for _, o := range omniverses {
		_, _ = fmt.Fprintf(a.out, "%s (%d multiverses)\n", o.Name, len(o.Multiverses))
	}
	return 0
}

func (a *App) omniverseShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("omniverse show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "omniverse name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	o, err := a.omniverses.GetOmniverse(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, o.Name)
	for _, m := range o.Multiverses {
		_, _ = fmt.Fprintf(a.out, "  - %s\n", m.Name)
	}
	return 0
}

func (a *App) runReality(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge reality <add|add-omniverse|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.realityAdd(ctx, args[1:])
	case "add-omniverse":
		return a.realityAddOmniverse(ctx, args[1:])
	case "list":
		return a.realityList(ctx)
	case "show":
		return a.realityShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown reality subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) realityAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("reality add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "reality name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	r, err := a.realities.CreateReality(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added reality %q\n", r.Name)
	return 0
}

func (a *App) realityAddOmniverse(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("reality add-omniverse", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddRealityOmniverseInput
	fs.StringVar(&in.Reality, "reality", "", "reality name")
	fs.StringVar(&in.Omniverse, "omniverse", "", "omniverse to add")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if _, err := a.realities.AddOmniverse(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added omniverse %q to reality %q\n", in.Omniverse, in.Reality)
	return 0
}

func (a *App) realityList(ctx context.Context) int {
	realities, err := a.realities.ListRealities(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(realities) == 0 {
		_, _ = fmt.Fprintln(a.out, "no realities")
		return 0
	}
	for _, r := range realities {
		_, _ = fmt.Fprintf(a.out, "%s (%d omniverses)\n", r.Name, len(r.Omniverses))
	}
	return 0
}

func (a *App) realityShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("reality show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "reality name")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	r, err := a.realities.GetReality(ctx, *name)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, r.Name)
	for _, o := range r.Omniverses {
		_, _ = fmt.Fprintf(a.out, "  - %s\n", o.Name)
	}
	return 0
}

func (a *App) runTimeline(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge timeline <show|add-event>")
		return 2
	}
	switch args[0] {
	case "show":
		return a.timelineShow(ctx, args[1:])
	case "add-event":
		return a.timelineAddEvent(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown timeline subcommand %q\n", args[0])
		return 2
	}
}

// locationRef builds the reference to a timeline-owning location from flags,
// reporting a usage error and returning false when the inputs are invalid.
func (a *App) locationRef(kind, name, universe string) (port.LocationRef, bool) {
	k := port.LocationKind(kind)
	if !k.Valid() {
		_, _ = fmt.Fprintf(a.err, "invalid location kind %q (want box|omniverse|multiverse|universe|realm)\n", kind)
		return port.LocationRef{}, false
	}
	if strings.TrimSpace(name) == "" {
		_, _ = fmt.Fprintln(a.err, "a location name is required (--name)")
		return port.LocationRef{}, false
	}
	if k == port.LocationRealm && strings.TrimSpace(universe) == "" {
		_, _ = fmt.Fprintln(a.err, "a realm timeline requires its universe (--universe)")
		return port.LocationRef{}, false
	}
	return port.LocationRef{Kind: k, Name: name, Universe: universe}, true
}

func (a *App) timelineShow(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("timeline show", flag.ContinueOnError)
	fs.SetOutput(a.err)
	kind := fs.String("kind", "", "location kind (box|omniverse|multiverse|universe|realm)")
	name := fs.String("name", "", "location name")
	universe := fs.String("universe", "", "owning universe (required for kind=realm)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	owner, ok := a.locationRef(*kind, *name, *universe)
	if !ok {
		return 2
	}
	t, err := a.timelines.GetTimeline(ctx, owner)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, t.Name)
	if len(t.Events) == 0 {
		_, _ = fmt.Fprintln(a.out, "  (no events)")
		return 0
	}
	for _, e := range t.Events {
		_, _ = fmt.Fprintf(a.out, "  %d. %s\n", e.Order, e.Description)
	}
	return 0
}

func (a *App) timelineAddEvent(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("timeline add-event", flag.ContinueOnError)
	fs.SetOutput(a.err)
	kind := fs.String("kind", "", "location kind (box|omniverse|multiverse|universe|realm)")
	name := fs.String("name", "", "location name")
	universe := fs.String("universe", "", "owning universe (required for kind=realm)")
	order := fs.Int("order", 0, "event order")
	description := fs.String("description", "", "event description")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	owner, ok := a.locationRef(*kind, *name, *universe)
	if !ok {
		return 2
	}
	in := port.AddTimelineEventInput{Owner: owner, Order: *order, Description: *description}
	if _, err := a.timelines.AddEvent(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added event %d to %s %q timeline\n", *order, owner.Kind, owner.Name)
	return 0
}

func (a *App) runSpecies(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(a.err, "usage: tge species <add|list>")
		return 2
	}
	switch args[0] {
	case "add":
		return a.speciesAdd(ctx, args[1:])
	case "list":
		return a.speciesList(ctx)
	default:
		_, _ = fmt.Fprintf(a.err, "unknown species subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) speciesAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("species add", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CreateSpeciesInput
	fs.StringVar(&in.Name, "name", "", "species name")
	fs.Float64Var(&in.Power, "power", 0, "base power")
	fs.IntVar(&in.Lifespan, "lifespan", 0, "base lifespan")
	fs.StringVar(&in.DefaultGender, "default-gender", "", "default gender for this species (Male|Female)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	s, err := a.species.CreateSpecies(ctx, in)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added species %q\n", s.Name)
	return 0
}

func (a *App) speciesList(ctx context.Context) int {
	list, err := a.species.ListSpecies(ctx)
	if err != nil {
		return a.fail(err)
	}
	if len(list) == 0 {
		_, _ = fmt.Fprintln(a.out, "no species")
		return 0
	}
	tw := tabwriter.NewWriter(a.out, 0, 2, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "NAME\tPOWER\tLIFESPAN\tDEFAULT-GENDER")
	for _, s := range list {
		gender := string(s.DefaultGender)
		if gender == "" {
			gender = "-"
		}
		_, _ = fmt.Fprintf(tw, "%s\t%g\t%d\t%s\n", s.Name, s.Power, s.Lifespan, gender)
	}
	tw.Flush()
	return 0
}
