import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiAlarmsFilters from '@/components/other/kpi/charts/partials/kpi-alarms-filters';

const stubs = {
  'c-quick-date-interval-field': true,
  'c-filter-field': true,
  'c-sampling-field': true,
  'c-alarm-metric-parameters-field': true,
};

describe('kpi-alarms-filters', () => {
  const nowTimestamp = 1386435600000;
  const initialQuery = {
    sampling: SAMPLINGS.day,
    filter: null,
    interval: {
      from: QUICK_RANGES.last30Days.start,
      to: QUICK_RANGES.last30Days.stop,
    },
  };

  mockDateNow(nowTimestamp);
  /**
   * Year ago date
   */
  const nowSubtractYearTimestamp = 1354834800;

  const factory = generateShallowRenderer(KpiAlarmsFilters, { stubs });

  it('Query changed after trigger a quick interval field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const quickIntervalField = wrapper.find('c-quick-date-interval-field-stub');

    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedInterval = {
      from: start,
      to: stop,
    };

    quickIntervalField.triggerCustomEvent('input', {
      from: start,
      to: stop,
    });

    await flushPromises();

    expect(wrapper).toEmitInput({
      ...initialQuery,
      interval: expectedInterval,
    });
  });

  it('Query changed after trigger a sampling field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const samplingField = wrapper.find('c-sampling-field-stub');

    samplingField.triggerCustomEvent('input', SAMPLINGS.month);

    await flushPromises();

    expect(wrapper).toEmitInput({
      ...initialQuery,
      sampling: SAMPLINGS.month,
    });
  });

  it('Query changed after trigger a sampling field to hour with large interval diff', async () => {
    const wrapper = factory({
      propsData: {
        query: {
          ...initialQuery,
          interval: {
            from: 0,
            to: initialQuery.interval.to,
          },
        },
      },
    });

    wrapper.find('c-sampling-field-stub').triggerCustomEvent('input', SAMPLINGS.hour);

    await flushPromises();

    expect(wrapper).toEmitInput({
      ...initialQuery,
      interval: {
        from: 1354921199, // We are not using nowSubtractYearTimestamp because `today` interval value is not `now`
        to: initialQuery.interval.to,
      },
      sampling: SAMPLINGS.hour,
    });
  });

  it('Query changed after trigger a sampling field to hour with normal interval diff', async () => {
    const wrapper = factory({
      propsData: {
        query: {
          ...initialQuery,
          interval: {
            from: nowSubtractYearTimestamp + 1,
            to: initialQuery.interval.to,
          },
        },
      },
    });

    const samplingField = wrapper.find('c-sampling-field-stub');

    samplingField.triggerCustomEvent('input', SAMPLINGS.hour);

    await flushPromises();

    expect(wrapper).toEmitInput({
      ...initialQuery,
      interval: {
        from: 1354921199, // We are not using nowSubtractYearTimestamp because `today` interval value is not `now`
        to: initialQuery.interval.to,
      },
      sampling: SAMPLINGS.hour,
    });
  });

  it('Query changed after trigger a parameters field', async () => {
    const newParameters = [
      ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      ALARM_METRIC_PARAMETERS.ackAlarms,
    ];
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    wrapper.find('c-alarm-metric-parameters-field-stub').triggerCustomEvent('input', newParameters);

    await flushPromises();

    expect(wrapper).toEmitInput({
      ...initialQuery,
      parameters: newParameters,
    });
  });

  it('Query changed after trigger a filters field', async () => {
    const expectedFilter = Faker.datatype.string();
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const filtersField = wrapper.find('c-filter-field-stub');

    filtersField.triggerCustomEvent('input', expectedFilter);

    await flushPromises();

    expect(wrapper).toEmitInput({
      ...initialQuery,
      filter: expectedFilter,
    });
  });
});
