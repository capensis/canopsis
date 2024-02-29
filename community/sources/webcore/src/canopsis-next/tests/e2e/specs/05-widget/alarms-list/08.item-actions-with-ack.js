// http://nightwatchjs.org/guide#usage

const { WIDGET_TYPES } = require('@/constants');

const {
  ALARMS_SHARED_ACTIONS_WITH_ACK,
  ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK,
} = require('../../../constants');
const { ALARM_STATES } = require('@/constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

const DROP_DOWN_ACTION_COUNT = 5;

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
      .clickCriticity(ALARM_STATES.major);

    alarmsWidget.waitFirstEventXHR(
      () => createChangeStateEventModal.clickSubmitButton(),
      ({ responseData }) => {
        browser.assert.equal(responseData.success, true);
        createChangeStateEventModal.verifyModalClosed();
      },
    );
  },
};
