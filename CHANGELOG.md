# Changelog

## [1.1.0](https://github.com/backendsouls/tge/compare/v1.0.0...v1.1.0) (2026-08-16)


### Features

* docs sync with code ([#5](https://github.com/backendsouls/tge/issues/5)) ([cd9e52b](https://github.com/backendsouls/tge/commit/cd9e52bfc81f49fcd527fe90f359ce34144e0433))

## 1.0.0 (2026-08-09)


### Features

* automate releases with goreleaser and add install script ([49f9a56](https://github.com/backendsouls/tge/commit/49f9a5640db1a9b42a6e24a2fe41317ab9443161))
* **character:** enforce unique character names via automatic iterative suffixing (N) ([0906ed4](https://github.com/backendsouls/tge/commit/0906ed4cbbb235032be6142c1f8aad12e400f133))
* **character:** implement CharacterRepository soft-deletion mechanism (.deleted files) ([1db43b4](https://github.com/backendsouls/tge/commit/1db43b44bd693f31d9cc57a97bd945a608e28fda))
* **character:** species and mortal characters ([6ba7411](https://github.com/backendsouls/tge/commit/6ba74119f942f39365b08d815f45acf77e8b6e57))
* **cli:** command-line adapter ([9499d19](https://github.com/backendsouls/tge/commit/9499d19dc347a37e7cf8aad38a22541534bd70d6))
* **cli:** dynamically generate auto-shorthand flags for all CLI commands based on available letters ([5376120](https://github.com/backendsouls/tge/commit/537612023ed3597fbba06e785f78a9056437201c))
* **cli:** introduce "tge character clean" subcommand to purge characters ([91618d9](https://github.com/backendsouls/tge/commit/91618d91a2a3576981ab40264da4855fcd62a372))
* **cli:** introduce distinct "--- Novel Log" and "--- Internal System Log" formatting boundaries ([40e3ce9](https://github.com/backendsouls/tge/commit/40e3ce93050a712d9750618c79f2b115905e7aa0))
* **cli:** strictly enforce POSIX double-dash standards for long flags (rejecting single-dash) ([a845aa1](https://github.com/backendsouls/tge/commit/a845aa18583bfdb3c47b34f634656e5ea78096b1))
* **cmd:** composition root, seeding and functional tests ([d95d6bc](https://github.com/backendsouls/tge/commit/d95d6bcb75593d228c9f7dbf86c6ca479297c31d))
* **config:** embedded defaults and starter catalog ([91607d8](https://github.com/backendsouls/tge/commit/91607d8c488adb814d0cf21358152f2da72b32c1))
* **cosmology:** per-location timelines ([9aad081](https://github.com/backendsouls/tge/commit/9aad0814005db1c9b79a82ad83d2ffdbe6853ab2))
* **cosmology:** reality-to-universe containment hierarchy ([16b0843](https://github.com/backendsouls/tge/commit/16b0843e0ac7f458ff5bcd89775d6dbf4a074fda))
* Implement dual logging, RPG grade system, and AI scripts ([9015f0d](https://github.com/backendsouls/tge/commit/9015f0dccdc7b083694849a49edcb4361d4654f6))
* **novel:** novels with volumes and chapters ([6c2172d](https://github.com/backendsouls/tge/commit/6c2172daac576ec54e475631d7037a8bcd1a903f))
* **port:** ports for cosmology, rpg and novels ([67cd018](https://github.com/backendsouls/tge/commit/67cd018321fc325ac7844d0ed8c7a03e5f1842f9))
* **port:** ports for progression and characters ([c02373b](https://github.com/backendsouls/tge/commit/c02373ba5f0bc9b1eb4a636fc2650f02b95a8ad1))
* **progression:** cultivation state and system kinds ([af97351](https://github.com/backendsouls/tge/commit/af973516f15bbc87c8489a2112c44fdce037ca26))
* **progression:** power trees and cultivation paths ([52f6c53](https://github.com/backendsouls/tge/commit/52f6c538a1265839e817bcfbfd9432c8b11f0fb4))
* **progression:** realms and their ordered levels ([603871e](https://github.com/backendsouls/tge/commit/603871e48fb6d729ba404441c7ccd580941a7034))
* **rpg:** equipment, classes, professions, quests and recipes ([bf0fc8d](https://github.com/backendsouls/tge/commit/bf0fc8d2df29b982342c8e400feb06ba6f168315))
* **rpg:** stats, inventory and core items ([48dde1d](https://github.com/backendsouls/tge/commit/48dde1de9cbb3aba4fd3d666b676eba4fa11daf4))
* **service:** cosmology, rpg and novel use cases ([b7a677b](https://github.com/backendsouls/tge/commit/b7a677bdf5b249a3f37e91d6384fbed374613f16))
* **service:** realm, power-system and character use cases ([d62dec7](https://github.com/backendsouls/tge/commit/d62dec77ed5b0957c465e53ef1bc9252dc005978))
* **sqlite:** character, rpg and novel repositories ([fbb6b68](https://github.com/backendsouls/tge/commit/fbb6b687d3ffa253230cb5faa859f0a3604e74e6))
* **sqlite:** progression and cosmology repositories ([7e0bad7](https://github.com/backendsouls/tge/commit/7e0bad77b3ddb38285abc69b6c1b96cd4feed29c))
* **sqlite:** schema and goose migrations ([3b9e865](https://github.com/backendsouls/tge/commit/3b9e86555d90f1dd3fa4cb71a133ac53c8de0d04))


### Bug Fixes

* **cli:** prevent empty-argument execution on status command to enforce usage menu ([2d8ec6c](https://github.com/backendsouls/tge/commit/2d8ec6c23783256a0e0c273c9e6e1a50f354857e))
