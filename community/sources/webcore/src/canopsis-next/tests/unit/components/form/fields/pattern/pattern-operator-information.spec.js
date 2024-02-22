import { generateRenderer } from '@unit/utils/vue';

import PatternOperatorInformation from '@/components/forms/fields/pattern/pattern-operator-information.vue';

const stubs = {
  'c-pattern-operator-chip': true,
};

describe('pattern-operator-information', () => {
  const snapshotFactory = generateRenderer(PatternOperatorInformation, { stubs });

  test('Renders `pattern-operator-information`', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
