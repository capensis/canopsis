// https://nightwatchjs.org/guide/#working-with-page-objects

const WAIT_PAUSE = 500;

const logoutPageCommands = {
  clickUserNavigationTopBarButton() {
    this.waitForElementVisible('@userNavigationTopBarButton')
      .click('@userNavigationTopBarButton')
      .api.pause(WAIT_PAUSE);

    return this;
  },

  clickLogoutButton() {
    this.waitForElementVisible('@logoutButton')
      .click('@logoutButton')
      .api.pause(WAIT_PAUSE);

    return this;
  },
};

module.exports = {
  url() {
    return process.env.VUE_DEV_SERVER_URL;
  },
  elements: {
    userNavigationTopBarButton: '.v-toolbar__content .v-toolbar__items .v-menu.v-menu--inline:nth-child(3) .v-btn',
    logoutButton: '#logout-btn',
  },
  commands: [logoutPageCommands],
};
