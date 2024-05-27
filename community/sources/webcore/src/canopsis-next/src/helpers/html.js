import { registerCustomProtocol } from 'linkifyjs';
import linkifyHtmlLib from 'linkify-html';
import sanitizeHtmlLib from 'sanitize-html';

import { DEFAULT_SANITIZE_OPTIONS, DEFAULT_LINKIFY_OPTIONS, LINKIFY_PROTOCOLS } from '@/config';

/**
 * Register custom protocols for linkify
 */
LINKIFY_PROTOCOLS.forEach(registerCustomProtocol);

/**
 * Sanitize HTML document
 *
 * @param {string} [html = '']
 * @param {Object} [options = DEFAULT_SANITIZE_OPTIONS]
 * @return {string}
 */
export const sanitizeHtml = (html = '', options = DEFAULT_SANITIZE_OPTIONS) => sanitizeHtmlLib(html, options);

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
