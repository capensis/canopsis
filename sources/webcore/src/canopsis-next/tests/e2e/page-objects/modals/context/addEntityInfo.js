// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('addEntityInfoModal');

const commands = {
  clickEntityInfoName() {
    return this.customClick('@addEntityInfoName');
  },

  clearEntityInfoName() {
    return this.customClearValue('@addEntityInfoName');
  },

  setEntityInfoName(value) {
    return this.customSetValue('@addEntityInfoName', value);
  },

  clickEntityInfoDescription() {
    return this.customClick('@addEntityInfoDescription');
  },

  clearEntityInfoDescription() {
    return this.customClearValue('@addEntityInfoDescription');
  },

  setEntityInfoDescription(value) {
    return this.customSetValue('@addEntityInfoDescription', value);
  },

  clickEntityInfoValue() {
    return this.customClick('@addEntityInfoValue');
  },

  clearEntityInfoValue() {
    return this.customClearValue('@addEntityInfoValue');
  },

  setEntityInfoValue(value) {
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
