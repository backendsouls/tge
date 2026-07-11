import os

def replace_in_file(path, old, new):
    if not os.path.exists(path):
        return
    with open(path, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(path, 'w') as f:
            f.write(content)

# cultivation test fixes
for f in ["level_test.go", "cultivation_path_test.go", "realm_test.go"]:
    path = "internal/core/domain/cultivation/" + f
    # First, restore progression.Err -> cultivation.Err
    replace_in_file(path, 'ErrInvalidLevelName', 'cultivation.ErrInvalidLevelName')
    replace_in_file(path, 'ErrInvalidLevelPoints', 'cultivation.ErrInvalidLevelPoints')
    replace_in_file(path, 'ErrLevelNumberExists', 'cultivation.ErrLevelNumberExists')
    replace_in_file(path, 'ErrLevelNumberExceedsMax', 'cultivation.ErrLevelNumberExceedsMax')
    
    replace_in_file(path, 'cultivation.cultivation.', 'cultivation.')

def replace_power_to_powersystem(path):
    replace_in_file(path, 'power.PowerSystemType', 'powersystem.PowerSystemType')
    replace_in_file(path, 'power.PowerSystem', 'powersystem.PowerSystem')
    replace_in_file(path, 'power.PowerNode', 'powersystem.PowerNode')

for root, dirs, files in os.walk("."):
    for name in files:
        if name.endswith(".go"):
            path = os.path.join(root, name)
            replace_power_to_powersystem(path)

replace_in_file('internal/adapter/cli/cli.go', 'powersystem.CultivationState', 'cultivation.CultivationState')
replace_in_file('internal/adapter/cli/cli.go', 'powersystem.SuperPowerState', 'superpower.SuperPowerState')
replace_in_file('internal/adapter/sqlite/realm_repository.go', 'progression.', 'cultivation.')

