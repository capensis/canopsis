import { isString } from 'lodash';

/**
 * Normalize handlebars NBSP
 *
 * Replace all NBSP in handlebars template to spaces.
 * Non-string values (Handlebars AST) are returned unchanged.
 *
 * @param {string | object} [html = '']
 * @returns {string | object}
 */
export const normalizeHandlebarsNbsp = (html = '') => (
  isString(html)
    ? html.replace(/{{{[\s\S]*?}}}|{{[\s\S]*?}}/g, block => block.replace(/&nbsp;|\u00A0/g, ' '))
    : html
);
