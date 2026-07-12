# Cosmology Domain Class Diagram

This diagram models the hierarchy of existence inside `internal/core/domain/cosmology`. The system scales from a foundational Reality down to specific alternative Timelines.

```mermaid
classDiagram
    class Reality {
        +String Name
    }
    
    class Omniverse {
        +String Name
        +String Reality
    }
    
    class Multiverse {
        +String Name
        +String Omniverse
    }
    
    class Universe {
        +String Name
        +String Multiverse
    }
    
    class Timeline {
        +String ID
        +String Name
        +String Universe
    }

    Reality "1" *-- "many" Omniverse : contains
    Omniverse "1" *-- "many" Multiverse : contains
    Multiverse "1" *-- "many" Universe : contains
    Universe "1" *-- "many" Timeline : contains
```
