import Faker from 'faker';
import flushPromises from 'flush-promises';
import { saveAs } from 'file-saver';
import { omit } from 'lodash';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow, mockModals, mockPopups } from '@unit/utils/mock-hooks';
import { createMockedStoreModule, createMockedStoreModules } from '@unit/utils/store';
import { fakeAlarmDetails, fakeStaticAlarms } from '@unit/data/alarm';

import {
  CANOPSIS_EDITION,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_STATUSES,
  MODALS,
  QUICK_RANGES,
  REMEDIATION_INSTRUCTION_TYPES,
  TIME_UNITS,
  USERS_PERMISSIONS,
} from '@/constants';

import AlarmsList from '@/components/widgets/alarm/alarms-list.vue';
import { generateDefaultAlarmListWidgetForm } from '@/helpers/entities';

jest.mock('file-saver', () => ({
  saveAs: jest.fn(),
}));

const localVue = createVueInstance();

const stubs = {
  'c-advanced-search-field': true,
  'c-entity-category-field': true,
  'v-switch': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'alarms-list-remediation-instructions-filters': true,
  'c-action-btn': true,
  'c-pagination': true,
  'c-density-btn-toggle': true,
  'c-table-pagination': true,
  'alarms-expand-panel-tour': true,
  'alarms-list-table': {
    template: `
      <div class="alarms-list-table">
        <slot />
      </div>
    `,
  },
};

const snapshotStubs = {
  'c-advanced-search-field': true,
  'c-entity-category-field': true,
  'v-switch': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'alarms-list-remediation-instructions-filters': true,
  'c-action-btn': true,
  'c-pagination': true,
  'alarms-list-table': true,
  'c-table-pagination': true,
  'c-density-btn-toggle': true,
  'alarms-expand-panel-tour': true,
};

