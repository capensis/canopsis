import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import {
  createActiveViewModule,
  createAlarmModule,
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import {
  ENTITY_TYPES,
  MAP_TYPES,
  MODALS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

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
      alarmsColumns,
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

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { activeViewModule, registerEditingOffHandler, unregisterEditingOffHandler } = createActiveViewModule();
  const { alarmModule, fetchOpenAlarmsListWithoutStore } = createAlarmModule();

  const filters = [{ _id: 'filter' }];
  const { userPreferenceModule, updateUserPreference, getUserPreferenceByWidgetId } = createUserPreferenceModule();
  const { serviceModule, fetchServiceAlarmsWithoutStore } = createServiceModule();
  const { queryModule, updateQuery } = createQueryModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    mapModule,
    activeViewModule,
    alarmModule,
    serviceModule,
    queryModule,
  ]);

  const factory = generateShallowRenderer(MapWidget, { stubs });
  const snapshotFactory = generateRenderer(MapWidget, { stubs });

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

    wrapper.destroy();

    expect(unregisterEditingOffHandler).toHaveBeenCalledTimes(1);
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

    const alarmsListWidget = generatePreparedDefaultAlarmListWidget();

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
    fetchMapStateWithoutStore.mockReturnValue({
      type: MAP_TYPES.geo,
      _id: mapId,
    });

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
      },
      store,
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
