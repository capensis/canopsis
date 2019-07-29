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
    browser.completed.createUser(undefined, user => browser.globals.user = user);
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Test delete user completed': (browser) => {
    browser.completed.deleteUser(browser.globals.user._id);
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Test create view completed': (browser) => {
    browser.completed.createView(undefined, view => browser.globals.view = view);
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Test delete view completed': (browser) => {
    browser.completed.deleteView(browser.globals.view.group_id, browser.globals.view._id);
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },
};
