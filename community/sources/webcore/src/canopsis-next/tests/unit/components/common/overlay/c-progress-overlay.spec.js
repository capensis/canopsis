import { mount, createVueInstance } from '@unit/utils/vue';

import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';

const localVue = createVueInstance();

describe('c-progress-overlay', () => {
  it('Renders `c-progress-overlay` with default props correctly', () => {
    const wrapper = mount(CProgressOverlay, { localVue });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-progress-overlay` with pending true correctly', () => {
    const wrapper = mount(CProgressOverlay, {
      localVue,
      propsData: {
        pending: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-progress-overlay` with custom props correctly', () => {
    const wrapper = mount(CProgressOverlay, {
      localVue,
      propsData: {
        pending: true,
        opacity: 0.3,
        backgroundColor: 'black',
        color: 'secondary',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
