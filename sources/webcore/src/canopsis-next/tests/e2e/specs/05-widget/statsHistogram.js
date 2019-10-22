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
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, removeWidgetView } = require('../../helpers/api');
const { generateTemporaryStatsHistogramWidget } = require('../../helpers/entities');

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

    delete browser.globals.credentials;
    delete browser.globals.temporary;

    done();
  },

  'Create widget stats histogram with some name': (browser) => {
    const statsHistogram = {
      ...generateTemporaryStatsHistogramWidget(),
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
          periodValue: 2,
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
        statsSelector: {
          newStats: [{
            type: STAT_TYPES.RESOLVED_TIME_SLA,
            title: 'title',
            recursive: false,
            states: [
              { index: STAT_STATES.OK, checked: true },
              { index: STAT_STATES.MINOR, checked: true },
            ],
            authors: ['first', 'second'],
            sla: '<=2',
          }],
        },
        statsColors: [{
          title: 'title',
          color: '#666666',
        }],
        annotationLine: {
          isEnabled: true,
          value: 12,
          label: 'Label',
          labelColor: '#111111',
          lineColor: '#444444',
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
      .clickWidget(WIDGET_TYPES.statsHistogram)
      .verifyModalClosed();

    browser.completed.widget.createStatsHistogram(statsHistogram, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget stats histogram with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 6,
        md: 6,
        lg: 6,
      },
      title: 'Stats histogram widget(edited)',
    });

    browser.page.widget.statsHistogram()
      .clickSubmitStatsHistogram();
  },

  'Delete widget stats histogram with some name': (browser) => {
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
