import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { KPI_SLI_GRAPH_DATA_TYPE, SAMPLINGS } from '@/constants';

import KpiSliChart from '@/components/other/kpi/charts/partials/kpi-sli-chart';

const stubs = {
  'kpi-chart-export-actions': true,
};

describe('kpi-sli-chart', () => {
  const metricsInPercent = [
    {
      timestamp: 1633737600,
      uptime: 70,
      downtime: 20,
      maintenance: 10,
    },
    {
      timestamp: 1633824000,
      uptime: 71.5,
      downtime: 19.5,
      maintenance: 9,
    },
    {
      timestamp: 1633910400,
      uptime: 71.5,
      downtime: 19.25,
      maintenance: 9.25,
    },
    {
      timestamp: 1633910400,
      uptime: 71,
      downtime: 29.75,
      maintenance: 9.25,
    },
  ];
  const metricsByHour = [
    {
      timestamp: 1633737600,
      uptime: 2520,
      downtime: 720,
      maintenance: 360,
    },
    {
      timestamp: 1633741200,
      uptime: 2574,
      downtime: 702,
      maintenance: 324,
    },
    {
      timestamp: 1633744800,
      uptime: 1287,
      downtime: 1980,
      maintenance: 333,
    },
    {
      timestamp: 1633748400,
      uptime: 2196,
      downtime: 1071,
      maintenance: 333,
    },
  ];
  const metricsByDay = [
    {
      timestamp: 1633737600,
      uptime: 60480,
      downtime: 17280,
      maintenance: 8640,
    },
    {
      timestamp: 1633824000,
      uptime: 61776,
      downtime: 16848,
      maintenance: 7776,
    },
    {
      timestamp: 1633910400,
      uptime: 61776,
      downtime: 16632,
      maintenance: 7992,
    },
    {
      timestamp: 1633910400,
      uptime: 61344,
      downtime: 16632,
      maintenance: 7992,
    },
  ];
  const metricsByWeek = [
    {
      timestamp: 1633305600,
      uptime: 423360,
      downtime: 120960,
      maintenance: 60480,
    },
    {
      timestamp: 1633910400,
      uptime: 432432,
      downtime: 117936,
      maintenance: 54432,
    },
    {
      timestamp: 1634515200,
      uptime: 429432,
      downtime: 119424,
      maintenance: 55944,
    },
    {
      timestamp: 1635120000,
      uptime: 424408,
      downtime: 121424,
      maintenance: 55944,
    },
  ];
  const metricsByMonth = [
    {
      timestamp: 1631145600,
      uptime: 846720,
      downtime: 362880,
      maintenance: 241920,
    },
    {
      timestamp: 1633737600,
      uptime: 864864,
      downtime: 825552,
      maintenance: 326592,
    },
  ];

  const snapshotFactory = generateRenderer(KpiSliChart, {
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

    kpiChartExportActions.triggerCustomEvent('export:csv');

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

    kpiChartExportActions.triggerCustomEvent('export:png');

    expect(exportPng).toHaveBeenCalledTimes(1);
  });

  it('Renders `kpi-sli-chart` with default props.', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-sli-chart` with percent metrics', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsInPercent,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-sli-chart` with time metrics by hour', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsByHour,
        dataType: KPI_SLI_GRAPH_DATA_TYPE.time,
        sampling: SAMPLINGS.hour,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-sli-chart` with time metrics by day', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsByDay,
        dataType: KPI_SLI_GRAPH_DATA_TYPE.time,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-sli-chart` with time metrics by week', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsByWeek,
        dataType: KPI_SLI_GRAPH_DATA_TYPE.time,
        sampling: SAMPLINGS.week,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-sli-chart` with time metrics by month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsByMonth,
        dataType: KPI_SLI_GRAPH_DATA_TYPE.time,
        sampling: SAMPLINGS.month,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `kpi-sli-chart` with downloading true', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: metricsByMonth,
        dataType: KPI_SLI_GRAPH_DATA_TYPE.time,
        sampling: SAMPLINGS.month,
        downloading: true,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
