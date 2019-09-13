// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  selectSelectedColumn(index = 1) {
    return this.customClick('@selectedColumn')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  setTemplate(value) {
    return this.customClick('@template')
      .sendKeys('@template', value);addInfoPopupLayout
  },

  el,
};

const modalSelector = sel('addInfoPopup');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      selectedColumn: `${sel('addInfoPopupLayout')} .v-input__slot`,
      template: `${sel('addInfoPopupLayout')} .jodit_wysiwyg`,

      submitButton: sel('addInfoSubmitButton'),
      cancelButton: sel('addInfoCancelButton'),
    }),
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
