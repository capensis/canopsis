import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

const prepareAttributes = attributes => Object.entries(attributes)
  .map(([key, value]) =>
    `${Handlebars.escapeExpression(key)}="${Handlebars.escapeExpression(value)}"`)
  .join(' ');

/**
 *
 * @param template
 * @param context
 * @returns {*}
 */
export function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);

  return handleBarFunction(context);
}

export function registerHelper(name, helper) {
  return Handlebars.registerHelper(name, helper);
}

export function unregisterHelper(name) {
  return Handlebars.unregisterHelper(name);
}

registerHelper('timestamp', (date) => {
  if (date) {
    return dateFilter(date, 'long');
  }

  return '';
});

registerHelper('internal-link', (options) => {
  const { href, text, ...attributes } = options.hash;
  const path = href.replace(window.location.origin, '');

  const link = `<router-link to="${path}" ${prepareAttributes(attributes)}>${text}</router-link>`;
  return new Handlebars.SafeString(link);
});
