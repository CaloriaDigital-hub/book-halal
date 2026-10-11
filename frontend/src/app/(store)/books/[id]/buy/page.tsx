import { Suspense } from "react";
import { Spinner } from "@/shared/ui";
import { BuyScreen } from "./BuyScreen";

export default function BuyPage() {
  return (
    <Suspense fallback={<Spinner label="Загрузка…" />}>
      <BuyScreen />
    </Suspense>
  );
}
