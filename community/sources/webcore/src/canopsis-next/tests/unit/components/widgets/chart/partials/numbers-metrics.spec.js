import { generateRenderer } from '@unit/utils/vue';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

import NumbersMetrics from '@/components/widgets/chart/partials/numbers-metrics.vue';

const stubs = {
  'numbers-metrics-item': true,
};

describe('numbers-metrics', () => {
  const snapshotFactory = generateRenderer(NumbersMetrics, { stubs });

  test('Renders `numbers-metrics` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `numbers-metrics` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        metrics: [
          { title: ALARM_METRIC_PARAMETERS.createdAlarms },
          { title: ALARM_METRIC_PARAMETERS.ackAlarms },
        ],
        title: 'Custom title',
        showTrend: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
