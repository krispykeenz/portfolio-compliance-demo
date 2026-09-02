import { describe, expect, it } from 'vitest'
import { parseMoney } from './money'

describe('parseMoney', () => {
  it.each([['9,999.00',999900],['10000',1000000],['10,001.00',1000100]])('parses %s exactly', (input, expected) => {
    expect(parseMoney(input)).toBe(expected)
  })
  it.each(['', '0', '-1.00', '1.001', 'abc'])('rejects %s', input => {
    expect(() => parseMoney(input)).toThrow('Enter a positive amount with at most two decimal places.')
  })
})
