import { mockModals } from '@unit/utils/mock-hooks';
import {
  createActiveViewModule,
  createAlarmModule,
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ENTITY_TYPES, MAP_TYPES, MODALS, USER_PERMISSIONS } from '@/constants';

import MapWidget from '@/components/widgets/map/map.vue';

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

describe('map', () => {
  const $modals = mockModals();
  const $periodicRefresh = {
    registerChild: vi.fn(),
    unregisterChild: vi.fn(),
  };

  const mapId = Faker.datatype.string();
  const tabId = Faker.datatype.string();
  const widgetId = Faker.datatype.string();
  const alarmsColumns = [{ value: 'column' }];

  const widget = {
    _id: widgetId,
    parameters: {
      map: mapId,
      alarmsColumns,
    },
  };

  const fetchMapStateWithoutStore = vi.fn().mockReturnValue({
    type: MAP_TYPES.geo,
    _id: mapId,
  });
  const mapModule = {
    name: 'map',
    actions: {
      fetchItemStateWithoutStore: fetchMapStateWithoutStore,
    },
  };

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { activeViewModule, registerEditingOffHandler, unregisterEditingOffHandler } =
    createActiveViewModule();
  const { alarmModule, fetchOpenAlarmsListWithoutStore } = createAlarmModule();

  const filters = [{ _id: 'filter' }];
  const { userPreferenceModule, updateUserPreference, getUserPreferenceByWidgetId } =
    createUserPreferenceModule();
  const { serviceModule, fetchServiceAlarmsWithoutStore } = createServiceModule();
  const { queryModule } = createQueryModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    mapModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
  ]);

  const factory = generateShallowRenderer(MapWidget, {
    stubs,
    propsData: {
      visible: true,
    },
    parentComponent: {
      provide: {
        $periodicRefresh,
      },
    },
  });

  afterEach(() => {
    $periodicRefresh.registerChild.mockClear();
    $periodicRefresh.unregisterChild.mockClear();
  });

  it('Register and unregister editing off handler is working', async () => {
    const wrapper = factory({
      propsData: {
        tabId,
        widget,
      },
      store,
    });

    await flushPromises();

    expect(registerEditingOffHandler).toHaveBeenCalledTimes(1);
    expect($periodicRefresh.registerChild).toHaveBeenCalledTimes(1);

    wrapper.destroy();

    expect(unregisterEditingOffHandler).toHaveBeenCalledTimes(1);
    expect($periodicRefresh.unregisterChild).toHaveBeenCalledTimes(1);
  });

  test('Map state fetched after component mount in editing mode', async () => {
    factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    expect(fetchMapStateWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: mapId,
      params: {
        category: undefined,
        filters: [],
      },
    });
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
        widget,
      },
    });

    await flushPromises();

    fetchMapStateWithoutStore.mockClear();

    const nextMap = { _id: Faker.datatype.string() };

    wrapper.vm.showMap(nextMap);

    await flushPromises();

    expect(fetchMapStateWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: nextMap._id,
      params: {
        category: undefined,
        filters: [],
      },
    });
  });

  test('Map changed after breadcrumbs was triggered', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    fetchMapStateWithoutStore.mockClear();

    const nextMap = { _id: Faker.datatype.string() };

    wrapper.vm.showMap(nextMap);

    await flushPromises();
    fetchMapStateWithoutStore.mockClear();

    wrapper.vm.backToBreadcrumb({ index: 0 });

    await flushPromises();

    expect(fetchMapStateWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: mapId,
      params: {
        category: undefined,
        filters: [],
      },
    });
  });

  test('Maps history cleared after trigger editing mode', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
    });

    await flushPromises();

    fetchMapStateWithoutStore.mockClear();

    const nextMap = { _id: Faker.datatype.string() };

    wrapper.vm.showMap(nextMap);

    await flushPromises();
    fetchMapStateWithoutStore.mockClear();

    const [, clearPreviousMaps] = registerEditingOffHandler.mock.calls[0];

    clearPreviousMaps();

    await flushPromises();

    expect(fetchMapStateWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: mapId,
      params: {
        category: undefined,
        filters: [],
      },
    });
  });

  test('Category query updated after trigger category field', async () => {
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.map.actions.category]: { actions: [] },
    });

    getUserPreferenceByWidgetId.mockReturnValue(() => ({
      content: {
        filters,
      },
    }));
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
        widget,
      },
    });

    await flushPromises();

    const category = {
      _id: Faker.datatype.string(),
    };

    wrapper.vm.updateCategory(category);

    expect(updateUserPreference).toHaveBeenCalledWith(expect.any(Object), {
      data: {
        content: {
          category: category._id,
          filters,
        },
      },
    });

    expect(wrapper.vm.query.category).toBe(category._id);
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
        widget,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const entityId = Faker.datatype.string();

    wrapper.vm.showAlarmListModal({
      entity: {
        _id: entityId,
        type: ENTITY_TYPES.service,
      },
    });

    expect($modals.show).toHaveBeenCalledWith(
      expect.objectContaining({
        name: MODALS.alarmsList,
        config: expect.objectContaining({
          widget: expect.objectContaining({
            _id: expect.any(String),
            parameters: expect.objectContaining({
              widgetColumns: expect.arrayContaining([
                expect.objectContaining({ value: alarmsColumns[0].value }),
              ]),
            }),
          }),
          fetchList: expect.any(Function),
        }),
      }),
    );

    const [modalArguments] = $modals.show.mock.calls[0];

    const params = {
      page: Faker.datatype.number(),
    };

    const response = await modalArguments.config.fetchList(params);

    expect(fetchServiceAlarmsWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: entityId,
      params,
    });

    expect(response).toEqual(serviceAlarmsResponse);
  });

  test('Opened alarms fetched correct after trigger fetch list from modal', async () => {
    const wrapper = factory({
      store,
      propsData: {
        tabId,
        widget,
      },
      mocks: {
        $modals,
      },
    });

    await flushPromises();

    const entityId = Faker.datatype.string();

    wrapper.vm.showAlarmListModal({
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

    expect(fetchOpenAlarmsListWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      params: {
        ...params,
        _id: entityId,
      },
    });
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
});
