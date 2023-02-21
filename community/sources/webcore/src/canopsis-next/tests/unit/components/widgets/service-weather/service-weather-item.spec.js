import flushPromises from 'flush-promises';
import Faker from 'faker';

import { createVueInstance, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createAuthModule,
  createMockedStoreModules,
  createServiceEntityModule, createServiceModule,
} from '@unit/utils/store';
import { mockModals } from '@unit/utils/mock-hooks';
import {
  COLOR_INDICATOR_TYPES,
  DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE,
  MODALS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  USERS_PERMISSIONS,
} from '@/constants';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities';
import { getImpactStateColor } from '@/helpers/color';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import ServiceWeatherItem from '@/components/widgets/service-weather/service-weather-item.vue';

const localVue = createVueInstance();

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'c-no-events-icon': true,
  'impact-state-indicator': true,
  'alarm-pbehavior-counters': true,
  'alarm-state-counters': true,
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('service-weather-item', () => {
  const $modals = mockModals();

  const widget = {
    _id: 'service-weather-id',
    parameters: {
      isPriorityEnabled: true,
      modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.both,
      serviceDependenciesColumns: [{ label: 'Label', value: 'value' }],
      counters: {
        pbehavior_enabled: true,
        state_enabled: true,
        pbehavior_types: ['pbh-type'],
        state_types: ['state-types'],
      },
      margin: {
        top: 1,
        right: 3,
        bottom: 5,
        left: 7,
      },
      heightFactor: 12,
      colorIndicator: COLOR_INDICATOR_TYPES.impactState,
      blockTemplate: DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE,
    },
  };
  const service = {
    _id: 'service-id',
    name: 'Service',
    is_action_required: true,
    impact_state: 3,
    last_update_date: 1111111111,
    counters: {
      pbh_types: [{}],
    },
  };

  const { authModule, currentUserPermissionsById } = createAuthModule();
  const { serviceModule, fetchServiceAlarmsWithoutStore } = createServiceModule();
  const { serviceEntityModule } = createServiceEntityModule();

  const store = createMockedStoreModules([
    serviceModule,
    authModule,
    serviceEntityModule,
  ]);

  const factory = generateShallowRenderer(ServiceWeatherItem, {
    localVue,
    stubs,
    attachTo: document.body,
    mocks: { $modals },
  });

  const snapshotFactory = generateRenderer(ServiceWeatherItem, {
    localVue,
    stubs,
    attachTo: document.body,
    mocks: { $modals },
  });

  test('Alarms list modal showed after click on button', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        serviceModule,
        authModule,
        serviceEntityModule,
      ]),
      propsData: {
        service,
        widget,
      },
    });

    selectButton(wrapper).vm.$emit('click', new MouseEvent('click'));

    const alarmListWidget = generatePreparedDefaultAlarmListWidget();
    alarmListWidget.parameters.serviceDependenciesColumns = widget.parameters.serviceDependenciesColumns;

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.alarmsList,
        config: {
          widget: {
            ...alarmListWidget,
            _id: expect.any(String),
          },
          title: 'Service - alarm list',
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
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
    });
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        serviceModule,
        authModule,
        serviceEntityModule,
      ]),
      propsData: {
        service,
        widget: {
          ...widget,
          parameters: {
            ...widget.parameters,
            modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList,
            blockTemplate: '<div class="custom-template"></div>',
          },
        },
      },
    });

    await flushPromises();

    await wrapper.find('.custom-template').trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.alarmsList,
        config: expect.any(Object),
      },
    );
  });

  test('Main information modal showed after click on card', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos]: { actions: [] },
    });

    const newParameters = {
      ...widget.parameters,
      blockTemplate: '<div class="custom-template"></div>',
    };
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        serviceModule,
        authModule,
        serviceEntityModule,
      ]),
      propsData: {
        service,
        widget: {
          ...widget,
          parameters: newParameters,
        },
      },
    });

    await flushPromises();

    await wrapper.find('.custom-template').trigger('click');

    expect($modals.show).toBeCalledWith(
      {
        name: MODALS.serviceEntities,
        config: {
          color: getImpactStateColor(3),
          service,
          widgetParameters: newParameters,
        },
      },
    );
  });

  test('Modal doesn\'t show after click on link', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos]: { actions: [] },
    });

    const newParameters = {
      ...widget.parameters,
      blockTemplate: '<a class="custom-template-link"></a>',
    };
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        serviceModule,
        authModule,
        serviceEntityModule,
      ]),
      propsData: {
        service,
        widget: {
          ...widget,
          parameters: newParameters,
        },
      },
    });

    await flushPromises();

    await wrapper.find('.custom-template-link').trigger('click');

    expect($modals.show).not.toBeCalledWith();
  });

  test('Renders `service-weather-item` with custom props', async () => {
    const wrapper = snapshotFactory({
      store,
      propsData: {
        service,
        widget,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-weather-item` with full access', async () => {
    currentUserPermissionsById.mockReturnValueOnce({
      [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: { actions: [] },
      [USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos]: { actions: [] },
    });
    const wrapper = snapshotFactory({
      propsData: {
        service: {
          _id: 'service-id',
        },
        widget,
      },
      store: createMockedStoreModules([
        serviceModule,
        authModule,
        serviceEntityModule,
      ]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
