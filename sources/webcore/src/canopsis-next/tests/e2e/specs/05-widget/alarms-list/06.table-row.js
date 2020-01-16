// http://nightwatchjs.org/guide#usage

const { WIDGET_TYPES } = require('@/constants');

const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

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
      parameters: {
        infoPopups: [{
          column: 'v.connector',
          template: 'Info popup',
        }],
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

  'Information pop-up can be shown': (browser) => {
    const alarmsTable = browser.page.tables.alarms();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    const column = 'Connector';

    alarmsTable
      .clickOnRowInfoPopupOpenButton(firstAlarm._id, column)
      .verifyRowInfoPopupVisible(firstAlarm._id, column);
  },

  'Information pop-up can be closed': (browser) => {
    const alarmsTable = browser.page.tables.alarms();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    const column = 'Connector';

    alarmsTable
      .clickOnRowInfoPopupCloseButton(firstAlarm._id, column)
      .verifyRowInfoPopupDeleted(firstAlarm._id, column);
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

  'Placing a cursor on signs in the column "Extra details" makes information pop-up show': (browser) => {
    const alarmsTable = browser.page.tables.alarms();
    const [firstAlarm] = browser.globals.temporary.alarmsList;

    const status = {
      [firstAlarm.pbehaviors.some(pbehavior => pbehavior.isActive)]: 'pbehaviors',
      [!!firstAlarm.v.canceled]: 'canceled',
      [!!firstAlarm.v.ticket]: 'ticket',
      [!!firstAlarm.v.ack]: 'ack',
    }.true;

    alarmsTable
      .moveToExtraDetailsOpenButton(firstAlarm._id, status)
      .verifyRowExtraDetailsVisible(firstAlarm._id)
      .moveOutsideToExtraDetailsOpenButton(firstAlarm._id, status)
      .verifyRowExtraDetailsDeleted(firstAlarm._id);
  },
};
