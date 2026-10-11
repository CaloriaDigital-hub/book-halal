import { Suspense } from "react";
import { RegisterForm } from "@/modules/auth";
import { storeRoutes } from "../../routes";
import styles from "../auth.module.css";

export default function RegisterPage() {
  return (
    <div className={styles.page}>
      <Suspense>
        <RegisterForm loginHref={storeRoutes.login} />
      </Suspense>
    </div>
  );
}
