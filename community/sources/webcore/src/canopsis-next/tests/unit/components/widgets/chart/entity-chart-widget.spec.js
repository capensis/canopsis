import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import {
  createAlarmModule,
  createAuthModule,
  createMetricsModule,
  createMockedStoreModules,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { mockDateNow } from '@unit/utils/mock-hooks';

import { AGGREGATE_FUNCTIONS, ALARM_METRIC_PARAMETERS, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import { widgetToForm } from '@/helpers/entities/widget/form';

import EntityChartWidget from '@/components/widgets/chart/entity-chart-widget.vue';

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
  const { metricsModule,
    fetchEntityAlarmsMetricsWithoutStore,
    fetchEntityAggregateMetricsWithoutStore } = createMetricsModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    alarmModule,
    serviceModule,
    metricsModule,
  ]);

  const factory = generateShallowRenderer(EntityChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(EntityChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Vector metrics fetched after mount', async () => {
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

    expect(fetchEntityAlarmsMetricsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
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
      },
      undefined,
    );
  });

  test('Aggregated metrics fetched after mount', async () => {
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

    expect(fetchEntityAggregateMetricsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          entity: entityId,
          from: 1383865200,
          to: 1386370800,
          with_history: true,
          sampling: SAMPLINGS.day,
          parameters: [
            {
              metric: ALARM_METRIC_PARAMETERS.createdAlarms,
              aggregate_func: AGGREGATE_FUNCTIONS.avg,
            },
          ],
        },
      },
      undefined,
    );
  });

  test('Renders `entity-chart-widget` with bar type', async () => {
    fetchEntityAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [{
        title: ALARM_METRIC_PARAMETERS.ackAlarms,
      }],
    });
    const barChartWidget = widgetToForm({
      type: WIDGET_TYPES.barChart,
      parameters: {
        metrics: [
          {
            metric: ALARM_METRIC_PARAMETERS.ackAlarms,
          },
        ],
      },
    });

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: barChartWidget,
        entity: {},
        availableMetrics: [ALARM_METRIC_PARAMETERS.ackAlarms],
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `entity-chart-widget` with line type', async () => {
    fetchEntityAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [{
        title: ALARM_METRIC_PARAMETERS.ackAlarms,
      }],
    });
    const lineChartWidget = widgetToForm({
      type: WIDGET_TYPES.lineChart,
      parameters: {
        metrics: [
          {
            metric: ALARM_METRIC_PARAMETERS.ackAlarms,
          },
        ],
      },
    });

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: lineChartWidget,
        entity: {},
        availableMetrics: [ALARM_METRIC_PARAMETERS.ackAlarms],
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `entity-chart-widget` with numbers type', async () => {
    fetchEntityAggregateMetricsWithoutStore.mockResolvedValue({
      data: [{
        title: ALARM_METRIC_PARAMETERS.ackAlarms,
      }, {
        title: ALARM_METRIC_PARAMETERS.averageAck,
      }],
    });
    const numbersWidget = widgetToForm({
      type: WIDGET_TYPES.numbers,
      parameters: {
        metrics: [
          {
            metric: ALARM_METRIC_PARAMETERS.ackAlarms,
          },
        ],
      },
    });

    const wrapper = snapshotFactory({
      store,
      propsData: {
        widget: numbersWidget,
        entity: {},
        availableMetrics: [ALARM_METRIC_PARAMETERS.ackAlarms, ALARM_METRIC_PARAMETERS.averageAck],
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
