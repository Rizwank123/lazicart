# LahziCart Implementation Memory

Purpose: maintain the implementation record for Plane.md. Update this file whenever a planned item is started, completed, deferred, or redesigned.

Last updated: 2026-11-25

## Status Legend

- [x] implemented and verified
- [~] partially implemented or needs hardening
- [ ] planned but not implemented
- [!] blocked or requires a decision

## Current Baseline

### Services

- [x] Auth: registration, login, OTP, refresh, logout, Redis state, JWT, and gRPC validation.
- [x] User: profile lifecycle and user lifecycle event consumers.
- [~] Catalog: products, images, reviews, ratings, pagination, and basic filtering. Categories, variants, normalized attributes, and real search remain.
- [~] Inventory: stock, adjustments, reservations, release, and commit. Expiry, idempotency, and Saga integration remain.
- [~] Cart: carts and cart items. Catalog price validation and checkout integration remain.
- [~] Order: order/item persistence, status, cancellation, and events. Client prices and distributed coordination remain.
- [~] Payment: payment records and status events. Provider integration, webhooks, authorization, capture, void, and refunds remain.
- [~] Notification: notification persistence and order/payment consumers. Providers, preferences, device tokens, and idempotency remain.
- [~] Delivery: shipment persistence and status APIs. Carrier integration and event-driven creation remain.
- [~] Gateway: routes and JWT protection. Full environment configuration and authorization hardening remain.
- [~] Mobile: Expo flows for auth, catalog, cart, checkout, orders, notifications, profile, and delivery. Wishlist is local-only.

### Shared Platform

- [x] Go workspace includes root, gateway, services, and shared module.
- [x] Shared Kafka producer/consumer with retry, exponential backoff, and DLQ naming.
- [x] Shared event topic constants and ecommerce event payload types.
- [x] Shared Echo middleware for request IDs, logging, recovery, timeout, CORS, validation, rate limiting, and metrics.
- [x] Goose migration scripts per service.
- [~] Shared observability exists, but tracing, Kafka lag, outbox backlog, and Saga metrics remain.
- [~] scripts/run-all-services.sh starts the services. Dependency readiness and per-service log files remain.

## Plan Progress

### Phase 0: Foundation

- [~] Event envelope and versioning. Implemented in `shared/events/envelope.go` and used by order outbox events.
- [~] Transactional outbox tables and relay workers. Implemented for order-service first; other event producers remain.
- [~] Inbox/idempotency tables. Implemented transactionally for the user-service `user.created` consumer; remaining consumers still need adoption.
- [~] Checkout idempotency keys. Order service persists an optional `Idempotency-Key` and request hash atomically with the order and outbox event; same-request replay returns the original order, conflicting reuse returns 409.
- [~] Configurable Kafka broker injection everywhere. User and notification consumers now use loaded `KAFKA_BROKERS`; remaining consumer hosts still need migration.
- [~] Readiness checks and graceful consumer shutdown. User and notification expose database-backed `/ready` checks and graceful Kafka shutdown; remaining hosts and richer dependency checks remain.
- [~] Request, correlation, and Saga IDs. Gateway and shared HTTP middleware establish/propagate request and correlation IDs, preserve optional Saga IDs, and order-created envelopes carry both IDs. Active Kafka consumers extract envelope IDs into handler context; legacy non-enveloped events still need migration.
- [~] OpenTelemetry tracing. Shared tracer initialization and Echo middleware are wired into API Gateway and all HTTP services; Kafka producers inject W3C trace context and consumers create linked spans. Collector-backed integration verification remains.

### Phase 1: Catalog and search

- [ ] Category tree and APIs.
- [ ] Product variants/SKUs and normalized attributes.
- [ ] Product change events.
- [ ] Dedicated search service.
- [ ] Lexical search, autocomplete, facets, filters, and ranking.
- [ ] pgvector embeddings and hybrid search.
- [ ] Reindexing and relevance metrics.

### Phase 2: Checkout Saga

