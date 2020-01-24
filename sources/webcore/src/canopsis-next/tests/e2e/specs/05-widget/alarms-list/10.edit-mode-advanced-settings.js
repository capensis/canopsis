// http://nightwatchjs.org/guide#usage

import {
  ALARMS_SHARED_ACTIONS,
  SORT_ORDERS,
  SORT_ORDERS_STRING,
  ALARMS_WIDGET_SORT_FIELD,
  PAGINATION_PER_PAGE_VALUES,
  ALARMS_RESOLVED_SHARED_ACTIONS,
  FILTERS_TYPE,
  FILTER_COLUMNS,
  VALUE_TYPES,
  FILTER_OPERATORS,
  INTERVAL_RANGES,
  INFO_POPUP_DEFAULT_COLUMNS,
} from '../../../constants';

const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const DEFAULT_COLUMN_COUNT = 8;
const ALARMS_COUNT = 402;
const NEW_COLUMN_NAME = 'New column';
const NEW_COLUMN_CHANGED_NAME = 'New renamed column';
const CONNECTOR_NAME_EQUAL_FILTER = {
  title: 'Default filter',
  groups: [{
    type: FILTERS_TYPE.OR,
    items: [{
      rule: FILTER_COLUMNS.CONNECTOR_NAME,
      operator: FILTER_OPERATORS.EQUAL,
      valueType: VALUE_TYPES.STRING,
      value: 'feeder2_inst0',
    }],
  }],
};
const RESOURCE_EQUAL_FILTER = {
  title: 'Resource equal value',
  groups: [{
    type: FILTERS_TYPE.OR,
    items: [{
      rule: FILTER_COLUMNS.RESOURCE,
      operator: FILTER_OPERATORS.EQUAL,
      valueType: VALUE_TYPES.STRING,
      value: 'feeder2_0',
    }],
  }],
};
const ALARMS_COUNT_WITH_RESOURCE_EQUAL_FILTER = 40;
const RESOURCE_NOT_EQUAL_FILTER = {
  title: 'Resource not equal value',
  groups: [{
    type: FILTERS_TYPE.OR,
    items: [{
      rule: FILTER_COLUMNS.RESOURCE,
      operator: FILTER_OPERATORS.NOT_EQUAL,
      valueType: VALUE_TYPES.STRING,
      value: 'feeder2_0',
    }],
  }],
};

const INTERVAL_START_DATE = '25/11/2019 00:00';
const INTERVAL_END_DATE = '26/11/2019 00:00';
const INTERVAL_ALARMS_COUNT = 0;
const LAST_SEVEN_DAY_ALARMS_COUNT = 0;
const LAST_YEAR_ALARMS_COUNT = 402;

