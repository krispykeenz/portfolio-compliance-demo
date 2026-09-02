const message = 'Enter a positive amount with at most two decimal places.'

export function parseMoney(input: string): number {
  const normalized = input.replaceAll(',', '').trim()
  const match = /^(\d+)(?:\.(\d{1,2}))?$/.exec(normalized)
  if (!match) throw new Error(message)
  const minor = Number(match[1]) * 100 + Number((match[2] ?? '').padEnd(2, '0'))
  if (!Number.isSafeInteger(minor) || minor <= 0) throw new Error(message)
  return minor
}
