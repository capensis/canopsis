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
  const { submitButton, cancelButton } = pageObject.elements;
  const preparedPageObjectCommands = pageObject.commands && pageObject.commands.length ? pageObject.commands : [{}];

  return {
    ...pageObject,

    commands: [
      ...preparedPageObjectCommands.map((commandsItem) => {
        const commands = {
          verifyModalOpened() {
            return this.waitForElementVisible(selector);
          },

          verifyModalClosed() {
            return this.waitForElementNotPresent(selector);
          },

          ...commandsItem,
        };

        if (submitButton) {
          commands.clickSubmitButton = function clickSubmitButton() {
            return this.customClick(submitButton);
          };
        }

        if (cancelButton) {
          commands.clickCancelButton = function clickCancelButton() {
            return this.customClick(cancelButton);
          };
        }

        return commands;
      }),
    ],
  };
};
