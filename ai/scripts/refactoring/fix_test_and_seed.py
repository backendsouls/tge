import os

def replace_in_file(path, old, new):
    if not os.path.exists(path): return
    with open(path, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(path, 'w') as f:
            f.write(content)

replace_in_file('cmd/tge/seed.go', 'progression.', 'powersystem.')
replace_in_file('test/functional/cycle_detection_ft_test.go', 'progression.', 'powersystem.')

# Fix AddEdge return value in cycle_detection_ft_test.go
with open('test/functional/cycle_detection_ft_test.go', 'r') as f:
    content = f.read()
content = content.replace('err := sysSvc.AddEdge(', '_, err := sysSvc.AddEdge(')
content = content.replace('err = sysSvc.AddEdge(', '_, err = sysSvc.AddEdge(')
with open('test/functional/cycle_detection_ft_test.go', 'w') as f:
    f.write(content)

