import { mount, createVueInstance } from '@unit/utils/vue';

import CPatternOperatorInformation from '@/components/common/block/c-pattern-operator-information.vue';

const localVue = createVueInstance();

const stubs = {
  'c-pattern-operator-chip': true,
};

const snapshotFactory = (options = {}) => mount(CPatternOperatorInformation, {
  localVue,
  stubs,

  ...options,
});

describe('c-pattern-operator-information', () => {
  test('Renders `c-pattern-operator-information`', () => {
    const wrapper = snapshotFactory({
      slots: {
        default: '<div class="default-slot" />',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
