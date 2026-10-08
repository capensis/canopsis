import { mockDateNow } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { ALARM_METRIC_PARAMETERS, QUICK_RANGES, SAMPLINGS } from '@/constants';

import KpiAlarms from '@/components/other/kpi/charts/kpi-alarms.vue';

const metrics = vi.hoisted(() => ({
  createKpiAlarmExport: vi.fn(),
  fetchAlarmsMetricsWithoutStore: vi.fn(),
  fetchMetricExport: vi.fn(),
}));

vi.mock('@/hooks/store/modules/metrics', () => ({
  useMetrics: () => metrics,
}));

const stubs = {
  'c-progress-overlay': true,
  'kpi-alarms-filters': true,
  'kpi-alarms-chart': true,
  'kpi-error-overlay': true,
};

describe('kpi-alarms', () => {
  const nowTimestamp = 1386435600000;

  mockDateNow(nowTimestamp);

  const factory = generateShallowRenderer(KpiAlarms, {
    stubs,
    parentComponent: {
      provide: {
        $system: { timezone: 'Europe/Paris' },
      },
    },
  });

  test('Metrics fetched after mount', async () => {
    const expectedDefaultParams = {
      /* now - 7d  */
      from: 1385852400,
      parameters: [{ metric: ALARM_METRIC_PARAMETERS.createdAlarms }],
      sampling: SAMPLINGS.day,
      filter: null,
      to: 1386370800,
    };
    metrics.fetchAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: { min_date: expectedDefaultParams.from },
    });

    factory();

    await flushPromises();

    expect(metrics.fetchAlarmsMetricsWithoutStore).toHaveBeenCalledTimes(1);
    expect(metrics.fetchAlarmsMetricsWithoutStore).toHaveBeenCalledWith({
      params: expectedDefaultParams,
    });
  });

  test('Metrics refreshed after change interval', async () => {
    const { start, stop } = QUICK_RANGES.last2Days;
    const expectedParamsAfterUpdate = {
      /* now - 2d  */
      from: 1386284400,
      parameters: [{ metric: ALARM_METRIC_PARAMETERS.createdAlarms }],
      sampling: SAMPLINGS.day,
      filter: null,
      to: 1386370800,
    };
    metrics.fetchAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: { min_date: expectedParamsAfterUpdate.from },
    });

    const wrapper = factory();

    await flushPromises();

    const kpiSliFiltersElement = wrapper.find('kpi-alarms-filters-stub');

    kpiSliFiltersElement.triggerCustomEvent('update:query', {
      parameters: [ALARM_METRIC_PARAMETERS.createdAlarms],
      sampling: SAMPLINGS.day,
      filter: null,
      interval: {
        from: start,
        to: stop,
      },
    });

    await flushPromises();

    expect(metrics.fetchAlarmsMetricsWithoutStore).toHaveBeenCalledTimes(2);
    expect(metrics.fetchAlarmsMetricsWithoutStore).toHaveBeenLastCalledWith({
      params: expectedParamsAfterUpdate,
    });
  });

  test('Metrics are not refreshed if min date is later than from', async () => {
    metrics.fetchAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [],
      meta: {
        min_date: 1385930800,
      },
    });

    factory();

    await flushPromises();

    metrics.fetchAlarmsMetricsWithoutStore.mockClear();

    await flushPromises();

    expect(metrics.fetchAlarmsMetricsWithoutStore).not.toHaveBeenCalled();
  });
});
