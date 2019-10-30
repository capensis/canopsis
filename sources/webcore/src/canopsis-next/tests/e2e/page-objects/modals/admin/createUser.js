// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const commands = {
  verifyModalOpened() {
    return this.waitForElementVisible('@modalSelector')
      .waitForElementNotVisible('@progressOverlay');
  },

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

  selectRole(index = 1) {
    return this.customClick('@roleField')
      .waitForElementVisible(this.el('@roleItemOption', index))
      .customClick(this.el('@roleItemOption', index));
  },

  selectLanguage(index = 1) {
    return this.customClick('@languageField')
      .waitForElementVisible(this.el('@languageItemOption', index))
      .customClick(this.el('@languageItemOption', index));
  },

  selectNavigationType(index = 1) {
    return this.customClick('@navigationTypeField')
      .waitForElementVisible(this.el('@navigationTypeOption', index))
      .customClick(this.el('@navigationTypeOption', index));
  },

  clickSelectDefaultViewButton() {
    return this.customClick('@selectDefaultViewButton');
  },

  clickRemoveDefaultViewButton() {
    return this.customClick('@removeDefaultViewButton');
  },

  clickEnabled() {
    return this.customClick('@enabled');
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
      languageField: `${sel('languageLayout')} .v-input__slot`,
      navigationTypeField: `${sel('navigationTypeLayout')} .v-input__slot`,
      selectDefaultViewButton: sel('selectDefaultViewButton'),
      removeDefaultViewButton: sel('removeDefaultViewButton'),
      enabled: `.v-input${sel('enabled')} .v-input--selection-controls__ripple`,
      submitButton: sel('submitButton'),
      progressOverlay: sel('progressOverlay'),
    }),

    modalSelector,
    roleItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    languageItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    navigationTypeOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
