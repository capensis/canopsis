import Faker from 'faker';
import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { stubDateNow } from '@unit/utils/stub-hooks';

import { ALARM_METRIC_PARAMETERS, KPI_RATING_CRITERIA, QUICK_RANGES } from '@/constants';

import KpiRatingFilters from '@/components/other/kpi/charts/partials/kpi-rating-filters';

const localVue = createVueInstance();

const stubs = {
  'c-quick-date-interval-field': true,
  'c-filters-field': true,
  'kpi-rating-criteria-field': true,
  'kpi-rating-metric-field': true,
  'c-records-per-page-field': true,
};

const factory = (options = {}) => shallowMount(KpiRatingFilters, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiRatingFilters, {
  localVue,
  stubs,

  ...options,
});

describe('kpi-rating-filters', () => {
  const nowTimestamp = 1386435600000;
  const initialQuery = {
    filter: null,
    criteria: undefined,
    metric: ALARM_METRIC_PARAMETERS.ackAlarms,
    rowsPerPage: 5,
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

  it('Query changed after trigger a criteria field', async () => {
    const ratingSetting = {
      id: 1,
      label: KPI_RATING_CRITERIA.role,
    };
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const criteriaField = wrapper.find('kpi-rating-criteria-field-stub');

    criteriaField.vm.$emit('input', ratingSetting);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.criteria).toEqual(ratingSetting);
    expect(eventData).toEqual({
      ...initialQuery,
      criteria: ratingSetting,
    });
  });

  it('Query changed after change criteria to entity criteria', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const criteriaField = wrapper.find('kpi-rating-criteria-field-stub');

    criteriaField.vm.$emit('input', KPI_RATING_CRITERIA.impactLevel);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.criteria).toEqual(KPI_RATING_CRITERIA.impactLevel);
    expect(eventData).toEqual({
      ...initialQuery,
      metric: ALARM_METRIC_PARAMETERS.totalAlarms,
      criteria: KPI_RATING_CRITERIA.impactLevel,
    });
  });

  it('Query changed after trigger a parameters field', async () => {
    const wrapper = factory({
      propsData: {
        query: initialQuery,
      },
    });

    const metricParametersField = wrapper.find('kpi-rating-metric-field-stub');

    metricParametersField.vm.$emit('input', ALARM_METRIC_PARAMETERS.cancelAckAlarms);

    await flushPromises();

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];

    expect(eventData.metric).toEqual(ALARM_METRIC_PARAMETERS.cancelAckAlarms);
    expect(eventData).toEqual({
      ...initialQuery,
      metric: ALARM_METRIC_PARAMETERS.cancelAckAlarms,
    });
  });

  it('Renders `kpi-rating-filters` with query', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        query: {
          filter: null,
          criteria: undefined,
          metric: ALARM_METRIC_PARAMETERS.ticketAlarms,
          rowsPerPage: 5,
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
