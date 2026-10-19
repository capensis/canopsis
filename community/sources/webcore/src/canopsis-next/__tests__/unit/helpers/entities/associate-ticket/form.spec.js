import { isTicketIdMatchingFormats } from '@/helpers/entities/associate-ticket/event/form';

describe('associate ticket helpers', () => {
  test('allows every ticket when no formats are configured', () => {
    expect(isTicketIdMatchingFormats('ANY-TICKET')).toBe(true);
  });

  test('accepts a ticket matching at least one configured format', () => {
    expect(isTicketIdMatchingFormats('INC-42', ['^REQ-[0-9]+$', '^INC-[0-9]+$'])).toBe(true);
  });

  test('rejects a ticket not matching any configured format', () => {
    expect(isTicketIdMatchingFormats('OTHER-42', ['^REQ-[0-9]+$', '^INC-[0-9]+$'])).toBe(false);
  });

  test('treats malformed configured regexps as non-matching', () => {
    expect(isTicketIdMatchingFormats('INC-42', ['[', '^INC-[0-9]+$'])).toBe(true);
    expect(isTicketIdMatchingFormats('INC-42', ['['])).toBe(false);
  });
});
