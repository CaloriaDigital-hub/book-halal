import { Suspense } from "react";
import { ReadScreen } from "./ReadScreen";

export default function ReadPage() {
  return (
    <Suspense>
      <ReadScreen />
    </Suspense>
  );
}
