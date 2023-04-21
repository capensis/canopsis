import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import {
  createAggregatedMetricsModule,
  createAlarmModule,
  createAuthModule, createMetricsModule,
  createMockedStoreModules,
  createServiceModule,
  createUserPreferenceModule,
} from '@unit/utils/store';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { widgetToForm } from '@/helpers/forms/widgets/common';
import { AGGREGATE_FUNCTIONS, ALARM_METRIC_PARAMETERS, SAMPLINGS, WIDGET_TYPES } from '@/constants';

import AlarmChartWidget from '@/components/widgets/chart/alarm-chart-widget.vue';

const stubs = {
  'chart-widget-filters': true,
  'chart-loader': true,
  'bar-chart-metrics': true,
  'line-chart-metrics': true,
  'numbers-metrics': true,
};

describe('alarm-chart-widget', () => {
  mockDateNow(1386435500000);

  const { authModule } = createAuthModule();
  const { alarmModule } = createAlarmModule();
  const { userPreferenceModule } = createUserPreferenceModule();
  const { serviceModule } = createServiceModule();
  const { aggregatedMetricsModule, fetchAggregatedMetricsWithoutStore } = createAggregatedMetricsModule();
  const { metricsModule, fetchAlarmsMetricsWithoutStore } = createMetricsModule();

  const store = createMockedStoreModules([
    authModule,
    userPreferenceModule,
    alarmModule,
    serviceModule,
    aggregatedMetricsModule,
    metricsModule,
  ]);

  const factory = generateShallowRenderer(AlarmChartWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });
  const snapshotFactory = generateRenderer(AlarmChartWidget, {
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
        ],
      },
    });

    factory({
      store,
      propsData: {
        widget: barChartWidget,
        alarm: {
          entity: {
            _id: entityId,
          },
        },
      },
    });

    await flushPromises();

    expect(fetchAlarmsMetricsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          entity: entityId,
          from: 1383843500,
          to: 1386435500,
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
        alarm: {
          entity: {
            _id: entityId,
          },
        },
      },
    });

    await flushPromises();

    expect(fetchAggregatedMetricsWithoutStore).toBeCalledWith(
      expect.any(Object),
      {
        params: {
          entity: entityId,
          from: 1383843500,
          to: 1386435500,
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

  test('Renders `alarm-chart-widget` with bar type', async () => {
    fetchAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [{}, {}],
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
        alarm: {
          entity: {},
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `alarm-chart-widget` with line type', async () => {
    fetchAlarmsMetricsWithoutStore.mockResolvedValue({
      data: [{}, {}, {}],
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
        alarm: {
          entity: {},
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `alarm-chart-widget` with numbers type', async () => {
    fetchAggregatedMetricsWithoutStore.mockResolvedValue({
      data: [{}, {}, {}],
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
        alarm: {
          entity: {},
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
