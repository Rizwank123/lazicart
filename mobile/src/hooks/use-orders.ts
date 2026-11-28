import { useInfiniteQuery, useMutation, useQueryClient } from "@tanstack/react-query";

import { orderApi, paymentApi } from "@/api/endpoints";
import { env } from "@/config/env";
import { useAuthStore } from "@/store/auth-store";
import { useCheckoutStore } from "@/store/checkout-store";
import type { Address, CartItem } from "@/types/api";

const ORDERS_PAGE_SIZE = 20;

export function useOrders() {
  const userId = useAuthStore((state) => state.userId);
  return useInfiniteQuery({
    queryKey: ["orders", userId],
    initialPageParam: 1,
    queryFn: ({ pageParam }) => orderApi.byUser(userId as string, pageParam, ORDERS_PAGE_SIZE),
    // The order service returns a plain array, so a short page means the end
    // of the list was reached.
    getNextPageParam: (last, all) => (last.length < ORDERS_PAGE_SIZE ? undefined : all.length + 1),
    enabled: !!userId,
  });
}

export function useCancelOrder(orderId?: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await orderApi.cancel(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["order", orderId] });
      queryClient.invalidateQueries({ queryKey: ["orders"] });
    },
  });
}

export function useCreateCheckout() {
  const queryClient = useQueryClient();
  const userId = useAuthStore((state) => state.userId);
  return useMutation({
    mutationFn: async ({ address, items }: { address: Address; items: CartItem[] }) => {
      if (!userId) throw new Error("Your session is not ready yet. Please try again.");
      // The key is kept in the store, so a retry of the same checkout reuses
      // it and the backend returns the order that already exists.
      const key = useCheckoutStore.getState().beginCheckout();
      const currency = items[0]?.currency ?? env.defaultCurrency;
      const order = await orderApi.create(
        {
          user_id: userId,
          currency,
          shipping_address: address,
          items: items.map((item) => ({
            product_id: item.product_id,
            quantity: item.quantity,
            unit_price: item.unit_price,
          })),
        },
        key,
      );
      // The order exists now, so the list has to show it even when the
      // payment fails afterwards.
      queryClient.invalidateQueries({ queryKey: ["orders", userId] });
      const payment = await paymentApi.create(
        {
          order_id: order.id,
          user_id: userId,
          amount: order.total,
          currency: order.currency,
          provider: "manual",
        },
        `${key}-payment`,
      );
      useCheckoutStore.getState().finishCheckout();
      return { order, payment };
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders", userId] });
      queryClient.invalidateQueries({ queryKey: ["cart", userId] });
      queryClient.invalidateQueries({ queryKey: ["notifications", userId] });
    },
  });
}
