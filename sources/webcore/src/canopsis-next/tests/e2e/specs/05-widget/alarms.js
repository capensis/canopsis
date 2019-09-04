// http://nightwatchjs.org/guide#usage

const uid = require('uid');
const { API_ROUTES } = require('../../../../src/config');
const { SERVICE_ALARMS_WIDGET_SORT_FIELD, SORT_ORDERS, PAGINATION_PER_PAGE_VALUES } = require('../../constants');
const { generateTemporaryView } = require('../../helpers/entities');

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

  'Create test tab': (browser) => {
    const { temporary } = browser.globals;
    const tab = `tab-${uid()}`;

    browser.page.layout.groupsSideBar()
      .clickPanelHeader(temporary.view.group_id)
      .clickLinkView(temporary.view._id);

    browser.page.view()
      .clickMenuViewButton()
      .clickAddViewButton();

    browser.page.modals.common.textFieldEditor()
      .verifyModalOpened()
      .setField(tab);

    browser.waitForFirstXHR(
      `${API_ROUTES.view}/${temporary.view._id}`,
      5000,
      () => browser.page.modals.common.textFieldEditor()
        .clickSubmitButton(),
      ({ responseData, requestData }) => temporary.view = {
        tab,
        ...temporary.view,
        tabId: JSON.parse(requestData).tabs
          .filter(item => item.title === tab)[0]._id,
        ...JSON.parse(requestData),
        ...JSON.parse(responseData),
      },
    );

    browser.page.modals.common.textFieldEditor()
      .verifyModalClosed();
  },

  'Open test tab': (browser) => {
    const { temporary } = browser.globals;

    browser.page.view()
      .clickTab(temporary.view.tabId);
  },

  'Create widget alarms with some name': (browser) => {
    const common = browser.page.widget.common();
    const alarms = browser.page.widget.alarms();

    browser.page.view()
      .clickEditViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget('AlarmsList')
      .verifyModalClosed();

    browser.completed.widget.createAlarmsList({
      row: 'row',
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      title: 'Alarms widget',
      advanced: true,
      periodicRefresh: 140,
      parameters: {
        sort: {
          order: SORT_ORDERS.desc,
          orderBy: SERVICE_ALARMS_WIDGET_SORT_FIELD.component,
        },
        elementsPerPage: PAGINATION_PER_PAGE_VALUES.HUNDRED,
        openedResolvedFilter: {
          open: true,
          resolve: true,
        },
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
      },
    });

    common
      .clickColumnNames()
      .editColumnName(1, {
        value: 'alarm.v.changeConnector',
        label: 'Connector(changed)',
        isHtml: true,
      })
      .clickColumnNameDownWard(1)
      .clickColumnNameUpWard(2)
      .clickDeleteColumnName(2)
      .clickAddColumnName()
      .editColumnName(8, {
        value: 'alarm.v.connector',
        label: 'New column',
        isHtml: true,
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

    browser.waitForFirstXHR(
      API_ROUTES.userPreferences,
      5000,
      () => alarms.clickSubmitAlarms(),
      ({ responseData }) => {
        browser.globals.temporary.widgetId = JSON.parse(responseData).data[0].widget_id;
      },
    );
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
          orderBy: SERVICE_ALARMS_WIDGET_SORT_FIELD.connector,
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
