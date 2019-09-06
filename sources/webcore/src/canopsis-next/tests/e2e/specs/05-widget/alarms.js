// http://nightwatchjs.org/guide#usage

const {
  ALARMS_WIDGET_INFO_POPUP_COLUMNS,
  ALARMS_WIDGET_SORT_FIELD,
  SORT_ORDERS,
  PAGINATION_PER_PAGE_VALUES,
} = require('../../constants');
const { generateTemporaryView, generateTemporaryAlarms } = require('../../helpers/entities');

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

  'Create widget alarms with some name': (browser) => {
    const alarmsWidget = {
      ...generateTemporaryAlarms(),
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      advanced: true,
      periodicRefresh: 140,
      parameters: {
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: ALARMS_WIDGET_SORT_FIELD.component,
        },
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.HUNDRED,
        openedResolvedFilter: {
          open: true,
          resolve: true,
        },
        infoPopups: [{
          column: ALARMS_WIDGET_INFO_POPUP_COLUMNS.connectorName,
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
    const { temporary } = browser.globals;
    const view = browser.page.view();
    const groupsSideBar = browser.page.layout.groupsSideBar();

    groupsSideBar.clickPanelHeader(temporary.view.group_id)
      .clickLinkView(temporary.view._id);

    view.clickMenuViewButton()
      .clickEditViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget('AlarmsList')
      .verifyModalClosed();

    browser.completed.widget.createAlarmsList(alarmsWidget, ({ response }) => {
      browser.globals.temporary.widgetId = response.data[0].widget_id;
    });

    // browser.completed.widget.createAlarmsList({
    //   advanced: {
    //     filters: {
    //       add: {
    //         title: 'FilterTitle',
    //         or: true,
    //         rule: {
    //           field: 2,
    //           operator: 2,
    //         },
    //       },
    //     },
    //     infoPopap: {
    //       add: {
    //         column: 2,
    //         template: 'Template',
    //       },
    //     },
    //   },
    // });
  },

  'Edit widget alarms with some name': (browser) => {
    browser.page.view()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      size: {
        sm: 10,
        md: 10,
        lg: 10,
      },
      title: 'Alarms widget(edited)',
      advanced: true,
      periodicRefresh: 180,
      parameters: {
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: ALARMS_WIDGET_SORT_FIELD.connector,
        },
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.TWENTY,
      },
    });

    browser.page.widget.alarms()
      .clickSubmitAlarms();
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
