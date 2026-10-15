// http://nightwatchjs.org/guide#usage

const { LANGUAGES_POSITIONS } = require('../../constants');
const {
  createAdminUser, removeUser, createWidgetView, removeWidgetView,
} = require('../../helpers/api');

module.exports = {
  async before(browser, done) {
    browser.globals.defaultViewData = await createWidgetView();
    const { data } = await createAdminUser();

    browser.globals.credentials = {
      password: data.password,
      username: data._id,
    };

    await browser.maximizeWindow()
      .completed.login(browser.globals.credentials.username, browser.globals.credentials.password);

    done();
  },

  async after(browser, done) {
    const { viewId, groupId } = browser.globals.defaultViewData;

    browser.completed.logout()
      .end();

    await removeUser(browser.globals.credentials.username);
    await removeWidgetView(viewId, groupId);

    delete browser.globals.credentials;
    delete browser.globals.defaultViewData;

    done();
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
      });
  },

  'Reset default view': (browser) => {
    const topBar = browser.page.layout.topBar();
    const createUserModal = browser.page.modals.admin.createUser();
    const selectViewModal = browser.page.modals.view.selectView();

    const { currentUser: { defaultview: defaultViewId } } = browser.globals;

    topBar.clickUserProfileButton();

    createUserModal.verifyModalOpened();

    if (defaultViewId) {
      createUserModal.clickSelectDefaultViewButton();

      selectViewModal.verifyModalOpened()
        .browseGroupByViewId(defaultViewId)
        .browseViewById(defaultViewId)
        .verifyModalClosed();
    } else {
      createUserModal.clickRemoveDefaultViewButton();
    }

    createUserModal.clickSubmitButton()
      .verifyModalClosed();
  },
};
