// http://nightwatchjs.org/guide#usage

const uid = require('uid');
const { API_ROUTES } = require('../../../../src/config');
const { generateTemporaryView } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    const { view } = browser.globals.temporary;

    browser.completed.view.delete(view.group_id, view._id);

    browser.completed.logout()
      .end(done);

    delete browser.globals.temporary;
  },

  'Create test view': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.temporary.view = view;
    });
  },

  'Create test tab': (browser) => {
    const { temporary } = browser.globals;
    const tab = `tab-${uid()}`;

    browser.page.layout.groupsSideBar()
      .clickPanelHeader(temporary.view.group_id)
      .clickLinkView(temporary.view._id);

    browser.page.view()
      .clickMenuViewButton()
      .clickAddViewButton();

    browser.page.modals.common.textFieldEditor()
      .verifyModalOpened()
      .setField(tab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${temporary.view._id}`,
      5000,
      () => browser.page.modals.common.textFieldEditor()
        .clickSubmitButton(),
      ({ responseData, requestData }) => temporary.view = {
        tab,
        ...temporary.view,
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
    const { temporary } = browser.globals;

    browser.page.view()
      .clickTab(temporary.view.tabId);
  },

  'Create widget weather with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget('ServiceWeather')
      .verifyModalClosed();

    browser.completed.widget.setCommonField({
      row: 'row',
      sm: 13,
      md: 13,
      lg: 13,
      title: 'Weather widget',
      periodRefresh: 120,
    });

    browser.page.widget.weather()
      .clickSubmitWeather();
  },

  'Edit widget weather with some name': (browser) => {
    browser.page.view()
      .clickEditWidgetButton();

    browser.completed.widget.setCommonField({
      sm: 10,
      md: 10,
      lg: 10,
      title: 'Weather widget(edited)',
      periodRefresh: 180,
    });

    browser.page.widget.weather()
      .clickSubmitWeather();
  },

  'Delete widget weather with some name': (browser) => {
    browser.page.view()
      .clickDeleteWidgetButton();

    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();
  },

  'Delete row with some name': (browser) => {
    browser.page.view()
      .clickDeleteRowButton();

    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();
  },
};
