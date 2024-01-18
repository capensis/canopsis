const { serialize } = require('jest-snapshot/build/utils');

const isHtmlString = received => received && typeof received === 'string' && received[0] === '<';
const isVueWrapper = received => (
  received
  && typeof received === 'object'
  && typeof received.isVueInstance === 'function'
);

module.exports = {
  test(received) {
    return isHtmlString(received) || isVueWrapper(received);
  },
  serialize(received) {
    const html = (isVueWrapper(received) ? received.element?.outerHTML : received);

    if (!html) {
      return 'undefined';
    }

    const preparedHTML = html
      .replace(/ aria-owns="[-\w]+"/g, '')
      .replace(/ id="input-[-\d]+"/g, '')
      .replace(/ aria-labelledby="input-[\d]+"/g, '')
      .replace(/ id="list-[-\d]+"/g, '')
      .replace(/ id="list-item-[-\d]+"/g, '')
      .replace(/ data-v[-\w]+=""/g, '')
      .replace(/ for="input-[-\d]+"/g, '')
      .replace(/ v-text-field--is-booted/g, '')
      .replace(/ data-booted="[-\w]+"/g, '');

    const element = document.createElement(preparedHTML.startsWith('<tr') ? 'tbody' : 'body');

    element.innerHTML = preparedHTML;

    const serializedHtml = serialize(element.firstChild);

    return serializedHtml ? serializedHtml.replaceAll(/\n\s+\n/g, '\n') : serializedHtml;
  },
};
