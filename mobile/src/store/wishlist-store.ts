import AsyncStorage from "@react-native-async-storage/async-storage";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

type WishlistState = {
  ids: string[];
  toggle: (id: string) => void;
  has: (id: string) => boolean;
};

export const useWishlistStore = create<WishlistState>()(
  persist(
    (set, get) => ({
      ids: [],
      toggle: (id) =>
        set((state) => ({
          ids: state.ids.includes(id) ? state.ids.filter((item) => item !== id) : [id, ...state.ids],
        })),
      has: (id) => get().ids.includes(id),
    }),
    { name: "lahzicart.wishlist", storage: createJSONStorage(() => AsyncStorage) },
  ),
);
