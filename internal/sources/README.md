Every module should contain

source.go
Implementation Source interface, alongside module-dependant logics

sqlite.go
Queries and DTOs behind database access and interaction

handlers.go
Code for handling http requests that are attached to th signalboard server for data serving



Explanation of sources:
commute/
├── domain.go         // business entities
├── dto.go            // HTTP responses
├── mapping.go        // domain <-> dto
├── handlers.go       // HTTP handlers
├── route\_matrix.go   // Google Maps integration
├── source.go         // source orchestration
├── sqlite.go         // persistence
└── store.go          // domain persistence logic
