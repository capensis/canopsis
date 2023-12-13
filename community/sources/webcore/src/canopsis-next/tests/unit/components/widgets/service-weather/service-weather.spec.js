import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { MODALS, SERVICE_WEATHER_WIDGET_MODAL_TYPES, USERS_PERMISSIONS, WIDGET_TYPES } from '@/constants';
import {
  generateDefaultServiceWeatherWidget,
  generatePreparedDefaultAlarmListWidget,
} from '@/helpers/entities/widget/form';
import { DEFAULT_WEATHER_LIMIT } from '@/config';

import ServiceWeatherWidget from '@/components/widgets/service-weather/service-weather.vue';
import { mockModals } from '@unit/utils/mock-hooks';

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
const selectServiceWeatherItemByIndex = (wrapper, index) => wrapper.findAll('service-weather-item-stub').at(index);

describe('service-weather', () => {
  const $modals = mockModals();
  const tabId = Faker.datatype.string();

  const defaultQuery = {
    category: undefined,
    filter: undefined,
    lockedFilter: null,
    sortDir: null,
    sortKey: null,
    limit: DEFAULT_WEATHER_LIMIT,
    hide_grey: false,
  };

  const widget = {
    ...generateDefaultServiceWeatherWidget(),
    _id: 'service-weather-id',
  };

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const {
    userPreferenceModule,
    fetchUserPreference,
    updateUserPreference,
  } = createUserPreferenceModule();
  const {
    serviceModule,
    getServicesListByWidgetId,
    getServicesErrorByWidgetId,
    fetchServicesList,
    fetchServiceAlarmsWithoutStore,
  } = createServiceModule();
  const { queryModule, updateQuery } = createQueryModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    serviceModule,
    queryModule,
  ]);

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

  const snapshotFactory = generateRenderer(ServiceWeatherWidget, {
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

  test('Query updated after mount', async () => {
    factory({ store });

    await flushPromises();

    expect(fetchUserPreference).toBeCalledWith(
      expect.any(Object),
      { id: widget._id },
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          ...defaultQuery,
          search: '',
        },
      },
      undefined,
    );
  });

  test('Services list fetched with correct query', async () => {
    const wrapper = factory({ store });

    await wrapper.vm.fetchList();

    expect(fetchServicesList).toBeCalledWith(
      expect.any(Object),
      {
        widgetId: widget._id,
        params: { limit: 10 },
      },
      undefined,
    );
  });

  test('Category query updated after trigger category field', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.category]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
    });

    await flushPromises();

    updateQuery.mockClear();

    const category = {
      _id: Faker.datatype.string(),
    };

    selectEntityCategoryField(wrapper).vm.$emit('input', category);

    expect(updateUserPreference).toBeCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            category: category._id,
          },
        },
      },
      undefined,
    );

    expect(updateQuery).toBeCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          category: category._id,
        },
      },
      undefined,
    );
  });

  test('Filter updated after trigger filter field', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.userFilter]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
    });

    await flushPromises();

    updateQuery.mockClear();

    const selectedFilter = {
      _id: Faker.datatype.string(),
      title: Faker.datatype.string(),
      filter: {},
    };

    selectFilterSelectorField(wrapper).vm.$emit('input', selectedFilter._id);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            mainFilter: selectedFilter._id,
          },
        },
      },
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          page: 1,
          filter: selectedFilter._id,
        },
      },
      undefined,
    );
  });

  test('Gray filter updated after trigger enabled field', async () => {
    const wrapper = factory({ store });

    await flushPromises();

    updateQuery.mockClear();

    selectEnabledField(wrapper).vm.$emit('input', true);

    await flushPromises();

    expect(updateUserPreference).toHaveBeenCalledWith(
      expect.any(Object),
      {
        data: {
          content: {
            hide_grey: true,
          },
        },
      },
      undefined,
    );
    expect(updateQuery).toHaveBeenCalledWith(
      expect.any(Object),
      {
        id: widget._id,
        query: {
          hide_grey: true,
        },
      },
      undefined,
    );
  });

  test('Alarms list modal showed after click on button', async () => {
    const service = {
      name: Faker.datatype.string(),
    };

    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValueOnce([service]);

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
      propsData: {
        widget,
      },
    });

    selectServiceWeatherItemByIndex(wrapper, 0).vm.$emit('show:alarms');

    const alarmListWidget = generatePreparedDefaultAlarmListWidget();
    alarmListWidget.parameters.serviceDependenciesColumns = widget.parameters.serviceDependenciesColumns;
    alarmListWidget.parameters.widgetColumns = widget.parameters.alarmsList.widgetColumns;

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.alarmsList,
        config: {
          widget: {
            ...alarmListWidget,
            _id: expect.any(String),
          },
          title: `${service.name} - alarm list`,
          fetchList: expect.any(Function),
        },
      },
    );

    const [modalArguments] = $modals.show.mock.calls[0];
    const params = { param: Faker.datatype.string() };

    await modalArguments.config.fetchList(params);

    expect(fetchServiceAlarmsWithoutStore).toBeCalledWith(
      expect.any(Object),
      { id: service._id, params },
      undefined,
    );
  });

  test('Alarms list modal showed after click on card', async () => {
    const service = {
      name: Faker.datatype.string(),
    };
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValueOnce([service]);

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
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

    await selectServiceWeatherItemByIndex(wrapper, 0).vm.$emit('show:service');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.alarmsList,
        config: expect.any(Object),
      },
    );
  });

  test('Main information modal showed after click on card', async () => {
    const service = {
      name: Faker.datatype.string(),
      is_grey: true,
    };
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValueOnce([service]);

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
      propsData: {
        widget,
      },
    });

    await flushPromises();

    await selectServiceWeatherItemByIndex(wrapper, 0).vm.$emit('show:service');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.serviceEntities,
        config: {
          color: 'var(--v-state-pause-base)',
          service,
          widgetParameters: widget.parameters,
        },
      },
    );
  });

  test('Renders `service-weather` with default props', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
      propsData: {
        tabId: 'tab-id',
        widget: {
          _id: 'service-weather-widget-id',
          type: WIDGET_TYPES.serviceWeather,
          title: 'Default service weather',
          parameters: {},
        },
        editing: false,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  test('Renders `service-weather` with full access', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.category]: { actions: [] },
      [USERS_PERMISSIONS.business.serviceWeather.actions.userFilter]: { actions: [] },
      [USERS_PERMISSIONS.business.serviceWeather.actions.addFilter]: { actions: [] },
      [USERS_PERMISSIONS.business.serviceWeather.actions.editFilter]: { actions: [] },
    });
    getServicesListByWidgetId.mockReturnValueOnce([{}]);
    const wrapper = snapshotFactory({
      propsData: {
        tabId: 'tab-id',
        widget: {
          _id: 'service-weather-widget-id',
          type: WIDGET_TYPES.serviceWeather,
          title: 'Default service weather',
          parameters: {
            columnDesktop: 2,
            margin: {},
            isHideGrayEnabled: false,
          },
        },
        editing: false,
      },
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  test('Renders `service-weather` with errors', async () => {
    getServicesErrorByWidgetId.mockReturnValueOnce({
      name: 'Service name error',
      description: 'Service description error',
    });
    const wrapper = snapshotFactory({
      propsData: {
        tabId: 'tab-id',
        widget: {
          _id: 'service-weather-widget-id',
          type: WIDGET_TYPES.serviceWeather,
          title: 'Default service weather',
          parameters: {},
        },
        editing: false,
      },
      store: createMockedStoreModules([
        authModule,
        userPreferenceModule,
        serviceModule,
        queryModule,
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});
