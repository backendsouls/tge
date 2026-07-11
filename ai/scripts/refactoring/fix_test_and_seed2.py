import os

def replace_in_file(path, old, new):
    if not os.path.exists(path): return
    with open(path, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(path, 'w') as f:
            f.write(content)

replace_in_file('cmd/tge/seed.go', 'powersystem.ErrLevelNumberExists', 'cultivation.ErrLevelNumberExists')
replace_in_file('test/functional/progression_ft_test.go', 'progression.', 'cultivation.')

# Fix missing mock method List in test/functional
def fix_mock_species(path):
    if not os.path.exists(path): return
    with open(path, 'r') as f:
        content = f.read()
    if 'func (m mockSpecies) Save(' in content and 'func (m mockSpecies) List(' not in content:
        content = content.replace('func (m mockSpecies) FindByName(ctx context.Context, name string) (character.Species, error) { return character.Species{Name: name}, nil }',
                                  'func (m mockSpecies) FindByName(ctx context.Context, name string) (character.Species, error) { return character.Species{Name: name}, nil }\nfunc (m mockSpecies) List(ctx context.Context) ([]character.Species, error) { return nil, nil }')
        with open(path, 'w') as f:
            f.write(content)

fix_mock_species('test/functional/character_ft_test.go')
fix_mock_species('test/functional/mechanic_state_ft_test.go')
fix_mock_species('test/functional/progression_ft_test.go')

# character.Item -> rpg.Item
replace_in_file('test/functional/character_ft_test.go', 'character.Item', 'rpg.Item')
replace_in_file('test/functional/progression_ft_test.go', 'port.DefaultSpecies', 'service.CharacterDefaults{}')

