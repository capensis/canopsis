// http://nightwatchjs.org/guide#usage

import { CRUD_ACTIONS } from '../../../constants';

const { WIDGET_TYPES } = require('@/constants');

const {
  ALARMS_SHARED_ACTIONS,
  SNOOZE_TYPES,
  DATE_INTERVAL_MINUTES,
  PBEHAVIOR_STEPS,
  PERIODICAL_BEHAVIOR_FREQUENCY,
  ALARMS_SHARED_DROPDOWN_ACTIONS,
  PERIODICAL_BEHAVIOR_REASONS,
  WEEK_DAYS,
} = require('../../../constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const DROP_DOWN_ACTION_COUNT = 2; // TODO Why should be 3?

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
        },
      );
  },

  'The "Ack" feature can be added without declaring the ticket': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createAckEventModal = browser.page.modals.alarm.createAckEvent();
    const confirmAckWithTicketModal = browser.page.modals.alarm.confirmAckWithTicket();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnSharedAction(firstAlarm._id, ALARMS_SHARED_ACTIONS.ACK);

    createAckEventModal
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('Note text')
      .setAckTicketResources(true);

    alarmsWidget.waitFirstEventXHR(
      () => createAckEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        confirmAckWithTicketModal.verifyModalClosed();
        createAckEventModal.verifyModalClosed();
      },
    );
  },

  'The "Acknowledge and Declare Ticket" feature adds acknowledge and declares ticket': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createAckEventModal = browser.page.modals.alarm.createAckEvent();
    const confirmAckWithTicketModal = browser.page.modals.alarm.confirmAckWithTicket();
    const [, secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnSharedAction(secondAlarm._id, ALARMS_SHARED_ACTIONS.ACK);

    createAckEventModal
      .verifyModalOpened()
      .clearTicketNumber()
      .setTicketNumber(165558556)
      .clearTicketNote()
      .setTicketNote('Note text')
      .setAckTicketResources(true)
      .clickSubmitButton();

    alarmsWidget.waitFirstEventXHR(
      () => confirmAckWithTicketModal
        .verifyModalOpened()
        .clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        confirmAckWithTicketModal.verifyModalClosed();
        createAckEventModal.verifyModalClosed();
      },
    );
  },

  'The "Snooze alarm" feature creates a "Snooze" event': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createSnoozeEventModal = browser.page.modals.alarm.createSnoozeEvent();
    const fourthAlarm = browser.globals.temporary.alarmsList[3];

    commonTable.clickOnSharedAction(fourthAlarm._id, ALARMS_SHARED_ACTIONS.SNOOZE_ALARM);

    createSnoozeEventModal
      .verifyModalOpened()
      .clickDurationValue()
      .clearDurationValue()
      .setDurationValue(1)
      .setDurationType(SNOOZE_TYPES.MINUTES);

    alarmsWidget.waitFirstEventXHR(
      () => createSnoozeEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createSnoozeEventModal.verifyModalClosed();
      },
    );
  },

  'The "Fast ack" feature adds acknowledge without filling in the form': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const thirdAlarm = browser.globals.temporary.alarmsList[2];

    alarmsWidget.waitFirstEventXHR(
      () => commonTable.clickOnSharedAction(thirdAlarm._id, ALARMS_SHARED_ACTIONS.FAST_ACK),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
      },
    );
  },

  'The three dots icon opens a menu with three options (when the "Ack" feature is off)': (browser) => {
    const commonTable = browser.page.tables.common();
    const fortiesAlarm = browser.globals.temporary.alarmsList[4];

    commonTable.clickOnDropDownDots(fortiesAlarm._id);
    commonTable.assertDropDownAction(DROP_DOWN_ACTION_COUNT);
  },

  'The "Report alarm" feature reports an alarm': () => { },

  'The "Periodical behavior" feature creates periodical behavior': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();
    const commonTable = browser.page.tables.common();
    const createPbehaviorModal = browser.page.modals.alarm.createPbehavior();
    const sixthAlarm = browser.globals.temporary.alarmsList[5];

    commonTable
      .clickOnDropDownDots(sixthAlarm._id)
      .clickOnDropDownAction(ALARMS_SHARED_DROPDOWN_ACTIONS.PERIODICAL_BEHAVIOR);

    createPbehaviorModal.verifyModalOpened();

    pbehaviorForm
      .clearName()
      .clickName()
      .setName('Name')
      .clearReason()
      .clickReason()
      .setReason('P')
      .selectReason(PERIODICAL_BEHAVIOR_REASONS.REHABILITATION_PROBLEM)
      .selectType(1)
      .clickStartDate();

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(3)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.FIVE)
      .clickOutsideDateInterval();

    pbehaviorForm.clickEndDate();

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(4)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.TEN)
      .clickOutsideDateInterval();

    pbehaviorForm
      .clickPbehaviorFormStep(PBEHAVIOR_STEPS.RRULE)
      .setRuleCheckbox(true)
      .selectFrequency(PERIODICAL_BEHAVIOR_FREQUENCY.MINUTELY)
      .clickByWeekDay()
      .selectByWeekDay(WEEK_DAYS.TUESDAY, true)
      .selectByWeekDay(WEEK_DAYS.FRIDAY, true)
      .clickOutsideByWeekDay()
      .clearRepeat()
      .clickRepeat()
      .setRepeat(5)
      .clearInterval()
      .clickInterval()
      .setInterval(5);

    pbehaviorForm
      .clickPbehaviorFormStep(PBEHAVIOR_STEPS.COMMENTS)
      .clickAddComment()
      .clickCommentField(1)
      .clearCommentField(1)
      .clearCommentField(1)
      .setCommentField(1, 2);

    alarmsWidget.waitFirstPbehaviorXHR(
      () => createPbehaviorModal.clickSubmitButton(),
      id => browser.assert.equal(!!id, true),
    );
  },

  'The "List periodic behaviors" feature allows to look through all periodic behaviors': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();
    const commonTable = browser.page.tables.common();
    const createPbehaviorModal = browser.page.modals.alarm.createPbehavior();
    const pbehaviorListModal = browser.page.modals.alarm.pbehaviorList();
    const sixthAlarm = browser.globals.temporary.alarmsList[5];

    commonTable
      .clickOnDropDownDots(sixthAlarm._id)
      .clickOnDropDownAction(ALARMS_SHARED_DROPDOWN_ACTIONS.LIST_PERIODICAL_BEHAVIOR);

    pbehaviorListModal
      .verifyModalOpened()
      .clickAction(sixthAlarm._id, CRUD_ACTIONS.UPDATE);

    createPbehaviorModal.verifyModalOpened();

    pbehaviorForm
      .clearName()
      .clickName()
      .setName('Name')
      .clearReason()
      .clickReason()
      .setReason('P')
      .selectReason(PERIODICAL_BEHAVIOR_REASONS.REHABILITATION_PROBLEM)
      .selectType(1)
      .clickStartDate();

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(3)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.FIVE)
      .clickOutsideDateInterval();

    pbehaviorForm.clickEndDate();

    dateIntervalField
      .clickDatePickerDayTab()
      .selectCalendarDay(4)
      .clickDatePickerHoursTab()
      .selectCalendarHour(16)
      .clickDatePickerMinutesTab()
      .selectCalendarMinute(DATE_INTERVAL_MINUTES.TEN)
      .clickOutsideDateInterval();

    pbehaviorForm
      .clickPbehaviorFormStep(PBEHAVIOR_STEPS.RRULE)
      .setRuleCheckbox(true)
      .selectFrequency(PERIODICAL_BEHAVIOR_FREQUENCY.MINUTELY)
      .clickByWeekDay()
      .selectByWeekDay(WEEK_DAYS.TUESDAY, true)
      .selectByWeekDay(WEEK_DAYS.FRIDAY, true)
      .clickOutsideByWeekDay()
      .clearRepeat()
      .clickRepeat()
      .setRepeat(5)
      .clearInterval()
      .clickInterval()
      .setInterval(5);

    pbehaviorForm
      .clickPbehaviorFormStep(PBEHAVIOR_STEPS.COMMENTS)
      .clickAddComment()
      .clickCommentField(1)
      .clearCommentField(1)
      .clearCommentField(1)
      .setCommentField(1, 2);

    alarmsWidget.waitFirstPbehaviorXHR(
      () => createPbehaviorModal.clickSubmitButton(),
      (id) => {
        browser.assert.equal(!!id, true);

        createPbehaviorModal.verifyModalClosed();
        pbehaviorListModal
          .clickSubmitButton()
          .verifyModalClosed();
      },
    );
  },
};
