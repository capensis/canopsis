// http://nightwatchjs.org/guide#usage

module.exports = {
  before(browser, done) {
    browser.maximizeWindow(done);
  },

  after(browser, done) {
    browser.end(done);
  },

  'Correct user credentials login': (browser) => {
    browser.finished.login('root', 'root'); // TODO: use from some constants file
  },

  'Authorized user logout': (browser) => {
    browser.finished.logout();
  },
};
