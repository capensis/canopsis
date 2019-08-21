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

  'Create widget alarms with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget('AlarmsList')
      .verifyModalClosed();

    browser.completed.widget.setCommonField({
      row: 'row',
      sm: 13,
      md: 13,
      lg: 13,
      title: 'Alarms widget',
      periodRefresh: 120,
    });

    browser.page.widget.alarms()
      .clickAdvancedSettings()
      .clickDefaultSortColumn()
      .selectSortColumn(2)
      .selectSortOrder(2)
      .clickColumnNames()
      .clickColumnDown(1)
      .clearColumnLabel(2)
      .setColumnLabel(2, 'Connector')
      .clearColumnValue(2)
      .setColumnValue(2, 'alarm.v.connector')
      .clickColumnHtml(2)
      .clickColumnUp(2)
      .clickColumnClose(1)
      .clickColumnAdd()
      .setColumnLabel(8, 'Connector')
      .setColumnValue(8, 'alarm.v.connector')
      .clickDefaultNumberOfElementsPerPage()
      .selectElementsPerPage(3)
      .clickFilterOnOpenResolved()
      .clickOpenFilter()
      .clickResolvedFilter()
      .clickFilters()
      .clickAddFilter();

    browser.page.modals.common.createFilter()
      .verifyModalOpened()
      .setFilterTitle('FilterTitle1')
      .clickRadioOr()
      .clickAddRule()
      .selectFieldRule(2)
      .selectOperatorRule(2)
      .clickSubmitFilter()
      .verifyModalClosed();

    browser.page.widget.alarms()
      .clickAddFilter();

    browser.page.modals.common.createFilter()
      .verifyModalOpened()
      .setFilterTitle('FilterTitle2')
      .clickRadioOr()
      .clickAddRule()
      .selectFieldRule(1)
      .selectOperatorRule(1)
      .clickSubmitFilter()
      .verifyModalClosed();

    browser.page.widget.alarms()
      .clickMixFilters()
      .clickOrFilters()
      .selectFilters(1)
      .selectFilters(2)
      .clickInfoPopupButton();

    browser.page.modals.alarm.infoPopupSetting()
      .verifyModalOpened()
      .clickAddPopup();

    browser.page.modals.alarm.addInfoPopup()
      .verifyModalOpened()
      .selectSelectedColumn(2)
      .setTemplate('Template')
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.modals.alarm.infoPopupSetting()
      .clickEditPopup();

    browser.page.modals.alarm.addInfoPopup()
      .verifyModalOpened()
      .selectSelectedColumn(1)
      .setTemplate('End')
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.modals.alarm.infoPopupSetting()
      .clickDeletePopup()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.widget.alarms()
      .clickCreateEditMore();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .setRTE('More')
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.widget.alarms()
      .clickCreateEditMore();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .setRTE(' info...')
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.widget.alarms()
      .clickDeleteMore();

    browser.page.modals.confirmation()
      .verifyModalOpened()
      .clickConfirmButton()
      .verifyModalClosed();

    browser.page.widget.alarms()
      .clickEnableHtml()
      .clickAckGroup()
      .clickIsAckNoteRequired()
      .clickIsMultiAckEnabled()
      .clickFastAckOutput()
      .clickFastAckOutputSwitch()
      .setFastAckOutputText('test');

    browser.page.widget.alarms()
      .clickSubmitAlarms();
  },

  'Edit widget weather with some name': (browser) => {
    browser.page.view()
      .clickEditWidgetButton();

    browser.completed.widget.setCommonField({
      sm: 10,
      md: 10,
      lg: 10,
      title: 'Alarms widget(edited)',
      periodRefresh: 180,
    });

    browser.page.widget.alarms()
      .clickSubmitAlarms();
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
