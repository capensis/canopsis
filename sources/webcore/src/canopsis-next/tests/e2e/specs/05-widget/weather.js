// http://nightwatchjs.org/guide#usage

const {
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SERVICE_WEATHER_WIDGET_SORT_FIELD,
  PAGINATION_PER_PAGE_VALUES,
  FILTERS_TYPE,
  FILTER_OPERATORS,
  FILTER_COLUMNS,
  VALUE_TYPES,
  SORT_ORDERS,
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { generateTemporaryView, generateTemporaryWeatherWidget } = require('../../helpers/entities');

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

  'Create widget weather with some name': (browser) => {
    const weatherWidget = {
      ...generateTemporaryWeatherWidget(),
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
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
    const view = browser.page.view();
    const groupsSideBar = browser.page.layout.groupsSideBar();

    groupsSideBar.clickPanelHeader(groupId)
      .clickLinkView(viewId);

    view.clickMenuViewButton()
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
    const weatherWidget = {
      parameters: {
        limit: 180,
        sort: {
          order: SORT_ORDERS.asc,
          orderBy: SERVICE_WEATHER_WIDGET_SORT_FIELD.criticity,
        },
        margin: {
          top: 2,
          right: 2,
          bottom: 2,
          left: 2,
        },
        advanced: true,
        alarmsList: true,
        columnSM: 6,
        columnMD: 6,
        columnLG: 6,
        heightFactor: 10,
        modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo,
      },
      size: {
        sm: 10,
        md: 10,
        lg: 10,
      },
      title: 'Weather widget(edited)',
      periodicRefresh: 180,
      filter: {},
      moreInfos: 'More infos popup(edited)',
      blockTemplate: 'Template weather item text(edited)',
      modalTemplate: 'Template modal text(edited)',
      newColumnNames: [{
        index: 8,
        data: {
          value: 'alarm.v.connector',
          label: 'New column',
          isHtml: true,
        },
      }],
      editColumnNames: [{
        index: 1,
        data: {
          value: 'alarm.v.connector',
          label: 'Connector(edited)',
          isHtml: true,
        },
      }, {
        index: 8,
        data: {
          value: 'alarm.v.connector_name',
          label: 'New column(edited)',
          isHtml: false,
        },
      }],
    };

    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.createServiceWeather(weatherWidget);
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

  'Delete test view': (browser) => {
    const { groupId, viewId } = browser.globals.defaultViewData;

    browser.completed.view.delete(groupId, viewId);
    browser.completed.view.deleteGroup(groupId);
  },
};
