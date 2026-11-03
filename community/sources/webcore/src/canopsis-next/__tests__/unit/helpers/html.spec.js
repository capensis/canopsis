import DOMPurify from 'dompurify';

import { sanitizeHtml } from '@/helpers/html';

import { TEMPLATE_VUE_ATTRIBUTES } from '@/constants/sanitize';

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

  it('does not restore attacker-controlled comment placeholders after sanitization', () => {
    const encoded = btoa('--><img src=x onerror="window.commentXss = true"><!--');
    const result = sanitizeHtml(`<div class="html-comment" data-c="${encoded}"></div>`);

    expect(result).not.toContain('<img');
    expect(result).not.toContain('onerror');
  });

  it('removes comments without decoding them into markup', () => {
    expect(sanitizeHtml('before<!-- <img src=x onerror="alert(1)"> -->after')).toBe('beforeafter');
  });

  it('preserves plain and Unicode comments', () => {
    const html = '<div><!-- Условие шаблона: {{#if enabled}} -->content<!-- {{/if}} --></div>';

    expect(sanitizeHtml(html)).toBe(html);
  });

  it('preserves Handlebars condition comments between table rows without moving them', () => {
    const html = '<table><tbody><!--{{#if enabled}}--><tr><td>yes</td></tr><!--{{else}}--><tr><td>no</td></tr><!--{{/if}}--></tbody></table>';

    expect(sanitizeHtml(html)).toBe(html);
    expect(sanitizeHtml(sanitizeHtml(html))).toBe(html);
  });

  it('preserves comments before and after a table', () => {
    const html = '<!--{{#if enabled}}--><table><tbody><tr><td>yes</td></tr></tbody></table><!--{{/if}}-->';

    expect(sanitizeHtml(html)).toBe(html);
  });

  it('preserves a standalone comment without interpreting its contents', () => {
    expect(sanitizeHtml('<!--{{#if enabled}}-->')).toBe('<!--{{#if enabled}}-->');
  });

  it('sanitizes markup after a malformed comment terminator', () => {
    const html = '<!--comment--!><img src="x" onerror="window.commentXss = true">';

    expect(sanitizeHtml(html)).toBe('<!--comment--><img src="x">');
  });

  it.each(Object.entries(TEMPLATE_VUE_ATTRIBUTES))('preserves only the exact helper attributes on %s', (tag, attributes) => {
    const htmlAttributes = Object.entries(attributes)
      .map(([name, value]) => `${name}="${value}"`)
      .join(' ');
    const element = document.createElement('template');

    element.innerHTML = sanitizeHtml(`<${tag} ${htmlAttributes}></${tag}>`);

    Object.entries(attributes).forEach(([name, value]) => {
      expect(element.content.firstChild.getAttribute(name)).toBe(value);
    });
  });

  it.each(Object.entries(TEMPLATE_VUE_ATTRIBUTES))('rejects arbitrary expressions in helper attributes on %s', (tag, attributes) => {
    const htmlAttributes = Object.keys(attributes)
      .map(name => `${name}="constructor.constructor('alert(1)')()"`)
      .join(' ');
    const element = document.createElement('template');

    element.innerHTML = sanitizeHtml(`<${tag} ${htmlAttributes}></${tag}>`);

    expect(element.content.firstChild.attributes).toHaveLength(0);
  });
});
