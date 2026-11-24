import { useInfiniteQuery, useQuery } from "@tanstack/react-query";

import { productApi } from "@/api/endpoints";
import type { Product } from "@/types/api";

export function useProducts() {
  return useInfiniteQuery({
    queryKey: ["products"],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => {
      const response = await productApi.list(pageParam, 12);
      return {
        items: response.data.data ?? [],
        page: response.data.pagination?.page ?? pageParam,
        total: response.data.pagination?.total ?? 0,
        limit: response.data.pagination?.limit ?? 12,
      };
    },
    getNextPageParam: (last) => {
      const loaded = last.page * last.limit;
      return loaded < last.total ? last.page + 1 : undefined;
    },
  });
}

export function useProduct(id?: string) {
  return useQuery({
    queryKey: ["product", id],
    queryFn: () => productApi.detail(id as string),
    enabled: !!id,
  });
}

export function flattenProducts(pages?: { items: Product[] }[]) {
  return pages?.flatMap((page) => page.items) ?? [];
}
