import { Suspense } from "react";
import { LoginForm } from "@/modules/auth";
import { storeRoutes } from "../../routes";
import styles from "../auth.module.css";

export default function LoginPage() {
  return (
    <div className={styles.page}>
      <Suspense>
        <LoginForm registerHref={storeRoutes.register} />
      </Suspense>
    </div>
  );
}
