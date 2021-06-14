const util = require('util');

module.exports = function el(elementName, ...args) {
  const element = elementName[0] === '@' ? this.elements[elementName.slice(1)] : elementName;

  if (args.length) {
    return util.format(element.selector, ...args);
  }

  return element.selector;
};
