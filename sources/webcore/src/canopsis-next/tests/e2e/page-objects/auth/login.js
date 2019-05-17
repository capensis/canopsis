// https://nightwatchjs.org/guide/#working-with-page-objects

const loginPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@loginForm')
      .assert.visible('@usernameField');
  },

  verifyPageElementsAfter() {
    return this.waitForElementVisible('@userTopBarDropdownButton');
  },

  setUsername(username) {
    return this.customSetValue('@usernameField', username);
  },

  setPassword(password) {
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
  },
  commands: [loginPageCommands],
};
