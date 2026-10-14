import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { KPI_SLI_GRAPH_DATA_TYPE, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiSli from '@/components/other/kpi/charts/kpi-sli.vue';

const metrics = vi.hoisted(() => ({
  createKpiSliExport: vi.fn(),
  fetchMetricExport: vi.fn(),
  fetchSliMetricsWithoutStore: vi.fn(),
}));

vi.mock('@/hooks/store/modules/metrics', () => ({
  useMetrics: () => metrics,
}));

const stubs = {
  'c-progress-overlay': true,
  'kpi-sli-filters': true,
  'kpi-sli-chart': true,
  'kpi-error-overlay': true,
};

describe('kpi-sli', () => {
  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiSli, {
    stubs,
    parentComponent: {
      provide: {
        $system: { timezone: 'Europe/Paris' },
      },
    },
  });

  it('Metrics fetched after mount', async () => {
    const expectedDefaultParams = {
      /* now - 7d  */
      from: 1385852400,
      in_percents: true,
      sampling: SAMPLINGS.day,
      to: 1386370800,
      filter: null,
    };
    metrics.fetchSliMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: { min_date: expectedDefaultParams.from },
    });

    factory();

    await flushPromises();

    expect(metrics.fetchSliMetricsWithoutStore).toHaveBeenCalledTimes(1);
    expect(metrics.fetchSliMetricsWithoutStore).toHaveBeenCalledWith({
      params: expectedDefaultParams,
    });
  });

  it('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* now - 2d  */
      from: 1386284400,
      in_percents: true,
      sampling: SAMPLINGS.day,
      to: 1386370800,
      filter: null,
    };
    metrics.fetchSliMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: { min_date: expectedParamsAfterUpdate.from },
    });

    const wrapper = factory();

    await flushPromises();

    const kpiSliFiltersElement = wrapper.find('kpi-sli-filters-stub');

    kpiSliFiltersElement.triggerCustomEvent('update:query', {
      sampling: SAMPLINGS.day,
      filter: null,
      type: KPI_SLI_GRAPH_DATA_TYPE.percent,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(metrics.fetchSliMetricsWithoutStore).toHaveBeenCalledTimes(2);
    expect(metrics.fetchSliMetricsWithoutStore).toHaveBeenLastCalledWith({
      params: expectedParamsAfterUpdate,
    });
  });

  it('Metrics are not refreshed if min date is later than from', async () => {
    metrics.fetchSliMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: { min_date: 1385930800 },
    });

    factory();

    await flushPromises();

    metrics.fetchSliMetricsWithoutStore.mockClear();

    await flushPromises();

    expect(metrics.fetchSliMetricsWithoutStore).not.toHaveBeenCalled();
  });
});
