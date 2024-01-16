import { generateRenderer } from '@unit/utils/vue';

import CPatternOperatorChip from '@/components/common/chips/c-pattern-operator-chip.vue';

describe('c-pattern-operator-chip', () => {
  const snapshotFactory = generateRenderer(CPatternOperatorChip);

  test('Renders `c-pattern-operator-chip`', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
