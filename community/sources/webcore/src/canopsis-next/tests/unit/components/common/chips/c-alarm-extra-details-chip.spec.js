import { generateRenderer } from '@unit/utils/vue';

import CAlarmExtraDetailsChip from '@/components/common/chips/c-alarm-extra-details-chip.vue';

describe('c-alarm-extra-details-chip', () => {
  const icon = 'help';
  const whiteColor = '#fff';
  const blackColor = '#000';
  const snapshotFactory = generateRenderer(CAlarmExtraDetailsChip);

  test('Renders `c-alarm-extra-details-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-extra-details-chip` with white color', () => {
    const wrapper = snapshotFactory({
      propsData: {
        icon,
        color: whiteColor,
        iconColor: blackColor,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-extra-details-chip` with black color', () => {
    const wrapper = snapshotFactory({
      propsData: {
        icon,
        color: blackColor,
        iconColor: whiteColor,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
