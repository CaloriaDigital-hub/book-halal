"use client";

import { useSyncExternalStore } from "react";
import { getSession, type Session } from "../api/session-api";

// Один запрос /api/session на вкладку; компоненты подписываются на результат.
type State = { status: "loading" } | { status: "ready"; session: Session };

const LOADING: State = { status: "loading" };

let state: State = LOADING;
let started = false;
const listeners = new Set<() => void>();

function emit(next: State) {
  state = next;
  listeners.forEach((l) => l());
}

function load() {
  started = true;
  getSession()
    .then((session) => emit({ status: "ready", session }))
    .catch(() => emit({ status: "ready", session: { authenticated: false } }));
}

function subscribe(listener: () => void) {
  listeners.add(listener);
  if (!started) load();
  return () => {
    listeners.delete(listener);
  };
}

export function useSession(): State {
  return useSyncExternalStore(subscribe, () => state, () => LOADING);
}

/** Обновить состояние после входа/выхода без лишнего запроса. */
export function setSession(session: Session) {
  started = true;
  emit({ status: "ready", session });
}
