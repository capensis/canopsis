// http://nightwatchjs.org/guide#usage

const { LANGUAGES_POSITIONS } = require('../../constants');
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
    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened();
  },

  'Select current user default view and interface language': (browser) => {
    const createUserModal = browser.page.modals.admin.createUser();
    const selectViewModal = browser.page.modals.view.selectView();
    const { viewId, groupId } = browser.globals.defaultViewData;

    createUserModal.selectLanguage(LANGUAGES_POSITIONS.fr)
      .clickSelectDefaultViewButton();

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

  'Caching of user first interface language': (browser) => {
    browser.page.layout.topBar()
      .clickUserDropdown()
      .getLogoutButtonText(({ value }) => {
        browser.globals.logoutButtonText = value;
      });
  },

  'Change user interface language': (browser) => {
    browser.page.layout.topBar()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectLanguage(LANGUAGES_POSITIONS.en)
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();
  },

  'Check user second interface language': (browser) => {
    const { logoutButtonText } = browser.globals;

    browser.page.layout.topBar()
      .clickUserDropdown()
      .getLogoutButtonText(({ value }) => {
        browser.assert.notEqual(value, logoutButtonText);
      })
      .clickUserDropdown();
  },

  'Reset default view': (browser) => {
    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .clickRemoveDefaultViewButton()
      .verifyModalClosed();
  },

  'Delete test view': (browser) => {
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.completed.view.delete(groupId, viewId);
    browser.completed.view.deleteGroup(groupId);
  },
};
