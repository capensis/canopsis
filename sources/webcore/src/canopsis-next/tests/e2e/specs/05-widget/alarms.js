// http://nightwatchjs.org/guide#usage

const {
  INFO_POPUP_DEFAULT_COLUMNS,
  ALARMS_WIDGET_SORT_FIELD,
  SORT_ORDERS,
  SORT_ORDERS_STRING,
  PAGINATION_PER_PAGE_VALUES,
  FILTERS_TYPE,
  VALUE_TYPES,
  INTERVAL_RANGES,
  FILTER_OPERATORS,
  ALARMS_MASS_ACTIONS,
  MONTH,
  WEEK_DAYS,
  PBEHAVIOR_TYPES,
  DATE_INTERVAL_MINUTES,
  ALARMS_SHARED_ACTIONS,
  PERIODICAL_BEHAVIOR_RESONES,
  PERIODICAL_BEHAVIOR_FREQUENCY,
  ALARMS_SHARED_DROPDOWN_ACTIONS,
  FILTER_COLUMNS,
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, removeWidgetView } = require('../../helpers/api');
const { generateTemporaryAlarmsWidget } = require('../../helpers/entities');
const getPaginationFirstIndex = require('../../helpers/getPaginationFirstIndex');

const SEARCH_STRING = 'feeder2_inst3';
const SEARCH_RESULT_COUNT = 16;
const COMPONENT_EQUAL_RESULT_COUNT = 50;
const COMPONENT_AND_RESOURCE_RESULT_COUNT = 5;
const COMPONENT_OR_RESOURCE_RESULT_COUNT = 57;
const CONNECTOR_NAME_EQUAL_VALUE_RESULT_COUNT = 30;
const CONNECTOR_NAME_NOT_EQUAL_VALUE_RESULT_COUNT = 106;
const ALARMS_COUNT = 136;
const INTERVAL_START_DATE = '25/11/2019 00:00';
const INTERVAL_END_DATE = '26/11/2019 00:00';
const INTERVAL_ITEMS_COUNT = 136;
const LAST_SEVEN_DAY_ITEMS_COUNT = 136;

