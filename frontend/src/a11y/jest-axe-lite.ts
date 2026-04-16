// 최소 axe 러너. jest-axe 의존을 피하고 axe-core 직접 사용.
import axeCore, { type AxeResults, type Result } from 'axe-core';

export async function axe(container: Element): Promise<AxeResults> {
  return new Promise<AxeResults>((resolve, reject) => {
    axeCore.run(
      container,
      {
        rules: {
          // jsdom 에서 색상 대비 검사 신뢰할 수 없음; dev 런타임 axe가 별도 감사.
          'color-contrast': { enabled: false },
        },
      },
      (err, results) => {
        if (err) reject(err);
        else resolve(results);
      },
    );
  });
}

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
                (v: Result) =>
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
