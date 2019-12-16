// http://nightwatchjs.org/guide#usage

const { WIDGET_TYPES } = require('@/constants');

const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const ALARMS_COUNT = 134;
const SEARCH_STRING = 'feeder2_inst3';
const SEARCH_RESULT_COUNT = 16;
const REQUEST_COUNT_WITH_EMPTY_SEARCH = 0;

module.exports = {
  async before(browser, done) {
    browser.globals.defaultViewData = await createWidgetView();
    const { groupId, viewId } = browser.globals.defaultViewData;

    const widgetInfo = {
      type: WIDGET_TYPES.alarmList,
      row: {
        title: 'Row',
      },
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
    };

    const { widgetId } = await createWidgetForView(browser.globals.defaultViewData.viewId, widgetInfo);

    browser.globals.defaultViewData.widgetId = widgetId;

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .clickPanelHeader(groupId)
      .clickLinkView(viewId);

    done();
  },

  async after(browser, done) {
    const { viewId, groupId } = browser.globals.defaultViewData;

    browser.completed.logout()
      .end(done);

    await removeWidgetView(viewId, groupId);

    delete browser.globals.defaultViewData;
    delete browser.globals.temporary;

    done();
  },

  'The empty search shows no results': (browser) => {
    const commonTable = browser.page.tables.common();
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitAllAlarmsListXHR(
      () => commonTable.keyupSearchEnter(),
      (xhrs) => {
        browser.assert.equal(xhrs.length, REQUEST_COUNT_WITH_EMPTY_SEARCH);
      },
    );
  },

  'The search with magnifier button displays relevant results': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable
        .setSearchInput(SEARCH_STRING)
        .clickSearchButton(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.total, SEARCH_RESULT_COUNT);

        commonTable
          .clearSearchInput()
          .clickSearchButton();
      },
    );
  },

  'The search with button "Enter" displays relevant results': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable
        .clickSearchInput()
        .clearSearchInput()
        .setSearchInput(SEARCH_STRING)
        .keyupSearchEnter(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.total, SEARCH_RESULT_COUNT);
      },
    );
  },

  'The button with cross cancels current search': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickSearchResetButton(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.notEqual(response.total, SEARCH_RESULT_COUNT);
        browser.assert.equal(response.total, ALARMS_COUNT);
      },
    );
  },

  'The click on the button with question mark shows pop-up with additional information': (browser) => {
    browser.page.tables.common()
      .moveToSearchInformation()
      .verifySearchInformationVisible();
  },

  'Removing a cursor from pop-up with additional information makes it disappear': (browser) => {
    browser.page.tables.common()
      .moveOutsideSearchInformation()
      .verifySearchInformationHidden();
  },
};
