// http://nightwatchjs.org/guide#usage

const { WIDGET_TYPES } = require('@/constants');

const {
  FILTERS_TYPE,
  VALUE_TYPES,
  INTERVAL_RANGES,
  FILTER_OPERATORS,
  FILTER_COLUMNS,
} = require('../../../constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const COMPONENT_EQUAL_RESULT_COUNT = 50;
const COMPONENT_AND_RESOURCE_RESULT_COUNT = 5;
const COMPONENT_OR_RESOURCE_RESULT_COUNT = 56;
const CONNECTOR_NAME_EQUAL_VALUE_RESULT_COUNT = 30;
const CONNECTOR_NAME_NOT_EQUAL_VALUE_RESULT_COUNT = 104;
const ALARMS_COUNT = 134;
const INTERVAL_START_DATE = '25/11/2019 00:00';
const INTERVAL_END_DATE = '26/11/2019 00:00';
const INTERVAL_ITEMS_COUNT = 134;
const LAST_SEVEN_DAY_ITEMS_COUNT = 0;

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

  'Add new filters': (browser) => {
    const commonTable = browser.page.tables.common();
    const filtersListModal = browser.page.modals.common.filtersList();
    const createFilterModal = browser.page.modals.common.createFilter();

    commonTable.showFiltersList();
    filtersListModal
      .verifyModalOpened()
      .clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .setFilterTitle('Connector name not equal value')
      .fillFilterGroups([{
        type: FILTERS_TYPE.OR,
        items: [{
          type: FILTERS_TYPE.AND,
          rule: FILTER_COLUMNS.CONNECTOR_NAME,
          operator: FILTER_OPERATORS.NOT_EQUAL,
          valueType: VALUE_TYPES.STRING,
          value: 'feeder2_inst2',
        }],
      }])
      .clickSubmitButton()
      .verifyModalClosed();

    filtersListModal.clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .setFilterTitle('Connector name equal value')
      .fillFilterGroups([{
        type: FILTERS_TYPE.OR,
        items: [{
          rule: FILTER_COLUMNS.CONNECTOR_NAME,
          operator: FILTER_OPERATORS.EQUAL,
          valueType: VALUE_TYPES.STRING,
          value: 'feeder2_inst2',
        }],
      }])
      .clickSubmitButton()
      .verifyModalClosed();

    filtersListModal
      .clickOutside()
      .verifyModalClosed();
  },

  'A filter can be selected': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable
        .clickFilter('Connector name equal value')
        .clickOutsideFilter(),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, CONNECTOR_NAME_EQUAL_VALUE_RESULT_COUNT);
      },
    );
  },

  'A selection of filter can be changed': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable
        .clickFilter('Connector name not equal value')
        .clickOutsideFilter(),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, CONNECTOR_NAME_NOT_EQUAL_VALUE_RESULT_COUNT);
      },
    );
  },

  'The button with cross cancels the selection of filters': (browser) => {
    browser.page.tables.common()
      .clearFilters()
      .assertActiveFilters(0);
  },

  'The "disjunction" (OR) option of "Mix filters" works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    commonTable
      .clickFilter('Connector name not equal value')
      .setMixFilters(true)
      .checkSelectedFilter('Connector name not equal value', true)
      .selectFilter('Connector name equal value', true)
      .checkSelectedFilter('Connector name equal value', true);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.OR),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, ALARMS_COUNT);
      },
    );
  },

  'The "conjunction" (AND) option of "Mix filters" works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.AND),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, 0);
      },
    );
  },

  'The deletion of filter can be canceled': (browser) => {
    const filtersListModal = browser.page.modals.common.filtersList();

    browser.page.tables.common()
      .showFiltersList();

    filtersListModal
      .verifyModalOpened()
      .clickDeleteFilterByName('Connector name equal value');

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();

    filtersListModal.verifyFilterVisibleByName('Connector name equal value');
  },

  'Filter can be deleted': (browser) => {
    const filtersListModal = browser.page.modals.common.filtersList();

    filtersListModal.clickDeleteFilterByName('Connector name equal value');

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    filtersListModal
      .verifyFilterDeletedByName('Connector name equal value')
      .clickOutside()
      .verifyModalClosed();
  },

  'The filter can be changed': (browser) => {
    const createFilterModal = browser.page.modals.common.createFilter();
    const filtersListModal = browser.page.modals.common.filtersList();

    browser.page.tables.common()
      .selectFilter('Connector name not equal value', false)
      .clickOutsideFiltersOptions()
      .showFiltersList();

    filtersListModal
      .verifyModalOpened()
      .clickEditFilterByName('Connector name not equal value');

    createFilterModal
      .verifyModalOpened()
      .clickDeleteRule(createFilterModal.selectGroup([1]), 1)
      .clearFilterTitle()
      .setFilterTitle('Component equal value')
      .fillFilterGroup([1], {
        type: FILTERS_TYPE.AND,
        items: [{
          rule: FILTER_COLUMNS.COMPONENT,
          operator: FILTER_OPERATORS.EQUAL,
          valueType: VALUE_TYPES.STRING,
          value: 'feeder2_0',
        }],
      })
      .clickSubmitButton()
      .verifyModalClosed();

    filtersListModal
      .verifyFilterVisibleByName('Component equal value')
      .clickOutside()
      .verifyModalClosed();
  },

  'The changed filter works in a new way': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.common()
        .clickFilter('Component equal value'),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, COMPONENT_EQUAL_RESULT_COUNT);
      },
    );
  },

  'A new filter can be created': (browser) => {
    const createFilterModal = browser.page.modals.common.createFilter();
    const filtersListModal = browser.page.modals.common.filtersList();

    browser.page.tables.common()
      .clickOutsideFiltersOptions()
      .showFiltersList();

    filtersListModal.clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .setFilterTitle('Resource equal value')
      .fillFilterGroups([{
        type: FILTERS_TYPE.OR,
        items: [{
          rule: FILTER_COLUMNS.RESOURCE,
          operator: FILTER_OPERATORS.EQUAL,
          valueType: VALUE_TYPES.STRING,
          value: 'feeder2_0',
        }],
      }])
      .clickSubmitButton()
      .verifyModalClosed();

    filtersListModal
      .clickOutside()
      .verifyModalClosed();
  },

  'A new filter works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    commonTable
      .setMixFilters(true)
      .setFiltersType(FILTERS_TYPE.AND);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable
        .selectFilter('Resource equal value', true)
        .clickOutsideFilter(),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, COMPONENT_AND_RESOURCE_RESULT_COUNT);
      },
    );

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.OR),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, COMPONENT_OR_RESOURCE_RESULT_COUNT);
      },
    );

    commonTable
      .clearFilters()
      .assertActiveFilters(0);
  },

  '"Live reporting" can be created for period of time selected by user': (browser) => {
    const dateIntervalField = browser.page.fields.dateInterval();
    const alarmsWidget = browser.page.widget.alarms();
    const liveReportingModal = browser.page.modals.common.liveReporting();
    const alarmsTable = browser.page.tables.alarms();

    alarmsTable.clickLiveReportingOpenButton();

    liveReportingModal.verifyModalOpened();

    dateIntervalField
      .selectRange(INTERVAL_RANGES.CUSTOM)
      .clearStartDate()
      .setStartDate(INTERVAL_START_DATE)
      .clearEndDate()
      .setEndDate(INTERVAL_END_DATE);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => liveReportingModal.clickSubmitButton(),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, INTERVAL_ITEMS_COUNT);
      },
    );

    liveReportingModal.verifyModalClosed();
  },

  '"Live reporting" can be created for determined period of time': (browser) => {
    const dateIntervalField = browser.page.fields.dateInterval();
    const alarmsWidget = browser.page.widget.alarms();
    const liveReportingModal = browser.page.modals.common.liveReporting();
    const alarmsTable = browser.page.tables.alarms();

    alarmsTable.clickLiveReportingOpenButton();

    liveReportingModal.verifyModalOpened();

    dateIntervalField.selectRange(INTERVAL_RANGES.LAST_SEVEN_DAY);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => liveReportingModal.clickSubmitButton(),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, LAST_SEVEN_DAY_ITEMS_COUNT);
      },
    );

    liveReportingModal.verifyModalClosed();
  },

  '"Live reporting" can be deleted': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.alarms()
        .clickLiveReportingResetButton(),
      ({ responseData: { data: [alarms] } }) => {
        browser.assert.equal(alarms.total, ALARMS_COUNT);
      },
    );
  },
};
