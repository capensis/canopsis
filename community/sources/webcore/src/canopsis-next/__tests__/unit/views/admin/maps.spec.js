import { mockModals } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';
import { omit } from 'lodash';

import { MAP_TYPES, MODALS } from '@/constants';

import Maps from '@/views/admin/custom-objects/maps-page.vue';

const stubs = {
  'c-page': {
    name: 'CPage',
    emits: ['create', 'refresh'],
    template: '<div class="c-page-stub"><slot /></div>',
  },
  'maps-list': {
    name: 'MapsList',
    emits: ['edit', 'remove', 'duplicate', 'refresh', 'update:options'],
    template: '<div class="maps-list-stub" />',
  },
};

const selectPage = wrapper => wrapper.findComponent({ name: 'CPage' });
const selectMapsList = wrapper => wrapper.findComponent({ name: 'MapsList' });

describe('maps', () => {
  const $modals = mockModals();

  const fetchMapsList = vi.fn();
  const fetchMapsListWithPreviousParams = vi.fn();
  const fetchMapWithoutStore = vi.fn();
  const createMap = vi.fn();
  const updateMap = vi.fn();
  const removeMap = vi.fn();
  const mapsPending = vi.fn(() => false);
  const mapsItems = vi.fn(() => []);
  const mapsMeta = vi.fn(() => ({
    total_count: 0,
  }));
  const mapModule = {
    name: 'map',
    getters: {
      pending: mapsPending,
      items: mapsItems,
      meta: mapsMeta,
    },
    actions: {
      fetchList: fetchMapsList,
      fetchListWithPreviousParams: fetchMapsListWithPreviousParams,
      fetchItemWithoutStore: fetchMapWithoutStore,
      create: createMap,
      update: updateMap,
      remove: removeMap,
    },
  };

  const currentUserPermissionsById = vi.fn().mockReturnValue({});
  const authModule = {
    name: 'auth',
    getters: {
      currentUser: {},
      currentUserPermissionsById,
    },
  };
  const store = createMockedStoreModules([authModule, mapModule]);

  const factory = generateShallowRenderer(Maps, { stubs });

  afterEach(() => {
    currentUserPermissionsById.mockClear();
    fetchMapsListWithPreviousParams.mockClear();
    createMap.mockClear();
    updateMap.mockClear();
    removeMap.mockClear();
  });

  test('Maps fetched after mounted', async () => {
    factory({
      store,
    });

    await flushPromises();

    expect(fetchMapsList).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        with_flags: true,
        page: 1,
        limit: 10,
      },
    });
  });

  test('Maps re-fetched after trigger refresh button', async () => {
    const wrapper = factory({ store });

    await flushPromises();

    fetchMapsList.mockClear();

    const page = selectPage(wrapper);

    page.vm.$emit('refresh');
    await flushPromises();

    expect(fetchMapsList).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        with_flags: true,
        page: 1,
        limit: 10,
      },
    });
  });

  test('Choose map type modal showed after trigger create button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const page = selectPage(wrapper);

    page.vm.$emit('create');
    await flushPromises();

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.createMap,
      config: {
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    const newMap = {
      type: 'type',
      parameters: {},
    };

    modalArguments.config.action(newMap);

    expect(createMap).toHaveBeenCalledWith(expect.any(Object), {
      data: newMap,
    });
    expect(fetchMapsList).toHaveBeenCalled();
  });

  test.each(Object.values(MAP_TYPES))(
    'Edit %s map modal showed after trigger edit button',
    async (value) => {
      const map = {
        _id: Faker.datatype.string(),
        type: value,
      };
      fetchMapWithoutStore.mockReturnValue(map);
      const wrapper = factory({
        store: createMockedStoreModules([mapModule, authModule]),
        mocks: {
          $modals,
        },
      });

      await flushPromises();

      const mapsList = selectMapsList(wrapper);

      mapsList.vm.$emit('edit', { _id: map._id });
      await flushPromises();

      expect(fetchMapWithoutStore).toHaveBeenCalledWith(expect.any(Object), { id: map._id });

      const title = {
        [MAP_TYPES.geo]: 'Edit a geomap',
        [MAP_TYPES.flowchart]: 'Edit a flowchart',
        [MAP_TYPES.treeOfDependencies]: 'Edit a tree of dependencies diagram',
        [MAP_TYPES.mermaid]: 'Edit a mermaid diagram',
      }[value];

      const modal = {
        [MAP_TYPES.geo]: MODALS.createGeoMap,
        [MAP_TYPES.flowchart]: MODALS.createFlowchartMap,
        [MAP_TYPES.treeOfDependencies]: MODALS.createTreeOfDependenciesMap,
        [MAP_TYPES.mermaid]: MODALS.createMermaidMap,
      }[value];

      expect($modals.show).toHaveBeenCalledTimes(1);
      expect($modals.show).toHaveBeenCalledWith({
        name: modal,
        config: {
          map,
          title,
          action: expect.any(Function),
        },
      });
      const [modalArguments] = $modals.show.mock.calls[0];

      const newMap = {
        type: value,
        parameters: {},
      };

      modalArguments.config.action(newMap);

      expect(updateMap).toHaveBeenCalledWith(expect.any(Object), {
        id: map._id,
        data: newMap,
      });
      expect(fetchMapsList).toHaveBeenCalled();
    },
  );

  test.each(Object.values(MAP_TYPES))(
    'Duplicate %s map modal showed after trigger edit button',
    async (value) => {
      const map = {
        _id: Faker.datatype.string(),
        type: value,
      };
      fetchMapWithoutStore.mockReturnValue(map);
      const wrapper = factory({
        store: createMockedStoreModules([mapModule, authModule]),
        mocks: {
          $modals,
        },
      });

      await flushPromises();

      const mapsList = selectMapsList(wrapper);

      mapsList.vm.$emit('duplicate', { _id: map._id });
      await flushPromises();

      expect(fetchMapWithoutStore).toHaveBeenCalledWith(expect.any(Object), { id: map._id });

      const title = {
        [MAP_TYPES.geo]: 'Duplicate a geomap',
        [MAP_TYPES.flowchart]: 'Duplicate a flowchart',
        [MAP_TYPES.treeOfDependencies]: 'Duplicate a tree of dependencies diagram',
        [MAP_TYPES.mermaid]: 'Duplicate a mermaid diagram',
      }[value];

      const modal = {
        [MAP_TYPES.geo]: MODALS.createGeoMap,
        [MAP_TYPES.flowchart]: MODALS.createFlowchartMap,
        [MAP_TYPES.treeOfDependencies]: MODALS.createTreeOfDependenciesMap,
        [MAP_TYPES.mermaid]: MODALS.createMermaidMap,
      }[value];

      expect($modals.show).toHaveBeenCalledTimes(1);
      expect($modals.show).toHaveBeenCalledWith({
        name: modal,
        config: {
          map: omit(map, ['_id']),
          title,
          action: expect.any(Function),
        },
      });
      const [modalArguments] = $modals.show.mock.calls[0];

      const newMap = {
        type: value,
        parameters: {},
      };

      modalArguments.config.action(newMap);

      expect(createMap).toHaveBeenCalledWith(expect.any(Object), {
        data: newMap,
      });
      expect(fetchMapsList).toHaveBeenCalled();
    },
  );

  test('Confirmation modal showed after trigger remove map button', async () => {
    const map = {
      _id: Faker.datatype.string(),
      type: MAP_TYPES.geo,
    };
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const mapsList = selectMapsList(wrapper);

    mapsList.vm.$emit('remove', map._id);
    await flushPromises();

    expect($modals.show).toHaveBeenCalledTimes(1);
    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        action: expect.any(Function),
      },
    });
    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(removeMap).toHaveBeenCalledWith(expect.any(Object), {
      id: map._id,
    });
    expect(fetchMapsList).toHaveBeenCalled();
  });
});
