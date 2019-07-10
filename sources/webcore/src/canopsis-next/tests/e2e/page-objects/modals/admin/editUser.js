// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

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

  clearUsername() {
    return this.customClearValue('@usernameField');
  },

  clearFirstName() {
    return this.customClearValue('@firstNameField');
  },

  clearLastName() {
    return this.customClearValue('@lastNameField');
  },

  clearEmail() {
    return this.customClearValue('@emailField');
  },

  selectRole(idx = 1) {
    return this.customClick('@roleField')
      .waitForElementVisible(this.el('@roleItemOption', idx))
      .customClick(this.el('@roleItemOption', idx));
  },

  selectLanguage(idx = 1) {
    return this.customClick('@languageField')
      .waitForElementVisible(this.el('@languageItemOption', idx))
      .customClick(this.el('@languageItemOption', idx));
  },

  clickSelectDefaultViewButton() {
    return this.customClick('@selectDefaultViewButton');
  },

  clickSubmitButton() {
    return this.customClick('@submitButton');
  },

  el,
};

const modalSelector = sel('createUserModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      usernameField: sel('username'),
      firstNameField: sel('firstName'),
      lastNameField: sel('lastName'),
      emailField: sel('email'),
      passwordField: sel('password'),
      roleField: `${sel('roleLayout')} .v-input__slot`,
      languageField: `${sel('roleLanguage')} .v-input__slot`,
      selectDefaultViewButton: sel('selectDefaultViewButton'),
      submitButton: sel('submitButton'),
    }),

    roleItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    languageItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
