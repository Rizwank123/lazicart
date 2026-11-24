import { api, unwrap } from "@/api/client";
import type {
  ApiResponse,
  Address,
  AuthTokens,
  Cart,
  Notification,
  Order,
  Payment,
  Product,
  RegisterPayload,
  Review,
  Shipment,
  User,
  UUID,
} from "@/types/api";

export const authApi = {
  register: (payload: RegisterPayload) => unwrap<AuthTokens | User>(api.post<ApiResponse<AuthTokens | User>>("/api/auth/register", payload)),
  login: (email: string, password: string) =>
    unwrap<AuthTokens>(api.post<ApiResponse<AuthTokens>>("/api/auth/login", { email, password })),
  sendOtp: (phone: string) => unwrap<null>(api.post<ApiResponse<null>>("/api/auth/send-otp", undefined, { params: { phone } })),
  loginWithPhone: (phone: string, otp: string) =>
    unwrap<AuthTokens>(api.post<ApiResponse<AuthTokens>>("/api/auth/login-with-phone", { phone, otp })),
  refresh: (refreshToken: string) =>
    unwrap<AuthTokens>(api.post<ApiResponse<AuthTokens>>("/api/auth/refresh-token", undefined, { params: { refresh_token: refreshToken } })),
  logout: (jti: string) => unwrap<null>(api.get<ApiResponse<null>>("/api/auth/logout", { params: { jti } })),
};

export const userApi = {
  byId: (id: UUID) => unwrap<User>(api.get<ApiResponse<User>>(`/api/users/${id}`)),
  byEmail: (email: string) => unwrap<User>(api.get<ApiResponse<User>>("/api/users/email", { params: { email } })),
  update: (id: UUID, payload: Partial<User>) => unwrap<User>(api.put<ApiResponse<User>>(`/api/users/${id}`, payload)),
};

export const productApi = {
  list: (page = 1, limit = 20) =>
    api.get<ApiResponse<Product[]>>(
      "/api/products",
      { params: { page, limit } },
    ),
  detail: (id: UUID) => unwrap<Product>(api.get<ApiResponse<Product>>(`/api/products/${id}`)),
  reviews: (id: UUID) => unwrap<Review[]>(api.get<ApiResponse<Review[]>>(`/api/products/${id}/reviews`)),
  addReview: (payload: { product_id: UUID; user_id: UUID; rating: number; comment?: string; is_approved: boolean }) =>
    unwrap<Product>(api.post<ApiResponse<Product>>("/api/products/reviews", payload)),
};

export const inventoryApi = {
  byProduct: (productId: UUID) => unwrap<unknown>(api.get<ApiResponse<unknown>>(`/api/inventory/products/${productId}`)),
};

export const cartApi = {
  get: (userId: UUID) => unwrap<Cart>(api.get<ApiResponse<Cart>>(`/api/carts/${userId}`)),
  addItem: (payload: { user_id: UUID; product_id: UUID; quantity: number; unit_price: number; currency: string }) =>
    unwrap<Cart>(api.post<ApiResponse<Cart>>("/api/carts/items", payload)),
  updateItem: (itemId: UUID, quantity: number) => unwrap<null>(api.put<ApiResponse<null>>(`/api/carts/items/${itemId}`, { quantity })),
  removeItem: (itemId: UUID) => unwrap<null>(api.delete<ApiResponse<null>>(`/api/carts/items/${itemId}`)),
  clear: (userId: UUID) => unwrap<null>(api.delete<ApiResponse<null>>(`/api/carts/${userId}`)),
};

export const orderApi = {
  create: (payload: {
    user_id: UUID;
    currency: string;
    shipping_address: Address;
    items: { product_id: UUID; quantity: number; unit_price: number }[];
  }) => unwrap<Order>(api.post<ApiResponse<Order>>("/api/orders", payload)),
  byId: (id: UUID) => unwrap<Order>(api.get<ApiResponse<Order>>(`/api/orders/${id}`)),
  byUser: (userId: UUID, page = 1, limit = 20) =>
    unwrap<Order[]>(api.get<ApiResponse<Order[]>>(`/api/orders/users/${userId}`, { params: { page, limit } })),
  cancel: (id: UUID) => unwrap<null>(api.post<ApiResponse<null>>(`/api/orders/${id}/cancel`)),
};

export const paymentApi = {
  create: (payload: { order_id: UUID; user_id: UUID; amount: number; currency: string; provider: string }) =>
    unwrap<Payment>(api.post<ApiResponse<Payment>>("/api/payments", payload)),
  byOrder: (orderId: UUID) => unwrap<Payment>(api.get<ApiResponse<Payment>>(`/api/payments/orders/${orderId}`)),
};

export const notificationApi = {
  byUser: (userId: UUID, page = 1, limit = 20) =>
    unwrap<Notification[]>(api.get<ApiResponse<Notification[]>>(`/api/notifications/users/${userId}`, { params: { page, limit } })),
};

export const deliveryApi = {
  byOrder: (orderId: UUID) => unwrap<Shipment>(api.get<ApiResponse<Shipment>>(`/api/deliveries/orders/${orderId}`)),
  byUser: (userId: UUID, page = 1, limit = 20) =>
    unwrap<Shipment[]>(api.get<ApiResponse<Shipment[]>>(`/api/deliveries/users/${userId}`, { params: { page, limit } })),
};
