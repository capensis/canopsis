// http://nightwatchjs.org/guide#usage

const {
  INFO_POPUP_DEFAULT_COLUMNS,
  ALARMS_WIDGET_SORT_FIELD,
  SORT_ORDERS,
  PAGINATION_PER_PAGE_VALUES,
  FILTERS_TYPE,
  VALUE_TYPES,
  INTERVAL_RANGES,
  FILTER_OPERATORS,
  ALARMS_MASS_ACTIONS,
  FILTER_COLUMNS,
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, removeWidgetView } = require('../../helpers/api');
const { generateTemporaryAlarmsWidget } = require('../../helpers/entities');

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

    delete browser.globals.defaultViewData;
    delete browser.globals.temporary;

    done();
  },

  'Create widget alarms with some name': (browser) => {
    const alarmsWidget = {
      ...generateTemporaryAlarmsWidget(),
      periodicRefresh: 140,
      parameters: {
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: ALARMS_WIDGET_SORT_FIELD.component,
        },
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.TWENTY,
        openedResolvedFilter: {
          open: true,
          resolve: true,
        },
        infoPopups: [{
          column: INFO_POPUP_DEFAULT_COLUMNS.connectorName,
          template: 'Info popup template',
        }],
        ack: {
          isAckNoteRequired: true,
          isMultiAckEnabled: true,
          fastAckOutput: {
            enabled: true,
            output: 'Output',
          },
        },
        moreInfos: 'More infos popup',
        enableHtml: true,
        liveReporting: {
          // calendarStartDate: {
        //     minute: 0,
        //     hour: 12,
        //     day: 12,
        //   },
        //   endDate: '13/09/2019 00:00',
          range: INTERVAL_RANGES.LAST_YEAR,
        },
        filters: {
          isMix: true,
          type: FILTERS_TYPE.OR,
          title: 'Filter title',
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

    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .clickPanelHeader(groupId)
      .clickLinkView(viewId);

    browser.page.view()
      .clickMenuViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget(WIDGET_TYPES.alarmList)
      .verifyModalClosed();

    browser.completed.widget.createAlarmsList(alarmsWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Table widget alarms': (browser) => {
    const alarmsTable = browser.page.tables.alarms();

    browser.page.view()
      .clickMenuViewButton();


    alarmsTable
      .clickAlarmListHeaderCell('Connector')
      .setAllCheckbox(true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.ACK);

    browser.page.modals.alarm.createAckEvent()
      .verifyModalOpened()
      .clickTicketNumber()
      .clearTicketNumber()
      .setTicketNumber(1223333)
      .clickTicketNote()
      .clearTicketNote()
      .setTicketNote('note')
      .setAckTicketResources(true)
      .clickCancelButton()
      .verifyModalClosed();

    alarmsTable
      .setAllCheckbox(true)
      .clickOnMassAction(ALARMS_MASS_ACTIONS.CANCEL_ACK);

    browser.page.modals.alarm.createCancelEvent()
      .verifyModalOpened()
      .clickTicketNote()
      .clearTicketNote()
      .setTicketNote('note')
      .clickCancelButton()
      .verifyModalClosed();

    alarmsTable
      .clickSearchInput()
      .clearSearchInput()
      .setSearchInput('search string')
      .clickSearchButton()
      .clickSearchResetButton()
      .moveToSearchInformation()
      .setMixFilters(true)
      .selectFilter(1)
      .setFiltersType(FILTERS_TYPE.AND)
      .clickNextPageTopPagination()
      .clickPreviousPageTopPagination()
      .clickNextPageBottomPagination()
      .clickPreviousPageBottomPagination()
      .clickOnPageBottomPagination(2)
      .setItemPerPage(PAGINATION_PER_PAGE_VALUES.FIVE);

    browser.page.view()
      .clickMenuViewButton();
  },

  'Delete widget alarms with some name': (browser) => {
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
