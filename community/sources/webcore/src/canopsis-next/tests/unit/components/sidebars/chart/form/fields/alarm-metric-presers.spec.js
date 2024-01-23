import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import AlarmMetricPresets from '@/components/sidebars/chart/form/fields/alarm-metric-presets.vue';

const stubs = {
  'widget-settings-item': true,
  'c-alarm-metric-presets-field': true,
};

const selectAlarmMetricPresetsField = wrapper => wrapper.find('c-alarm-metric-presets-field-stub');

describe('alarm-metric-presets', () => {
  const factory = generateShallowRenderer(AlarmMetricPresets, { stubs });
  const snapshotFactory = generateRenderer(AlarmMetricPresets, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    const newValue = [{
      metric: ALARM_METRIC_PARAMETERS.ackActiveAlarms,
    }];

    selectAlarmMetricPresetsField(wrapper).triggerCustomEvent('input', newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Renders `alarm-metric-presets` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-metric-presets` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [{
          metric: ALARM_METRIC_PARAMETERS.nonDisplayedAlarms,
        }],
        withColor: true,
        withAggregateFunction: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
