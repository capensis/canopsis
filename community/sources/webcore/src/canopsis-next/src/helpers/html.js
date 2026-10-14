import DOMPurify from 'dompurify';
import linkifyHtmlLib from 'linkify-html';
import { registerCustomProtocol } from 'linkifyjs';
import { escape, unescape } from 'lodash';

import { DEFAULT_SANITIZE_OPTIONS, DEFAULT_LINKIFY_OPTIONS, LINKIFY_PROTOCOLS } from '@/config';

import { uid } from '@/helpers/uid';

/**
 * Encodes Vue template attributes so DOMPurify does not strip them during sanitization.
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

const COMMENT_REGEX = /<!--([\s\S]*?)-->/g;
const COMMENT_DIV_REGEX = /<div class="html-comment"[^>]*data-c="([^"]+)"[^>]*><\/div>/g;

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
 * Replaces HTML comments with placeholder elements so DOMPurify preserves them.
 *
 * Each comment is base64-encoded, escaped, and stored in a `div.html-comment` element.
 *
 * @param {string} [html = ''] - HTML string that may contain comments.
 * @returns {string} HTML with comments replaced by placeholder divs.
 */
export const replaceHtmlComments = (html = '') =>
  html.replace(COMMENT_REGEX, (_, content) => {
    const encoded = btoa(content);
    const safe = escape(encoded);
    const id = uid('cmt_');

    return `<div class="html-comment" data-id="${id}" data-c="${safe}"></div>`;
  });

/**
 * Restores HTML comments from placeholder elements created by {@link replaceHtmlComments}.
 *
 * @param {string} [html = ''] - HTML string with comment placeholder divs.
 * @returns {string} HTML with original comments restored.
 */
export const restoreHtmlComments = (html = '') =>
  html.replace(COMMENT_DIV_REGEX, (_, encodedEscaped) => {
    const encoded = unescape(encodedEscaped);
    const content = atob(encoded);

    return `<!--${content}-->`;
  });

/**
 * Sanitizes HTML while preserving Vue attributes and comments.
 *
 * Pipeline: encode Vue attributes → replace comments → DOMPurify → restore comments → decode Vue attributes.
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
    restoreHtmlComments(
      DOMPurify.sanitize(replaceHtmlComments(encodeVueAttributes(html)), options),
    ),
  );

  if (cacheable) {
    cacheSanitizedHtml(html, sanitizedHtml);
  }

  return sanitizedHtml;
};

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
