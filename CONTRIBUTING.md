# Contributing to TGE

Welcome! We are excited to have you contribute to The Great Emulator (TGE). 

## Branching Strategy
We follow a standard Git branching model:
- `main` is our stable branch.
- Feature work should be done on isolated branches (e.g., `feature/add-magic`, `fix/cli-crash`).
- Once your work is ready, open a Pull Request against `main`.

## Continuous Integration (CI)
When you open a Pull Request, GitHub Actions will automatically run two pipelines:
1. `[Push]`: Verifies your isolated branch code.
2. `[Merge]`: Verifies a simulated merge with `main`.
Both must pass (Lint, Vet, Build, Unit Tests, Functional Tests) before merging.

## Conventional Commits & Versioning
TGE uses **Google's Release Please** to entirely automate semantic versioning and changelog generation. **You must strictly follow the Conventional Commits specification** for your PR titles and commit messages.

Release Please analyzes your commit prefixes to determine the next version bump:

- `fix: ...` : Triggers a **PATCH** version bump (e.g., `1.0.0` -> `1.0.1`). Use for bug fixes.
- `feat: ...` : Triggers a **MINOR** version bump (e.g., `1.0.0` -> `1.1.0`). Use for new features.
- `feat!: ...` or `fix!: ...` : Triggers a **MAJOR** version bump (e.g., `1.0.0` -> `2.0.0`). Use for breaking changes.

Other accepted prefixes that do **not** trigger a release bump (but appear in history):
- `chore:`, `docs:`, `style:`, `refactor:`, `perf:`, `test:`

## Releasing
You do not manually create releases or tags!
When code is merged to `main`, Release Please will automatically open a "Release PR". It will continually update this PR as more commits land in `main`.

**To publish a release to the world, simply merge the Release PR!** 
GitHub Actions will automatically generate the Git tag and compile/upload the binaries via GoReleaser.
