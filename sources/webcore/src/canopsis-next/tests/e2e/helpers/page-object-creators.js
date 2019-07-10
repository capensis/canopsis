const { isString } = require('lodash');

module.exports.elementsWrapperCreator = function elementsWrapperCreator(selector, elements) {
  return Object.entries(elements).reduce((acc, [key, value]) => {
    if (isString(value) && !value.startsWith(selector)) {
      acc[key] = `${selector} ${value}`;
    } else {
      acc[key] = value;
    }

    return acc;
  }, {});
};

module.exports.modalCreator = function modalCreator(selector, pageObject) {
  return {
    ...pageObject,

    commands: [
      ...pageObject.commands.map(commandsItem => ({
        verifyModalOpened() {
          return this.waitForElementVisible(selector);
        },

        verifyModalClosed() {
          return this.waitForElementNotPresent(selector);
        },

        ...commandsItem,
      })),
    ],
  };
};
