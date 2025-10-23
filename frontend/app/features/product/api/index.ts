import { apiClient } from "@/lib/apiClient";
import { GetProductsResponse } from "@/packages/types/src";

export async function getProducts(limit: number, offset: number) {
  const res = await apiClient.get<GetProductsResponse>("products", {
    params: {
      limit: limit,
      offset: offset,
    },
  });
  return res.data;
}
