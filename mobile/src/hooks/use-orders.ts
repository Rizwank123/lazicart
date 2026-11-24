import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { orderApi, paymentApi } from "@/api/endpoints";
import { env } from "@/config/env";
import { useAuthStore } from "@/store/auth-store";
import type { Address, CartItem } from "@/types/api";

export function useOrders() {
  const userId = useAuthStore((state) => state.userId);
  return useQuery({
    queryKey: ["orders", userId],
    queryFn: () => orderApi.byUser(userId as string),
    enabled: !!userId,
  });
}

export function useCreateCheckout() {
  const queryClient = useQueryClient();
  const userId = useAuthStore((state) => state.userId);
  return useMutation({
    mutationFn: async ({ address, items }: { address: Address; items: CartItem[] }) => {
      const currency = items[0]?.currency ?? env.defaultCurrency;
      const order = await orderApi.create({
        user_id: userId as string,
        currency,
        shipping_address: address,
        items: items.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
          unit_price: item.unit_price,
        })),
      });
      const payment = await paymentApi.create({
        order_id: order.id,
        user_id: userId as string,
        amount: order.total,
        currency: order.currency,
        provider: "manual",
      });
      return { order, payment };
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["orders", userId] });
      queryClient.invalidateQueries({ queryKey: ["cart", userId] });
    },
  });
}