function tableColumnSortAssertHelper(browser, columnName) {
  const alarmsWidget = browser.page.widget.alarms();
  const commonTable = browser.page.tables.common();

  const sortFunction = () => commonTable.clickTableHeaderCell(columnName);

  alarmsWidget.waitFirstAlarmsListXHR(sortFunction, ({ responseData }) => {
    browser.assert.equal(responseData.success, true);
    commonTable.checkTableHeaderSort(columnName, SORT_ORDERS_STRING.asc);
  });

  alarmsWidget.waitFirstAlarmsListXHR(sortFunction, ({ responseData }) => {
    browser.assert.equal(responseData.success, true);
    commonTable.checkTableHeaderSort(columnName, SORT_ORDERS_STRING.desc);
  });
}

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
    browser.globals.defaultViewData = await createWidgetView();
    browser.globals.tablePageNumber = 1;

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

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

  'Create widget alarms with some name': (browser) => {
    const alarmsWidget = {
      ...generateTemporaryAlarmsWidget(),
      periodicRefresh: 140,
      parameters: {
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: ALARMS_WIDGET_SORT_FIELD.component,
        },
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.TWENTY,
        openedResolvedFilter: {
          open: true,
          resolve: true,
        },
        infoPopups: [{
          column: INFO_POPUP_DEFAULT_COLUMNS.connectorName,
          template: 'Info popup template',
        }],
        ack: {
          isAckNoteRequired: true,
          isMultiAckEnabled: true,
          fastAckOutput: {
            enabled: true,
            output: 'Output',
          },
        },
        moreInfos: 'More infos popup',
        enableHtml: true,
        liveReporting: {
          range: INTERVAL_RANGES.LAST_YEAR,
        },
        filters: {
          isMix: true,
          type: FILTERS_TYPE.OR,
          title: 'Filter 1',
          groups: [{
            type: FILTERS_TYPE.OR,
            items: [{
              rule: FILTER_COLUMNS.CONNECTOR,
              operator: FILTER_OPERATORS.EQUAL,
              valueType: VALUE_TYPES.STRING,
              value: 'value',
            }, {
              type: FILTERS_TYPE.AND,
              rule: FILTER_COLUMNS.CONNECTOR_NAME,
              operator: FILTER_OPERATORS.NOT_EQUAL,
              valueType: VALUE_TYPES.NUMBER,
              value: 136,
            }],
          }],
        },
      },
    };
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .clickPanelHeader(groupId)
      .clickLinkView(viewId);

    browser.page.view()
      .clickMenuViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget(WIDGET_TYPES.alarmList)
      .verifyModalClosed();

    browser.completed.widget.createAlarmsList(alarmsWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'The empty search shows no results': (browser) => {
    const commonTable = browser.page.tables.common();
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitAllAlarmsListXHR(
      () => commonTable.keyupSearchEnter(),
      (xhrs) => {
        browser.assert.equal(xhrs.length, 0);
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
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, SEARCH_RESULT_COUNT);

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
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, SEARCH_RESULT_COUNT);
      },
    );
  },

  'The button with cross cancels current search': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickSearchResetButton(),
      ({ responseData: { data } }) => {
        const response = data[0];

        browser.assert.notEqual(response.total, SEARCH_RESULT_COUNT);
        browser.assert.equal(response.total, ALARMS_COUNT);

        browser.globals.temporary.alarmsList = response.alarms;
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

  'Right arrow opens the next page': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    browser.globals.tablePageNumber += 1;

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickNextPageTopPagination(),
      ({ responseData: { data } }) => {
        browser.assert.equal(
          data[0].first,
          getPaginationFirstIndex(browser.globals.tablePageNumber, 20),
        );

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
      ({ responseData: { data } }) => {
        browser.assert.equal(
          data[0].first,
          getPaginationFirstIndex(browser.globals.tablePageNumber, 20),
        );

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(page, browser.globals.tablePageNumber);
        });
      },
    );
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

    commonTable
      .clickFilter('Connector name equal value')
      .clickOutsideFilter();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setMixFilters(false),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, CONNECTOR_NAME_EQUAL_VALUE_RESULT_COUNT);
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
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, CONNECTOR_NAME_NOT_EQUAL_VALUE_RESULT_COUNT);
      },
    );
  },

  'The button with cross cancels the selection of filters': (browser) => {
    browser.page.tables.common()
      .clearFilters()
      .assertActiveFilters(0);
  },

  'The "conjunction" (AND) option of "Mix filters" works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    commonTable
      .clickFilter('Connector name not equal value')
      .setMixFilters(true)
      .checkSelectedFilter('Connector name not equal value', true)
      .selectFilter('Connector name equal value', true)
      .checkSelectedFilter('Connector name equal value', true);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.AND),
      ({ responseData: { data } }) => {
        browser.assert.equal(0, data[0].total);
      },
    );
  },

  'The "disjunction" (OR) option of "Mix filters" works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.OR),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, ALARMS_COUNT);
      },
    );
  },

  'The deletion of filter can be canceled': (browser) => {
    const filtersListModal = browser.page.modals.common.filtersList();

    browser.page.tables.common()
      .showFiltersList();

    filtersListModal
      .verifyModalOpened()
      .clickDeleteFilter('Connector name equal value');

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();

    filtersListModal.verifyFilterVisible('Connector name equal value');
  },

  'Filter can be deleted': (browser) => {
    const filtersListModal = browser.page.modals.common.filtersList();

    filtersListModal.clickDeleteFilter('Connector name equal value');

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    filtersListModal
      .verifyFilterDeleted('Connector name equal value')
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
      .clickEditFilter('Connector name not equal value');

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
      .verifyFilterVisible('Component equal value')
      .clickOutside()
      .verifyModalClosed();
  },

  'The changed filter works in a new way': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.common()
        .selectFilter('Component equal value', true),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, COMPONENT_EQUAL_RESULT_COUNT);
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

    commonTable.setFiltersType(FILTERS_TYPE.AND);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.selectFilter('Resource equal value', true),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, COMPONENT_AND_RESOURCE_RESULT_COUNT);
      },
    );

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.OR),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, COMPONENT_OR_RESOURCE_RESULT_COUNT);
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

    alarmsTable.openLiveReporting();

    liveReportingModal.verifyModalOpened();

    dateIntervalField
      .selectRange(INTERVAL_RANGES.CUSTOM)
      .clearStartDate()
      .setStartDate(INTERVAL_START_DATE)
      .clearEndDate()
      .setEndDate(INTERVAL_END_DATE);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => liveReportingModal.clickSubmitButton(),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, INTERVAL_ITEMS_COUNT);
      },
    );

    liveReportingModal.verifyModalClosed();
  },

  '"Live reporting" can be created for determined period of time': (browser) => {
    const dateIntervalField = browser.page.fields.dateInterval();
    const alarmsWidget = browser.page.widget.alarms();
    const liveReportingModal = browser.page.modals.common.liveReporting();
    const alarmsTable = browser.page.tables.alarms();

    alarmsTable.openLiveReporting();

    liveReportingModal.verifyModalOpened();

    dateIntervalField.selectRange(INTERVAL_RANGES.LAST_SEVEN_DAY);

    alarmsWidget.waitFirstAlarmsListXHR(
      () => liveReportingModal.clickSubmitButton(),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, LAST_SEVEN_DAY_ITEMS_COUNT);
      },
    );

    liveReportingModal.verifyModalClosed();
  },

  '"Live reporting" can be deleted': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.alarms()
        .clickResetLiveReporting(),
      ({ responseData: { data } }) => {
        browser.assert.equal(data[0].total, ALARMS_COUNT);
      },
    );
  },

  'All elements in the table can be selected': (browser) => {
    browser.page.tables.common()
      .setAllCheckbox(true)
      .assertActiveCheckboxCount(20)
      .moveOutsideMassActionsPanel()
      .setAllCheckbox(false)
      .assertActiveCheckboxCount(0);
  },

  'The only one element in the table can be selected': (browser) => {
    const commonTable = browser.page.tables.common();
    const [firstAlarm, secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .checkRowCheckboxValue(firstAlarm._id, (value) => {
        browser.assert.equal(value, 'true');
      })
      .setRowCheckbox(secondAlarm._id, true)
      .checkRowCheckboxValue(secondAlarm._id, (value) => {
        browser.assert.equal(value, 'true');
      })
      .assertActiveCheckboxCount(2)
      .setRowCheckbox(firstAlarm._id, false)
      .checkRowCheckboxValue(firstAlarm._id, (value) => {
        browser.assert.equal(value, 'false');
      })
      .setRowCheckbox(secondAlarm._id, false)
      .checkRowCheckboxValue(secondAlarm._id, (value) => {
        browser.assert.equal(value, 'false');
      })
      .assertActiveCheckboxCount(0);
  },

  'Pressing on button "Periodical behavior" creates periodical behavior': () => {},

  'An acknowledge without a ticket can be created': (browser) => {
    const commonTable = browser.page.tables.common();
    const createAckEventModal = browser.page.modals.alarm.createAckEvent();
    const confirmAckWithTcketModal = browser.page.modals.alarm.confirmAckWithTcket();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.ACK);

    createAckEventModal
      .verifyModalOpened()
      .clearTicketNumber()
      .setTicketNumber(165558555)
      .clearTicketNote()
      .setTicketNote('Note text')
      .setAckTicketResources(true)
      .clickSubmitButton();

    confirmAckWithTcketModal
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    createAckEventModal.verifyModalClosed();
  },

  'An acknowledge with ticket can be created': (browser) => {
    const commonTable = browser.page.tables.common();
    const createAckEventModal = browser.page.modals.alarm.createAckEvent();
    const [secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(secondAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.ACK);

    createAckEventModal
      .verifyModalOpened()
      .clearTicketNumber()
      .setTicketNumber(165558556)
      .clearTicketNote()
      .setTicketNote('Note text')
      .setAckTicketResources(true)
      .clickSubmitButtonWithTicket()
      .verifyModalClosed();
  },

  '"Fast ack" can be created': (browser) => {
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.FAST_ACK);
  },

  'An "ack" can be canceled': (browser) => {
    const commonTable = browser.page.tables.common();
    const [firstAlarm, secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ACK);

    browser.page.modals.alarm.createCancelEvent()
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('Cancel first alarm ack')
      .clickSubmitButton()
      .verifyModalClosed();

    commonTable
      .setRowCheckbox(secondAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ACK);

    browser.page.modals.alarm.createCancelEvent()
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('Cancel second alarm ack')
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'An alarm can be canceled': (browser) => {
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ALARM);

    browser.page.modals.alarm.createCancelEvent()
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('Cancel first alarm')
      .clickSubmitButton()
      .verifyModalClosed();

    commonTable.setRowCheckbox(firstAlarm._id, false);
  },

  'Elements can be sorted by connector name': (browser) => {
    tableColumnSortAssertHelper(browser, 'Connector name');
  },

  'Elements can be sorted by connector': (browser) => {
    tableColumnSortAssertHelper(browser, 'Connector');
  },

  'Elements can be sorted by component': (browser) => {
    tableColumnSortAssertHelper(browser, 'Component');
  },

  'Elements can be sorted by resource': (browser) => {
    tableColumnSortAssertHelper(browser, 'Resource');
  },

  'Elements can be sorted by output': (browser) => {
    tableColumnSortAssertHelper(browser, 'Output');
  },

  'Elements can be sorted by extra details': (browser) => {
    tableColumnSortAssertHelper(browser, 'Extra details');
  },

  'Elements can be sorted by state': (browser) => {
    tableColumnSortAssertHelper(browser, 'State');
  },

  'Elements can be sorted by status': (browser) => {
    tableColumnSortAssertHelper(browser, 'Status');
  },

  'Information pop-up can be shown': (browser) => {
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .clickOnRowInfoPopupOpen(firstAlarm._id)
      .verifyRowInfoPopupVisible(firstAlarm._id);
  },

  'Information pop-up can be closed': (browser) => {
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .clickOnRowInfoPopupClose(firstAlarm._id)
      .verifyRowInfoPopupDeleted(firstAlarm._id);
  },

  'Pressing on element shows details about this element': (browser) => {
    const commonTable = browser.page.tables.common();
    const alarmsTable = browser.page.tables.alarms();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnRow(firstAlarm._id);
    alarmsTable.verifyAlarmTimeLineVisible(firstAlarm._id);
  },

  'Pressing on element hidden details about this element': (browser) => {
    const commonTable = browser.page.tables.common();
    const alarmsTable = browser.page.tables.alarms();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnRow(firstAlarm._id);
    alarmsTable.verifyAlarmTimeLineDeleted(firstAlarm._id);
  },

  'Placing a cursor on signs in the column "Extra details" makes information pop-up show': () => {},

  'Table widget alarms': (browser) => {
    const commonTable = browser.page.tables.common();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();

    const firstId = '6770ba94-51d9-4b8c-ae85-7c62fab18a54';

    browser.page.view()
      .clickMenuViewButton();

    commonTable.clickOnSharedAction(firstId, ALARMS_SHARED_ACTIONS.SNOOZE_ALARM);

    browser.page.modals.alarm.createSnoozeEvent()
      .verifyModalOpened()
      .clickDurationValue()
      .clearDurationValue()
      .setDurationValue(10)
      .setDurationType(PBEHAVIOR_TYPES.MAINTENANCE)
      .clickCancelButton()
      .verifyModalClosed();

    commonTable
      .clickOnDropDownActions(firstId, ALARMS_SHARED_DROPDOWN_ACTIONS.PERIODICAL_BEHAVIOR);

    browser.page.modals.alarm.createPbehavior()
      .verifyModalOpened();

    pbehaviorForm
      .clearName()
      .clickName()
      .setName('Name')
      .clickStartDate();

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(3)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.FIVE);

    pbehaviorForm
      .clickName()
      .clickEndDate();

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(3)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.TEN);

    pbehaviorForm
      .clickName()
      .selectType(1)
      .clearReason()
      .clickReason()
      .setReason('P')
      .selectReason(PERIODICAL_BEHAVIOR_RESONES.REHABILITATION_PROBLEM)
      .setRuleCheckbox(true)
      .selectFrequency(PERIODICAL_BEHAVIOR_FREQUENCY.MINUTELY)
      .selectByWeekDay(WEEK_DAYS.TUESDAY, false)
      .selectByWeekDay(WEEK_DAYS.TUESDAY, true)
      .clearRepeat()
      .clickRepeat()
      .setRepeat(5)
      .clearInterval()
      .clickInterval()
      .setInterval(5)
      .clickAdvanced(5)
      .selectWeekStart(WEEK_DAYS.TUESDAY)
      // .selectByMonth(WEEK_DAYS.TUESDAY, false)
      .selectByMonth(MONTH.JUNUARY, true)
      .clearBySetPosition()
      .clickBySetPosition()
      .setBySetPosition(15)
      .clearByMonthDay()
      .clickByMonthDay()
      .setByMonthDay(12)
      .clearByYearDay()
      .clickByYearDay()
      .setByYearDay(23)
      .clearByWeekNo()
      .clickByWeekNo()
      .setByWeekNo(1)
      .clearByHour()
      .clearByHour()
      .clickByHour()
      .setByHour(2)
      .clearByMinute()
      .clickByMinute()
      .setByMinute(2)
      .clearBySecond()
      .clickBySecond()
      .setBySecond(2)
      .clickAddExdate()
      .clickExdateField(1);

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(3)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.ZERO);

    browser.page.forms.pbehavior()
      .clickSimple()
      .clickAddComment()
      .clickCommentField(1)
      .clearCommentField(1)
      .clearCommentField(1)
      .setCommentField(1, 2)
      .clickCommentDelete(1);

    browser.page.modals.alarm.createPbehavior()
      .clickCancelButton()
      .verifyModalClosed();
  },

  'Delete widget alarms with some name': (browser) => {
    browser.page.view()
      .clickMenuViewButton()
      .clickDeleteWidgetButton(browser.globals.temporary.widgetId);

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },

  'Delete row with some name': (browser) => {
    browser.page.view()
      .clickDeleteRowButton(1);

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();
  },
};
