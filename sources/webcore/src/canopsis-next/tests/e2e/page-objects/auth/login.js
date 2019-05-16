// https://nightwatchjs.org/guide/#working-with-page-objects

const WAIT_PAUSE = 500;

const loginPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@loginForm')
      .assert.visible('@usernameField');
  },

  verifyPageElementsAfter() {
    return this.waitForElementVisible('@navigationTopBarButton')
      .assert.containsText('@navigationTopBarButton', 'menu');
  },

  enterUsername(username) {
    this.waitForElementVisible('@usernameField')
      .click('@usernameField')
      .setValue('@usernameField', username)
      .api.pause(WAIT_PAUSE);

    return this;
  },

  enterPassword(password) {
    this.waitForElementVisible('@passwordField')
      .click('@passwordField')
      .setValue('@passwordField', password)
      .api.pause(WAIT_PAUSE);

    return this;
  },

  clickSubmitButton() {
    this.waitForElementVisible('@submitButton')
      .click('@submitButton')
      .api.pause(WAIT_PAUSE);

    return this;
  },
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}login`;
  },
  elements: {
    loginForm: '.loginContainer form',
    navigationTopBarButton: '.v-toolbar__content .v-btn__content',
    usernameField: 'input[name=username]',
    passwordField: 'input[name=password]',
    submitButton: 'button[type=submit]',
  },
  commands: [loginPageCommands],
};
