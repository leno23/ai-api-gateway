"use client";

import { useEffect, type DependencyList } from "react";

/** Run callback after mount without synchronous setState in the effect body (eslint react-hooks/set-state-in-effect). */
export function useDeferredEffect(
  effect: () => void | Promise<void>,
  deps: DependencyList,
): void {
  useEffect(() => {
    queueMicrotask(() => {
      void effect();
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps -- caller controls deps
  }, deps);
}
