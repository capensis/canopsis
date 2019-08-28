// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clickField() {
    return this.customSetValue('@textareaField');
  },

  setField(value) {
    return this.customSetValue('@textareaField', value);
  },
};

const modalSelector = sel('textEditorModal');

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
