import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { ENTITY_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS, ENTITY_TYPES, MODALS } from '@/constants';

import { getWidgetColumnLabel } from '@/helpers/entities/widget/list';

import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';

const service = vi.hoisted(() => ({
  fetchServiceDependenciesWithoutStore: vi.fn(),
  fetchServiceImpactsWithoutStore: vi.fn(),
}));

vi.mock('@/stores/modules/service', async () => {
  const { defineStore } = await vi.importActual('pinia');

  return {
    useServiceStore: defineStore('service', {
      actions: service,
    }),
  };
});

const stubs = {
  'service-dependencies-show-type-field': true,
  'state-settings-summary': true,
  'service-dependencies-row': true,
};

describe('service-dependencies', () => {
  const $modals = mockModals();
  const data = [
    {
      _id: 'data-alarm-1-entity',
      name: 'Data alarm 1 entity',
      type: ENTITY_TYPES.service,
      state: 0,
      impact_level: 5,
      impact_state: 0,
      has_impacts: false,
      cycle: false,
    },
    {
      _id: 'data-alarm-2-entity',
      name: 'Data alarm 2 entity',
      type: ENTITY_TYPES.service,
      state: 1,
      impact_level: 1,
      impact_state: 0,
      depends_count: 3,
      has_impacts: false,
    },
    {
      _id: 'data-alarm-3-entity',
      name: 'Data alarm 3 entity',
      type: ENTITY_TYPES.connector,
      state: 2,
      impact_level: 5,
      impact_state: 0,
      has_impacts: false,
    },
  ];
  const meta = {
    page: 1,
    per_page: 10,
    page_count: 1,
    total_count: data.length,
  };
  const root = {
    _id: 'root-alarm-entity',
    name: 'Root alarm entity',
    impact_state: 7,
  };
  const columns = [{ value: ENTITY_FIELDS.name }, { value: ENTITY_FIELDS.type }].map(column => ({
    ...column,
    sortable: false,
    text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
  }));
  const factory = generateShallowRenderer(ServiceDependencies, {
    stubs,
    propsData: {
      columns,
      root,
    },
    mocks: {
      $modals,
    },
  });

  beforeEach(() => {
    service.fetchServiceDependenciesWithoutStore.mockResolvedValue({ data, meta });
    service.fetchServiceImpactsWithoutStore.mockResolvedValue({ data, meta });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it('fetches root dependencies after mount', async () => {
    factory();

    await flushPromises();

    expect(service.fetchServiceDependenciesWithoutStore).toHaveBeenCalledWith({
      id: root._id,
      params: {
        limit: 10,
        with_flags: true,
        define_state: false,
      },
    });
  });

  it('loads dependencies for a child entity', async () => {
    const wrapper = factory();

    await flushPromises();
    service.fetchServiceDependenciesWithoutStore.mockClear();

    await wrapper.vm.fetchDependenciesById(data[1]._id);

    expect(service.fetchServiceDependenciesWithoutStore).toHaveBeenCalledWith({
      id: data[1]._id,
      params: {
        limit: 10,
        with_flags: true,
        define_state: false,
      },
    });
  });

  it('loads the next dependencies page', async () => {
    service.fetchServiceDependenciesWithoutStore.mockResolvedValueOnce({
      data: data.slice(0, 2),
      meta: {
        page: 1,
        per_page: 2,
        page_count: 2,
        total_count: data.length,
      },
    });
    const wrapper = factory();

    await flushPromises();
    service.fetchServiceDependenciesWithoutStore.mockClear();

    await wrapper.vm.loadMore({ parentId: root._id });

    expect(service.fetchServiceDependenciesWithoutStore).toHaveBeenCalledWith({
      id: root._id,
      params: {
        page: 2,
        limit: 10,
        with_flags: true,
        define_state: false,
      },
    });
  });

  it('shows the dependencies modal for an entity with children', async () => {
    const wrapper = factory();

    await flushPromises();
    wrapper.vm.showTreeOfDependenciesModal({ entity: data[1] });

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.serviceDependencies,
      config: {
        root: data[1],
        columns,
        impact: false,
      },
    });
  });

  it('does not show the dependencies modal for an entity without children', async () => {
    const wrapper = factory();

    await flushPromises();
    wrapper.vm.showTreeOfDependenciesModal({ entity: data[2] });

    expect($modals.show).not.toHaveBeenCalled();
  });
});
