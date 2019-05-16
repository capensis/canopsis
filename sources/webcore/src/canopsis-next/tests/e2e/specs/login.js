// For authoring Nightwatch tests, see
// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.login('root', 'root');

    done();
  },
  after(browser, done) {
    browser.end(done);
  },
  'default e2e tests': (browser) => {
    browser.url(process.env.VUE_DEV_SERVER_URL)
      .waitForElementVisible('.v-toolbar__content .v-btn__content')
      .assert.containsText('.v-toolbar__content .v-btn__content', 'menu');
  },
};
