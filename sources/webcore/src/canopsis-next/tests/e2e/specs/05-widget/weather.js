// http://nightwatchjs.org/guide#usage

const uid = require('uid');
const { API_ROUTES } = require('../../../../src/config');
const {
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SERVICE_WEATHER_WIDGET_SORT_FIELD,
  PAGINATION_PER_PAGE_VALUES,
  SORT_ORDERS,
} = require('../../constants');
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
          .find(item => item.title === tab)._id,
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

  'Create widget weather with some name': (browser) => {
    const common = browser.page.widget.common();
    const weather = browser.page.widget.weather();

    browser.page.view()
      .clickEditViewButton()
      .clickAddWidgetButton();

    browser.page.modals.view.createWidget()
      .verifyModalOpened()
      .clickWidget('ServiceWeather')
      .verifyModalClosed();

    browser.completed.widget.setCommonFields({
      row: 'row',
      size: {
        sm: 12,
        md: 12,
        lg: 12,
      },
      title: 'Weather widget',
      periodicRefresh: 140,
      advanced: true,
      alarmsList: true,
      parameters: {
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
      })
      .clickCreateFilter();

    browser.page.modals.view.createFilter()
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();

    common.clickCreateMoreInfos();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .setField('More infos popup')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickTemplateWeatherItem();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .setField('Template weather item text')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickTemplateModal();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .setField('Template modal text')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickTemplateEntities();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .setField('Template entities text')
      .clickSubmitButton()
      .verifyModalClosed();

    browser.waitForFirstXHR(
      API_ROUTES.userPreferences,
      5000,
      () => weather.clickSubmitWeather(),
      ({ responseData }) => {
        browser.globals.temporary.widgetId = JSON.parse(responseData).data[0].widget_id;
      },
    );
  },

  'Edit widget weather with some name': (browser) => {
    const common = browser.page.widget.common();
    const weather = browser.page.widget.weather();

    browser.page.view()
      .clickEditWidgetButton(browser.globals.temporary.widgetId);

    browser.completed.widget.setCommonFields({
      advanced: true,
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
        alarmsList: {
          perPage: PAGINATION_PER_PAGE_VALUES.HUNDRED,
        },
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
    });

    common
      .clickColumnNames()
      .editColumnName(1, {
        value: 'alarm.v.connector',
        label: 'Connector(edited)',
        isHtml: true,
      })
      .editColumnName(8, {
        value: 'alarm.v.connector_name',
        label: 'New column(edited)',
        isHtml: false,
      })
      .clickCreateFilter();

    browser.page.modals.view.createFilter()
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();

    common.clickEditMoreInfos();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .clearField()
      .setField('More infos popup(edited)')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickTemplateWeatherItem();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .clearField()
      .setField('Template weather item text(edited)')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickTemplateModal();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .clearField()
      .setField('Template modal text(edited)')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickTemplateEntities();

    browser.page.modals.common.textEditor()
      .verifyModalOpened()
      .clickField()
      .clearField()
      .setField('Template entities text(edited)')
      .clickSubmitButton()
      .verifyModalClosed();

    weather.clickSubmitWeather();
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
