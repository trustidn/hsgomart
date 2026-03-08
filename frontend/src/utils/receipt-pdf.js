import { jsPDF } from 'jspdf'
import { parseDate } from './index'

const FMT = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 })
function fmt(v) { return FMT.format(v ?? 0) }

function fmtDate(d) {
  const dt = parseDate(d)
  if (!dt || isNaN(dt)) return '—'
  return dt.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
}

export function generateReceiptPDF(data) {
  const doc = new jsPDF({ unit: 'mm', format: [80, 200] })
  const w = 80
  let y = 8

  doc.setFontSize(12)
  doc.text(data.storeName || 'HSMart', w / 2, y, { align: 'center' })
  y += 5
  doc.setFontSize(7)
  doc.text(fmtDate(data.date), w / 2, y, { align: 'center' })
  y += 4
  if (data.transactionId) {
    doc.text(`#${data.transactionId.slice(0, 8)}`, w / 2, y, { align: 'center' })
    y += 4
  }
  if (data.cashier) {
    doc.text(`Cashier: ${data.cashier}`, w / 2, y, { align: 'center' })
    y += 4
  }

  doc.setDrawColor(180)
  doc.setLineDashPattern([1, 1], 0)
  doc.line(4, y, w - 4, y)
  y += 4

  doc.setFontSize(7)
  const items = data.items || []
  for (const item of items) {
    const name = (item.name || item.product_name || '').substring(0, 20)
    const qty = item.quantity
    const price = item.price ?? item.sell_price ?? 0
    const sub = price * qty
    doc.text(name, 4, y)
    doc.text(`${qty}x`, 46, y, { align: 'right' })
    doc.text(fmt(price), 60, y, { align: 'right' })
    doc.text(fmt(sub), w - 4, y, { align: 'right' })
    y += 3.5
  }

  doc.line(4, y, w - 4, y)
  y += 4

  doc.setFontSize(8)
  const rows = [
    ['Total', fmt(data.total)],
    ['Payment', fmt(data.paidAmount)],
    ['Change', fmt(data.change)],
  ]
  for (const [label, val] of rows) {
    doc.text(label, 4, y)
    doc.text(val, w - 4, y, { align: 'right' })
    y += 4
  }

  y += 2
  doc.line(4, y, w - 4, y)
  y += 4
  doc.setFontSize(7)
  doc.text('Thank you', w / 2, y, { align: 'center' })

  return doc
}

export function buildReceiptText(data) {
  const lines = []
  lines.push(`*${data.storeName || 'HSMart'}*`)
  lines.push(fmtDate(data.date))
  if (data.transactionId) lines.push(`Trx: #${data.transactionId.slice(0, 8)}`)
  if (data.cashier) lines.push(`Cashier: ${data.cashier}`)
  lines.push('---')
  for (const item of (data.items || [])) {
    const name = item.name || item.product_name || ''
    const qty = item.quantity
    const price = item.price ?? item.sell_price ?? 0
    lines.push(`${name} x${qty} = ${fmt(price * qty)}`)
  }
  lines.push('---')
  lines.push(`*Total: ${fmt(data.total)}*`)
  lines.push(`Paid: ${fmt(data.paidAmount)}`)
  lines.push(`Change: ${fmt(data.change)}`)
  lines.push('')
  lines.push('Thank you!')
  return lines.join('\n')
}
