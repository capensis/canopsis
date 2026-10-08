import { mockModals } from '@unit/utils/mock-hooks';
import {
  createAlarmTagModule,
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
  createActiveViewModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { DEFAULT_WEATHER_LIMIT } from '@/config';
import {
  COLOR_INDICATOR_TYPES,
  MODALS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  USER_PERMISSIONS,
  WIDGET_TYPES,
} from '@/constants';

import { formToWidget, widgetToForm } from '@/helpers/entities/widget/form';

import ServiceWeatherWidget from '@/components/widgets/service-weather/service-weather.vue';

vi.mock('@/stores/modules/weather-service-widget', async () => {
  const { defineStore } = await vi.importActual('pinia');
  const dispatch = (type, payload) => globalThis.__TEST_STORE__?.dispatch(type, payload);
  const getter = (type, fallback) => globalThis.__TEST_STORE__?.getters?.[type] ?? fallback;

  return {
    useWeatherServiceWidgetStore: defineStore('weatherServiceWidget', {
      getters: {
        getListByWidgetId: () => getter('service/getListByWidgetId', () => []),
        getPendingByWidgetId: () => getter('service/getPendingByWidgetId', () => false),
        getErrorByWidgetId: () => getter('service/getErrorByWidgetId', () => undefined),
      },
      actions: {
        fetchList: payload => dispatch('service/fetchList', payload),
      },
    }),
  };
});

const stubs = {
  'c-entity-category-field': true,
  'filter-selector': true,
  'filters-list-btn': true,
  'c-enabled-field': true,
  'service-weather-item': true,
  'c-help-icon': true,
};

const selectEntityCategoryField = wrapper => wrapper.find('c-entity-category-field-stub');
const selectFilterSelectorField = wrapper => wrapper.find('filter-selector-stub');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectServiceWeatherItemByIndex = (wrapper, index) =>
  wrapper.findAll('service-weather-item-stub').at(index);

describe('service-weather', () => {
  const $modals = mockModals();
  const tabId = Faker.datatype.string();

  const defaultQuery = {
    category: undefined,
    filter: undefined,
    lockedFilter: null,
    sortDesc: [],
    sortBy: [],
    itemsPerPage: DEFAULT_WEATHER_LIMIT,
    with_tag_colors: true,
    hide_grey: false,
  };

  const widget = {
    ...formToWidget(widgetToForm({ type: WIDGET_TYPES.serviceWeather })),
    _id: 'service-weather-id',
  };

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { userPreferenceModule, fetchUserPreference, updateUserPreference } =
    createUserPreferenceModule();
  const {
    serviceModule,
    getServicesListByWidgetId,

    fetchServicesList,
    fetchServiceAlarmsWithoutStore,
  } = createServiceModule();
  const { queryModule, updateQuery, getQueryById } = createQueryModule();
  const { alarmTagModule } = createAlarmTagModule();
  const { activeViewModule } = createActiveViewModule();
  const createStore = () =>
    createMockedStoreModules([
      authModule,
      userPreferenceModule,
      serviceModule,
      queryModule,
      alarmTagModule,
      activeViewModule,
    ]);

  const store = createStore();

  const factory = generateShallowRenderer(ServiceWeatherWidget, {
    stubs,
    propsData: {
      widget,
      tabId,
    },
    mocks: {
      $mq: 'l',
      $modals,
    },
  });

  beforeEach(() => {
    currentUserPermissionsById.mockReturnValue({});
    getQueryById.mockReturnValue(() => ({}));
    getServicesListByWidgetId.mockReturnValue([]);
  });

  test('Query updated after mount', async () => {
    factory({ store });

    await flushPromises();

    expect(fetchUserPreference).toBeDispatchedWith({ id: widget._id });
    expect(updateQuery).toBeDispatchedWith({
      id: widget._id,
      query: {
        ...defaultQuery,
        search: '',
      },
    });
  });

  test('Services list fetched with correct query', async () => {
    const wrapper = factory({ store });

    await wrapper.vm.fetchList();

    expect(fetchServicesList).toBeDispatchedWith({
      widgetId: widget._id,
      params: { limit: 10 },
    });
  });

  test('Services list fetched with custom limit query', async () => {
    const limit = Faker.datatype.number({ min: 0, max: 100 });
    getQueryById.mockReturnValue(() => ({
      itemsPerPage: limit,
    }));

    const wrapper = factory({
      store: createStore(),
    });

    await flushPromises();

    await wrapper.vm.fetchList();

    expect(fetchServicesList).toBeDispatchedWith({
      widgetId: widget._id,
      params: { limit },
    });
  });

  test('Category query updated after trigger category field', async () => {
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.serviceWeather.actions.category]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
        activeViewModule,
      ]),
    });

    await flushPromises();

    updateQuery.mockClear();

    const category = {
      _id: Faker.datatype.string(),
    };

    selectEntityCategoryField(wrapper).triggerCustomEvent('update:category', category);

    expect(updateUserPreference).toHaveBeenCalledWith(expect.any(Object), {
      data: {
        content: {
          category: category._id,
        },
      },
    });

    expect(updateQuery).toHaveBeenCalledWith(expect.any(Object), {
      id: widget._id,
      query: {
        category: category._id,
      },
    });
  });

  test('Filter updated after trigger filter field', async () => {
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.serviceWeather.actions.userFilter]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
        activeViewModule,
      ]),
    });

    await flushPromises();

    updateQuery.mockClear();

    const selectedFilter = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      filter: {},
    };

    selectFilterSelectorField(wrapper).triggerCustomEvent('update:modelValue', selectedFilter._id);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(expect.any(Object), {
      data: {
        content: {
          mainFilter: selectedFilter._id,
        },
      },
    });
    expect(updateQuery).toHaveBeenCalledWith(expect.any(Object), {
      id: widget._id,
      query: {
        filter: selectedFilter._id,
      },
    });
  });

  test('Gray filter updated after trigger enabled field', async () => {
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.serviceWeather.actions.userFilter]: { actions: [] },
    });
    const wrapper = factory({
      store,
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            isHideGrayEnabled: true,
          },
        },
      },
    });

    await flushPromises();

    updateQuery.mockClear();

    selectEnabledField(wrapper).triggerCustomEvent('update:modelValue', true);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(expect.any(Object), {
      data: {
        content: {
          hide_grey: true,
        },
      },
    });
    expect(updateQuery).toHaveBeenCalledWith(expect.any(Object), {
      id: widget._id,
      query: {
        hide_grey: true,
      },
    });
  });

  test('Alarms list modal showed after click on button', async () => {
    const service = {
      name: Faker.datatype.string(),
    };

    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValue([service]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
        activeViewModule,
      ]),
      propsData: {
        widget,
      },
    });

    selectServiceWeatherItemByIndex(wrapper, 0).triggerCustomEvent('show:alarms');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.alarmsList,
      config: expect.objectContaining({
        widget: expect.objectContaining({
          _id: expect.any(String),
          parameters: expect.objectContaining({
            serviceDependenciesColumns: widget.parameters.serviceDependenciesColumns,
            widgetColumns: widget.parameters.alarmsList.widgetColumns,
          }),
        }),
        title: 'modals.alarmsList.prefixTitle',
        fetchList: expect.any(Function),
      }),
    });

    const [modalArguments] = $modals.show.mock.calls[0];
    const params = { param: Faker.datatype.string() };

    await modalArguments.config.fetchList(params);

    expect(fetchServiceAlarmsWithoutStore).toHaveBeenCalledWith(expect.any(Object), {
      id: service._id,
      params,
    });
  });

  test('Alarms list modal showed after click on card', async () => {
    const service = {
      name: Faker.datatype.string(),
    };
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValue([service]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
        activeViewModule,
      ]),
      propsData: {
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList,
          },
        },
      },
    });

    await flushPromises();

    await selectServiceWeatherItemByIndex(wrapper, 0).triggerCustomEvent('show:service');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.alarmsList,
      config: expect.any(Object),
    });
  });

  test('Root cause modal showed after trigger show:root-cause event', async () => {
    const service = {
      name: Faker.datatype.string(),
    };
    getServicesListByWidgetId.mockReturnValue([service]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
        activeViewModule,
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    await selectServiceWeatherItemByIndex(wrapper, 0).triggerCustomEvent('show:root-cause');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.entitiesRootCauseDiagram,
      config: {
        entity: service,
        colorIndicator: COLOR_INDICATOR_TYPES.state,
      },
    });
  });

  test('Main information modal showed after click on card', async () => {
    const service = {
      name: Faker.datatype.string(),
      is_grey: true,
    };
    currentUserPermissionsById.mockReturnValue({
      [USER_PERMISSIONS.business.serviceWeather.actions.moreInfos]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValue([service]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
        activeViewModule,
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    await selectServiceWeatherItemByIndex(wrapper, 0).triggerCustomEvent('show:service');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.serviceEntities,
      config: {
        titleClass: 'state-pause',
        service,
        widgetParameters: widget.parameters,
      },
    });
  });
});
