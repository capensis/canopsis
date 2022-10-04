import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance, shallowMount } from '@unit/utils/vue';
import { createMockedStoreModules } from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import { ENTITY_TYPES, MAP_TYPES, MODALS, USERS_PERMISSIONS, WIDGET_TYPES } from '@/constants';

import MapWidget from '@/components/widgets/map/map.vue';
import { generateDefaultAlarmListWidget } from '@/helpers/entities';

const localVue = createVueInstance();

const stubs = {
  'c-entity-category-field': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'map-breadcrumbs': true,
  'geomap-preview': true,
  'flowchart-preview': true,
  'mermaid-preview': true,
  'tree-of-dependencies-preview': true,
};

const factory = (options = {}) => shallowMount(MapWidget, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(MapWidget, {
  localVue,
  stubs,

  ...options,
});

const selectMapBreadcrumbs = wrapper => wrapper.find('map-breadcrumbs-stub');
const selectEntityCategoryField = wrapper => wrapper.find('c-entity-category-field-stub');
const selectMapPreview = wrapper => wrapper.find(
  [
    'geomap-preview-stub',
    'flowchart-preview-stub',
    'mermaid-preview-stub',
    'tree-of-dependencies-preview-stub',
  ].join(','),
);

describe('map', () => {
  const $modals = mockModals();

  const mapId = Faker.datatype.string();
  const tabId = Faker.datatype.string();
  const widgetId = Faker.datatype.string();
  const alarmsColumns = [{ value: 'column' }];

  const widget = {
    _id: widgetId,
    parameters: {
      map: mapId,
      alarms_columns: alarmsColumns,
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
  const fetchMapStateWithoutStore = jest.fn().mockReturnValue({
    type: MAP_TYPES.geo,
    _id: mapId,
  });
  const mapModule = {
    name: 'map',
    actions: {
      fetchItemStateWithoutStore: fetchMapStateWithoutStore,
    },
  };
  const registerEditingOffHandler = jest.fn();
  const unregisterEditingOffHandler = jest.fn();
  const activeViewModule = {
    name: 'activeView',
    actions: {
      registerEditingOffHandler,
      unregisterEditingOffHandler,
    },
  };
  const fetchOpenAlarmsListWithoutStore = jest.fn();
  const alarmModule = {
    name: 'alarm',
    actions: {
      fetchOpenAlarmsListWithoutStore,
    },
  };
  const getUserPreferenceByWidgetId = jest.fn().mockReturnValue({});
  const filters = [{ _id: 'filter' }];
  const getItemByWidgetId = jest.fn().mockReturnValue(() => ({
    content: {
      filters,
    },
  }));
  const updateUserPreference = jest.fn();
  const fetchItem = jest.fn();
  const userPreferenceModule = {
    name: 'userPreference',
    getters: {
      getItemByWidgetId,
      getUserPreferenceByWidgetId,
    },
    actions: {
      update: updateUserPreference,
      fetchItem,
    },
  };
  const fetchServiceAlarmsWithoutStore = jest.fn();
  const serviceModule = {
    name: 'service',
    actions: {
      fetchAlarmsWithoutStore: fetchServiceAlarmsWithoutStore,
    },
  };
  const getQueryById = jest.fn().mockReturnValue(() => ({}));
  const getQueryNonceById = jest.fn().mockReturnValue(() => 'nonce');
  const updateQuery = jest.fn();
  const queryModule = {
    name: 'query',
    getters: {
      getQueryById,
      getQueryNonceById,
    },
    actions: {
      update: updateQuery,
    },
  };

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    mapModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
  ]);

  beforeEach(() => {
    fetchMapStateWithoutStore.mockClear();
    registerEditingOffHandler.mockClear();
    unregisterEditingOffHandler.mockClear();
  });

  it('Register and unregister editing off handler is working', async () => {
    const wrapper = factory({
      propsData: {
        tabId,
        editing: true,
        widget,
      },
      store,
    });

    await flushPromises();

    expect(registerEditingOffHandler).toHaveBeenCalledTimes(1);

    wrapper.destroy();

    expect(unregisterEditingOffHandler).toHaveBeenCalledTimes(1);
  });

  test('Map state fetched after component mount in editing mode', async () => {
    factory({
      store,
      propsData: {
        tabId,
        editing: true,
        widget,
      },
    });

    await flushPromises();

    expect(fetchMapStateWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: mapId,
        params: {},
      },
      undefined,
    );
  });

  test('Map changed after map preview was triggered', async () => {
    fetchMapStateWithoutStore.mockReturnValue({
      type: MAP_TYPES.geo,
      _id: mapId,
    });
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        editing: true,
        widget,
      },
    });

    await flushPromises();

    fetchMapStateWithoutStore.mockClear();

    const mapPreview = selectMapPreview(wrapper);

    const nextMap = { _id: Faker.datatype.string() };

    mapPreview.vm.$emit('show:map', nextMap);

    await flushPromises();

    expect(fetchMapStateWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: nextMap._id,
        params: {},
      },
      undefined,
    );
  });

  test('Map changed after breadcrumbs was triggered', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        editing: true,
        widget,
      },
    });

    await flushPromises();

    fetchMapStateWithoutStore.mockClear();

    const nextMap = { _id: Faker.datatype.string() };

    const mapPreview = selectMapPreview(wrapper);
    mapPreview.vm.$emit('show:map', nextMap);

    await flushPromises();
    fetchMapStateWithoutStore.mockClear();

    const breadcrumbs = selectMapBreadcrumbs(wrapper);
    breadcrumbs.vm.$emit('click', { index: 0 });

    await flushPromises();

    expect(fetchMapStateWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: mapId,
        params: {},
      },
      undefined,
    );
  });

  test('Maps history cleared after trigger editing mode', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        editing: true,
        widget,
      },
    });

    await flushPromises();

    fetchMapStateWithoutStore.mockClear();

    const nextMap = { _id: Faker.datatype.string() };

    const mapPreview = selectMapPreview(wrapper);
    mapPreview.vm.$emit('show:map', nextMap);

    await flushPromises();
    fetchMapStateWithoutStore.mockClear();

    const [, clearPreviousMaps] = registerEditingOffHandler.mock.calls[0];

    clearPreviousMaps();

    await flushPromises();

    expect(fetchMapStateWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: mapId,
        params: {},
      },
      undefined,
    );
  });

  test('Category query updated after trigger category field', async () => {
    currentUserPermissionsById.mockReturnValue({
      [USERS_PERMISSIONS.business.map.actions.category]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        mapModule,
        activeViewModule,
        alarmModule,
        serviceModule,
        queryModule,
      ]),
      propsData: {
        tabId,
        editing: true,
        widget,
      },
    });

    await flushPromises();

    const entityCategoryField = selectEntityCategoryField(wrapper);

    const category = {
      _id: Faker.datatype.string(),
    };

    entityCategoryField.vm.$emit('input', category);

    expect(updateUserPreference).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            category: category._id,
            filters,
          },
        },
      },
      undefined,
    );

    expect(updateQuery).toBeCalledWith(
      expect.any(Object),
      {
        id: widgetId,
        query: {
          category: category._id,
        },
      },
      undefined,
    );
  });

  test('Alarms list modal showed after trigger map preview', async () => {
    const serviceAlarmsResponse = {
      meta: {
        total_count: 1,
      },
      data: [{ _id: Faker.datatype.string() }],
    };
    fetchServiceAlarmsWithoutStore.mockResolvedValueOnce(serviceAlarmsResponse);
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        editing: true,
        widget,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const mapPreview = selectMapPreview(wrapper);

    const entityId = Faker.datatype.string();

    mapPreview.vm.$emit('show:alarms', {
      entity: {
        _id: entityId,
        type: ENTITY_TYPES.service,
      },
    });

    const alarmsListWidget = generateDefaultAlarmListWidget();

    expect($modals.show).toBeCalledWith({
      name: MODALS.alarmsList,
      config: {
        widget: {
          ...alarmsListWidget,
          _id: expect.any(String),
          parameters: {
            ...alarmsListWidget.parameters,
            widgetColumns: alarmsColumns,
          },
        },
        fetchList: expect.any(Function),
      },
    });

    const [modalArguments] = $modals.show.mock.calls[0];

    const params = {
      page: Faker.datatype.number(),
    };

    const response = await modalArguments.config.fetchList(params);

    expect(fetchServiceAlarmsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        id: entityId,
        params,
      },
      undefined,
    );

    expect(response).toEqual(serviceAlarmsResponse);
  });

  test('Opened alarms fetched correct after trigger fetch list from modal', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        editing: true,
        widget,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const mapPreview = selectMapPreview(wrapper);

    const entityId = Faker.datatype.string();

    mapPreview.vm.$emit('show:alarms', {
      entity: {
        _id: entityId,
        type: ENTITY_TYPES.component,
      },
    });

    expect($modals.show).toBeCalled();

    const [modalArguments] = $modals.show.mock.calls[0];

    const params = {
      page: Faker.datatype.number(),
    };

    const response = await modalArguments.config.fetchList(params);

    expect(fetchOpenAlarmsListWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          ...params,
          _id: entityId,
        },
      },
      undefined,
    );
    fetchOpenAlarmsListWithoutStore.mockClear();

    expect(response).toEqual({
      meta: { total_count: 0 },
      data: [],
    });

    const alarm = { _id: Faker.datatype.string() };

    fetchOpenAlarmsListWithoutStore.mockResolvedValueOnce(alarm);

    const responseWithAlarm = await modalArguments.config.fetchList(params);

    expect(fetchOpenAlarmsListWithoutStore).toBeCalled();
    expect(responseWithAlarm).toEqual({
      meta: { total_count: 1 },
      data: [alarm],
    });
  });

  test.each(Object.values(MAP_TYPES))('Renders %s `map` with state', async (type) => {
    fetchMapStateWithoutStore.mockReturnValue({
      type,
      _id: `${type}_id`,
    });
    currentUserPermissionsById.mockReturnValue({
      [USERS_PERMISSIONS.business.map.actions.category]: { actions: [] },
      [USERS_PERMISSIONS.business.map.actions.userFilter]: { actions: [] },
      [USERS_PERMISSIONS.business.map.actions.listFilters]: { actions: [] },
      [USERS_PERMISSIONS.business.map.actions.editFilter]: { actions: [] },
      [USERS_PERMISSIONS.business.map.actions.addFilter]: { actions: [] },
    });
    const wrapper = snapshotFactory({
      propsData: {
        tabId: 'tab-id',
        widget: {
          _id: 'map-widget-id',
          type: WIDGET_TYPES.map,
          title: 'Title',
          parameters: {
            map: 'map',
            entity_info_template: 'entity_info_template',
            color_indicator: 'color_indicator',
            entities_under_pbehavior_enabled: true,
          },
        },
        editing: true,
      },
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        mapModule,
        activeViewModule,
        alarmModule,
        serviceModule,
        queryModule,
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `map` with default props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        tabId: 'tab-id',
        widget: {
          _id: 'map-widget-id',
          type: WIDGET_TYPES.map,
          title: 'Default map',
          parameters: {
            map: 'map',
          },
        },
        editing: false,
      },
      store,
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
