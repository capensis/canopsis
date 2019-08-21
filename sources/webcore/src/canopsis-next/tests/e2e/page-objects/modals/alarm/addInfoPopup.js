// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },

  selectSelectedColumn(index = 1) {
    return this.customClick('@selectedColumn')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  setTemplate(value) {
    return this.customClick('@template')
      .sendKeys('@template', value);
  },
};

const modalSelector = sel('addInfoPopup');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      selectedColumn: `${sel('addInfoPopupFields')} .v-input__slot`,
      template: `${sel('addInfoPopupFields')} .jodit_wysiwyg`,
      submitButton: sel('submitButton'),
    }),
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
