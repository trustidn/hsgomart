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
