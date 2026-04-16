import { Canvas, useFrame } from '@react-three/fiber';
import { useRef } from 'react';
import type { Mesh } from 'three';

function Banana({ position, speed }: { position: [number, number, number]; speed: number }) {
  const ref = useRef<Mesh>(null!);
  useFrame((state) => {
    const t = state.clock.elapsedTime * speed;
    ref.current.rotation.x = t * 0.4;
    ref.current.rotation.z = t * 0.6;
    ref.current.position.y = position[1] + Math.sin(t) * 0.2;
  });
  return (
    <mesh ref={ref} position={position}>
      <torusGeometry args={[0.6, 0.22, 16, 48, Math.PI * 1.4]} />
      <meshStandardMaterial color="#fbbf24" metalness={0.1} roughness={0.5} />
    </mesh>
  );
}

function MonkeyHead({ position, speed }: { position: [number, number, number]; speed: number }) {
  const ref = useRef<Mesh>(null!);
  useFrame((state) => {
    const t = state.clock.elapsedTime * speed;
    ref.current.rotation.y = Math.sin(t) * 0.3;
    ref.current.position.y = position[1] + Math.sin(t * 1.2) * 0.15;
  });
  return (
    <group ref={ref as never} position={position}>
      {/* 머리 */}
      <mesh>
        <sphereGeometry args={[0.55, 32, 32]} />
        <meshStandardMaterial color="#8B4513" roughness={0.6} />
      </mesh>
      {/* 얼굴 */}
      <mesh position={[0, -0.1, 0.45]}>
        <sphereGeometry args={[0.3, 24, 24]} />
        <meshStandardMaterial color="#DEB887" roughness={0.7} />
      </mesh>
      {/* 왼쪽 귀 */}
      <mesh position={[-0.5, 0.15, 0]}>
        <sphereGeometry args={[0.2, 16, 16]} />
        <meshStandardMaterial color="#8B4513" roughness={0.6} />
      </mesh>
      {/* 오른쪽 귀 */}
      <mesh position={[0.5, 0.15, 0]}>
        <sphereGeometry args={[0.2, 16, 16]} />
        <meshStandardMaterial color="#8B4513" roughness={0.6} />
      </mesh>
      {/* 왼쪽 눈 */}
      <mesh position={[-0.18, 0.1, 0.5]}>
        <sphereGeometry args={[0.07, 12, 12]} />
        <meshStandardMaterial color="#1a1a1a" />
      </mesh>
      {/* 오른쪽 눈 */}
      <mesh position={[0.18, 0.1, 0.5]}>
        <sphereGeometry args={[0.07, 12, 12]} />
        <meshStandardMaterial color="#1a1a1a" />
      </mesh>
    </group>
  );
}

export default function Hero3D() {
  return (
    <div
      className="hero-banner relative h-56 w-full overflow-hidden rounded-xl border border-edge-base bg-gradient-to-br from-brand-50 to-surface-subtle dark:from-brand-900/30 dark:to-surface-subtle"
      aria-hidden
    >
      <Canvas camera={{ position: [0, 0, 5], fov: 50 }} dpr={[1, 2]}>
        <ambientLight intensity={0.6} />
        <directionalLight position={[3, 4, 2]} intensity={1.1} />
        <MonkeyHead position={[-0.8, 0, 0]} speed={0.7} />
        <Banana position={[1.2, 0.2, -0.3]} speed={0.9} />
        <Banana position={[2.0, -0.3, -0.8]} speed={1.1} />
      </Canvas>
    </div>
  );
}
