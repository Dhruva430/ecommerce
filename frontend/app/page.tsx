import { getProducts } from "./features/product/api";

export default async function Home() {
  const products = await getProducts(10, 20);
  return <div></div>;
}
