package cli

import (
	"context"
	"flag"
	"fmt"

	"tge/internal/core/domain/rpg"
	"tge/internal/core/port"
)

// fail reports an error and returns the conventional exit code 1.
func (a *App) fail(err error) int {
	_, _ = fmt.Fprintf(a.err, "error: %v\n", err)
	return 1
}

// --- Ability ---

func (a *App) runAbility(ctx context.Context, args []string) int {
	return runSimpleRPGCommand(a, ctx, args, "ability", "abilities",
		func(c context.Context, n, d, g string) (rpg.Ability, error) {
			return a.rpg.Abilities.CreateAbility(c, port.CreateAbilityInput{Name: n, Description: d, Grade: g})
		},
		a.rpg.Abilities.ListAbilities,
		a.rpg.Abilities.GetAbility,
		func(i rpg.Ability) string { return i.Name },
		func(i rpg.Ability) string { return i.Description },
	)
}

// --- Skill ---

func (a *App) runSkill(ctx context.Context, args []string) int {
	return runSimpleRPGCommand(a, ctx, args, "skill", "skills",
		func(c context.Context, n, d, g string) (rpg.Skill, error) {
			return a.rpg.Skills.CreateSkill(c, port.CreateSkillInput{Name: n, Description: d, Grade: g})
		},
		a.rpg.Skills.ListSkills,
		a.rpg.Skills.GetSkill,
		func(i rpg.Skill) string { return i.Name },
		func(i rpg.Skill) string { return i.Description },
	)
}

// --- Item ---

func (a *App) runItem(ctx context.Context, args []string) int {
	return runSimpleRPGCommand(a, ctx, args, "item", "items",
		func(c context.Context, n, d, g string) (rpg.Item, error) {
			return a.rpg.Items.CreateItem(c, port.CreateItemInput{Name: n, Description: d, Grade: g})
		},
		a.rpg.Items.ListItems,
		a.rpg.Items.GetItem,
		func(i rpg.Item) string { return i.Name },
		func(i rpg.Item) string { return i.Description },
	)
}

// --- Profession ---

func (a *App) runProfession(ctx context.Context, args []string) int {
	return runSimpleRPGCommand(a, ctx, args, "profession", "professions",
		func(c context.Context, n, d, g string) (rpg.Profession, error) {
			return a.rpg.Professions.CreateProfession(c, port.CreateProfessionInput{Name: n, Description: d, Grade: g})
		},
		a.rpg.Professions.ListProfessions,
		a.rpg.Professions.GetProfession,
		func(i rpg.Profession) string { return i.Name },
		func(i rpg.Profession) string { return i.Description },
	)
}

// --- Class ---

func (a *App) runClass(ctx context.Context, args []string) int {
	return runSimpleRPGCommand(a, ctx, args, "class", "classes",
		func(c context.Context, n, d, g string) (rpg.Class, error) {
			return a.rpg.Classes.CreateClass(c, port.CreateClassInput{Name: n, Description: d, Grade: g})
		},
		a.rpg.Classes.ListClasses,
		a.rpg.Classes.GetClass,
		func(i rpg.Class) string { return i.Name },
		func(i rpg.Class) string { return i.Description },
	)
}

// --- Effect ---

const effectHelp = `tge effect — Manage effects

Usage:
  tge effect <command> [arguments]

Commands:
  create       Create a new effect
  list         List all effects
  show         Show details of an effect
`

func (a *App) runEffect(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, effectHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, effectHelp)
		return 0
	case "create":
		fs := flag.NewFlagSet("effect create", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateEffectInput
		fs.StringVar(&in.Name, "name", "", "effect name")
		fs.StringVar(&in.Name, "n", "", "effect name (shorthand)")
		fs.StringVar(&in.Kind, "kind", "", "effect kind (Buff|Debuff|Status)")
		fs.StringVar(&in.Description, "description", "", "effect description")
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		e, err := a.rpg.Effects.CreateEffect(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added %s effect %q\n", e.Kind, e.Name)
		return 0
	case "list":
		if noArgsHelp(a.out, args[1:], "effect list") {
			return 0
		}
		list, err := a.rpg.Effects.ListEffects(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			_, _ = fmt.Fprintln(a.out, "no effects")
			return 0
		}
		for _, e := range list {
			_, _ = fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Kind)
		}
		return 0
	case "show":
		name := flagName(a, "effect show", args[1:])
		if name == "" {
			return 2
		}
		e, err := a.rpg.Effects.GetEffect(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Kind)
		if e.Description != "" {
			_, _ = fmt.Fprintf(a.out, "  %s\n", e.Description)
		}
		return 0
	default:
		_, _ = fmt.Fprintf(a.err, "unknown effect subcommand %q\n", args[0])
		return 2
	}
}

// --- Equipment ---

const equipmentHelp = `tge equipment — Manage equipment

Usage:
  tge equipment <command> [arguments]

Commands:
  create       Create new equipment
  list         List all equipment
  show         Show details of equipment
`

