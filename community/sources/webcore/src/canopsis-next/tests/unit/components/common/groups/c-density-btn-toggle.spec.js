import { mount, createVueInstance } from '@unit/utils/vue';

import CDensityBtnToggle from '@/components/common/groups/c-density-btn-toggle.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(CDensityBtnToggle, {
  localVue,

  ...options,
});

describe('c-density-btn-toggle', () => {
  it('Renders `c-density-btn-toggle` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-density-btn-toggle` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
