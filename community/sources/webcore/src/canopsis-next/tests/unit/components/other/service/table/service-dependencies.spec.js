import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';

import { mockModals } from '@unit/utils/mock-hooks';
import { ENTITY_TYPES, MODALS } from '@/constants';
import CTreeviewDataTable from '@/components/common/table/c-treeview-data-table.vue';
import ServiceDependencies from '@/components/other/service/table/service-dependencies.vue';

const localVue = createVueInstance();

const stubs = {
  'c-treeview-data-table': CTreeviewDataTable,
  'c-no-events-icon': true,
  'color-indicator-wrapper': true,
};

const factory = (options = {}) => shallowMount(ServiceDependencies, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(ServiceDependencies, {
  localVue,
  stubs,

  ...options,
});

const selectTreeviewTable = wrapper => wrapper.find('.service-dependencies');
const selectNodeExpandPanel = wrapper => wrapper.findAll('.v-treeview-node__label');
const selectMoreButton = wrapper => wrapper.find('button.v-btn');
const selectDependenciesModalButton = wrapper => wrapper.find('button.v-btn');

describe('service-dependencies', () => {
  const $modals = mockModals();
  const fetchDependenciesWithoutStore = jest.fn()
    .mockResolvedValue({
      data: [],
      meta: {},
    });
  const fetchImpactsWithoutStore = jest.fn()
    .mockResolvedValue({
      data: [],
      meta: {},
    });

  const data = [
    {
      _id: 'data-alarm-1-entity',
      name: 'Data alarm 1 entity',
      type: ENTITY_TYPES.service,
      impact_level: 5,
      impact_state: 0,
      has_impacts: false,
      cycle: false,
    },
    {
      _id: 'data-alarm-2-entity',
      name: 'Data alarm 2 entity',
      type: ENTITY_TYPES.service,
      impact_level: 1,
      impact_state: 0,
      has_impacts: false,
    },
    {
      _id: 'data-alarm-3-entity',
      name: 'Data alarm 3 entity',
      type: ENTITY_TYPES.connector,
      impact_level: 5,
      impact_state: 0,
      has_impacts: false,
    },
    {
      _id: 'data-alarm-4-entity',
      name: 'Data alarm 4 entity',
      type: ENTITY_TYPES.service,
      impact_level: 1,
      impact_state: 0,
      has_impacts: false,
    },
  ];
  const meta = {
    page: 1,
    per_page: 10,
    page_count: 1,
    total_count: 4,
  };
  const entity = {
    _id: 'root-alarm-entity',
    name: 'Root alarm entity',
    impact: [
      'data-alarm-entity-1',
      'data-alarm-entity-2',
      'data-alarm-entity-3',
      'data-alarm-entity-4',
    ],
  };

  const serviceModule = {
    name: 'service',
    actions: {
      fetchDependenciesWithoutStore,
      fetchImpactsWithoutStore,
    },
  };
  const columns = [
    {
      label: 'common.name',
      value: 'entity.name',
    },
    {
      label: 'common.type',
      value: 'type',
    },
  ];

  const store = createMockedStoreModules([
    serviceModule,
  ]);

  it('Dependencies fetched after mount', async () => {
    factory({
      store,
      propsData: {
        root: entity,
      },
    });

    await flushPromises();

    expect(fetchDependenciesWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: entity._id,
        params: {
          limit: 10,
          with_flags: true,
        },
      },
      undefined,
    );
  });

  it('Children loaded after trigger load children', async () => {
    const wrapper = factory({
      store,
      propsData: {
        root: entity,
      },
    });

    await flushPromises();

    const treeviewTable = selectTreeviewTable(wrapper);

    const [, entityWithDeps] = data;

    treeviewTable.vm.loadChildren(entityWithDeps);

    expect(fetchDependenciesWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: entityWithDeps._id,
        params: {
          limit: 10,
          with_flags: true,
        },
      },
      undefined,
    );
  });

  it('More dependencies fetched after trigger more button', async () => {
    fetchDependenciesWithoutStore.mockResolvedValueOnce({
      data: data.slice(0, 2),
      meta: {
        page: 1,
        per_page: 2,
        page_count: 2,
        total_count: 4,
      },
    });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        root: entity,
      },
    });

    await flushPromises();

    fetchDependenciesWithoutStore.mockClear();

    const expandPanel = selectNodeExpandPanel(wrapper).at(2);
    const moreButton = selectMoreButton(expandPanel);

    moreButton.trigger('click');

    await flushPromises();

    expect(fetchDependenciesWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: entity._id,
        params: {
          limit: 10,
          page: 2,
          with_flags: true,
        },
      },
      undefined,
    );
  });

  it('Dependencies modal showed after trigger dependencies button', async () => {
    fetchDependenciesWithoutStore.mockResolvedValueOnce({
      data,
      meta,
    });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        root: entity,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const expandPanel = selectNodeExpandPanel(wrapper).at(1);
    const dependenciesModalButton = selectDependenciesModalButton(expandPanel);

    dependenciesModalButton.trigger('click');

    const [, entityWithDeps] = data;

    expect($modals.show).toBeCalledWith({
      name: MODALS.serviceDependencies,
      config: {
        columns: undefined,
        impact: false,
        root: {
          ...entityWithDeps,
          impact_state: 0,
        },
      },
    });
  });

  it('Dependencies modal not showed after trigger dependencies button with impact', async () => {
    fetchDependenciesWithoutStore.mockResolvedValueOnce({
      data,
      meta,
    });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        root: entity,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const expandPanel = selectNodeExpandPanel(wrapper).at(2);
    const dependenciesModalButton = selectDependenciesModalButton(expandPanel);

    dependenciesModalButton.trigger('click');

    expect($modals.show).not.toBeCalled();
  });

  it('Renders `service-dependencies` with required props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        root: entity,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `service-dependencies` with custom props', async () => {
    fetchImpactsWithoutStore.mockResolvedValueOnce({
      data,
      meta,
    });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        includeRoot: true,
        dark: true,
        light: true,
        impact: true,
        columns: [
          {
            label: 'common.name',
            value: 'entity.name',
          },
          {
            label: 'common.type',
            value: 'entity.type',
          },
        ],
        root: entity,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `service-dependencies` with excluded root', async () => {
    fetchDependenciesWithoutStore.mockResolvedValueOnce({
      data,
      meta,
    });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        columns: [
          {
            label: 'common.name',
            value: 'entity.name',
          },
          {
            label: 'common.type',
            value: 'entity.type',
          },
        ],
        root: entity,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `service-dependencies` with pages', async () => {
    fetchDependenciesWithoutStore.mockResolvedValueOnce({
      data: data.slice(0, 2),
      meta: {
        page: 1,
        per_page: 2,
        page_count: 2,
        total_count: 4,
      },
    });
    const wrapper = snapshotFactory({
      store,
      propsData: {
        columns,
        root: entity,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});