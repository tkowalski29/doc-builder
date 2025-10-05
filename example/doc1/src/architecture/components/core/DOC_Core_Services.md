---
category: architecture/components/core
title: Core Services Catalog
---

# Core Services Catalog

The core service layer exposes essential APIs that power scheduling, reporting,
and campaign optimisation. Each service below includes its owning team and key
SLAs.

## Services

| Service | Team | SLA |
|---------|------|-----|
| Campaign Planner | Growth | 99.9% |
| Feed Transformer | Data Platform | 99.5% |
| Billing Adapter | Finance | 99.9% |

## Operational Notes

- Each service publishes metrics to the shared Prometheus stack.
- Ownership is reflected in the architecture decision records.
