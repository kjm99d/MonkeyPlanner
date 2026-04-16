type Props = {
  className?: string;
  count?: number;
};

export function Skeleton({ className = 'h-4 w-full', count = 1 }: Props) {
  return (
    <>
      {Array.from({ length: count }, (_, i) => (
        <div
          key={i}
          className={`animate-pulse rounded-md bg-surface-muted ${className}`}
          aria-hidden
        />
      ))}
    </>
  );
}

export function CardSkeleton() {
  return (
    <div className="flex flex-col gap-3 rounded-lg border border-edge-base bg-surface-subtle p-4">
      <Skeleton className="h-5 w-3/4" />
      <Skeleton className="h-3 w-1/2" />
    </div>
  );
}

export function BoardSkeleton() {
  return (
    <div className="grid gap-4 lg:grid-cols-4 md:grid-cols-2">
      {[1, 2, 3, 4].map((i) => (
        <div key={i} className="flex min-h-[16rem] flex-col gap-3 rounded-lg border border-edge-base bg-surface-subtle p-3">
          <Skeleton className="h-4 w-20" />
          <CardSkeleton />
          <CardSkeleton />
        </div>
      ))}
    </div>
  );
}
