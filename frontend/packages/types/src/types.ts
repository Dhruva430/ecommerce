export interface Product {
  id: string;
  title: string;
  category: string;
  details: string[];
  image: string;
  price: {
    payable: string;
    savings: string;
  };
}
