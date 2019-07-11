// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setUsername(value) {
    return this.customSetValue('@usernameField', value);
  },

  setFirstName(value) {
    return this.customSetValue('@firstNameField', value);
  },

  setLastName(value) {
    return this.customSetValue('@lastNameField', value);
  },

  setEmail(value) {
    return this.customSetValue('@emailField', value);
  },

  setPassword(value) {
    return this.customSetValue('@passwordField', value);
  },

  selectRole() {
    return this.customClick('@roleField')
      .waitForElementVisible('@roleItemOption')
      .customClick('@roleItemOption');
  },

  clickSelectDefaultViewButton() {
    return this.customClick('@selectDefaultViewButton');
  },

  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

const modalSelector = sel('createRoleModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      usernameField: sel('username'),
      firstNameField: sel('firstName'),
      lastNameField: sel('lastName'),
      emailField: sel('email'),
      passwordField: sel('password'),
      languageField: sel('language'),
      roleField: `${sel('roleLayout')} .v-input__slot`,
      selectDefaultViewButton: sel('selectDefaultViewButton'),
      submitButton: sel('submitButton'),
    }),

    roleItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(1) .v-list__tile',
  },
  commands: [commands],
});
