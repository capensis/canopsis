// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('addEntityInfoModal');

const commands = {
  clickName() {
    return this.customClick('@addEntityInfoName');
  },

  clearName() {
    return this.customClearValue('@addEntityInfoName');
  },

  setName(value) {
    return this.customSetValue('@addEntityInfoName', value);
  },

  clickDescription() {
    return this.customClick('@addEntityInfoDescription');
  },

  clearDescription() {
    return this.customClearValue('@addEntityInfoDescription');
  },

  setDescription(value) {
    return this.customSetValue('@addEntityInfoDescription', value);
  },

  clickValue() {
    return this.customClick('@addEntityInfoValue');
  },

  clearValue() {
    return this.customClearValue('@addEntityInfoValue');
  },

  setValue(value) {
    return this.customSetValue('@addEntityInfoValue', value);
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('addEntityInfoCancelButton'),
      submitButton: sel('addEntityInfoSubmitButton'),
    }),

    addEntityInfoName: sel('addEntityInfoName'),
    addEntityInfoDescription: sel('addEntityInfoDescription'),
    addEntityInfoValue: sel('addEntityInfoValue'),
  },
  commands: [commands],
});
