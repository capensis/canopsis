import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { fakeStaticAlarms } from '@unit/data/alarm';
import { generateDefaultAlarmListWidget } from '@/helpers/entities';

import { ALARMS_OPENED_VALUES, SORT_ORDERS } from '@/constants';

import GroupAlarmsList from '@/components/widgets/alarm/group-alarms-list.vue';

jest.mock('file-saver', () => ({
  saveAs: jest.fn(),
}));

const localVue = createVueInstance();

const stubs = {
  'alarms-list-table': true,
  'c-table-pagination': true,
};

const factory = (options = {}) => shallowMount(GroupAlarmsList, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(GroupAlarmsList, {
  localVue,
  stubs,

  ...options,
});

const selectTablePaginationField = wrapper => wrapper.find('c-table-pagination-stub');
const selectAlarmsListTable = wrapper => wrapper.find('alarms-list-table-stub');

describe('group-alarms-list', () => {
  const nowTimestamp = 1386435600000;

  const consequencesAlarms = fakeStaticAlarms({
    totalItems: 5,
    timestamp: nowTimestamp,
  });

  const causesAlarms = fakeStaticAlarms({
    totalItems: 7,
    timestamp: nowTimestamp,
  });

  const widget = generateDefaultAlarmListWidget();

  const defaultQuery = {
    limit: 13,
    category: 'category-id',
    opened: null,
    search: 'search',
  };
  const updateQuery = jest.fn();
  const updateLockedQuery = jest.fn();
  const queryModule = {
    name: 'query',
    getters: {
      getQueryById: () => () => defaultQuery,
    },
    actions: {
      update: updateQuery,
      updateLocked: updateLockedQuery,
    },
  };
  const tabId = 'tab-id';

  const store = createMockedStoreModules([
    queryModule,
  ]);

  afterEach(() => {
    updateLockedQuery.mockClear();
    updateQuery.mockClear();
  });

  it('Query updated after mounted', async () => {
    factory({
      store,
      propsData: {
        alarm: {},
        tabId,
        widget,
      },
    });

    await flushPromises();

    expect(updateQuery).toBeCalledTimes(1);
    expect(updateQuery).toBeCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          active_columns: widget.parameters.widgetColumns.map(v => v.value),
          limit: widget.parameters.itemsPerPage,
          opened: ALARMS_OPENED_VALUES.opened,
          with_instructions: true,
          multiSortBy: [],
          page: 1,
        },
      },
      undefined,
    );
  });

  it('Query updated after change pagination on the table', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarm: {},
        tabId,
        widget,
      },
    });

    updateQuery.mockClear();

    const alarmsListTable = selectAlarmsListTable(wrapper);

    const multiSortBy = [
      {
        sortBy: Faker.datatype.string(),
        descending: false,
      },
    ];
    const newQuery = {
      sortBy: Faker.datatype.string(),
      descending: true,
      multiSortBy,
    };

    alarmsListTable.vm.$emit('update:pagination', newQuery);

    await flushPromises();

    expect(updateQuery).toBeCalledTimes(1);
    expect(updateQuery).toBeCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          sortDir: SORT_ORDERS.desc,
          sortKey: newQuery.sortBy,
          opened: ALARMS_OPENED_VALUES.all,
          multiSortBy,
        },
      },
      undefined,
    );
  });

  it('Query updated after trigger update page on the table pagination', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarm: {},
        tabId,
        widget,
      },
    });

    updateQuery.mockClear();

    const tablePagination = selectTablePaginationField(wrapper);

    const newPage = Faker.datatype.number();

    tablePagination.vm.$emit('update:page', newPage);

    await flushPromises();

    expect(updateQuery).toBeCalledWith(
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

  it('Query updated after trigger update records per page on the table pagination', async () => {
    const wrapper = factory({
      store,
      propsData: {
        alarm: {},
        tabId,
        widget,
      },
    });

    updateQuery.mockClear();

    const tablePagination = selectTablePaginationField(wrapper);

    const newLimit = Faker.datatype.number();

    tablePagination.vm.$emit('update:rows-per-page', newLimit);

    await flushPromises();

    expect(updateLockedQuery).toBeCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          limit: newLimit,
        },
      },
      undefined,
    );
  });

  it('Renders `group-alarms-list` with default props', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {},
        tabId,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `group-alarms-list` with consequences alarms', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          consequences: {
            data: consequencesAlarms,
          },
        },
        tabId,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `group-alarms-list` with causes alarms', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          causes: {
            data: causesAlarms,
          },
        },
        tabId,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `group-alarms-list` with sort key', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        {
          ...queryModule,
          getters: {
            getQueryById: () => () => ({
              ...defaultQuery,
              sortKey: 'name',
              sortDir: 'desc',
            }),
          },
        },
      ]),
      propsData: {
        alarm: {
          causes: {
            data: causesAlarms,
          },
        },
        tabId,
        widget,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `group-alarms-list` without group columns', () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        alarm: {
          causes: {
            data: causesAlarms,
          },
        },
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            widgetGroupColumns: undefined,
          },
        },
        tabId,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
