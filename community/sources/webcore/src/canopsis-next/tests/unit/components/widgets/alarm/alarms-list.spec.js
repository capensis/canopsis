import Faker from 'faker';
import flushPromises from 'flush-promises';
import { saveAs } from 'file-saver';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow, mockPopups } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { fakeStaticAlarms } from '@unit/data/alarm';
import { alarmListWidgetToForm } from '@/helpers/forms/widgets/alarm';
import {
  CANOPSIS_EDITION,
  EXPORT_STATUSES,
  FILTER_DEFAULT_VALUES,
  FILTER_MONGO_OPERATORS,
  USERS_PERMISSIONS,
} from '@/constants';

import AlarmsList from '@/components/widgets/alarm/alarms-list.vue';

jest.mock('file-saver', () => ({
  saveAs: jest.fn(),
}));

const localVue = createVueInstance();

const stubs = {
  'c-advanced-search-field': true,
  'c-entity-category-field': true,
  'v-switch': true,
  'filter-selector': true,
  'alarms-list-remediation-instructions-filters': true,
  'c-action-btn': true,
  'c-pagination': true,
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
  'alarms-list-remediation-instructions-filters': true,
  'c-action-btn': true,
  'c-pagination': true,
  'alarms-list-table': true,
  'c-table-pagination': true,
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
const selectExportButtonField = wrapper => wrapper.findAll('c-action-btn-stub').at(1);

describe('alarms-list', () => {
  const $popups = mockPopups();

  const nowTimestamp = 1386435600000;
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
  const widget = alarmListWidgetToForm();
  const defaultQuery = {
    active_columns: widget.parameters.widgetColumns.map(v => v.value),
    correlation: userPreferences.content.isCorrelationEnabled,
    category: userPreferences.content.category,
    limit: userPreferences.content.itemsPerPage,
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
  const createAlarmsListExport = jest.fn().mockReturnValue(exportAlarmData);
  const fetchAlarmsListExport = jest.fn().mockReturnValue(exportAlarmData);
  const fetchAlarmsListCsvFile = jest.fn().mockReturnValue(exportAlarmFile);
  const sideBarModule = {
    name: 'sideBar',
    actions: {
      hide: hideSideBar,
    },
  };
  const infoModule = {
    name: 'info',
    getters: { edition: CANOPSIS_EDITION.cat },
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
  const userPreferenceModule = {
    name: 'userPreference',
    getters: {
      getItemByWidget: () => () => userPreferences,
    },
    actions: {
      update: updateUserPreference,
      fetchItem: jest.fn(),
    },
  };
  const authModule = {
    name: 'auth',
    getters: {
      currentUser: {},
      currentUserPermissionsById: {},
    },
  };
  const alarmModule = {
    name: 'alarm',
    getters: {
      getMetaByWidgetId: () => () => ({
        total_count: totalItems,
      }),
      getListByWidgetId: () => () => alarms,
      getPendingByWidgetId: () => () => false,
      getExportByWidgetId: () => () => ({}),
    },
    actions: {
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
    updateUserPreference.mockReset();
    updateView.mockReset();
    updateQuery.mockReset();
    hideSideBar.mockReset();
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

    const filterSelectorField = selectFilterSelectorField(wrapper);

    const newFilter = {
      title: Faker.datatype.string(),
      filter: {},
    };

    filterSelectorField.vm.$emit('input', newFilter);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            mainFilter: newFilter,
            mainFilterUpdatedAt: nowTimestamp,
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
          filter: newFilter.filter,
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

    const filterSelectorField = selectFilterSelectorField(wrapper);

    const newFilter = {
      title: Faker.datatype.string(),
      filter: {},
    };

    filterSelectorField.vm.$emit('input', newFilter);

    await flushPromises();

    expect(updateUserPreference).not.toHaveBeenCalled();
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          page: 1,
          filter: newFilter.filter,
        },
      },
      undefined,
    );
  });

  it('Filter condition updated after trigger filter field', async () => {
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

    const filterSelectorField = selectFilterSelectorField(wrapper);

    filterSelectorField.vm.$emit('update:condition', FILTER_MONGO_OPERATORS.or);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            mainFilterCondition: FILTER_MONGO_OPERATORS.or,
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
        },
      },
      undefined,
    );
  });

  it('Filter condition updated after trigger filter field without value', async () => {
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

    const filterSelectorField = selectFilterSelectorField(wrapper);

    filterSelectorField.vm.$emit('update:condition');

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            mainFilterCondition: FILTER_DEFAULT_VALUES.condition,
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
        },
      },
      undefined,
    );
  });

  it('Filters updated after trigger filter field', async () => {
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

    const mainFilter = {
      title: 'main-filter',
      filter: {},
    };
    const filters = [
      {
        title: Faker.datatype.string(),
        filter: {},
      },
      mainFilter,
    ];

    const filterSelectorField = selectFilterSelectorField(wrapper);

    filterSelectorField.vm.$emit('update:filters', filters, mainFilter);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            ...userPreferences.content,
            viewFilters: filters,
            mainFilter,
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
          filter: mainFilter.filter,
          page: 1,
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
    const dateSpy = jest
      .spyOn(global, 'Date')
      .mockImplementation(() => nowDate);
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

    const exportButton = selectExportButtonField(wrapper);

    exportButton.vm.$emit('click');

    expect(createAlarmsListExport).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          search: defaultQuery.search,
          category: defaultQuery.category,
          correlation: defaultQuery.correlation,
          opened: defaultQuery.opened,
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
    dateSpy.mockReset();
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

    const exportButton = selectExportButtonField(wrapper);

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

    const exportButton = selectExportButtonField(wrapper);

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

    const exportButton = selectExportButtonField(wrapper);

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

    const exportButton = selectExportButtonField(wrapper);

    exportButton.vm.$emit('click');

    await flushPromises();

    jest.runAllTimers();

    await flushPromises();

    expect($popups.error).toHaveBeenCalledWith({
      text: 'errors.default',
    });

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
