import { mockDateNow } from '@unit/utils/mock-hooks';
import {
  createAlarmModule,
  createAuthModule,
  createMockedStoreModules,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { AGGREGATE_FUNCTIONS, ALARM_METRIC_PARAMETERS, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import { widgetToForm } from '@/helpers/entities/widget/form';

import EntityChartWidget from '@/components/widgets/chart/entity-chart-widget.vue';

const metrics = vi.hoisted(() => ({
  createKpiAlarmExport: vi.fn(),
  fetchEntityAggregateMetricsWithoutStore: vi.fn(),
  fetchEntityAlarmsMetricsWithoutStore: vi.fn(),
  fetchMetricExport: vi.fn(),
}));

vi.mock('@/stores/modules/metrics', async () => {
  const { defineStore } = await vi.importActual('pinia');

  return {
    useMetricsStore: defineStore('metrics', {
      actions: metrics,
    }),
  };
});

const stubs = {
  'kpi-widget-filters': true,
  'chart-loader': true,
  'bar-chart-metrics': true,
  'line-chart-metrics': true,
  'numbers-metrics': true,
};

describe('entity-chart-widget', () => {
  mockDateNow(1386435500000);

  const { authModule } = createAuthModule();
  const { alarmModule } = createAlarmModule();
  const { userPreferenceModule } = createUserPreferenceModule();
  const { serviceModule } = createServiceModule();
  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    alarmModule,
    serviceModule,
  ]);

  const factory = generateShallowRenderer(EntityChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: { timezone: 'Europe/Paris' },
      },
    },
  });

  test('Vector metrics fetched after mount', async () => {
    metrics.fetchEntityAlarmsMetricsWithoutStore.mockResolvedValue({ data: [], meta: {} });
    const entityId = Faker.datatype.string();
    const barChartWidget = widgetToForm({
      type: WIDGET_TYPES.barChart,
      parameters: {
        metrics: [
          {
            metric: ALARM_METRIC_PARAMETERS.ackAlarms,
          },
          {
            metric: '.*',
          },
        ],
      },
    });

    factory({
      store,
      propsData: {
        widget: barChartWidget,
        entity: {
          _id: entityId,
        },
        availableMetrics: [ALARM_METRIC_PARAMETERS.ackAlarms],
      },
    });

    await flushPromises();

    expect(metrics.fetchEntityAlarmsMetricsWithoutStore).toHaveBeenCalledWith({
      params: {
        entity: entityId,
        from: 1383865200,
        to: 1386370800,
        with_history: false,
        sampling: SAMPLINGS.day,
        parameters: [
          {
            metric: ALARM_METRIC_PARAMETERS.ackAlarms,
            aggregate_func: '',
          },
        ],
      },
    });
  });

  test('Aggregated metrics fetched after mount', async () => {
    metrics.fetchEntityAggregateMetricsWithoutStore.mockResolvedValue({ data: [], meta: {} });
    const entityId = Faker.datatype.string();
    const numbersWidget = widgetToForm({
      type: WIDGET_TYPES.numbers,
      parameters: {
        show_trend: true,
        metrics: [
          {
            metric: ALARM_METRIC_PARAMETERS.createdAlarms,
            aggregate_func: AGGREGATE_FUNCTIONS.avg,
          },
        ],
      },
    });

    factory({
      store,
      propsData: {
        widget: numbersWidget,
        entity: {
          _id: entityId,
        },
        availableMetrics: [ALARM_METRIC_PARAMETERS.createdAlarms],
      },
    });

    await flushPromises();

    expect(metrics.fetchEntityAggregateMetricsWithoutStore).toHaveBeenCalledWith({
      params: {
        entity: entityId,
        from: 1383865200,
        prev_from: 1383865200,
        to: 1386370800,
        prev_to: 1386370800,
        with_history: true,
        sampling: SAMPLINGS.day,
        parameters: [
          {
            metric: ALARM_METRIC_PARAMETERS.createdAlarms,
            aggregate_func: AGGREGATE_FUNCTIONS.avg,
          },
        ],
      },
    });
  });
});
