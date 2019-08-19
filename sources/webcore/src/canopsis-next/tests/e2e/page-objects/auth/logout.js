// https://nightwatchjs.org/guide/#working-with-page-objects

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@userTopBarDropdownButton');
  },

  verifyPageElementsAfter() {
    return this.waitForElementVisible('@loginForm');
  },
};

module.exports = {
  elements: {
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    loginForm: sel('loginForm'),
  },
  commands: [commands],
};
