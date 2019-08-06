// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clearRightID() {
    return this.customClearValue('@rightID');
  },
  clearRightDescription() {
    return this.customClearValue('@rightDescription');
  },
  setRightID(value) {
    return this.customSetValue('@rightID', value);
  },
  setRightDescription(value) {
    return this.customSetValue('@rightDescription', value);
  },
  selectRightType(index = 1) {
    return this.customClick('@rightTypeField')
      .waitForElementVisible(this.el('@rightTypeOption', index))
      .customClick(this.el('@rightTypeOption', index));
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('createRightModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      rightID: sel('rightID'),
      rightDescription: sel('rightDescription'),
      rightTypeField: `${sel('typeLayout')} .v-input__slot`,
      submitButton: sel('submitButton'),
    }),

    rightTypeOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
