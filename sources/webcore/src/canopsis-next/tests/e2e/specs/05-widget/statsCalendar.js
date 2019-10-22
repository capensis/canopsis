// http://nightwatchjs.org/guide#usage

const {
  FILTERS_TYPE,
  VALUE_TYPES,
  FILTER_OPERATORS,
  CONTEXT_FILTER_COLUMNS,
  PAGINATION_PER_PAGE_VALUES,
  INFO_POPUP_DEFAULT_COLUMNS,
} = require('../../constants');
const { WIDGET_TYPES, STATS_CRITICITY } = require('@/constants');
const { createWidgetView, removeWidgetView } = require('../../helpers/api');
const { generateTemporaryStatsCalendarWidget } = require('../../helpers/entities');

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

  'Create widget stats calendar with some name': (browser) => {
    const statsCalendar = {
      ...generateTemporaryStatsCalendarWidget(),
      parameters: {
        newColumnNames: [{
          index: 9,
          data: {
            value: 'connector',
            label: 'New column',
          },
        }],
        editColumnNames: [{
          index: 1,
          data: {
            value: 'connector',
            label: 'Connector(changed)',
          },
        }],
        moveColumnNames: [{
          index: 1,
          down: true,
        }, {
          index: 2,
          up: true,
        }],
        infoPopups: [{
          column: INFO_POPUP_DEFAULT_COLUMNS.connectorName,
          template: 'Info popup template',
        }],
        openedResolvedFilter: {
          open: true,
          resolve: true,
        },
        deleteColumnNames: [2],
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.HUNDRED,
        moreInfos: 'More infos popup',
        criticityLevels: {
          minor: 20,
          major: 30,
          critical: 40,
        },
        colorsSelector: {
          [STATS_CRITICITY.ok]: '#111111',
          [STATS_CRITICITY.minor]: '#444444',
          [STATS_CRITICITY.major]: '#666666',
          [STATS_CRITICITY.critical]: '#ffffff',
        },
        considerPbehaviors: false,
        filters: {
          title: 'Filter title',
          groups: [{
            type: FILTERS_TYPE.OR,
            items: [{
              rule: CONTEXT_FILTER_COLUMNS.NAME,
              operator: FILTER_OPERATORS.EQUAL,
              valueType: VALUE_TYPES.STRING,
              value: 'value',
              groups: [{
                type: FILTERS_TYPE.OR,
                items: [{
                  rule: CONTEXT_FILTER_COLUMNS.NAME,
                  operator: FILTER_OPERATORS.IN,
                  valueType: VALUE_TYPES.BOOLEAN,
                  value: true,
                }],
              }],
            }, {
              type: FILTERS_TYPE.AND,
              rule: CONTEXT_FILTER_COLUMNS.TYPE,
              operator: FILTER_OPERATORS.NOT_EQUAL,
              valueType: VALUE_TYPES.NUMBER,
              value: 136,
            }],
          }],
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
      .clickWidget(WIDGET_TYPES.statsCalendar)
      .verifyModalClosed();

    browser.completed.widget.createStatsCalendar(statsCalendar, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget stats calendar with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 6,
        md: 6,
        lg: 6,
      },
      title: 'Stats calendar widget(edited)',
      parameters: {},
    });

    browser.page.widget.statsCalendar()
      .clickSubmitStatsCalendar();
  },

  'Delete widget stats calendar with some name': (browser) => {
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
