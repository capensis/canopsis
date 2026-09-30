import { generateRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
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

  const factory = generateRenderer(ServiceWeatherItem, {
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
    const wrapper = factory({
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
    const wrapper = factory({
      propsData: {
        service,
        showRootCauseByStateClick: true,
      },
    });

    await selectToolbarButton(wrapper).trigger('click');

    expect(wrapper).toHaveBeenEmit('show:root-cause');
  });

  test('Modal doesn\'t show after click on link', async () => {
    const wrapper = factory({
      propsData: {
        service,
        template: '<a class="custom-template-link" href="#123"></a>',
      },
    });

    await flushPromises();

    await wrapper.find('.custom-template-link').trigger('click');

    expect(wrapper).not.toHaveBeenEmit('show:service');
  });
});