- [ ] Exact catalog price snapshot.
- [ ] Order Saga state machine.
- [ ] Reservation lease and expiry.
- [ ] Idempotent inventory commands.
- [ ] Payment authorization/capture/void/refund interface.
- [ ] Inventory and payment compensation.
- [ ] Shipment creation after correct order state.
- [ ] Last-unit concurrency tests.

### Phase 3: Payments and post-order

- [ ] Provider adapter.
- [ ] Signed, replay-protected webhooks.
- [ ] Payment attempts.
- [ ] Refunds and partial refunds.
- [ ] Returns and return shipments.
- [ ] Customer order timeline.

### Phase 4: AI

- [ ] Versioned product/review embeddings.
- [ ] Query intent and attribute extraction.
- [ ] Semantic/hybrid search.
- [ ] Similar-product recommendations.
- [ ] Frequently-bought-together recommendations.
- [ ] Personalized ranking with cold-start fallback.
- [ ] Grounded comparison and shopping assistant.

### Phase 5: Customer features

- [ ] Server-side wishlist.
- [ ] Coupons and promotions.
- [ ] Saved addresses.
- [ ] Push tokens and notification preferences.
- [ ] Product Q&A and review moderation.
- [ ] Seller/merchant capabilities if required.
- [ ] Search, recommendation, analytics, and admin APIs.

### Phase 6: Operations

- [ ] Distributed tracing.
- [ ] Kafka lag and DLQ dashboards.
- [ ] Outbox/Saga alerts.
- [ ] Secrets management.
- [ ] CI and security scanning.
- [ ] Container/deployment manifests.
- [ ] Load and failure-injection tests.

## Decisions

### Search

Build search as a separate service. Start with PostgreSQL full-text/trigram search plus pgvector, then move to OpenSearch if scale or facet requirements justify it.

### Probabilistic filter

Use Cuckoo only for dynamic lookup acceleration where deletion is needed. Do not use Bloom or Cuckoo for exact inventory, pricing, authorization, or order validation.

### Distributed transaction

Use an orchestrated Saga with transactional outbox and inbox/idempotency. Do not use 2PC across independent service databases and Kafka. Order service is the initial coordinator; a future checkout-orchestrator may be extracted later.

Compensations:

- Inventory reserve failure: cancel order.
- Payment failure: release inventory and cancel order.
- Payment timeout: retry, then expire and release inventory.
- Cancellation after authorization: void or refund payment and release inventory.
- Delivery failure: retry shipment or create an operations task.

## Change Log

### 2026-11-25

