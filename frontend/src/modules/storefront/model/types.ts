export type CatalogBook = {
  id: string;
  title: string;
  author: string;
  description: string;
  price: number;
  cover_url: string;
};

export type BookDetails = CatalogBook & {
  total_pages: number;
  status: string;
  /** Можно ли читать: бесплатная, куплена или админ. Для гостя — false. */
  has_access: boolean;
};
