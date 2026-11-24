# LahziCart Backend Analysis

Analyzed services: API gateway, auth, user, catalog, inventory, cart, order, payment, notification, delivery, shared response/events/middleware.

## Gateway

Base URL defaults to `http://localhost:7070`.

Public:
- `GET /health`
- `POST /api/auth/*`
- `GET /api/products`
- `GET /api/products/*`

JWT protected:
- `/api/users/*`
- mutating `/api/products/*`
- `/api/inventory/*`
- `/api/carts/*`
- `/api/orders/*`
- `/api/payments/*`
- `/api/notifications/*`
- `/api/deliveries/*`

Gateway proxies `/api/products` to catalog `/products`, `/api/carts` to cart `/carts`, and similarly for other services.

## API Endpoints

Auth:
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/send-otp?phone=...`
- `POST /api/auth/login-with-phone`
- `POST /api/auth/refresh-token?refresh_token=...`
- `GET /api/auth/logout?jti=...`

Users:
- `GET /api/users/:id`
- `GET /api/users/email?email=...`
- `GET /api/users/phone?phone=...`
- `PUT /api/users/:id`
- `DELETE /api/users/:id`

Products:
- `GET /api/products?page=&limit=`
- `GET /api/products/:id`
- protected create/update/delete product
- reviews: `POST /api/products/reviews`, `GET /api/products/:id/reviews`, update/delete review
- images: `POST /api/products/images`, `GET /api/products/:id/images`, update/delete image

Cart:
- `GET /api/carts/:user_id`
- `POST /api/carts/items`
- `PUT /api/carts/items/:item_id`
- `DELETE /api/carts/items/:item_id`
- `DELETE /api/carts/:user_id`

Orders:
- `POST /api/orders`
- `GET /api/orders/:id`
- `GET /api/orders/users/:user_id?page=&limit=`
- `POST /api/orders/:id/cancel`

Payments:
- `POST /api/payments`
- `GET /api/payments/:id`
- `GET /api/payments/orders/:order_id`
- `PUT /api/payments/:id/status`

Inventory:
- `GET /api/inventory`
- `GET /api/inventory/products/:product_id`
- protected upsert/adjust/reserve/release/commit

Notifications:
- `POST /api/notifications`
- `GET /api/notifications/users/:user_id`
- `POST /api/notifications/:id/sent`

Deliveries:
- `POST /api/deliveries`
- `GET /api/deliveries/:id`
- `GET /api/deliveries/orders/:order_id`
- `GET /api/deliveries/users/:user_id`
- `PUT /api/deliveries/:id/status`

## Business Rules

- Auth issues access/refresh JWTs with `sub`, `role`, `jti`, `exp`.
- Registration creates auth record and publishes `user.created`; user-service consumes this to create profile.
- OTP is stored in Redis for five minutes and logged server-side.
- Products require `name`, positive `price`; default currency is `INR`.
- Reviews require rating `1..5`.
- Cart requires user/product ids and positive quantity; duplicate cart product increments quantity.
- Orders require user id and items; total is computed from client-supplied unit prices.
- Payment is manual/pending by default; no external provider redirect or SDK is implemented.
- Delivery is a shipment record with pending status.
- Notifications are persisted records; no mobile push token registration endpoint exists.

## Missing Or Incomplete APIs

- No category CRUD/list endpoint, despite `category_id` on products.
- No wishlist service.
- No coupons/promo codes.
- No password reset/change endpoint.
- No social login.
- No saved-address CRUD; user profile has a single embedded address.
- No recommendation/trending/banner/promotion endpoints.
- Product filter service method exists but no route exposes it.
- No true payment gateway integration or checkout session endpoint.
- No push notification device-token registration endpoint.
- Gateway JWT middleware validates only signature/expiry and does not enforce ownership or roles.

## Mobile Architecture

- Expo SDK 55, Expo Router, TypeScript strict mode.
- Central Axios API client with token injection and refresh-token retry.
- SecureStore for tokens, AsyncStorage-backed local wishlist.
- TanStack Query for server data caching/pagination.
- Zustand for auth, checkout draft, and local wishlist state.
- React Hook Form + Zod on auth forms.
- Dark-mode-aware theme primitives under `src/components/ui`.

## Screen Flow

- Auth: login, register, phone OTP.
- Tabs: home/search, cart, orders, notifications, profile.
- Stack: product details, checkout, order details.
- Checkout creates order, creates pending manual payment, clears cart locally through the cart API.

## Extension Strategy

Missing backend capabilities are isolated behind frontend-only UX or typed API modules so future endpoints can be added without reshaping screens: categories can replace inferred brand chips, wishlist can swap from local Zustand to server mutations, and payment can replace manual pending creation with a provider checkout session.
