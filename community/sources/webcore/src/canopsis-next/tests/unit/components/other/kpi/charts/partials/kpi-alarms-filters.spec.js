import Faker from 'faker';
import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { mockDateNow } from '@unit/utils/mock-hooks';

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
  const snapshotFactory = generateRenderer(KpiAlarmsFilters, { stubs });

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

    quickIntervalField.vm.$emit('input', {
      from: start,
      to: stop,
    });

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.interval).toEqual(expectedInterval);
    expect(eventData).toEqual({
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

    samplingField.vm.$emit('input', SAMPLINGS.month);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.sampling).toEqual(SAMPLINGS.month);
    expect(eventData).toEqual({
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

    const samplingField = wrapper.find('c-sampling-field-stub');

    samplingField.vm.$emit('input', SAMPLINGS.hour);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.sampling).toEqual(SAMPLINGS.hour);
    expect(eventData).toEqual({
      ...initialQuery,
      interval: {
        from: nowSubtractYearTimestamp,
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

    samplingField.vm.$emit('input', SAMPLINGS.hour);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.sampling).toEqual(SAMPLINGS.hour);
    expect(eventData).toEqual({
      ...initialQuery,
      interval: {
        from: nowSubtractYearTimestamp + 1,
        to: initialQuery.interval.to,
      },
      sampling: SAMPLINGS.hour,
    });
  });

  it('Query changed after trigger a parameters field', async () => {
    const newParameters = [ALARM_METRIC_PARAMETERS.ticketActiveAlarms, ALARM_METRIC_PARAMETERS.ackAlarms];
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const metricParametersField = wrapper.find('c-alarm-metric-parameters-field-stub');

    metricParametersField.vm.$emit('input', newParameters);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.parameters).toEqual(newParameters);
    expect(eventData).toEqual({
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

    filtersField.vm.$emit('input', expectedFilter);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.filter).toEqual(expectedFilter);
    expect(eventData).toEqual({
      ...initialQuery,
      filter: expectedFilter,
    });
  });

  it('Renders `kpi-alarms-filters` with query', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        query: {
          sampling: SAMPLINGS.day,
          filter: null,
          parameters: [ALARM_METRIC_PARAMETERS.ticketActiveAlarms, ALARM_METRIC_PARAMETERS.ackAlarms],
          interval: {
            from: QUICK_RANGES.last30Days.start,
            to: QUICK_RANGES.last30Days.stop,
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-alarms-filters` with hour sampling', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        minDate: 0,
        query: {
          sampling: SAMPLINGS.hour,
          filter: null,
          parameters: [ALARM_METRIC_PARAMETERS.ticketActiveAlarms, ALARM_METRIC_PARAMETERS.ackAlarms],
          interval: {
            from: QUICK_RANGES.last30Days.start,
            to: QUICK_RANGES.last30Days.stop,
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-alarms-filters` with hour sampling and normal interval', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        minDate: nowSubtractYearTimestamp + 1,
        query: {
          sampling: SAMPLINGS.hour,
          filter: null,
          parameters: [ALARM_METRIC_PARAMETERS.ticketActiveAlarms, ALARM_METRIC_PARAMETERS.ackAlarms],
          interval: {
            from: QUICK_RANGES.last30Days.start,
            to: QUICK_RANGES.last30Days.stop,
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
