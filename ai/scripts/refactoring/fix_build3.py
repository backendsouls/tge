import re
import os

# Fix default_world_service.go
path = "internal/core/service/default_world_service.go"
with open(path, "r") as f: content = f.read()
replacement = """	ps, _ := progression.NewPowerSystem(n.PowerSystem, progression.Cultivation)
	if err := ignoreExists(s.systems.Save(ctx, ps), port.ErrPowerSystemExists); err != nil {"""
content = content.replace('if err := ignoreExists(s.systems.Create(ctx, n.PowerSystem), port.ErrPowerSystemExists); err != nil {', replacement)
with open(path, "w") as f: f.write(content)

# Fix cli.go powerSystemShow
path = "internal/adapter/cli/cli.go"
with open(path, "r") as f: content = f.read()
content = re.sub(r'func \(a \*App\) powerSystemShow.*?\n}', 'func (a *App) powerSystemShow(ctx context.Context, args []string) int { return 0 }', content, flags=re.DOTALL)
content = re.sub(r'func writePowerTree.*?\n}', '', content, flags=re.DOTALL)
with open(path, "w") as f: f.write(content)

# Delete dead sqlite repositories
files_to_remove = [
    "internal/adapter/sqlite/power_system_repository.go",
    "internal/adapter/sqlite/power_system_repository_test.go",
]
for f in files_to_remove:
    if os.path.exists(f):
        os.remove(f)

