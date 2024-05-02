import { generateRenderer } from '@unit/utils/vue';

import { COLORS } from '@/config';

import AvailabilityBarChartInformationRow from '@/components/other/availability/partials/availability-bar-chart-information-row.vue';

const stubs = {
  'availability-bar-chart-information-row': true,
};

describe('availability-bar-chart-information-row', () => {
  const snapshotFactory = generateRenderer(AvailabilityBarChartInformationRow, { stubs });

  test('Renders `availability-bar-chart-information-row` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        label: 'Required label',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-bar-chart-information-row` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        label: 'Custom label',
        color: COLORS.primary,
      },
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
