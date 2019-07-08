// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const sel = require('../../helpers/sel');

const usersPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyCreateUserModal() {
    return this.waitForElementVisible('@createUserModal');
  },

  verifyPageUserBefore(user) {
    return this.waitForElementVisible('@dataTable')
      .assert.visible(this.sel(user));
  },

  verifyCreateConfirmModal() {
    return this.waitForElementVisible('@createConfirmModal')
      .assert.visible('@confirmButton');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },

  clickEditButton(user) {
    return this.customClick(this.el('@editButton', this.sel(user)));
  },

  clickDeleteButton(user) {
    return this.customClick(this.el('@deleteButton', this.sel(user)));
  },

  clickConfirmButton() {
    return this.customClick('@confirmButton');
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

  clickSubmitButton() {
    return this.customClick('@submitButton');
  },

  el,
  sel,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/users`;
  },
  elements: {
    dataTable: '.v-datatable',
    addButton: sel('addButton'),
    editButton: `%s ${sel('editButton')}`,
    deleteButton: `%s ${sel('deleteButton')}`,
    confirmButton: sel('confirmButton'),
    createConfirmModal: sel('createConfirmModal'),
    createUserModal: sel('createUserModal'),
    usernameField: sel('username'),
    firstNameField: sel('firstName'),
    lastNameField: sel('lastName'),
    emailField: sel('email'),
    passwordField: sel('password'),
    roleField: `${sel('roleLayout')} .v-input__slot`,
    roleItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    languageField: `${sel('roleLanguage')} .v-input__slot`,
    languageItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    submitButton: sel('submitButton'),
  },
  commands: [usersPageCommands],
};
