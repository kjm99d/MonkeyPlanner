// prefers-reduced-motion 또는 3D 로드 실패 시 대체될 정적 Hero.

export default function Hero3DFallback() {
  return (
    <div
      role="img"
      aria-label="몽키 플래너 브랜드 일러스트"
      className="hero-banner flex h-56 w-full items-center justify-center overflow-hidden rounded-xl border border-edge-base bg-gradient-to-br from-brand-50 to-surface-subtle dark:from-brand-900/30 dark:to-surface-subtle"
    >
      <div className="flex items-center gap-4 text-6xl motion-safe:animate-pulse">
        <span aria-hidden>🍌</span>
        <span aria-hidden className="text-4xl text-brand-500">✳</span>
        <span aria-hidden>🐒</span>
      </div>
    </div>
  );
}
