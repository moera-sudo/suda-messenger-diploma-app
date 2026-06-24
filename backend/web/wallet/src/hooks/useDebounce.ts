import { useCallback, useEffect, useRef } from 'react';

// useDebounce returns a stable debounced wrapper around fn. The timer lives in a
// ref (not a render-scoped closure) so re-renders never lose the pending call,
// and it is cleared on unmount to avoid firing after the component is gone.
export function useDebounce<Args extends unknown[]>(
  fn: (...args: Args) => void,
  ms: number,
): (...args: Args) => void {
  const fnRef = useRef(fn);
  fnRef.current = fn;

  const timer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  useEffect(() => () => clearTimeout(timer.current), []);

  return useCallback(
    (...args: Args) => {
      clearTimeout(timer.current);
      timer.current = setTimeout(() => fnRef.current(...args), ms);
    },
    [ms],
  );
}
