import os

# fix mockItem
with open('test/functional/character_ft_test.go', 'r') as f:
    content = f.read()
if 'func (mockItem) List' not in content:
    content = content.replace('func (mockItem) FindByName(ctx context.Context, name string) (rpg.Item, error) {\n\treturn rpg.Item{Name: name}, nil\n}', 'func (mockItem) FindByName(ctx context.Context, name string) (rpg.Item, error) {\n\treturn rpg.Item{Name: name}, nil\n}\nfunc (mockItem) List(ctx context.Context) ([]rpg.Item, error) { return nil, nil }\nfunc (mockItem) Save(ctx context.Context, i rpg.Item) error { return nil }')
    with open('test/functional/character_ft_test.go', 'w') as f:
        f.write(content)

# fix progression_ft_test.go
with open('test/functional/progression_ft_test.go', 'r') as f:
    content = f.read()
content = content.replace('Species: service.CharacterDefaults{}', 'Species: character.Species{Name: "Human"}')
with open('test/functional/progression_ft_test.go', 'w') as f:
    f.write(content)

