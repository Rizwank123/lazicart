# LahziCart Implementation Plan

Status: proposed roadmap
Last reviewed: 2026-11-25

## Objective

Evolve LahziCart from a working ecommerce microservice foundation into a reliable, searchable, AI-assisted commerce platform with:

- complete catalog and category management;
- product search, autocomplete, faceted filters, and ranking;
- AI query understanding, recommendations, comparison, and shopping assistance;
- correct order, inventory, payment, and delivery consistency;
- reliable event delivery, idempotent consumers, observability, and recovery;
- customer features expected from a Flipkart/Amazon-style shopping experience.

## Current Baseline

Implemented:

- Auth: registration, login, OTP login, refresh, logout, Redis token state, JWT, and gRPC validation.
- User: profile lifecycle and user lifecycle event consumers.
- Catalog: products, images, reviews, ratings, pagination, and basic filtering.
- Inventory: stock, adjustments, reservations, release, and commit.
- Cart: carts and cart items.
- Order: order/item persistence, status, cancellation, and order events.
- Payment: payment records, status updates, and payment events.
- Notification: persisted notifications and order/payment consumers.
- Delivery: shipment persistence and status APIs.
- Gateway: reverse proxy and JWT protection.
- Shared: Kafka producer/consumer, retry/DLQ behavior, events, middleware, responses, metrics, and auth proto.
- Mobile: Expo client with auth, catalog, cart, checkout, orders, notifications, profile, and delivery flows.

Important gaps:

- Order creation trusts client-provided prices and is not a complete inventory/payment workflow.
- Kafka publishing happens directly after database writes; there is no transactional outbox.
- Consumers do not yet have a shared inbox/idempotency table.
- There is no Saga coordinator connecting order, inventory, payment, and delivery.
- Catalog has no category CRUD, variants, normalized attributes, or dedicated search service.
- Payment has no provider adapter, webhook verification, authorization, capture, void, or refund flow.
- Notification providers, device-token registration, user preferences, and delivery retries are incomplete.
- Wishlist is local-only in the mobile client.
- Automated testing is mostly compile coverage.

## Architecture Decisions

### Search and filtering

Create a dedicated search-service instead of placing search-specific behavior in catalog-service.

Initial implementation:

1. PostgreSQL full-text search and trigram matching.
2. pgvector for product and query embeddings.
3. Denormalized search documents built from catalog, inventory, rating, and promotion events.
4. Hybrid retrieval: lexical match plus vector similarity.
5. Exact filters and facet aggregation.
6. Ranking from text score, vector score, availability, price, rating, popularity, and business rules.
7. Move to OpenSearch only when traffic, facet complexity, or indexing latency requires it.

Bloom and Cuckoo filters are not product-filter engines. Product price, stock, category, and permission checks must use exact data stores or indexes.

### Bloom versus Cuckoo

Recommendation: use Cuckoo only as an optional dynamic membership accelerator.

Cuckoo is appropriate when membership changes frequently and deletion is needed. Bloom is appropriate for append-mostly or periodically rebuilt sets where deletion is not required. Both can return false positives, so a positive result must fall through to an authoritative lookup.

For LahziCart, a Cuckoo filter may protect hot product or availability-cache lookups. It must never authorize a purchase or decide exact stock. Inventory remains an atomic PostgreSQL update with a conditional quantity check.

### AI

Implement AI in increasing-risk order:

1. Query normalization, typo correction, synonyms, and attribute extraction.
2. Semantic product search using embeddings and hybrid ranking.
3. Autosuggest from popular, recent, trending, and session queries.
4. Similar products and frequently-bought-together recommendations.
5. Product comparison and a grounded shopping assistant.
6. Personalized ranking after enough click, cart, purchase, and conversion data exists.

AI must retrieve from approved catalog data and must not invent price, stock, delivery promises, review claims, or payment status. Every AI result needs an exact product and availability lookup before purchase.

### Distributed transactions

Use an orchestrated Saga with local ACID transactions, transactional outboxes, inbox/idempotency handling, and compensating actions. Do not use two-phase commit across service databases and Kafka.

