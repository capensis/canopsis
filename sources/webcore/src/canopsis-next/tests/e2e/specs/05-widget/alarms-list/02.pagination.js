// http://nightwatchjs.org/guide#usage

const { WIDGET_TYPES } = require('@/constants');
const { PAGINATION_LIMIT } = require('@/config');

const { PAGINATION_PER_PAGE_VALUES } = require('../../../constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const BOTTOM_PAGINATION_PAGE_NUMBER = 12;

module.exports = {
  async before(browser, done) {
    browser.globals.tablePageNumber = 1;
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

  'Right arrow opens the next page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber += 1;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickNextPageTopPagination(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, PAGINATION_LIMIT * browser.globals.tablePageNumber);

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(page, browser.globals.tablePageNumber);
        });
      },
    );
  },

  'Left arrow opens the previous page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber -= 1;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickPreviousPageTopPagination(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, PAGINATION_LIMIT * browser.globals.tablePageNumber);

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(page, browser.globals.tablePageNumber);
        });
      },
    );
  },

  'Right arrow at the bottom of the widget opens the next page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber += 1;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickNextPageBottomPagination(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, PAGINATION_LIMIT * browser.globals.tablePageNumber);

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(page, browser.globals.tablePageNumber);
        });
      },
    );
  },

  'Left arrow at the bottom of the widget opens the previous page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber -= 1;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickPreviousPageBottomPagination(),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, PAGINATION_LIMIT * browser.globals.tablePageNumber);

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(page, browser.globals.tablePageNumber);
        });
      },
    );
  },

  'The button with page number at the bottom of the widget leads to the selected page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber = BOTTOM_PAGINATION_PAGE_NUMBER;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickOnPageBottomPagination(browser.globals.tablePageNumber),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, PAGINATION_LIMIT * browser.globals.tablePageNumber);

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(page, browser.globals.tablePageNumber);

          browser.globals.tablePageNumber = 1;

          commonTable.clickOnPageBottomPagination(browser.globals.tablePageNumber);
        });
      },
    );
  },

  '5 alarms can be shown on the page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setItemPerPage(PAGINATION_PER_PAGE_VALUES.FIVE),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, 5);
      },
    );
  },

  '10 alarms can be shown on the page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber = BOTTOM_PAGINATION_PAGE_NUMBER;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setItemPerPage(PAGINATION_PER_PAGE_VALUES.TEN),
      ({ responseData: { data: [response] } }) => {
        browser.assert.equal(response.last, 10);
      },
    );
  },
};
