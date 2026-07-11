import os

def replace_in_file(path, old, new):
    with open(path, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(path, 'w') as f:
            f.write(content)

replace_in_file('internal/core/port/power_system.go', 'progression.', 'powersystem.')
replace_in_file('internal/core/port/power_system.go', 'power.PowerSystem', 'powersystem.PowerSystem')

replace_in_file('internal/core/domain/cosmology/universe_test.go', 'progression.', 'powersystem.')
replace_in_file('internal/core/domain/cultivation/level_test.go', 'progression.', '') # It's in cultivation package now, so progression.ErrInvalidLevelName becomes ErrInvalidLevelName

# Let's run a generic replacement for progression.Err* in cultivation package
replace_in_file('internal/core/domain/cultivation/level_test.go', 'progression.', '')
replace_in_file('internal/core/domain/cultivation/cultivation_path_test.go', 'progression.', '')
replace_in_file('internal/core/domain/cultivation/realm_test.go', 'progression.', '')

