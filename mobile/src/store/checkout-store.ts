import { create } from "zustand";

import type { Address } from "@/types/api";

type CheckoutState = {
  address: Address;
  setAddress: (address: Address) => void;
};

export const useCheckoutStore = create<CheckoutState>((set) => ({
  address: { city: "", state: "", country: "India", pincode: "" },
  setAddress: (address) => set({ address }),
}));
