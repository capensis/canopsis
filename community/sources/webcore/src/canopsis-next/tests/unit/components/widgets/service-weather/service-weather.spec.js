import flushPromises from 'flush-promises';
import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createAuthModule,
  createMockedStoreModules,
  createQueryModule,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { USERS_PERMISSIONS, WIDGET_TYPES } from '@/constants';
import { generateDefaultServiceWeatherWidget } from '@/helpers/entities';
import { DEFAULT_WEATHER_LIMIT } from '@/config';

import ServiceWeatherWidget from '@/components/widgets/service-weather/service-weather.vue';

const localVue = createVueInstance();

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

describe('service-weather', () => {
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
  } = createServiceModule();
  const { queryModule, updateQuery } = createQueryModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    serviceModule,
    queryModule,
  ]);

  const factory = generateShallowRenderer(ServiceWeatherWidget, {
    localVue,
    stubs,
    propsData: {
      widget,
      tabId,
    },
    mocks: {
      $mq: 'l',
    },
  });

  const snapshotFactory = generateRenderer(ServiceWeatherWidget, {
    localVue,
    stubs,
    propsData: {
      widget,
      tabId,
    },
    mocks: {
      $mq: 'l',
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
