// https://nightwatchjs.org/guide/#working-with-page-objects

const WAIT_PAUSE = 500;

const logoutPageCommands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@userTopBarDropdownButton');
  },

  verifyPageElementsAfter() {
    return this.waitForElementVisible('@loginForm');
  },

  clickUserNavigationTopBarButton() {
    this.waitForElementVisible('@userTopBarDropdownButton')
      .click('@userTopBarDropdownButton')
      .api.pause(WAIT_PAUSE);

    return this;
  },

  clickLogoutButton() {
    this.waitForElementNotPresent('@activePopup', 15000)
      .waitForElementVisible('@logoutButton')
      .click('@logoutButton')
      .api.pause(WAIT_PAUSE);

    return this;
  },
};

module.exports = {
  elements: {
    userTopBarDropdownButton: sel('userTopBarDropdownButton'),
    logoutButton: sel('logoutButton'),
    loginForm: sel('loginForm'),
    activePopup: `${sel('popupsWrapper')} .v-alert`,
  },
  commands: [logoutPageCommands],
};
