// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  clearRoleName() {
    return this.customClearValue('@roleName');
  },
  clearRoleDescription() {
    return this.customClearValue('@roleDescription');
  },
  setRoleName(value) {
    return this.customSetValue('@roleName', value);
  },
  setRoleDescription(value) {
    return this.customSetValue('@roleDescription', value);
  },
  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('createRoleModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      roleName: sel('roleName'),
      roleDescription: sel('roleDescription'),
      submitButton: sel('submitButton'),
    }),
  },
  commands: [commands],
});
