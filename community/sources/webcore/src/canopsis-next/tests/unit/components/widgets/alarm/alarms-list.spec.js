import Faker from 'faker';
import { omit } from 'lodash';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { mockModals, mockPopups, mockSocket } from '@unit/utils/mock-hooks';
import {
  createActiveViewModule,
  createMockedStoreModule,
  createMockedStoreModules,
  createServiceModule,
  createWidgetModule,
} from '@unit/utils/store';
import { fakeAlarmDetails, fakeStaticAlarms } from '@unit/data/alarm';

import { API_HOST, API_ROUTES } from '@/config';
import {
  CANOPSIS_EDITION,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_STATUSES,
  MODALS,
  LIVE_REPORTING_QUICK_RANGES,
  TIME_UNITS,
  USER_PERMISSIONS,
} from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';
import { formToAdvancedSearch } from '@/helpers/search/advanced-search';

import AlarmsList from '@/components/widgets/alarm/alarms-list.vue';

const stubs = {
  'c-alarm-advanced-search': true,
  'c-entity-category-field': true,
  'v-switch': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'alarms-list-remediation-instructions-filter': true,
  'c-action-btn': true,
  'c-pagination': true,
  'c-density-btn-toggle': true,
  'c-table-pagination': true,
  'c-alarm-tag-field': true,
  'mass-actions-panel': true,
  'alarms-list-table': {
    template: `
      <div class="alarms-list-table">
        <slot />
      </div>
    `,
  },
};

const snapshotStubs = {
  'c-alarm-advanced-search': true,
  'c-entity-category-field': true,
  'v-switch': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'alarms-list-remediation-instructions-filter': true,
  'c-action-btn': true,
  'c-pagination': true,
  'alarms-list-table': true,
  'c-table-pagination': true,
  'c-density-btn-toggle': true,
  'c-alarm-tag-field': true,
  'mass-actions-panel': true,
};

const selectVSwitch = wrapper => wrapper.find('v-switch-stub');
const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');
const selectCategoryField = wrapper => wrapper.find('c-entity-category-field-stub');
const selectAlarmTagField = wrapper => wrapper.find('c-alarm-tag-field-stub');
const selectExportButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);
const selectLiveReportingButton = wrapper => wrapper.findAll('c-action-btn-stub').at(0);
const selectInstructionsFiltersField = wrapper => wrapper.find('alarms-list-remediation-instructions-filter-stub');
const selectRemoveHistoryButton = wrapper => wrapper.find('v-chip-stub');
const selectAlarmAdvancedSearchField = wrapper => wrapper.find('c-alarm-advanced-search-stub');
const selectAlarmsListTable = wrapper => wrapper.find('.alarms-list-table');

