import os
import glob

# Rename packages in the moved files
def rename_package(directory, pkg_name):
    for f in glob.glob(f"internal/core/domain/{directory}/*.go"):
        with open(f, 'r') as file:
            content = file.read()
        content = content.replace("package progression", f"package {pkg_name}")
        with open(f, 'w') as file:
            file.write(content)

rename_package("power", "power")
rename_package("powersystem", "powersystem")
rename_package("cultivation", "cultivation")
rename_package("superpower", "superpower")

# Find all go files and replace progression imports and usages
for root, dirs, files in os.walk("."):
    for name in files:
        if name.endswith(".go") and name != "fix_imports.py":
            f = os.path.join(root, name)
            with open(f, 'r') as file:
                content = file.read()
            
            orig = content
            
            # Replace import
            if '"tge/internal/core/domain/progression"' in content:
                content = content.replace('"tge/internal/core/domain/progression"', 
                    '"tge/internal/core/domain/power"\n\t"tge/internal/core/domain/powersystem"\n\t"tge/internal/core/domain/cultivation"')
            
            # Now replace usage of progression.XXX
            # power.go
            content = content.replace("progression.PowerState", "power.PowerState")
            content = content.replace("progression.Power", "power.Power")
            content = content.replace("progression.NewPower(", "power.NewPower(")
            
            # mechanic_state.go
            content = content.replace("progression.MechanicState", "power.MechanicState")
            content = content.replace("progression.NewMechanicState(", "power.NewMechanicState(")
            
            # power_system.go
            content = content.replace("progression.PowerSystemType", "powersystem.PowerSystemType")
            content = content.replace("progression.PowerSystem", "powersystem.PowerSystem")
            content = content.replace("progression.NewPowerSystem(", "powersystem.NewPowerSystem(")
            content = content.replace("progression.Cultivation", "powersystem.Cultivation")
            content = content.replace("progression.Magic", "powersystem.Magic")
            content = content.replace("progression.SuperPower", "powersystem.SuperPower")
            
            # power_node.go
            content = content.replace("progression.PowerNode", "powersystem.PowerNode")
            content = content.replace("progression.NewPowerNode(", "powersystem.NewPowerNode(")
            
            # cultivation related
            content = content.replace("progression.CultivationState", "cultivation.CultivationState")
            content = content.replace("progression.NewCultivationState(", "cultivation.NewCultivationState(")
            content = content.replace("progression.CultivationPath", "cultivation.CultivationPath")
            content = content.replace("progression.NewCultivationPath(", "cultivation.NewCultivationPath(")
            content = content.replace("progression.Realm", "cultivation.Realm")
            content = content.replace("progression.NewRealm(", "cultivation.NewRealm(")
            content = content.replace("progression.Level", "cultivation.Level")
            content = content.replace("progression.NewLevel(", "cultivation.NewLevel(")
            content = content.replace("progression.ErrLevelNotFound", "cultivation.ErrLevelNotFound")

            if content != orig:
                with open(f, 'w') as file:
                    file.write(content)
