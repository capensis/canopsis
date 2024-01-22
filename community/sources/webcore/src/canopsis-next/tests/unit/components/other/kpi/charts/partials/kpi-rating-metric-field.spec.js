import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ALARM_METRIC_PARAMETERS, KPI_RATING_SETTINGS_TYPES } from '@/constants';

import KpiRatingMetricField from '@/components/other/kpi/charts/form/fields/kpi-rating-metric-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
};

const snapshotStubs = {
  'c-select-field': CSelectField,
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('kpi-rating-metric-field', () => {
  const factory = generateShallowRenderer(KpiRatingMetricField, { stubs });
  const snapshotFactory = generateRenderer(KpiRatingMetricField, { stubs: snapshotStubs });

  it('Metric changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: ALARM_METRIC_PARAMETERS.ackAlarms,
        type: KPI_RATING_SETTINGS_TYPES.entity,
      },
    });

    const valueElement = selectSelectField(wrapper);

    valueElement.setValue(ALARM_METRIC_PARAMETERS.instructionAlarms);

    expect(wrapper).toEmit('input', ALARM_METRIC_PARAMETERS.instructionAlarms);
  });

  it('Renders `kpi-rating-metric-field` without props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_METRIC_PARAMETERS.ackAlarms,
        type: KPI_RATING_SETTINGS_TYPES.user,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