func (a *App) runEquipment(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, equipmentHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, equipmentHelp)
		return 0
	case "create":
		fs := flag.NewFlagSet("equipment create", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateEquipmentInput
		fs.StringVar(&in.Name, "name", "", "equipment name")
		fs.StringVar(&in.Name, "n", "", "equipment name (shorthand)")
		fs.StringVar(&in.Slot, "slot", "", "slot (Weapon|Armor|Accessory)")
		bindStatFlags(fs, &in.Bonus)
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		e, err := a.rpg.Equipment.CreateEquipment(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added %s %q\n", e.Slot, e.Name)
		return 0
	case "list":
		if noArgsHelp(a.out, args[1:], "equipment list") {
			return 0
		}
		list, err := a.rpg.Equipment.ListEquipment(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			_, _ = fmt.Fprintln(a.out, "no equipment")
			return 0
		}
		for _, e := range list {
			_, _ = fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Slot)
		}
		return 0
	case "show":
		name := flagName(a, "equipment show", args[1:])
		if name == "" {
			return 2
		}
		e, err := a.rpg.Equipment.GetEquipment(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Slot)
		_, _ = fmt.Fprintf(a.out, "  bonus:%s\n", formatStats(e.Bonus))
		return 0
	default:
		_, _ = fmt.Fprintf(a.err, "unknown equipment subcommand %q\n", args[0])
		return 2
	}
}

// --- Quest ---

const questHelp = `tge quest — Manage quests

Usage:
  tge quest <command> [arguments]

Commands:
  create         Create a new quest
  create-objective Create an objective to a quest
  list           List all quests
  show           Show details of a quest
`

