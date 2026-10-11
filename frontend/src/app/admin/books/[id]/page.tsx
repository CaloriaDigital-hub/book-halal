import { Suspense } from "react";
import { Spinner } from "@/shared/ui";
import { BookScreen } from "./BookScreen";

export default function BookPage() {
  return (
    <Suspense fallback={<Spinner label="Загрузка…" />}>
      <BookScreen />
    </Suspense>
  );
}
