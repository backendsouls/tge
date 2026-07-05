import os

file_path = "internal/adapter/cli/rpg.go"
with open(file_path, "r") as f:
    content = f.read()

helper = """
func runSimpleRPGCommand[T any](
	a *App, ctx context.Context, args []string,
	entityName, pluralName string,
	create func(context.Context, string, string) (T, error),
	list func(context.Context) ([]T, error),
	get func(context.Context, string) (T, error),
	getName func(T) string,
	getDesc func(T) string,
) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(a.err, "usage: tge %s <add|list|show>\\n", entityName)
		return 2
	}
	switch args[0] {
	case "add":
		fs := flag.NewFlagSet(entityName+" add", flag.ContinueOnError)
		fs.SetOutput(a.err)
		var name, desc string
		fs.StringVar(&name, "name", "", entityName+" name")
		fs.StringVar(&desc, "description", "", entityName+" description")
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		item, err := create(ctx, name, desc)
		if err != nil {
			return a.fail(err)
		}
		_, _ = fmt.Fprintf(a.out, "added %s %q\\n", entityName, getName(item))
		return 0
	case "list":
		items, err := list(ctx)
		if err != nil {
			return a.fail(err)
		}
		if len(items) == 0 {
			_, _ = fmt.Fprintf(a.out, "no %s\\n", pluralName)
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
		_, _ = fmt.Fprintf(a.err, "unknown %s subcommand %q\\n", entityName, args[0])
		return 2
	}
}
"""

def generate_impl(name, singular, plural, port_create_input, get_name_field, get_desc_field, svc_field, svc_create, svc_list, svc_get):
    return f"""func (a *App) run{name}(ctx context.Context, args []string) int {{
	return runSimpleRPGCommand(a, ctx, args, "{singular}", "{plural}",
		func(c context.Context, n, d string) ({get_name_field}, error) {{
			return a.rpg.{svc_field}.{svc_create}(c, {port_create_input}{{Name: n, Description: d}})
		}},
		a.rpg.{svc_field}.{svc_list},
		a.rpg.{svc_field}.{svc_get},
		func(i {get_name_field}) string {{ return i.Name }},
		func(i {get_name_field}) string {{ return i.Description }},
	)
}}
"""

import re
# We need to replace runAbility, runSkill, runItem, runProfession, runClass

def replace_func(func_name, code):
    global content
    pattern = r"func \(a \*App\) " + func_name + r"\(ctx context\.Context, args \[\]string\) int \{.*?\n\}\n"
    content = re.sub(pattern, code, content, flags=re.DOTALL)

replace_func("runAbility", generate_impl("Ability", "ability", "abilities", "port.CreateAbilityInput", "rpg.Ability", "Description", "Abilities", "CreateAbility", "ListAbilities", "GetAbility"))
replace_func("runSkill", generate_impl("Skill", "skill", "skills", "port.CreateSkillInput", "rpg.Skill", "Description", "Skills", "CreateSkill", "ListSkills", "GetSkill"))
replace_func("runItem", generate_impl("Item", "item", "items", "port.CreateItemInput", "rpg.Item", "Description", "Items", "CreateItem", "ListItems", "GetItem"))
replace_func("runProfession", generate_impl("Profession", "profession", "professions", "port.CreateProfessionInput", "rpg.Profession", "Description", "Professions", "CreateProfession", "ListProfessions", "GetProfession"))
replace_func("runClass", generate_impl("Class", "class", "classes", "port.CreateClassInput", "rpg.Class", "Description", "Classes", "CreateClass", "ListClasses", "GetClass"))

# append helper to the end
content += helper

with open(file_path, "w") as f:
    f.write(content)
print("done")
