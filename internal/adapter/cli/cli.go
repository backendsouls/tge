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
	"tge/internal/core/domain/cultivation"
	"tge/internal/core/domain/powersystem"
	"tge/internal/core/port"
	"tge/internal/logger"
)

var errNoArgs = errors.New("no arguments provided")

// parseFlags intercepts "help" as the first argument and converts it to a flag usage request.
func parseFlags(fs *flag.FlagSet, args []string) error {
	taken := make(map[string]bool)
	fs.VisitAll(func(f *flag.Flag) {
		taken[f.Name] = true
	})

	var toAdd []struct {
		name, usage string
		value       flag.Value
	}

	fs.VisitAll(func(f *flag.Flag) {
		if len(f.Name) == 1 || strings.HasSuffix(f.Usage, "(shorthand)") {
			return
		}
		short := ""
		for _, r := range f.Name {
			c := string(r)
			if c >= "a" && c <= "z" && !taken[c] {
				short = c
				taken[c] = true
				break
			}
		}
		if short != "" {
			toAdd = append(toAdd, struct {
				name, usage string
				value       flag.Value
			}{
				name:  short,
				usage: f.Usage + " (shorthand)",
				value: f.Value,
			})
		}
	})

	for _, a := range toAdd {
		fs.Var(a.value, a.name, a.usage)
	}

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		var shortFlags = make(map[flag.Value]*flag.Flag)
		var orderedLong []*flag.Flag
		fs.VisitAll(func(f *flag.Flag) {
			if strings.HasSuffix(f.Usage, "(shorthand)") {
				shortFlags[f.Value] = f
			} else {
				orderedLong = append(orderedLong, f)
			}
		})

		printFlag := func(f *flag.Flag, short *flag.Flag) {
			typ := fmt.Sprintf("%T", f.Value)
			typeStr := " string"
			if strings.Contains(typ, "intValue") || strings.Contains(typ, "int64Value") {
				typeStr = " int"
			} else if strings.Contains(typ, "float64Value") {
				typeStr = " float64"
			} else if strings.Contains(typ, "boolValue") {
				typeStr = ""
			}

			usage := strings.TrimSuffix(f.Usage, " (shorthand)")

			var prefix string
			if short != nil {
				prefix = fmt.Sprintf("  -%s, --%s%s", short.Name, f.Name, typeStr)
			} else {
				if len(f.Name) == 1 {
					prefix = fmt.Sprintf("  -%s%s", f.Name, typeStr)
				} else {
					prefix = fmt.Sprintf("  --%s%s", f.Name, typeStr)
				}
			}

			fmt.Fprintf(fs.Output(), "%s\n", prefix)
			fmt.Fprintf(fs.Output(), "    \t%s\n", usage)
		}

		for _, f := range orderedLong {
			printFlag(f, shortFlags[f.Value])
		}
	}

	if len(args) > 0 && (args[0] == "help" || args[0] == "-h" || args[0] == "--help") {
		fs.Usage()
		return flag.ErrHelp
	}
	if len(args) == 0 {
		fs.Usage()
		return errNoArgs
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			name := arg[1:]
			if strings.Contains(name, "=") {
				name = strings.SplitN(name, "=", 2)[0]
			}
			if len(name) > 1 {
				err := fmt.Errorf("invalid flag format %q: long flags must use double dashes (e.g. --%s)", arg, name)
				fmt.Fprintln(fs.Output(), err)
				fs.Usage()
				return err
			}
		}
	}
	return fs.Parse(args)
}

func isHelpErr(err error) bool {
	return errors.Is(err, flag.ErrHelp)
}

// noArgsHelp checks if the user requested help for a command with no arguments.
func noArgsHelp(out io.Writer, args []string, cmdName string) bool {
	if len(args) > 0 && (args[0] == "help" || args[0] == "-h" || args[0] == "--help") {
		_, _ = fmt.Fprintf(out, "usage: tge %s\n\nThis command takes no arguments or options.\n", cmdName)
		return true
	}
	return false
}

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
		a.usage(a.err)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		a.usage(a.out)
		return 0
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
		a.usage(a.err)
		return 2
	}
}

