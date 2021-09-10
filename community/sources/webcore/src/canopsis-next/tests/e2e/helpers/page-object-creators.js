const { isString } = require('lodash');

const PageUtils = require('nightwatch/lib/page-object/base-object');

function elementsWrapperCreator(selector, elements) {
  return Object.entries(elements).reduce((acc, [key, value]) => {
    if (isString(value) && !value.startsWith(selector)) {
      acc[key] = `${selector} ${value}`;
    } else {
      acc[key] = value;
    }

    return acc;
  }, {});
}

function modalCreator(selector, pageObject) {
  const { submitButton, cancelButton, deleteButton } = pageObject.elements;
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

        if (deleteButton) {
          commands.clickDeleteButton = function clickDeleteButton() {
            return this.customClick(deleteButton);
          };
        }

        return commands;
      }),
    ],
  };
}

function scopedPageObject(pageObject) {
  const preparedPageObjectCommands = pageObject.commands && pageObject.commands.length ? pageObject.commands : [{}];

  return {
    ...pageObject,
    commands: [
      ...preparedPageObjectCommands.map(commandsItem => ({
        setSelectorScope(sectionSelector) {
          const elements = sectionSelector
            ? elementsWrapperCreator(sectionSelector, pageObject.elements)
            : pageObject.elements;

          PageUtils.createElements(this, elements || {});

          return this;
        },

        ...commandsItem,
      })),
    ],
  };
}

module.exports.elementsWrapperCreator = elementsWrapperCreator;
module.exports.modalCreator = modalCreator;
module.exports.modalCreator = modalCreator;
module.exports.scopedPageObject = scopedPageObject;
