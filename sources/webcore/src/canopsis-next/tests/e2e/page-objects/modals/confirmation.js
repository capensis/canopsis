// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../helpers/page-object-creators');

const commands = {
  clickConfirmButton() {
    return this.customClick('@confirmButton');
  },
};

const modalSelector = sel('confirmationModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      confirmButton: sel('confirmButton'),
    }),
  },
  commands: [commands],
});