const realmHelp = `tge realm — Manage cultivation realms

Usage:
  tge realm <command> [arguments]

Commands:
  create       Create a new cultivation realm
  create-level Create a level to a realm
  list         List all realms
  show         Show details of a realm
`

func (a *App) runRealm(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, realmHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, realmHelp)
		return 0
	case "create":
		return a.realmAdd(ctx, args[1:])
	case "create-level":
		return a.realmAddLevel(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "realm list") {
			return 0
		}
		return a.realmList(ctx)
	case "show":
		return a.realmShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown realm subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) realmAddLevel(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("realm create-level", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddLevelInput
	fs.StringVar(&in.Realm, "realm", "", "realm name")
	fs.IntVar(&in.Number, "number", 0, "level number (e.g. 1 for the First Level)")
	fs.StringVar(&in.Name, "name", "", "level name (e.g. \"First Level\")")
	fs.StringVar(&in.Name, "n", "", "level name (shorthand)")
	fs.IntVar(&in.BreakthroughPoints, "breakthrough", 0, "breakthrough points to advance to the next level")
	fs.IntVar(&in.BottleneckPoints, "bottleneck", 0, "bottleneck points lost when a breakthrough fails")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs.StringVar(name, "n", "", "realm name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("realm create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var cfg cultivation.RealmConfig
	fs.StringVar(&cfg.Name, "name", "", "realm name")
	fs.StringVar(&cfg.Name, "n", "", "realm name (shorthand)")
	fs.IntVar(&cfg.Tier, "tier", 0, "realm tier ordering (1 = lowest; 0 = unordered)")
	fs.Float64Var(&cfg.PowerMultiplier, "power-mult", 0, "power multiplier (a in ax+b)")
	fs.Float64Var(&cfg.PowerAdder, "power-add", 0, "power adder (b in ax+b)")
	fs.Float64Var(&cfg.LifespanMultiplier, "lifespan-mult", 0, "lifespan multiplier (a in ax+b)")
	fs.Float64Var(&cfg.LifespanAdder, "lifespan-add", 0, "lifespan adder (b in ax+b)")
	fs.IntVar(&cfg.BottleneckPoints, "bottleneck", 0, "realm-wide bottleneck points")
	fs.IntVar(&cfg.MaxLevels, "max-levels", 0, "max levels a normal character may reach (0 = unlimited)")
	fs.IntVar(&cfg.MainCharacterMaxLevels, "max-levels-main", 0, "max levels the main character may reach (0 = same as --max-levels)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const characterHelp = `tge character — Manage characters

Usage:
  tge character <command> [arguments]

Commands:
  create       Create a new character
  list         List all characters
  clean        Soft delete all characters
  give-item    Give an item to a character
  train-node   Train a cultivation node
  idle         Assign an idle activity
  pass-time    Pass time for a character
`

func (a *App) runCharacter(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, characterHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, characterHelp)
		return 0
	case "create":
		return a.characterCreate(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "character list") {
			return 0
		}
		return a.characterList(ctx)
	case "clean":
		if noArgsHelp(a.out, args[1:], "character clean") {
			return 0
		}
		return a.characterClean(ctx)
	case "give-item":
		return a.characterGiveItem(ctx, args[1:])
	case "train-node":
		return a.characterTrainNode(ctx, args[1:])
	case "idle":
		return a.characterIdle(ctx, args[1:])
	case "pass-time":
		return a.characterPassTime(ctx, args[1:])
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
	fs.StringVar(&in.Character, "n", "", "character name (shorthand)")
	fs.StringVar(&in.Item, "item", "", "item name")
	fs.IntVar(&in.Quantity, "quantity", 1, "quantity to give")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}
	if _, err := a.characters.GiveItem(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "gave %d x %q to %q\n", in.Quantity, in.Item, in.Character)
	return 0
}

func (a *App) characterIdle(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character idle", flag.ContinueOnError)
	fs.SetOutput(a.err)
	charName := fs.String("name", "", "character name")
	fs.StringVar(charName, "n", "", "character name (shorthand)")
	activity := fs.String("activity", "", "activity to assign (e.g. 'rest', 'secluded_cultivation', 'none')")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}
	if *charName == "" {
		_, _ = fmt.Fprintln(a.err, "character name is required")
		return 2
	}
	c, err := a.characters.AssignIdleActivity(ctx, *charName, *activity)
	if err != nil {
		return a.fail(err)
	}
	if *activity == "" {
		*activity = "none"
	}
	_, _ = fmt.Fprintf(a.out, "assigned activity %q to character %q\n", *activity, c.Name)
	return 0
}

func (a *App) characterPassTime(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character pass-time", flag.ContinueOnError)
	fs.SetOutput(a.err)
	charName := fs.String("name", "", "character name")
	fs.StringVar(charName, "n", "", "character name (shorthand)")
	days := fs.Int64("days", 0, "days to pass")
	hours := fs.Int64("hours", 0, "hours to pass")
	minutes := fs.Int64("minutes", 0, "minutes to pass")

	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}
	if *charName == "" {
		_, _ = fmt.Fprintln(a.err, "character name is required")
		return 2
	}

	seconds := (*days * 86400) + (*hours * 3600) + (*minutes * 60)
	if seconds <= 0 {
		_, _ = fmt.Fprintln(a.err, "time to pass must be positive")
		return 2
	}

	c, err := a.characters.PassTime(ctx, *charName, seconds)
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "passed %d seconds for character %q. Current NovelTime: %d\n", seconds, c.Name, c.NovelTime)
	return 0
}

// findCultivation returns the CultivationState at a character's (system, path)
// node. Empty system/path match the first system and its node respectively.

func (a *App) characterCreate(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CreateCharacterInput
	var systems stringSlice
	fs.StringVar(&in.Name, "name", "", "character name")
	fs.StringVar(&in.Name, "n", "", "character name (shorthand)")
	fs.StringVar(&in.Type, "type", "MainCharacter", "character type (MainCharacter|SideCharacter|SupportCharacter|Hero|Heroine)")
	fs.StringVar(&in.Gender, "gender", "", "character gender (Male|Female); a main character defaults to Male")
	fs.StringVar(&in.Species, "species", "", "character species; a main character defaults to the Human base")
	fs.StringVar(&in.Class, "class", "", "character class (optional; must exist)")
	fs.StringVar(&in.Profession, "profession", "", "character profession (optional; must exist)")
	fs.Var(&systems, "system", "power system the character belongs to (repeatable)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}

	c, err := a.characters.CreateCharacter(ctx, in)
	if err != nil {
		return a.fail(err)
	}
	fmt.Fprintln(a.out, "--- Novel Log")
	logger.System("System binding complete. Welcome, Host %s!", c.Name)
	fmt.Fprintln(a.out, "")
	fmt.Fprintln(a.out, "--- Internal System Log")
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
	_, _ = fmt.Fprintln(tw, "NAME\tTYPE\tGENDER\tSPECIES\tPOWER")
	for _, c := range chars {
		_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			c.Name, c.Type, c.Gender, speciesName(c.Species), c.PowerValue)
	}
	tw.Flush()
	return 0
}

func (a *App) status(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "character to show (defaults to the main character)")
	fs.StringVar(name, "n", "", "character to show (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	_, _ = fmt.Fprintf(a.out, "  Power: %s\n", c.PowerValue)
	_, _ = fmt.Fprintf(a.out, "  Age %d/%d\n", c.Mortal.Age, c.Mortal.Lifespan)
	if c.Class.Name != "" {
		_, _ = fmt.Fprintf(a.out, "  Class: %s\n", c.Class.Name)
	}
	if c.Profession.Name != "" {
		_, _ = fmt.Fprintf(a.out, "  Profession: %s\n", c.Profession.Name)
	}
	if c.IdleState.ActiveActivity != "" && c.IdleState.ActiveActivity != "none" {
		_, _ = fmt.Fprintf(a.out, "  Activity: %s\n", c.IdleState.ActiveActivity)
	}
	pools := c.CurrentEnergyPools(c.NovelTime)
	if len(pools) > 0 {
		_, _ = fmt.Fprintf(a.out, "  Pools:\n")
		for k, v := range pools {
			_, _ = fmt.Fprintf(a.out, "    %s: %d\n", k, v)
		}
	}
	_, _ = fmt.Fprintf(a.out, "  Stats:%s\n", formatStats(c.Stats))
	return 0
}

// statusTarget returns the named character, or the main character if name is empty.
func (a *App) statusTarget(ctx context.Context, name string) (character.Character, error) {
	if name == "" {
		return a.characters.MainCharacter(ctx)
	}
	return a.characters.Character(ctx, name)
}

// speciesName renders a character's species as a comma-separated list of names.
func speciesName(species []character.Species) string {
	names := make([]string, len(species))
	for i, s := range species {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}

const powerSystemHelp = `tge powersystem — Manage power systems

Usage:
  tge powersystem <command> [arguments]

Commands:
  create       Create a new power system
  add-power    Add a power node to a system
  list         List all power systems
  show         Show details of a power system
`

func (a *App) runPowerSystem(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, powerSystemHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, powerSystemHelp)
		return 0
	case "create":
		return a.powerSystemAdd(ctx, args[1:])

	case "list":
		if noArgsHelp(a.out, args[1:], "powersystem list") {
			return 0
		}
		return a.powerSystemList(ctx)
	case "show":
		return a.powerSystemShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown powersystem subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) powerSystemAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("powersystem create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "power system name")
	fs.StringVar(name, "n", "", "power system name (shorthand)")
	kind := fs.String("kind", "Cultivation", "power system kind (Cultivation, Magic, SuperPower)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}
	ps, err := a.powerSystems.CreateSystem(ctx, *name, powersystem.PowerSystemType(*kind))
	if err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added power system %q\n", ps.Name)
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
	_, _ = fmt.Fprintln(a.err, "powersystem show is temporarily disabled due to mechanic state refactor")
	return 1
}

// writePowerTree renders a power and its children with indentation.

const universeHelp = `tge universe — Manage universes

Usage:
  tge universe <command> [arguments]

Commands:
  create       Create a new universe
  create-system Create a power system to a universe
  list         List all universes
  show         Show details of a universe
`

func (a *App) runUniverse(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, universeHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, universeHelp)
		return 0
	case "create":
		return a.universeAdd(ctx, args[1:])
	case "create-system":
		return a.universeAddSystem(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "universe list") {
			return 0
		}
		return a.universeList(ctx)
	case "show":
		return a.universeShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown universe subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) universeAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("universe create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "universe name")
	fs.StringVar(name, "n", "", "universe name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("universe create-system", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddUniverseSystemInput
	fs.StringVar(&in.Universe, "universe", "", "universe name")
	fs.StringVar(&in.System, "system", "", "power system to add")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs.StringVar(name, "n", "", "universe name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const novelHelp = `tge novel — Manage novels

Usage:
  tge novel <command> [arguments]

Commands:
  create       Create a new novel
  create-volume Create a volume to a novel
  create-chapter Create a chapter to a volume
  list         List all novels
  show         Show details of a novel
`

func (a *App) runNovel(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, novelHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, novelHelp)
		return 0
	case "create":
		return a.novelAdd(ctx, args[1:])
	case "create-volume":
		return a.novelAddVolume(ctx, args[1:])
	case "create-chapter":
		return a.novelAddChapter(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "novel list") {
			return 0
		}
		return a.novelList(ctx)
	case "show":
		return a.novelShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown novel subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) novelAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("novel create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CreateNovelInput
	fs.StringVar(&in.Title, "title", "", "novel title")
	fs.StringVar(&in.MainCharacter, "main-character", "", "name of the main character")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("novel create-volume", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddVolumeInput
	fs.StringVar(&in.Novel, "novel", "", "novel title")
	fs.IntVar(&in.Number, "number", 0, "volume number")
	fs.StringVar(&in.Title, "title", "", "volume title (optional)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}
	if _, err := a.novels.AddVolume(ctx, in); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintf(a.out, "added volume %d to %q\n", in.Number, in.Novel)
	return 0
}

func (a *App) novelAddChapter(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("novel create-chapter", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddChapterInput
	fs.StringVar(&in.Novel, "novel", "", "novel title")
	fs.IntVar(&in.Volume, "volume", 0, "volume number")
	fs.IntVar(&in.Number, "number", 0, "chapter number")
	fs.StringVar(&in.Title, "title", "", "chapter title (optional)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const helpText = `tge — Cultivation Game CLI

tge is a command-line tool for authoring and managing the elements of a
cultivation-style game or web novel: characters, their power systems,
cosmology, RPG building blocks, and novels.

Usage:
  tge <command> [arguments]

Core Commands:
  character    Manage characters (create, list, give-item, train-node, idle, pass-time)
  status       Show current state and attributes of a character
  novel        Manage novels (add, add-volume, add-chapter, list, show)

World & Cosmology:
  reality      Manage realities
  omniverse    Manage omniverses
  multiverse   Manage multiverses
  universe     Manage universes
  realm        Manage cultivation realms and levels
  powersystem  Manage power systems
  species      Manage species
  timeline     Manage timelines and events

RPG Toolkit:
  class        Manage RPG classes
  profession   Manage professions
  item         Manage items
  equipment    Manage equipment
  ability      Manage abilities
  skill        Manage skills
  effect       Manage status effects
  quest        Manage quests
  recipe       Manage crafting recipes

Run 'tge <command>' with no arguments to see its sub-commands.
`

func (a *App) usage(w io.Writer) {
	_, _ = fmt.Fprint(w, helpText)
}

const multiverseHelp = `tge multiverse — Manage multiverses

Usage:
  tge multiverse <command> [arguments]

Commands:
  create       Create a new multiverse
  create-universe Create a universe to a multiverse
  list         List all multiverses
  show         Show details of a multiverse
`

func (a *App) runMultiverse(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, multiverseHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, multiverseHelp)
		return 0
	case "create":
		return a.multiverseAdd(ctx, args[1:])
	case "create-universe":
		return a.multiverseAddUniverse(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "multiverse list") {
			return 0
		}
		return a.multiverseList(ctx)
	case "show":
		return a.multiverseShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown multiverse subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) multiverseAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("multiverse create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "multiverse name")
	fs.StringVar(name, "n", "", "multiverse name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("multiverse create-universe", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddMultiverseUniverseInput
	fs.StringVar(&in.Multiverse, "multiverse", "", "multiverse name")
	fs.StringVar(&in.Universe, "universe", "", "universe to add")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs.StringVar(name, "n", "", "multiverse name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const omniverseHelp = `tge omniverse — Manage omniverses

Usage:
  tge omniverse <command> [arguments]

Commands:
  create         Create a new omniverse
  create-multiverse Create a multiverse to an omniverse
  list           List all omniverses
  show           Show details of an omniverse
`

func (a *App) runOmniverse(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, omniverseHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, omniverseHelp)
		return 0
	case "create":
		return a.omniverseAdd(ctx, args[1:])
	case "create-multiverse":
		return a.omniverseAddMultiverse(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "omniverse list") {
			return 0
		}
		return a.omniverseList(ctx)
	case "show":
		return a.omniverseShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown omniverse subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) omniverseAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("omniverse create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "omniverse name")
	fs.StringVar(name, "n", "", "omniverse name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("omniverse create-multiverse", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddOmniverseMultiverseInput
	fs.StringVar(&in.Omniverse, "omniverse", "", "omniverse name")
	fs.StringVar(&in.Multiverse, "multiverse", "", "multiverse to add")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs.StringVar(name, "n", "", "omniverse name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const realityHelp = `tge reality — Manage realities

Usage:
  tge reality <command> [arguments]

Commands:
  create         Add a new reality
  create-omniverse  Add an omniverse to a reality
  list           List all realities
  show           Show details of a reality
`

func (a *App) runReality(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, realityHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, realityHelp)
		return 0
	case "create":
		return a.realityAdd(ctx, args[1:])
	case "create-omniverse":
		return a.realityAddOmniverse(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "reality list") {
			return 0
		}
		return a.realityList(ctx)
	case "show":
		return a.realityShow(ctx, args[1:])
	default:
		_, _ = fmt.Fprintf(a.err, "unknown reality subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) realityAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("reality create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	name := fs.String("name", "", "reality name")
	fs.StringVar(name, "n", "", "reality name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("reality create-omniverse", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.AddRealityOmniverseInput
	fs.StringVar(&in.Reality, "reality", "", "reality name")
	fs.StringVar(&in.Omniverse, "omniverse", "", "omniverse to add")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs.StringVar(name, "n", "", "reality name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const timelineHelp = `tge timeline — Manage timelines

Usage:
  tge timeline <command> [arguments]

Commands:
  show         Show details of a timeline
  create-event Create an event to a timeline
`

func (a *App) runTimeline(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, timelineHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, timelineHelp)
		return 0
	case "show":
		return a.timelineShow(ctx, args[1:])
	case "create-event":
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
	fs.StringVar(name, "n", "", "location name (shorthand)")
	universe := fs.String("universe", "", "owning universe (required for kind=realm)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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
	fs := flag.NewFlagSet("timeline create-event", flag.ContinueOnError)
	fs.SetOutput(a.err)
	kind := fs.String("kind", "", "location kind (box|omniverse|multiverse|universe|realm)")
	name := fs.String("name", "", "location name")
	fs.StringVar(name, "n", "", "location name (shorthand)")
	universe := fs.String("universe", "", "owning universe (required for kind=realm)")
	order := fs.Int("order", 0, "event order")
	description := fs.String("description", "", "event description")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

const speciesHelp = `tge species — Manage species

Usage:
  tge species <command> [arguments]

Commands:
  create       Create a new species
  list         List all species
`

func (a *App) runSpecies(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, speciesHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, speciesHelp)
		return 0
	case "create":
		return a.speciesAdd(ctx, args[1:])
	case "list":
		if noArgsHelp(a.out, args[1:], "species list") {
			return 0
		}
		return a.speciesList(ctx)
	default:
		_, _ = fmt.Fprintf(a.err, "unknown species subcommand %q\n", args[0])
		return 2
	}
}

func (a *App) speciesAdd(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("species create", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.CreateSpeciesInput
	fs.StringVar(&in.Name, "name", "", "species name")
	fs.StringVar(&in.Name, "n", "", "species name (shorthand)")
	fs.Float64Var(&in.Power, "power", 0, "base power")
	fs.IntVar(&in.Lifespan, "lifespan", 0, "base lifespan")
	fs.StringVar(&in.DefaultGender, "default-gender", "", "default gender for this species (Male|Female)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
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

func (a *App) characterClean(ctx context.Context) int {
	if err := a.characters.CleanCharacters(ctx); err != nil {
		return a.fail(err)
	}
	_, _ = fmt.Fprintln(a.out, "all characters soft deleted")
	return 0
}

func (a *App) characterTrainNode(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("character train-node", flag.ContinueOnError)
	fs.SetOutput(a.err)
	var in port.TrainNodeInput
	fs.StringVar(&in.Character, "name", "", "character name")
	fs.StringVar(&in.Character, "n", "", "character name (shorthand)")
	fs.StringVar(&in.System, "system", "", "power system name")
	fs.StringVar(&in.NodeID, "node", "", "node ID to train/unlock")

	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return 0
		}
		return 2
	}

	if in.Character == "" || in.System == "" || in.NodeID == "" {
		_, _ = fmt.Fprintln(a.err, "name, system, and node are required")
		return 2
	}

	c, err := a.characters.TrainNode(ctx, in)
	if err != nil {
		return a.fail(err)
	}

	fmt.Fprintln(a.out, "--- Novel Log")
	logger.System("Training complete! %s advanced in %s.", c.Name, in.System)
	fmt.Fprintln(a.out, "")
	fmt.Fprintln(a.out, "--- Internal System Log")
	_, _ = fmt.Fprintf(a.out, "%s trained node %s in %s. Power is now %s\n", c.Name, in.NodeID, in.System, c.PowerValue)
	return 0
}
