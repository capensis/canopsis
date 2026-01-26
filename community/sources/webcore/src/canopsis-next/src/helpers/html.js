import { escape, unescape } from 'lodash';
import { registerCustomProtocol } from 'linkifyjs';
import linkifyHtmlLib from 'linkify-html';
import sanitizeHtmlLib from 'sanitize-html';

import { DEFAULT_SANITIZE_OPTIONS, DEFAULT_LINKIFY_OPTIONS, LINKIFY_PROTOCOLS } from '@/config';

import { uid } from '@/helpers/uid';

/**
 * Register custom protocols for linkify
 */
LINKIFY_PROTOCOLS.forEach(registerCustomProtocol);

const COMMENT_REGEX = /<!--([\s\S]*?)-->/g;
const COMMENT_DIV_REGEX = /<div class="html-comment"[^>]*data-c="([^"]+)"[^>]*><\/div>/g;

/**
 * Replace HTML comments with encoded div elements
 *
 * This function finds all HTML comments (<!-- -->) in the provided HTML string and replaces them
 * with div elements containing the encoded comment content. The encoding process:
 * 1. Extracts the comment content
 * 2. Encodes it to base64 using btoa()
 * 3. Escapes the base64 string using lodash escape() to make it safe for HTML attributes
 * 4. Generates a unique ID for each comment
 * 5. Returns a div element with class "html-comment" containing the encoded data
 *
 * This is useful for preserving HTML comments during sanitization, as most HTML sanitizers
 * remove comments. The comments can later be restored using restoreHtmlComments().
 *
 * @param {string} [html = ''] - HTML string that may contain comments
 * @return {string} - HTML string with comments replaced by encoded div elements
 *
 * @example
 * replaceHtmlComments('<!-- This is a comment --><p>Text</p>')
 * // Returns: '<div class="html-comment" data-id="cmt_123" data-c="VGhpcyBpcyBhIGNvbW1lbnQ="></div><p>Text</p>'
 */
export const replaceHtmlComments = (html = '') => html.replace(COMMENT_REGEX, (_, content) => {
  const encoded = btoa(content);
  const safe = escape(encoded);
  const id = uid('cmt_');

  return `<div class="html-comment" data-id="${id}" data-c="${safe}"></div>`;
});

/**
 * Restore HTML comments from encoded div elements
 *
 * This function finds all div elements with class "html-comment" that were created by
 * replaceHtmlComments() and restores them back to their original HTML comment format.
 * The restoration process:
 * 1. Extracts the escaped base64 string from the data-c attribute
 * 2. Unescapes it using lodash unescape()
 * 3. Decodes the base64 string using atob()
 * 4. Returns the original HTML comment format (<!-- -->)
 *
 * This function is the inverse of replaceHtmlComments() and should be called after
 * sanitization to restore the original comments.
 *
 * @param {string} [html = ''] - HTML string containing encoded comment div elements
 * @return {string} - HTML string with encoded divs restored to original comments
 *
 * @example
 * restoreHtmlComments(
 *   '<div class="html-comment" data-id="cmt_123" data-c="VGhpcyBpcyBhIGNvbW1lbnQ="></div><p>Text</p>'
 * )
 * // Returns: '<!-- This is a comment --><p>Text</p>'
 */
export const restoreHtmlComments = (html = '') => html.replace(COMMENT_DIV_REGEX, (_, encodedEscaped) => {
  const encoded = unescape(encodedEscaped);
  const content = atob(encoded);

  return `<!--${content}-->`;
});

/**
 * Sanitize HTML document
 *
 * @param {string} [html = '']
 * @param {Object} [options = DEFAULT_SANITIZE_OPTIONS]
 * @return {string}
 */
export const sanitizeHtml = (html = '', options = DEFAULT_SANITIZE_OPTIONS) => (
  restoreHtmlComments(sanitizeHtmlLib(replaceHtmlComments(html), options))
);

/**
 * Convert all links in html to tag <a>
 *
 * @param {string} [html = '']
 * @param {Object} [options = DEFAULT_LINKIFY_OPTIONS]
 * @return {string}
 */
export const linkifyHtml = (html = '', options = DEFAULT_LINKIFY_OPTIONS) => linkifyHtmlLib(html, options);

/**
 * Normilize HTML (close not closed tags and etc.)
 *
 * @param {string} [html = '']
 * @return {string}
 */
export const normalizeHtml = (html = '') => {
  const element = document.createElement('div');

  element.innerHTML = html;

  return element.innerHTML;
};
