import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { axe, toHaveNoViolations } from './jest-axe-lite';

import { Button } from '../components/Button';
import { StatusBadge } from '../components/StatusBadge';
import { Input } from '../components/Input';

expect.extend(toHaveNoViolations);

function wrap(ui: React.ReactNode) {
  const qc = new QueryClient();
  return (
    <MemoryRouter>
      <QueryClientProvider client={qc}>{ui}</QueryClientProvider>
    </MemoryRouter>
  );
}

describe('a11y smoke', () => {
  it('Button has no axe violations', async () => {
    const { container } = render(wrap(<Button>클릭</Button>));
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('Input with label has no axe violations', async () => {
    const { container } = render(wrap(<Input label="제목" />));
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('StatusBadge has no axe violations', async () => {
    const { container } = render(wrap(<StatusBadge status="Pending" />));
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });
});
