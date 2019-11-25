// http://nightwatchjs.org/guide#usage

const {
  INFO_POPUP_DEFAULT_COLUMNS,
  ALARMS_WIDGET_SORT_FIELD,
  SORT_ORDERS,
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
const SEARCH_RESULT_COUNT = 64;
const COMPONENT_EQUAL_RESULT_COUNT = 46;
const COMPONENT_AND_RESOURCE_RESULT_COUNT = 5;
const COMPONENT_OR_RESOURCE_RESULT_COUNT = 91;
const ALARMS_COUNT = 496;

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
          // calendarStartDate: {
          //   minute: 15,
          //   hour: 12,
          //   day: 12,
          // },
        //   endDate: '13/09/2019 00:00',
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
        browser.assert.equal(xhrs.length, 0); // TODO test break
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
        browser.assert.equal(SEARCH_RESULT_COUNT, data[0].total);

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
        browser.assert.equal(SEARCH_RESULT_COUNT, data[0].total);
      },
    );
  },

  'The button with cross cancels current search': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.clickSearchResetButton(),
      ({ responseData: { data } }) => {
        browser.assert.notEqual(SEARCH_RESULT_COUNT, data[0].total);
        browser.assert.equal(ALARMS_COUNT, data[0].total);
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
          getPaginationFirstIndex(browser.globals.tablePageNumber, 20),
          data[0].first,
        );

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(browser.globals.tablePageNumber, page);
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
          getPaginationFirstIndex(browser.globals.tablePageNumber, 20),
          data[0].first,
        );

        commonTable.getTopPaginationPage((page) => {
          browser.assert.equal(browser.globals.tablePageNumber, page);
        });
      },
    );
  },

  'Add new filters': (browser) => {
    const commonWidget = browser.page.widget.common();
    const commonTable = browser.page.tables.common();
    const filtersListModal = browser.page.modals.common.filtersList();
    const createFilterModal = browser.page.modals.common.createFilter();

    commonTable.showFiltersList();
    filtersListModal.verifyModalOpened();
    commonWidget.clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .setFilterTitle('Connector name not equal value')
      .fillFilterGroups([{
        type: FILTERS_TYPE.OR,
        items: [{
          type: FILTERS_TYPE.AND,
          rule: FILTER_COLUMNS.CONNECTOR_NAME,
          operator: FILTER_OPERATORS.EQUAL,
          valueType: VALUE_TYPES.STRING,
          value: 'feeder2_inst2',
        }],
      }])
      .clickSubmitButton()
      .verifyModalClosed();

    commonWidget.clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .setFilterTitle('Connector name equal value')
      .fillFilterGroups([{
        type: FILTERS_TYPE.OR,
        items: [{
          rule: FILTER_COLUMNS.CONNECTOR_NAME,
          operator: FILTER_OPERATORS.NOT_EQUAL,
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
    browser.page.tables.common()
      .selectFilter('Connector name equal value', true)
      .checkSelectedFilter('Connector name equal value', true);
  },

  'A selection of filter can be changed': (browser) => {
    browser.page.tables.common()
      .selectFilter('Connector name equal value', false)
      .selectFilter('Connector name not equal value', true)
      .checkSelectedFilter('Connector name equal value', false)
      .checkSelectedFilter('Connector name not equal value', true)
      .selectFilter('Connector name equal value', true)
      .checkSelectedFilter('Connector name equal value', true)
      .clickOutsideFilter();
  },

  'The "conjunction" (AND) option of "Mix filters" works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.common()
        .setMixFilters(true)
        .setFiltersType(FILTERS_TYPE.AND),
      ({ responseData: { data } }) => {
        browser.assert.equal(0, data[0].total);
      },
    );
  },

  'The "disjunction" (OR) option of "Mix filters" works correctly': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.common()
        .setMixFilters(true)
        .setFiltersType(FILTERS_TYPE.OR),
      ({ responseData: { data } }) => {
        browser.assert.equal(ALARMS_COUNT, data[0].total);
      },
    );
  },

  'The deletion of filter can be canceled': (browser) => {
    const commonWidget = browser.page.widget.common();

    browser.page.tables.common()
      .showFiltersList();

    browser.page.modals.common.filtersList()
      .verifyModalOpened();

    commonWidget.clickDeleteFilter('Connector name equal value');

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();

    commonWidget.verifyFilterVisible('Connector name equal value');
  },

  'Filter can be deleted': (browser) => {
    browser.page.widget.common()
      .clickDeleteFilter('Connector name equal value');

    browser.page.modals.common.confirmation()
      .verifyModalOpened()
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.widget.common()
      .verifyFilterDeleted('Connector name equal value');
  },

  'The filter can be changed': (browser) => {
    const createFilterModal = browser.page.modals.common.createFilter();

    browser.page.widget.common()
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
          value: 'feeder2_9',
        }],
      })
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.widget.common()
      .verifyFilterVisible('Component equal value');

    browser.page.modals.common.filtersList()
      .clickOutside()
      .verifyModalClosed();
  },

  'The changed filter works in a new way': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();

    alarmsWidget.waitFirstAlarmsListXHR(
      () => browser.page.tables.common()
        .selectFilter('Component equal value', true),
      ({ responseData: { data } }) => {
        browser.assert.equal(COMPONENT_EQUAL_RESULT_COUNT, data[0].total);
      },
    );
  },

  'A new filter can be created': (browser) => {
    const commonWidget = browser.page.widget.common();
    const createFilterModal = browser.page.modals.common.createFilter();

    browser.page.tables.common()
      .showFiltersList();

    commonWidget.clickAddFilter();

    createFilterModal
      .verifyModalOpened()
      .setFilterTitle('Resource equal value')
      .fillFilterGroups([{
        type: FILTERS_TYPE.OR,
        items: [{
          rule: FILTER_COLUMNS.RESOURCE,
          operator: FILTER_OPERATORS.EQUAL,
          valueType: VALUE_TYPES.STRING,
          value: 'feeder2_1',
        }],
      }])
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.modals.common.filtersList()
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
        browser.assert.equal(COMPONENT_AND_RESOURCE_RESULT_COUNT, data[0].total);
      },
    );

    alarmsWidget.waitFirstAlarmsListXHR(
      () => commonTable.setFiltersType(FILTERS_TYPE.OR),
      ({ responseData: { data } }) => {
        browser.assert.equal(COMPONENT_OR_RESOURCE_RESULT_COUNT, data[0].total);
      },
    );
  },

  'Table widget alarms': (browser) => {
    // const alarmsTable = browser.page.tables.alarms();
    const commonTable = browser.page.tables.common();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();

    const firstId = '6770ba94-51d9-4b8c-ae85-7c62fab18a54';
    // const secondId = 'b7f65652-e53e-4cda-8c3b-be9bbf600ca0';

    browser.page.view()
      .clickMenuViewButton();

    // alarmsTable
    //   .openLiveReporting()
    //   .clickResetLiveReporting();

    commonTable
      .setRowCheckbox(firstId, true)
      .setAllCheckbox(true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.ACK);

    browser.page.modals.alarm.createAckEvent()
      .verifyModalOpened()
      .clickTicketNumber()
      .clearTicketNumber()
      .setTicketNumber(1223333)
      .clickTicketNote()
      .clearTicketNote()
      .setTicketNote('note')
      .setAckTicketResources(true)
      .clickCancelButton()
      .verifyModalClosed();

    commonTable
      .setAllCheckbox(true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ACK);

    browser.page.modals.alarm.createCancelEvent()
      .verifyModalOpened()
      .clickTicketNote()
      .clearTicketNote()
      .setTicketNote('note')
      .clickCancelButton()
      .verifyModalClosed();

    commonTable
      // .clickAlarmListHeaderCell('Connector')
      .clickOnRow(firstId)
      .clickOnSharedAction(firstId, ALARMS_SHARED_ACTIONS.ACK);

    browser.page.modals.alarm.createAckEvent()
      .verifyModalOpened()
      .clickTicketNumber()
      .clearTicketNumber()
      .setTicketNumber(1223333)
      .clickTicketNote()
      .clearTicketNote()
      .setTicketNote('note')
      .setAckTicketResources(true)
      .clickCancelButton()
      .verifyModalClosed();

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
