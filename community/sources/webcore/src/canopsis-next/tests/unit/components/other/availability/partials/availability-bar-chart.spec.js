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
        availability: {
          uptime_duration: 20000,
          uptime_share: 66.67,
          downtime_duration: 10000,
          downtime_share: 33.33,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-bar-chart` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        availability: {
          uptime_duration: 30000,
          uptime_share: 60,
          downtime_duration: 20000,
          downtime_share: 40,
          inactive_time: 2000,
        },
        uptimeColor: COLORS.primary,
        downtimeColor: COLORS.secondary,
        showType: AVAILABILITY_SHOW_TYPE.duration,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
