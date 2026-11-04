import { mount } from '@vue/test-utils';

import { DEFAULT_SANITIZE_OPTIONS } from '@/config';

import { sanitizeHtml, sanitizeVueTemplate, sanitizeSvg, normalizeHtml } from '@/helpers/html';

import { TEMPLATE_VUE_ATTRIBUTES, TEXT_WIDGET_SANITIZE_OPTIONS } from '@/constants/sanitize';

const expression = "$root.constructor.constructor('window.templateXss = true')()";
const renderTemplate = html => mount({ template: sanitizeVueTemplate(`<div>${html}</div>`) });

describe('HTML and Vue template security boundaries', () => {
  let wrapper;

  afterEach(() => {
    wrapper?.destroy();
    delete window.templateXss;
  });

  it.each([
    'v-if', 'v-else-if', 'v-else', 'v-for', 'v-html', 'v-text', 'v-bind', 'v-on',
    ':title', 'v-bind:title', '@click', 'v-on:click', '@click.stop', 'v-on:click.once',
    ':[title]', 'v-bind:[title]', '@[event]', 'v-slot', '#default', '.title',
    'data-v-bind-title', 'data-v-on-click', 'is', 'tag', 'innerHTML', 'onload', 'onerror', 'inline-template',
  ])('rejects arbitrary attribute %s even with custom options', (attribute) => {
    const result = sanitizeVueTemplate(`<span ${attribute}="${expression}">text</span>`, {
      allowedAttributes: { '*': ['*'] },
    });
    wrapper = mount({ template: `<div>${result}</div>` });

    expect(wrapper.find('span').attributes()).toEqual({});
    expect(window.templateXss).toBeUndefined();
  });

  it.each([
    `{{ ${expression} }}`,
    `&#123;&#123; ${expression} &#125;&#125;`,
    `&#x7b;{ ${expression} }}`,
    `&lbrace;{ ${expression} }}`,
    `{{{ ${expression} }}}`,
    `{<!---->{ ${expression} }}`,
  ])('keeps interpolation as text after normalization: %s', (value) => {
    wrapper = renderTemplate(normalizeHtml(`<p>${value}</p>`));

    expect(wrapper.text()).toContain(expression);
    expect(wrapper.text()).toContain('{{');
    expect(window.templateXss).toBeUndefined();
  });

  it('also protects template contents allowed by a special profile', () => {
    const html = sanitizeVueTemplate(`<template><span v-if="${expression}">{{ ${expression} }}</span></template>`, {
      allowedTags: ['div', 'template', 'span'],
      allowedAttributes: { '*': ['*'] },
    });
    wrapper = mount({ template: `<div>${html}</div>` });

    expect(wrapper.text()).toBe(`{{ ${expression} }}`);
    expect(window.templateXss).toBeUndefined();
  });

  it('validates decoded values, encoded attribute names and custom transforms', () => {
    const html = sanitizeVueTemplate('<c-links-list data-v-bind-links="links" v-if="links &amp;&amp; hasAccessForLinks" :links="&#97;lert(1)"></c-links-list>', {
      allowedAttributes: { '*': ['*'] },
      transformTags: { '*': (tagName, attribs) => ({ tagName, attribs: { ...attribs, onload: 'bad' } }) },
    });

    expect(html).toContain('v-if="links &amp;&amp; hasAccessForLinks"');
    expect(html).not.toContain('data-v-');
    expect(html).not.toContain(':links');
    expect(html).not.toContain('onload');
  });

  it('retains only the exact helper tag, attribute and expression triples', () => {
    Object.entries(TEMPLATE_VUE_ATTRIBUTES).forEach(([tag, attributes]) => {
      Object.entries(attributes).forEach(([name, value]) => {
        const escaped = value.replace(/&/g, '&amp;').replace(/"/g, '&quot;');
        const fragment = document.createElement('template');
        fragment.innerHTML = sanitizeHtml(`<${tag} ${name}="${escaped}"></${tag}>`);
        expect(fragment.content.firstChild.getAttribute(name)).toBe(value);
        expect(sanitizeHtml(`<${tag} ${name}="${expression}"></${tag}>`)).not.toContain(name);
        expect(sanitizeHtml(`<span ${name}="${escaped}"></span>`)).not.toContain(name);
      });
    });
  });

  it('preserves ordinary rich text, links and a special iframe profile', () => {
    const source = '<p style="color:red"><strong>rich</strong><a href="https://example.com">link</a></p>';
    wrapper = renderTemplate(source);

    expect(wrapper.find('strong').text()).toBe('rich');
    expect(wrapper.find('a').attributes('href')).toBe('https://example.com');
    expect(sanitizeHtml('<a href="javascript:alert(1)">bad</a>')).not.toContain('href');
    expect(sanitizeHtml(
      '<iframe src="https://example.com" srcdoc="bad" v-if="bad"></iframe>',
      TEXT_WIDGET_SANITIZE_OPTIONS,
    )).toBe('<iframe src="https://example.com"></iframe>');
    expect(DEFAULT_SANITIZE_OPTIONS.allowedTags).not.toContain('svg');
  });

  it.each([undefined, TEXT_WIDGET_SANITIZE_OPTIONS, { disallowedTagsMode: 'discard' }])('preserves safe Unicode comments in table positions with options %s', (options) => {
    const source = '<!--до--><table><tbody><!--{{#each rows}}--><!--{{#if enabled}}--><tr><!--строка--><td>yes</td></tr><!--{{else}}--><tr><td>no</td></tr><!--{{/if}}--><!--{{/each}}--></tbody></table><!--после-->';
    const result = sanitizeHtml(source, options);

    expect(result).toBe(source);
    expect(normalizeHtml(result)).toBe(source);
    expect(sanitizeHtml(result, options)).toBe(source);
  });

  it.each([
    '<div class="html-comment" data-c="LS0+PGltZyBzcmM9eCBvbmVycm9yPWFsZXJ0KDEpPg=="></div>',
    '<!-- --!><img src="x" onerror="window.templateXss=true">',
    '<!--a--!><span v-if="bad">text</span><!--b-->',
    '<svg><foreignObject><!-- --!><img src="x" onerror="window.templateXss=true"></foreignObject></svg>',
    '<!--a<!--b--><img src="x" onerror="window.templateXss=true">',
  ])('cannot restore unsanitized HTML from comments or forged placeholders: %s', (html) => {
    wrapper = renderTemplate(html);

    expect(wrapper.findAll('[onerror], [v-if]')).toHaveLength(0);
    expect(window.templateXss).toBeUndefined();
    expect(sanitizeHtml(sanitizeHtml(html))).not.toContain('onerror="');
  });

  it('does not recognize comments in attributes or raw-text contents', () => {
    const result = sanitizeHtml('<p title="<!--{{#if enabled}}-->">text</p><script><!--unsafe--></script>');

    expect(result).toContain('title="&lt;!--{{#if enabled}}--&gt;"');
    expect(result).not.toContain('canopsisComment');
    expect(result).not.toContain('<!--unsafe-->');
  });
});

describe('separate SVG policy', () => {
  it('keeps case-sensitive geometry, gradients, clipping and local references', () => {
    const source = '<svg viewBox="0 0 24 24" class="icon"><defs><linearGradient id="paint" gradientTransform="rotate(30)"><stop offset="0" stop-color="red"/></linearGradient><clipPath id="clip"><path d="M0 0"/></clipPath><filter id="shadow"><feGaussianBlur stdDeviation="2"/></filter></defs><g transform="scale(2)" clip-path="url(#clip)"><path id="shape" d="M0 0" fill="url(#paint)" stroke="currentColor"/><use href="#shape" xlink:href="#shape"/></g></svg>';
    const result = sanitizeSvg(source);

    expect(result).toContain('viewBox="0 0 24 24"');
    expect(result).toContain('linearGradient');
    expect(result).toContain('gradientTransform="rotate(30)"');
    expect(result).toContain('clipPath');
    expect(result).toContain('feGaussianBlur stdDeviation="2"');
    expect(result).toContain('transform="scale(2)"');
    expect(result).toContain('fill="url(#paint)"');
    expect(result).toContain('xlink:href="#shape"');
    expect(sanitizeSvg(result)).toBe(result);
  });

  it('removes executable HTML, Vue attributes, animation and external resources', () => {
    const result = sanitizeSvg(`<svg onload="bad" v-if="${expression}"><foreignObject><img src="x" onerror="bad"/></foreignObject><script>bad</script><animate attributeName="href"/><use href="javascript:bad"/><use xlink:href="https://example.com/icon"/><path fill="url(https://example.com/paint)" style="fill:red" d="M0 0"/></svg>`);

    expect(result).not.toMatch(/foreignObject|script|animate|onload|onerror|v-if|javascript:|https:|style=/i);
    expect(result).toContain('d="M0 0"');
    expect(sanitizeSvg('<p>no SVG</p>')).toBe('');
  });
});
