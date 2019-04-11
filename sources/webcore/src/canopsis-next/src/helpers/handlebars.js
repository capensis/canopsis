import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

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

export default {
  compile,
  registerHelper,
  unregisterHelper,
};
