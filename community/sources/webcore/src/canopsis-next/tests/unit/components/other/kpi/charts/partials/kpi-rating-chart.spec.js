import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import KpiRatingChart from '@/components/other/kpi/charts/partials/kpi-rating-chart';

const stubs = {
  'kpi-chart-export-actions': true,
};

describe('kpi-rating-chart', () => {
  const metricsInPercent = [
    {
      label: 'Percent label 1',
      value: 70,
    },
    {
      label: 'Percent label 2',
      value: 3,
    },
    {
      label: 'Percent label 2',
      value: 17,
    },
    {
      label: 'Percent label 2',
      value: 18,
    },
  ];
  const metricsInSeconds = [
    {
      label: 'Seconds label 1',
      value: 70054,
    },
    {
      label: 'Seconds label 2',
      value: 70156,
    },
    {
      label: 'Seconds label 2',
      value: 13200,
    },
    {
      label: 'Seconds label 2',
      value: 84254,
    },
  ];
  const metricsInCount = [
    {
      label: 'Counter label 1',
      value: 522,
    },
    {
      label: 'Counter label 1',
      value: 633,
    },
    {
      label: 'Counter label 1',
      value: 156,
    },
    {
      label: 'Counter label 1',
      value: 359,
    },
  ];

  const snapshotFactory = generateRenderer(KpiRatingChart, {
    stubs,
    attachTo: document.body,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Export csv event emitted', async () => {
    const exportCsv = jest.fn();

    const wrapper = snapshotFactory({
      listeners: {
        'export:csv': exportCsv,
      },
    });

    await flushPromises();

    const kpiChartExportActions = wrapper.find('kpi-chart-export-actions-stub');

    kpiChartExportActions.vm.$emit('export:csv');

    expect(exportCsv).toHaveBeenCalledTimes(1);
  });

  it('Export png event emitted', async () => {
    const exportPng = jest.fn();

    const wrapper = snapshotFactory({
      listeners: {
        'export:png': exportPng,
      },
    });

    await flushPromises();

    const kpiChartExportActions = wrapper.find('kpi-chart-export-actions-stub');

    kpiChartExportActions.vm.$emit('export:png');

    expect(exportPng).toHaveBeenCalledTimes(1);
  });

  it('Renders `kpi-rating-chart` with default props.', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-rating-chart` with percent metrics', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsInPercent,
        metric: ALARM_METRIC_PARAMETERS.ratioTickets,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-rating-chart` with seconds metrics', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsInSeconds,
        metric: ALARM_METRIC_PARAMETERS.averageAck,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-rating-chart` with counter metrics', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsInCount,
        metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-rating-chart` with downloading true', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsInCount,
        metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
        downloading: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
