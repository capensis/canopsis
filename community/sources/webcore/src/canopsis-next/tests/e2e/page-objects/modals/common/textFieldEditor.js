// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clearField() {
    return this.customClearValue('@textField');
  },
  setField(value) {
    return this.customSetValue('@textField', value);
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('textFieldEditorModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      textField: sel('textField'),
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
