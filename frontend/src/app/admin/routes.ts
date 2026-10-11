// Адреса экранов админки — единственное место, где модули узнают про роутинг.
export const adminRoutes = {
  overview: "/admin",
  books: "/admin/books",
  book: (id: string) => `/admin/books/${encodeURIComponent(id)}`,
};
