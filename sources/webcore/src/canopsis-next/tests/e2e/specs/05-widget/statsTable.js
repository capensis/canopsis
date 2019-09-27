// http://nightwatchjs.org/guide#usage

const {
  SORT_ORDERS,
  FILTERS_TYPE,
  VALUE_TYPES,
  INTERVAL_RANGES,
  INTERVAL_PERIODS,
  FILTER_OPERATORS,
  FILTER_COLUMNS,
  STAT_TYPES,
  STAT_STATES,
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { generateTemporaryView, generateTemporaryStatsTableWidget } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.completed.logout()
      .end(done);

    delete browser.globals.temporary;
  },

  'Create test view': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.defaultViewData = {
        viewId: view._id,
        groupId: view.group_id,
      };
    });
  },

  'Create widget stats table with some name': (browser) => {
    const statsTableWidget = {
      ...generateTemporaryStatsTableWidget(),
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      parameters: {
        advanced: true,
        dateInterval: {
          calendarStartDate: {
            minute: 0,
            hour: 12,
            day: 12,
          },
          endDate: '13/09/2019 00:00',
          range: INTERVAL_RANGES.CUSTOM,
          period: INTERVAL_PERIODS.HOUR,
        },
        statsSelector: {
          newStats: [{
            type: STAT_TYPES.RESOLVED_TIME_SLA,
            title: 'title',
            trend: false,
            recursive: false,
            states: [
              { index: STAT_STATES.OK, checked: true },
              { index: STAT_STATES.MINOR, checked: true },
            ],
            authors: ['first', 'second'],
            sla: '<=2',
          }],
        },
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: 1,
        },
        filter: {
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
    const view = browser.page.view();
    const groupsSideBar = browser.page.layout.groupsSideBar();

    groupsSideBar.clickPanelHeader(groupId)
      .clickLinkView(viewId);

    view.clickMenuViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget(WIDGET_TYPES.statsTable)
      .verifyModalClosed();

    browser.completed.widget.createStatsTable(statsTableWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget stats table with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 10,
        md: 10,
        lg: 10,
      },
      title: 'Stats table widget(edited)',
    });

    browser.page.widget.statsTable()
      .clickSubmitStatsTable();
  },

  'Delete widget stats table with some name': (browser) => {
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

  'Delete test view': (browser) => {
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.completed.view.delete(groupId, viewId);
    browser.completed.view.deleteGroup(groupId);
  },
};
