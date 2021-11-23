import Faker from 'faker';
import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { stubDateNow } from '@unit/utils/stub-hooks';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiAlarmsFilters from '@/components/other/kpi/charts/partials/kpi-alarms-filters';

const localVue = createVueInstance();

const stubs = {
  'c-quick-date-interval-field': true,
  'c-filters-field': true,
  'c-sampling-field': true,
  'c-alarm-metric-parameters-field': true,
};

const factory = (options = {}) => shallowMount(KpiAlarmsFilters, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiAlarmsFilters, {
  localVue,
  stubs,

  ...options,
});

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

  stubDateNow(nowTimestamp);

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

  it('Query changed after trigger a parameters field', async () => {
    const newParameters = [ALARM_METRIC_PARAMETERS.ticketAlarms, ALARM_METRIC_PARAMETERS.ackAlarms];
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

    const filtersField = wrapper.find('c-filters-field-stub');

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
          parameters: [ALARM_METRIC_PARAMETERS.ticketAlarms, ALARM_METRIC_PARAMETERS.ackAlarms],
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
