// http://nightwatchjs.org/guide#usage

const uid = require('uid');
const { API_ROUTES } = require('../../../../src/config');
const { generateTemporaryView } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.views = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  async after(browser, done) {
    browser.completed.logout()
      .end(done);

    delete browser.globals.views;
  },

  'Create test view': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.views.view = view;
    });
  },

  'Create test tab': (browser) => {
    const { views } = browser.globals;
    const tab = `tab-${uid()}`;

    browser.page.layout.groupsSideBar()
      .clickPanelHeader(views.view.group_id)
      .clickLinkView(views.view._id);

    browser.page.view()
      .clickMenuViewButton()
      .clickAddViewButton();

    browser.page.modals.common.textFieldEditor()
      .verifyModalOpened()
      .setField(tab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.view._id}`,
      5000,
      () => browser.page.modals.common.textFieldEditor()
        .clickSubmitButton(),
      ({ responseData, requestData }) => views.view = {
        tab,
        ...views.view,
        tabId: JSON.parse(requestData).tabs
          .filter(item => item.title === tab)[0]._id,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    browser.page.modals.common.textFieldEditor()
      .verifyModalClosed();
  },

  'Open test tab': (browser) => {
    const { views } = browser.globals;

    browser.page.view()
      .clickTab(views.view.tabId);
  },

  'Edit test tab': (browser) => {
    const { views } = browser.globals;
    const tab = `tab-${uid()}`;

    browser.page.view()
      .clickEditViewButton()
      .clickEditTab(views.view.tabId);

    browser.page.modals.common.textFieldEditor()
      .verifyModalOpened()
      .clearField()
      .setField(tab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.view._id}`,
      5000,
      () => browser.page.modals.common.textFieldEditor()
        .clickSubmitButton(),
      ({ responseData, requestData }) => views.view = {
        tab,
        ...views.view,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    browser.page.modals.common.textFieldEditor()
      .verifyModalClosed();
  },

  'Copy test tab': (browser) => {
    const { views } = browser.globals;
    const copyTab = `tab-${uid()}`;

    browser.page.view()
      .clickCopyTab(views.view.tabId);

    browser.page.modals.common.textFieldEditor()
      .verifyModalOpened()
      .setField(copyTab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.view._id}`,
      5000,
      () => browser.page.modals.common.textFieldEditor()
        .clickSubmitButton(),
      ({ responseData, requestData }) => views.view = {
        copyTab,
        ...views.view,
        copyTabId: JSON.parse(requestData).tabs
          .filter(item => item.title === copyTab)[0]._id,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    browser.page.modals.common.textFieldEditor()
      .verifyModalClosed();
  },

  'Move tab by dragdrop': (browser) => {
    const { copyTabId } = browser.globals.views.view;
    browser.page.view()
      .moveTab(copyTabId);
  },

  'Delete test tabs': (browser) => {
    const { tabId, copyTabId } = browser.globals.views.view;

    browser.page.view()
      .clickDeleteTab(tabId);
    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.view()
      .clickDeleteTab(copyTabId);
    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Delete test view': (browser) => {
    const { view } = browser.globals.views;

    browser.completed.view.delete(view.group_id, view._id);
    browser.completed.view.deleteGroup(view.group_id);
  },
};
