// http://nightwatchjs.org/guide#usage

const { generateTemporaryView } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.defaultViewData = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.completed.logout()
      .end(done);
  },

  'Create test view': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.defaultViewData = {
        viewId: view._id,
        groupId: view.group_id,
      };
    });
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
    const { viewId, groupId } = browser.globals.defaultViewData;

    createUserModal.clickSelectDefaultViewButton();

    selectViewModal.verifyModalOpened()
      .browseGroupById(groupId)
      .browseViewById(viewId)
      .verifyModalClosed();

    createUserModal.clickSubmitButton()
      .verifyModalClosed();
  },

  'Check default view': (browser) => {
    const { viewId } = browser.globals.defaultViewData;

    browser.url(process.env.VUE_DEV_SERVER_URL)
      .page.view()
      .verifyPageElementsBeforeById(viewId);
  },

  'Delete test view': (browser) => {
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.completed.view.delete(groupId, viewId);
    browser.completed.view.deleteGroup(groupId);
  },
};
