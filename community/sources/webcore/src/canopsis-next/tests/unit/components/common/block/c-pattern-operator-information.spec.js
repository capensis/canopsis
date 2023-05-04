import { generateRenderer } from '@unit/utils/vue';

import CPatternOperatorInformation from '@/components/common/block/c-pattern-operator-information.vue';

const stubs = {
  'c-pattern-operator-chip': true,
};

describe('c-pattern-operator-information', () => {
  const snapshotFactory = generateRenderer(CPatternOperatorInformation, { stubs });

  test('Renders `c-pattern-operator-information`', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
