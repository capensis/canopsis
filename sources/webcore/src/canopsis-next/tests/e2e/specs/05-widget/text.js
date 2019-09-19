// http://nightwatchjs.org/guide#usage

const {
  FILTERS_TYPE,
  VALUE_TYPES,
  INTERVAL_RANGES,
  FILTER_OPERATORS,
  FILTER_COLUMNS,
  WIDGETS_TYPES,
  STAT_TYPES,
  STAT_STATES,
  INTERVAL_PERIODS,
} = require('../../constants');
const { generateTemporaryView, generateTemporaryTextWidget } = require('../../helpers/entities');

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

  'Create widget text with some name': (browser) => {
    const textWidget = {
      ...generateTemporaryTextWidget(),
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      parameters: {
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
        template: 'Template',
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
      .clickWidget(WIDGETS_TYPES.text)
      .verifyModalClosed();

    browser.completed.widget.createText(textWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget text with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 10,
        md: 10,
        lg: 10,
      },
      title: 'Text widget(edited)',
    });

    browser.page.widget.text()
      .clickSubmitText();
  },

  'Delete widget text with some name': (browser) => {
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
