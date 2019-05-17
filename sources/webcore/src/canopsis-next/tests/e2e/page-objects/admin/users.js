// https://nightwatchjs.org/guide/#working-with-page-objects

const usersPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyCreateUserModal() {
    return this.waitForElementVisible('@createUserModal');
  },

  clickAddButton() {
    return this.customClick('@addButton');
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
    return this.customSetValue('@usernameField', value);
  },

  selectRole() {
    return this.customClick('@roleField')
      .waitForElementVisible('@roleItemOption')
      .customClick('@roleItemOption');
  },
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/users`;
  },
  elements: {
    dataTable: '.v-datatable',
    addButton: sel('addButton'),
    createUserModal: sel('createUserModal'),
    usernameField: sel('username'),
    firstNameField: sel('firstName'),
    lastNameField: sel('lastName'),
    emailField: sel('email'),
    passwordField: sel('password'),
    roleField: sel('role'),
    languageField: sel('language'),
    roleItemOption: '.v-select-list v-list__tile',
  },
  commands: [usersPageCommands],
};
