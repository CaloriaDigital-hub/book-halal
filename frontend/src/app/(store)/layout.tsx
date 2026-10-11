import { AccountMenu } from "@/modules/auth";
import { StoreShell } from "@/shared/ui";
import { storeRoutes } from "../routes";

export default function StoreLayout({ children }: LayoutProps<"/">) {
  return (
    <StoreShell
      brand="Book Halal"
      homeHref={storeRoutes.catalog}
      actions={<AccountMenu loginHref={storeRoutes.login} registerHref={storeRoutes.register} adminHref={storeRoutes.admin} />}
    >
      {children}
    </StoreShell>
  );
}
