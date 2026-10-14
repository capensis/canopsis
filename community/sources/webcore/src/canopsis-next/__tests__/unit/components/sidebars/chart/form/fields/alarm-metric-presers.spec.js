import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import AlarmMetricPresets from '@/components/sidebars/chart/form/fields/alarm-metric-presets.vue';

const stubs = {
  'widget-settings-item': true,
  'c-alarm-metric-presets-field': true,
};

const selectAlarmMetricPresetsField = wrapper => wrapper.find('c-alarm-metric-presets-field-stub');

describe('alarm-metric-presets', () => {
  const factory = generateShallowRenderer(AlarmMetricPresets, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    const newValue = [
      {
        metric: ALARM_METRIC_PARAMETERS.ackActiveAlarms,
      },
    ];

    selectAlarmMetricPresetsField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });
});