const factory = (options = {}) => shallowMount(AlarmsList, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(AlarmsList, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectCorrelationField = wrapper => wrapper.find('v-switch-stub');
const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');
const selectCategoryField = wrapper => wrapper.find('c-entity-category-field-stub');
const selectTablePaginationField = wrapper => wrapper.find('c-table-pagination-stub');
const selectExportButton = wrapper => wrapper.findAll('c-action-btn-stub').at(1);
const selectLiveReportingButton = wrapper => wrapper.findAll('c-action-btn-stub').at(0);
const selectInstructionsFiltersField = wrapper => wrapper.find('alarms-list-remediation-instructions-filters-stub');
const selectRemoveHistoryButton = wrapper => wrapper.find('v-chip-stub');
const selectPagination = wrapper => wrapper.find('c-pagination-stub');
const selectAlarmsExpandPanelTour = wrapper => wrapper.find('alarms-expand-panel-tour-stub');

describe('alarms-list', () => {
  const $popups = mockPopups();
  const $modals = mockModals();

  const nowTimestamp = 1386435600000;
  const nowUnix = 1386435600;
  const nowSubtractOneYearUnix = 1354899600;

  mockDateNow(nowTimestamp);

  const totalItems = 10;
  const alarms = fakeStaticAlarms({
    totalItems,
    timestamp: nowTimestamp,
  });

  const userPreferences = {
    content: {
      isCorrelationEnabled: false,
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
  const exportAlarmFile = 'exportAlarmFile';
  const widget = {
    ...generateDefaultAlarmListWidgetForm(),

    _id: '880c5d0c-3f31-477c-8365-2f90389326cc',
  };
  const defaultQuery = {
    filters: [],
    active_columns: widget.parameters.widgetColumns.map(v => v.value),
    correlation: userPreferences.content.isCorrelationEnabled,
    category: userPreferences.content.category,
    limit: userPreferences.content.itemsPerPage,
    tstart: QUICK_RANGES.last1Year.start,
    tstop: QUICK_RANGES.last1Year.stop,
    opened: null,
    search: 'search',
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
  const fetchAlarmsListCsvFile = jest.fn().mockReturnValue(exportAlarmFile);
  const fetchAlarmDetails = jest.fn();
  const fetchAlarmsDetailsList = jest.fn();
  const updateAlarmDetailsQuery = jest.fn();
  const removeAlarmDetailsQuery = jest.fn();
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
      getExportByWidgetId: () => () => ({}),
    },
    actions: {
      fetchList: fetchAlarmsList,
      createAlarmsListExport,
      fetchAlarmsListExport,
      fetchAlarmsListCsvFile,
    },
  };

  const store = createMockedStoreModules([
    alarmModule,
    sideBarModule,
    infoModule,
    queryModule,
    viewModule,
    userPreferenceModule,
    authModule,
  ]);

  afterEach(() => {
    fetchUserPreference.mockClear();
    updateUserPreference.mockClear();
    updateView.mockClear();
    updateQuery.mockClear();
    hideSideBar.mockClear();
  });

  it('Query updated after mount', async () => {
    factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    expect(fetchUserPreference).toBeCalledWith(
      expect.any(Object),
      { id: widget._id },
      undefined,
    );

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...omit(defaultQuery, ['search', 'tstart', 'tstop', 'filters']),
          multiSortBy: [],
          page: 1,
          with_instructions: true,
          with_links: true,
          opened: true,
        },
      },
      undefined,
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

    expect(fetchUserPreference).not.toBeCalled();

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...omit(defaultQuery, ['search', 'tstart', 'tstop', 'filters']),
          multiSortBy: [],
          page: 1,
          with_instructions: true,
          with_links: true,
          opened: true,
        },
      },
      undefined,
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
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.correlation]: { actions: [] },
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

    const correlationField = selectCorrelationField(wrapper);

    correlationField.vm.$emit('change', !userPreferences.content.isCorrelationEnabled);

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
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          correlation: !userPreferences.content.isCorrelationEnabled,
        },
      },
      undefined,
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
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.userFilter]: { actions: [] },
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

    filterSelectorField.vm.$emit('input', selectedFilter._id);

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
      undefined,
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
      undefined,
    );
  });

  it('Filter not updated after trigger filter field without access', async () => {
    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        authModule,
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

    filterSelectorField.vm.$emit('input', selectedFilter._id);

    await flushPromises();

    expect(updateUserPreference).not.toHaveBeenCalled();
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
      undefined,
    );
  });

  it('Instruction filters updated after trigger filter field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const manualInstructionFilter = {
      manual: true,
      instructions: [{
        _id: 'manual-instruction-id',
      }],
      _id: 'id1',
    };
    const autoInstructionFilter = {
      auto: true,
      instructions: [{
        _id: 'auto-instruction-id',
      }],
      _id: 'id2',
    };
    const allAndWithInstructionFilter = {
      all: true,
      with: true,
      instructions: [{
        _id: 'all-and-with-instruction-id',
      }, {
        _id: 'all-instruction-id',
      }],
      _id: 'id3',
    };

    const newRemediationInstructionsFilters = [
      manualInstructionFilter,
      autoInstructionFilter,
      allAndWithInstructionFilter,
    ];
    const excludeInstructionsIds = [
      autoInstructionFilter.instructions[0]._id,
      manualInstructionFilter.instructions[0]._id,
    ];
    const includeInstructionsIds = [
      allAndWithInstructionFilter.instructions[0]._id,
      allAndWithInstructionFilter.instructions[1]._id,
    ];

    const instructionsFiltersField = selectInstructionsFiltersField(wrapper);

    instructionsFiltersField.vm.$emit('update:filters', newRemediationInstructionsFilters);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            remediationInstructionsFilters: newRemediationInstructionsFilters,
          },
        },
      },
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          include_instruction_types: [REMEDIATION_INSTRUCTION_TYPES.manual, REMEDIATION_INSTRUCTION_TYPES.auto],
          exclude_instruction_types: [REMEDIATION_INSTRUCTION_TYPES.manual, REMEDIATION_INSTRUCTION_TYPES.auto],
          exclude_instructions: excludeInstructionsIds,
          include_instructions: includeInstructionsIds,
          page: 1,
        },
      },
      undefined,
    );
  });

  it('Locked instruction filters updated after trigger filter field', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const manualInstructionFilter = {
      manual: true,
      instructions: [{
        _id: 'manual-instruction-id',
      }],
      _id: 'id1',
    };
    const autoInstructionFilter = {
      auto: true,
      disabled: true,
      instructions: [{
        _id: 'auto-instruction-id',
      }],
      _id: 'id2',
    };
    const disabledFilters = [autoInstructionFilter._id];

    const newRemediationInstructionsFilters = [
      manualInstructionFilter,
      autoInstructionFilter,
    ];
    const excludeInstructionsIds = [
      manualInstructionFilter.instructions[0]._id,
    ];

    const instructionsFiltersField = selectInstructionsFiltersField(wrapper);

    instructionsFiltersField.vm.$emit('update:locked-filters', newRemediationInstructionsFilters);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            disabledWidgetRemediationInstructionsFilters: disabledFilters,
          },
        },
      },
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          exclude_instruction_types: [REMEDIATION_INSTRUCTION_TYPES.manual],
          exclude_instructions: excludeInstructionsIds,
          page: 1,
        },
      },
      undefined,
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

    removeHistoryButton.vm.$emit('input');

    await flushPromises();

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: omit(defaultQuery, ['tstart', 'tstop']),
      },
      undefined,
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

    liveReportingButton.vm.$emit('click');

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
      tstart: QUICK_RANGES.last3Hour.start,
      tstop: QUICK_RANGES.last3Hour.stop,
    };

    modalArguments.config.action(actionValue);

    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          ...actionValue,
        },
      },
      undefined,
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
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.category]: { actions: [] },
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

    categoryField.vm.$emit('input', newCategory);

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
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          category: newCategory._id,
        },
      },
      undefined,
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

    const tablePagination = selectTablePaginationField(wrapper);

    const newLimit = Faker.datatype.number();

    tablePagination.vm.$emit('update:rows-per-page', newLimit);

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
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          limit: newLimit,
        },
      },
      undefined,
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

    const tablePagination = selectTablePaginationField(wrapper);

    const newPage = Faker.datatype.number();

    tablePagination.vm.$emit('update:page', newPage);

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
      undefined,
    );
  });

  it('Widget exported after trigger export button', async () => {
    const nowDate = new Date(nowTimestamp);
    const OriginalDate = Date;
    const dateSpy = jest
      .spyOn(global, 'Date')
      .mockImplementation(() => new OriginalDate(nowTimestamp));
    jest.useFakeTimers('legacy');

    const wrapper = factory({
      store: createMockedStoreModules([
        alarmModule,
        sideBarModule,
        infoModule,
        queryModule,
        viewModule,
        userPreferenceModule,
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
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

    exportButton.vm.$emit('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          filters: defaultQuery.filters,
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          opened: defaultQuery.opened,
          tstart: nowSubtractOneYearUnix,
          tstop: nowUnix,
          fields: widget.parameters.widgetExportColumns.map(({ label, value }) => ({
            label,
            name: value,
          })),
          separator: widget.parameters.exportCsvSeparator,
          time_format: widget.parameters.exportCsvDatetimeFormat,
        },
        widgetId: widget._id,
      },
      undefined,
    );

    await flushPromises();

    jest.runAllTimers();

    expect(fetchAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
      undefined,
    );

    await flushPromises();

    expect(fetchAlarmsListCsvFile).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
      undefined,
    );

    expect(saveAs).toHaveBeenCalledWith(
      expect.any(Blob),
      `${widget._id}-${nowDate.toLocaleString()}.csv`,
    );

    jest.useRealTimers();
    dateSpy.mockClear();
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
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
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

    exportButton.vm.$emit('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          filters: defaultQuery.filters,
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          opened: defaultQuery.opened,
          tstart: nowSubtractOneYearUnix,
          tstop: nowUnix,
          fields: widget.parameters.widgetExportColumns.map(({ label, value }) => ({
            label,
            name: value,
          })),
          separator: widget.parameters.exportCsvSeparator,
          time_format: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
        },
        widgetId: widget._id,
      },
      undefined,
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
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
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

    exportButton.vm.$emit('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          filters: defaultQuery.filters,
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          opened: defaultQuery.opened,
          tstart: nowSubtractOneYearUnix,
          tstop: nowUnix,
          fields: widget.parameters.widgetColumns.map(({ label, value }) => ({
            label,
            name: value,
          })),
          separator: widget.parameters.exportCsvSeparator,
          time_format: widget.parameters.exportCsvDatetimeFormat,
        },
        widgetId: widget._id,
      },
      undefined,
    );

    wrapper.destroy();
  });

  it('Widget exported after trigger export button with long request time', async () => {
    jest.useFakeTimers('legacy');

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
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport,
            fetchAlarmsListExport: longFetchAlarmsListExport,
            fetchAlarmsListCsvFile,
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
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

    exportButton.vm.$emit('click');

    await flushPromises();

    jest.runAllTimers();

    expect(longFetchAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
      undefined,
    );

    await flushPromises();

    jest.runAllTimers();

    expect(longFetchAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: exportAlarmData._id,
        widgetId: widget._id,
      },
      undefined,
    );

    await flushPromises();

    expect(saveAs).toHaveBeenCalled();

    jest.useRealTimers();
  });

  it('Error popup showed exported after trigger export button with failed create export', async () => {
    const rejectValue = { error: 'Create error' };

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
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport: jest.fn().mockRejectedValue(rejectValue),
            fetchAlarmsListExport,
            fetchAlarmsListCsvFile,
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    const exportButton = selectExportButton(wrapper);

    exportButton.vm.$emit('click');

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: rejectValue.error,
    });
  });

  it('Error popup showed exported after trigger export button with failed fetch export', async () => {
    jest.useFakeTimers('legacy');

    const rejectValue = { error: 'Fetch error' };

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
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport,
            fetchAlarmsListExport: jest.fn().mockRejectedValue(rejectValue),
            fetchAlarmsListCsvFile,
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    const exportButton = selectExportButton(wrapper);

    exportButton.vm.$emit('click');

    await flushPromises();

    jest.runAllTimers();

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: rejectValue.error,
    });

    jest.useRealTimers();
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
        {
          ...alarmModule,
          actions: {
            createAlarmsListExport,
            fetchAlarmsListExport: jest.fn().mockReturnValue(exportFailedAlarmData),
            fetchAlarmsListCsvFile,
          },
        },
        {
          ...authModule,
          getters: {
            currentUser: {},
            currentUserPermissionsById: {
              [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: { actions: [] },
            },
          },
        },
      ]),
      propsData: {
        widget,
      },
    });

    const exportButton = selectExportButton(wrapper);

    exportButton.vm.$emit('click');

    await flushPromises();

    jest.runAllTimers();

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'Something went wrong...',
    });

    jest.useRealTimers();
  });

  it('Query updated after trigger pagination', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    const pagination = selectPagination(wrapper);

    const newPage = Faker.datatype.number();

    pagination.vm.$emit('input', newPage);

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
      undefined,
    );
  });

  it('First alarm expanded after click on the prev step with first step', async () => {
    const wrapper = factory({
      store,
      propsData: {
        widget,
      },
    });

    const alarmsExpandPanelTour = selectAlarmsExpandPanelTour(wrapper);

    alarmsExpandPanelTour.vm.callbacks.onPreviousStep(1);
  });

  it('First alarm not expanded after click on the next step with already expanded alarm', async () => {
    const expanded = {
      [alarms[0]._id]: true,
    };
    const wrapper = factory({
      store,
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({
            expanded,
          }),
        },
      },
      propsData: {
        widget,
      },
    });

    const alarmsExpandPanelTour = selectAlarmsExpandPanelTour(wrapper);

    alarmsExpandPanelTour.vm.callbacks.onNextStep();

    expect(expanded).toBe(expanded);
  });

  it('First alarm not expanded after click on the prev step with second step', async () => {
    const expanded = {};
    const wrapper = factory({
      store,
      stubs: {
        ...stubs,
        'alarms-list-table': {
          template: '<div />',
          data: () => ({
            expanded,
          }),
        },
      },
      propsData: {
        widget,
      },
    });

    const alarmsExpandPanelTour = selectAlarmsExpandPanelTour(wrapper);

    alarmsExpandPanelTour.vm.callbacks.onPreviousStep(2);

    expect(expanded).toEqual({
      [alarms[0]._id]: true,
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
          ...defaultQuery,

          tstart: expect.any(Number),
          tstop: expect.any(Number),
        },
      },
      undefined,
    );
  });

  it('Alarms fetched after change query nonce', async () => {
    const firstAlarmId = alarms[0]._id;

    const expanded = {
      'non-exist-id': true,
      [firstAlarmId]: true,
    };
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
          ...defaultQuery,

          tstart: expect.any(Number),
          tstop: expect.any(Number),
        },
      },
      undefined,
    );
    expect(expanded).toEqual({
      'non-exist-id': false,
      [firstAlarmId]: true,
    });
  });

  it('Periodic started after mount with enabled value', async () => {
    jest.useFakeTimers();

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

    jest.runTimersToTime(1000);

    expect(fetchAlarmsList).toHaveBeenCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: {
          ...defaultQuery,

          tstart: expect.any(Number),
          tstop: expect.any(Number),
        },
      },
      undefined,
    );

    jest.useRealTimers();
  });

  it('Interval cleared after update periodic refresh', async () => {
    jest.useFakeTimers();

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

    jest.useRealTimers();
  });

  it('Interval cleared after destroy', async () => {
    jest.useFakeTimers();

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

    jest.useRealTimers();
  });

  it('Renders `alarms-list` with default props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
