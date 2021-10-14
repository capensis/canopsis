import { mount, createVueInstance } from '@unit/utils/vue';

import CEnabled from '@/components/icons/c-enabled.vue';

const localVue = createVueInstance();

describe('c-enabled', () => {
  it('Renders `c-enabled` correctly.', () => {
    const wrapper = mount(CEnabled, { localVue });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-enabled` with enabled prop correctly.', () => {
    const wrapper = mount(CEnabled, {
      localVue,
      propsData: { value: true },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
