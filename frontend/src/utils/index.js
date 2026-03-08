// Shared utilities for HSMart POS frontend
export function formatCurrency(value) {
  if (value == null) return '—'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(value)
}

export function formatPrice(value) {
  return formatCurrency(value ?? 0)
}

/**
 * Parse backend timestamp correctly.
 * The DB stores "timestamp without time zone" but pgx reads it as UTC.
 * Strip timezone marker so the browser treats it as local time.
 */
export function parseDate(iso) {
  if (!iso) return null
  if (iso instanceof Date) return iso
  const s = String(iso).replace('Z', '').replace(/[+-]\d{2}:\d{2}$/, '')
  return new Date(s)
}

export function formatDate(iso) {
  const d = parseDate(iso)
  if (!d || isNaN(d)) return '—'
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

export function formatDateTime(iso) {
  const d = parseDate(iso)
  if (!d || isNaN(d)) return '—'
  return d.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
}
