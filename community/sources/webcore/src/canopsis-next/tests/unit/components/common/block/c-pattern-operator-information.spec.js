import { mount, createVueInstance } from '@unit/utils/vue';

import CPatternOperatorInformation from '@/components/common/block/c-pattern-operator-information.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CPatternOperatorInformation, {
  localVue,

  ...options,
});

describe('c-pattern-operator-information', () => {
  test('Renders `c-pattern-operator-information`', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });
});
