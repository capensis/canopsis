import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createSelectInputStub } from '@unit/stubs/input';
import { ALARM_METRIC_PARAMETERS, KPI_RATING_CRITERIA } from '@/constants';

import KpiRatingMetricField from '@/components/other/kpi/charts/partials/kpi-rating-metric-field';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(KpiRatingMetricField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiRatingMetricField, {
  localVue,

  ...options,
});

describe('kpi-rating-metric-field', () => {
  it('Metric changed after trigger select field', () => {
    const wrapper = factory({
      propsData: {
        value: ALARM_METRIC_PARAMETERS.ackAlarms,
        criteria: {
          id: 1,
          label: KPI_RATING_CRITERIA.category,
        },
      },
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.setValue(ALARM_METRIC_PARAMETERS.instructionAlarms);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toBe(ALARM_METRIC_PARAMETERS.instructionAlarms);
  });

  it('Renders `kpi-rating-metric-field` without props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_METRIC_PARAMETERS.ackAlarms,
        criteria: {
          id: 1,
          label: KPI_RATING_CRITERIA.user,
        },
      },
      mocks: {
        $te: path => path.includes('users'),
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });
});
