import Faker from 'faker';
import { omit } from 'lodash';

import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';

import { CRUD_ACTIONS, MAP_TYPES, MODALS, USERS_PERMISSIONS } from '@/constants';

import Maps from '@/views/admin/maps.vue';

const stubs = {
  'c-page-header': true,
  'maps-list': true,
  'c-fab-btn': true,
};

const selectFabButton = wrapper => wrapper.find('c-fab-btn-stub');
const selectMapsList = wrapper => wrapper.find('maps-list-stub');

describe('maps', () => {
  const $modals = mockModals();

  const fetchMapsList = jest.fn();
  const fetchMapsListWithPreviousParams = jest.fn();
  const fetchMapWithoutStore = jest.fn();
  const createMap = jest.fn();
  const updateMap = jest.fn();
  const removeMap = jest.fn();
  const bulkRemoveMap = jest.fn();
  const mapsPending = jest.fn(() => false);
  const mapsItems = jest.fn(() => []);
  const mapsMeta = jest.fn(() => ({
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
      bulkRemove: bulkRemoveMap,
    },
  };

  const currentUserPermissionsById = jest.fn().mockReturnValue({});
  const authModule = {
    name: 'auth',
    getters: {
      currentUser: {},
      currentUserPermissionsById,
    },
  };
  const store = createMockedStoreModules([
    authModule,
    mapModule,
  ]);

  const factory = generateShallowRenderer(Maps, { stubs });
  const snapshotFactory = generateRenderer(Maps, { stubs });

  afterEach(() => {
    currentUserPermissionsById.mockClear();
    fetchMapsListWithPreviousParams.mockClear();
    createMap.mockClear();
    updateMap.mockClear();
    removeMap.mockClear();
    bulkRemoveMap.mockClear();
  });

  test('Maps fetched after mounted', async () => {
    factory({
      store,
    });

    await flushPromises();

    expect(fetchMapsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          with_flags: true,
          page: 1,
          limit: 10,
        },
      },
    );
  });

  test('Maps re-fetched after trigger refresh button', async () => {
    const wrapper = factory({ store });

    await flushPromises();

    fetchMapsList.mockClear();

    const fabButton = selectFabButton(wrapper);

    fabButton.triggerCustomEvent('refresh');

    expect(fetchMapsList).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          with_flags: true,
          page: 1,
          limit: 10,
        },
      },
    );
  });

  test('Choose map type modal showed after trigger create button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const fabButton = selectFabButton(wrapper);

    fabButton.triggerCustomEvent('create');

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.createMap,
        config: {
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    const newMap = {
      type: 'type',
      parameters: {},
    };

    modalArguments.config.action(newMap);

    expect(createMap).toBeCalledWith(
      expect.any(Object),
      {
        data: newMap,
      },
    );
    expect(fetchMapsList).toBeCalled();
  });

  test.each(Object.values(MAP_TYPES))('Edit %s map modal showed after trigger edit button', async (value) => {
    const map = {
      _id: Faker.datatype.string(),
      type: value,
    };
    fetchMapWithoutStore.mockReturnValue(map);
    const wrapper = factory({
      store: createMockedStoreModules([
        mapModule,
        authModule,
      ]),
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const mapsList = selectMapsList(wrapper);

    await mapsList.triggerCustomEvent('edit', { _id: map._id });

    expect(fetchMapWithoutStore).toBeCalledWith(
      expect.any(Object),
      { id: map._id },
    );

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

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: modal,
        config: {
          map,
          title,
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    const newMap = {
      type: value,
      parameters: {},
    };

    modalArguments.config.action(newMap);

    expect(updateMap).toBeCalledWith(
      expect.any(Object),
      {
        id: map._id,
        data: newMap,
      },
    );
    expect(fetchMapsList).toBeCalled();
  });

  test.each(Object.values(MAP_TYPES))('Duplicate %s map modal showed after trigger edit button', async (value) => {
    const map = {
      _id: Faker.datatype.string(),
      type: value,
    };
    fetchMapWithoutStore.mockReturnValue(map);
    const wrapper = factory({
      store: createMockedStoreModules([
        mapModule,
        authModule,
      ]),
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const mapsList = selectMapsList(wrapper);

    await mapsList.triggerCustomEvent('duplicate', { _id: map._id });

    expect(fetchMapWithoutStore).toBeCalledWith(
      expect.any(Object),
      { id: map._id },
    );

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

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: modal,
        config: {
          map: omit(map, ['_id']),
          title,
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    const newMap = {
      type: value,
      parameters: {},
    };

    modalArguments.config.action(newMap);

    expect(createMap).toBeCalledWith(
      expect.any(Object),
      {
        data: newMap,
      },
    );
    expect(fetchMapsList).toBeCalled();
  });

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

    await mapsList.triggerCustomEvent('remove', map._id);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(removeMap).toBeCalledWith(
      expect.any(Object),
      {
        id: map._id,
      },
    );
    expect(fetchMapsList).toBeCalled();
  });

  test('Confirmation modal showed after trigger remove selected maps button', async () => {
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

    await mapsList.triggerCustomEvent('remove-selected', [map]);

    expect($modals.show).toBeCalledTimes(1);
    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.confirmation,
        config: {
          action: expect.any(Function),
        },
      },
    );
    const [modalArguments] = $modals.show.mock.calls[0];

    modalArguments.config.action();

    expect(bulkRemoveMap).toBeCalledWith(
      expect.any(Object),
      {
        data: [
          {
            _id: map._id,
          },
        ],
      },
    );
    expect(fetchMapsList).toBeCalled();
  });

  test('Renders `maps` without permissions', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `maps` with permissions', () => {
    currentUserPermissionsById.mockReturnValueOnce(({
      [USERS_PERMISSIONS.technical.map]: {
        actions: [
          CRUD_ACTIONS.create,
          CRUD_ACTIONS.update,
          CRUD_ACTIONS.read,
          CRUD_ACTIONS.delete,
        ],
      },
    }));
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        mapModule,
        authModule,
      ]),
    });

    expect(wrapper).toMatchSnapshot();
  });
});
