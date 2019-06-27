import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

const prepareAttributes = attributes => Object.entries(attributes)
  .map(([key, value]) =>
    `${Handlebars.escapeExpression(key)}="${Handlebars.escapeExpression(value)}"`)
  .join(' ');

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

registerHelper('external-link', (options) => {
  const { href, text, ...attributes } = options.hash;
  const link = `<a href="${href}" ${prepareAttributes(attributes)}>${text}</a>`;
  return new Handlebars.SafeString(link);
});

registerHelper('internal-link', (options) => {
  const { href, text, ...attributes } = options.hash;
  const link = `<router-link to="${href}" ${prepareAttributes(attributes)}>${text}</router-link>`;
  return new Handlebars.SafeString(link);
});

export default {
  compile,
  registerHelper,
  unregisterHelper,
};
