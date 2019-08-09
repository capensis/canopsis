// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('textEditorModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
