"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { Button } from "@/shared/ui";
import { signOut } from "../api/session-api";
import { setSession } from "../model/use-session";

export function SignOutButton() {
  const router = useRouter();
  const [pending, setPending] = useState(false);

  async function onClick() {
    setPending(true);
    try {
      await signOut();
    } finally {
      setSession({ authenticated: false });
      router.replace("/login");
    }
  }

  return (
    <Button variant="ghost" size="sm" loading={pending} onClick={onClick}>
      Выйти
    </Button>
  );
}
