// For authoring Nightwatch tests, see
// http://nightwatchjs.org/guide#usage

const ELEMENTS_WAITING_DELAY = 5000;

module.exports = {
  before(browser) {
    browser
      .url(`${process.env.VUE_DEV_SERVER_URL}login`)
      .waitForElementVisible('form', ELEMENTS_WAITING_DELAY)
      .setValue('input[name=username]', 'root')
      .setValue('input[type=password]', 'root')
      .click('button[type=submit]')
      .pause(ELEMENTS_WAITING_DELAY);
  },
  'default e2e tests': (browser) => {
    browser
      .url(process.env.VUE_DEV_SERVER_URL)
      .waitForElementVisible('.v-toolbar__content .v-btn__content', ELEMENTS_WAITING_DELAY)
      .assert.containsText('.v-toolbar__content .v-btn__content', 'menu')
      .end();
  },
};
