// 다크 모드 상태를 localStorage + OS preference와 함께 관리.

import { useEffect, useState, useCallback } from 'react';

const KEY = 'monkey-planner.theme';

export type ThemeMode = 'light' | 'dark';

function readInitialMode(): ThemeMode {
  if (typeof window === 'undefined') return 'light';
  const stored = window.localStorage.getItem(KEY) as ThemeMode | null;
  if (stored === 'light' || stored === 'dark') return stored;
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function applyMode(mode: ThemeMode) {
  const root = document.documentElement;
  root.classList.toggle('dark', mode === 'dark');
  root.dataset.theme = mode;
}

export function useTheme(): { mode: ThemeMode; toggle: () => void; setMode: (m: ThemeMode) => void } {
  const [mode, setModeState] = useState<ThemeMode>(readInitialMode);

  useEffect(() => {
    applyMode(mode);
    window.localStorage.setItem(KEY, mode);
  }, [mode]);

  const toggle = useCallback(() => {
    setModeState((m) => (m === 'dark' ? 'light' : 'dark'));
  }, []);

  return { mode, toggle, setMode: setModeState };
}