describe('alarms-list', () => {
  const $popups = mockPopups();
  const $modals = mockModals();
  const $socket = mockSocket();

  const nowTimestamp = 1386435600000;
  const nowSubtractOneYearUnix = 1354834800;

  const totalItems = 10;
  const alarms = fakeStaticAlarms({
    totalItems,
    timestamp: nowTimestamp,
  });

  const userPreferences = {
    content: {
      isCorrelationEnabled: false,
      onlyBookmarks: false,
      itemsPerPage: 13,
      category: 'category-id',
    },
  };

  const exportAlarmData = {
    _id: 'export-alarm-id',
    status: EXPORT_STATUSES.completed,
  };

  const exportFailedAlarmData = {
    _id: 'export-alarm-id',
    status: EXPORT_STATUSES.failed,
  };
  const widget = {
    ...generatePreparedDefaultAlarmListWidget(),

    _id: '880c5d0c-3f31-477c-8365-2f90389326cc',
  };

  const defaultWidgetColumnsActiveColumns = widget.parameters.widgetColumns.map(v => v.value);
  const defaultMoreInfoTemplateActiveColumns = [
    'v.display_name',
    'is_meta_alarm',
    'v.max_state',
    'v.initial_state',
    'v.initial_output',
    'v.creation_date',
    'v.ticket.ticket_url',
    'v.ticket.ticket',
    'v.ticket',
    'v.duration',
    'v.events_count',
    'v.total_state_changes',
    'pbehavior.name',
    'pbehavior.type.name',
    'pbehavior',
    'v.comments',
    'v.tickets',
    'entity.infos',
    'entity.component_infos',
    'v.infos',
  ];

  const defaultQuery = {
    page: 1,
    filters: [],
    sortBy: [],
    sortDesc: [],
    instructionsFilter: {},
    active_columns: [...defaultWidgetColumnsActiveColumns, ...defaultMoreInfoTemplateActiveColumns],
    correlation: userPreferences.content.isCorrelationEnabled,
    only_bookmarks: userPreferences.content.onlyBookmarks,
    category: userPreferences.content.category,
    itemsPerPage: userPreferences.content.itemsPerPage,
    tstart: LIVE_REPORTING_QUICK_RANGES.last1Year.start,
    tstop: LIVE_REPORTING_QUICK_RANGES.last1Year.stop,
    opened: null,
    search: 'search',
    tags: undefined,
  };
  const view = {
    enabled: true,
    title: 'Alarm widgets',
    description: 'Alarm widgets',
    tabs: [
      {
        widgets: [widget],
      },
    ],
    tags: ['alarm'],
    periodic_refresh: {
      value: 1,
      unit: 's',
      enabled: false,
    },
    author: 'root',
    group: {
      _id: 'text-widget-group',
    },
  };
  const updateUserPreference = jest.fn();
  const updateView = jest.fn();
  const updateQuery = jest.fn();
  const hideSideBar = jest.fn();
  const fetchAlarmsList = jest.fn();
  const createAlarmsListExport = jest.fn().mockReturnValue(exportAlarmData);
  const fetchAlarmsListExport = jest.fn().mockReturnValue(exportAlarmData);
  const fetchAlarmDetails = jest.fn();
  const fetchAlarmsDetailsList = jest.fn();
  const updateAlarmDetailsQuery = jest.fn();
  const removeAlarmDetailsQuery = jest.fn();
  const fetchTagsList = jest.fn();
  const sideBarModule = {
    name: 'sideBar',
    actions: {
      hide: hideSideBar,
    },
  };
  const infoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.pro },
  };
  const queryModule = {
    name: 'query',
    getters: {
      getQueryNonceById: () => () => ({}),
      getQueryById: () => () => defaultQuery,
    },
    actions: {
      update: updateQuery,
    },
  };
  const viewModule = {
    name: 'view',
    getters: {
      item: view,
    },
    actions: {
      update: updateView,
    },
  };
  const fetchUserPreference = jest.fn();
  const userPreferenceModule = {
    name: 'userPreference',
    getters: {
      getItemByWidgetId: () => () => userPreferences,
    },
    actions: {
      update: updateUserPreference,
      fetchItem: fetchUserPreference,
    },
  };
  const authModule = {
    name: 'auth',
    getters: {
      currentUser: {},
      currentUserPermissionsById: {},
    },
  };

  const alarmDetailsModule = createMockedStoreModule({
    name: 'details',
    getters: {
      getItem: () => () => fakeAlarmDetails(),
      getPending: () => () => false,
      getQuery: () => () => ({ page: 1, limit: 10 }),
      getQueries: () => () => [
        { page: 2, limit: 5 },
        { page: 1, limit: 10 },
      ],
    },
    actions: {
      fetchItem: fetchAlarmDetails,
      fetchList: fetchAlarmsDetailsList,
      updateQuery: updateAlarmDetailsQuery,
      removeQuery: removeAlarmDetailsQuery,
    },
  });

  const alarmModule = {
    name: 'alarm',
    modules: {
      details: alarmDetailsModule,
    },
    getters: {
      getMetaByWidgetId: () => () => ({
        total_count: totalItems,
      }),
      getListByWidgetId: () => () => alarms,
      getPendingByWidgetId: () => () => false,
    },
    actions: {
      fetchList: fetchAlarmsList,
      createAlarmsListExport,
      fetchAlarmsListExport,
    },
  };

  const alarmTagModule = {
    name: 'alarmTag',
    getters: {
      pending: () => false,
    },
    actions: {
      fetchList: fetchTagsList,
    },
  };

  const { serviceModule, fetchEntityInfosKeysWithoutStore } = createServiceModule();
  const { widgetModule } = createWidgetModule();
  const { activeViewModule } = createActiveViewModule();

  const store = createMockedStoreModules([
    alarmModule,
    sideBarModule,
    infoModule,
    queryModule,
    viewModule,
    userPreferenceModule,
    authModule,
    alarmTagModule,
    serviceModule,
    widgetModule,
    activeViewModule,
  ]);

  const factory = generateShallowRenderer(AlarmsList, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
    mocks: {
      $socket,
    },
  });
  const snapshotFactory = generateRenderer(AlarmsList, {
    stubs: snapshotStubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
    mocks: {
      $socket,
    },
  });

  beforeEach(() => {
    jest.useFakeTimers({ now: nowTimestamp });
  });

  afterEach(() => {
    fetchUserPreference.mockClear();
    updateUserPreference.mockClear();
    updateView.mockClear();
    updateQuery.mockClear();
    hideSideBar.mockClear();
    fetchTagsList.mockClear();
    fetchEntityInfosKeysWithoutStore.mockClear();
  });

  it('Query updated after mount', async () => {
    factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    expect(fetchUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      { id: widget._id },
    );

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...omit(defaultQuery, ['tstart', 'tstop', 'filters']),
          filter: undefined,
          lockedFilter: null,
          page: 1,
          with_instructions: true,
          with_tag_colors: true,
          with_declare_tickets: true,
          with_links: true,
          opened: true,
        },
      },
    );
  });

  it('User preferences not fetched after mount with local widget prop', async () => {
    factory({
      store,
      propsData: {
        widget,
        localWidget: true,
      },
    });

    await flushPromises();

    expect(fetchUserPreference).not.toHaveBeenCalled();

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...omit(defaultQuery, ['tstart', 'tstop', 'filters']),
          filter: undefined,
          lockedFilter: null,
          page: 1,
          with_instructions: true,
          with_tag_colors: true,
          with_declare_tickets: true,
          with_links: true,
          opened: true,
        },
      },
    );
  });

  it('Correlation updated after trigger correlation field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.correlation]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const correlationField = selectVSwitch(wrapper);

    correlationField.triggerCustomEvent('change', !userPreferences.content.isCorrelationEnabled);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            isCorrelationEnabled: !userPreferences.content.isCorrelationEnabled,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,

          page: 1,
          correlation: !userPreferences.content.isCorrelationEnabled,
        },
      },
    );
  });

  it('Only bookmark updated after trigger filter by bookmark field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const filterByBookmarkField = selectVSwitch(wrapper);

    filterByBookmarkField.triggerCustomEvent('change', !userPreferences.content.onlyBookmarks);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            onlyBookmarks: !userPreferences.content.onlyBookmarks,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,

          page: 1,
          only_bookmarks: !userPreferences.content.onlyBookmarks,
        },
      },
    );
  });

  it('Tags updated after trigger alarm tag field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    const newTags = [Faker.datatype.string(), Faker.datatype.string()];

    await flushPromises();

    updateQuery.mockClear();

    const tagField = selectAlarmTagField(wrapper);

    tagField.triggerCustomEvent('input', newTags);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            tags: newTags,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,

          page: 1,
          tags: newTags,
        },
      },
    );
  });

  it('Filter updated after trigger filter field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.userFilter]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const filterSelectorField = selectFilterSelectorField(wrapper);

    const selectedFilter = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      filter: {},
    };

    filterSelectorField.triggerCustomEvent('input', selectedFilter._id);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            mainFilter: selectedFilter._id,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          page: 1,
          filter: selectedFilter._id,
        },
      },
    );
  });

  it('Instruction filters updated after trigger filter field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newInstructionsFilter = {
      instruction_filter_type: 'manual',
      instruction_type: 'auto',
      instruction_statuses: ['pending', 'done'],
      instruction_ids: ['id1', 'id2'],
    };

    const filterComponent = selectInstructionsFiltersField(wrapper);

    filterComponent.triggerCustomEvent('input', newInstructionsFilter);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            instructionsFilter: newInstructionsFilter,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          instructionsFilter: newInstructionsFilter,
        },
      },
    );
  });

  it('Interval query removed after click on the button', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const removeHistoryButton = selectRemoveHistoryButton(wrapper);

    removeHistoryButton.triggerCustomEvent('click:close');

    await flushPromises();

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...omit(defaultQuery, ['tstart', 'tstop']),

          page: 1,
        },
      },
    );
  });

  it('Interval modal showed after click on the live reporting button', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const liveReportingButton = selectLiveReportingButton(wrapper);

    liveReportingButton.triggerCustomEvent('click');

    await flushPromises();

    expect($modals.show).toHaveBeenCalledWith(
      {
        name: MODALS.editLiveReporting,
        config: {
          action: expect.any(Function),
          tstart: defaultQuery.tstart,
          tstop: defaultQuery.tstop,
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const actionValue = {
      tstart: LIVE_REPORTING_QUICK_RANGES.last3Hour.start,
      tstop: LIVE_REPORTING_QUICK_RANGES.last3Hour.stop,
    };

    modalArguments.config.action(actionValue);

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          ...actionValue,

          page: 1,
        },
      },
    );
  });

  it('Category updated after trigger category field', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.category]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const categoryField = selectCategoryField(wrapper);

    const newCategory = {
      _id: Faker.datatype.string(),
    };

    categoryField.triggerCustomEvent('input', newCategory);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            category: newCategory._id,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,

          page: 1,
          category: newCategory._id,
        },
      },
    );
  });

  it('Limit updated after trigger table pagination', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newLimit = Faker.datatype.number();
    selectAlarmsListTable(wrapper).triggerCustomEvent('update:items-per-page', newLimit);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            itemsPerPage: newLimit,
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,

          page: 1,
          itemsPerPage: newLimit,
        },
      },
    );
  });

  it('Page updated after trigger table pagination', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const newPage = Faker.datatype.number();
    selectAlarmsListTable(wrapper).triggerCustomEvent('update:page', newPage);

    await flushPromises();

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          page: newPage,
        },
      },
    );
  });

  it('Widget exported after trigger export button', async () => {
    const originalWindowOpen = window.open;
    window.open = jest.fn();

    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          filters: defaultQuery.filters,
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          only_bookmarks: defaultQuery.only_bookmarks,
          opened: defaultQuery.opened,
          tstart: nowSubtractOneYearUnix,
          tstop: 1386435600,
          fields: widget.parameters.widgetExportColumns.map(({ text, value }) => ({
            label: text,
            name: value,
            template: undefined,
          })),
          separator: widget.parameters.exportCsvSeparator,
          time_format: widget.parameters.exportCsvDatetimeFormat,
        },
        widgetId: widget._id,
      },
    );

    await flushPromises();

    jest.runAllTimers();

    expect(fetchAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
    );

    await flushPromises();

    expect(window.open).toHaveBeenCalledWith(
      `${API_HOST}${API_ROUTES.alarmListExport}/${exportAlarmData._id}/download`,
      '_blank',
    );

    window.open = originalWindowOpen;
  });

  /**
   * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/3997
   * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/4102
   */
  it('Widget exported after trigger export button with invalid structure', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds,
          },
        },
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          filters: defaultQuery.filters,
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          only_bookmarks: defaultQuery.only_bookmarks,
          opened: defaultQuery.opened,
          tstart: nowSubtractOneYearUnix,
          tstop: 1386435600,
          fields: widget.parameters.widgetExportColumns.map(({ text, value }) => ({
            label: text,
            name: value,
            template: undefined,
          })),
          separator: widget.parameters.exportCsvSeparator,
          time_format: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
        },
        widgetId: widget._id,
      },
    );

    wrapper.destroy();
  });

  it('Widget exported after trigger export button without export columns', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            widgetExportColumns: [],
          },
        },
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          filters: defaultQuery.filters,
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          only_bookmarks: defaultQuery.only_bookmarks,
          opened: defaultQuery.opened,
          tstart: nowSubtractOneYearUnix,
          tstop: 1386435600,
          fields: widget.parameters.widgetColumns.map(({ text, value }) => ({
            label: text,
            name: value,
            template: undefined,
          })),
          separator: widget.parameters.exportCsvSeparator,
          time_format: widget.parameters.exportCsvDatetimeFormat,
        },
        widgetId: widget._id,
      },
    );

    wrapper.destroy();
  });

  it('Widget exported after trigger export button with long request time', async () => {
    const originalWindowOpen = window.open;
    window.open = jest.fn();

    const longFetchAlarmsListExport = jest.fn()
      .mockReturnValueOnce({
        _id: exportAlarmData._id,
        status: EXPORT_STATUSES.running,
      })
      .mockReturnValueOnce(exportAlarmData);

    const wrapper = factory({
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport,
            fetchAlarmsListExport: longFetchAlarmsListExport,
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    await flushPromises();

    jest.runAllTimers();

    expect(longFetchAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
    );

    await flushPromises();

    jest.runAllTimers();

    expect(longFetchAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
    );

    await flushPromises();

    expect(window.open).toHaveBeenCalled();

    window.open = originalWindowOpen;
  });

  it('Error popup showed exported after trigger export button with failed create export', async () => {
    const wrapper = factory({
      mocks: {
        $popups,
      },
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport: jest.fn().mockRejectedValue(),
            fetchAlarmsListExport,
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'Failed to export alarms list in CSV format',
    });
  });

  it('Error popup showed exported after trigger export button with failed fetch export', async () => {
    jest.useFakeTimers('legacy');

    const wrapper = factory({
      mocks: {
        $popups,
      },
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport,
            fetchAlarmsListExport: jest.fn().mockRejectedValue(),
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    await flushPromises();

    jest.runAllTimers();

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'Failed to export alarms list in CSV format',
    });
  });

  it('Error popup showed exported after trigger export button with failed status fetch export', async () => {
    jest.useFakeTimers('legacy');

    const wrapper = factory({
      mocks: {
        $popups,
      },
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport,
            fetchAlarmsListExport: jest.fn().mockReturnValue(exportFailedAlarmData),
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    const exportButton = selectExportButton(wrapper);

    exportButton.triggerCustomEvent('click');

    await flushPromises();

    jest.runAllTimers();

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'Failed to export alarms list in CSV format',
    });
  });

  it('Alarms not fetched after change query without columns', async () => {
    const wrapper = factory({
      store,
      data: () => ({
        testQuery: {},
      }),
      computed: {
        query: {
          get() {
            return this.testQuery;
          },
          set() {},
        },
      },
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            widgetColumns: [],
          },
        },
      },
    });

    wrapper.vm.testQuery = defaultQuery;

    await flushPromises();

    expect(fetchAlarmsList).not.toHaveBeenCalled();
  });

  it('Alarms fetched after change query', async () => {
    const expanded = {};
    const wrapper = factory({
      store,
      data: () => ({
        testQuery: {},
      }),
      computed: {
        query: {
          get() {
            return this.testQuery;
          },
          set() {},
        },
      },
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({ expanded }),
        },
      },
      propsData: {
        widget,
      },
    });

    wrapper.vm.testQuery = defaultQuery;

    await flushPromises();

    expect(fetchAlarmsList).toHaveBeenCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: {
          ...omit(defaultQuery, ['sortBy', 'sortDesc', 'itemsPerPage', 'instructionsFilter']),

          limit: defaultQuery.itemsPerPage,
          tstart: expect.any(Number),
          tstop: expect.any(Number),
        },
      },
    );
  });

  it('Alarms fetched after change query nonce', async () => {
    const firstAlarmId = alarms[0]._id;

    const expanded = [{ _id: 'non-exist-id' }, { _id: firstAlarmId }];
    const wrapper = factory({
      store,
      data: () => ({
        testTabQueryNonce: 0,
      }),
      computed: {
        tabQueryNonce: {
          get() {
            return this.testTabQueryNonce;
          },
          set() {},
        },
      },
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({ expanded }),
        },
      },
      propsData: {
        widget,
      },
    });

    wrapper.vm.testTabQueryNonce += 1;

    await flushPromises();

    expect(fetchAlarmsList).toHaveBeenCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: {
          ...omit(defaultQuery, ['sortBy', 'sortDesc', 'itemsPerPage', 'instructionsFilter']),

          limit: defaultQuery.itemsPerPage,
          tstart: expect.any(Number),
          tstop: expect.any(Number),
        },
      },
    );
    expect(expanded).toEqual([
      { _id: 'non-exist-id' },
      { _id: firstAlarmId },
    ]);
  });

  it('Periodic started after mount with enabled value', async () => {
    jest.useFakeTimers();
    jest.spyOn(window, 'setInterval');

    const expanded = {};
    factory({
      store,
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({ expanded }),
        },
      },
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            periodic_refresh: {
              enabled: true,
              unit: TIME_UNITS.second,
              value: 1,
            },
          },
        },
      },
    });

    expect(setInterval).toHaveBeenCalledTimes(1);
    expect(setInterval).toHaveBeenCalledWith(
      expect.any(Function),
      1000,
    );

    jest.advanceTimersByTime(1000);

    expect(fetchAlarmsList).toHaveBeenCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: {
          ...omit(defaultQuery, ['sortBy', 'sortDesc', 'itemsPerPage', 'instructionsFilter']),

          limit: defaultQuery.itemsPerPage,
          tstart: expect.any(Number),
          tstop: expect.any(Number),
        },
      },
    );
  });

  it('Interval cleared after update periodic refresh', async () => {
    jest.spyOn(global, 'setInterval');
    jest.spyOn(global, 'clearInterval');

    const expanded = {};
    const wrapper = factory({
      store,
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({ expanded }),
        },
      },
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            periodic_refresh: {
              enabled: true,
              unit: TIME_UNITS.minute,
              value: 1,
            },
          },
        },
      },
    });

    expect(setInterval).toHaveBeenCalled();
    setInterval.mockClear();

    await wrapper.setProps({
      widget: {
        ...widget,
        parameters: {
          ...widget.parameters,
          periodic_refresh: {
            enabled: true,
            unit: TIME_UNITS.minute,
            value: 2,
          },
        },
      },
    });

    expect(clearInterval).toHaveBeenCalledTimes(1);

    expect(setInterval).toHaveBeenCalledTimes(1);
    expect(setInterval).toHaveBeenCalledWith(
      expect.any(Function),
      120000,
    );
  });

  it('Interval cleared after destroy', async () => {
    jest.spyOn(global, 'setInterval');
    jest.spyOn(global, 'clearInterval');

    const expanded = {};
    const wrapper = factory({
      store,
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({ expanded }),
        },
      },
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            periodic_refresh: {
              enabled: true,
              unit: TIME_UNITS.minute,
              value: 1,
            },
          },
        },
      },
    });

    expect(setInterval).toHaveBeenCalledTimes(1);

    wrapper.destroy();

    expect(clearInterval).toHaveBeenCalledTimes(1);
  });

  it('Search updated after alarm-advanced-search-field submit trigger', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const alarmAdvancedSearchField = selectAlarmAdvancedSearchField(wrapper);

    const advancedSearchValue = formToAdvancedSearch();

    advancedSearchValue._id = Faker.datatype.string();
    advancedSearchValue.search = Faker.datatype.string();
    advancedSearchValue.pinned = false;

    alarmAdvancedSearchField.triggerCustomEvent('submit', advancedSearchValue);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            searches: [advancedSearchValue],
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,

          page: 1,
          search: advancedSearchValue.search,
        },
      },
    );
  });

  it('Patterns updated after alarm-advanced-search-field submit trigger', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const alarmAdvancedSearchField = selectAlarmAdvancedSearchField(wrapper);

    const advancedSearchValue = formToAdvancedSearch();

    advancedSearchValue._id = Faker.datatype.string();
    advancedSearchValue.pinned = false;
    advancedSearchValue.alarm_pattern = [[{ field: 'connector', cond: 'eq', value: 'test' }]];
    advancedSearchValue.entity_pattern = [[{ field: 'connector', cond: 'eq', value: 'test1' }]];
    advancedSearchValue.pbehavior_pattern = [[{ field: 'connector', cond: 'eq', value: 'test2' }]];

    alarmAdvancedSearchField.triggerCustomEvent('submit', advancedSearchValue);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            searches: [advancedSearchValue],
          },
        },
      },
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...omit(defaultQuery, ['search']),

          page: 1,
          alarm_pattern: JSON.stringify(advancedSearchValue.alarm_pattern),
          entity_pattern: JSON.stringify(advancedSearchValue.entity_pattern),
          pbehavior_pattern: JSON.stringify(advancedSearchValue.pbehavior_pattern),
        },
      },
    );
  });

  it('Search field cleared after alarm-advanced-search-field reset trigger', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const alarmAdvancedSearchField = selectAlarmAdvancedSearchField(wrapper);

    alarmAdvancedSearchField.triggerCustomEvent('reset');

    await flushPromises();

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: omit(defaultQuery, ['search']),
      },
    );
  });

  it('Search pinned after alarm-advanced-search-field toggle-pin:search trigger', async () => {
    const firstSearch = formToAdvancedSearch();
    const secondSearch = formToAdvancedSearch();

    firstSearch._id = Faker.datatype.string();
    firstSearch.search = 'first';
    firstSearch.pinned = true;

    secondSearch._id = Faker.datatype.string();
    secondSearch.search = 'second';
    secondSearch.pinned = false;

    const wrapper = factory({
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        alarmTagModule,
        serviceModule,
        alarmModule,
        authModule,
        widgetModule,
        activeViewModule,
        {
          ...userPreferenceModule,
          getters: {
            getItemByWidgetId: () => () => ({
              ...userPreferences,
              content: {
                ...userPreferences.content,
                searches: [firstSearch, secondSearch],
              },
            }),
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const alarmAdvancedSearchField = selectAlarmAdvancedSearchField(wrapper);

    alarmAdvancedSearchField.triggerCustomEvent('toggle-pin:search', secondSearch._id);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            searches: [{ ...secondSearch, pinned: true }, firstSearch],
          },
        },
      },
    );
  });

  it('Search removed after alarm-advanced-search-field remove:search trigger', async () => {
    const firstSearch = formToAdvancedSearch();
    const secondSearch = formToAdvancedSearch();

    firstSearch._id = Faker.datatype.string();
    firstSearch.search = Faker.datatype.string();
    firstSearch.pinned = true;

    secondSearch._id = Faker.datatype.string();
    secondSearch.search = Faker.datatype.string();
    secondSearch.pinned = false;

    const wrapper = factory({
      store: createMockedStoreModules([
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        alarmTagModule,
        serviceModule,
        alarmModule,
        authModule,
        widgetModule,
        activeViewModule,
        {
          ...userPreferenceModule,
          getters: {
            getItemByWidgetId: () => () => ({
              ...userPreferences,
              content: {
                ...userPreferences.content,
                searches: [firstSearch, secondSearch],
              },
            }),
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const alarmAdvancedSearchField = selectAlarmAdvancedSearchField(wrapper);

    alarmAdvancedSearchField.triggerCustomEvent('remove:search', firstSearch._id);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            searches: [secondSearch],
          },
        },
      },
    );
  });

  it('Renders `alarms-list` with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list` with default props and user filter permissions and correlation and bookmark', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget,
      },
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: { actions: [] },
              [USER_PERMISSIONS.business.alarmsList.actions.userFilter]: { actions: [] },
              [USER_PERMISSIONS.business.alarmsList.actions.correlation]: { actions: [] },
              [USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark]: { actions: [] },
            },
          },
        },
      ]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `alarms-list` with clear filter disabled props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,

            clearFilterEnabled: true,
          },
        },
      },
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        alarmTagModule,
        serviceModule,
        widgetModule,
        activeViewModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USER_PERMISSIONS.business.alarmsList.actions.userFilter]: { actions: [] },
            },
          },
        },
        {
          ...userPreferenceModule,
          getters: {
            getItemByWidgetId: () => () => ({
              content: {
                ...userPreferences,

                searches: ['item 1', 'item 2'],
              },
            }),
          },
        },
      ]),
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
