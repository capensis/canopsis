// http://nightwatchjs.org/guide#usage

const { PAGINATION_LIMIT } = require('@/config');
const { WIDGET_TYPES } = require('@/constants');

const {
  ALARMS_MASS_ACTIONS,
  WEEK_DAYS,
  DATE_INTERVAL_MINUTES,
  PERIODICAL_BEHAVIOR_REASONS,
  PERIODICAL_BEHAVIOR_FREQUENCY,
  PBEHAVIOR_STEPS,
} = require('../../../constants');
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

  'All elements in the table can be selected': (browser) => {
    browser.page.tables.common()
      .setAllCheckbox(true)
      .assertActiveCheckboxCount(PAGINATION_LIMIT)
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

  'Pressing on button "Periodical behavior" creates periodical behavior': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const dateIntervalField = browser.page.fields.dateInterval();
    const pbehaviorForm = browser.page.forms.pbehavior();
    const commonTable = browser.page.tables.common();
    const createPbehaviorModal = browser.page.modals.alarm.createPbehavior();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.PERIODICAL_BEHAVIOR);

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

  'An acknowledge without a ticket can be created': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createAckEventModal = browser.page.modals.alarm.createAckEvent();
    const confirmAckWithTicketModal = browser.page.modals.alarm.confirmAckWithTicket();
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

  'An acknowledge with ticket can be created': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
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
      .setAckTicketResources(true);

    alarmsWidget.waitFirstEventXHR(
      () => createAckEventModal
        .clickSubmitButtonWithTicket()
        .verifyModalClosed(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createAckEventModal.verifyModalClosed();
      },
    );
  },

  '"Fast ack" can be created': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    alarmsWidget.waitFirstEventXHR(
      () => commonTable
        .setRowCheckbox(firstAlarm._id, true)
        .clickOnMassAction(ALARMS_MASS_ACTIONS.FAST_ACK),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
      },
    );
  },

  'An "ack" can be canceled': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createCancelEventModal = browser.page.modals.alarm.createCancelEvent();
    const [firstAlarm, secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .setRowCheckbox(secondAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ACK);

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

  'An alarm can be canceled': (browser) => {
    const alarmsWidget = browser.page.widget.alarms();
    const commonTable = browser.page.tables.common();
    const createCancelEventModal = browser.page.modals.alarm.createCancelEvent();
    const [firstAlarm, secondAlarm] = browser.globals.temporary.alarmsList;

    commonTable
      .setRowCheckbox(firstAlarm._id, true)
      .setRowCheckbox(secondAlarm._id, true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ALARM);

    createCancelEventModal
      .verifyModalOpened()
      .clearTicketNote()
      .setTicketNote('First and second alarm cancel');

    alarmsWidget.waitFirstEventXHR(
      () => createCancelEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createCancelEventModal.verifyModalClosed();
      },
    );
  },
};
