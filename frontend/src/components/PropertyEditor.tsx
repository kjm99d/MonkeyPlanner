import { useState } from 'react';
import { Plus, X, Tag, Hash, Calendar, CheckSquare, List, Type } from 'lucide-react';
import type { BoardProperty, PropertyType } from '../api/types';

const typeIcons: Record<PropertyType, typeof Type> = {
  text: Type, number: Hash, select: List, multi_select: Tag, date: Calendar, checkbox: CheckSquare,
};

const typeLabels: Record<PropertyType, string> = {
  text: '텍스트', number: '숫자', select: '선택', multi_select: '다중 선택', date: '날짜', checkbox: '체크박스',
};

type Props = {
  properties: BoardProperty[];
  values: Record<string, unknown>;
  onChange: (propId: string, value: unknown) => void;
};

export function PropertyEditor({ properties, values, onChange }: Props) {
  if (properties.length === 0) return null;

  return (
    <section className="flex flex-col gap-3 rounded-xl border border-edge-base bg-surface-subtle p-4">
      <h3 className="text-sm font-semibold text-ink-secondary">속성</h3>
      <div className="flex flex-col gap-2.5">
        {properties.map((prop) => (
          <PropertyField key={prop.id} prop={prop} value={values[prop.id]} onChange={(v) => onChange(prop.id, v)} />
        ))}
      </div>
    </section>
  );
}

function PropertyField({ prop, value, onChange }: { prop: BoardProperty; value: unknown; onChange: (v: unknown) => void }) {
  const Icon = typeIcons[prop.type];

  return (
    <div className="flex items-center gap-3">
      <div className="flex w-32 shrink-0 items-center gap-1.5 text-xs text-ink-muted">
        <Icon size={13} />
        <span className="truncate">{prop.name}</span>
      </div>
      <div className="flex-1">
        {prop.type === 'text' && (
          <input
            type="text"
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value)}
            className="h-8 w-full rounded-md border border-edge-base bg-surface-base px-2 text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none"
            placeholder="입력..."
          />
        )}
        {prop.type === 'number' && (
          <input
            type="number"
            value={(value as number) ?? ''}
            onChange={(e) => onChange(e.target.value ? Number(e.target.value) : null)}
            className="h-8 w-full rounded-md border border-edge-base bg-surface-base px-2 text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none"
            placeholder="0"
          />
        )}
        {prop.type === 'date' && (
          <input
            type="date"
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value || null)}
            className="h-8 w-full rounded-md border border-edge-base bg-surface-base px-2 text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none"
          />
        )}
        {prop.type === 'checkbox' && (
          <label className="flex cursor-pointer items-center gap-2">
            <input
              type="checkbox"
              checked={!!value}
              onChange={(e) => onChange(e.target.checked)}
              className="h-4 w-4 rounded border-edge-base accent-brand-500"
            />
            <span className="text-sm text-ink-secondary">{value ? '완료' : '미완료'}</span>
          </label>
        )}
        {prop.type === 'select' && (
          <select
            value={(value as string) ?? ''}
            onChange={(e) => onChange(e.target.value || null)}
            className="h-8 w-full rounded-md border border-edge-base bg-surface-base px-2 text-sm text-ink-primary focus-visible:border-brand-500 focus-visible:outline-none"
          >
            <option value="">선택...</option>
            {prop.options.map((opt) => (
              <option key={opt} value={opt}>{opt}</option>
            ))}
          </select>
        )}
        {prop.type === 'multi_select' && (
          <MultiSelectField options={prop.options} value={(value as string[]) ?? []} onChange={onChange} />
        )}
      </div>
    </div>
  );
}

function MultiSelectField({ options, value, onChange }: { options: string[]; value: string[]; onChange: (v: unknown) => void }) {
  const toggle = (opt: string) => {
    const next = value.includes(opt) ? value.filter((v) => v !== opt) : [...value, opt];
    onChange(next.length > 0 ? next : null);
  };

  return (
    <div className="flex flex-wrap gap-1.5">
      {options.map((opt) => {
        const active = value.includes(opt);
        return (
          <button
            key={opt}
            type="button"
            onClick={() => toggle(opt)}
            className={`rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors ${
              active
                ? 'border-brand-500 bg-brand-500/15 text-brand-500'
                : 'border-edge-base text-ink-muted hover:border-brand-500/30 hover:text-ink-secondary'
            }`}
          >
            {opt}
          </button>
        );
      })}
    </div>
  );
}

// 속성 추가 폼 (보드 설정용)
type AddPropertyProps = {
  onAdd: (name: string, type: PropertyType, options: string[]) => void;
};

export function AddPropertyForm({ onAdd }: AddPropertyProps) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState('');
  const [type, setType] = useState<PropertyType>('text');
  const [optStr, setOptStr] = useState('');

  const submit = () => {
    if (!name.trim()) return;
    const options = (type === 'select' || type === 'multi_select')
      ? optStr.split(',').map((s) => s.trim()).filter(Boolean)
      : [];
    onAdd(name.trim(), type, options);
    setName('');
    setType('text');
    setOptStr('');
    setOpen(false);
  };

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="flex items-center gap-1.5 rounded-lg border border-dashed border-edge-base px-3 py-2 text-xs text-ink-muted transition-colors hover:border-brand-500/30 hover:text-brand-500"
      >
        <Plus size={14} /> 속성 추가
      </button>
    );
  }

  return (
    <div className="flex flex-col gap-2 rounded-lg border border-edge-base bg-surface-subtle p-3">
      <div className="flex items-center justify-between">
        <span className="text-xs font-medium text-ink-secondary">새 속성</span>
        <button type="button" onClick={() => setOpen(false)} className="text-ink-muted hover:text-ink-primary">
          <X size={14} />
        </button>
      </div>
      <input
        placeholder="속성 이름"
        value={name}
        onChange={(e) => setName(e.target.value)}
        className="h-8 rounded-md border border-edge-base bg-surface-base px-2 text-sm focus-visible:border-brand-500 focus-visible:outline-none"
      />
      <select
        value={type}
        onChange={(e) => setType(e.target.value as PropertyType)}
        className="h-8 rounded-md border border-edge-base bg-surface-base px-2 text-sm focus-visible:border-brand-500 focus-visible:outline-none"
      >
        {(Object.keys(typeLabels) as PropertyType[]).map((t) => (
          <option key={t} value={t}>{typeLabels[t]}</option>
        ))}
      </select>
      {(type === 'select' || type === 'multi_select') && (
        <input
          placeholder="옵션 (쉼표 구분: P0, P1, P2)"
          value={optStr}
          onChange={(e) => setOptStr(e.target.value)}
          className="h-8 rounded-md border border-edge-base bg-surface-base px-2 text-sm focus-visible:border-brand-500 focus-visible:outline-none"
        />
      )}
      <button
        type="button"
        onClick={submit}
        className="h-8 rounded-md bg-brand-500 text-sm font-medium text-white hover:bg-brand-600 transition-colors"
      >
        추가
      </button>
    </div>
  );
}
