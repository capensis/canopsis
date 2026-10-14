import DOMPurify from 'dompurify';

import { sanitizeHtml } from '@/helpers/html';

describe('sanitizeHtml', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('caches sanitized HTML when the default options are used', () => {
    const sanitizeSpy = vi.spyOn(DOMPurify, 'sanitize');
    const html = '<strong>cached sanitize value</strong>';

    const firstResult = sanitizeHtml(html);
    const secondResult = sanitizeHtml(html);

    expect(secondResult).toBe(firstResult);
    expect(sanitizeSpy).toHaveBeenCalledTimes(1);
  });

  it('does not cache HTML when custom options are used', () => {
    const sanitizeSpy = vi.spyOn(DOMPurify, 'sanitize');
    const html = '<strong>custom sanitize options value</strong>';
    const options = { ALLOWED_TAGS: ['strong'] };

    sanitizeHtml(html, options);
    sanitizeHtml(html, options);

    expect(sanitizeSpy).toHaveBeenCalledTimes(2);
  });

  it('normalizes malformed HTML while sanitizing it', () => {
    expect(sanitizeHtml('<div><span>value</div>')).toBe('<div><span>value</span></div>');
  });
});