const ALARM_INFO_POPUP_TEXT = 'Info popup template';
const ALARM_INFO_POPUP_UPDATED_TEXT = 'Info popup template';
const MORE_INFOS_TEXT = 'More infos text';
const MORE_INFOS_CHANGED_TEXT = 'More infos changed text';
const FAST_ACK_OUTPUT_TEXT = 'Fast ack output text';

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
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const newColumn = {
      value: 'alarm.v.connector',
      label: NEW_COLUMN_NAME,
    };

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    const newColumnValue = 'alarm.v.connector_name';
    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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

  'HTML mode can be set for column': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickColumnNames()
      .setColumnNameIsHtml(lastColumnIndex, true)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
        },
      );
  },

  'A column can be deleted from the table': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    const lastColumnIndex = DEFAULT_COLUMN_COUNT + 1;

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

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

  '5 can be set as default number of elements per page': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickElementsPerPage()
      .selectElementsPerPage(PAGINATION_PER_PAGE_VALUES.FIVE);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { data: [response], success } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.last, 5);
      },
    );
  },

  '10 can be set as default number of elements per page': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickElementsPerPage()
      .selectElementsPerPage(PAGINATION_PER_PAGE_VALUES.TEN);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { data: [response], success } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.last, 10);
      },
    );
  },

  'Filter on Open/Resolved can be turn off': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilterOnOpenResolved()
      .setOpenFilter(false)
      .setResolvedFilter(false);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { data: [response], success } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, 0);
      },
    );
  },

  'Filter on Open can be set': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilterOnOpenResolved()
      .setOpenFilter(true)
      .setResolvedFilter(false);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { data: [response], success } }) => {
        browser.globals.temporary.openedAlarmsCount = response.total;
        browser.assert.equal(success, true);
      },
    );
  },

  'Filter on Resolved can be set': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilterOnOpenResolved()
      .setOpenFilter(false)
      .setResolvedFilter(true);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { data: [response], success } }) => {
        browser.globals.temporary.resolvedAlarmsCount = response.total;
        browser.assert.equal(success, true);
      },
    );
  },

  'Filter on Open and Resolved can be set': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilterOnOpenResolved()
      .setOpenFilter(true)
      .setResolvedFilter(true);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { data: [response], success } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, ALARMS_COUNT);
      },
    );
  },

  'Default filter can be created in advanced settings': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const createFilterModal = browser.page.modals.common.createFilter();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilters()
      .clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .clearFilterTitle()
      .setFilterTitle(CONNECTOR_NAME_EQUAL_FILTER.title)
      .fillFilterGroups(CONNECTOR_NAME_EQUAL_FILTER.groups)
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
        commonTable.verifyFilterVisible(CONNECTOR_NAME_EQUAL_FILTER.title);
      },
    );
  },

  'Default filter can be edited in advanced settings': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const createFilterModal = browser.page.modals.common.createFilter();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilters()
      .clickEditFilter(CONNECTOR_NAME_EQUAL_FILTER.title);

    createFilterModal
      .verifyModalOpened()
      .clearFilterTitle()
      .setFilterTitle(RESOURCE_EQUAL_FILTER.title)
      .clickDeleteRule(createFilterModal.selectGroup([1]), 1)
      .fillFilterGroups(RESOURCE_EQUAL_FILTER.groups)
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
        commonTable.verifyFilterVisible(RESOURCE_EQUAL_FILTER.title);
      },
    );

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickFilter(RESOURCE_EQUAL_FILTER.title),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, ALARMS_COUNT_WITH_RESOURCE_EQUAL_FILTER);
        commonTable
          .clickFilter(RESOURCE_EQUAL_FILTER.title)
          .clearFilters();
      },
    );
  },

  'Default filter can be deleted in advanced settings': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilters()
      .clickDeleteFilter(RESOURCE_EQUAL_FILTER.title);

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
        commonTable.verifyFilterDeleted(RESOURCE_EQUAL_FILTER.title);
      },
    );
  },

  'Default filter can be set in advanced settings': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createFilterModal = browser.page.modals.common.createFilter();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilters()
      .clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .clearFilterTitle()
      .setFilterTitle(RESOURCE_EQUAL_FILTER.title)
      .fillFilterGroups(RESOURCE_EQUAL_FILTER.groups)
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.selectFilterByName(RESOURCE_EQUAL_FILTER.title);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, ALARMS_COUNT_WITH_RESOURCE_EQUAL_FILTER);
        commonTable.verifyFilterVisible(RESOURCE_EQUAL_FILTER.title);
      },
    );
  },

  'Two default filters can be set with AND-rule': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createFilterModal = browser.page.modals.common.createFilter();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilters()
      .clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .clearFilterTitle()
      .setFilterTitle(RESOURCE_NOT_EQUAL_FILTER.title)
      .fillFilterGroups(RESOURCE_NOT_EQUAL_FILTER.groups)
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget
      .selectFilterByName(RESOURCE_EQUAL_FILTER.title)
      .setMixFilters(true)
      .setFiltersType(FILTERS_TYPE.AND)
      .selectFilterByName(RESOURCE_NOT_EQUAL_FILTER.title);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        commonTable
          .verifyFilterVisible(RESOURCE_EQUAL_FILTER.title)
          .verifyFilterVisible(RESOURCE_NOT_EQUAL_FILTER.title);
        browser.assert.equal(response.total, 0);
      },
    );
  },

  'Two default filters can be set with OR-rule': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilters()
      .setMixFilters(true)
      .setFiltersType(FILTERS_TYPE.OR);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        commonTable
          .verifyFilterVisible(RESOURCE_EQUAL_FILTER.title)
          .verifyFilterVisible(RESOURCE_NOT_EQUAL_FILTER.title);
        browser.assert.equal(response.total, ALARMS_COUNT);
        commonTable.clearFilters();
      },
    );
  },

  'Live reporting with custom dates can be created': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const liveReportingModal = browser.page.modals.common.liveReporting();
    const dateIntervalField = browser.page.fields.dateInterval();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget.clickCreateLiveReporting();

    liveReportingModal.verifyModalOpened();

    dateIntervalField
      .selectRange(INTERVAL_RANGES.CUSTOM)
      .clearStartDate()
      .setStartDate(INTERVAL_START_DATE)
      .clearEndDate()
      .setEndDate(INTERVAL_END_DATE);

    liveReportingModal
      .clickSubmitButton()
      .verifyModalClosed();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, INTERVAL_ALARMS_COUNT);
      },
    );
  },

  'Live reporting with determined dates can be created': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const liveReportingModal = browser.page.modals.common.liveReporting();
    const dateIntervalField = browser.page.fields.dateInterval();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget.clickEditLiveReporting();

    liveReportingModal.verifyModalOpened();

    dateIntervalField.selectRange(INTERVAL_RANGES.LAST_SEVEN_DAY);

    liveReportingModal
      .clickSubmitButton()
      .verifyModalClosed();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, LAST_SEVEN_DAY_ALARMS_COUNT);
      },
    );
  },

  'Live reporting can be edited': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const liveReportingModal = browser.page.modals.common.liveReporting();
    const dateIntervalField = browser.page.fields.dateInterval();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget.clickEditLiveReporting();

    liveReportingModal.verifyModalOpened();

    dateIntervalField.selectRange(INTERVAL_RANGES.LAST_YEAR);

    liveReportingModal
      .clickSubmitButton()
      .verifyModalClosed();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, LAST_YEAR_ALARMS_COUNT);
      },
    );
  },

  'Live reporting can be deleted': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget.clickDeleteLiveReporting();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.assert.equal(response.total, ALARMS_COUNT);
        browser.globals.temporary.alarmsList = response.alarms;
      },
    );
  },

  'Info popup can be created': (browser) => {
    const addInfoPopupModal = browser.page.modals.common.addInfoPopup();
    const infoPopupModal = browser.page.modals.common.infoPopupSetting();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickCreateOrEditInfoPopup();

    infoPopupModal
      .verifyModalOpened()
      .clickAddPopup();

    addInfoPopupModal
      .verifyModalOpened()
      .selectSelectedColumn(INFO_POPUP_DEFAULT_COLUMNS.connectorName)
      .setTemplate(ALARM_INFO_POPUP_TEXT)
      .clickSubmitButton()
      .verifyModalClosed();

    infoPopupModal
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
      },
    );
  },

  'Info popup should displayed on table column': (browser) => {
    const alarmsTable = browser.page.tables.alarms();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    const column = 'Connector name';

    alarmsTable
      .clickOnRowInfoPopupOpenButton(firstAlarm._id, column)
      .verifyRowInfoPopupVisible(firstAlarm._id, column)
      .getRowInfoPopupText(firstAlarm._id, column, (text) => {
        browser.assert.equal(text, ALARM_INFO_POPUP_TEXT);
      });
  },

  'Info popup can be edited': (browser) => {
    const addInfoPopupModal = browser.page.modals.common.addInfoPopup();
    const infoPopupModal = browser.page.modals.common.infoPopupSetting();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const alarmsTable = browser.page.tables.alarms();

    const column = 'Connector name';

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickCreateOrEditInfoPopup();

    infoPopupModal
      .verifyModalOpened()
      .clickEditPopup(1);

    addInfoPopupModal
      .verifyModalOpened()
      .selectSelectedColumn(INFO_POPUP_DEFAULT_COLUMNS.connectorName)
      .clickTemplate()
      .clearTemplate()
      .setTemplate(ALARM_INFO_POPUP_UPDATED_TEXT)
      .clickSubmitButton()
      .verifyModalClosed();

    infoPopupModal
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
        const [firstAlarm] = browser.globals.temporary.alarmsList;

        alarmsTable
          .clickOnRowInfoPopupOpenButton(firstAlarm._id, column)
          .verifyRowInfoPopupVisible(firstAlarm._id, column)
          .getRowInfoPopupText(firstAlarm._id, column, (text) => {
            browser.assert.equal(text, ALARM_INFO_POPUP_UPDATED_TEXT);
          });
      },
    );
  },

  'More than one info popup can be added': (browser) => {
    const infoPopupModal = browser.page.modals.common.infoPopupSetting();
    const addInfoPopupModal = browser.page.modals.common.addInfoPopup();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const alarmsTable = browser.page.tables.alarms();

    const column = 'Connector';

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickCreateOrEditInfoPopup();

    infoPopupModal
      .verifyModalOpened()
      .clickAddPopup();

    addInfoPopupModal
      .verifyModalOpened()
      .selectSelectedColumn(INFO_POPUP_DEFAULT_COLUMNS.connector)
      .setTemplate(ALARM_INFO_POPUP_TEXT)
      .clickSubmitButton()
      .verifyModalClosed();

    infoPopupModal
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
        const [firstAlarm] = browser.globals.temporary.alarmsList;

        alarmsTable
          .clickOnRowInfoPopupOpenButton(firstAlarm._id, column)
          .verifyRowInfoPopupVisible(firstAlarm._id, column)
          .getRowInfoPopupText(firstAlarm._id, column, (text) => {
            browser.assert.equal(text, ALARM_INFO_POPUP_UPDATED_TEXT);
          });
      },
    );
  },

  'Info popup can be deleted': (browser) => {
    const infoPopupModal = browser.page.modals.common.infoPopupSetting();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const alarmsTable = browser.page.tables.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickCreateOrEditInfoPopup();

    infoPopupModal
      .verifyModalOpened()
      .clickDeletePopup(2)
      .clickDeletePopup(1)
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
        const [firstAlarm] = browser.globals.temporary.alarmsList;

        alarmsTable.verifyRowInfoPopupOpenButtonDeleted(firstAlarm._id, 'Connector name');
        alarmsTable.verifyRowInfoPopupOpenButtonDeleted(firstAlarm._id, 'Connector');
      },
    );
  },

  '"More infos" popup can be created': (browser) => {
    const textEditorModal = browser.page.modals.common.textEditor();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilterOnOpenResolved()
      .setOpenFilter(false)
      .setResolvedFilter(true)
      .clickCreateMoreInfos();

    textEditorModal
      .verifyModalOpened()
      .clickField()
      .setField(MORE_INFOS_TEXT)
      .clickSubmitButton()
      .verifyModalClosed();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.assert.equal(success, true);
        browser.globals.temporary.alarmsList = response.alarms;
      },
    );
  },

  '"More infos" modal window is showing': (browser) => {
    const moreInfosModal = browser.page.modals.alarm.moreInfos();
    const commonTable = browser.page.tables.common();

    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnSharedAction(firstAlarm._id, ALARMS_RESOLVED_SHARED_ACTIONS.MORE_INFOS);
    moreInfosModal
      .verifyModalOpened()
      .getContentText((text) => {
        browser.assert.equal(text, MORE_INFOS_TEXT);
      })
      .clickOutside()
      .verifyModalClosed();
  },

  '"More infos" popup can be edited': (browser) => {
    const textEditorModal = browser.page.modals.common.textEditor();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickEditMoreInfos();

    textEditorModal
      .verifyModalOpened()
      .clickField()
      .clearField()
      .setField(MORE_INFOS_CHANGED_TEXT)
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
      },
    );
  },

  '"More infos" modal window is showing with edited data': (browser) => {
    const moreInfosModal = browser.page.modals.alarm.moreInfos();
    const commonTable = browser.page.tables.common();

    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .verifySharedActionVisible(firstAlarm._id, ALARMS_RESOLVED_SHARED_ACTIONS.MORE_INFOS)
      .clickOnSharedAction(firstAlarm._id, ALARMS_RESOLVED_SHARED_ACTIONS.MORE_INFOS);

    moreInfosModal
      .verifyModalOpened()
      .getContentText((text) => {
        browser.assert.equal(text, MORE_INFOS_CHANGED_TEXT);
      })
      .clickOutside()
      .verifyModalClosed();
  },

  '"More infos" popup can be deleted': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickDeleteMoreInfos();

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);

        const [firstAlarm] = browser.globals.temporary.alarmsList;
        commonTable.verifySharedActionDeleted(firstAlarm._id, ALARMS_RESOLVED_SHARED_ACTIONS.MORE_INFOS);
      },
    );
  },

  '"HTML enabled on timeline?" checkbox can be turn on': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget
      .setEnableHtml(true)
      .checkEnableHtmlValue(true);

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
      },
    );
  },

  '"Note field required when ack?" checkbox can be turn on': (browser) => {
    const createAckEventModal = browser.page.modals.alarm.createAckEvent();
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget
      .clickAdvancedSettings()
      .clickFilterOnOpenResolved()
      .setOpenFilter(true)
      .setResolvedFilter(false);

    alarmsWidget
      .clickAckGroup()
      .setIsAckNoteRequired(true);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success, data: [response] } }) => {
        browser.globals.temporary.alarmsList = response.alarms;

        browser.assert.equal(success, true);

        const [firstAlarm] = browser.globals.temporary.alarmsList;

        commonTable.clickOnSharedAction(firstAlarm._id, ALARMS_SHARED_ACTIONS.ACK);

        createAckEventModal
          .verifyModalOpened()
          .clickSubmitButton()
          .verifyModalOpened()
          .clickCancelButton()
          .verifyModalClosed();
      },
    );
  },

  '"Multiple ack" checkbox can be turn on': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget
      .clickAckGroup()
      .clickFastAckOutput()
      .setIsMultiAckEnabled(true)
      .checkIsMultiAckEnabled(true);
  },

  'Fast-ack output (optional) can be enabled': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    browser.page.view()
      .openWidgetSettings(browser.globals.defaultViewData.widgetId);

    commonWidget.clickAdvancedSettings();
    alarmsWidget
      .clickAckGroup()
      .clickFastAckOutput()
      .setFastAckOutputSwitch(true)
      .checkFastAckOutputSwitch(true);
  },

  'Text of fast-ack output (optional) can be edited': (browser) => {
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget
      .clickFastAckOutputText()
      .clearFastAckOutputText()
      .setFastAckOutputText(FAST_ACK_OUTPUT_TEXT);

    commonWidget.waitFirstUserPreferencesXHR(
      () => alarmsWidget.clickSubmitAlarms(),
      ({ responseData: { success } }) => {
        browser.assert.equal(success, true);
      },
    );
  },
};
