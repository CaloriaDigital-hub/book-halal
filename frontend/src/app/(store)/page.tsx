"use client";

import { CatalogView } from "@/modules/storefront";
import { storeRoutes } from "../routes";

export default function CatalogPage() {
  return <CatalogView bookHref={storeRoutes.book} />;
}
