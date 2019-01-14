import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';

Handlebars.registerHelper('timestamp', (date) => {
  if (date) {
    return dateFilter(date, 'long');
  }

  return '';
});

export function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);

  return handleBarFunction(context);
}

export default {
  compile,
};
