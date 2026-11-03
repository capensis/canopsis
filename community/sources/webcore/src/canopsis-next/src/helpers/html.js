import DOMPurify from 'dompurify';
import linkifyHtmlLib from 'linkify-html';
import { registerCustomProtocol } from 'linkifyjs';

import { DEFAULT_SANITIZE_OPTIONS, DEFAULT_LINKIFY_OPTIONS, LINKIFY_PROTOCOLS } from '@/config';

import { TEMPLATE_VUE_ATTRIBUTES } from '@/constants/sanitize';

const VUE_ATTRIBUTE = /^(?:v-|:|\.|@|#|data-v-)/i;
const UNSAFE_TEMPLATE_ATTRIBUTE = /^(?:on|is$|tag$|ref$|innerhtml$|textcontent$|srcdoc$|slot$|slot-scope$|inline-template$)/i;

DOMPurify.addHook('uponSanitizeAttribute', (node, hookEvent) => {
  const data = hookEvent;
  const { attrName, attrValue } = data;

  if (VUE_ATTRIBUTE.test(attrName)) {
    const name = attrName
      .replace(/^data-v-bind-/, ':')
      .replace(/^data-v-on-/, '@');
    const allowed = TEMPLATE_VUE_ATTRIBUTES[node.localName];

    data.keepAttr = Boolean(allowed && Object.hasOwn(allowed, name) && allowed[name] === attrValue);

    if (data.keepAttr) {
      /**
       * Encoded event names may contain a colon (e.g. @update:options).
       * Approve only the exact pair above; ordinary data-attribute rules reject those names.
       */
      data.allowedAttributes[attrName] = true;
    }
  } else if (UNSAFE_TEMPLATE_ATTRIBUTE.test(attrName)) {
    data.keepAttr = false;
  }
});

/**
 * Encodes Vue attributes for DOMPurify. The hook above validates their tag, name and value.
 *
 * Replaces `:attr` with `data-v-bind-attr` and `@event` with `data-v-on-event`.
 *
 * @param {string} [html = ''] - HTML string that may contain Vue bindings.
 * @returns {string} HTML with Vue attributes encoded as data attributes.
 */
export const encodeVueAttributes = (html = '') =>
  html
    .replace(/(\s):([^\s=>/]+)(?=\s*=)/g, '$1data-v-bind-$2')
    .replace(/(\s)@([^\s=>/]+)(?=\s*=)/g, '$1data-v-on-$2');

/**
 * Restores Vue template attributes from their encoded data-attribute form.
 *
 * Reverses {@link encodeVueAttributes} by converting `data-v-bind-attr` to `:attr`
 * and `data-v-on-event` to `@event`.
 *
 * @param {string} [html = ''] - HTML string with encoded Vue attributes.
 * @returns {string} HTML with Vue attributes restored.
 */
export const decodeVueAttributes = (html = '') =>
  html
    .replace(/(\s)data-v-bind-([^\s=>/]+)(?=\s*=)/g, '$1:$2')
    .replace(/(\s)data-v-on-([^\s=>/]+)(?=\s*=)/g, '$1@$2');

LINKIFY_PROTOCOLS.forEach(registerCustomProtocol);

const SANITIZE_HTML_CACHE_LIMIT = 500;
const SANITIZE_HTML_MAX_CACHEABLE_LENGTH = 20_000;
const sanitizeHtmlCache = new Map();

/**
 * Returns a cached sanitized value and refreshes its LRU position.
 *
 * @param {string} html
 * @returns {string | undefined}
 */
const getCachedSanitizedHtml = (html) => {
  if (!sanitizeHtmlCache.has(html)) {
    return undefined;
  }

  const sanitizedHtml = sanitizeHtmlCache.get(html);

  sanitizeHtmlCache.delete(html);
  sanitizeHtmlCache.set(html, sanitizedHtml);

  return sanitizedHtml;
};

/**
 * Stores a sanitized value in the bounded LRU cache.
 *
 * @param {string} html
 * @param {string} sanitizedHtml
 */
const cacheSanitizedHtml = (html, sanitizedHtml) => {
  sanitizeHtmlCache.set(html, sanitizedHtml);

  if (sanitizeHtmlCache.size > SANITIZE_HTML_CACHE_LIMIT) {
    sanitizeHtmlCache.delete(sanitizeHtmlCache.keys().next().value);
  }
};

/**
 * Sanitizes HTML while preserving helper attributes and allowed HTML comments.
 *
 * Use sanitizeVueTemplate after Handlebars compilation before passing HTML to Vue.
 * DOMPurify preserves safe comments via #comment in the allowlist. Do not decode
 * comment placeholders or restore untrusted markup after sanitization.
 *
 * @param {string} [html = ''] - Raw HTML to sanitize.
 * @param {import('dompurify').Config} [options = DEFAULT_SANITIZE_OPTIONS] - DOMPurify configuration.
 * @returns {string} Sanitized HTML safe for rendering.
 */
export const sanitizeHtml = (html = '', options = DEFAULT_SANITIZE_OPTIONS) => {
  const cacheable =
    options === DEFAULT_SANITIZE_OPTIONS && html.length <= SANITIZE_HTML_MAX_CACHEABLE_LENGTH;
  const cachedHtml = cacheable ? getCachedSanitizedHtml(html) : undefined;

  if (cachedHtml !== undefined) {
    return cachedHtml;
  }

  const sanitizedHtml = decodeVueAttributes(
    DOMPurify.sanitize(encodeVueAttributes(html), options),
  );

  if (cacheable) {
    cacheSanitizedHtml(html, sanitizedHtml);
  }

  return sanitizedHtml;
};

/**
 * Sanitizes compiled Handlebars output for the Vue template boundary.
 * HTML entities keep interpolation delimiters as literal text in Vue's parser.
 * Do this last: another HTML normalization pass would decode them back to {{.
 *
 * @param {string} [html = '']
 * @param {import('dompurify').Config} [options = DEFAULT_SANITIZE_OPTIONS]
 * @returns {string}
 */
export const sanitizeVueTemplate = (html = '', options = DEFAULT_SANITIZE_OPTIONS) =>
  sanitizeHtml(html, options).replace(/\{\{/g, '&#123;&#123;');

/**
 * Converts plain-text URLs in HTML into clickable anchor tags.
 *
 * Custom protocols from {@link LINKIFY_PROTOCOLS} are registered at module load time.
 *
 * @param {string} [html = ''] - HTML string to linkify.
 * @param {Object} [options = DEFAULT_LINKIFY_OPTIONS] - linkify-html options.
 * @returns {string} HTML with detected URLs wrapped in anchor tags.
 */
export const linkifyHtml = (html = '', options = DEFAULT_LINKIFY_OPTIONS) =>
  linkifyHtmlLib(html, options);

/**
 * Normalizes HTML by parsing and re-serializing it through the browser DOM.
 *
 * Useful for canonicalizing markup before registration or comparison (e.g. custom Vuetify icons).
 *
 * @param {string} [html = ''] - HTML string to normalize.
 * @returns {string} Browser-normalized HTML.
 */
export const normalizeHtml = (html = '') => {
  const element = document.createElement('div');

  element.innerHTML = html;

  return element.innerHTML;
};
