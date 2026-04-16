import { Canvas, useFrame } from '@react-three/fiber';
import { useRef } from 'react';
import type { Mesh } from 'three';

// 홈 Hero: 코숭이 브랜드 실루엣 대신 "부유하는 바나나색 원환(torus) + 큐브" 2개 메시.
// R3F + three 는 이 파일을 import 한 청크에서만 로드되며, lazy() 로 분리됨.

function FloatingShape({
  position,
  geometry,
  color,
  speed,
}: {
  position: [number, number, number];
  geometry: 'torus' | 'box';
  color: string;
  speed: number;
}) {
  const ref = useRef<Mesh>(null!);
  useFrame((state) => {
    const t = state.clock.elapsedTime * speed;
    ref.current.rotation.x = t * 0.5;
    ref.current.rotation.y = t * 0.7;
    ref.current.position.y = position[1] + Math.sin(t) * 0.25;
  });
  return (
    <mesh ref={ref} position={position}>
      {geometry === 'torus' ? (
        <torusGeometry args={[0.75, 0.28, 16, 64]} />
      ) : (
        <boxGeometry args={[1, 1, 1]} />
      )}
      <meshStandardMaterial color={color} metalness={0.2} roughness={0.45} />
    </mesh>
  );
}

export default function Hero3D() {
  return (
    <div
      className="relative h-56 w-full overflow-hidden rounded-xl border border-edge-base bg-gradient-to-br from-brand-50 to-surface-subtle dark:from-brand-900/30 dark:to-surface-subtle"
      aria-hidden
    >
      <Canvas camera={{ position: [0, 0, 5], fov: 50 }} dpr={[1, 2]}>
        <ambientLight intensity={0.6} />
        <directionalLight position={[3, 4, 2]} intensity={1.1} />
        <FloatingShape position={[-1.5, 0, 0]} geometry="torus" color="#f97316" speed={0.8} />
        <FloatingShape position={[1.5, 0.3, -0.4]} geometry="box" color="#fbbf24" speed={1.2} />
      </Canvas>
    </div>
  );
}
