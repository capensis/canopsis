import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE, QUICK_RANGES } from '@/constants';

import AvailabilityGraphSettings from '@/components/sidebars/availability/form/fields/availability-graph-settings.vue';

const stubs = {
  'widget-settings-item': true,
  'c-enabled-field': true,
  'c-quick-date-interval-type-field': true,
  'c-availability-show-type-field': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectQuickDateIntervalTypeField = wrapper => wrapper.find('c-quick-date-interval-type-field-stub');
const selectAvailabilityShowTypeField = wrapper => wrapper.find('c-availability-show-type-field-stub');

describe('availability-graph-settings', () => {
  const factory = generateShallowRenderer(AvailabilityGraphSettings, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityGraphSettings, { stubs });

  test('Enabled field changed after trigger enabled field', () => {
    const form = {
      enabled: false,
      default_show_type: AVAILABILITY_SHOW_TYPE.percent,
      default_time_range: QUICK_RANGES.last2Days.value,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', true);

    expect(wrapper).toEmitInput({
      ...form,
      enabled: true,
    });
  });

  test('Default show type changed after trigger availability show type field', () => {
    const form = {
      enabled: true,
      default_show_type: AVAILABILITY_SHOW_TYPE.percent,
      default_time_range: QUICK_RANGES.last2Days.value,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectAvailabilityShowTypeField(wrapper).triggerCustomEvent('input', AVAILABILITY_SHOW_TYPE.duration);

    expect(wrapper).toEmitInput({
      ...form,
      default_show_type: AVAILABILITY_SHOW_TYPE.duration,
    });
  });

  test('Default time range changed after trigger quick date interval field', () => {
    const form = {
      enabled: true,
      default_show_type: AVAILABILITY_SHOW_TYPE.percent,
      default_time_range: QUICK_RANGES.last2Days.value,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectQuickDateIntervalTypeField(wrapper).triggerCustomEvent('input', QUICK_RANGES.last7Days.value);

    expect(wrapper).toEmitInput({
      ...form,
      default_time_range: QUICK_RANGES.last7Days.value,
    });
  });

  test('Renders `availability-graph-settings` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          enabled: false,
          default_show_type: AVAILABILITY_SHOW_TYPE.percent,
          default_time_range: QUICK_RANGES.last2Days.value,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-graph-settings` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          enabled: true,
          default_show_type: AVAILABILITY_SHOW_TYPE.duration,
          default_time_range: QUICK_RANGES.last7Days.value,
        },
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
