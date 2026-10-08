import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, USER_METRIC_PARAMETERS } from '@/constants';

import KpiRating from '@/components/other/kpi/charts/kpi-rating.vue';

const metrics = vi.hoisted(() => ({
  createKpiRatingExport: vi.fn(),
  fetchMetricExport: vi.fn(),
  fetchRatingMetricsWithoutStore: vi.fn(),
}));

vi.mock('@/hooks/store/modules/metrics', () => ({
  useMetrics: () => metrics,
}));

const stubs = {
  'c-progress-overlay': true,
  'kpi-rating-filters': true,
  'kpi-rating-chart': true,
  'kpi-error-overlay': true,
};

describe('kpi-rating', () => {
  const nowTimestamp = 1386435600000;
  const nowUnix = nowTimestamp / 1000;

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiRating, {
    stubs,
    parentComponent: {
      provide: {
        $system: { timezone: 'Europe/Paris' },
      },
    },
  });

  it('Metrics are not fetched after mount without criteria', async () => {
    metrics.fetchRatingMetricsWithoutStore.mockResolvedValue([]);

    factory();

    await flushPromises();

    expect(metrics.fetchRatingMetricsWithoutStore).not.toHaveBeenCalled();
  });

  it('Metrics fetched after setting criteria', async () => {
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      criteria: 1,
      metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      limit: 10,
      filter: undefined,
      to: nowUnix,
    };
    metrics.fetchRatingMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: {
        min_date: expectedDefaultParams.from,
      },
    });

    const wrapper = factory();

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('update:query', {
      criteria: {
        id: expectedDefaultParams.criteria,
      },
      metric: expectedDefaultParams.metric,
      itemsPerPage: expectedDefaultParams.limit,
      interval: {
        from: expectedDefaultParams.from,
        to: expectedDefaultParams.to,
      },
    });

    await flushPromises();

    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalledTimes(1);
    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalledWith({
      params: {
        ...expectedDefaultParams,

        from: 1383778800,
        to: 1386370800,
      },
    });
  });

  it('Metrics are not refreshed if min date is later than from', async () => {
    const expectedDefaultParams = {
      /* now - 30d  */
      from: 1383843600,
      criteria: 1,
      metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      limit: 10,
      filter: undefined,
      to: nowUnix,
    };
    metrics.fetchRatingMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: {
        min_date: 1383943600,
      },
    });

    const wrapper = factory();

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('update:query', {
      criteria: {
        id: expectedDefaultParams.criteria,
      },
      metric: expectedDefaultParams.metric,
      itemsPerPage: expectedDefaultParams.itemsPerPage,
      interval: {
        from: expectedDefaultParams.from,
        to: expectedDefaultParams.to,
      },
    });

    await flushPromises();

    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalled();

    metrics.fetchRatingMetricsWithoutStore.mockClear();

    await flushPromises();

    expect(metrics.fetchRatingMetricsWithoutStore).not.toHaveBeenCalled();
  });

  it('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* today - 1d  */
      from: 1386284400,
      criteria: 1,
      filter: Faker.datatype.string(),
      metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      limit: 10,
      to: 1386370800,
    };
    metrics.fetchRatingMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: {
        min_date: expectedParamsAfterUpdate.from,
      },
    });

    const wrapper = factory();

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('update:query', {
      criteria: {
        id: expectedParamsAfterUpdate.criteria,
      },
      filter: expectedParamsAfterUpdate.filter,
      metric: expectedParamsAfterUpdate.metric,
      itemsPerPage: expectedParamsAfterUpdate.limit,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalledTimes(1);
    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalledWith({
      params: expectedParamsAfterUpdate,
    });
  });

  it('Metrics refreshed without filter with total active time metric', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* today - 1d  */
      from: 1386284400,
      criteria: 1,
      filter: undefined,
      metric: USER_METRIC_PARAMETERS.totalUserActivity,
      limit: 10,
      to: 1386370800,
    };
    metrics.fetchRatingMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: {
        min_date: expectedParamsAfterUpdate.from,
      },
    });

    const wrapper = factory();

    const kpiRatingFiltersElement = wrapper.find('kpi-rating-filters-stub');

    kpiRatingFiltersElement.triggerCustomEvent('update:query', {
      criteria: {
        id: expectedParamsAfterUpdate.criteria,
      },
      metric: expectedParamsAfterUpdate.metric,
      itemsPerPage: expectedParamsAfterUpdate.limit,
      filter: Faker.datatype.string(),
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalledTimes(1);
    expect(metrics.fetchRatingMetricsWithoutStore).toHaveBeenCalledWith({
      params: expectedParamsAfterUpdate,
    });
  });
});
