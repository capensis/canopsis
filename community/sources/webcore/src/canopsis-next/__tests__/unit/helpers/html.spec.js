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

  it('replaces legacy icon tags while preserving icon attributes and content', () => {
    const html = '<V-ICON color="red" small class="custom-icon">home</V-ICON><c-icon size="32">close</c-icon>';
    const expected = '<c-icon color="red" small="" class="custom-icon">home</c-icon><c-icon size="32">close</c-icon>';

    expect(sanitizeHtml(html)).toBe(expected);
    expect(sanitizeHtml(expected)).toBe(expected);
  });

  it('does not replace icon-like text in attributes or escaped content', () => {
    const html = '<div title="<v-icon>home</v-icon>">text</div>&lt;v-icon&gt;home&lt;/v-icon&gt;';

    expect(sanitizeHtml(html)).toBe(html);
  });

  it('sanitizes attributes and descendants of replaced icons', () => {
    const html = '<v-icon tag="script" @click="alert(1)" :innerHTML="html" onload="alert(1)"><img src="x" onerror="alert(1)"><script>alert(1)</script>home</v-icon>';

    expect(sanitizeHtml(html)).toBe('<c-icon><img src="x">home</c-icon>');
  });

  it('applies custom tag restrictions to replaced icons', () => {
    expect(sanitizeHtml('<v-icon>home</v-icon>', { ALLOWED_TAGS: ['strong'] })).toBe('home');
  });

  it('sanitizes nested and adjacent replacement icons without skipping nodes', () => {
    const html = '<v-icon onload="alert(1)"><v-icon onclick="alert(1)">home</v-icon></v-icon><v-icon onmouseover="alert(1)">close</v-icon><img src="x" onerror="alert(1)">';

    expect(sanitizeHtml(html)).toBe('<c-icon><c-icon>home</c-icon></c-icon><c-icon>close</c-icon><img src="x">');
  });

  it('removes replacement icons when their tag and contents are forbidden', () => {
    expect(sanitizeHtml('<v-icon><img src="x">home</v-icon><strong>text</strong>', {
      FORBID_TAGS: ['c-icon'],
      KEEP_CONTENT: false,
    })).toBe('<strong>text</strong>');
  });

  it('does not leave an unsanitized replacement outside an in-place sanitization root', () => {
    const host = document.createElement('div');

    host.innerHTML = '<v-icon onclick="alert(1)">home</v-icon>';
    DOMPurify.sanitize(host.firstChild, {
      IN_PLACE: true,
      ADD_TAGS: ['v-icon', 'c-icon'],
    });

    expect(host.innerHTML).not.toContain('onclick');
    expect(host.textContent).toBe('home');
  });

  it.each([
    '<v-icon onclick="alert(1)"><img src="x" onerror="alert(1)"></v-icon>',
    '<v-icon><a href="jav&#x61;script:alert(1)">link</a></v-icon>',
    '<svg><v-icon onload="alert(1)"><a xlink:href="javascript:alert(1)">link</a></v-icon></svg>',
    '<math><mtext><v-icon><img src="x" onerror="alert(1)"></v-icon></mtext></math>',
    '<v-icon><svg><foreignObject><img src="x" onerror="alert(1)"></foreignObject></svg></v-icon>',
    '<v-icon><math><mtext><table><mglyph><style><!--</style><img title="--><img src=x onerror=alert(1)>"></v-icon>',
    '<v-icon><svg><p><style><g title="</style><img src=x onerror=alert(1)>"></v-icon>',
    '<v-icon><template><v-icon onclick="alert(1)">home</v-icon></template></v-icon>',
  ])('keeps replaced icon markup safe after reparsing: %s', (html) => {
    const element = document.createElement('template');
    const reparsed = document.createElement('template');
    const sanitized = sanitizeHtml(`<strong>safe</strong>${html}`);

    element.innerHTML = sanitized;
    reparsed.innerHTML = element.innerHTML;

    expect(reparsed.content.querySelector('strong').textContent).toBe('safe');
    expect(reparsed.content.querySelector('script, iframe, object, embed, foreignObject')).toBeNull();
    reparsed.content.querySelectorAll('*').forEach((node) => {
      Array.from(node.attributes).forEach((attribute) => {
        expect(attribute.name).not.toMatch(/^(?:on|v-|:|\.|@|#|data-v-)/i);
      });
      Array.from(node.attributes)
        .filter(attribute => /^(?:href|src|xlink:href)$/i.test(attribute.name))
        .forEach((attribute) => {
          expect(attribute.value).not.toMatch(/^\s*(?:javascript|vbscript|data):/i);
        });
    });
  });

  it('preserves Handlebars table comments when replacing icons', () => {
    const html = '<!--{{#if enabled}}--><table><tbody><!--{{#each items}}--><tr><td><v-icon>home</v-icon></td></tr><!--{{/each}}--></tbody></table><!--{{/if}}-->';
    const expected = '<!--{{#if enabled}}--><table><tbody><!--{{#each items}}--><tr><td><c-icon>home</c-icon></td></tr><!--{{/each}}--></tbody></table><!--{{/if}}-->';

    expect(sanitizeHtml(html)).toBe(expected);
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
