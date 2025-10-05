---
category: platform/mobile/features
title: Offline Mode Playbook
---

# Offline Mode Playbook

Offline-first features cache user actions locally and replay them when the
connection is restored. This playbook explains the queueing mechanism and
conflict resolution rules.

## Queue Design

- Every action is serialized into a deterministic JSON envelope.
- Envelopes are persisted in an encrypted SQLite database on device.
- Once connectivity resumes, envelopes are replayed in chronological order.

## Conflict Handling

- Latest write wins for idempotent resources.
- Merge strategies are registered per domain for complex aggregates.
