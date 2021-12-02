import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';
import { ALARM_METRIC_PARAMETERS, SAMPLINGS } from '@/constants';

import KpiAlarmsChart from '@/components/other/kpi/charts/partials/kpi-alarms-chart';

const localVue = createVueInstance();

const stubs = {
  'kpi-chart-export-actions': true,
};

const snapshotFactory = (options = {}) => mount(KpiAlarmsChart, {
  localVue,
  stubs,
  attachTo: document.body,

  ...options,
});

describe('kpi-alarms-chart', () => {
  const metricsInPercentByHour = [
    {
      timestamp: 1633737600,
      value: 63,
    },
    {
      timestamp: 1633741200,
      value: 26,
    },
    {
      timestamp: 1633744800,
      value: 28,
    },
  ];
  const metricsInPercentByDay = [
    {
      timestamp: 1633737600,
      value: 99,
    },
    {
      timestamp: 1633824000,
      value: 11,
    },
    {
      timestamp: 1633910400,
      value: 22,
    },
  ];
  const metricsInPercentByWeek = [
    {
      timestamp: 1633305600,
      value: 54,
    },
    {
      timestamp: 1633910400,
      value: 75,
    },
    {
      timestamp: 1634515200,
      value: 100,
    },
    {
      timestamp: 1635120000,
      value: 64,
    },
  ];
  const metricsInPercentByMonth = [
    {
      timestamp: 1631145600,
      value: 20,
    },
    {
      timestamp: 1633737600,
      value: 30,
    },
  ];

  const metricsInTimeByHour = [
    {
      timestamp: 1633737600,
      value: 7105.62,
    },
    {
      timestamp: 1633741200,
      value: 12300.03,
    },
    {
      timestamp: 1633744800,
      value: 9296.74,
    },
  ];
  const metricsInTimeByDay = [
    {
      timestamp: 1633737600,
      value: 3666,
    },
    {
      timestamp: 1633824000,
      value: 6200,
    },
    {
      timestamp: 1633910400,
      value: 11356,
    },
  ];
  const metricsInTimeByWeek = [
    {
      timestamp: 1633305600,
      value: 66558,
    },
    {
      timestamp: 1633910400,
      value: 95696,
    },
    {
      timestamp: 1634515200,
      value: 96663,
    },
  ];
  const metricsInTimeByMonth = [
    {
      timestamp: 1631145600,
      value: 268412,
    },
    {
      timestamp: 1633737600,
      value: 228412,
    },
  ];

  const metricsInCountByHour = [
    {
      timestamp: 1633737600,
      value: 233,
    },
    {
      timestamp: 1633741200,
      value: 346,
    },
    {
      timestamp: 1633744800,
      value: 1233,
    },
  ];
  const metricsInCountByDay = [
    {
      timestamp: 1633737600,
      value: 2312,
    },
    {
      timestamp: 1633824000,
      value: 5678,
    },
    {
      timestamp: 1633910400,
      value: 9612,
    },
  ];
  const metricsInCountByWeek = [
    {
      timestamp: 1633305600,
      value: 23125,
    },
    {
      timestamp: 1633910400,
      value: 34561,
    },
    {
      timestamp: 1634515200,
      value: 1238,
    },
  ];
  const metricsInCountByMonth = [
    {
      timestamp: 1631145600,
      value: 31627,
    },
    {
      timestamp: 1633737600,
      value: 2312,
    },
  ];

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

  it('Renders `kpi-alarms-chart` with default props.', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by hour', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioInstructions,
          data: metricsInPercentByHour,
        }],
        sampling: SAMPLINGS.hour,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by day', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
          data: metricsInPercentByDay,
        }],
        sampling: SAMPLINGS.day,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by week', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioTickets,
          data: metricsInPercentByWeek,
        }],
        sampling: SAMPLINGS.week,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioCorrelation,
          data: metricsInPercentByMonth,
        }],
        sampling: SAMPLINGS.month,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by hour', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByHour,
        }],
        sampling: SAMPLINGS.hour,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by day', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByDay,
        }],
        sampling: SAMPLINGS.day,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by week', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByWeek,
        }],
        sampling: SAMPLINGS.week,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByMonth,
        }],
        sampling: SAMPLINGS.month,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by hour', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.createdAlarms,
          data: metricsInCountByHour,
        }],
        sampling: SAMPLINGS.hour,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by day', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
          data: metricsInCountByDay,
        }],
        sampling: SAMPLINGS.day,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by week', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ackAlarms,
          data: metricsInCountByWeek,
        }],
        sampling: SAMPLINGS.week,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.correlationAlarms,
          data: metricsInCountByMonth,
        }],
        sampling: SAMPLINGS.month,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by hour', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioInstructions,
          data: metricsInPercentByHour,
        }, {
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByHour,
        }, {
          title: ALARM_METRIC_PARAMETERS.createdAlarms,
          data: metricsInCountByHour,
        }],
        sampling: SAMPLINGS.hour,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by day', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
          data: metricsInPercentByDay,
        }, {
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByDay,
        }, {
          title: ALARM_METRIC_PARAMETERS.cancelAckAlarms,
          data: metricsInCountByDay,
        }],
        sampling: SAMPLINGS.day,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by week', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioCorrelation,
          data: metricsInPercentByWeek,
        }, {
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByWeek,
        }, {
          title: ALARM_METRIC_PARAMETERS.instructionAlarms,
          data: metricsInCountByWeek,
        }],
        sampling: SAMPLINGS.week,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioTickets,
          data: metricsInPercentByMonth,
        }, {
          title: ALARM_METRIC_PARAMETERS.averageResolve,
          data: metricsInTimeByMonth,
        }, {
          title: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
          data: metricsInCountByMonth,
        }],
        sampling: SAMPLINGS.month,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `kpi-alarms-chart` with downloading true', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [],
        sampling: SAMPLINGS.month,
        downloading: true,
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
