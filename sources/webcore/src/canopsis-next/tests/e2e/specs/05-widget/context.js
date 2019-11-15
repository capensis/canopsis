// http://nightwatchjs.org/guide#usage

const {
  CONTEXT_WIDGET_SORT_FIELD,
  SORT_ORDERS,
  FILTERS_TYPE,
  VALUE_TYPES,
  FILTER_OPERATORS,
  CONTEXT_FILTER_COLUMNS,
  CONTEXT_TYPE_OF_ENTITIES,
  CONTEXT_MASS_ACTIONS,
  CONTEXT_SHARED_ACTIONS,
  PAGINATION_PER_PAGE_VALUES,
  CONTEXT_CREATE_ENTITY_TAB,
  WIDGET_TYPES,
} = require('../../constants');
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

  'Table widget context': (browser) => {
    const commonTable = browser.page.tables.common();
    const entitiesSelectForm = browser.page.forms.entitiesSelect();
    const manageInfosForm = browser.page.forms.manageInfos();
    const confirmation = browser.page.modals.common.confirmation();
    const createEntityModal = browser.page.modals.context.createEntity();
    const addEntityInfoModal = browser.page.modals.context.addEntityInfo();

    const firstId = '8cd9c153f596';

    browser.page.view()
      .clickMenuViewButton();

    commonTable
      .setAllCheckbox(true)
      .clickOnMassAction(CONTEXT_MASS_ACTIONS.DELETE_ENTITY);

    confirmation
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();

    commonTable
      // .clickTableHeaderCell('Name')
      .clickOnRow(firstId)
      .clickOnSharedAction(firstId, CONTEXT_SHARED_ACTIONS.EDIT);

    createEntityModal
      .verifyModalOpened()
      .clickName()
      .clearName()
      .setName('Create entity name')
      .setEnabled(true)
      .setType(1)
      .clickImpact();

    entitiesSelectForm
      .setItemPerPage(PAGINATION_PER_PAGE_VALUES.TEN)
      .clickSearch()
      .clearSearch()
      .setSearch('a')
      .clickSubmitSearch()
      .setRowCheckbox('Engine_watcher/c1a24f0183a8')
      .clickAddEntity('Engine_watcher/c1a24f0183a8')
      .clickRemoveEntity('Engine_watcher/c1a24f0183a8')
      .setAllCheckbox()
      .clickAddCollection()
      .clickAddCollection()
      .clickClearEntities();

    createEntityModal
      .clickImpact()
      .clickDepends();

    entitiesSelectForm
      .setItemPerPage(PAGINATION_PER_PAGE_VALUES.TEN)
      .clickSearch()
      .clearSearch()
      .setSearch('a')
      .clickSubmitSearch()
      .setRowCheckbox('Engine_watcher/c1a24f0183a8')
      .clickAddEntity('Engine_watcher/c1a24f0183a8')
      .clickRemoveEntity('Engine_watcher/c1a24f0183a8')
      .setAllCheckbox()
      .clickAddCollection()
      .clickAddCollection()
      .clickClearEntities();

    createEntityModal
      .clickTab(CONTEXT_CREATE_ENTITY_TAB.MANAGE_INFOS);

    manageInfosForm
      .clickAddInfo();

    addEntityInfoModal
      .verifyModalOpened()
      .clickEntityInfoName()
      .clearEntityInfoName()
      .setEntityInfoName('Information name')
      .clickEntityInfoDescription()
      .clearEntityInfoDescription()
      .setEntityInfoDescription('Information description')
      .clickEntityInfoValue()
      .clearEntityInfoValue()
      .setEntityInfoValue('Information value')
      .clickSubmitButton()
      .verifyModalClosed();

    manageInfosForm
      .setItemPerPage(PAGINATION_PER_PAGE_VALUES.TEN)
      .clickRowEditInfo('Information value');

    addEntityInfoModal
      .verifyModalOpened()
      .clickEntityInfoValue()
      .clearEntityInfoValue()
      .setEntityInfoValue('New information value')
      .clickSubmitButton()
      .verifyModalClosed();

    manageInfosForm.clickRowDeleteInfo('New information value');

    createEntityModal
      .clickTab(CONTEXT_CREATE_ENTITY_TAB.FORM)
      .clickCancelButton()
      .verifyModalClosed();

    commonTable
      .clickSearchInput()
      .clearSearchInput()
      .setSearchInput('search string')
      .clickSearchButton()
      .clickSearchResetButton()
      .setMixFilters(false)
      .selectFilter(1)
      .setFiltersType(FILTERS_TYPE.AND)
      .clickNextPageTopPagination()
      .clickPreviousPageTopPagination()
      .clickNextPageBottomPagination()
      .clickPreviousPageBottomPagination()
      .clickOnPageBottomPagination(2)
      .setItemPerPage(PAGINATION_PER_PAGE_VALUES.FIVE);
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
