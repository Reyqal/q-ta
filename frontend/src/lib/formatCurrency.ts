export function formatRupiah(amount: number): string {
  return 'Rp ' + amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, '.');
}

export function formatRupiahShort(amount: number): string {
  if (amount >= 1_000_000_000) {
    return 'Rp ' + (amount / 1_000_000_000).toFixed(1).replace('.', ',') + ' M';
  }
  if (amount >= 1_000_000) {
    return 'Rp ' + (amount / 1_000_000).toFixed(1).replace('.', ',') + ' Jt';
  }
  if (amount >= 1_000) {
    return 'Rp ' + (amount / 1_000).toFixed(0) + ' Rb';
  }
  return formatRupiah(amount);
}
