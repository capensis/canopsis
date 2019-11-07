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

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
    browser.globals.defaultViewData = await createWidgetView();

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
          title: 'Filter title',
          groups: [{
            type: FILTERS_TYPE.OR,
            items: [{
              rule: FILTER_COLUMNS.CONNECTOR,
              operator: FILTER_OPERATORS.EQUAL,
              valueType: VALUE_TYPES.STRING,
              value: 'value',
              groups: [{
                type: FILTERS_TYPE.OR,
                items: [{
                  rule: FILTER_COLUMNS.CONNECTOR_NAME,
                  operator: FILTER_OPERATORS.IN,
                  valueType: VALUE_TYPES.BOOLEAN,
                  value: true,
                }],
              }],
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

  'Table widget alarms': (browser) => {
    // const alarmsTable = browser.page.tables.alarms();
    const commonTable = browser.page.tables.common();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();

    const firstId = '6e7fb977-7eab-4e03-9c31-4593a1bf892b';
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

    commonTable
      .clickSearchInput()
      .clearSearchInput()
      .setSearchInput('search string')
      .clickSearchButton()
      .clickSearchResetButton()
      .moveToSearchInformation()
      .setMixFilters(true)
      .selectFilter(1)
      .setFiltersType(FILTERS_TYPE.AND)
      .clickNextPageTopPagination()
      .clickPreviousPageTopPagination()
      .clickNextPageBottomPagination()
      .clickPreviousPageBottomPagination()
      .clickOnPageBottomPagination(2)
      .setItemPerPage(PAGINATION_PER_PAGE_VALUES.FIVE);

    browser.page.view()
      .clickMenuViewButton();
  },

  'Delete widget alarms with some name': (browser) => {
    browser.page.view()
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
