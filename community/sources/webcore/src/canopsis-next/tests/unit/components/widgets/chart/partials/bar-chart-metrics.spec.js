import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS, SAMPLINGS } from '@/constants';

import BarChartMetrics from '@/components/widgets/chart/partials/bar-chart-metrics.vue';

const stubs = {
  'kpi-chart-export-actions': true,
};

describe('bar-chart-metrics', () => {
  const metrics = [
    {
      title: ALARM_METRIC_PARAMETERS.createdAlarms,
      data: [
        { timestamp: 1000000, value: 10 },
        { timestamp: 1100000, value: 15 },
        { timestamp: 1200000, value: 17 },
        { timestamp: 1300000, value: 18 },
      ],
    },
    {
      title: ALARM_METRIC_PARAMETERS.ackAlarms,
      data: [
        { timestamp: 1000000, value: 12 },
        { timestamp: 1100000, value: 16 },
        { timestamp: 1200000, value: 2 },
        { timestamp: 1300000, value: 15 },
      ],
    },
  ];
  const metricsWithHistory = metrics.map(metric => ({
    ...metric,
    data: metric.data.map(item => ({
      ...item,
      history_timestamp: item.timestamp - 100000,
      history_value: item.value - 2,
    })),
  }));

  const snapshotFactory = generateRenderer(BarChartMetrics, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `bar-chart-metrics` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `bar-chart-metrics` with metrics', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics,
        height: 500,
        title: 'Custom title',
        sampling: SAMPLINGS.hour,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper.element).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  test('Renders `bar-chart-metrics` with stacked metrics', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics,
        height: 500,
        title: 'Custom title',
        sampling: SAMPLINGS.hour,
        stacked: true,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper.element).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  test('Renders `bar-chart-metrics` with history data', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsWithHistory,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper.element).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  test('Renders `bar-chart-metrics` with stacked history data', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsWithHistory,
        stacked: true,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper.element).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  test('Renders `bar-chart-metrics` with empty data', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [
          { title: ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms, data: [] },
        ],
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper.element).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });
});
