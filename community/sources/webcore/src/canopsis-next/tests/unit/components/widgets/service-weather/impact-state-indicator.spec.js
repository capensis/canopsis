import { mount, createVueInstance } from '@unit/utils/vue';

import ImpactStateIndicator from '@/components/widgets/service-weather/impact-state-indicator.vue';

const localVue = createVueInstance();

describe('impact-state-indicator', () => {
  it('Renders `impact-state-indicator` with default props correctly', () => {
    const wrapper = mount(ImpactStateIndicator, { localVue });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `impact-state-indicator` with custom props correctly', () => {
    const wrapper = mount(ImpactStateIndicator, {
      localVue,
      propsData: {
        value: 8,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
