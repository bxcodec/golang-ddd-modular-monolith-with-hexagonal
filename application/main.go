// Package main is the entry point for the modular monolith application.
//
// This application demonstrates a modular monolith architecture with hexagonal design:
//
// Modular Monolith:
//   - Multiple independent modules (payment, payment-settings) in a single deployable unit
//   - Each module has clear boundaries and communicates via well-defined interfaces
//   - Modules can be developed, tested, and evolved independently
//   - All modules share the same process and database, simplifying deployment and transactions
//
// Hexagonal Architecture (Ports & Adapters):
//   - Domain logic (hexagon core) is isolated from external concerns
//   - Inbound adapters: REST controllers, cron jobs (drive the application)
//   - Outbound adapters: Repositories, external APIs (driven by the application)
//   - Ports: Interfaces that define contracts between core and adapters
//
// Benefits:
//   - Simpler than microservices while maintaining modularity
//   - Easy refactoring to microservices if needed (modules are already isolated)
//   - Shared infrastructure reduces operational complexity
//   - Domain-driven design principles throughout
package main

import (
	"github.com/bxcodec/golang-ddd-modular-monolith-with-hexagonal/cmd"
)

func main() {
	cmd.Execute()
}
