const WEI_PER_SUDA = 10n ** 18n;

export function weiToSuda(weiStr: string, decimals = 2): string {
  const wei = BigInt(weiStr);
  const whole = wei / WEI_PER_SUDA;
  const frac = wei % WEI_PER_SUDA;
  if (decimals === 0) return whole.toString();
  const fracVal = (frac * 10n ** BigInt(decimals)) / WEI_PER_SUDA;
  const fracStr = fracVal.toString().padStart(decimals, '0');
  return `${whole}.${fracStr}`;
}

export function sudaToWei(suda: string): string {
  if (!suda || suda === '.' || suda === '') return '0';
  const [wholeRaw = '0', fracRaw = ''] = suda.split('.');
  const whole = wholeRaw || '0';
  const frac = (fracRaw + '0'.repeat(18)).slice(0, 18);
  return (BigInt(whole) * WEI_PER_SUDA + BigInt(frac)).toString();
}

export function truncateAddr(addr: string): string {
  if (!addr || addr.length < 10) return addr;
  return `${addr.slice(0, 6)}…${addr.slice(-4)}`;
}

export function formatAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  if (diff < 60_000) return 'just now';
  if (diff < 3_600_000) return `${Math.floor(diff / 60_000)}m ago`;
  if (diff < 86_400_000) {
    return `Today · ${new Date(iso).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })}`;
  }
  return `Yesterday · ${new Date(iso).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })}`;
}
