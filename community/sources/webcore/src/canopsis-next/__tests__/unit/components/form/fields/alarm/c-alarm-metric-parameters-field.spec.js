import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import KpiAlarmMetricParametersField from '@/components/other/kpi/charts/form/fields/kpi-alarm-metric-parameters-field.vue';

const stubs = {
  'c-autocomplete': createSelectInputStub('c-autocomplete'),
};

const selectAutocomplete = wrapper => wrapper.find('select.c-autocomplete');

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

  it('Multiple parameters changed after selecting another metric', () => {
    const parameters = [ALARM_METRIC_PARAMETERS.createdAlarms, ALARM_METRIC_PARAMETERS.activeAlarms];
    const wrapper = factory({
      propsData: {
        value: [parameters[0]],
      },
    });

    selectAutocomplete(wrapper).triggerCustomEvent('update:modelValue', parameters);

    expect(wrapper).toEmitInput(parameters);
  });

  it('Last selected parameter cannot be deselected', () => {
    const selectedParameter = ALARM_METRIC_PARAMETERS.createdAlarms;
    const wrapper = factory({
      propsData: {
        value: [selectedParameter],
      },
    });
    const parameter = wrapper.vm.availableParameters.find(({ value }) => value === selectedParameter);

    expect(wrapper.vm.getItemProps(parameter)).toEqual({ disabled: true });
    expect(wrapper.vm.getSelectionLabel(parameter)).toBe(`alarm.metrics.${selectedParameter}`);
  });

  it('Disabled parameters remain unavailable while active parameters can be deselected', () => {
    const selectedParameters = [ALARM_METRIC_PARAMETERS.createdAlarms, ALARM_METRIC_PARAMETERS.activeAlarms];
    const disabledParameter = ALARM_METRIC_PARAMETERS.ackAlarms;
    const wrapper = factory({
      propsData: {
        value: selectedParameters,
        disabledParameters: [disabledParameter, selectedParameters[0]],
      },
    });
    const selectedItem = wrapper.vm.availableParameters.find(({ value }) => value === selectedParameters[0]);
    const disabledItem = wrapper.vm.availableParameters.find(({ value }) => value === disabledParameter);

    expect(wrapper.vm.getItemProps(selectedItem)).toEqual({ disabled: false });
    expect(wrapper.vm.getItemProps(disabledItem)).toEqual({ disabled: true });
    expect(wrapper.vm.getSelectionLabel(selectedItem)).toBe('common.parametersToDisplay');
  });
});
