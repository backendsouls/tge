# Adapters Diagram

This diagram visualizes the external adapter layers implementing the core application interfaces. 
- The **CLI Adapter** translates terminal user inputs into domain structures and triggers Services.
- The **SQLite Adapter** implements the Core Driven Ports (Repositories) to persist the domain state to disk.

```mermaid
classDiagram
    %% Driving Adapter
    class CLIApp {
        -RPGService rpg
        -CharacterService char
        +Run(args List~String~) int
        +runItem(ctx Context, args List~String~) int
    }
    
    %% Driven Adapter
    class SQLiteItemRepository {
        -sql.DB db
        +Save(Context, Item) error
        +FindByName(Context, String) Item
        +List(Context) List~Item~
    }
    
    %% Core Abstractions
    class ItemService {
        <<service>>
    }
    
    class ItemRepository {
        <<interface>>
    }

    CLIApp --> ItemService : invokes commands
    SQLiteItemRepository ..|> ItemRepository : implements
```
