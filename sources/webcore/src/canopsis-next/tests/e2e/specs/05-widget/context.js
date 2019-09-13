// http://nightwatchjs.org/guide#usage

const {
  CONTEXT_WIDGET_SORT_FIELD,
  SORT_ORDERS,
  FILTERS_TYPE,
  VALUE_TYPES,
  FILTER_OPERATORS,
  CONTEXT_FILTER_COLUMNS,
} = require('../../constants');
const { generateTemporaryView, generateTemporaryContext } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.temporary = {};
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    const { view } = browser.globals.temporary;

    browser.completed.view.delete(view.group_id, view._id);

    browser.completed.logout()
      .end(done);

    delete browser.globals.temporary;
  },

  'Create test view': (browser) => {
    browser.completed.view.create(generateTemporaryView(), (view) => {
      browser.globals.temporary.view = view;
    });
  },

  'Create widget context with some name': (browser) => {
    const contextWidget = {
      ...generateTemporaryContext(),
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      parameters: {
        advanced: true,
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
          index: 1,
          value: true,
        }, {
          index: 2,
          value: true,
        }],
      },
    };
    const { temporary } = browser.globals;
    const view = browser.page.view();
    const groupsSideBar = browser.page.layout.groupsSideBar();

    groupsSideBar.clickPanelHeader(temporary.view.group_id)
      .clickLinkView(temporary.view._id);

    view.clickMenuViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget('Context')
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
      sm: 10,
      md: 10,
      lg: 10,
      title: 'Context widget(edited)',
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
