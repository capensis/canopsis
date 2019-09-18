// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setName(value) {
    return this.customSetValue('@nameField', value);
  },

  setDescription(value) {
    return this.customSetValue('@descriptionField', value);
  },

  clickSelectDefaultViewButton() {
    return this.customClick('@selectDefaultViewButton');
  },

  clickRemoveDefaultViewButton() {
    return this.customClick('@removeDefaultViewButton');
  },
};

const modalSelector = sel('createRoleModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      nameField: sel('name'),
      descriptionField: sel('description'),
      selectDefaultViewButton: sel('selectDefaultViewButton'),
      removeDefaultViewButton: sel('removeDefaultViewButton'),
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
