import flushPromises from 'flush-promises';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { stubDateNow } from '@unit/utils/stub-hooks';

import { createMockedStoreModules } from '@unit/utils/store';
import { ALARM_METRIC_PARAMETERS, KPI_RATING_CRITERIA, QUICK_RANGES } from '@/constants';

import KpiRating from '@/components/other/kpi/charts/kpi-rating';

const localVue = createVueInstance();

const stubs = {
  'c-progress-overlay': true,
  'kpi-rating-filters': true,
  'kpi-rating-chart': true,
};

const factory = (options = {}) => shallowMount(KpiRating, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(KpiRating, {
  localVue,
  stubs,

  ...options,
});

describe('kpi-rating', () => {
  const nowTimestamp = 1386435600000;
  const nowUnix = nowTimestamp / 1000;

  stubDateNow(nowTimestamp);

  it('Metrics fetched after mount', async () => {
    const fetchRatingMetrics = jest.fn(() => []);
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      criteria: KPI_RATING_CRITERIA.user,
      metric: ALARM_METRIC_PARAMETERS.ticketAlarms,
      limit: 10,
      to: nowUnix,
    };

    factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    expect(fetchRatingMetrics).toBeCalledTimes(1);
    expect(fetchRatingMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedDefaultParams },
      undefined,
    );
  });

  it('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* now - 30d  */
      from: 1386262800,
      criteria: KPI_RATING_CRITERIA.user,
      metric: ALARM_METRIC_PARAMETERS.ticketAlarms,
      limit: 10,
      to: nowUnix,
    };
    const fetchRatingMetrics = jest.fn(() => []);

    const wrapper = factory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: fetchRatingMetrics,
        },
      }]),
    });

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.vm.$emit('input', {
      filter: null,
      criteria: KPI_RATING_CRITERIA.user,
      metric: ALARM_METRIC_PARAMETERS.ticketAlarms,
      limit: 10,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(fetchRatingMetrics).toBeCalledTimes(2);
    expect(fetchRatingMetrics).toBeCalledWith(
      expect.any(Object),
      { params: expectedParamsAfterUpdate },
      undefined,
    );
  });

  it('Renders `kpi-rating` without metrics', async () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreModules([{
        name: 'metrics',
        actions: {
          fetchRatingMetricsWithoutStore: jest.fn(() => []),
        },
      }]),
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
