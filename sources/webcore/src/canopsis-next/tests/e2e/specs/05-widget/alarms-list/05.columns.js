// http://nightwatchjs.org/guide#usage

const { WIDGET_TYPES } = require('@/constants');

const { SORT_ORDERS_STRING } = require('../../../constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
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

    const { widgetId } = await createWidgetForView(viewId, widgetInfo);

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

  'Elements can be sorted by column': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const columnName = 'Connector name';
    const sortFunction = () => commonTable.clickTableHeaderCell(columnName);

    alarmsWidget.waitFirstAlarmsListXHR(sortFunction, ({ responseData }) => {
      browser.assert.equal(responseData.success, true);
      commonTable.checkTableHeaderSort(columnName, SORT_ORDERS_STRING.asc);
    });

    alarmsWidget.waitFirstAlarmsListXHR(sortFunction, ({ responseData }) => {
      browser.assert.equal(responseData.success, true);
      commonTable.checkTableHeaderSort(columnName, SORT_ORDERS_STRING.desc);
    });
  },
};
