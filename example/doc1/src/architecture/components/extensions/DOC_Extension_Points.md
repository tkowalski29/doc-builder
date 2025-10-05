---
category: architecture/components/extensions
title: Extension Points & Webhooks
---

# Extension Points & Webhooks

Extension endpoints allow partners to react to lifecycle events within the
platform. This document catalogues each webhook, expected payloads, and retry
behaviour.

## Events

- `order.created`
- `order.cancelled`
- `campaign.paused`
- `campaign.reactivated`

## Reliability

All webhooks are delivered with exponential backoff and are signed with an HMAC
hash derived from the client secret.
