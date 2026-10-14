import { flushPromises, generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_SHOW_TYPE, SAMPLINGS } from '@/constants';

import AvailabilityLineChart from '@/components/other/availability/partials/availability-line-chart.vue';

const stubs = {
  'chart-export-actions': true,
};

const selectChartExportActions = wrapper => wrapper.find('chart-export-actions-stub');

describe('availability-line-chart', () => {
  const availabilities = [
    {
      timestamp: 1709517600,
      uptime_duration: 3142,
      downtime_duration: 7101,
      uptime_share: '30.67',
      downtime_share: 69.33,
    },
    {
      timestamp: 1709521200,
      uptime_duration: 6410,
      downtime_duration: 4355,
      uptime_share: '59.54',
      downtime_share: 40.46,
    },
    {
      timestamp: 1709524800,
      uptime_duration: 1409,
      downtime_duration: 682,
      uptime_share: '67.38',
      downtime_share: 32.620000000000005,
    },
    {
      timestamp: 1709528400,
      uptime_duration: 4072,
      downtime_duration: 6664,
      uptime_share: '37.93',
      downtime_share: 62.07,
    },
    {
      timestamp: 1709532000,
      uptime_duration: 6479,
      downtime_duration: 7925,
      uptime_share: '44.98',
      downtime_share: 55.02,
    },
  ];

  const snapshotFactory = generateRenderer(AvailabilityLineChart, {
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

    selectChartExportActions(wrapper).triggerCustomEvent('export:csv');

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

    selectChartExportActions(wrapper).triggerCustomEvent('export:png');

    expect(exportPng).toHaveBeenCalledTimes(1);
  });

  it('Renders `availability-line-chart` with default props.', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchCanvasSnapshot();
  });

  it('Renders `availability-line-chart` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        chartId: 'custom-chart-id',
        availabilities,
        sampling: SAMPLINGS.hour,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchCanvasSnapshot();
  });

  it.each(Object.entries(AVAILABILITY_SHOW_TYPE))('Renders `availability-line-chart` with "%s" show type props', async (_, showType) => {
    const wrapper = snapshotFactory({
      propsData: {
        availabilities,
        showType,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchCanvasSnapshot();
  });

  it.each(Object.entries(AVAILABILITY_DISPLAY_PARAMETERS))('Renders `availability-line-chart` with "%s" display parameter props', async (_, displayParameter) => {
    const wrapper = snapshotFactory({
      propsData: {
        availabilities,
        displayParameter,
        responsive: false,
        animation: false,
      },
    });

    await flushPromises();

    expect(wrapper).toMatchCanvasSnapshot();
  });
});
