# Legacy runtime

This directory contains the original AutoSEL parser, event monitor, policy
loader, and container replacement implementation while the migration is in
progress.

It is not part of the supported runtime contract. Nothing under this
directory may be called by the read-only `autosel plan` path. In particular,
the legacy code can stop, remove, recreate, or relabel containers and has not
been qualified for data preservation or rollback.

New code must use `internal/spec`, `internal/policy`, and explicit adapters
instead. Remove legacy packages only after their required behavior has been
replaced and covered by tests.
