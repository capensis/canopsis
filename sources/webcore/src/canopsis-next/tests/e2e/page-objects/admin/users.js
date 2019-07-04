// https://nightwatchjs.org/guide/#working-with-page-objects

const usersPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyCreateUserModal() {
    return this.waitForElementVisible('@createUserModal');
  },

  verifyCreateConfirmModal() {
    return this.waitForElementVisible('@createConfirmModal');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },

  clickEditButton() {
    return this.customClick('@editButton');
  },

  clickDeleteButton() {
    return this.customClick('@deleteButton');
  },

  clickConfirmButton() {
    return this.customClick('@confirmButton');
  },

  clearUsername() {
    return this.customClearValue('@usernameField');
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

  selectRole() {
    return this.customClick('@roleField')
      .waitForElementVisible('@roleItemOption')
      .customClick('@roleItemOption');
  },

  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/users`;
  },
  elements: {
    dataTable: '.v-datatable',
    addButton: sel('addButton'),
    editButton: `.v-datatable tbody tr:last-child ${sel('editButton')}`,
    deleteButton: `.v-datatable tbody tr:last-child ${sel('deleteButton')}`,
    confirmButton: sel('confirmButton'),
    createConfirmModal: sel('createConfirmModal'),
    createUserModal: sel('createUserModal'),
    usernameField: sel('username'),
    firstNameField: sel('firstName'),
    lastNameField: sel('lastName'),
    emailField: sel('email'),
    passwordField: sel('password'),
    languageField: sel('language'),
    roleField: `${sel('roleLayout')} .v-input__slot`,
    roleItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(1) .v-list__tile',
    submitButton: sel('submitButton'),
  },
  commands: [usersPageCommands],
};
