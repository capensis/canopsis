import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import AvailabilityGraphSettings from '@/components/sidebars/availability/form/fields/availability-graph-settings.vue';

const stubs = {
  'widget-settings-item': true,
  'v-expand-transition': {
    template: '<div><slot /></div>',
  },
  'c-enabled-field': true,
  'c-quick-date-interval-type-field': true,
  'availability-show-type-radio-field': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
describe('availability-graph-settings', () => {
  const factory = generateShallowRenderer(AvailabilityGraphSettings, { stubs });

  test('Enabled field changed after trigger enabled field', () => {
    const modelValue = {
      enabled: false,
      default_show_type: AVAILABILITY_SHOW_TYPE.percent,
      default_time_range: QUICK_RANGES.last2Days.value,
    };
    const wrapper = factory({
      propsData: {
        modelValue,
      },
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', true);

    expect(wrapper).toEmitInput({
      ...modelValue,
      enabled: true,
    });
  });

  test('Default show type changed after trigger availability show type field', () => {
    const modelValue = {
      enabled: true,
      default_show_type: AVAILABILITY_SHOW_TYPE.percent,
      default_time_range: QUICK_RANGES.last2Days.value,
    };
    const wrapper = factory({
      propsData: {
        modelValue,
      },
    });

    wrapper.vm.updateField('default_show_type', AVAILABILITY_SHOW_TYPE.duration);

    expect(wrapper).toEmitInput({
      ...modelValue,
      default_show_type: AVAILABILITY_SHOW_TYPE.duration,
    });
  });

  test('Default time range changed after trigger quick date interval field', () => {
    const modelValue = {
      enabled: true,
      default_show_type: AVAILABILITY_SHOW_TYPE.percent,
      default_time_range: QUICK_RANGES.last2Days.value,
    };
    const wrapper = factory({
      propsData: {
        modelValue,
      },
    });

    wrapper.vm.updateField('default_time_range', QUICK_RANGES.last7Days.value);

    expect(wrapper).toEmitInput({
      ...modelValue,
      default_time_range: QUICK_RANGES.last7Days.value,
    });
  });
});
