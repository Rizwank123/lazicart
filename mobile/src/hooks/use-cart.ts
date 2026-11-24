import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { cartApi } from "@/api/endpoints";
import { useAuthStore } from "@/store/auth-store";

export function useCart() {
  const userId = useAuthStore((state) => state.userId);
  return useQuery({
    queryKey: ["cart", userId],
    queryFn: () => cartApi.get(userId as string),
    enabled: !!userId,
  });
}

export function useCartActions() {
  const queryClient = useQueryClient();
  const userId = useAuthStore((state) => state.userId);
  const invalidate = () => queryClient.invalidateQueries({ queryKey: ["cart", userId] });

  const add = useMutation({
    mutationFn: (payload: { product_id: string; quantity: number; unit_price: number; currency: string }) =>
      cartApi.addItem({ user_id: userId as string, ...payload }),
    onSuccess: invalidate,
  });
  const update = useMutation({
    mutationFn: ({ itemId, quantity }: { itemId: string; quantity: number }) => cartApi.updateItem(itemId, quantity),
    onSuccess: invalidate,
  });
  const remove = useMutation({
    mutationFn: (itemId: string) => cartApi.removeItem(itemId),
    onSuccess: invalidate,
  });
  const clear = useMutation({
    mutationFn: () => cartApi.clear(userId as string),
    onSuccess: invalidate,
  });

  return { add, update, remove, clear };
}
