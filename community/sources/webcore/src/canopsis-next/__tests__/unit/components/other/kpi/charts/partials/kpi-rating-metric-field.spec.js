import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS, KPI_RATING_SETTINGS_TYPES } from '@/constants';

import KpiRatingMetricField from '@/components/other/kpi/charts/form/fields/kpi-rating-metric-field.vue';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('kpi-rating-metric-field', () => {
  const factory = generateShallowRenderer(KpiRatingMetricField, { stubs });

  it('Metric changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: ALARM_METRIC_PARAMETERS.ackAlarms,
        type: KPI_RATING_SETTINGS_TYPES.entity,
      },
    });

    const valueElement = selectSelectField(wrapper);

    valueElement.setValue(ALARM_METRIC_PARAMETERS.instructionAlarms);

    expect(wrapper).toEmitInput(ALARM_METRIC_PARAMETERS.instructionAlarms);
  });
});
