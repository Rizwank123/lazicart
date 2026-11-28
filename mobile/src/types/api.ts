export type UUID = string;

export type ApiResponse<T> = {
  message: string;
  data?: T;
  pagination?: Pagination;
};

export type Pagination = {
  total?: number;
  page?: number;
  limit?: number;
};

export type Address = {
  city?: string;
  state?: string;
  country?: string;
  pincode?: string;
  coordinates?: number[];
  [key: string]: unknown;
};

export type BaseEntity = {
  id: UUID;
  created_at?: string;
  updated_at?: string;
  deleted_at?: string | null;
};

export type AuthTokens = {
  access_token: string;
  refresh_token: string;
  jti: string;
};

/**
 * The auth-service account record returned by registration.
 *
 * It is not the profile of the user-service: the profile is created
 * asynchronously from the user.created event and is fetched separately.
 */
export type AuthUser = {
  id: UUID;
  email?: string;
  phone?: string;
  is_active?: boolean;
  created_at?: string;
  updated_at?: string;
};

export type RegisterPayload = {
  email: string;
  phone: string;
  password: string;
  first_name: string;
  last_name: string;
  address: Address;
  dob: string;
  gender: string;
};

export type User = BaseEntity & {
  first_name?: string;
  last_name?: string;
  email?: string;
  auth_id: UUID;
  phone?: string;
  is_active?: boolean;
  dob?: string;
  gender?: string;
  address?: Address;
};

/**
 * Fields the user-service accepts on update.
 *
 * `Partial<User>` is not correct here: the backend ignores id, auth_id and
 * the audit timestamps, and it replaces the whole address when the field is
 * present.
 */
export type UpdateUserPayload = {
  first_name?: string;
  last_name?: string;
  email?: string;
  phone?: string;
  dob?: string;
  gender?: string;
  address?: Address;
};

export type InventoryItem = BaseEntity & {
  product_id: UUID;
  sku?: string | null;
  quantity_available: number;
  quantity_reserved: number;
  reorder_level: number;
  is_active: boolean;
};

export type ProductImage = BaseEntity & {
  product_id: UUID;
  image_url: string;
  thumbnail_url?: string | null;
  is_primary: boolean;
  display_order: number;
};

export type Review = BaseEntity & {
  product_id: UUID;
  user_id: UUID;
  rating: number;
  comment?: string | null;
  is_approved: boolean;
};

export type Product = BaseEntity & {
  name: string;
  description?: string | null;
  price: number;
  currency: string;
  category_id?: UUID | null;
  brand?: string | null;
  attributes: Record<string, unknown>;
  is_active: boolean;
  images?: ProductImage[];
  review?: Review[];
};

export type CartItem = BaseEntity & {
  cart_id: UUID;
  product_id: UUID;
  quantity: number;
  unit_price: number;
  currency: string;
};

export type Cart = BaseEntity & {
  user_id: UUID;
  status: string;
  items?: CartItem[];
};

export type OrderItem = BaseEntity & {
  order_id: UUID;
  product_id: UUID;
  quantity: number;
  unit_price: number;
  currency: string;
};

export type Order = BaseEntity & {
  user_id: UUID;
  status: string;
  total: number;
  currency: string;
  shipping_address: Address;
  items?: OrderItem[];
};

export type Payment = BaseEntity & {
  order_id: UUID;
  user_id: UUID;
  amount: number;
  currency: string;
  provider: string;
  status: string;
  transaction_id?: string | null;
  failure_reason?: string | null;
};

export type Notification = BaseEntity & {
  user_id?: UUID | null;
  channel: string;
  type: string;
  recipient?: string | null;
  subject?: string | null;
  body: string;
  status: string;
  metadata: Record<string, unknown>;
};

export type Shipment = BaseEntity & {
  order_id: UUID;
  user_id: UUID;
  carrier?: string | null;
  tracking_number?: string | null;
  status: string;
  shipping_address: Address;
  estimated_delivery_at?: string | null;
  delivered_at?: string | null;
};
