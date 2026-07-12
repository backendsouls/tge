# Core Ports & Services Diagram

This diagram visualizes the Hexagonal Architecture pattern used to separate driving logic (Services) from driven logic (Repositories) inside the core domain. Note that this pattern is replicated identically across all domain models via Generics, but mapped out here using the RPG `Item` as a reference.

```mermaid
classDiagram
    class ItemService {
        -ItemRepository repo
        +CreateItem(CreateItemInput) Item
        +GetItem(String) Item
        +ListItems() List~Item~
    }
    
    class ItemRepository {
        <<interface>>
        +Save(Context, Item) error
        +FindByName(Context, String) Item
        +List(Context) List~Item~
    }
    
    class CreateItemInput {
        <<struct>>
        +String Name
        +String Description
        +String Grade
    }

    ItemService ..> ItemRepository : drives
    ItemService ..> CreateItemInput : consumes
```
