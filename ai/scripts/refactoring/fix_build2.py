import re

# Fix default_world_service.go
path = "internal/core/service/default_world_service.go"
with open(path, "r") as f: content = f.read()
content = content.replace('s.systems.Create(ctx, ps)', 's.systems.Save(ctx, ps)')
with open(path, "w") as f: f.write(content)

# Fix power_system_repository.go
path = "internal/adapter/sqlite/power_system_repository.go"
with open(path, "r") as f: content = f.read()
content = content.replace('func (r *PowerSystemRepository) Create(', 'func (r *PowerSystemRepository) Save(')
content = re.sub(r'system\.Powers', 'system.Nodes', content)
content = re.sub(r'ps\.Powers', 'ps.Nodes', content)
with open(path, "w") as f: f.write(content)

# Fix cli.go
path = "internal/adapter/cli/cli.go"
with open(path, "r") as f: content = f.read()
content = content.replace('case "add-power":\n\t\treturn a.powerSystemAddPower(ctx, args[1:])', '')
content = re.sub(r'func \(a \*App\) powerSystemAddPower.*?\n}', '', content, flags=re.DOTALL)
content = content.replace('for _, p := range ps.Powers {', 'for _, p := range ps.Nodes {')
# also fix c.Power missing
content = re.sub(r'func findCultivation.*?\n}', '', content, flags=re.DOTALL)
with open(path, "w") as f: f.write(content)
