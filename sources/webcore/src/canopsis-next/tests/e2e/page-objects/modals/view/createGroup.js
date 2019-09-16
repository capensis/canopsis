// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clearGroupName() {
    return this.customClearValue('@modalGroupNameField');
  },
  setGroupName(value) {
    return this.customSetValue('@modalGroupNameField', value);
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
  clickDeleteButton() {
    return this.customClick('@deleteButton');
  },
};

const modalSelector = sel('createGroupViewModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      modalGroupNameField: sel('modalGroupNameField'),
      submitButton: sel('createGroupSubmitButton'),
      deleteButton: sel('createGroupDeleteButton'),
    }),
  },
  commands: [commands],
});
