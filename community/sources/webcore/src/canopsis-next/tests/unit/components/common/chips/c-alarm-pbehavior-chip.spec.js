import { generateRenderer } from '@unit/utils/vue';

import CAlarmPbehaviorChip from '@/components/common/chips/c-alarm-pbehavior-chip.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
};

describe('c-alarm-pbehavior-chip', () => {
  const icon = 'help';
  const whiteColor = '#fff';
  const blackColor = '#000';
  const snapshotFactory = generateRenderer(CAlarmPbehaviorChip, { stubs });

  test('Renders `c-alarm-pbehavior-chip` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-pbehavior-chip` with white color', () => {
    const wrapper = snapshotFactory({
      propsData: {
        icon,
        color: whiteColor,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-alarm-pbehavior-chip` with black color', () => {
    const wrapper = snapshotFactory({
      propsData: {
        icon,
        color: blackColor,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
