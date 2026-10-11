"use client";

import { useCallback, useEffect, useEffectEvent, useState } from "react";

type State<T> = { key: string | null; data?: T; error?: Error };

/**
 * Загрузка данных виджетом. `key` — всё, от чего зависит запрос
 * (смена ключа = новый запрос). Старые данные остаются на экране, пока идёт новый.
 */
export function useAsync<T>(key: string, load: () => Promise<T>) {
  const [nonce, setNonce] = useState(0);
  const [state, setState] = useState<State<T>>({ key: null });
  const requestKey = `${key}#${nonce}`;
  const runLoad = useEffectEvent(load);

  useEffect(() => {
    let active = true;
    runLoad().then(
      (data) => {
        if (active) setState({ key: requestKey, data });
      },
      (error: unknown) => {
        if (active) setState((prev) => ({ key: requestKey, data: prev.data, error: toError(error) }));
      },
    );
    return () => {
      active = false;
    };
  }, [requestKey]);

  const reload = useCallback(() => setNonce((n) => n + 1), []);

  return {
    data: state.data,
    error: state.key === requestKey ? state.error : undefined,
    loading: state.key !== requestKey,
    reload,
  };
}

export function toError(error: unknown): Error {
  return error instanceof Error ? error : new Error(String(error));
}
