// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.login('root', 'root'); // TODO: use from some constants file

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Browse exploitation event-filter': (browser) => {
    browser.page.layout.topBar()
      .clickUserProfileButton()
      .api.pause(5000);
  },

  // 'Edit user with some username': (browser) => {},
  //
  // 'Remove user with some username': (browser) => {},
};
