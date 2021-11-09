import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';

import KpiRatingChart from '@/components/other/kpi/charts/partials/kpi-rating-chart';
import { ALARM_METRIC_PARAMETERS } from '@/constants';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(KpiRatingChart, {
  localVue,
  attachTo: document.body,

  ...options,
});

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
        dataType: ALARM_METRIC_PARAMETERS.ticketAlarms,
      },
    });

    await flushPromises();

    const canvas = wrapper.find('canvas');

    expect(canvas.element).toMatchCanvasSnapshot();
  });
});
