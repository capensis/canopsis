/**
 * Normalize handlebars NBSP
 *
 * Replace all NBSP in handlebars template to spaces.
 * This is needed because handlebars does not support NBSP.
 *
 * @param {string} [html = '']
 * @returns {string}
 */
export const normalizeHandlebarsNbsp = (html = '') => html.replace(/{{{[\s\S]*?}}}|{{[\s\S]*?}}/g, block => block.replace(/&nbsp;|\u00A0/g, ' '));
