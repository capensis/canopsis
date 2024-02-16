import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { AGGREGATE_FUNCTIONS } from '@/constants';

import AlarmMetricAggregateFunction from '@/components/sidebars/chart/form/fields/alarm-metric-aggregate-function.vue';

const stubs = {
  'widget-settings-item': true,
  'c-alarm-metric-aggregate-function-field': true,
};

const selectAlarmMetricPresetsField = wrapper => wrapper.find('c-alarm-metric-aggregate-function-field-stub');

describe('alarm-metric-aggregate-function', () => {
  const factory = generateShallowRenderer(AlarmMetricAggregateFunction, { stubs });
  const snapshotFactory = generateRenderer(AlarmMetricAggregateFunction, { stubs });

  test('Value changed after trigger color indicator field', () => {
    const wrapper = factory({
      propsData: {
        value: AGGREGATE_FUNCTIONS.max,
      },
    });

    selectAlarmMetricPresetsField(wrapper).vm.$emit('input', AGGREGATE_FUNCTIONS.min);

    expect(wrapper).toEmit('input', AGGREGATE_FUNCTIONS.min);
  });

  test('Renders `alarm-metric-aggregate-function` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AGGREGATE_FUNCTIONS.avg,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-metric-aggregate-function` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AGGREGATE_FUNCTIONS.sum,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
