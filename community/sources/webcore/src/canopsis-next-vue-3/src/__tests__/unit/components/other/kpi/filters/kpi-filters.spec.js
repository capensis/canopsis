import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { MODALS } from '@/constants';

import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';

const stubs = {
  'kpi-filters-list': true,
};

const selectFiltersList = wrapper => wrapper.find('kpi-filters-list-stub');

describe('kpi-filters', () => {
  const factory = generateShallowRenderer(KpiFilters, { stubs });

  it('Filters fetched after mount', async () => {
    const fetchFilters = vi.fn();
    factory({
      store: createMockedStoreModules([
        {
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
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(1);
    expect(fetchFilters).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        limit: 10,
        page: 1,
      },
    });
  });

  it('Filters fetched after change query', async () => {
    const fetchFilters = vi.fn();
    const initialItemsPerPage = Faker.datatype.number();
    const wrapper = factory({
      data() {
        return {
          query: {
            itemsPerPage: initialItemsPerPage,
          },
        };
      },
      store: createMockedStoreModules([
        {
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
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
    });

    await flushPromises();

    fetchFilters.mockReset();

    const kpiFiltersListElement = selectFiltersList(wrapper);

    const itemsPerPage = Faker.datatype.number({ max: initialItemsPerPage });
    const page = Faker.datatype.number();

    kpiFiltersListElement.triggerCustomEvent('update:options', {
      itemsPerPage,
      page,
    });

    await flushPromises();

    expect(fetchFilters).toBeCalledTimes(1);
    expect(fetchFilters).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        limit: itemsPerPage,
        page,
      },
    });
  });

  it('Patterns modal showed after trigger edit', async () => {
    const showModal = vi.fn();
    const updateFilter = vi.fn();
    const fetchWithPrevious = vi.fn();
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: 'filter',
          getters: {
            pending: false,
            items: [],
            meta: {
              total_count: 0,
            },
          },
          actions: {
            fetchList: vi.fn(),
            update: updateFilter,
            fetchListWithPreviousParams: fetchWithPrevious,
          },
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
      mocks: {
        $modals: {
          show: showModal,
        },
      },
    });

    await flushPromises();

    const kpiFiltersListElement = selectFiltersList(wrapper);

    const filter = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      entity_pattern: [],
    };
    kpiFiltersListElement.triggerCustomEvent('edit', filter);

    expect(showModal).toBeCalledTimes(1);
    expect(showModal).toHaveBeenCalledWith({
      name: MODALS.createKpiFilter,
      config: {
        action: expect.any(Function),
        filter,
        title: 'Edit filter',
      },
    });

    const [modalArguments] = showModal.mock.calls[0];

    const newFilterData = {
      name: Faker.datatype.string(),
      entity_pattern: [],
    };

    await modalArguments.config.action(newFilterData);

    expect(updateFilter).toBeCalledTimes(1);
    expect(fetchWithPrevious).toBeCalledTimes(1);
    expect(updateFilter).toHaveBeenCalledWith(expect.any(Object), {
      data: newFilterData,
      id: filter._id,
    });
  });

  it('Patterns modal showed after trigger duplicate', async () => {
    const showModal = vi.fn();
    const createFilter = vi.fn();
    const fetchWithPrevious = vi.fn();
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: 'filter',
          getters: {
            pending: false,
            items: [],
            meta: {
              total_count: 0,
            },
          },
          actions: {
            fetchList: vi.fn(),
            create: createFilter,
            fetchListWithPreviousParams: fetchWithPrevious,
          },
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
      mocks: {
        $modals: {
          show: showModal,
        },
      },
    });

    await flushPromises();

    const kpiFiltersListElement = selectFiltersList(wrapper);

    const filter = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      entity_pattern: [[]],
    };
    kpiFiltersListElement.triggerCustomEvent('duplicate', filter);

    expect(showModal).toBeCalledTimes(1);
    expect(showModal).toHaveBeenCalledWith({
      name: MODALS.createKpiFilter,
      config: {
        action: expect.any(Function),
        filter: {
          name: filter.name,
          entity_pattern: filter.entity_pattern,
        },
        title: 'Duplicate filter',
      },
    });

    const [modalArguments] = showModal.mock.calls[0];

    const newFilterData = {
      name: Faker.datatype.string(),
      entity_pattern: [],
    };

    await modalArguments.config.action(newFilterData);

    expect(createFilter).toBeCalledTimes(1);
    expect(fetchWithPrevious).toBeCalledTimes(1);
    expect(createFilter).toHaveBeenCalledWith(expect.any(Object), {
      data: newFilterData,
    });
  });

  it('Patterns modal showed after trigger remove', async () => {
    const showModal = vi.fn();
    const removeFilter = vi.fn();
    const fetchWithPrevious = vi.fn();
    const wrapper = factory({
      store: createMockedStoreModules([
        {
          name: 'filter',
          getters: {
            pending: false,
            items: [],
            meta: {
              total_count: 0,
            },
          },
          actions: {
            fetchList: vi.fn(),
            remove: removeFilter,
            fetchListWithPreviousParams: fetchWithPrevious,
          },
        },
        {
          name: 'auth',
          getters: {
            currentUserPermissionsById: {},
          },
        },
      ]),
      mocks: {
        $modals: {
          show: showModal,
        },
      },
    });

    await flushPromises();

    const kpiFiltersListElement = selectFiltersList(wrapper);

    const filter = {
      _id: Faker.datatype.string(),
      name: Faker.datatype.string(),
      entity_pattern: [],
    };
    kpiFiltersListElement.triggerCustomEvent('remove', filter._id);

    expect(showModal).toBeCalledTimes(1);
    expect(showModal).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });

    const [modalArguments] = showModal.mock.calls[0];

    await modalArguments.config.action();

    expect(removeFilter).toBeCalledTimes(1);
    expect(fetchWithPrevious).toBeCalledTimes(1);
    expect(removeFilter).toHaveBeenCalledWith(expect.any(Object), {
      id: filter._id,
    });
  });
});