func (a *App) runQuest(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, questHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, questHelp)
		return 0
	case "create":
		fs := flag.NewFlagSet("quest create", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateQuestInput
		fs.StringVar(&in.Name, "name", "", "quest name")
		fs.StringVar(&in.Name, "n", "", "quest name (shorthand)")
		fs.StringVar(&in.Description, "description", "", "quest description")
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		q, err := a.rpg.Quests.CreateQuest(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added quest %q\n", q.Name)
		return 0
	case "create-objective":
		fs := flag.NewFlagSet("quest create-objective", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.AddQuestObjectiveInput
		fs.StringVar(&in.Quest, "quest", "", "quest name")
		fs.IntVar(&in.Order, "order", 0, "objective order")
		fs.StringVar(&in.Description, "description", "", "objective description")
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		if _, err := a.rpg.Quests.AddObjective(ctx, in); err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added objective %d to quest %q\n", in.Order, in.Quest)
		return 0
	case "list":
		if noArgsHelp(a.out, args[1:], "quest list") {
			return 0
		}
		list, err := a.rpg.Quests.ListQuests(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			_, _ = fmt.Fprintln(a.out, "no quests")
			return 0
		}
		for _, q := range list {
			_, _ = fmt.Fprintln(a.out, q.Name)
		}
		return 0
	case "show":
		name := flagName(a, "quest show", args[1:])
		if name == "" {
			return 2
		}
		q, err := a.rpg.Quests.GetQuest(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintln(a.out, q.Name)
		if q.Description != "" {
			_, _ = fmt.Fprintf(a.out, "  %s\n", q.Description)
		}
		for _, o := range q.Objectives {
			_, _ = fmt.Fprintf(a.out, "  %d. %s\n", o.Order, o.Description)
		}
		return 0
	default:
		_, _ = fmt.Fprintf(a.err, "unknown quest subcommand %q\n", args[0])
		return 2
	}
}

// --- Recipe ---

const recipeHelp = `tge recipe — Manage recipes

Usage:
  tge recipe <command> [arguments]

Commands:
  create         Create a new recipe
  create-input Create an input item to a recipe
  list           List all recipes
  show           Show details of a recipe
`

func (a *App) runRecipe(ctx context.Context, args []string) int {
	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, recipeHelp)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, recipeHelp)
		return 0
	case "create":
		fs := flag.NewFlagSet("recipe create", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateRecipeInput
		fs.StringVar(&in.Name, "name", "", "recipe name")
		fs.StringVar(&in.Name, "n", "", "recipe name (shorthand)")
		fs.StringVar(&in.Output, "output", "", "output item")
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		r, err := a.rpg.Recipes.CreateRecipe(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added recipe %q -> %q\n", r.Name, r.Output)
		return 0
	case "create-input":
		fs := flag.NewFlagSet("recipe create-input", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.AddIngredientInput
		fs.StringVar(&in.Recipe, "recipe", "", "recipe name")
		fs.StringVar(&in.Item, "item", "", "input item")
		fs.IntVar(&in.Quantity, "quantity", 1, "quantity required")
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		if _, err := a.rpg.Recipes.AddIngredient(ctx, in); err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added %d x %q to recipe %q\n", in.Quantity, in.Item, in.Recipe)
		return 0
	case "list":
		if noArgsHelp(a.out, args[1:], "recipe list") {
			return 0
		}
		list, err := a.rpg.Recipes.ListRecipes(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			_, _ = fmt.Fprintln(a.out, "no recipes")
			return 0
		}
		for _, r := range list {
			_, _ = fmt.Fprintf(a.out, "%s -> %s\n", r.Name, r.Output)
		}
		return 0
	case "show":
		name := flagName(a, "recipe show", args[1:])
		if name == "" {
			return 2
		}
		r, err := a.rpg.Recipes.GetRecipe(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "%s -> %s\n", r.Name, r.Output)
		for _, in := range r.Inputs {
			_, _ = fmt.Fprintf(a.out, "  %d x %s\n", in.Quantity, in.Item)
		}
		return 0
	default:
		_, _ = fmt.Fprintf(a.err, "unknown recipe subcommand %q\n", args[0])
		return 2
	}
}

// flagName parses a lone --name flag, printing usage on error and returning ""
// when missing or invalid.
func flagName(a *App, name string, args []string) string {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(a.err)
	n := fs.String("name", "", "name")
	fs.StringVar(n, "n", "", "name (shorthand)")
	if err := parseFlags(fs, args); err != nil {
		if isHelpErr(err) {
			return ""
		}
		return ""
	}
	if *n == "" {
		_, _ = fmt.Fprintln(a.err, "a name is required (--name)")
		return ""
	}
	return *n
}

// printNamed renders a name with an optional indented description.
func (a *App) printNamed(name, description string) {
	_, _ = fmt.Fprintln(a.out, name)
	if description != "" {
		_, _ = fmt.Fprintf(a.out, "  %s\n", description)
	}
}

// bindStatFlags binds the eight stat attributes to flags on fs.
func bindStatFlags(fs *flag.FlagSet, s *rpg.Stats) {
	fs.Float64Var(&s.STR, "str", 0, "Strength")
	fs.Float64Var(&s.AGI, "agi", 0, "Agility")
	fs.Float64Var(&s.INT, "int", 0, "Intelligence")
	fs.Float64Var(&s.VIT, "vit", 0, "Vitality")
	fs.Float64Var(&s.DEX, "dex", 0, "Dexterity")
	fs.Float64Var(&s.WIS, "wis", 0, "Wisdom")
	fs.Float64Var(&s.CHA, "cha", 0, "Charisma")
	fs.Float64Var(&s.LUK, "luk", 0, "Luck")
}

// formatStats renders a stat block multiline.
func formatStats(s rpg.Stats) string {
	return fmt.Sprintf("\n   STR: %g\n   AGI: %g\n   INT: %g\n   VIT: %g\n   DEX: %g\n   WIS: %g\n   CHA: %g\n   LUK: %g",
		s.STR, s.AGI, s.INT, s.VIT, s.DEX, s.WIS, s.CHA, s.LUK)
}

func runSimpleRPGCommand[T any](
	a *App, ctx context.Context, args []string,
	entityName, pluralName string,
	create func(context.Context, string, string, string) (T, error),
	list func(context.Context) ([]T, error),
	get func(context.Context, string) (T, error),
	getName func(T) string,
	getDesc func(T) string,
) int {
	helpText := fmt.Sprintf(`tge %s — Manage %s

Usage:
  tge %s <command> [arguments]

Commands:
  create       Create a new %s
  list         List all %s
  show         Show details of a %s
`, entityName, pluralName, entityName, entityName, pluralName, entityName)

	if len(args) == 0 {
		_, _ = fmt.Fprint(a.err, helpText)
		return 2
	}
	switch args[0] {
	case "help", "-h", "--help":
		_, _ = fmt.Fprint(a.out, helpText)
		return 0
	case "create":
		fs := flag.NewFlagSet(entityName+" create", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var name, desc, grade string
		fs.StringVar(&name, "name", "", entityName+" name")
		fs.StringVar(&name, "n", "", entityName+" name (shorthand)")
		fs.StringVar(&desc, "description", "", entityName+" description")
		fs.StringVar(&grade, "grade", "", entityName+" grade (optional)")
		if err := parseFlags(fs, args[1:]); err != nil {
			if isHelpErr(err) {
				return 0
			}
			return 2
		}
		item, err := create(ctx, name, desc, grade)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added %s %q\n", entityName, getName(item))
		return 0
	case "list":
		if noArgsHelp(a.out, args[1:], entityName+" list") {
			return 0
		}
		items, err := list(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(items) == 0 {
			_, _ = fmt.Fprintf(a.out, "no %s\n", pluralName)
			return 0
		}
		for _, item := range items {
			_, _ = fmt.Fprintln(a.out, getName(item))
		}
		return 0
	case "show":
		name := flagName(a, entityName+" show", args[1:])
		if name == "" {
			return 2
		}
		item, err := get(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		a.printNamed(getName(item), getDesc(item))
		return 0
	default:
		_, _ = fmt.Fprintf(a.err, "unknown %s subcommand %q\n", entityName, args[0])
		return 2
	}
}
