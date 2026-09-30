import { generateShallowRenderer } from '@unit/utils/vue';

import { AGGREGATE_FUNCTIONS } from '@/constants';

import AlarmMetricAggregateFunction from '@/components/sidebars/chart/form/fields/alarm-metric-aggregate-function.vue';

const stubs = {
  'widget-settings-item': true,
  'c-alarm-metric-aggregate-function-field': true,
};

const selectAlarmMetricPresetsField = wrapper =>
  wrapper.find('c-alarm-metric-aggregate-function-field-stub');

describe('alarm-metric-aggregate-function', () => {
  const factory = generateShallowRenderer(AlarmMetricAggregateFunction, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const wrapper = factory({
      propsData: {
        value: AGGREGATE_FUNCTIONS.max,
      },
    });

    selectAlarmMetricPresetsField(wrapper).triggerCustomEvent('input', AGGREGATE_FUNCTIONS.min);

    expect(wrapper).toEmitInput(AGGREGATE_FUNCTIONS.min);
  });
});