The order service should initially own checkout Saga coordination. It can later become a separate checkout-orchestrator service without changing participant contracts.

Saga:

    CREATED
      -> INVENTORY_RESERVING
      -> INVENTORY_RESERVED
      -> PAYMENT_AUTHORIZING
      -> PAYMENT_AUTHORIZED
      -> CONFIRMED
      -> FULFILLING
      -> SHIPPED
      -> DELIVERED

Compensation:

    inventory failure  -> cancel order
    payment failure    -> release inventory and cancel order
    payment timeout    -> retry, then expire and release inventory
    cancellation       -> release inventory and void/refund payment
    delivery failure   -> retry shipment or create an operations task

Every command and event must contain event_id, event_type, version, aggregate_id, saga_id, correlation_id, occurred_at, producer, and payload.

## Implementation Phases

### Phase 0: Foundation and correctness

Priority: P0

- [~] Add a shared event envelope and versioning. Implemented in `shared/events` and used by order events.
- [~] Add outbox_events tables and relay workers to event-producing services. Order service is the first vertical slice; remaining event producers still need rollout.
- [~] Add inbox_messages or equivalent idempotency tables to consumers. Implemented transactionally for the auth `user.created` -> user-service consumer path; remaining consumers still need adoption.
- [~] Add checkout Idempotency-Key storage. Order service now persists an optional key and request hash atomically with the order; same-request replay returns the original order and conflicting reuse returns 409.
- [~] Make Kafka brokers configurable in every consumer. User and notification consumers now use `KAFKA_BROKERS`; remaining consumer hosts still need migration.
- [~] Add readiness checks and graceful consumer shutdown. User and notification now expose database-backed `/ready` checks and graceful Kafka shutdown; remaining hosts and richer dependency checks still need work.
- [~] Add request, correlation, and Saga ID propagation. Gateway and shared HTTP middleware establish/propagate request and correlation IDs, preserve optional Saga IDs, and order-created envelopes carry both IDs. Active Kafka consumers now extract envelope IDs into handler context; legacy non-enveloped events still need migration.
- [~] Add OpenTelemetry tracing. Shared tracer initialization and Echo middleware are wired into API Gateway and all HTTP services; Kafka producers now inject W3C trace context and consumers create linked spans. Collector-backed integration verification remains.

Acceptance: database writes and event intent commit atomically; retries do not duplicate business effects; correlated logs follow checkout across services.

Current implementation slice: order creation and cancellation now persist their event intent in the same PostgreSQL transaction as the order mutation, and a retrying order outbox relay publishes the versioned envelope to Kafka. Consumer inbox/idempotency and rollout to the remaining producers are still pending.

### Phase 1: Catalog and search

Priority: P0

- [ ] Add category tree and category APIs.
- [ ] Add product variants/SKUs and normalized attributes.
- [ ] Add product change events.
- [ ] Create search-service.
- [ ] Implement lexical search, autocomplete, facets, filters, sort, and ranking.
- [ ] Add pgvector embeddings and hybrid search.
- [ ] Add reindexing and search relevance metrics.

Initial API:

    GET /api/search?q=&category_id=&brand=&min_price=&max_price=&rating=&in_stock=&sort=&page=&limit=
    GET /api/search/suggest?q=&session_id=
    POST /api/search/reindex

Acceptance: exact terms, typo tolerance, category, brand, price, rating, stock filters, zero-result metrics, and zero-downtime reindexing work.

### Phase 2: Inventory and checkout Saga

Priority: P0

- [ ] Replace client prices with catalog price snapshots.
- [ ] Add order Saga state and transition tables.
- [ ] Add inventory reservation lease and expiry.
- [ ] Make reserve, release, and commit commands idempotent.
- [ ] Add payment authorization, capture, void, and refund interfaces.
- [ ] Add inventory/payment compensation.
- [ ] Create shipment only after the correct order state.
- [ ] Add last-unit concurrency tests.

Checkout:

    checkout intent -> validate price -> reserve inventory -> authorize payment
      -> confirm order -> capture payment -> create shipment

Acceptance: no overselling, same result for repeated idempotency keys, inventory release after payment failure, and restart-safe Saga recovery.

