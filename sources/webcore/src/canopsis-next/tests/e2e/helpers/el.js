const util = require('util');

module.exports = function el(elementName, data) {
  const element = elementName[0] === '@' ? this.elements[elementName.slice(1)] : elementName;

  if (data) {
    return util.format(element.selector, data);
  }

  return element.selector;
};
