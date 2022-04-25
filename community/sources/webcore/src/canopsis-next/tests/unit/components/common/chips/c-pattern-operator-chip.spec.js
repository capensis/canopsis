import { mount, createVueInstance } from '@unit/utils/vue';

import CPatternOperatorChip from '@/components/common/chips/c-pattern-operator-chip.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CPatternOperatorChip, {
  localVue,

  ...options,
});

describe('c-pattern-operator-chip', () => {
  test('Renders `c-pattern-operator-chip`', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
