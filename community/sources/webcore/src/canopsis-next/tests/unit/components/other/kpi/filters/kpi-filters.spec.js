import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { CRUD_ACTIONS, MODALS, USERS_PERMISSIONS } from '@/constants';

import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';

const localVue = createVueInstance();

const stubs = {
  'kpi-filters-list': true,
};

const defaultStore = createMockedStoreModules([{
  name: 'filter',
  getters: {
    pending: false,
    items: [],
    meta: {
      total_count: 0,
    },
  },
  actions: {
    fetchList: jest.fn(),
  },
}, {
  name: 'auth',
  getters: {
    currentUserPermissionsById: {},
  },
}]);

const factory = (options = {}) => shallowMount(KpiFilters, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiFilters, {
  localVue,
  stubs,

  ...options,
});

describe('kpi-filters', () => {
  it('Filters fetched after mount', async () => {
    const fetchFilters = jest.fn();
    factory({
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: fetchFilters,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(1);
    expect(fetchFilters).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: 10,
          page: 1,
        },
      },
      undefined,
    );
  });

  it('Filters fetched after change query', async () => {
    const fetchFilters = jest.fn();
    const initialRowsPerPage = Faker.datatype.number();
    const wrapper = factory({
      data() {
        return {
          query: {
            rowsPerPage: initialRowsPerPage,
          },
        };
      },
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: fetchFilters,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
    });

    await flushPromises();

    fetchFilters.mockReset();

    const kpiFiltersListElement = wrapper.find('kpi-filters-list-stub');

    const rowsPerPage = Faker.datatype.number({ max: initialRowsPerPage });
    const page = Faker.datatype.number();

    kpiFiltersListElement.vm.$emit('update:pagination', {
      rowsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(1);
    expect(fetchFilters).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          limit: rowsPerPage,
          page,
        },
      },
      undefined,
    );
  });

  it('Patterns modal showed after trigger edit', async () => {
    const showModal = jest.fn();
    const updateFilter = jest.fn();
    const fetchWithPrevious = jest.fn();
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: jest.fn(),
          update: updateFilter,
          fetchListWithPreviousParams: fetchWithPrevious,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
      mocks: {
        $modals: {
          show: showModal,
        },
      },
    });

    await flushPromises();

    const kpiFiltersListElement = wrapper.find('kpi-filters-list-stub');

    const filter = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      entity_patterns: [{}],
    };
    kpiFiltersListElement.vm.$emit('edit', filter);

    expect(showModal).toBeCalledTimes(1);
    expect(showModal).toBeCalledWith(
      {
        name: MODALS.patterns,
        config: {
          action: expect.any(Function),
          entity: true,
          name: true,
          patterns: filter,
          title: 'modals.filter.edit.title',
        },
      },
    );

    const [modalArguments] = showModal.mock.calls[0];

    const newFilterData = {
      name: Faker.datatype.string(),
      entity_patterns: [],
    };

    await modalArguments.config.action(newFilterData);

    expect(updateFilter).toBeCalledTimes(1);
    expect(fetchWithPrevious).toBeCalledTimes(1);
    expect(updateFilter).toBeCalledWith(
      expect.any(Object),
      {
        data: newFilterData,
        id: filter._id,
      },
      undefined,
    );
  });

  it('Patterns modal showed after trigger duplicate', async () => {
    const showModal = jest.fn();
    const createFilter = jest.fn();
    const fetchWithPrevious = jest.fn();
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: jest.fn(),
          create: createFilter,
          fetchListWithPreviousParams: fetchWithPrevious,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
      mocks: {
        $modals: {
          show: showModal,
        },
      },
    });

    await flushPromises();

    const kpiFiltersListElement = wrapper.find('kpi-filters-list-stub');

    const filter = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      entity_patterns: [{}],
    };
    kpiFiltersListElement.vm.$emit('duplicate', filter);

    expect(showModal).toBeCalledTimes(1);
    expect(showModal).toBeCalledWith(
      {
        name: MODALS.patterns,
        config: {
          action: expect.any(Function),
          entity: true,
          name: true,
          patterns: {
            name: filter.name,
            entity_patterns: filter.entity_patterns,
          },
          title: 'modals.filter.duplicate.title',
        },
      },
    );

    const [modalArguments] = showModal.mock.calls[0];

    const newFilterData = {
      name: Faker.datatype.string(),
      entity_patterns: [],
    };

    await modalArguments.config.action(newFilterData);

    expect(createFilter).toBeCalledTimes(1);
    expect(fetchWithPrevious).toBeCalledTimes(1);
    expect(createFilter).toBeCalledWith(
      expect.any(Object),
      {
        data: newFilterData,
      },
      undefined,
    );
  });

  it('Patterns modal showed after trigger remove', async () => {
    const showModal = jest.fn();
    const removeFilter = jest.fn();
    const fetchWithPrevious = jest.fn();
    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 0,
          },
        },
        actions: {
          fetchList: jest.fn(),
          remove: removeFilter,
          fetchListWithPreviousParams: fetchWithPrevious,
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: {},
        },
      }]),
      mocks: {
        $modals: {
          show: showModal,
        },
      },
    });

    await flushPromises();

    const kpiFiltersListElement = wrapper.find('kpi-filters-list-stub');

    const filter = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      entity_patterns: [{}],
    };
    kpiFiltersListElement.vm.$emit('remove', filter._id);

    expect(showModal).toBeCalledTimes(1);
    expect(showModal).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );

    const [modalArguments] = showModal.mock.calls[0];

    await modalArguments.config.action();

    expect(removeFilter).toBeCalledTimes(1);
    expect(fetchWithPrevious).toBeCalledTimes(1);
    expect(removeFilter).toBeCalledWith(
      expect.any(Object),
      {
        id: filter._id,
      },
      undefined,
    );
  });

  it('Renders `kpi-filters` with default props', () => {
    const wrapper = snapshotFactory({
      store: defaultStore,
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-filters` with full permissions', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'filter',
        getters: {
          pending: false,
          items: [],
          meta: {
            total_count: 10,
          },
        },
        actions: {
          fetchList: jest.fn(),
        },
      }, {
        name: 'auth',
        getters: {
          currentUserPermissionsById: () => ({
            [USERS_PERMISSIONS.technical.kpiFilters]: {
              actions: [
                CRUD_ACTIONS.create,
                CRUD_ACTIONS.update,
                CRUD_ACTIONS.read,
                CRUD_ACTIONS.delete,
              ],
            },
          }),
        },
      }]),
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
