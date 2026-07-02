// Package integration holds cross-cutting integration tests for the gloo
// framework. They exercise behavior that spans the framework's primitives
// rather than a single command — in particular its concurrency guarantees:
//
//   - Stateless commands (patterns.Map/Filter) are safe to reuse: the same
//     Command value can be executed from many goroutines concurrently.
//   - Stateful commands (patterns.StatefulMap/StatefulFilter) get fresh state
//     per Execute, so a reused value never leaks state across goroutines.
//
// The tests carry the "integration" build tag (they are cross-cutting, not
// 1:1 unit tests). Run the race detector to validate those guarantees:
//
//	go test -race -tags integration ./tests.integration/
package integration
