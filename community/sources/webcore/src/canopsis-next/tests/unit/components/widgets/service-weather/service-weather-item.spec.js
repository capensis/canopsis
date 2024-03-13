import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { COLOR_INDICATOR_TYPES, DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE } from '@/constants';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import ServiceWeatherItem from '@/components/widgets/service-weather/service-weather-item.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
  'c-no-events-icon': true,
  'impact-state-indicator': true,
  'alarm-pbehavior-counters': true,
  'alarm-state-counters': true,
  'card-with-see-alarms-btn': true,
};

const selectCard = wrapper => wrapper.find('card-with-see-alarms-btn-stub');
const selectToolbar = wrapper => wrapper.find('.service-weather-item__toolbar');
const selectToolbarButton = wrapper => selectToolbar(wrapper).find('.v-btn, v-btn-stub');

describe('service-weather-item', () => {
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

  const factory = generateShallowRenderer(ServiceWeatherItem, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  const snapshotFactory = generateRenderer(ServiceWeatherItem, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Alarms list modal showed after click on button', async () => {
    const wrapper = factory({
      propsData: {
        service,
        margin: {},
        showAlarmsButton: true,
      },
    });

    selectCard(wrapper).triggerCustomEvent('show:alarms', new MouseEvent('click'));

    expect(wrapper).toHaveBeenEmit('show:alarms');
  });

  test('Main information modal showed after click on card', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        service,
        template: '<div class="custom-template"></div>',
        showAlarmsButton: true,
      },
    });

    await flushPromises();

    await wrapper.find('.custom-template').trigger('click');

    expect(wrapper).toHaveBeenEmit('show:service');
  });

  test('Show root cause emitted after trigger click on show root cause button', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        service,
        showRootCauseByStateClick: true,
      },
    });

    await selectToolbarButton(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('show:root-cause');
  });

  test('Modal doesn\'t show after click on link', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        service,
        template: '<a class="custom-template-link" href="#123"></a>',
      },
    });

    await flushPromises();

    await wrapper.find('.custom-template-link').trigger('click');

    expect(wrapper).not.toHaveBeenEmit('show:service');
  });

  test('Renders `service-weather-item` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        service,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-weather-item` with full access', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        service: {
          _id: 'service-id',
        },
        priorityEnabled: true,
        countersSettings: {
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
        template: DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE,
        showVariablesHelpButton: true,
        showAlarmsButton: true,
        showRootCauseByStateClick: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
