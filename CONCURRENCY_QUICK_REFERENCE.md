# Concurrency Quick Reference

> **Superseded.** See [README.md](README.md) → "Safe patterns (current API)".
> Use `gloo.Run` / `gloo.Pipe` / `gloo.Chain`; build a command once and reuse the
> value across goroutines (stateless and stateful are both safe).
