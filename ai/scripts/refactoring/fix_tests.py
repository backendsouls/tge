import os
import re

# 1. cli_test.go
path = 'internal/adapter/cli/cli_test.go'
with open(path, 'r') as f: content = f.read()
if '"strings"' not in content:
    content = content.replace('"testing"', '"strings"\n\t"testing"')

content = content.replace(
    'func (f *fakeSpeciesService) GetSpecies(context.Context, string) (character.Species, error) {\n\treturn character.Species{}, nil\n}',
    'func (f *fakeSpeciesService) Species(ctx context.Context, string) (character.Species, error) {\n\treturn character.Species{}, nil\n}'
)
if 'func (f *fakeUniverseService) AddRealms(' not in content:
    content += '\nfunc (f *fakeUniverseService) AddRealms(ctx context.Context, in port.AddRealmsInput) (cosmology.Universe, error) { return cosmology.Universe{}, nil }\n'
with open(path, 'w') as f: f.write(content)


# 2. character_service_test.go
path = 'internal/core/service/character_service_test.go'
with open(path, 'r') as f: content = f.read()
content = re.sub(
    r'NewCharacterService\(([^,]+),\s*([^)]+)\)',
    r'NewCharacterService(\1, \2, &stubSpeciesRepo{})',
    content
)
if 'type stubSpeciesRepo struct{}' not in content:
    content += '''
type stubSpeciesRepo struct{}
func (s *stubSpeciesRepo) Save(ctx context.Context, sp character.Species) error { return nil }
func (s *stubSpeciesRepo) FindByName(ctx context.Context, name string) (character.Species, error) { return character.Species{Name: name}, nil }
func (s *stubSpeciesRepo) List(ctx context.Context) ([]character.Species, error) { return nil, nil }
'''
with open(path, 'w') as f: f.write(content)

# 3. universe_service_test.go
path = 'internal/core/service/universe_service_test.go'
with open(path, 'r') as f: content = f.read()
if '"tge/internal/core/domain/progression"' not in content:
    content = content.replace('"tge/internal/core/domain/cosmology"', '"tge/internal/core/domain/cosmology"\n\t"tge/internal/core/domain/progression"')
with open(path, 'w') as f: f.write(content)
