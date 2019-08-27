// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.completed.logout()
      .end(done);
  },

  'Open current user modal': (browser) => {
    const topBar = browser.page.layout.topBar();
    const popup = browser.page.layout.popup();
    const createUserModal = browser.page.modals.admin.createUser();

    popup.clickOnEveryPopupsCloseIcons();

    topBar.clickUserDropdown()
      .clickUserProfileButton();

    createUserModal.verifyModalOpened();
  },

  'Select current user default view': (browser) => {
    const createUserModal = browser.page.modals.admin.createUser();
    const selectViewModal = browser.page.modals.view.selectView();

    createUserModal.clickSelectDefaultViewButton();

    selectViewModal.verifyModalOpened()
      .browseGroupById('05b2e049-b3c4-4c5b-94a5-6e7ff142b28c') // TODO: use from some constants file when we will use fixtures
      .browseViewById('875df4c2-027b-4549-8add-e20ed7ff7d4f')
      .verifyModalClosed();

    createUserModal.clickSubmitButton()
      .verifyModalClosed();
  },

  'Check default view': (browser) => {
    browser.url(process.env.VUE_DEV_SERVER_URL)
      .page.view()
      .verifyPageElementsBeforeById('875df4c2-027b-4549-8add-e20ed7ff7d4f'); // TODO: use from some constants file when we will use fixtures
  },
};
