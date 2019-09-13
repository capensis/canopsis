// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('textEditorModal');

const commands = {
  clickField() {
    return this.customClick('@textareaField');
  },

  setField(value) {
    return this.sendKeys('@textareaField', value);
  },

  clearField() {
    return this.customClearRTE('@textareaField');
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('textEditorSubmitButton'),
      cancelButton: sel('textEditorCancelButton'),
    }),
    textareaField: `${sel('textEditorModal')} .jodit_wysiwyg`,
  },
  commands: [commands],
});
