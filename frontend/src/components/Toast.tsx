import { createContext, useCallback, useContext, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { CheckCircle2, XCircle, X } from 'lucide-react';

type ToastType = 'success' | 'error';
type ToastItem = { id: number; type: ToastType; message: string };

const Ctx = createContext<{
  toast: (type: ToastType, message: string) => void;
}>({ toast: () => {} });

export const useToast = () => useContext(Ctx);

let nextId = 0;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([]);
  const { t } = useTranslation();

  const toast = useCallback((type: ToastType, message: string) => {
    const id = nextId++;
    const duration = type === 'error' ? 8000 : 4000;
    setToasts((prev) => [...prev, { id, type, message }]);
    setTimeout(() => setToasts((prev) => prev.filter((item) => item.id !== id)), duration);
  }, []);

  const dismiss = useCallback((id: number) => {
    setToasts((prev) => prev.filter((item) => item.id !== id));
  }, []);

  return (
    <Ctx.Provider value={{ toast }}>
      {children}
      <div
        aria-live="polite"
        className="fixed bottom-6 right-6 z-50 flex flex-col gap-2"
      >
        {toasts.map((item) => (
          <div
            key={item.id}
            className={`flex items-center gap-2 rounded-lg border px-4 py-3 shadow-lg backdrop-blur-sm transition-all duration-300 animate-in slide-in-from-right ${
              item.type === 'success'
                ? 'border-emerald-200 bg-emerald-50/90 text-emerald-800 dark:border-emerald-800 dark:bg-emerald-950/90 dark:text-emerald-200'
                : 'border-red-200 bg-red-50/90 text-red-800 dark:border-red-800 dark:bg-red-950/90 dark:text-red-200'
            }`}
          >
            {item.type === 'success' ? (
              <CheckCircle2 size={18} className="shrink-0" />
            ) : (
              <XCircle size={18} className="shrink-0" />
            )}
            <span className="text-sm font-medium">{item.message}</span>
            <button
              onClick={() => dismiss(item.id)}
              className="ml-2 shrink-0 rounded p-0.5 opacity-60 hover:opacity-100 transition-opacity"
              aria-label={t('common.closeToast')}
            >
              <X size={14} />
            </button>
          </div>
        ))}
      </div>
    </Ctx.Provider>
  );
}
