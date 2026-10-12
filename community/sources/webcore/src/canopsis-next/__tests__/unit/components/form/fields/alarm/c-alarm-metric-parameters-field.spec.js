import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import KpiAlarmMetricParametersField from '@/components/other/kpi/charts/form/fields/kpi-alarm-metric-parameters-field.vue';

const stubs = {
  'v-autocomplete': createSelectInputStub('v-autocomplete'),
};

const selectAutocomplete = wrapper => wrapper.find('select.v-autocomplete');

describe('c-alarm-metric-parameters-field', () => {
  const factory = generateShallowRenderer(KpiAlarmMetricParametersField, { stubs });

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    selectAutocomplete(wrapper).triggerCustomEvent(
      'update:modelValue',
      ALARM_METRIC_PARAMETERS.maxAck,
    );

    expect(wrapper).toEmitInput(ALARM_METRIC_PARAMETERS.maxAck);
  });
});
