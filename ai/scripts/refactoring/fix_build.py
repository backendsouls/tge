import re
import os

cli_path = "internal/adapter/cli/cli.go"
with open(cli_path, "r") as f:
    content = f.read()

# Remove awaken, cultivate, train commands
content = re.sub(r'func \(a \*App\) characterAwaken.*?return 0\n}', '', content, flags=re.DOTALL)
content = re.sub(r'func \(a \*App\) characterCultivate.*?return 0\n}', '', content, flags=re.DOTALL)
content = re.sub(r'func \(a \*App\) characterTrain.*?return 0\n}', '', content, flags=re.DOTALL)

# Remove the switch cases
content = re.sub(r'case "cultivate":\s*return a\.characterCultivate\(ctx, args\[1:\]\)', '', content)
content = re.sub(r'case "train":\s*return a\.characterTrain\(ctx, args\[1:\]\)', '', content)
content = re.sub(r'case "awaken":\s*return a\.characterAwaken\(ctx, args\[1:\]\)', '', content)

# Update usage
content = content.replace('<create|list|give-item|cultivate|train|awaken>', '<create|list|give-item>')

# Remove in.Systems = systems
content = content.replace('in.Systems = systems', '')

# Remove systemNames
content = re.sub(r'func systemNames.*?return names\n}', '', content, flags=re.DOTALL)

# Remove PowerSystems usage from list
content = content.replace('NAME\\tTYPE\\tGENDER\\tSPECIES\\tPOWER\\tSYSTEMS', 'NAME\\tTYPE\\tGENDER\\tSPECIES\\tPOWER')
content = re.sub(r'c\.PowerValue, strings\.Join\(systemNames\(c\.PowerSystems\), ", "\)\)', 'c.PowerValue)', content)
content = content.replace('%s\\t%s\\t%s\\t%s\\t%s\\t%s\\n', '%s\\t%s\\t%s\\t%s\\t%s\\n')

# Fix status
status_body_old = r'systems := systemNames\(c\.PowerSystems\).*?return 0'
status_body_new = r'return 0'
content = re.sub(status_body_old, status_body_new, content, flags=re.DOTALL)

with open(cli_path, "w") as f:
    f.write(content)
