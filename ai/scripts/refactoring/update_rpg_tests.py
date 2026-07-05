import os
import re

files = [
    "internal/core/domain/rpg/rpg_test.go",
    "internal/adapter/sqlite/rpg_repository_test.go"
]

for file_path in files:
    with open(file_path, "r") as f:
        content = f.read()

    # Replace rpg.NewAbility("X", "Y") with rpg.NewAbility("X", "Y", "Common")
    content = re.sub(r'rpg\.NewAbility\(([^,]+), ([^,]+)\)', r'rpg.NewAbility(\1, \2, "Common")', content)
    content = re.sub(r'rpg\.NewSkill\(([^,]+), ([^,]+)\)', r'rpg.NewSkill(\1, \2, "Common")', content)
    content = re.sub(r'rpg\.NewItem\(([^,]+), ([^,]+)\)', r'rpg.NewItem(\1, \2, "Common")', content)
    content = re.sub(r'rpg\.NewClass\(([^,]+), ([^,]+)\)', r'rpg.NewClass(\1, \2, "Common")', content)
    content = re.sub(r'rpg\.NewProfession\(([^,]+), ([^,]+)\)', r'rpg.NewProfession(\1, \2, "Common")', content)

    # Some tests might pass port.CreateAbilityInput{Name: "X", Description: "Y"}
    # This might not break tests as Grade will just be empty "", which is valid, but let's see.

    with open(file_path, "w") as f:
        f.write(content)

print("done")
