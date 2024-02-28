import { generateRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';
import { COLORS } from '@/config';

import AvailabilityBarChart from '@/components/other/availability/partials/availability-bar-chart.vue';

const stubs = {
  'availability-bar-chart-information-row': true,
};

describe('availability-bar-chart', () => {
  const snapshotFactory = generateRenderer(AvailabilityBarChart, { stubs });

  test('Renders `availability-bar-chart` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        uptime: 20000,
        downtime: 10000,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-bar-chart` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        uptime: 30000,
        downtime: 20000,
        inactiveTime: 2000,
        uptimeColor: COLORS.primary,
        downtimeColor: COLORS.secondary,
        showType: AVAILABILITY_SHOW_TYPE.duration,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
