import { create } from "zustand";

import type { Address } from "@/types/api";

type CheckoutState = {
  address: Address;
  setAddress: (address: Address) => void;
  /**
   * Idempotency key of the checkout that is in progress.
   *
   * The key lives in the store rather than in the screen so that a retry
   * after a partial failure reuses it. The backend then returns the order
   * that was already created instead of creating a second one.
   */
  idempotencyKey: string | null;
  beginCheckout: () => string;
  finishCheckout: () => void;
  /** Drop the stored address and key, for example when the account changes. */
  reset: () => void;
};

const EMPTY_ADDRESS: Address = { city: "", state: "", country: "India", pincode: "" };

function newIdempotencyKey() {
  const random = Math.random().toString(36).slice(2, 12);
  return `checkout-${Date.now().toString(36)}-${random}`;
}

export const useCheckoutStore = create<CheckoutState>((set, get) => ({
  address: { ...EMPTY_ADDRESS },
  idempotencyKey: null,
  setAddress: (address) => set({ address }),
  beginCheckout: () => {
    const existing = get().idempotencyKey;
    if (existing) return existing;
    const key = newIdempotencyKey();
    set({ idempotencyKey: key });
    return key;
  },
  finishCheckout: () => set({ idempotencyKey: null }),
  reset: () => set({ address: { ...EMPTY_ADDRESS }, idempotencyKey: null }),
}));
