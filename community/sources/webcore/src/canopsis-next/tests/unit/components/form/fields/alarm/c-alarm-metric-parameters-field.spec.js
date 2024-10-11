import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ALARM_METRIC_PARAMETERS } from '@/constants';

import CAlarmMetricParametersField from '@/components/forms/fields/kpi/c-alarm-metric-parameters-field.vue';

const stubs = {
  'v-autocomplete': createSelectInputStub('v-autocomplete'),
};

const selectAutocompleteNode = wrapper => wrapper.vm.$children[0];

describe('c-alarm-metric-parameters-field', () => {
  const factory = generateShallowRenderer(CAlarmMetricParametersField, { stubs });
  const snapshotFactory = generateRenderer(CAlarmMetricParametersField);

  it('Value changed after trigger the input', () => {
    const wrapper = factory({
      propsData: {
        value: [],
      },
    });

    selectAutocompleteNode(wrapper).$emit('input', ALARM_METRIC_PARAMETERS.maxAck);

    expect(wrapper).toEmit('input', ALARM_METRIC_PARAMETERS.maxAck);
  });

  it('Renders `c-alarm-metric-parameters-field` with default props', () => {
    snapshotFactory({
      propsData: {
        value: [ALARM_METRIC_PARAMETERS.createdAlarms],
      },
    });

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `c-alarm-metric-parameters-field` with custom props', () => {
    snapshotFactory({
      propsData: {
        value: [ALARM_METRIC_PARAMETERS.createdAlarms, ALARM_METRIC_PARAMETERS.ratioInstructions],
        min: 2,
        name: 'customName',
        hideDetails: true,
        parameters: [
          ALARM_METRIC_PARAMETERS.createdAlarms,
          ALARM_METRIC_PARAMETERS.ratioInstructions,
          ALARM_METRIC_PARAMETERS.ratioTickets,
        ],
        disabledParameters: [
          ALARM_METRIC_PARAMETERS.ratioTickets,
        ],
      },
    });

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `c-alarm-metric-parameters-field` with external metrics', async () => {
    snapshotFactory({
      propsData: {
        value: [ALARM_METRIC_PARAMETERS.createdAlarms, ALARM_METRIC_PARAMETERS.ratioInstructions],
        parameters: [],
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});
