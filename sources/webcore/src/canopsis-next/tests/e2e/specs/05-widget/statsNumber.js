// http://nightwatchjs.org/guide#usage

const {
  FILTERS_TYPE,
  VALUE_TYPES,
  FILTER_OPERATORS,
  FILTER_COLUMNS,
  INTERVAL_RANGES,
  STAT_TYPES,
  STAT_STATES,
  INTERVAL_PERIODS,
  PAGINATION_PER_PAGE_VALUES,
  SORT_ORDERS,
  STATS_DISPLAY_MODE,
} = require('../../constants');
const { WIDGET_TYPES, STATS_CRITICITY } = require('@/constants');
const { generateTemporaryView, generateTemporaryStatsNumberWidget } = require('../../helpers/entities');

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

  'Create widget stats number with some name': (browser) => {
    const statsNumber = {
      ...generateTemporaryStatsNumberWidget(),
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      parameters: {
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
        filter: {
          groups: [{
            type: FILTERS_TYPE.OR,
            items: [{
              rule: FILTER_COLUMNS.NAME,
              operator: FILTER_OPERATORS.EQUAL,
              valueType: VALUE_TYPES.STRING,
              value: 'value',
              groups: [{
                type: FILTERS_TYPE.OR,
                items: [{
                  rule: FILTER_COLUMNS.NAME,
                  operator: FILTER_OPERATORS.IN,
                  valueType: VALUE_TYPES.BOOLEAN,
                  value: true,
                }],
              }],
            }, {
              type: FILTERS_TYPE.AND,
              rule: FILTER_COLUMNS.TYPE,
              operator: FILTER_OPERATORS.NOT_EQUAL,
              valueType: VALUE_TYPES.NUMBER,
              value: 136,
            }],
          }],
        },
        statSelector: {
          type: STAT_TYPES.RESOLVED_TIME_SLA,
          title: 'title',
          recursive: false,
          states: [
            { index: STAT_STATES.OK, checked: true },
            { index: STAT_STATES.MINOR, checked: true },
          ],
          authors: ['first', 'second'],
          sla: '<=2',
        },
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.HUNDRED,
        sortOrder: SORT_ORDERS.asc,
        displayMode: {
          type: STATS_DISPLAY_MODE.VALUE,
          parameters: {
            [STATS_CRITICITY.ok]: {
              value: 1,
              color: '#111111',
            },
            [STATS_CRITICITY.minor]: {
              value: 2,
              color: '#444444',
            },
            [STATS_CRITICITY.major]: {
              value: 2,
              color: '#666666',
            },
            [STATS_CRITICITY.critical]: {
              value: 2,
              color: '#ffffff',
            },
          },
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
      .clickWidget(WIDGET_TYPES.statsNumber)
      .verifyModalClosed();

    browser.completed.widget.createStatsNumber(statsNumber, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget stats number with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 6,
        md: 6,
        lg: 6,
      },
      title: 'Stats number widget(edited)',
    });

    browser.page.widget.statsNumber()
      .clickSubmitStatsNumber();
  },

  'Delete widget stats number with some name': (browser) => {
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
