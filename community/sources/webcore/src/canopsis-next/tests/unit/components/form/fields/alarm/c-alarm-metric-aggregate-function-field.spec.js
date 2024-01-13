import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createCheckboxInputStub } from '@unit/stubs/input';
import { AGGREGATE_FUNCTIONS } from '@/constants';

import CAlarmMetricAggregateFunctionField from '@/components/forms/fields/kpi/c-alarm-metric-aggregate-function-field.vue';

const stubs = {
  'v-radio-group': createCheckboxInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('c-alarm-metric-aggregate-function-field', () => {
  const factory = generateShallowRenderer(CAlarmMetricAggregateFunctionField, { stubs });
  const snapshotFactory = generateRenderer(CAlarmMetricAggregateFunctionField);

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: AGGREGATE_FUNCTIONS.sum,
      },
    });

    selectRadioGroup(wrapper).vm.$emit('input', AGGREGATE_FUNCTIONS.max);

    expect(wrapper).toEmit('input', AGGREGATE_FUNCTIONS.max);
  });

  it('Renders `c-alarm-metric-aggregate-function-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-alarm-metric-aggregate-function-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AGGREGATE_FUNCTIONS.avg,
        label: 'Custom label',
        name: 'custom_name',
        hideDetails: true,
        aggregateFunctions: [
          AGGREGATE_FUNCTIONS.avg,
          AGGREGATE_FUNCTIONS.sum,
          AGGREGATE_FUNCTIONS.min,
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
