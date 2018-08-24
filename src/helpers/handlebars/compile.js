import Handlebars from 'handlebars';

export default function compile(template, context) {
  const handleBarFunction = Handlebars.compile(template);

  return handleBarFunction(context);
}
