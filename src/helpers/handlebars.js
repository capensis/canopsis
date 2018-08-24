import Handlebars from 'handlebars';
import i18n from '@/i18n';

Handlebars.registerHelper('timestamp', (date) => {
  if (date) {
    return i18n.d(new Date(date * 1000), 'long');
  }

  return '';
});

export default function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);

  return handleBarFunction(context);
}