### Phase 3: Payments and post-order lifecycle

Priority: P1

- [ ] Add provider adapter and sandbox provider.
- [ ] Add signed, replay-protected payment webhooks.
- [ ] Add payment attempt history.
- [ ] Add refunds and partial refunds.
- [ ] Add cancellation windows and returns.
- [ ] Add return shipment and inspection states.
- [ ] Add customer order timeline.

### Phase 4: AI and recommendations

Priority: P1

- [ ] Add versioned product/review embedding pipeline.
- [ ] Add query intent and attribute extraction.
- [ ] Add semantic and hybrid search.
- [ ] Add similar-product and frequently-bought-together recommendations.
- [ ] Add personalized ranking with cold-start fallback.
- [ ] Add grounded comparison and shopping assistant.

AI API candidates:

    GET /api/recommendations/home
    GET /api/recommendations/products/:id/similar
    POST /api/assistant/product-help
    POST /api/compare

### Phase 5: Customer and marketplace features

Priority: P1/P2

- [ ] Server-side wishlist.
- [ ] Coupons, promotions, taxes, and shipping fees.
- [ ] Saved addresses.
- [ ] Push device-token registration and notification preferences.
- [ ] Product Q&A, review moderation, and abuse reporting.
- [ ] Seller/merchant capabilities if marketplace behavior is required.
- [ ] Search, recommendation, analytics, and admin APIs.

### Phase 6: Production operations

Priority: P1

- [ ] Distributed tracing and Kafka lag metrics.
- [ ] Dashboards and alerts for outbox backlog, DLQ growth, reservation expiry, payment failures, and zero-result search spikes.
- [ ] Secrets management; remove development credentials from committed environment files.
- [ ] CI for formatting, tests, migration checks, security scanning, and container builds.
- [ ] Deployment manifests and environment-specific configuration.
- [ ] Load tests for search, reservation, and checkout.

## Testing Strategy

1. Unit tests for domain rules, state transitions, ranking, and compensation.
2. PostgreSQL repository integration tests.
3. Kafka event schema and compatibility tests.
4. Duplicate, out-of-order, and retry consumer tests.
5. Last-unit inventory concurrency tests.
6. Gateway end-to-end tests.
7. Search relevance evaluation with curated queries.
8. Failure injection between database commit and event relay.

## Recommended Order

1. Outbox, inbox, event envelope, idempotency, readiness, and tracing.
2. Correct checkout with exact price snapshots, inventory reservation, payment state, and Saga.
3. Catalog categories, variants, normalized attributes, and product events.
4. Search service with lexical search, facets, autocomplete, and ranking.
5. Payment provider, webhooks, refunds, returns, and delivery lifecycle.
6. Hybrid AI search, recommendations, comparison, and assistant.
7. Wishlist, promotions, saved addresses, Q&A, analytics, and admin.
8. Load testing, security hardening, deployment, and operations.

## Research References

- Amazon Search: The joy of ranking products: https://www.amazon.science/publications/amazon-search-the-joy-of-ranking-products
- Amazon ecommerce search and recommendations research: https://www.amazon.science/publications/challenges-and-research-opportunities-in-ecommerce-search-and-recommendations
- Flipkart personalized autosuggest: https://blog.flipkart.tech/building-personalized-autosuggestion-9e705d5bf5f8
- Flipkart distributed inventory tracing: https://blog.flipkart.tech/distributed-entity-tracing-9e293940b086
- Flipkart Commerce Cloud search: https://docs.flipkartcommercecloud.com/docs/digital-commerce/SEO/docs/overview
- Redis Cuckoo filter documentation: https://redis.io/docs/latest/data-types/probabilistic/cuckoo-filter/
- CMU Cuckoo Filter paper: https://www.cs.cmu.edu/~binfan/papers/conext14_cuckoofilter.pdf
- Microsoft Saga pattern: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- AWS transactional outbox: https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html

## Definition of Done

A phase is complete only when its code, migrations, API contracts, events, tests, gateway registration, configuration, observability, and documentation are updated together. Happy-path endpoints alone are not production readiness; retry, duplicate, timeout, crash, and recovery behavior must also be verified.
