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
	fmt.Fprintf(a.err, "error: %v\n", err)
	return 1
}

// --- Ability ---

func (a *App) runAbility(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge ability <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("ability add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateAbilityInput
		fs.StringVar(&in.Name, "name", "", "ability name")
		fs.StringVar(&in.Description, "description", "", "ability description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		ab, err := a.rpg.Abilities.CreateAbility(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added ability %q\n", ab.Name)
		return 0
	case "list":
		list, err := a.rpg.Abilities.ListAbilities(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no abilities")
			return 0
		}
		for _, ab := range list {
			fmt.Fprintln(a.out, ab.Name)
		}
		return 0
	case "show":
		name := flagName(a, "ability show", args[1:])
		if name == "" {
			return 2
		}
		ab, err := a.rpg.Abilities.GetAbility(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		a.printNamed(ab.Name, ab.Description)
		return 0
	default:
		fmt.Fprintf(a.err, "unknown ability subcommand %q\n", args[0])
		return 2
	}
}

// --- Skill ---

func (a *App) runSkill(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge skill <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("skill add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateSkillInput
		fs.StringVar(&in.Name, "name", "", "skill name")
		fs.StringVar(&in.Description, "description", "", "skill description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		sk, err := a.rpg.Skills.CreateSkill(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added skill %q\n", sk.Name)
		return 0
	case "list":
		list, err := a.rpg.Skills.ListSkills(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no skills")
			return 0
		}
		for _, sk := range list {
			fmt.Fprintln(a.out, sk.Name)
		}
		return 0
	case "show":
		name := flagName(a, "skill show", args[1:])
		if name == "" {
			return 2
		}
		sk, err := a.rpg.Skills.GetSkill(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		a.printNamed(sk.Name, sk.Description)
		return 0
	default:
		fmt.Fprintf(a.err, "unknown skill subcommand %q\n", args[0])
		return 2
	}
}

// --- Item ---

func (a *App) runItem(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge item <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("item add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateItemInput
		fs.StringVar(&in.Name, "name", "", "item name")
		fs.StringVar(&in.Description, "description", "", "item description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		it, err := a.rpg.Items.CreateItem(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added item %q\n", it.Name)
		return 0
	case "list":
		list, err := a.rpg.Items.ListItems(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no items")
			return 0
		}
		for _, it := range list {
			fmt.Fprintln(a.out, it.Name)
		}
		return 0
	case "show":
		name := flagName(a, "item show", args[1:])
		if name == "" {
			return 2
		}
		it, err := a.rpg.Items.GetItem(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		a.printNamed(it.Name, it.Description)
		return 0
	default:
		fmt.Fprintf(a.err, "unknown item subcommand %q\n", args[0])
		return 2
	}
}

// --- Profession ---

func (a *App) runProfession(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge profession <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("profession add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateProfessionInput
		fs.StringVar(&in.Name, "name", "", "profession name")
		fs.StringVar(&in.Description, "description", "", "profession description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		p, err := a.rpg.Professions.CreateProfession(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added profession %q\n", p.Name)
		return 0
	case "list":
		list, err := a.rpg.Professions.ListProfessions(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no professions")
			return 0
		}
		for _, p := range list {
			fmt.Fprintln(a.out, p.Name)
		}
		return 0
	case "show":
		name := flagName(a, "profession show", args[1:])
		if name == "" {
			return 2
		}
		p, err := a.rpg.Professions.GetProfession(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		a.printNamed(p.Name, p.Description)
		return 0
	default:
		fmt.Fprintf(a.err, "unknown profession subcommand %q\n", args[0])
		return 2
	}
}

// --- Class ---

func (a *App) runClass(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge class <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("class add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateClassInput
		fs.StringVar(&in.Name, "name", "", "class name")
		fs.StringVar(&in.Description, "description", "", "class description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		c, err := a.rpg.Classes.CreateClass(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added class %q\n", c.Name)
		return 0
	case "list":
		list, err := a.rpg.Classes.ListClasses(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no classes")
			return 0
		}
		for _, c := range list {
			fmt.Fprintln(a.out, c.Name)
		}
		return 0
	case "show":
		name := flagName(a, "class show", args[1:])
		if name == "" {
			return 2
		}
		c, err := a.rpg.Classes.GetClass(ctx, name)
		if err != nil {
			return a.fail(err)
		}
		a.printNamed(c.Name, c.Description)
		return 0
	default:
		fmt.Fprintf(a.err, "unknown class subcommand %q\n", args[0])
		return 2
	}
}

// --- Effect ---

func (a *App) runEffect(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge effect <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("effect add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateEffectInput
		fs.StringVar(&in.Name, "name", "", "effect name")
		fs.StringVar(&in.Kind, "kind", "", "effect kind (Buff|Debuff|Status)")
		fs.StringVar(&in.Description, "description", "", "effect description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		e, err := a.rpg.Effects.CreateEffect(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added %s effect %q\n", e.Kind, e.Name)
		return 0
	case "list":
		list, err := a.rpg.Effects.ListEffects(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no effects")
			return 0
		}
		for _, e := range list {
			fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Kind)
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
		fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Kind)
		if e.Description != "" {
			fmt.Fprintf(a.out, "  %s\n", e.Description)
		}
		return 0
	default:
		fmt.Fprintf(a.err, "unknown effect subcommand %q\n", args[0])
		return 2
	}
}

// --- Equipment ---

func (a *App) runEquipment(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge equipment <add|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("equipment add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateEquipmentInput
		fs.StringVar(&in.Name, "name", "", "equipment name")
		fs.StringVar(&in.Slot, "slot", "", "slot (Weapon|Armor|Accessory)")
		bindStatFlags(fs, &in.Bonus)
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		e, err := a.rpg.Equipment.CreateEquipment(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added %s %q\n", e.Slot, e.Name)
		return 0
	case "list":
		list, err := a.rpg.Equipment.ListEquipment(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no equipment")
			return 0
		}
		for _, e := range list {
			fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Slot)
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
		fmt.Fprintf(a.out, "%s (%s)\n", e.Name, e.Slot)
		fmt.Fprintf(a.out, "  bonus: %s\n", formatStats(e.Bonus))
		return 0
	default:
		fmt.Fprintf(a.err, "unknown equipment subcommand %q\n", args[0])
		return 2
	}
}

// --- Quest ---

func (a *App) runQuest(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge quest <add|add-objective|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("quest add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateQuestInput
		fs.StringVar(&in.Name, "name", "", "quest name")
		fs.StringVar(&in.Description, "description", "", "quest description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		q, err := a.rpg.Quests.CreateQuest(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added quest %q\n", q.Name)
		return 0
	case "add-objective":
		fs := flag.NewFlagSet("quest add-objective", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.AddQuestObjectiveInput
		fs.StringVar(&in.Quest, "quest", "", "quest name")
		fs.IntVar(&in.Order, "order", 0, "objective order")
		fs.StringVar(&in.Description, "description", "", "objective description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		if _, err := a.rpg.Quests.AddObjective(ctx, in); err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added objective %d to quest %q\n", in.Order, in.Quest)
		return 0
	case "list":
		list, err := a.rpg.Quests.ListQuests(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no quests")
			return 0
		}
		for _, q := range list {
			fmt.Fprintln(a.out, q.Name)
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
		fmt.Fprintln(a.out, q.Name)
		if q.Description != "" {
			fmt.Fprintf(a.out, "  %s\n", q.Description)
		}
		for _, o := range q.Objectives {
			fmt.Fprintf(a.out, "  %d. %s\n", o.Order, o.Description)
		}
		return 0
	default:
		fmt.Fprintf(a.err, "unknown quest subcommand %q\n", args[0])
		return 2
	}
}

// --- Recipe ---

func (a *App) runRecipe(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(a.err, "usage: tge recipe <add|add-input|list|show>")
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet("recipe add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.CreateRecipeInput
		fs.StringVar(&in.Name, "name", "", "recipe name")
		fs.StringVar(&in.Output, "output", "", "output item")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		r, err := a.rpg.Recipes.CreateRecipe(ctx, in)
		if err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added recipe %q -> %q\n", r.Name, r.Output)
		return 0
	case "add-input":
		fs := flag.NewFlagSet("recipe add-input", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var in port.AddIngredientInput
		fs.StringVar(&in.Recipe, "recipe", "", "recipe name")
		fs.StringVar(&in.Item, "item", "", "input item")
		fs.IntVar(&in.Quantity, "quantity", 1, "quantity required")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		if _, err := a.rpg.Recipes.AddIngredient(ctx, in); err != nil {
			return a.fail(err)
		}
		fmt.Fprintf(a.out, "added %d x %q to recipe %q\n", in.Quantity, in.Item, in.Recipe)
		return 0
	case "list":
		list, err := a.rpg.Recipes.ListRecipes(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(list) == 0 {
			fmt.Fprintln(a.out, "no recipes")
			return 0
		}
		for _, r := range list {
			fmt.Fprintf(a.out, "%s -> %s\n", r.Name, r.Output)
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
		fmt.Fprintf(a.out, "%s -> %s\n", r.Name, r.Output)
		for _, in := range r.Inputs {
			fmt.Fprintf(a.out, "  %d x %s\n", in.Quantity, in.Item)
		}
		return 0
	default:
		fmt.Fprintf(a.err, "unknown recipe subcommand %q\n", args[0])
		return 2
	}
}

// flagName parses a lone --name flag, printing usage on error and returning ""
// when missing or invalid.
func flagName(a *App, name string, args []string) string {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(a.err)
	n := fs.String("name", "", "name")
	if err := fs.Parse(args); err != nil {
		return ""
	}
	if *n == "" {
		fmt.Fprintln(a.err, "a name is required (--name)")
		return ""
	}
	return *n
}

// printNamed renders a name with an optional indented description.
func (a *App) printNamed(name, description string) {
	fmt.Fprintln(a.out, name)
	if description != "" {
		fmt.Fprintf(a.out, "  %s\n", description)
	}
}

// bindStatFlags binds the eight stat attributes to flags on fs.
func bindStatFlags(fs *flag.FlagSet, s *rpg.Stats) {
	fs.IntVar(&s.STR, "str", 0, "Strength")
	fs.IntVar(&s.AGI, "agi", 0, "Agility")
	fs.IntVar(&s.INT, "int", 0, "Intelligence")
	fs.IntVar(&s.VIT, "vit", 0, "Vitality")
	fs.IntVar(&s.DEX, "dex", 0, "Dexterity")
	fs.IntVar(&s.WIS, "wis", 0, "Wisdom")
	fs.IntVar(&s.CHA, "cha", 0, "Charisma")
	fs.IntVar(&s.LUK, "luk", 0, "Luck")
}

// formatStats renders a stat block on one line.
func formatStats(s rpg.Stats) string {
	return fmt.Sprintf("STR %d, AGI %d, INT %d, VIT %d, DEX %d, WIS %d, CHA %d, LUK %d",
		s.STR, s.AGI, s.INT, s.VIT, s.DEX, s.WIS, s.CHA, s.LUK)
}
