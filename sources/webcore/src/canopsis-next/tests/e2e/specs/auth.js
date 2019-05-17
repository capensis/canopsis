// http://nightwatchjs.org/guide#usage

module.exports = {
  after(browser, done) {
    browser.end(done);
  },
  'Correct user credentials login': (browser) => {
    browser.login('root', 'root');
  },
  'Authorized user logout': (browser) => {
    browser.logout();
  },
};
