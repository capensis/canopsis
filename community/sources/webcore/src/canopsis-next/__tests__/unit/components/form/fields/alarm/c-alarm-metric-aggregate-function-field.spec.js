import { createCheckboxInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AGGREGATE_FUNCTIONS } from '@/constants';

import AlarmMetricAggregateFunction from '@/components/sidebars/chart/form/fields/alarm-metric-aggregate-function.vue';

const stubs = {
  'widget-settings-item': {
    template: '<div class="widget-settings-item"><slot /></div>',
  },
  'c-radio-group': createCheckboxInputStub('c-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.c-radio-group');

describe('c-alarm-metric-aggregate-function-field', () => {
  const factory = generateShallowRenderer(AlarmMetricAggregateFunction, { stubs });

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        modelValue: AGGREGATE_FUNCTIONS.sum,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('update:modelValue', AGGREGATE_FUNCTIONS.max);

    expect(wrapper).toEmitInput(AGGREGATE_FUNCTIONS.max);
  });
});
