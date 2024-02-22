import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS, SAMPLINGS } from '@/constants';

import KpiAlarmsChart from '@/components/other/kpi/charts/partials/kpi-alarms-chart';

const stubs = {
  'kpi-chart-export-actions': true,
};

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

  const snapshotFactory = generateRenderer(KpiAlarmsChart, {
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

  it('Renders `kpi-alarms-chart` with default props.', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by hour', async () => {
    const firstMetric = metricsInPercentByHour[0];
    const lastMetric = metricsInPercentByHour[metricsInPercentByHour.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioInstructions,
          data: metricsInPercentByHour,
        }],
        sampling: SAMPLINGS.hour,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by day', async () => {
    const firstMetric = metricsInPercentByDay[0];
    const lastMetric = metricsInPercentByDay[metricsInPercentByDay.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
          data: metricsInPercentByDay,
        }],
        sampling: SAMPLINGS.day,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by week', async () => {
    const firstMetric = metricsInPercentByWeek[0];
    const lastMetric = metricsInPercentByWeek[metricsInPercentByWeek.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioTickets,
          data: metricsInPercentByWeek,
        }],
        sampling: SAMPLINGS.week,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with percent by month', async () => {
    const firstMetric = metricsInPercentByMonth[0];
    const lastMetric = metricsInPercentByMonth[metricsInPercentByMonth.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ratioCorrelation,
          data: metricsInPercentByMonth,
        }],
        sampling: SAMPLINGS.month,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by hour', async () => {
    const firstMetric = metricsInTimeByHour[0];
    const lastMetric = metricsInTimeByHour[metricsInTimeByHour.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByHour,
        }],
        sampling: SAMPLINGS.hour,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by day', async () => {
    const firstMetric = metricsInTimeByDay[0];
    const lastMetric = metricsInTimeByDay[metricsInTimeByDay.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByDay,
        }],
        sampling: SAMPLINGS.day,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by week', async () => {
    const firstMetric = metricsInTimeByWeek[0];
    const lastMetric = metricsInTimeByWeek[metricsInTimeByWeek.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByWeek,
        }],
        sampling: SAMPLINGS.week,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with time by month', async () => {
    const firstMetric = metricsInTimeByMonth[0];
    const lastMetric = metricsInTimeByMonth[metricsInTimeByMonth.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.averageAck,
          data: metricsInTimeByMonth,
        }],
        sampling: SAMPLINGS.month,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by hour', async () => {
    const firstMetric = metricsInCountByHour[0];
    const lastMetric = metricsInCountByHour[metricsInCountByHour.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.createdAlarms,
          data: metricsInCountByHour,
        }],
        sampling: SAMPLINGS.hour,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by day', async () => {
    const firstMetric = metricsInCountByDay[0];
    const lastMetric = metricsInCountByDay[metricsInCountByHour.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
          data: metricsInCountByDay,
        }],
        sampling: SAMPLINGS.day,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by week', async () => {
    const firstMetric = metricsInCountByWeek[0];
    const lastMetric = metricsInCountByWeek[metricsInCountByWeek.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.ackAlarms,
          data: metricsInCountByWeek,
        }],
        sampling: SAMPLINGS.week,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with count by month', async () => {
    const firstMetric = metricsInCountByMonth[0];
    const lastMetric = metricsInCountByMonth[metricsInCountByMonth.length - 1];
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [{
          title: ALARM_METRIC_PARAMETERS.correlationAlarms,
          data: metricsInCountByMonth,
        }],
        sampling: SAMPLINGS.month,
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by hour', async () => {
    const firstMetric = metricsInPercentByHour[0];
    const lastMetric = metricsInPercentByHour[metricsInPercentByHour.length - 1];
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
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by day', async () => {
    const firstMetric = metricsInPercentByDay[0];
    const lastMetric = metricsInPercentByDay[metricsInPercentByDay.length - 1];
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
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by week', async () => {
    const firstMetric = metricsInPercentByWeek[0];
    const lastMetric = metricsInPercentByWeek[metricsInPercentByWeek.length - 1];
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
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });

  it('Renders `kpi-alarms-chart` with all types by month', async () => {
    const firstMetric = metricsInPercentByMonth[0];
    const lastMetric = metricsInPercentByMonth[metricsInPercentByMonth.length - 1];
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
        interval: {
          from: firstMetric.timestamp,
          to: lastMetric.timestamp,
        },
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
