# LahziCart Task Status

## Added: `sample.env`
- Comprehensive environment configuration with all 9 service database names
- Common DB settings, Kafka, Redis, and gateway URLs
- Each service can copy to its `.env` and set `DB_NAME_<SERVICE>`

## Created: `scripts/create-databases.sh`
- Checks and creates all 9 service databases
- Proper existence checking with `psql -lqt`
- Handles hyphenated database names (e.g., `user-service`)
- Summary reporting: created / already existed / failed

## Databases Created/Verified
All 9 databases confirmed in PostgreSQL:
- `user_auth` (Auth Service)
- `user-service` (User Service)
- `catalog` (Catalog Service)
- `inventory` (Inventory Service)
- `cart` (Cart Service)
- `orders` (Order Service)
- `payment` (Payment Service)
- `notification` (Notification Service)
- `delivery` (Delivery Service)

## Migrations Applied
- `make migrate SERVICE=<name> ACTION=up` ran successfully for all 9 services
- 2 migration files fixed for Goose annotation format

## Services Running (earlier)
All 10 services were running and healthy:
- API Gateway (7070), Auth (7700), User (7701), Catalog (7702), 
- Inventory (7703), Cart (7704), Order (7705), Payment (7706), 
- Notification (7707), Delivery (7708)

## Added: frontend and API integration fixes

### Backend
- Gateway: every proxied service URL is configured, empty targets answer 503, and proxy failures answer with a JSON 502.
- Gateway: the caller identity is taken from the verified JWT, any client supplied `X-User-Id` is removed, and a request that addresses another user is rejected with 403.
- Auth: refresh tokens are validated with the refresh key and checked against Redis, so a rotated or revoked token is rejected.
- Auth: phone login resolves the auth record and signs the auth-service id; the password hash is no longer serialized.
- Gateway: a revoked session is effective right away. The auth-service keeps one server side session entry per token pair, and the gateway reads that state, so a token that was signed out or already replaced by a refresh is rejected with a 401 instead of staying usable until it expires. The key format is shared by the writer and the reader through `shared/session`, and the gateway falls back to the signature check alone when Redis cannot be reached.
- Catalog: product payloads include images and approved reviews; the moderation flag is decided by the server.
- User: the profile id is the auth-service id, partial updates no longer wipe the address, and the delete endpoint really deletes.
- Payment: creation is idempotent through the `Idempotency-Key` header.
- Kafka: the producer creates a missing topic and retries a failed write, and `make kafka-topics` bootstraps the full list.

### Mobile
- A single shared refresh operation, so concurrent 401 responses can no longer erase a refreshed session.
- Profile resolution runs in the background with a bounded deadline and a session generation, so a slow service cannot block sign-in and a stale response cannot restore a previous account.
- The checkout key is kept in the checkout store, payment creation is idempotent, and a cart cleanup failure no longer turns a completed order into a failed checkout.
- Mixed currency carts are grouped by currency and blocked at checkout instead of adding incompatible prices.
- Protected screens redirect to sign-in, errors are shown instead of empty states, and the cart, orders and notifications screens paginate and refresh.
- Registration collects a real date of birth, gender and state instead of sending placeholders.

### Web build
- The app loads in the browser. Zustand v5 ships an ESM build that reads `import.meta.env.MODE`, and the web bundle is served as a classic script, so `import.meta` was a syntax error that stopped the whole bundle from parsing and left only the splash screen. `metro.config.js` points the web build at the CommonJS build, which has no `import.meta`. The native build already used that build through the `react-native` export condition.
- The browser can call the API. The gateway does not answer the CORS preflight, so every request with a body was blocked before it reached the gateway. `scripts/dev-api-proxy.js` answers the preflight and adds the CORS headers, and `npm run web` starts it together with the web server. A device still talks to the gateway directly.
