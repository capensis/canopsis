// http://nightwatchjs.org/guide#usage

import { ROW_SIZE_CLASSES, ROW_SIZE_KEYS } from '../../../constants';

const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, createWidgetForView, removeWidgetView } = require('../../../helpers/api');

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
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

    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .clickPanelHeader(groupId)
      .clickLinkView(viewId);

    browser.page.view()
      .clickMenuViewButton()
      .clickEditViewButton();
  },

  'Widget is size can be changed in mobile version': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const rowSize = 5;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickRowGridSize()
      .setRowSize(ROW_SIZE_KEYS.SMARTPHONE, rowSize)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonWidget.assertWidgetRowClasses(widgetId, `${ROW_SIZE_CLASSES.SMARTPHONE}${rowSize}`);
        },
      );
  },

  'Widget is size can be changed in tablet version': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const rowSize = 7;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickRowGridSize()
      .setRowSize(ROW_SIZE_KEYS.TABLET, rowSize)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonWidget.assertWidgetRowClasses(widgetId, `${ROW_SIZE_CLASSES.TABLET}${rowSize}`);
        },
      );
  },

  'Widget is size can be changed in desktop version': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const rowSize = 6;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickRowGridSize()
      .setRowSize(ROW_SIZE_KEYS.DESKTOP, rowSize)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonWidget.assertWidgetRowClasses(widgetId, `${ROW_SIZE_CLASSES.DESKTOP}${rowSize}`);
        },
      );
  },

  'Widget is title can be changed': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const widgetTitle = 'Alarm title(changed)';

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickWidgetTitle()
      .clearWidgetTitle()
      .setWidgetTitle(widgetTitle)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
          commonWidget.getWidgetTitle(widgetId, ({ value }) => {
            browser.assert.equal(widgetTitle, value);
          });
        },
      );
  },

  'The periodic refresh can be set manually': (browser) => {
    const { widgetId } = browser.globals.defaultViewData;
    const commonWidget = browser.page.widget.common();
    const alarmsWidget = browser.page.widget.alarms();
    const refreshTime = 180;

    browser.page.view()
      .clickEditingMenu(widgetId)
      .clickEditWidgetButton(widgetId);

    commonWidget
      .clickPeriodicRefresh()
      .setPeriodicRefreshSwitch(true)
      .clickPeriodicRefreshField()
      .clearPeriodicRefreshField()
      .setPeriodicRefreshField(refreshTime)
      .waitFirstUserPreferencesXHR(
        () => alarmsWidget.clickSubmitAlarms(),
        ({ responseData }) => {
          browser.assert.equal(responseData.success, true);
        },
      );
  },
};
