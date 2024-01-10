import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';
import { ALARM_METRIC_PARAMETERS, KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import PieChartMetrics from '@/components/widgets/chart/partials/pie-chart-metrics.vue';

const stubs = {
  'kpi-chart-export-actions': true,
};

describe('pie-chart-metrics', () => {
  const snapshotFactory = generateRenderer(PieChartMetrics, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `pie-chart-metrics` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pie-chart-metrics` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [
          { title: ALARM_METRIC_PARAMETERS.createdAlarms, value: 12 },
          { title: ALARM_METRIC_PARAMETERS.ackAlarms, value: 12, color: '#010101' },
          { title: ALARM_METRIC_PARAMETERS.ticketActiveAlarms, value: 33 },
        ],
        title: 'Custom title',
        showMode: KPI_PIE_CHART_SHOW_MODS.percent,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  test('Renders `pie-chart-metrics` with number labels props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [
          { title: ALARM_METRIC_PARAMETERS.notAckedAlarms, value: 23 },
          { title: ALARM_METRIC_PARAMETERS.notAckedInDayAlarms, value: 1 },
          { title: ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms, value: 0 },
          { title: ALARM_METRIC_PARAMETERS.notAckedInHourAlarms, value: 22 },
        ],
        showMode: KPI_PIE_CHART_SHOW_MODS.numbers,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  test('Renders `pie-chart-metrics` with empty data', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [
          { title: ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms, value: 0 },
        ],
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(wrapper).toMatchSnapshot();

    expect(canvas.element).toMatchCanvasSnapshot();
  });
});
