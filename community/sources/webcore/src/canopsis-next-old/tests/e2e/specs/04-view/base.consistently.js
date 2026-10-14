// http://nightwatchjs.org/guide#usage

const uid = require('uid');
const { API_ROUTES } = require('../../../../src/config');
const { generateTemporaryView } = require('../../helpers/entities');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../constants');

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
      .clickAddTabButton();

    browser.page.modals.common.textFieldEditor()
      .verifyModalOpened()
      .setField(tab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.view._id}`,
      WAIT_FOR_FIRST_XHR_TIME,
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
      WAIT_FOR_FIRST_XHR_TIME,
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

  'Copy test tab into the same view': (browser) => {
    const viewPage = browser.page.view();
    const selectViewModal = browser.page.modals.view.selectView();
    const textFieldEditorModal = browser.page.modals.common.textFieldEditor();

    const { views } = browser.globals;
    const copyTab = `tab-${uid()}`;

    viewPage.clickCopyTab(views.view.tabId);

    selectViewModal.verifyModalOpened()
      .browseGroupById(views.view.group_id)
      .browseViewById(views.view._id);

    textFieldEditorModal.verifyModalOpened()
      .setField(copyTab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.view._id}`,
      WAIT_FOR_FIRST_XHR_TIME,
      () => textFieldEditorModal.clickSubmitButton(),
      ({ responseData, requestData }) => views.view = {
        copyTab,
        ...views.view,
        copyTabId: JSON.parse(requestData).tabs
          .filter(item => item.title === copyTab)[0]._id,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    textFieldEditorModal.verifyModalClosed();

    selectViewModal.verifyModalClosed();

    browser.perform(() => browser.assert.urlContains(`${views.view._id}?tabId=${views.view.copyTabId}`));
  },

  'Create test view for copying': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.views.viewForCopying = view;
    });
  },

  'Copy test tab into another view': (browser) => {
    const viewPage = browser.page.view();
    const selectViewModal = browser.page.modals.view.selectView();
    const textFieldEditorModal = browser.page.modals.common.textFieldEditor();

    const { views } = browser.globals;
    const copyTab = `tab-${uid()}`;

    viewPage.clickCopyTab(views.view.copyTabId);

    selectViewModal.verifyModalOpened()
      .browseGroupById(views.viewForCopying.group_id)
      .browseViewById(views.viewForCopying._id);

    textFieldEditorModal.verifyModalOpened()
      .setField(copyTab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${views.viewForCopying._id}`,
      WAIT_FOR_FIRST_XHR_TIME,
      () => textFieldEditorModal.clickSubmitButton(),
      ({ responseData, requestData }) => views.viewForCopying = {
        copyTab,
        ...views.viewForCopying,
        copyTabId: JSON.parse(requestData).tabs
          .filter(item => item.title === copyTab)[0]._id,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    textFieldEditorModal.verifyModalClosed();

    selectViewModal.verifyModalClosed();

    browser.perform(() => browser.assert.urlContains(`${views.viewForCopying._id}?tabId=${views.viewForCopying.copyTabId}`));
  },

  'Go back into test view': (browser) => {
    browser.page.layout.groupsSideBar()
      .clickLinkView(browser.globals.views.view._id);

    browser.page.view()
      .clickMenuViewButton()
      .clickEditViewButton();
  },

  'Move tab by dragdrop': (browser) => {
    const { copyTabId } = browser.globals.views.view;
    browser.page.view()
      .moveTab(copyTabId);
  },

  'Delete test tabs': (browser) => {
    const viewPage = browser.page.view();
    const confirmationModal = browser.page.modals.common.confirmation();

    const { tabId, copyTabId } = browser.globals.views.view;

    viewPage.clickDeleteTab(tabId);

    confirmationModal.verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    viewPage.clickDeleteTab(copyTabId);

    confirmationModal.verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Delete test view': (browser) => {
    const { view, viewForCopying } = browser.globals.views;

    browser.completed.view.delete(view.group_id, view._id);
    browser.completed.view.delete(viewForCopying.group_id, viewForCopying._id);

    browser.completed.view.deleteGroup(view.group_id);
    browser.completed.view.deleteGroup(viewForCopying.group_id);
  },
};
