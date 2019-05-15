// For authoring Nightwatch tests, see
// http://nightwatchjs.org/guide#usage

const nightWatchRecord = require('nightwatch-record');

module.exports = {
  before(browser, done) {
    browser.login('root', 'root').perform(() => done());
  },
  beforeEach(browser, done) {
    nightWatchRecord.start(browser, done);
  },
  afterEach(browser, done) {
    nightWatchRecord.stop(browser, done);
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
