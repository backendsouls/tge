import os
import re

file_path = "internal/adapter/cli/rpg.go"
with open(file_path, "r") as f:
    content = f.read()

# Update runSimpleRPGCommand signature
content = content.replace(
    "create func(context.Context, string, string) (T, error),",
    "create func(context.Context, string, string, string) (T, error),"
)

# Update flag parsing
old_flags = """		var name, desc string
		fs.StringVar(&name, "name", "", entityName+" name")
		fs.StringVar(&desc, "description", "", entityName+" description")
		if err := fs.Parse(args[1:]); err != nil {"""
new_flags = """		var name, desc, grade string
		fs.StringVar(&name, "name", "", entityName+" name")
		fs.StringVar(&desc, "description", "", entityName+" description")
		fs.StringVar(&grade, "grade", "", entityName+" grade (optional)")
		if err := fs.Parse(args[1:]); err != nil {"""
content = content.replace(old_flags, new_flags)

# Update create call
content = content.replace(
    "item, err := create(ctx, name, desc)",
    "item, err := create(ctx, name, desc, grade)"
)

# Update the anonymous functions mapping port inputs
def update_call(name_type, input_type):
    global content
    old_func = f"""		func(c context.Context, n, d string) (rpg.{name_type}, error) {{
			return a.rpg.{name_type}s.Create{name_type}(c, port.{input_type}{{Name: n, Description: d}})
		}},"""
    new_func = f"""		func(c context.Context, n, d, g string) (rpg.{name_type}, error) {{
			return a.rpg.{name_type}s.Create{name_type}(c, port.{input_type}{{Name: n, Description: d, Grade: g}})
		}},"""
    
    # Since Class and Profession use 'es' and 's', I'll just use regex
    pattern = r'func\(c context\.Context, n, d string\) \(rpg\.' + name_type + r', error\) \{\n\s*return a\.rpg\.[a-zA-Z]+\.Create' + name_type + r'\(c, port\.' + input_type + r'\{Name: n, Description: d\}\)\n\s*\},'
    
    def repl(m):
        return m.group(0).replace('n, d string', 'n, d, g string').replace('Description: d}', 'Description: d, Grade: g}')
        
    content = re.sub(pattern, repl, content)

update_call("Ability", "CreateAbilityInput")
update_call("Skill", "CreateSkillInput")
update_call("Item", "CreateItemInput")
update_call("Profession", "CreateProfessionInput")
update_call("Class", "CreateClassInput")


with open(file_path, "w") as f:
    f.write(content)

print("done")
