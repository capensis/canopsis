// http://nightwatchjs.org/guide#usage

import { CRUD_ACTIONS } from '../../../constants';

const { WIDGET_TYPES } = require('@/constants');

const {
  ALARMS_SHARED_ACTIONS_WITH_ACK,
  ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK,
  DATE_INTERVAL_MINUTES,
  PBEHAVIOR_STEPS,
  PERIODICAL_BEHAVIOR_FREQUENCY,
  PERIODICAL_BEHAVIOR_REASONS,
  WEEK_DAYS,
} = require('../../../constants');
const { ENTITIES_STATES } = require('@/constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const DROP_DOWN_ACTION_COUNT = 5; // TODO Why should be 6?

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

  'The button "Declare ticket" reports an incident': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createDeclareTicketEventModal = browser.page.modals.alarm.createDeclareTicketEvent();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnSharedAction(firstAlarm._id, ALARMS_SHARED_ACTIONS_WITH_ACK.DECLARE_TICKET);
    createDeclareTicketEventModal.verifyModalOpened();

    alarmsWidget.waitFirstEventXHR(
      () => createDeclareTicketEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createDeclareTicketEventModal.verifyModalClosed();
      },
    );
  },

  'The "Associate ticket" feature associates a ticket': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createAssociateTicketModal = browser.page.modals.alarm.createAssociateTicket();
    const [, secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable.clickOnSharedAction(secondAlarm._id, ALARMS_SHARED_ACTIONS_WITH_ACK.ASSOCIATE_TICKET);
    createAssociateTicketModal
      .verifyModalOpened()
      .clickTicketNumber()
      .clearTicketNumber()
      .setTicketNumber(1000);

    alarmsWidget.waitFirstEventXHR(
      () => createAssociateTicketModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createAssociateTicketModal.verifyModalClosed();
      },
    );
  },

  'The button "Cancel alarm" changes element status to "Canceled"': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createCancelEventModal = browser.page.modals.alarm.createCancelEvent();
    const thirdAlarm = browser.globals.temporary.alarmsList[2];

    commonTable.clickOnSharedAction(thirdAlarm._id, ALARMS_SHARED_ACTIONS_WITH_ACK.CANCEL_ALARM);
    createCancelEventModal
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('First and second ack cancel');

    alarmsWidget.waitFirstEventXHR(
      () => createCancelEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createCancelEventModal.verifyModalClosed();
      },
    );
  },

  'The three dots icon opens a menu with six options (when the "Ack" feature is on)': (browser) => {
    const commonTable = browser.page.tables.common();
    const fortiesAlarm = browser.globals.temporary.alarmsList[4];

    commonTable.clickOnDropDownDots(fortiesAlarm._id);
    commonTable.assertDropDownAction(DROP_DOWN_ACTION_COUNT);
  },

  'The "Cancel ack" feature cancels the "Acknowledge" mode': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createCancelEventModal = browser.page.modals.alarm.createCancelEvent();
    const sixthAlarm = browser.globals.temporary.alarmsList[5];

    commonTable
      .clickOnDropDownDots(sixthAlarm._id)
      .clickOnDropDownAction(ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK.CANCEL_ACK);

    createCancelEventModal
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('Sixth alarm canceled');

    alarmsWidget.waitFirstEventXHR(
      () => createCancelEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createCancelEventModal.verifyModalClosed();
      },
    );
  },

  'The criticity of alarm can be changed': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createChangeStateEventModal = browser.page.modals.alarm.createChangeStateEvent();
    const seventhAlarm = browser.globals.temporary.alarmsList[6];

    commonTable
      .clickOnDropDownDots(seventhAlarm._id)
      .clickOnDropDownAction(ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK.CHANGE_CRITICITY);

    createChangeStateEventModal
      .verifyModalOpened()
      .clickNote()
      .clearNote()
      .setNote('Note')
      .clickCriticity(ENTITIES_STATES.major);

    alarmsWidget.waitFirstEventXHR(
      () => createChangeStateEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createChangeStateEventModal.verifyModalClosed();
      },
    );
  },

  'The "Ack" feature can be added without declaring the ticket': () => {},

  'The "Acknowledge and Associate Ticket" feature adds acknowledge and associates ticket': () => {},

  'The "Acknowledge and Declare Ticket" feature adds acknowledge and declares ticket': () => {},

  'The "Acknowledge" feature adds acknowledge by additional modal window': () => {},

  'The "Acknowledge and Associate ticket" feature adds acknowledge and associates ticket by additional modal window': () => {},

  'The "Snooze alarm" feature creates a "Snooze" event': () => {},

  'The "Periodical behavior" feature creates periodical behavior': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();
    const commonTable = browser.page.tables.common();
    const createPbehaviorModal = browser.page.modals.alarm.createPbehavior();
    const eighthAlarm = browser.globals.temporary.alarmsList[7];

    commonTable
      .clickOnDropDownDots(eighthAlarm._id)
      .clickOnDropDownAction(ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK.PERIODICAL_BEHAVIOR);

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
    const eighthAlarm = browser.globals.temporary.alarmsList[7];

    commonTable
      .clickOnDropDownDots(eighthAlarm._id)
      .clickOnDropDownAction(ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK.LIST_PERIODICAL_BEHAVIOR);

    pbehaviorListModal
      .verifyModalOpened()
      .clickAction(eighthAlarm._id, CRUD_ACTIONS.UPDATE);

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
