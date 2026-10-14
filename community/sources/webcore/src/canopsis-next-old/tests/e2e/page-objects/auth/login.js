// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@loginForm')
      .assert.visible('@usernameField');
  },

  clearUsername() {
    return this.customClearValue('@usernameField');
  },

  clearPassword() {
    return this.customClearValue('@passwordField');
  },

  verifyPageElementsAfter() {
    return this.waitForElementVisible('@userTopBarDropdownButton');
  },

  verifyErrorDisabledUser() {
    return this.waitForElementVisible('@errorLogin');
  },

  setUsername(username) {
    return this.customSetValue('@usernameField', username);
  },

  customSetPassword(password) {
    return this.customSetValue('@passwordField', password);
  },

  clickSubmitButton() {
    return this.customClick('@submitButton');
  },
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}login`;
  },
  elements: {
    loginForm: sel('loginForm'),
    usernameField: sel('username'),
    passwordField: sel('password'),
    submitButton: sel('submitButton'),
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    errorLogin: sel('errorLogin'),
  },
  commands: [commands],
};
