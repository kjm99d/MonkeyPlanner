// 최소 axe 러너. jest-axe 의존을 피하고 axe-core 직접 사용.
import axeCore from 'axe-core';

export async function axe(container: Element) {
  return axeCore.run(container, {
    rules: {
      // 테스트 환경은 색상 대비 검사 무의미(jsdom); 실제 감사는 dev 모드 @axe-core/react 가 담당.
      'color-contrast': { enabled: false },
    },
  });
}

type AxeResults = Awaited<ReturnType<typeof axeCore.run>>;

export const toHaveNoViolations = {
  toHaveNoViolations(received: AxeResults) {
    const violations = received.violations;
    const pass = violations.length === 0;
    return {
      pass,
      message: () =>
        pass
          ? 'expected violations but found none'
          : `axe violations (${violations.length}):\n` +
            violations
              .map(
                (v) =>
                  `- [${v.id}] ${v.help} (${v.nodes.length} nodes)\n  ${v.helpUrl}`,
              )
              .join('\n'),
    };
  },
};

declare module 'vitest' {
  interface Assertion {
    toHaveNoViolations(): void;
  }
}
