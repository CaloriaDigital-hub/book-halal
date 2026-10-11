import { Suspense } from "react";
import { LoginForm } from "@/modules/auth";
import styles from "./page.module.css";

export default function LoginPage() {
  return (
    <main className={styles.page}>
      <Suspense>
        <LoginForm />
      </Suspense>
    </main>
  );
}
