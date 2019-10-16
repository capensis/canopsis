// http://nightwatchjs.org/guide#usage

module.exports = {
  before(browser, done) {
    browser.maximizeWindow(done);
  },

  after(browser, done) {
    browser.end(done);
  },

  'Correct user credentials login': (browser) => {
    browser.completed.loginAsAdmin();
  },

  'Authorized user logout': (browser) => {
    browser.completed.logout();
  },
};
