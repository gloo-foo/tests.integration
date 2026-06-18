# Concurrency Summary

> **Superseded.** See [README.md](README.md) for the current framework's
> concurrency model. Short version: stateless commands (`patterns.Map`/`Filter`)
> are safe to reuse across goroutines; stateful commands
> (`patterns.StatefulMap`/`StatefulFilter`) get fresh state per `Execute`.