- Audited backend and mobile repository structure.
- Confirmed implemented services: gateway, auth, user, catalog, inventory, cart, order, payment, notification, and delivery.
- Confirmed mobile flows for auth, catalog, cart, checkout, orders, notifications, profile, and delivery.
- Identified missing production features: outbox/inbox, Saga orchestration, idempotency, search service, categories, variants, AI retrieval, recommendations, provider payments, refunds, wishlist backend, promotions, addresses, push tokens, tracing, and operational tests.
- Researched Amazon and Flipkart search/autosuggest patterns, Cuckoo versus Bloom behavior, and Saga/outbox guidance.
- Selected Cuckoo only as an optional dynamic lookup accelerator.
- Selected hybrid lexical plus vector search for product discovery.
- Selected orchestrated Saga plus transactional outbox/inbox for distributed consistency.
- Implemented the first Phase 0 vertical slice: versioned event envelope, order outbox migration, atomic order/event persistence, retrying order outbox relay, and order event publication through the relay.
- Remaining for this slice: database/Kafka integration verification, inbox/idempotency, readiness and graceful shutdown hardening, and rollout to the other event-producing services.
- Implementation files: `shared/events/envelope.go`, `Services/order-service/internal/db/migration/20261125110729_add_order_outbox.sql`, `Services/order-service/internal/repository/order_repository.go`, `Services/order-service/internal/service/order_service.go`, `Services/order-service/internal/outbox/relay.go`, and `Services/order-service/cmd/server/main.go`.
- Verification: `go test ./...` passes in `Services/order-service`; shared packages also pass. Full multi-module workspace output is limited by the current Go workspace layout and still needs a dedicated CI command.
- Implemented the next Phase 0 slice for `user.created`: auth now publishes the shared event envelope, user-service has a Goose-generated `inbox_messages` migration, and event claim, user creation, and processed marking commit in one transaction.
- Verification: `go test ./...` passes in `Services/user-service`, `Services/auth-service`, and `shared`. Remaining work is to adopt the inbox boundary across all future Kafka consumers and add database/Kafka integration tests.
- Implemented checkout idempotency storage with Goose migration `Services/order-service/internal/db/migration/20261125113605_add_order_idempotency.sql`, request hashing, unique-key protection, and replay/conflict handling in the order repository and HTTP handler.
- Verification: `go test ./...` passes in `Services/order-service`; database concurrency tests are still pending.
- Implemented configurable Kafka brokers and graceful consumer shutdown for the shared Kafka consumer, user-service, and notification-service. Context cancellation now exits cleanly, retry backoff observes cancellation, and service shutdown waits for consumer goroutines.
- Verification: `go test ./...` passes in `shared`, `Services/user-service`, and `Services/notification-service`. Remaining work is readiness checks and adoption by any future consumer hosts.
- Added `/health` liveness and bounded database-backed `/ready` endpoints to user-service and notification-service. Readiness returns `503 NOT_READY` when the database dependency cannot be reached.
- Verification: `go test ./...` passes in user-service and notification-service after the readiness changes.
- Extended shared request-ID middleware to generate/propagate `X-Correlation-ID` and optional `X-Saga-ID`, and wired order checkout metadata into the `order.created` event envelope.
- Verification: `go test ./...` passes in `Services/order-service` and `shared`; gateway header propagation and Kafka consumer context propagation remain.
- Added API Gateway propagation middleware. It establishes valid `X-Request-ID` and `X-Correlation-ID` headers, preserves valid `X-Saga-ID`, and forwards all three through reverse-proxy requests.
- Verification: `go test ./...` passes in `Services/api-gateway`; Kafka consumer context propagation remains.
- Added `events.ContextFromMessage` and `Consumer.ConsumeWithContext` so Kafka handlers receive correlation and Saga IDs from versioned envelopes. User and notification consumers now use the context-aware API.
- Verification: `go test ./...` passes in `shared`, `Services/user-service`, and `Services/notification-service`. Legacy event producers that do not emit envelopes remain to be migrated.
- Added shared OpenTelemetry tracer initialization and HTTP tracing middleware. API Gateway and order-service now create spans; `OTEL_EXPORTER_OTLP_ENDPOINT` enables OTLP HTTP export and `OTEL_EXPORTER_OTLP_INSECURE=true` supports local plaintext collectors.
- Verification: `go test ./...` passes in `shared`, `Services/api-gateway`, and `Services/order-service` after tracing changes. Remaining work is service-wide rollout and Kafka producer/consumer span propagation.
- Rolled the shared tracing provider and HTTP middleware into auth, user, catalog, inventory, cart, payment, notification, and delivery services.
- Verification: all eight service modules passed `go test ./...`; the earlier user-service initialization-order issue was corrected. Kafka producer/consumer span links remain.
- Added W3C trace-context injection to the shared Kafka producer and extraction plus consumer spans to the shared Kafka consumer. Kafka retry handler execution now runs under the extracted trace context.
- Verification: `go test ./...` passes in `shared` and `Services/order-service`; end-to-end collector verification remains.

## Update Rules

1. Move a plan item from [ ] to [~] when implementation starts.
2. Record affected files, migrations, API/event changes, and tests.
3. Move it to [x] only after code, migration, gateway, configuration, and verification are complete.
4. Record blockers under the relevant phase.
5. Add a dated entry to the Change Log.
