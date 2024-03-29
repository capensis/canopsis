import { registerCustomProtocol } from 'linkifyjs';
import linkifyHtmlLib from 'linkify-html';
import sanitizeHtmlLib from 'sanitize-html';

import { LINKIFY_PROTOCOLS } from '@/config';

const DEFAULT_SANITIZE_OPTIONS = {
  allowedTags: sanitizeHtmlLib.defaults.allowedTags.concat([
    'h1', 'h2', 'u', 'nl', 'font', 'img', 'video', 'audio', 'area', 'map', 'strike', 'button', 'span', 'address',
    'bdo', 'cite', 'q', 'dfn', 'var', 'dl', 'dt', 'dd', 'section', 'article', 'colgroup', 'col',

    /**
     * VUE COMPONENTS
     */
    'router-link', 'c-alarm-chip', 'c-alarm-tags-chips', 'c-copy-wrapper', 'c-links-list', 'service-entities-list',
  ]),
  allowedAttributes: {
    '*': [
      'style', 'title', 'class', 'id', 'v-if', 'autoplay', 'colspan', 'controls', 'dir', 'align', 'width', 'height',
      'role',
    ],
    a: ['href', 'name', 'target'],
    img: ['src', 'alt'],
    font: ['color', 'size', 'face'],
    'router-link': ['href', 'name', 'target', 'to'],
    'c-alarm-chip': ['value'],
    'c-alarm-tags-chips': [':alarm', ':selected-tag', 'closable-active', 'inline-count', '@select', '@close'],
    'c-copy-wrapper': ['value'],
    'c-links-list': [':links', ':category'],
    'service-entities-list': [
      ':service', ':service-entities', ':widget-parameters', ':pagination', ':total-items', ':actions-requests',
      'entity-name-field', '@refresh', '@update:pagination', '@add:action',
    ],
  },
  allowedSchemes: sanitizeHtmlLib.defaults.allowedSchemes.concat(['data']),
};

const DEFAULT_LINKIFY_OPTIONS = {
  target: '_blank',
  ignoreTags: ['script', 'style'],
  validate: (str, type, token) => token.hasProtocol(),
};

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
