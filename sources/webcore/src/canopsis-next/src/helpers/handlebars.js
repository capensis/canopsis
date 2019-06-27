import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

const attributes = ({ hash }) => {
  const attrs = [];
  Object.keys(hash).forEach((prop) => {
    attrs.push(`${
      Handlebars.escapeExpression(prop)
    }="${
      Handlebars.escapeExpression(hash[prop])
    }"`);
  });
  return attrs.join(' ');
};

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

registerHelper('external-link', (text, href, options) => new Handlebars.SafeString(`<a href="${href}" ${attributes(options)}>${text}</a>`));

registerHelper('internal-link', (text, href, options) => new Handlebars.SafeString(`<router-link to="${href}" ${attributes(options)}>${text}</router-link>`));

export default {
  compile,
  registerHelper,
  unregisterHelper,
};
