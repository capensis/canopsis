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

  'Open current user modal': (browser) => {
    const topBar = browser.page.layout.topBar();
    const createUserModal = browser.page.modals.admin.createUser();

    topBar.clickUserDropdown()
      .clickUserProfileButton();

    createUserModal.verifyModalOpened()
      .selectNavigationType(2)
      .clickSubmitButton()
      .verifyModalClosed();
  },
};
