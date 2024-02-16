import { generateRenderer } from '@unit/utils/vue';
import { AGGREGATE_FUNCTIONS, ALARM_METRIC_PARAMETERS } from '@/constants';

import NumbersMetricsItem from '@/components/widgets/chart/partials/numbers-metrics-item.vue';

const stubs = {
  'c-help-icon': true,
};

describe('numbers-metrics-item', () => {
  const snapshotFactory = generateRenderer(NumbersMetricsItem, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `numbers-metrics-item` with metric', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.notAckedAlarms,
          value: 12,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics-item` with avg aggregated function', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.timeToAck,
          value: 156,
          aggregate_func: AGGREGATE_FUNCTIONS.avg,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics-item` with sum aggregated function', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.ackAlarms,
          value: 23,
          aggregate_func: AGGREGATE_FUNCTIONS.sum,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics-item` with min aggregated function', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.averageAck,
          value: 35,
          aggregate_func: AGGREGATE_FUNCTIONS.min,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics-item` with max aggregated function', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.averageAck,
          value: 86,
          aggregate_func: AGGREGATE_FUNCTIONS.max,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics-item` with trend up', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.averageAck,
          value: 26,
          aggregate_func: AGGREGATE_FUNCTIONS.max,
          previous_metric: 25,
          previous_interval: {
            from: 1383843600,
            to: 1383844600,
          },
        },
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics-item` with trend down', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metric: {
          title: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
          value: 26,
          aggregate_func: AGGREGATE_FUNCTIONS.avg,
          previous_metric: 35,
          previous_interval: {
            from: 1383843600,
            to: 1383844600,
          },
        },
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
