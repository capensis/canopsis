// http://nightwatchjs.org/guide#usage

const {
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SERVICE_WEATHER_WIDGET_SORT_FIELD,
  INFO_POPUP_DEFAULT_COLUMNS,
  PAGINATION_PER_PAGE_VALUES,
  FILTERS_TYPE,
  FILTER_OPERATORS,
  FILTER_COLUMNS,
  VALUE_TYPES,
  SORT_ORDERS,
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, removeWidgetView } = require('../../helpers/api');
const { generateTemporaryWeatherWidget } = require('../../helpers/entities');

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

  'Create widget weather with some name': (browser) => {
    const weatherWidget = {
      ...generateTemporaryWeatherWidget(),
      periodicRefresh: 140,
      parameters: {
        advanced: true,
        alarmsList: true,
        limit: 140,
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: SERVICE_WEATHER_WIDGET_SORT_FIELD.status,
        },
        margin: {
          top: 3,
          right: 3,
          bottom: 3,
          left: 3,
        },
        infoPopups: [{
          column: INFO_POPUP_DEFAULT_COLUMNS.connectorName,
          template: 'Info popup template',
        }],
        elementPerPage: PAGINATION_PER_PAGE_VALUES.HUNDRED,
        columnSM: 12,
        columnMD: 12,
        columnLG: 12,
        heightFactor: 20,
        modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList,
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
        moreInfos: 'More infos popup',
        blockTemplate: 'Template weather item text',
        modalTemplate: 'Template modal text',
        newColumnNames: [{
          index: 9,
          data: {
            value: 'alarm.v.connector',
            label: 'New column',
            isHtml: true,
          },
        }],
        editColumnNames: [{
          index: 1,
          data: {
            value: 'alarm.v.changeConnector',
            label: 'Connector(changed)',
            isHtml: true,
          },
        }],
        moveColumnNames: [{
          index: 1,
          down: true,
        }, {
          index: 2,
          up: true,
        }],
        deleteColumnNames: [2],
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
      .clickWidget(WIDGET_TYPES.weather)
      .verifyModalClosed();

    browser.completed.widget.createServiceWeather(weatherWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget weather with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 10,
        md: 10,
        lg: 10,
      },
      title: 'Weather widget(edited)',
      parameters: {
        advanced: true,
        limit: 180,
        sort: {
          order: SORT_ORDERS.asc,
          orderBy: SERVICE_WEATHER_WIDGET_SORT_FIELD.criticity,
        },
      },
    });

    browser.page.widget.weather()
      .clickSubmitWeather();
  },

  'Delete widget weather with some name': (browser) => {
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
