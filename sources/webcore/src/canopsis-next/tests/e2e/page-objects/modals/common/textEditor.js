// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setRTE(value) {
    return this.customClick('@textRTE')
      .sendKeys('@textRTE', value);
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('textEditorModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      textRTE: `${sel('jodit')} .jodit_wysiwyg`,
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
