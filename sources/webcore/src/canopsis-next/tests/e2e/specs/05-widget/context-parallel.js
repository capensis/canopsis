// http://nightwatchjs.org/guide#usage

const {
  CONTEXT_WIDGET_SORT_FIELD,
  SORT_ORDERS,
  FILTERS_TYPE,
  VALUE_TYPES,
  FILTER_OPERATORS,
  CONTEXT_FILTER_COLUMNS,
  CONTEXT_TYPE_OF_ENTITIES,
} = require('../../constants');
const { WIDGET_TYPES } = require('@/constants');
const { createWidgetView, removeWidgetView } = require('../../helpers/api');
const { generateTemporaryContextWidget } = require('../../helpers/entities');

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

  'Create widget context with some name': (browser) => {
    const contextWidget = {
      ...generateTemporaryContextWidget(),
      parameters: {
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: CONTEXT_WIDGET_SORT_FIELD.type,
        },
        newColumnNames: [{
          index: 3,
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
        deleteColumnNames: [2],
        filters: {
          isMix: true,
          type: FILTERS_TYPE.OR,
          title: 'Filter title',
          selected: [1],
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
        typeOfEntities: [{
          index: CONTEXT_TYPE_OF_ENTITIES.CONNECTOR,
          value: true,
        }, {
          index: CONTEXT_TYPE_OF_ENTITIES.COMPONENT,
          value: true,
        }],
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
      .clickWidget(WIDGET_TYPES.context)
      .verifyModalClosed();

    browser.completed.widget.createContext(contextWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });
  },

  'Edit widget context with some name': (browser) => {
    browser.page.view()
      .clickEditViewButton()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 6,
        md: 6,
        lg: 6,
      },
      title: 'Context widget(edited)',
      parameters: {
        advanced: true,
        sort: {
          order: SORT_ORDERS.asc,
          orderBy: CONTEXT_WIDGET_SORT_FIELD.name,
        },
        typeOfEntities: [{
          index: CONTEXT_TYPE_OF_ENTITIES.CONNECTOR,
          value: false,
        }],
      },
    });

    browser.page.widget.context()
      .clickSubmitContext();
  },

  'Delete widget context with some name': (browser) => {
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
