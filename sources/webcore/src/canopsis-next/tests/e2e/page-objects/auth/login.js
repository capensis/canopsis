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

  verifyPageElementsAfter(checkLogin = true) {
    return this.waitForElementVisible('@userTopBarDropdownButton', 5000, checkLogin, (result) => {
      if (!checkLogin) {
        if (result.value) {
          this.end();
        }
      }
      return this;
    });
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
  commands: [commands],
};
