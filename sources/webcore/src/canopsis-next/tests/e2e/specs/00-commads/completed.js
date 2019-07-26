// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  async after(browser, done) {
    await browser.completed.logout()
      .end(done);
  },

  'Test create user completed': (browser) => {
    browser.completed.createUser();
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Test delete user completed': (browser) => {
    browser.completed.deleteUser();
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Test create view completed': (browser) => {
    browser.completed.createView();
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Test delete view completed': (browser) => {
    browser.completed.deleteView();
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },
};
