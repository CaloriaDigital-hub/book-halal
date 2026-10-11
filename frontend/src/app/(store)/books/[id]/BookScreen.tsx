"use client";

import { useParams } from "next/navigation";
import { useCallback, useState } from "react";
import { useSession } from "@/modules/auth";
import { BuyBox } from "@/modules/purchase";
import { BookView } from "@/modules/storefront";
import { storeRoutes } from "../../../routes";

export function BookScreen() {
  const { id } = useParams<{ id: string }>();
  const session = useSession();
  const signedIn = session.status === "loading" ? null : session.session.authenticated;
  // Оплату подтвердили, пока страница открыта, — перечитываем книгу (has_access)
  const [version, setVersion] = useState(0);
  const onAccessGranted = useCallback(() => setVersion((v) => v + 1), []);

  return (
    <BookView
      key={version}
      bookId={id}
      backHref={storeRoutes.catalog}
      actions={(book) => (
        <BuyBox
          bookId={book.id}
          price={book.price}
          hasAccess={book.has_access}
          signedIn={signedIn}
          readHref={storeRoutes.read(book.id)}
          buyHref={storeRoutes.buy(book.id)}
          loginHref={storeRoutes.loginThen(storeRoutes.book(book.id))}
          onAccessGranted={onAccessGranted}
        />
      )}
    />
  );
}
