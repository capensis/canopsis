// http://nightwatchjs.org/guide#usage

import { SORT_ORDERS, SORT_ORDERS_STRING, ALARMS_WIDGET_SORT_FIELD } from '../../../constants';

const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const DEFAULT_COLUMN_COUNT = 8;
const NEW_COLUMN_NAME = 'New column';
const NEW_COLUMN_CHANGED_NAME = 'New renamed column';

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
    browser.globals.temporary.alarmsList = [];
    browser.globals.defaultViewData = await createWidgetView();

    const { viewId } = browser.globals.defaultViewData;

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

  'Open view': (browser) => {
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.widget.alarms()
      .waitFirstAlarmsListXHR(
        () => browser.page.layout.groupsSideBar()
          .clickGroupsSideBarButton()
          .clickPanelHeader(groupId)
          .clickLinkView(viewId),
        ({ responseData: { data: [response] } }) => {
          browser.globals.temporary.alarmsList = response.alarms;
          browser.page.view()
            .clickMenuViewButton()
            .clickEditViewButton();
        },
      );
  },

  'Default sort column can be set for table': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickDefaultSortColumn()
      .selectSortOrderBy(ALARMS_WIDGET_SORT_FIELD.connectorName)
      .selectSortOrders(SORT_ORDERS.asc)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonTable.checkTableHeaderSort('Connector name', SORT_ORDERS_STRING.asc);
        },
      );
  },

  'A column can be added to table': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const newColumn = {
      value: 'alarm.v.connector',
      label: NEW_COLUMN_NAME,
    };

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .clickAddColumnName()
      .editColumnName(DEFAULT_COLUMN_COUNT + 1, newColumn)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonTable.verifyTableColumnVisible(NEW_COLUMN_NAME);
        },
      );
  },

  'A column is name can be changed': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .clickColumnNameLabel(lastColumnIndex)
      .clearColumnNameLabel(lastColumnIndex)
      .setColumnNameLabel(lastColumnIndex, NEW_COLUMN_CHANGED_NAME)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonTable.verifyTableColumnVisible(NEW_COLUMN_CHANGED_NAME);
        },
      );
  },

  'A column is value can be changed': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    const newColumnValue = 'alarm.v.connector_name';
    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .clickColumnNameValue(lastColumnIndex)
      .clearColumnNameValue(lastColumnIndex)
      .setColumnNameValue(lastColumnIndex, newColumnValue)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);

          commonTable.getCellTextByColumnName(firstAlarm._id, NEW_COLUMN_CHANGED_NAME, ({ value }) => {
            commonTable.getCellTextByColumnName(firstAlarm._id, 'Connector name', ({ value: expectedValue }) => {
              browser.assert.equal(value, expectedValue);
            });
          });
        },
      );
  },

  'A column is card can be moved above': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .clickColumnNameUpWard(lastColumnIndex)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonTable.getTableHeaderTextByIndex(DEFAULT_COLUMN_COUNT, ({ value }) => {
            browser.assert.equal(value, NEW_COLUMN_CHANGED_NAME);
          });
        },
      );
  },

  'A column is card can be moved below': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .clickColumnNameDownWard(DEFAULT_COLUMN_COUNT)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonTable.getTableHeaderTextByIndex(lastColumnIndex, ({ value }) => {
            browser.assert.equal(value, NEW_COLUMN_CHANGED_NAME);
          });
        },
      );
  },

  'A column can be deleted from the table': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .clickDeleteColumnName(lastColumnIndex)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonTable.verifyTableColumnDeleted(NEW_COLUMN_CHANGED_NAME);
        },
      );
  },
};
