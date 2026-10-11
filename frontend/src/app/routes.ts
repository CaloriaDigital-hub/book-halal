// Адреса экранов магазина — модули получают их пропсами и не знают про роутинг.
export const storeRoutes = {
  catalog: "/",
  login: "/login",
  register: "/register",
  admin: "/admin",
  book: (id: string) => `/books/${encodeURIComponent(id)}`,
  buy: (id: string) => `/books/${encodeURIComponent(id)}/buy`,
  read: (id: string) => `/books/${encodeURIComponent(id)}/read`,
  loginThen: (next: string) => `/login?next=${encodeURIComponent(next)}`,
};
