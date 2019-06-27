import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

const attributes = ({ hash }, locked = []) => Object.entries(hash)
  .filter(([key]) => !locked.includes(key))
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

const lockedAttrs = ['text', 'href'];

registerHelper('external-link', (options) => {
  const { href, text } = options.hash;
  const link = `<a href="${href}" ${attributes(options, lockedAttrs)}>${text}</a>`;
  return new Handlebars.SafeString(link);
});

registerHelper('internal-link', (options) => {
  const { href, text } = options.hash;
  const link = `<router-link to="${href}" ${attributes(options, lockedAttrs)}>${text}</router-link>`;
  return new Handlebars.SafeString(link);
});

export default {
  compile,
  registerHelper,
  unregisterHelper,
};
