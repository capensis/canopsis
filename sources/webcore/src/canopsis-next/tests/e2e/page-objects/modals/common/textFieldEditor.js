// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setField(value) {
    return this.customClick('@textField')
      .sendKeys('@textField', value);
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('textFieldEditorModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      textField: `${sel('textField')} .jodit_wysiwyg`,
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
