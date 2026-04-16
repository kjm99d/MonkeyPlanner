// 개발 모드에서 접근성 위반을 실시간 감사. import.meta.env.DEV 에서만 동작.
// 빌드 산출물에는 포함되지 않음.

export async function enableAxe() {
  if (!import.meta.env.DEV) return;
  const React = await import('react');
  const ReactDOM = await import('react-dom');
  const axe = (await import('@axe-core/react')).default;
  axe(React, ReactDOM, 1000);
}
